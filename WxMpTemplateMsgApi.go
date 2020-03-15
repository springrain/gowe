package gowe

import (
	"errors"
)

//模板消息
//https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html

//WxMTemplateMsgSend 发送模板消息
func WxMTemplateMsgSend(wxMpConfig IWxMpConfig, wxMpTemplateMsg *WxMpTemplateMsg) (*APIResult, error) {

	apiurl := WxMpAPIURL + "/cgi-bin/message/template/send?access_token=" + wxMpConfig.GetAccessToken()
	resultMap, errMap := httpPostResultMap(apiurl, wxMpTemplateMsg.getTemplateMsgMap())
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMpTemplateMsgSetIndustry 设置所属行业
func WxMpTemplateMsgSetIndustry(wxMpConfig IWxMpConfig, industryId1 string, industryId2 string) (*APIResult, error) {
	if len(industryId1) < 1 || len(industryId2) < 1 {
		return nil, errors.New("industry_id1或者industry_id2不能为空")
	}

	apiurl := WxMpAPIURL + "/cgi-bin/template/api_set_industry?access_token=" + wxMpConfig.GetAccessToken()

	params := make(map[string]interface{})
	params["industry_id1"] = industryId1
	params["industry_id2"] = industryId2
	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMpTemplateMsgGetIndustry 获取帐号设置的行业信息
func WxMpTemplateMsgGetIndustry(wxMpConfig IWxMpConfig) (*APIResult, error) {

	apiurl := WxMpAPIURL + "/cgi-bin/template/get_industry?access_token=" + wxMpConfig.GetAccessToken()
	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMpTemplateMsgGetTemplateId 获得模板ID  templateIdShort 模板库中模板的编号，有“TM**”和“OPENTMTM**”等形式
//从行业模板库选择模板到帐号后台,获得模板ID的过程可在微信公众平台后台完成.为方便第三方开发者,提供通过接口调用的方式来获取模板ID
func WxMpTemplateMsgGetTemplateId(wxMpConfig IWxMpConfig, templateIdShort string) (*APIResult, error) {
	if len(templateIdShort) < 1 {
		return nil, errors.New("templateIdShort不能为空")
	}

	apiurl := WxMpAPIURL + "/cgi-bin/template/api_set_industry?access_token=" + wxMpConfig.GetAccessToken()

	params := make(map[string]interface{})
	params["template_id_short"] = templateIdShort
	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMpTemplateMsgGetAllTemplate 获取模板列表
func WxMpTemplateMsgGetAllTemplate(wxMpConfig IWxMpConfig) (*APIResult, error) {

	apiurl := WxMpAPIURL + "/cgi-bin/template/get_all_private_template?access_token=" + wxMpConfig.GetAccessToken()
	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMpTemplateMsgDelPrivateTemplate  删除模板
func WxMpTemplateMsgDelPrivateTemplate(wxMpConfig IWxMpConfig, templateId string) (*APIResult, error) {
	if len(templateId) < 1 {
		return nil, errors.New("templateId不能为空")
	}

	apiurl := WxMpAPIURL + "/cgi-bin/template/del_private_template?access_token=" + wxMpConfig.GetAccessToken()

	params := make(map[string]interface{})
	params["template_id"] = templateId
	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

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
