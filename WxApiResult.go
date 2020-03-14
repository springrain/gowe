package gowe

//APIResult  封装 API 响应结果，将 json 字符串转换成 java 数据类型
type APIResult struct {
	//ResultMap 接口返回的json,格式化成map
	ResultMap map[string]interface{}
	//FileData 文件类型[]byte
	FileData []byte
}

//错误的有效期时间
const errExpires = -100

//newAPIResult 使用返回的结果初始化APIResult
func newAPIResult(resultMap map[string]interface{}) APIResult {
	apiResult := APIResult{}
	apiResult.ResultMap = resultMap
	return apiResult
}

//IsSuccess 是否成功
func (apiResult *APIResult) IsSuccess() bool {
	errorCode := apiResult.GetErrorCode()
	// errorCode 为 0 时也可以表示为成功,详见: https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Global_Return_Code.html
	return errorCode == 0
}

//GetErrorCode 获取错误代码
func (apiResult *APIResult) GetErrorCode() int {
	return mapGetInt(apiResult.ResultMap, "errcode")
}

//GetErrorMsg 获取错误代码的说明
func (apiResult *APIResult) GetErrorMsg() string {
	errorCode := mapGetInt(apiResult.ResultMap, "errcode")
	errMsg, ok := errCodeToErrMsgMap[errorCode]
	if ok {
		return errMsg
	}
	return mapGetString(apiResult.ResultMap, "errmsg")
}

//GetOpenID 获取openID
func (apiResult *APIResult) GetOpenID() string {
	return mapGetString(apiResult.ResultMap, "openId")
}

//GetUnionID 获取unionID
func (apiResult *APIResult) GetUnionID() string {
	return mapGetString(apiResult.ResultMap, "unionid")
}

//GetSessionKey 获取sessionKey
func (apiResult *APIResult) GetSessionKey() string {
	return mapGetString(apiResult.ResultMap, "session_key")
}

//GetScope 获取scope
func (apiResult *APIResult) GetScope() string {
	return mapGetString(apiResult.ResultMap, "scope")
}

//GetAccessToken 获取accessToken
func (apiResult *APIResult) GetAccessToken() string {
	return mapGetString(apiResult.ResultMap, "access_token")
}

//isAccessTokenInvalid 返回的accessToken是否有效
func (apiResult *APIResult) isAccessTokenInvalid() bool {
	errorCode := mapGetInt(apiResult.ResultMap, "errcode")
	return errorCode == errExpires || errorCode == 40001 || errorCode == 42001 || errorCode == 42002 || errorCode == 40014
}

//getExpiresIn 获取expiresIn
func (apiResult *APIResult) getExpiresIn() int {
	return mapGetInt(apiResult.ResultMap, "expires_in")
}

//mapGetString 从map中获取一个字符串类型的值
func mapGetString(attrs map[string]interface{}, name string) string {
	value := attrs[name]
	stringValue, intOk := value.(string)
	if !intOk {
		return ""
	}
	return stringValue
}

//mapGetInt 从map中获取一个int类型的值
func mapGetInt(attrs map[string]interface{}, name string) int {
	value := attrs[name]
	intValue, intOk := value.(int)
	if !intOk {
		return errExpires
	}
	return intValue
}
