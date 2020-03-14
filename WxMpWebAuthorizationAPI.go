package gowe

import (
	"errors"
)

//网页授权获取 access_token API
//https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html

func WxMpWrapAuthorizeURL(wxMpConfig IWxMpConfig, redirectUri string, snsapiBase bool) (string, error) {
	return wrapAuthorizeURL(wxMpConfig, redirectUri, "", snsapiBase)
}

func wrapAuthorizeURL(wxMpConfig IWxMpConfig, redirectUri string, state string, snsapiBase bool) (string, error) {

	apiurl := WxmpApiUrl + "/connect/oauth2/authorize?appid=" + wxMpConfig.GetAppId() + "&response_type=code&redirect_uri=" + redirectUri

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

//WxMpAuthAccessTokenByCode 用code换取accessToken, 认证的accessToken 和API的accessToken不一样
func WxMpAuthAccessTokenByCode(wxMpConfig IWxMpConfig, code string) (*APIResult, error) {
	if len(code) < 1 {
		return nil, errors.New("code不能为空")
	}

	apiurl := WxmpApiUrl + "/sns/oauth2?appid=" + wxMpConfig.GetAppId() + "&secret=" + wxMpConfig.GetSecret() + "&code=" + code + "&grant_type=authorization_code"
	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}
