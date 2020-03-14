package gowe

import "time"

//WxAccessToken 微信accessToken
type WxAccessToken struct {
	AppId                  string
	AccessToken            string
	accessTokenExpiresTime int64
	ExpiresIn              int
}

// token是否过期
func (wxAccessToken *WxAccessToken) IsAccessTokenExpired() bool {
	return time.Now().Unix() > wxAccessToken.accessTokenExpiresTime
}

//WxCardTicket 微信卡券Ticket
type WxCardTicket struct {
	AppId                 string
	CardTicket            string
	cardTicketExpiresTime int64
	ExpiresIn             int
}

//微信卡券Ticket是否过期
func (wxCardTicket *WxCardTicket) IsCardTicketExpired() bool {
	return time.Now().Unix() > wxCardTicket.cardTicketExpiresTime
}

// 微信WxJsTicket
type WxJsTicket struct {
	AppId               string
	JsTicket            string
	jsTicketExpiresTime int64
	ExpiresIn           int
}

//WxJsTicket 是否过期
func (wxJsTicket *WxJsTicket) IsJsTicketExpired() bool {
	return time.Now().Unix() > wxJsTicket.jsTicketExpiresTime
}
