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

func httpPostResultMap(apiurl string, params map[string]interface{}) (map[string]interface{}, error) {
	//data := make(url.Values)
	//for k, v := range params {
	//	data.Add(k, v)
	//}
	byteparams, errparams := json.Marshal(params)
	if errparams != nil {
		return nil, errparams
	}
	resp, errPost := http.NewRequest("post", apiurl, strings.NewReader(string(byteparams)))
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
