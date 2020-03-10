package gowe

import (
	"context"
	"errors"
	"time"
)

/**
 * 认证并获取 access_token API
 * https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/Nontax_Bill/API_list.html
 * <p>
 * 生成签名之前必须先了解一下jsapi_ticket，jsapi_ticket是公众号用于调用微信JS接口的临时票据
 * https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html
 * https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#62
 * <p>
 * 微信卡券接口签名凭证 api_ticket
 * https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#54
 */

var accessTokenURL = WxmpApiUrl + "/cgi-bin/token"

//GetAccessToken 获取 access token，如果未取到或者 access token 不可用则先更新再获取
func GetAccessToken(ctx context.Context, wxConfig IWxConfig) (*WxAccessToken, error) {
	apiurl := accessTokenURL + "?grant_type=client_credential&appid=" + wxConfig.GetAppId() + "&secret=" + wxConfig.GetSecret()

	resultMap, errMap := httpGetMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}

	expires := mapGetInt(resultMap, "expires_in")
	if expires == errExpires {
		return nil, errors.New("expires_in超时时间错误")
	}

	accessToken := mapGetString(resultMap, "access_token")

	if len(accessToken) < 1 {
		return nil, errors.New("未能获得accessToken")
	}

	wxAccessToken := WxAccessToken{}
	wxAccessToken.AppId = wxConfig.GetAppId()
	wxAccessToken.AccessToken = accessToken
	wxAccessToken.ExpiresIn = expires
	// 生产遇到接近过期时间时,access_token在某些服务器上会提前失效,设置时间短一些
	// https://developers.weixin.qq.com/community/develop/doc/0008cc492503e8e04dc7d619754c00
	wxAccessToken.accessTokenExpiresTime = time.Now().Unix() + int64(wxAccessToken.ExpiresIn/2)

	return &wxAccessToken, nil
}
