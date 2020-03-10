package gowe

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

//httpGetResultMap 发起get请求,并返回json格式的结果
func httpGetResultMap(apiurl string) (map[string]interface{}, error) {
	resp, errGet := http.Get(apiurl)

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

func httpPostResultMap(apiurl string, parm map[string]interface{}) (map[string]interface{}, error) {
	//data := make(url.Values)
	//for k, v := range parm {
	//	data.Add(k, v)
	//}
	byteparm, errparm := json.Marshal(parm)
	if errparm != nil {
		return nil, errparm
	}
	resp, errPost := http.NewRequest("post", apiurl, strings.NewReader(string(byteparm)))
	if errPost != nil {
		return nil, errPost
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
