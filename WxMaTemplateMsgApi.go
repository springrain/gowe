package gowe

//WxMaTemplateMsgSend 发送模板消息
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/template-message/templateMessage.send.html

//WxMpTemplateMsg 微信公众号模板消息
type WxMpTemplateMsg struct {
	templateMsgMap map[string]interface{}
	dataMap        map[string]interface{}
}

//SetTouser 接收者(用户)的 openid
func (wxMpTemplateMsg *WxMpTemplateMsg) SetTouser(value string) {
	if wxMpTemplateMsg.templateMsgMap == nil {
		wxMpTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMpTemplateMsg.templateMsgMap["touser"] = value
}

//SetTemplateId 所需下发的模板消息的id
func (wxMpTemplateMsg *WxMpTemplateMsg) SetTemplateId(value string) {
	if wxMpTemplateMsg.templateMsgMap == nil {
		wxMpTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMpTemplateMsg.templateMsgMap["template_id"] = value
}

//SetPage 点击模板卡片后的跳转页面,仅限本小程序内的页面.支持带参数,(示例index?foo=bar).该字段不填则模板无跳转.
func (wxMpTemplateMsg *WxMpTemplateMsg) SetPage(value string) {
	if wxMpTemplateMsg.templateMsgMap == nil {
		wxMpTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMpTemplateMsg.templateMsgMap["page"] = value
}

//SetFormId 表单提交场景下,为 submit 事件带上的 formId；支付场景下,为本次支付的 prepay_id
func (wxMpTemplateMsg *WxMpTemplateMsg) SetFormId(value string) {
	if wxMpTemplateMsg.templateMsgMap == nil {
		wxMpTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMpTemplateMsg.templateMsgMap["form_id"] = value
}

//SetEmphasisKeyword 模板需要放大的关键词,不填则默认无放大
func (wxMpTemplateMsg *WxMpTemplateMsg) SetEmphasisKeyword(value string) {
	if wxMpTemplateMsg.templateMsgMap == nil {
		wxMpTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMpTemplateMsg.templateMsgMap["emphasis_keyword"] = value
}

//AddData 模板内容,不填则下发空模板.具体格式请参考示例.
func (wxMpTemplateMsg *WxMpTemplateMsg) AddData(key string, value string) {
	if wxMpTemplateMsg.dataMap == nil {
		wxMpTemplateMsg.dataMap = make(map[string]interface{})
	}
	valueMap := make(map[string]string)
	valueMap["value"] = value
	wxMpTemplateMsg.dataMap[key] = valueMap
}

//getTemplateMsgMap 获取参数的map,作为请求API的参数
func (wxMpTemplateMsg *WxMpTemplateMsg) getTemplateMsgMap() map[string]interface{} {
	if wxMpTemplateMsg.dataMap != nil {
		wxMpTemplateMsg.templateMsgMap["data"] = wxMpTemplateMsg.dataMap
	}
	return wxMpTemplateMsg.templateMsgMap
}

//WxMaTemplateMsgSend 发送模板消息
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/template-message/templateMessage.send.html
func WxMaTemplateMsgSend(wxMaConfig IWxMaConfig, wxMpTemplateMsg *WxMpTemplateMsg) (*APIResult, error) {
	apiurl := WxMpAPIURL + "/cgi-bin/message/wxopen/template/send?access_token=" + wxMaConfig.GetAccessToken()

	resultMap, errMap := httpPostResultMap(apiurl, wxMpTemplateMsg.getTemplateMsgMap())
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}
