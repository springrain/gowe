package gowe

import (
	"bytes"
	"context"
	"crypto"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"
)

// http请求的client
var client *http.Client

// 初始化 http连接信息
func init() {
	client = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:     3 * time.Minute,
			TLSHandshakeTimeout: 10 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 10 * time.Minute,
				DualStack: true,
			}).DialContext,
		},
	}
}

// httpGet 发起get请求
func httpGet(ctx context.Context, apiurl string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", apiurl, nil)
	if err != nil {
		return nil, err
	}
	resp, errGet := client.Do(req)
	//resp, errGet := client.Get(apiurl)

	if errGet != nil {
		return nil, errGet
	}
	defer resp.Body.Close()

	body, errRead := io.ReadAll(resp.Body)
	return body, errRead
}

// httpPost post请求,返回原始的字节数组
func httpPost(ctx context.Context, apiurl string, params map[string]interface{}) ([]byte, error) {
	//data := make(url.Values)
	//for k, v := range params {
	//	data.Add(k, v)
	//}
	byteparams, errparams := json.Marshal(params)
	if errparams != nil {
		return nil, errparams
	}
	req, err := http.NewRequestWithContext(ctx, "POST", apiurl, bytes.NewReader(byteparams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, errPost := client.Do(req)
	//resp, errPost := client.Post(apiurl, "application/json", bytes.NewReader(byteparams))
	if errPost != nil {
		return nil, errPost
	}
	defer resp.Body.Close()

	body, errRead := io.ReadAll(resp.Body)

	return body, errRead
}

// httpPostXml 发送Post请求,参数是XML格式的字符串
func httpPostXml(ctx context.Context, url string, xmlBody string) (body []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(xmlBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/xml")
	resp, err := client.Do(req)
	//resp, err := client.Post(url, "application/xml", strings.NewReader(xmlBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	return
}

// httpPostXmlWithCert 发送带证书的Post请求,参数是XML格式的字符串
func httpPostXmlWithCert(ctx context.Context, url string, xmlBody string, client *http.Client) (body []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(xmlBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/xml")
	resp, err := client.Do(req)
	//resp, err := client.Post(url, "application/xml", strings.NewReader(xmlBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	return
}

// mchType 0:普通商户接口, 1:特殊的商户接口(企业付款,微信找零),2:红包
func wxPayBuildBody(ctx context.Context, wxPayConfig IWxPayConfig, bodyObj interface{}, mchType int) (body map[string]interface{}, err error) {
	// 将bodyObj转换为map[string]interface{}类型
	bodyJson, _ := json.Marshal(bodyObj)
	body = make(map[string]interface{})
	_ = json.Unmarshal(bodyJson, &body)
	// 添加固定参数
	if mchType == 1 { //特殊的商户接口(企业付款,微信找零)
		body["mch_appid"] = wxPayConfig.GetAppId(ctx)
		body["mchid"] = wxPayConfig.GetMchID(ctx)
	} else if mchType == 2 { //红包
		body["wxappid"] = wxPayConfig.GetAppId(ctx) //微信分配的公众账号ID(企业号corpid即为此appId).接口传入的所有appid应该为公众号的appid(在mp.weixin.qq.com申请的),不能为APP的appid(在open.weixin.qq.com申请的)
		body["mch_id"] = wxPayConfig.GetMchID(ctx)
	} else { //普通微信支付
		body["appid"] = wxPayConfig.GetAppId(ctx)
		body["mch_id"] = wxPayConfig.GetMchID(ctx)
	}

	//如果是服务商模式
	if isWxPayFacilitator(wxPayConfig.GetServiceType(ctx)) {
		body["sub_appid"] = wxPayConfig.GetSubAppId(ctx)
		body["sub_mch_id"] = wxPayConfig.GetSubMchId(ctx)
	}
	//nonceStr := getRandomString(32)
	nonceStr := FuncGenerateRandomString(ctx, 32)
	body["nonce_str"] = nonceStr
	// 生成签名
	signType, _ := body["sign_type"].(string)
	var sign string
	if wxPayConfig.IsProd(ctx) {
		sign = wxPayLocalSign(body, signType, wxPayConfig.GetAPIKey(ctx))
	} else {
		body["sign_type"] = SignTypeMD5
		key, iErr := wxPaySandboxSign(ctx, wxPayConfig, nonceStr, SignTypeMD5)
		if err = iErr; iErr != nil {
			return
		}
		sign = wxPayLocalSign(body, SignTypeMD5, key)
	}
	body["sign"] = sign
	return
}

// 是否是服务商模式
func isWxPayFacilitator(serviceType int) bool {
	switch serviceType {
	case ServiceTypeFacilitatorDomestic, ServiceTypeFacilitatorAbroad, ServiceTypeBankServiceProvidor:
		return true
	default:
		return false
	}
}

// generateXml 生成请求XML的Body体
func generateXml(data map[string]interface{}) string {
	buffer := new(bytes.Buffer)
	buffer.WriteString("<xml>")
	for k, v := range data {
		buffer.WriteString(fmt.Sprintf("<%s><![CDATA[%v]]></%s>", k, v, k))
	}
	buffer.WriteString("</xml>")
	return buffer.String()
}

// jsonString 生成模型字符串
func jsonString(m interface{}) string {
	bs, _ := json.Marshal(m)
	return string(bs)
}

// formatDateTime 格式化时间,按照yyyyMMddHHmmss格式
func formatDateTime(t time.Time) string {
	return t.Format("20060102150405")
}

// encodePath 对URL进行Encode编码
func encodePath(u string) (path string, err error) {
	uriObj, err := url.Parse(u)
	if err != nil {
		return
	}
	path = uriObj.EscapedPath()
	return
}

// pkcs7UnPadding 解密填充模式(去除补全码) pkcs7UnPadding 解密时,需要在最后面去掉加密时添加的填充byte
func pkcs7UnPadding(plainText []byte) []byte {
	length := len(plainText)
	unpadding := int(plainText[length-1])   // 找到Byte数组最后的填充byte
	return plainText[:(length - unpadding)] // 只截取返回有效数字内的byte数组
}

// isValidAuthCode 18位纯数字,以10、11、12、13、14、15开头
func isValidAuthCode(authcode string) (ok bool) {
	pattern := "^1[0-5][0-9]{16}$"
	ok, _ = regexp.MatchString(pattern, authcode)
	return
}

// FuncGenerateRandomString 生成指定位数的随机字符串
var FuncGenerateRandomString func(context.Context, int) string = generateRandomString

// generateRandomString 获取随机字符串
func generateRandomString(ctx context.Context, length int) string {
	str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	b := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, b[r.Intn(len(b))])
	}
	return strings.ToLower(string(result))
}

// wxPayDoWeChat 向微信发送请求
func wxPayDoWeChat(ctx context.Context, wxPayConfig IWxPayConfig, apiuri string, bodyObj interface{}, mchType int) (bytes []byte, err error) {
	apiurl := WxPayMchAPIURL + apiuri
	if !wxPayConfig.IsProd(ctx) {
		apiurl = WxPaySanBoxAPIURL + apiuri
	}
	// 转换参数
	body, err := wxPayBuildBody(ctx, wxPayConfig, bodyObj, mchType)
	if err != nil {
		return
	}
	// 发起请求
	bytes, err = httpPostXml(ctx, apiurl, generateXml(body))
	return
}

// wxPayDoWeChatWithCert 向微信发送带证书请求
// mchType 0:普通商户接口, 1:特殊的商户接口(企业付款,微信找零),2:红包
func wxPayDoWeChatWithCert(ctx context.Context, wxPayConfig IWxPayConfig, apiuri string, bodyObj interface{}, mchType int) ([]byte, error) {
	// 转换参数
	body, err := wxPayBuildBody(ctx, wxPayConfig, bodyObj, mchType)
	if err != nil {
		return nil, err
	}
	// 设置证书和连接池
	client, err := wxPayGetCertHttpClient(ctx, wxPayConfig)
	if err != nil {
		return nil, err
	}

	apiurl := WxPayMchAPIURL + apiuri
	if !wxPayConfig.IsProd(ctx) {
		apiurl = WxPaySanBoxAPIURL + apiuri
	}

	// 发起请求
	bytes, err := httpPostXmlWithCert(ctx, apiurl, generateXml(body), client)
	return bytes, err
}

// wxPayGetCertHttpClient 获取带证数的httpClient
func wxPayGetCertHttpClient(ctx context.Context, wxPayConfig IWxPayConfig) (*http.Client, error) {
	certPath := wxPayConfig.GetCertificateFile(ctx)
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	client, err := wxPayBuildClient(ctx, certData, wxPayConfig.GetMchID(ctx))

	return client, err
}

// wxPayBuildClient 构建带证数的httpClient
func wxPayBuildClient(ctx context.Context, data []byte, mchId string) (client *http.Client, err error) {
	// 将pkcs12证书转成pem
	cert, err := wxPayPkc12ToPerm(data, mchId)
	if err != nil {
		return
	}
	// tls配置
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	// 带证书的客户端
	client = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:     3 * time.Minute,
			TLSHandshakeTimeout: 10 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 10 * time.Minute,
				DualStack: true,
			}).DialContext,
			TLSClientConfig:    config,
			DisableCompression: true,
		},
	}
	return
}

