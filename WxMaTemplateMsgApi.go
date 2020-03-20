package gowe

import "encoding/json"

//WxMaTemplateMsgSend 发送模板消息
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/template-message/templateMessage.send.html
//WxMaTemplateMsg 微信公众号模板消息
type WxMaTemplateMsg struct {
	templateMsgMap map[string]interface{}
	dataMap        map[string]interface{}
}

//SetTouser 接收者(用户)的 openid
func (wxMaTemplateMsg *WxMaTemplateMsg) SetTouser(value string) {
	if wxMaTemplateMsg.templateMsgMap == nil {
		wxMaTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMaTemplateMsg.templateMsgMap["touser"] = value
}

//SetTemplateId 所需下发的模板消息的id
func (wxMaTemplateMsg *WxMaTemplateMsg) SetTemplateId(value string) {
	if wxMaTemplateMsg.templateMsgMap == nil {
		wxMaTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMaTemplateMsg.templateMsgMap["template_id"] = value
}

//SetPage 点击模板卡片后的跳转页面,仅限本小程序内的页面.支持带参数,(示例index?foo=bar).该字段不填则模板无跳转.
func (wxMaTemplateMsg *WxMaTemplateMsg) SetPage(value string) {
	if wxMaTemplateMsg.templateMsgMap == nil {
		wxMaTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMaTemplateMsg.templateMsgMap["page"] = value
}

//SetFormId 表单提交场景下,为 submit 事件带上的 formId；支付场景下,为本次支付的 prepay_id
func (wxMaTemplateMsg *WxMaTemplateMsg) SetFormId(value string) {
	if wxMaTemplateMsg.templateMsgMap == nil {
		wxMaTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMaTemplateMsg.templateMsgMap["form_id"] = value
}

//SetEmphasisKeyword 模板需要放大的关键词,不填则默认无放大
func (wxMaTemplateMsg *WxMaTemplateMsg) SetEmphasisKeyword(value string) {
	if wxMaTemplateMsg.templateMsgMap == nil {
		wxMaTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMaTemplateMsg.templateMsgMap["emphasis_keyword"] = value
}

//AddData 模板内容,不填则下发空模板.具体格式请参考示例.
func (wxMaTemplateMsg *WxMaTemplateMsg) AddData(key string, value string) {
	if wxMaTemplateMsg.dataMap == nil {
		wxMaTemplateMsg.dataMap = make(map[string]interface{})
	}
	valueMap := make(map[string]string)
	valueMap["value"] = value
	wxMaTemplateMsg.dataMap[key] = valueMap
}

//getTemplateMsgMap 获取参数的map,作为请求API的参数
func (wxMaTemplateMsg *WxMaTemplateMsg) getTemplateMsgMap() map[string]interface{} {
	if wxMaTemplateMsg.dataMap != nil {
		wxMaTemplateMsg.templateMsgMap["data"] = wxMaTemplateMsg.dataMap
	}
	return wxMaTemplateMsg.templateMsgMap
}

//WxMaTemplateMsgSend 发送模板消息
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/template-message/templateMessage.send.html
func WxMaTemplateMsgSend(wxMaConfig IWxMaConfig, wxMaTemplateMsg *WxMaTemplateMsg) (res ResponseBase, err error) {
	apiurl := WxMpAPIURL + "/cgi-bin/message/wxopen/template/send?access_token=" + wxMaConfig.GetAccessToken()

	data, err := httpPost(apiurl, wxMaTemplateMsg.getTemplateMsgMap())
	// 发送请求
	if err != nil {
		return res, err
	}
	// 尝试解码
	_ = json.Unmarshal(data, &res)

	return res, nil
}
