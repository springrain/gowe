package gowe

import "encoding/json"

//模板消息
//https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html

//WxMpTemplateMsgSend 发送模板消息
func WxMpTemplateMsgSend(wxMpConfig IWxMpConfig, body *WxMpTemplateMsgSendBody) (*WxMpTemplateMsgSendResponse, error) {

	apiurl := WxMpAPIURL + "/cgi-bin/message/template/send?access_token=" + wxMpConfig.GetAccessToken()

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

	if body.MaAppid != "" {
		maMap := make(map[string]string)
		maMap["appid"] = body.MaAppid
		if body.MaPagepath != "" {
			maMap["pagepath"] = body.MaPagepath
		}
		params["miniprogram"] = maMap
	}

	data, err := httpPost(apiurl, params)
	// 发送请求
	if err != nil {
		return nil, err
	}
	// 尝试解码
	res := &WxMpTemplateMsgSendResponse{}
	err = json.Unmarshal(data, res)

	return res, err
}

//https://api.weixin.qq.com/cgi-bin/message/template/subscribe?access_token=ACCESS_TOKEN
//WxMpSubscribeMsgSend 发送一次订阅消息

func WxMpSubscribeMsgSend(wxMpConfig IWxMpConfig, body *WxMpSubscribeMsgSendBody) (*WxMpTemplateMsgSendResponse, error) {

	apiurl := WxMpAPIURL + "/cgi-bin/message/template/subscribe?access_token=" + wxMpConfig.GetAccessToken()

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

	if body.MaAppid != "" {
		maMap := make(map[string]string)
		maMap["appid"] = body.MaAppid
		if body.MaPagepath != "" {
			maMap["pagepath"] = body.MaPagepath
		}
		params["miniprogram"] = maMap
	}

	data, err := httpPost(apiurl, params)
	// 发送请求
	if err != nil {
		return nil, err
	}
	// 尝试解码
	res := &WxMpTemplateMsgSendResponse{}
	err = json.Unmarshal(data, res)

	return res, err
}

//WxMpSubscribeMsgSendBody 一次订阅消息的请求参数

type WxMpSubscribeMsgSendBody struct{

	Touser          string                 `json:"touser"`                     // 接收者(用户)的 openid
	TemplateId      string                 `json:"template_id"`                // 所需下发的模板消息的id
	URL             string                 `json:"url,omitempty"`              // 模板跳转链接(海外帐号没有跳转能力)
	MaAppid         string                 `json:"-"`                          //需要跳转的小程序APPID
	MaPagepath      string                 `json:"-"`                          //所需跳转到小程序的具体页面路径,支持带参数,(示例index?foo=bar),要求该小程序已发布,暂不支持小游戏
	scene           string                 `json:"scene"`                    //订阅场景值
	title           string                 `json:"title,omitempty"` // 消息标题，15字以内
	dataMap         map[string]interface{} `json:"-"`                          //消息正文，value为消息内容文本（200字以内），没有固定格式，可用\n换行，color为整段消息内容的字体颜色（目前仅支持整段消息为一种颜色）
}


//WxMpTemplateMsgSendBody 模板消息的请求参数
type WxMpTemplateMsgSendBody struct {
	Touser          string                 `json:"touser"`                     // 接收者(用户)的 openid
	TemplateId      string                 `json:"template_id"`                // 所需下发的模板消息的id
	URL             string                 `json:"url,omitempty"`              // 模板跳转链接(海外帐号没有跳转能力)
	MaAppid         string                 `json:"-"`                          //需要跳转的小程序APPID
	MaPagepath      string                 `json:"-"`                          //所需跳转到小程序的具体页面路径,支持带参数,(示例index?foo=bar),要求该小程序已发布,暂不支持小游戏
	EmphasisKeyword string                 `json:"emphasis_keyword,omitempty"` // 模板需要放大的关键词,不填则默认无放大
	dataMap         map[string]interface{} `json:"-"`                          //模板数据
}

//WxMpTemplateMsgSendResponse 发送模板消息的返回值
type WxMpTemplateMsgSendResponse struct {
	MsgId   int64  `json:"msgid"`   // 用户唯一标识,调用成功后返回
	ErrCode int    `json:"errcode"` // 错误码
	ErrMsg  string `json:"errmsg"`  // 错误信息
}

//AddData 模板内容,不填则下发空模板.具体格式请参考示例,color默认#173177
func (wxMpTemplateMsg *WxMpTemplateMsgSendBody) AddData(key string, value string, color string) {
	if wxMpTemplateMsg.dataMap == nil {
		wxMpTemplateMsg.dataMap = make(map[string]interface{})
	}
	if color == "" {
		color = "#173177"
	}
	valueMap := make(map[string]string)
	valueMap["value"] = value
	valueMap["color"] = color
	wxMpTemplateMsg.dataMap[key] = valueMap
}
