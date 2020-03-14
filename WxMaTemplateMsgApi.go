package gowe

import "context"

type WxMpTemplateMsg struct {
	templateMsgMap map[string]interface{}
	dataMap        map[string]interface{}
}

func (wxMpTemplateMsg *WxMpTemplateMsg) setTouser(value string) {
	if wxMpTemplateMsg.templateMsgMap == nil {
		wxMpTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMpTemplateMsg.templateMsgMap["touser"] = value
}
func (wxMpTemplateMsg *WxMpTemplateMsg) setTemplateId(value string) {
	if wxMpTemplateMsg.templateMsgMap == nil {
		wxMpTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMpTemplateMsg.templateMsgMap["template_id"] = value
}
func (wxMpTemplateMsg *WxMpTemplateMsg) setPage(value string) {
	if wxMpTemplateMsg.templateMsgMap == nil {
		wxMpTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMpTemplateMsg.templateMsgMap["page"] = value
}
func (wxMpTemplateMsg *WxMpTemplateMsg) setFormId(value string) {
	if wxMpTemplateMsg.templateMsgMap == nil {
		wxMpTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMpTemplateMsg.templateMsgMap["form_id"] = value
}
func (wxMpTemplateMsg *WxMpTemplateMsg) setEmphasisKeyword(value string) {
	if wxMpTemplateMsg.templateMsgMap == nil {
		wxMpTemplateMsg.templateMsgMap = make(map[string]interface{})
	}
	wxMpTemplateMsg.templateMsgMap["emphasis_keyword"] = value
}
func (wxMpTemplateMsg *WxMpTemplateMsg) AddData(key string, value string) {
	if wxMpTemplateMsg.dataMap == nil {
		wxMpTemplateMsg.dataMap = make(map[string]interface{})
	}
	valueMap := make(map[string]string)
	valueMap["value"] = value
	wxMpTemplateMsg.dataMap[key] = valueMap
}

func (wxMpTemplateMsg *WxMpTemplateMsg) getTemplateMsgMap() map[string]interface{} {
	if wxMpTemplateMsg.dataMap != nil {
		wxMpTemplateMsg.templateMsgMap["data"] = wxMpTemplateMsg.dataMap
	}
	return wxMpTemplateMsg.templateMsgMap
}

//WxMaTemplateMsgSend 发送模板消息
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/template-message/templateMessage.send.html
func WxMaTemplateMsgSend(ctx context.Context, wxMpTemplateMsg *WxMpTemplateMsg) (*APIResult, error) {
	apiurl := WxmpApiUrl + "/cgi-bin/message/wxopen/template/send?access_token=" + wxMaConfig.getAccessToken(ctx)

	resultMap, errMap := httpPostResultMap(apiurl, wxMpTemplateMsg.getTemplateMsgMap())
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}
