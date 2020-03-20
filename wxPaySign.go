package gowe

import (
	"bytes"
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

// 本地通过支付参数计算签名值
// 生成算法:https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_3
func wxPayLocalSign(body map[string]interface{}, signType string, apiKey string) string {
	signStr := wxPaySortSignParams(body, apiKey)
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

// 获取根据Key排序后的请求参数字符串
func wxPaySortSignParams(body map[string]interface{}, apiKey string) string {
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
	buffer.WriteString(fmt.Sprintf("key=%s", apiKey))
	return buffer.String()
}

// 获取沙盒签名Key的返回值
type getSignKeyResponse struct {
	ReturnCode     string `xml:"return_code"` // SUCCESS/FAIL 此字段是通信标识,非交易标识,交易是否成功需要查看result_code来判断
	ReturnMsg      string `xml:"return_msg"`  // 返回信息,如非空,为错误原因:签名失败/参数格式校验错误
	RetMsg         string `xml:"retmsg"`      // 沙盒时返回的错误信息
	Retcode        string `xml:"retcode"`
	MchId          string `xml:"mch_id"`
	SandboxSignkey string `xml:"sandbox_signkey"`
}

// 获取沙盒的签名
func wxPaySandboxSign(wxPayConfig IWxPayConfig, nonceStr string, signType string) (key string, err error) {
	body := make(BodyMap)
	body["mch_id"] = wxPayConfig.GetMchId()
	body["nonce_str"] = nonceStr
	// 计算沙箱参数Sign
	sanboxSign := wxPayLocalSign(body, signType, wxPayConfig.GetAPIKey())
	// 沙箱环境:获取key后,重新计算Sign
	key, err = getSandBoxSignKey(wxPayConfig.GetMchId(), nonceStr, sanboxSign)
	return
}

// 调用微信提供的接口获取SandboxSignkey
func getSandBoxSignKey(mchId string, nonceStr string, sign string) (key string, err error) {
	params := make(map[string]interface{})
	params["mch_id"] = mchId
	params["nonce_str"] = nonceStr
	params["sign"] = sign
	paramXml := GenerateXml(params)
	bytes, err := httpPostXml(WxMpPaySanBoxAPIURL+"/pay/getsignkey", paramXml)
	if err != nil {
		return
	}
	var keyResponse getSignKeyResponse
	if err = xml.Unmarshal(bytes, &keyResponse); err != nil {
		return
	}
	if keyResponse.ReturnCode == ResponseFail {
		err = errors.New(keyResponse.RetMsg)
		return
	}
	key = keyResponse.SandboxSignkey
	return
}

// 验证微信返回的结果签名
func wxPayDoVerifySign(wxPayConfig IWxPayConfig, xmlStr []byte, breakWhenFail bool) (err error) {
	// 生成XML文档
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(xmlStr); err != nil {
		return
	}
	root := doc.SelectElement("xml")
	// 验证return_code
	retCode := root.SelectElement("return_code").Text()
	if retCode != ResponseSuccess && breakWhenFail {
		return
	}
	// 遍历所有Tag,生成Map和Sign
	result, targetSign := make(map[string]interface{}), ""
	for _, elem := range root.ChildElements() {
		// 跳过空值
		if elem.Text() == "" || elem.Text() == "0" {
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
	// 生成签名
	var sign string
	if wxPayConfig.IsProd() {
		sign = wxPayLocalSign(result, signType, wxPayConfig.GetAPIKey())
	} else {
		key, iErr := wxPaySandboxSign(wxPayConfig, result["nonce_str"].(string), SignTypeMD5)
		if err = iErr; iErr != nil {
			return
		}
		sign = wxPayLocalSign(result, SignTypeMD5, key)
	}
	// 验证
	if targetSign != sign {
		err = errors.New("签名无效")
	}
	return
}

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

// 获取根据Key排序后的请求参数字符串
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
