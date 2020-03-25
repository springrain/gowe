package gowe

import "encoding/json"

//WxMaTemplateMsgSend 发送模板消息
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/template-message/templateMessage.send.html
func WxMaTemplateMsgSend(wxMaConfig IWxMaConfig, body *WxMaTemplateMsgSendBody) (*ResponseBase, error) {
	apiurl := WxMpAPIURL + "/cgi-bin/message/wxopen/template/send?access_token=" + wxMaConfig.GetAccessToken()
	return wxMaSendTemplateMsg(wxMaConfig, apiurl, body)
}

//WxMaSubscribeMessageSend 发送订阅消息
//https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html
func WxMaSubscribeMessageSend(wxMaConfig IWxMaConfig, body *WxMaTemplateMsgSendBody) (*ResponseBase, error) {
	apiurl := WxMpAPIURL + "/cgi-bin/message/subscribe/send?access_token=" + wxMaConfig.GetAccessToken()
	return wxMaSendTemplateMsg(wxMaConfig, apiurl, body)
}

func wxMaSendTemplateMsg(wxMaConfig IWxMaConfig, apiurl string, body *WxMaTemplateMsgSendBody) (*ResponseBase, error) {

	// 参数处理
	bodyStr, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	params := make(map[string]interface{})
	if err = json.Unmarshal(bodyStr, &params); err != nil {
		return nil, err
	}

	if body.dataMap != nil {
		params["data"] = body.dataMap
	}

	data, err := httpPost(apiurl, params)
	// 发送请求
	if err != nil {
		return nil, err
	}
	// 尝试解码
	res := &ResponseBase{}
	err = json.Unmarshal(data, res)
	return res, err
}

//WxMaTemplateMsgSendBody 小程序模板消息的请求参数
type WxMaTemplateMsgSendBody struct {
	Touser          string                 `json:"touser"`                     // 接收者(用户)的 openid
	TemplateId      string                 `json:"template_id"`                // 所需下发的模板消息的id
	Page            string                 `json:"page,omitempty"`             // 点击模板卡片后的跳转页面,仅限本小程序内的页面.支持带参数,(示例index?foo=bar).该字段不填则模板无跳转.
	FormId          string                 `json:"form_id"`                    //表单提交场景下,为 submit 事件带上的 formId;支付场景下,为本次支付的 prepay_id
	EmphasisKeyword string                 `json:"emphasis_keyword,omitempty"` // 模板需要放大的关键词,不填则默认无放大
	dataMap         map[string]interface{} `json:"-"`                          //模板数据

	//发送订阅消息的属性
	MiniprogramState string `json:"miniprogram_state,omitempty"` //跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
	Lang             string `json:"lang,omitempty"`              //进入小程序查看”的语言类型，支持zh_CN(简体中文)、en_US(英文)、zh_HK(繁体中文)、zh_TW(繁体中文)，默认为zh_CN
}

//AddData 模板内容,不填则下发空模板.具体格式请参考示例.
func (wxMaTemplateMsg *WxMaTemplateMsgSendBody) AddData(key string, value string) {
	if wxMaTemplateMsg.dataMap == nil {
		wxMaTemplateMsg.dataMap = make(map[string]interface{})
	}
	valueMap := make(map[string]string)
	valueMap["value"] = value
	wxMaTemplateMsg.dataMap[key] = valueMap
}
