package gowe

//APIResult  封装 API 响应结果，将 json 字符串转换成 java 数据类型
type APIResult struct {
	ResultMap map[string]interface{}
	FileData  []byte
}

//错误的有效期时间
const errExpires = -100

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

func (apiResult *APIResult) GetErrorCode() int {
	return mapGetInt(apiResult.ResultMap, "errcode")
}

func (apiResult *APIResult) GetErrorMsg() string {
	errorCode := mapGetInt(apiResult.ResultMap, "errcode")
	errMsg, ok := errCodeToErrMsgMap[errorCode]
	if ok {
		return errMsg
	}
	return mapGetString(apiResult.ResultMap, "errmsg")
}

func (apiResult *APIResult) GetOpenID() string {
	return mapGetString(apiResult.ResultMap, "openId")
}
func (apiResult *APIResult) GetUnionID() string {
	return mapGetString(apiResult.ResultMap, "unionid")
}
func (apiResult *APIResult) GetSessionKey() string {
	return mapGetString(apiResult.ResultMap, "session_key")
}
func (apiResult *APIResult) GetScope() string {
	return mapGetString(apiResult.ResultMap, "scope")
}
func (apiResult *APIResult) GetAccessToken() string {
	return mapGetString(apiResult.ResultMap, "access_token")
}
func (apiResult *APIResult) isAccessTokenInvalid() bool {
	errorCode := mapGetInt(apiResult.ResultMap, "errcode")
	return errorCode == errExpires || errorCode == 40001 || errorCode == 42001 || errorCode == 42002 || errorCode == 40014
}
func (apiResult *APIResult) getExpiresIn() int {
	return mapGetInt(apiResult.ResultMap, "expires_in")
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
		return errExpires
	}
	return intValue
}