// wxPayPkc12ToPerm 证数格式转化
func wxPayPkc12ToPerm(data []byte, mchId string) (cert tls.Certificate, err error) {
	blocks, err := pkcs12.ToPEM(data, mchId)
	if err != nil {
		return
	}
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}
	cert, err = tls.X509KeyPair(pemData, pemData)
	return
}

// 以上是v2
// --------------------------------------v3微信支付util--------------------------------------------------------
// ------------------创建JSAPI支付订单---------------------
func createJsapiOrder(ctx context.Context, wxPayConfig IWxPayConfig, openid, outTradeNo string, totalFee int, description string) (string, error) {
	appId := wxPayConfig.GetAppId(ctx)
	mchId := wxPayConfig.GetMchID(ctx)
	notifyUrl := wxPayConfig.GetNotifyURL(ctx)
	certSerialNo := wxPayConfig.GetCertSerialNo(ctx)
	reqData := JsapiOrderRequest{
		AppID:       appId,
		MchID:       mchId,
		Description: description,
		OutTradeNo:  outTradeNo,
		NotifyURL:   notifyUrl,
		TimeExpire:  time.Now().Add(30 * time.Minute).Format("2006-01-02T15:04:05Z07:00"),
	}
	reqData.Amount.Total = totalFee
	reqData.Amount.Currency = "CNY"
	reqData.Payer.OpenID = openid

	body, _ := json.Marshal(reqData)

	req, _ := http.NewRequest("POST", "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi", bytes.NewBuffer(body))

	timestamp := time.Now().Unix()
	nonceStr := generateNonceStr()
	signature, err := generateSignature(ctx, wxPayConfig, "POST", "/v3/pay/transactions/jsapi", timestamp, nonceStr, body)
	if err != nil {
		return "", fmt.Errorf("生成签名失败: %v", err)
	}

	authHeader := fmt.Sprintf(
		"WECHATPAY2-SHA256-RSA2048 mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%d\",serial_no=\"%s\",signature=\"%s\"",
		mchId, nonceStr, timestamp, certSerialNo, signature,
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("User-Agent", "go-wechatpay-client")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求微信API失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return "", fmt.Errorf("响应解析失败: %s", string(respBody))
		}
		return "", fmt.Errorf("微信API错误: %s(%s)", errResp.Message, errResp.Code)
	}

	var orderResp JsapiOrderResponse
	if err := json.Unmarshal(respBody, &orderResp); err != nil {
		return "", fmt.Errorf("响应解析失败: %v", err)
	}

	return orderResp.PrepayID, nil
}

