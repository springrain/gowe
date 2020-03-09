package gowe

import "time"

type WxAccessToken struct {
	AppId                  string
	AccessToken            string
	accessTokenExpiresTime int64
	ExpiresIn              int
}

func (wxAccessToken *WxAccessToken) isAccessTokenExpired() bool {
	return time.Now().Unix() > wxAccessToken.accessTokenExpiresTime
}

type WxCardTicket struct {
	AppId                 string
	CardTicket            string
	cardTicketExpiresTime int64
	ExpiresIn             int
}

func (wxCardTicket *WxCardTicket) isCardTicketExpired() bool {
	return time.Now().Unix() > wxCardTicket.cardTicketExpiresTime
}

type WxJsTicket struct {
	AppId               string
	JsTicket            string
	jsTicketExpiresTime int64
	ExpiresIn           int
}

func (wxJsTicket *WxJsTicket) isJsTicketExpired() bool {
	return time.Now().Unix() > wxJsTicket.jsTicketExpiresTime
}
