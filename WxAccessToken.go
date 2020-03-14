package gowe

import "time"

//WxAccessToken 微信accessToken
type WxAccessToken struct {
	AppId                  string
	AccessToken            string
	AccessTokenExpiresTime int64
	ExpiresIn              int
}

// token是否过期
func (wxAccessToken *WxAccessToken) IsAccessTokenExpired() bool {
	return time.Now().Unix() > wxAccessToken.AccessTokenExpiresTime
}

//WxCardTicket 微信卡券Ticket
type WxCardTicket struct {
	AppId                 string
	CardTicket            string
	CardTicketExpiresTime int64
	ExpiresIn             int
}

//微信卡券Ticket是否过期
func (wxCardTicket *WxCardTicket) IsCardTicketExpired() bool {
	return time.Now().Unix() > wxCardTicket.CardTicketExpiresTime
}

// 微信WxJsTicket
type WxJsTicket struct {
	AppId               string
	JsTicket            string
	JsTicketExpiresTime int64
	ExpiresIn           int
}

//WxJsTicket 是否过期
func (wxJsTicket *WxJsTicket) IsJsTicketExpired() bool {
	return time.Now().Unix() > wxJsTicket.JsTicketExpiresTime
}
