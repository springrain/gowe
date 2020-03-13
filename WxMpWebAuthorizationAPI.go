package gowe

import (
	"context"
	"errors"
)

//网页授权获取 access_token API
//https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html

func WxMpWrapAuthorizeURL(ctx context.Context, redirectUri string, snsapiBase bool) (string, error) {
	return wrapAuthorizeURL(ctx, redirectUri, "", snsapiBase)
}

func wrapAuthorizeURL(ctx context.Context, redirectUri string, state string, snsapiBase bool) (string, error) {
	wxMpConfig, errWxMpConfig := getWxMpConfig(ctx)
	if errWxMpConfig != nil {
		return "", errWxMpConfig
	}
	apiurl := WxmpApiUrl + "/connect/oauth2/authorize?appid=" + wxMpConfig.AppId + "&response_type=code&redirect_uri=" + redirectUri

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
func WxMpAuthAccessTokenByCode(ctx context.Context, code string) (*APIResult, error) {
	if len(code) < 1 {
		return nil, errors.New("code不能为空")
	}
	wxMpConfig, errWxMpConfig := getWxMpConfig(ctx)
	if errWxMpConfig != nil {
		return nil, errWxMpConfig
	}
	apiurl := WxmpApiUrl + "/sns/oauth2?appid=" + wxMpConfig.AppId + "&secret=" + wxMpConfig.Secret + "&code=" + code + "&grant_type=authorization_code"
	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}
