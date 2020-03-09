package gowe

import (
	"encoding/json"
	"os"
)

//APIResult  封装 API 响应结果，将 json 字符串转换成 java 数据类型
type APIResult struct {
	Attrs   map[string]interface{}
	jsonStr string
	file    *os.File
}

//newAPIResult 创建APIResult
func newAPIResult(jsonstr string) (*APIResult, error) {
	apiResult := APIResult{}
	apiResult.Attrs = make(map[string]interface{}, 0)
	apiResult.jsonStr = jsonstr
	if err := json.Unmarshal([]byte(jsonstr), &(apiResult.Attrs)); err != nil {
		return nil, err
	}
	return &apiResult, nil
}

//IsSuccess 是否成功
func (apiResult *APIResult) IsSuccess() {

}
