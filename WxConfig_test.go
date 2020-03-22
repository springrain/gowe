package gowe

import (
	"fmt"
	"testing"
)

type WxConfig struct {
	Id     string
	AppId  string
	Secret string
}

var wxConfig = &WxConfig{
	Id:     "test",
	AppId:  "XXXXXXXXXXXXXXxxx",
	Secret: "XXXXXXXXXXXXXXX",
}

func (wxConfig *WxConfig) GetId() string {
	return wxConfig.Id
}

func (wxConfig *WxConfig) GetAppId() string {
	return wxConfig.AppId
}

func (wxConfig *WxConfig) GetAccessToken() string {
	//正常应该把wxAccessToken缓存起来,从缓存中获取
	wxAccessToken, err := GetAccessToken(wxConfig)
	fmt.Println(err)
	return wxAccessToken.AccessToken
}

func (wxConfig *WxConfig) GetSecret() string {
	return wxConfig.Secret
}

func TestToken(t *testing.T) {
	wxConfig.GetAccessToken()
}
