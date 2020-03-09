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
func (apiResult *APIResult) IsSuccess() bool {
	errorCode := apiResult.getErrorCode()
	// errorCode 为 0 时也可以表示为成功,详见: https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Global_Return_Code.html
	return errorCode == 0
}

func (apiResult *APIResult) getErrorCode() int {
	return apiResult.getInt("errcode")
}

func (apiResult *APIResult) getErrorMsg() string {
	errorCode := apiResult.getInt("errcode")
	errMsg, ok := errCodeToErrMsgMap[errorCode]
	if ok {
		return errMsg
	}
	return apiResult.get("errmsg")
}

func (apiResult *APIResult) getOpenID() string {
	return apiResult.get("openId")
}
func (apiResult *APIResult) getUnionID() string {
	return apiResult.get("unionid")
}
func (apiResult *APIResult) getSessionKey() string {
	return apiResult.get("session_key")
}
func (apiResult *APIResult) getScope() string {
	return apiResult.get("scope")
}
func (apiResult *APIResult) getAccessToken() string {
	return apiResult.get("access_token")
}
func (apiResult *APIResult) isAccessTokenInvalid() bool {
	errorCode := apiResult.getInt("errcode")
	return errorCode == -2 || errorCode == 40001 || errorCode == 42001 || errorCode == 42002 || errorCode == 40014
}
func (apiResult *APIResult) getExpiresIn() int {
	return apiResult.getInt("expires_in")
}

func (apiResult *APIResult) get(name string) string {
	value := apiResult.Attrs[name]
	stringValue, intOk := value.(string)
	if !intOk {
		return ""
	}
	return stringValue
}
func (apiResult *APIResult) getInt(name string) int {
	value := apiResult.Attrs[name]
	intValue, intOk := value.(int)
	if !intOk {
		return -2
	}
	return intValue
}