// 生成商户订单号
func generateOutTradeNo() string {
	return fmt.Sprintf("ORDER%d", time.Now().UnixNano())
}

func generateJsapiPayParams(ctx context.Context, wxPayConfig IWxPayConfig, prepayID string) (*JsapiPayParams, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonceStr := generateNonceStr()
	packageStr := fmt.Sprintf("prepay_id=%s", prepayID)
	appId := wxPayConfig.GetAppId(ctx)
	signStr := fmt.Sprintf("%s\n%s\n%s\n%s\n", appId, timestamp, nonceStr, packageStr)

	signature, err := signWithPrivateKey(ctx, wxPayConfig, signStr)
	if err != nil {
		return nil, err
	}

	return &JsapiPayParams{
		AppID:     appId,
		Timestamp: timestamp,
		NonceStr:  nonceStr,
		Package:   packageStr,
		SignType:  "RSA",
		PaySign:   signature,
	}, nil
}

// 生成随机字符串
func generateNonceStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 生成V3 API签名
func generateSignature(ctx context.Context, wxPayConfig IWxPayConfig, method, path string, timestamp int64, nonceStr string, body []byte) (string, error) {
	signStr := fmt.Sprintf("%s\n%s\n%d\n%s\n%s\n",
		method, path, timestamp, nonceStr, string(body),
	)
	return signWithPrivateKey(ctx, wxPayConfig, signStr)
}

