package gowe

import "encoding/json"

//模板消息
//https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html

//WxMTemplateMsgSend 发送模板消息
func WxMTemplateMsgSend(wxMpConfig IWxMpConfig, wxMpTemplateMsg *WxMaTemplateMsg) (res WxMTemplateMsgSendResponse, err error) {

	apiurl := WxMpAPIURL + "/cgi-bin/message/template/send?access_token=" + wxMpConfig.GetAccessToken()
	data, err := httpPost(apiurl, wxMpTemplateMsg.getTemplateMsgMap())
	// 发送请求
	if err != nil {
		return res, err
	}
	// 尝试解码
	_ = json.Unmarshal(data, &res)

	return res, nil
}

type WxMTemplateMsgSendResponse struct {
	MsgID   int64  `json:"msgid"`   // 用户唯一标识，调用成功后返回
	ErrCode int    `json:"errcode"` // 错误码
	ErrMsg  string `json:"errmsg"`  // 错误信息
}
