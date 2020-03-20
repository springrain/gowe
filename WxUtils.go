package gowe

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"
)

//http请求的client
var client *http.Client

//初始化 http连接信息
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

//httpGet 发起get请求
func httpGet(apiurl string) ([]byte, error) {
	resp, errGet := client.Get(apiurl)

	if errGet != nil {
		return nil, errGet
	}
	defer resp.Body.Close()

	body, errRead := ioutil.ReadAll(resp.Body)
	return body, errRead
}

//httpPost post请求,返回原始的字节数组
func httpPost(apiurl string, params map[string]interface{}) ([]byte, error) {
	//data := make(url.Values)
	//for k, v := range params {
	//	data.Add(k, v)
	//}
	byteparams, errparams := json.Marshal(params)
	if errparams != nil {
		return nil, errparams
	}
	resp, errPost := client.Post(apiurl, "application/json", bytes.NewReader(byteparams))
	if errPost != nil {
		return nil, errPost
	}
	defer resp.Body.Close()

	body, errRead := ioutil.ReadAll(resp.Body)

	return body, errRead
}

//httpPostXml 发送Post请求,参数是XML格式的字符串
func httpPostXml(url string, xmlBody string) (body []byte, err error) {
	resp, err := client.Post(url, "application/xml", strings.NewReader(xmlBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

//httpPostXmlWithCert 发送带证书的Post请求,参数是XML格式的字符串
func httpPostXmlWithCert(url string, xmlBody string, client *http.Client) (body []byte, err error) {
	resp, err := client.Post(url, "application/xml", strings.NewReader(xmlBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func wxPayBuildBody(wxPayConfig IWxPayConfig, bodyObj interface{}) (body map[string]interface{}, err error) {
	// 将bodyObj转换为map[string]interface{}类型
	bodyJson, _ := json.Marshal(bodyObj)
	body = make(map[string]interface{})
	_ = json.Unmarshal(bodyJson, &body)
	// 添加固定参数
	if wxPayConfig.IsMch() {
		body["mch_appid"] = wxPayConfig.GetAppId()
		body["mchid"] = wxPayConfig.GetMchId()
	} else {
		body["appid"] = wxPayConfig.GetAppId()
		body["mch_id"] = wxPayConfig.GetMchId()
	}
	if isWxPayFacilitator(wxPayConfig.GetServiceType()) {
		body["sub_appid"] = wxPayConfig.GetSubAppId()
		body["sub_mch_id"] = wxPayConfig.GetSubMchId()
	}
	nonceStr := GetRandomString(32)
	body["nonce_str"] = nonceStr
	// 生成签名
	signType, _ := body["sign_type"].(string)
	var sign string
	if wxPayConfig.IsProd() {
		sign = wxPayLocalSign(body, signType, wxPayConfig.GetAPIKey())
	} else {
		body["sign_type"] = signTypeMD5
		key, iErr := wxPaySandboxSign(wxPayConfig, nonceStr, signTypeMD5)
		if err = iErr; iErr != nil {
			return
		}
		sign = wxPayLocalSign(body, signTypeMD5, key)
	}
	body["sign"] = sign
	return
}

// 是否是服务商模式
func isWxPayFacilitator(serviceType int) bool {
	switch serviceType {
	case serviceTypeFacilitatorDomestic, serviceTypeFacilitatorAbroad, serviceTypeBankServiceProvidor:
		return true
	default:
		return false
	}
}

//GenerateXml 生成请求XML的Body体
func GenerateXml(data map[string]interface{}) string {
	buffer := new(bytes.Buffer)
	buffer.WriteString("<xml>")
	for k, v := range data {
		buffer.WriteString(fmt.Sprintf("<%s><![CDATA[%v]]></%s>", k, v, k))
	}
	buffer.WriteString("</xml>")
	return buffer.String()
}

//JsonString 生成模型字符串
func JsonString(m interface{}) string {
	bs, _ := json.Marshal(m)
	return string(bs)
}

//FormatDateTime 格式化时间,按照yyyyMMddHHmmss格式
func FormatDateTime(t time.Time) string {
	return t.Format("20060102150405")
}

//EscapedPath 对URL进行Encode编码
func EscapedPath(u string) (path string, err error) {
	uriObj, err := url.Parse(u)
	if err != nil {
		return
	}
	path = uriObj.EscapedPath()
	return
}

//PKCS7UnPadding 解密填充模式(去除补全码) PKCS7UnPadding 解密时,需要在最后面去掉加密时添加的填充byte
func PKCS7UnPadding(plainText []byte) []byte {
	length := len(plainText)
	unpadding := int(plainText[length-1])   // 找到Byte数组最后的填充byte
	return plainText[:(length - unpadding)] // 只截取返回有效数字内的byte数组
}

//IsValidAuthCode 18位纯数字,以10、11、12、13、14、15开头
func IsValidAuthCode(authcode string) (ok bool) {
	pattern := "^1[0-5][0-9]{16}$"
	ok, _ = regexp.MatchString(pattern, authcode)
	return
}

//GetRandomString 获取随机字符串
func GetRandomString(length int) string {
	str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	b := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, b[r.Intn(len(b))])
	}
	return string(result)
}

//wxPayDoWeChat 向微信发送请求
func wxPayDoWeChat(wxPayConfig IWxPayConfig, apiurl string, bodyObj interface{}) (bytes []byte, err error) {
	// 转换参数
	body, err := wxPayBuildBody(wxPayConfig, bodyObj)
	if err != nil {
		return
	}
	// 发起请求
	bytes, err = httpPostXml(apiurl, GenerateXml(body))
	return
}

//wxPayDoWeChatWithCert 向微信发送带证书请求
func wxPayDoWeChatWithCert(wxPayConfig IWxPayConfig, apiurl string, bodyObj interface{}) ([]byte, error) {
	// 转换参数
	body, err := wxPayBuildBody(wxPayConfig, bodyObj)
	if err != nil {
		return nil, err
	}
	// 设置证书和连接池
	client, err := wxPayGetCertHttpClient(wxPayConfig)
	if err != nil {
		return nil, err
	}
	// 发起请求
	bytes, err := httpPostXmlWithCert(apiurl, GenerateXml(body), client)
	return bytes, err
}

//wxPayGetCertHttpClient 获取带证数的httpClient
func wxPayGetCertHttpClient(wxPayConfig IWxPayConfig) (*http.Client, error) {
	certPath := wxPayConfig.GetCertificateFile()
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	client, err := wxPayBuildClient(certData, wxPayConfig.GetMchId())

	return client, err
}

//wxPayBuildClient 构建带证数的httpClient
func wxPayBuildClient(data []byte, mchId string) (client *http.Client, err error) {
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

//wxPayPkc12ToPerm 证数格式转化
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