// 使用私钥签名
func signWithPrivateKey(ctx context.Context, wxPayConfig IWxPayConfig, data string) (string, error) {
	// 1. 读取私钥文件
	apiPrivateKey := wxPayConfig.GetPrivateKey(ctx)
	keyBytes, err := os.ReadFile(apiPrivateKey)
	if err != nil {
		return "", fmt.Errorf("读取私钥文件失败: %v", err)
	}
	//2 解析 PEM
	block, _ := pem.Decode([]byte(keyBytes))
	if block == nil {
		return "", fmt.Errorf("私钥解析失败")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("私钥解析失败: %v", err)
	}
	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("不是有效的RSA私钥")
	}

	hashed := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(crand.Reader, rsaKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("签名失败: %v", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// -------------------- 创建Native支付订单--------------------------------------
func createNativeOrder(ctx context.Context, wxPayConfig IWxPayConfig, ip string, storeId string, outTradeNo string, totalFee int, description string) (string, error) {
	appId := wxPayConfig.GetAppId(ctx)
	mchId := wxPayConfig.GetMchID(ctx)
	notifyUrl := wxPayConfig.GetNotifyURL(ctx)
	certSerialNo := wxPayConfig.GetCertSerialNo(ctx)
	reqData := NativeOrderRequest{
		AppID:       appId,
		MchID:       mchId,
		Description: description,
		OutTradeNo:  outTradeNo,
		NotifyURL:   notifyUrl,
		TimeExpire:  time.Now().Add(30 * time.Minute).Format("2006-01-02T15:04:05Z07:00"),
	}
	reqData.Amount.Total = totalFee
	reqData.Amount.Currency = "CNY"
	reqData.SceneInfo.PayerClientIp = ip
	reqData.SceneInfo.StoreInfo.ID = storeId

	body, _ := json.Marshal(reqData)

	req, _ := http.NewRequest("POST", "https://api.mch.weixin.qq.com/v3/pay/transactions/native", bytes.NewBuffer(body))

	timestamp := time.Now().Unix()
	nonceStr := generateNonceStr()
	signature, err := generateSignature(ctx, wxPayConfig, "POST", "/v3/pay/transactions/native", timestamp, nonceStr, body)
	if err != nil {
		return "", fmt.Errorf("生成签名失败: %v", err)
	}

	authHeader := fmt.Sprintf(
		"WECHATPAY2-SHA256-RSA2048 mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%d\",serial_no=\"%s\",signature=\"%s\"",
		mchId, nonceStr, timestamp, certSerialNo, signature,
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("User-Agent", "go-wechatpay-client")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求微信API失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return "", fmt.Errorf("响应解析失败: %s", string(respBody))
		}
		return "", fmt.Errorf("微信API错误: %s(%s)", errResp.Message, errResp.Code)
	}

	var orderResp NativeOrderResponse
	if err := json.Unmarshal(respBody, &orderResp); err != nil {
		return "", fmt.Errorf("响应解析失败: %v", err)
	}

	return orderResp.CodeUrl, nil
}
