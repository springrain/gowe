package gowe

import (
	"encoding/json"
	"errors"
)

//网页授权获取 access_token API
//https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html

//WxMpWrapAuthorizeURL 包装网页授权的url
func WxMpWrapAuthorizeURL(wxMpConfig IWxMpConfig, redirectUri string, snsapiBase bool) (string, error) {
	return wrapAuthorizeURL(wxMpConfig, redirectUri, "", snsapiBase)
}

//wrapAuthorizeURL 包装网页授权的url
func wrapAuthorizeURL(wxMpConfig IWxMpConfig, redirectUri string, state string, snsapiBase bool) (string, error) {
	apiurl := WxMpAPIURL + "/connect/oauth2/authorize?appid=" + wxMpConfig.GetAppId() + "&response_type=code&redirect_uri=" + redirectUri
	if snsapiBase {
		apiurl = apiurl + "&scope=snsapi_base"
	} else {
		apiurl = apiurl + "&scope=snsapi_userinfo"
	}
	if state == "" {
		state = "wx"
	}
	apiurl = apiurl + "&state=" + state + "#wechat_redirect"

	return apiurl, nil

}

//WxMpWebAuthAccessToken 用code换取accessToken. 认证的accessToken 和API的accessToken不一样,暂时使用同一个struct进行接收
func WxMpWebAuthAccessToken(wxMpConfig IWxMpConfig, code string) (*WxAccessToken, error) {
	if len(code) < 1 {
		return nil, errors.New("code不能为空")
	}
	apiurl := WxMpAPIURL + "/sns/oauth2?appid=" + wxMpConfig.GetAppId() + "&secret=" + wxMpConfig.GetSecret() + "&code=" + code + "&grant_type=authorization_code"
	body, err := httpGet(apiurl)
	if err != nil {
		return nil, err
	}
	wxAccessToken := WxAccessToken{}
	err = json.Unmarshal(body, &wxAccessToken)
	return &wxAccessToken, err
}

//WxMpWebAuthRefreshAccessToken 刷新认证的AccessToken.认证的accessToken 和API的accessToken不一样,暂时使用同一个struct进行接收
func WxMpWebAuthRefreshAccessToken(wxMpConfig IWxMpConfig, refreshToken string) (*WxAccessToken, error) {
	if len(refreshToken) < 1 {
		return nil, errors.New("refreshToken不能为空")
	}
	apiurl := WxMpAPIURL + "/sns/oauth2/refresh_token?appid=" + wxMpConfig.GetAppId() + "&grant_type=refresh_token&refresh_token=" + refreshToken
	body, err := httpGet(apiurl)
	if err != nil {
		return nil, err
	}
	wxAccessToken := WxAccessToken{}
	err = json.Unmarshal(body, &wxAccessToken)
	return &wxAccessToken, err
}
