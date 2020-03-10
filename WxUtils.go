package gowe

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//httpGetMap 发起get请求,并返回json格式的结果
func httpGetMap(url string) (map[string]interface{}, error) {
	resp, errGet := http.Get(url)

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
