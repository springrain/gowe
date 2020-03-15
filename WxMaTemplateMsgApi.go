package gowe

//WxMaTemplateMsgSend 发送模板消息
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/template-message/templateMessage.send.html

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
