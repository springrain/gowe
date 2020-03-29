package gowe

import (
	"encoding/json"
	"net/url"
)

//模板消息
//https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
//https://developers.weixin.qq.com/doc/offiaccount/Message_Management/One-time_subscription_info.html

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

//WxMpSubscribeMsgURL 构造订阅模板消息的授权URL
func WxMpSubscribeMsgURL(body *WxMpSubscribeMsgURLBody) (string, error) {

	apiurl := WxMpWeiXinURL + "/mp/subscribemsg?action=get_confirm&appid=" + body.AppId + "&scene=" + body.Scene + "&template_id=" + body.TemplateId + "&redirect_url=" + url.QueryEscape(body.RedirectURL)
	if len(body.Reserved) >= 0 {
		apiurl = apiurl + "&reserved=" + body.Reserved
	}
	apiurl = apiurl + "#wechat_redirect"
	return apiurl, nil

}

//WxMpSubscribeMsgSend 发送一次订阅消息
//https://developers.weixin.qq.com/doc/offiaccount/Message_Management/One-time_subscription_info.html
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

//用于生成授权URL的实体类
type WxMpSubscribeMsgURLBody struct {
	AppId       string //公众号的唯一标识
	Scene       string //重定向后会带上scene参数，开发者可以填0-10000的整形值，用来标识订阅场景值
	TemplateId  string //订阅消息模板ID，登录公众平台后台，在接口权限列表处可查看订阅模板ID
	RedirectURL string //授权后重定向的回调地址，请使用UrlEncode对链接进行处理。 注：要求redirect_url的域名要跟登记的业务域名一致，且业务域名不能带路径。 业务域名需登录公众号，在设置-公众号设置-功能设置里面对业务域名设置。
	Reserved    string //用于保持请求和回调的状态，授权请后原样带回给第三方。该参数可用于防止csrf攻击（跨站请求伪造攻击），建议第三方带上该参数，可设置为简单的随机数加session进行校验，开发者可以填写a-zA-Z0-9的参数值，最多128字节，要求做urlencode
}

//WxMpSubscribeMsgSendBody 一次订阅消息的请求参数
type WxMpSubscribeMsgSendBody struct {
	wxMpTemplateMsgBody
	Scene string `json:"scene"` //订阅场景值
	Title string `json:"title"` // 消息标题，15字以内
}

//WxMpTemplateMsgSendBody 模板消息的请求参数
type WxMpTemplateMsgSendBody struct {
	wxMpTemplateMsgBody
	EmphasisKeyword string `json:"emphasis_keyword,omitempty"` // 模板需要放大的关键词,不填则默认无放大

}

//wxMpTemplateMsgBody 公用的模板消息参数
type wxMpTemplateMsgBody struct {
	Touser     string                 `json:"touser"`        // 接收者(用户)的 openid
	TemplateId string                 `json:"template_id"`   // 所需下发的模板消息的id
	URL        string                 `json:"url,omitempty"` // 模板跳转链接(海外帐号没有跳转能力)
	MaAppid    string                 `json:"-"`             //需要跳转的小程序APPID
	MaPagepath string                 `json:"-"`             //所需跳转到小程序的具体页面路径,支持带参数,(示例index?foo=bar),要求该小程序已发布,暂不支持小游戏
	dataMap    map[string]interface{} `json:"-"`             //模板数据
}

//WxMpTemplateMsgSendResponse 发送模板消息的返回值
type WxMpTemplateMsgSendResponse struct {
	MsgId   int64  `json:"msgid"`   // 用户唯一标识,调用成功后返回
	ErrCode int    `json:"errcode"` // 错误码
	ErrMsg  string `json:"errmsg"`  // 错误信息
}

//AddData 模板内容,不填则下发空模板.具体格式请参考示例,color默认#173177
func (wxMpTemplateMsg *wxMpTemplateMsgBody) AddData(key string, value string, color string) {
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
