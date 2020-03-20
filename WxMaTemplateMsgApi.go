package gowe

import "encoding/json"

//WxMaTemplateMsgSend 发送模板消息
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/template-message/templateMessage.send.html
func WxMaTemplateMsgSend(wxMaConfig IWxMaConfig, body WxMaTemplateMsgSendBody) (res ResponseBase, err error) {
	apiurl := WxMpAPIURL + "/cgi-bin/message/wxopen/template/send?access_token=" + wxMaConfig.GetAccessToken()

	// 参数处理
	bodyStr, err := json.Marshal(body)
	if err != nil {
		return
	}

	params := make(map[string]interface{})
	if err = json.Unmarshal(bodyStr, &params); err != nil {
		return
	}

	if body.dataMap != nil {
		params["data"] = body.dataMap
	}

	data, err := httpPost(apiurl, params)
	// 发送请求
	if err != nil {
		return res, err
	}
	// 尝试解码
	_ = json.Unmarshal(data, &res)

	return res, nil
}

type WxMaTemplateMsgSendBody struct {
	Touser          string                 `json:"touser"`                     // 接收者(用户)的 openid
	TemplateId      string                 `json:"template_id"`                // 所需下发的模板消息的id
	Page            string                 `json:"page,omitempty"`             // 点击模板卡片后的跳转页面,仅限本小程序内的页面.支持带参数,(示例index?foo=bar).该字段不填则模板无跳转.
	FormId          string                 `json:"form_id"`                    //表单提交场景下,为 submit 事件带上的 formId；支付场景下,为本次支付的 prepay_id
	EmphasisKeyword string                 `json:"emphasis_keyword,omitempty"` // 模板需要放大的关键词,不填则默认无放大
	dataMap         map[string]interface{} `json:"-"`                          //模板数据
}

//AddData 模板内容,不填则下发空模板.具体格式请参考示例.
func (wxMaTemplateMsg WxMaTemplateMsgSendBody) AddData(key string, value string) {
	if wxMaTemplateMsg.dataMap == nil {
		wxMaTemplateMsg.dataMap = make(map[string]interface{})
	}
	valueMap := make(map[string]string)
	valueMap["value"] = value
	wxMaTemplateMsg.dataMap[key] = valueMap
}
