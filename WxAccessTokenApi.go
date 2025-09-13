package gowe

import (
	"context"
	"encoding/json"
	"time"
)

/**
 * 认证并获取 access_token API
 * https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/Nontax_Bill/API_list.html
 * <p>
 * 生成签名之前必须先了解一下jsapi_ticket,jsapi_ticket是公众号用于调用微信JS接口的临时票据
 * https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html
 * https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDKs.html#62
 * <p>
 * 微信卡券接口签名凭证 api_ticket
 * https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#54
 */

// WxAccessToken 微信accessToken
type WxAccessToken struct {
	AccessToken            string `json:"access_token"`  // 获取到的凭证
	ExpiresIn              int    `json:"expires_in"`    // SessionKey超时时间(秒)
	RefreshToken           string `json:"refresh_token"` // 用户刷新access_tokenOpenId
	OpenId                 string `json:"openid"`        // 用户唯一标识
	Scope                  string `json:"scope"`         // 用户授权的作用域
	ErrCode                int    `json:"errcode"`       // 错误码
	ErrMsg                 string `json:"errmsg"`        // 错误信息
	AccessTokenExpiresTime int64  //过期的时间,验证是否过期
}

// IsAccessTokenExpired token是否过期
func (wxAccessToken *WxAccessToken) IsAccessTokenExpired() bool {
	return time.Now().Unix() > wxAccessToken.AccessTokenExpiresTime
}

// WxCardTicket 微信卡券Ticket
type WxCardTicket struct {
	CardTicket            string `json:"ticket"`     // 获取到的凭证
	ExpiresIn             int    `json:"expires_in"` // SessionKey超时时间(秒)
	ErrCode               int    `json:"errcode"`    // 错误码
	ErrMsg                string `json:"errmsg"`     // 错误信息
	CardTicketExpiresTime int64
}

// IsCardTicketExpired 微信卡券Ticket是否过期
func (wxCardTicket *WxCardTicket) IsCardTicketExpired() bool {
	return time.Now().Unix() > wxCardTicket.CardTicketExpiresTime
}

// WxJsTicket 微信WxJsTicket
type WxJsTicket struct {
	JsTicket            string `json:"ticket"`     // 获取到的凭证
	ExpiresIn           int    `json:"expires_in"` // SessionKey超时时间(秒)
	ErrCode             int    `json:"errcode"`    // 错误码
	ErrMsg              string `json:"errmsg"`     // 错误信息
	JsTicketExpiresTime int64
}

// IsJsTicketExpired WxJsTicket 是否过期
func (wxJsTicket *WxJsTicket) IsJsTicketExpired() bool {
	return time.Now().Unix() > wxJsTicket.JsTicketExpiresTime
}

// GetAccessToken 获取 access token,如果未取到或者 access token 不可用则先更新再获取
func GetAccessToken(ctx context.Context, wxConfig IWxConfig) (*WxAccessToken, error) {
	apiurl := WxMpAPIURL + "/cgi-bin/token?grant_type=client_credential&appid=" + wxConfig.GetAppId(ctx) + "&secret=" + wxConfig.GetSecret(ctx)
	body, err := httpGet(ctx, apiurl)
	if err != nil {
		return nil, err
	}

	wxAccessToken := &WxAccessToken{}
	err = json.Unmarshal(body, wxAccessToken)
	if err != nil {
		return wxAccessToken, err
	}
	// 生产遇到接近过期时间时,access_token在某些服务器上会提前失效,设置时间短一些
	// https://developers.weixin.qq.com/community/develop/doc/0008cc492503e8e04dc7d619754c00
	wxAccessToken.AccessTokenExpiresTime = time.Now().Unix() + int64(wxAccessToken.ExpiresIn/2)

	return wxAccessToken, err
}

// GetJsTicket 获取jsTicket
func GetJsTicket(ctx context.Context, wxConfig IWxConfig) (*WxJsTicket, error) {
	apiurl := WxMpAPIURL + "/cgi-bin/ticket/getticket?access_token=" + wxConfig.GetAccessToken(ctx) + "&type=jsapi"
	body, err := httpGet(ctx, apiurl)
	if err != nil {
		return nil, err
	}
	wxJsTicket := &WxJsTicket{}
	err = json.Unmarshal(body, wxJsTicket)
	if err != nil {
		return wxJsTicket, err
	}

	// 生产遇到接近过期时间时,access_token在某些服务器上会提前失效,设置时间短一些
	// https://developers.weixin.qq.com/community/develop/doc/0008cc492503e8e04dc7d619754c00
	wxJsTicket.JsTicketExpiresTime = time.Now().Unix() + int64(wxJsTicket.ExpiresIn/2)
	return wxJsTicket, err
}

// GetCardTicket 获取cardTicket
func GetCardTicket(ctx context.Context, wxConfig IWxConfig) (*WxCardTicket, error) {
	apiurl := WxMpAPIURL + "/cgi-bin/ticket/getticket?access_token=" + wxConfig.GetAccessToken(ctx) + "&type=wx_card"
	body, err := httpGet(ctx, apiurl)
	if err != nil {
		return nil, err
	}
	wxCardTicket := &WxCardTicket{}
	err = json.Unmarshal(body, wxCardTicket)
	if err != nil {
		return wxCardTicket, err
	}
	// 生产遇到接近过期时间时,access_token在某些服务器上会提前失效,设置时间短一些
	// https://developers.weixin.qq.com/community/develop/doc/0008cc492503e8e04dc7d619754c00
	wxCardTicket.CardTicketExpiresTime = time.Now().Unix() + int64(wxCardTicket.ExpiresIn/2)
	return wxCardTicket, err
}
