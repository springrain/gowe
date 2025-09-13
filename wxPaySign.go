package gowe

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/beevik/etree"
)

// wxPayLocalSign 本地通过支付参数计算签名值, 生成算法:https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_3
func wxPayLocalSign(body map[string]interface{}, signType string, apiKey string) string {
	signStr := wxPaySortSignParams(body)
	signStr = signStr + "&key=" + apiKey
	var hashSign []byte
	if signType == SignTypeHmacSHA256 {
		hash := hmac.New(sha256.New, []byte(apiKey))
		hash.Write([]byte(signStr))
		hashSign = hash.Sum(nil)
	} else {
		hash := md5.New()
		hash.Write([]byte(signStr))
		hashSign = hash.Sum(nil)
	}
	return strings.ToUpper(hex.EncodeToString(hashSign))
}

// wxPaySortSignParams 获取根据Key排序后的请求参数字符串
func wxPaySortSignParams(body map[string]interface{}) string {
	keyList := make([]string, 0)
	for k := range body {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	buffer := new(bytes.Buffer)
	keyLen := len(keyList)
	for i, k := range keyList {
		s := fmt.Sprintf("%s=%s", k, fmt.Sprintf("%v", body[k]))
		buffer.WriteString(s)
		if (i + 1) < keyLen { //不是最后一个
			buffer.WriteString("&")
		}
	}
	//buffer.WriteString(fmt.Sprintf("key=%s", apiKey))
	return buffer.String()
}

// getSignKeyResponse 获取沙盒签名Key的返回值
type getSignKeyResponse struct {
	ReturnCode     string `xml:"return_code"` // SUCCESS/FAIL 此字段是通信标识,非交易标识,交易是否成功需要查看result_code来判断
	ReturnMsg      string `xml:"return_msg"`  // 返回信息,如非空,为错误原因:签名失败/参数格式校验错误
	RetMsg         string `xml:"retmsg"`      // 沙盒时返回的错误信息
	Retcode        string `xml:"retcode"`
	MchId          string `xml:"mch_id"`
	SandboxSignkey string `xml:"sandbox_signkey"`
}

// wxPaySandboxSign 获取沙盒的签名
func wxPaySandboxSign(ctx context.Context, wxPayConfig IWxPayConfig, nonceStr string, signType string) (key string, err error) {
	body := make(map[string]interface{})
	body["mch_id"] = wxPayConfig.GetMchId(ctx)
	body["nonce_str"] = nonceStr
	// 计算沙箱参数Sign
	sanboxSign := wxPayLocalSign(body, signType, wxPayConfig.GetAPIKey(ctx))
	// 沙箱环境:获取key后,重新计算Sign
	key, err = getSandBoxSignKey(ctx, wxPayConfig.GetMchId(ctx), nonceStr, sanboxSign)
	return
}

// getSandBoxSignKey 调用微信提供的接口获取SandboxSignkey
func getSandBoxSignKey(ctx context.Context, mchId string, nonceStr string, sign string) (key string, err error) {
	params := make(map[string]interface{})
	params["mch_id"] = mchId
	params["nonce_str"] = nonceStr
	params["sign"] = sign
	paramXml := generateXml(params)
	bytes, err := httpPostXml(ctx, WxPaySanBoxAPIURL+"/pay/getsignkey", paramXml)
	if err != nil {
		return
	}
	var keyResponse getSignKeyResponse
	if err = xml.Unmarshal(bytes, &keyResponse); err != nil {
		return
	}
	if keyResponse.ReturnCode == responseFail {
		err = errors.New(keyResponse.RetMsg)
		return
	}
	key = keyResponse.SandboxSignkey
	return
}

// wxPayDoVerifySign 验证微信返回的结果签名
func wxPayDoVerifySign(ctx context.Context, wxPayConfig IWxPayConfig, xmlStr []byte, breakWhenFail bool) error {
	// 生成XML文档
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlStr); err != nil {
		return err
	}
	root := doc.SelectElement("xml")
	// 验证return_code
	retCode := root.SelectElement("return_code").Text()
	if retCode != responseSuccess && breakWhenFail {
		return errors.New(retCode)
	}
	// 遍历所有Tag,生成Map和Sign
	result, targetSign := make(map[string]interface{}), ""
	for _, elem := range root.ChildElements() {
		// 跳过空值
		if elem.Text() == "" {
			continue
		}
		if elem.Tag != "sign" {
			result[elem.Tag] = elem.Text()
		} else {
			targetSign = elem.Text()
		}
	}
	// 获取签名类型
	signType := SignTypeMD5
	if result["sign_type"] != nil {
		signType = result["sign_type"].(string)
	}

	key := wxPayConfig.GetAPIKey(ctx)
	if !wxPayConfig.IsProd(ctx) { //测试的沙箱环境
		var errSandboxSign error
		key, errSandboxSign = wxPaySandboxSign(ctx, wxPayConfig, result["nonce_str"].(string), signType)
		if errSandboxSign != nil {
			return errSandboxSign
		}
	}
	// 生成签名
	sign := wxPayLocalSign(result, signType, key)
	// 验证
	if targetSign != sign {
		return errors.New("签名无效")
	}
	return nil
}

