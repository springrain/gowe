package gowe

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
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

//httpGetResultMap 发起get请求,并返回json格式的结果
func httpGetResultMap(apiurl string) (map[string]interface{}, error) {
	resp, errGet := client.Get(apiurl)

	if errGet != nil {
		return nil, errGet
	}
	defer resp.Body.Close()

	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}

	resultMap := make(map[string]interface{})
	errJSON := json.Unmarshal(body, &resultMap)
	if errJSON != nil {
		return nil, errJSON
	}

	return resultMap, nil
}

//httpPostResultMap post请求,返回的json封装成map
func httpPostResultMap(apiurl string, params map[string]interface{}) (map[string]interface{}, error) {

	body, errRead := httpPost(apiurl, params)
	if errRead != nil {
		return nil, errRead
	}

	resultMap := make(map[string]interface{})
	errJSON := json.Unmarshal(body, &resultMap)
	if errJSON != nil {
		return nil, errJSON
	}

	return resultMap, nil
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

//httpPostXml 发送Post请求，参数是XML格式的字符串
func httpPostXml(url string, xmlBody string) (body []byte, err error) {
	resp, err := client.Post(url, "application/xml", strings.NewReader(xmlBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

//httpPostXmlWithCert 发送带证书的Post请求，参数是XML格式的字符串
func httpPostXmlWithCert(url string, xmlBody string, client *http.Client) (body []byte, err error) {
	resp, err := client.Post(url, "application/xml", strings.NewReader(xmlBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}
