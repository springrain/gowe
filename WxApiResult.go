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
	return mapGetInt(apiResult.Attrs, "errcode")
}

func (apiResult *APIResult) getErrorMsg() string {
	errorCode := mapGetInt(apiResult.Attrs, "errcode")
	errMsg, ok := errCodeToErrMsgMap[errorCode]
	if ok {
		return errMsg
	}
	return mapGetString(apiResult.Attrs, "errmsg")
}

func (apiResult *APIResult) getOpenID() string {
	return mapGetString(apiResult.Attrs, "openId")
}
func (apiResult *APIResult) getUnionID() string {
	return mapGetString(apiResult.Attrs, "unionid")
}
func (apiResult *APIResult) getSessionKey() string {
	return mapGetString(apiResult.Attrs, "session_key")
}
func (apiResult *APIResult) getScope() string {
	return mapGetString(apiResult.Attrs, "scope")
}
func (apiResult *APIResult) getAccessToken() string {
	return mapGetString(apiResult.Attrs, "access_token")
}
func (apiResult *APIResult) isAccessTokenInvalid() bool {
	errorCode := mapGetInt(apiResult.Attrs, "errcode")
	return errorCode == -2 || errorCode == 40001 || errorCode == 42001 || errorCode == 42002 || errorCode == 40014
}
func (apiResult *APIResult) getExpiresIn() int {
	return mapGetInt(apiResult.Attrs, "expires_in")
}

func mapGetString(attrs map[string]interface{}, name string) string {
	value := attrs[name]
	stringValue, intOk := value.(string)
	if !intOk {
		return ""
	}
	return stringValue
}
func mapGetInt(attrs map[string]interface{}, name string) int {
	value := attrs[name]
	intValue, intOk := value.(int)
	if !intOk {
		return -2
	}
	return intValue
}