// WxPayJSAPISign 统一下单获取prepay_id参数后,再次计算出JSAPI需要的sign
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=7_7&index=6
func WxPayJSAPISign(appId, nonceStr, packages, signType, timeStamp, apiKey string) (paySign string) {
	signParams := make(map[string]interface{})
	signParams["appId"] = appId
	signParams["nonceStr"] = nonceStr
	signParams["package"] = packages
	signParams["signType"] = signType
	signParams["timeStamp"] = timeStamp
	return wxPayLocalSign(signParams, signType, apiKey)
}

// WxPayMaSign 统一下单获取prepay_id参数后,再次计算出小程序需要的sign
// https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=7_7&index=5
func WxPayMaSign(appId, nonceStr, packages, signType, timeStamp, apiKey string) (paySign string) {
	return WxPayJSAPISign(appId, nonceStr, packages, signType, timeStamp, apiKey)
}

// WxPayAppSign APP支付,统一下单获取支付参数后,再次计算APP支付所需要的的sign
// https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=8_3
func WxPayAppSign(appId, nonceStr, partnerId, prepayId, signType, timeStamp, apiKey string) (paySign string) {
	signParams := make(map[string]interface{}, 0)
	signParams["appId"] = appId
	signParams["nonceStr"] = nonceStr
	signParams["package"] = "Sign=WXPay"
	signParams["partnerid"] = partnerId
	signParams["timeStamp"] = timeStamp
	return wxPayLocalSign(signParams, signType, apiKey)
}

// WxJSSDKTicketSign 生成JS-SDK权限验证的签名
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#62
func WxJSSDKTicketSign(nonceStr, ticket, timeStamp, url string) (ticketSign string) {
	signParams := make(map[string]interface{})
	signParams["noncestr"] = nonceStr
	signParams["jsapi_ticket"] = ticket
	signParams["timestamp"] = timeStamp
	signParams["url"] = url

	// 生成参数排序并拼接
	signStr := wxPaySortSignParams(signParams)
	// 加密签名
	h := sha1.New()
	h.Write([]byte(signStr))
	ticketSign = hex.EncodeToString(h.Sum([]byte("")))
	return
}

/*
//WxPayH5Sign JSAPI支付,统一下单获取支付参数后,再次计算出微信内H5支付需要用的paySign
func WxPayH5Sign(appId, nonceStr, packages, signType, timeStamp, apiKey string) (paySign string) {
	// 原始字符串
	raw := fmt.Sprintf("appId=%s&nonceStr=%s&package=%s&signType=%s&timeStamp=%s&key=%s",
		appId, nonceStr, packages, signType, timeStamp, apiKey)
	buffer := new(bytes.Buffer)
	buffer.WriteString(raw)
	signStr := buffer.String()
	// 加密签名
	var hashSign []byte
	if signType == SignTypeHmacSHA256 {
		hash := hmac.New(sha256.New, []byte(apiKey))
		hash.Write([]byte(signStr))
		hashSign = hash.Sum(nil)
	} else {
		hash := md5.New()
		hash.Write([]byte(signStr))
		hashSign = hash.Sum(nil)
	}
	paySign = strings.ToUpper(hex.EncodeToString(hashSign))
	return
}

//WxPayAppSign APP支付,统一下单获取支付参数后,再次计算APP支付所需要的的sign
func WxPayAppSign(appId, nonceStr, partnerId, prepayId, signType, timeStamp, apiKey string) (paySign string) {
	// 原始字符串
	raw := fmt.Sprintf("appId=%s&nonceStr=%s&package==Sign=WXPay&partnerid=%s&prepayid=%s&timeStamp=%s&key=%s",
		appId, nonceStr, partnerId, prepayId, timeStamp, apiKey)
	buffer := new(bytes.Buffer)
	buffer.WriteString(raw)
	// 加密签名
	signStr := buffer.String()
	var hashSign []byte
	if signType == SignTypeHmacSHA256 {
		hash := hmac.New(sha256.New, []byte(apiKey))
		hash.Write([]byte(signStr))
		hashSign = hash.Sum(nil)
	} else {
		hash := md5.New()
		hash.Write([]byte(signStr))
		hashSign = hash.Sum(nil)
	}
	paySign = strings.ToUpper(hex.EncodeToString(hashSign))
	return
}

//WxPayAPPTicketSign 生成JS-SDK权限验证的签名
func WxPayAPPTicketSign(nonceStr, ticket, timeStamp, url string) (ticketSign string) {
	// 生成参数排序并拼接
	signStr := wxPayAPPsortSignParams(nonceStr, ticket, timeStamp, url)

	// 加密签名
	h := sha1.New()
	h.Write([]byte(signStr))
	ticketSign = hex.EncodeToString(h.Sum([]byte("")))
	return
}

//wxPayAPPsortSignParams 获取根据Key排序后的请求参数字符串
func wxPayAPPsortSignParams(nonceStr, ticket, timeStamp, url string) string {
	body := make(map[string]interface{})
	body["noncestr"] = nonceStr
	body["jsapi_ticket"] = ticket
	body["timestamp"] = timeStamp
	body["url"] = url

	keyList := make([]string, 0)
	for k := range body {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	buffer := new(bytes.Buffer)
	for _, k := range keyList {
		s := fmt.Sprintf("%s=%s&", k, fmt.Sprintf("%v", body[k]))
		buffer.WriteString(s)
	}
	return buffer.String()
}
*/
