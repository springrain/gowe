package gowe

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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

func httpPostResultMap(apiurl string, parm map[string]string) (map[string]interface{}, error) {
	data := make(url.Values)
	for k, v := range parm {
		data.Add(k, v)
	}
	resp, errPost := http.NewRequest("post", apiurl, strings.NewReader(data.Encode()))
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
