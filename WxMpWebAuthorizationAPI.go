package gowe

import "context"

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
