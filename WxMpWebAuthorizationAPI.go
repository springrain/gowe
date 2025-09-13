package gowe

import (
	"context"
	"encoding/json"
	"errors"
)

//网页授权获取 access_token API
//https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html

// WxMpWrapAuthorizeURL 包装网页授权的url
func WxMpWrapAuthorizeURL(ctx context.Context, wxMpConfig IWxMpConfig, redirectUri string, snsapiBase bool) (string, error) {
	return wrapAuthorizeURL(ctx, wxMpConfig, redirectUri, "", snsapiBase)
}

// wrapAuthorizeURL 包装网页授权的url
func wrapAuthorizeURL(ctx context.Context, wxMpConfig IWxMpConfig, redirectUri string, state string, snsapiBase bool) (string, error) {
	apiurl := WxMpAPIURL + "/connect/oauth2/authorize?appid=" + wxMpConfig.GetAppId(ctx) + "&response_type=code&redirect_uri=" + redirectUri
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

// WxMpWebAuthAccessToken 用code换取accessToken. 认证的accessToken 和API的accessToken不一样,暂时使用同一个struct进行接收
func WxMpWebAuthAccessToken(ctx context.Context, wxMpConfig IWxMpConfig, code string) (*WxAccessToken, error) {
	if len(code) < 1 {
		return nil, errors.New("code不能为空")
	}
	apiurl := WxMpAPIURL + "/sns/oauth2?appid=" + wxMpConfig.GetAppId(ctx) + "&secret=" + wxMpConfig.GetSecret(ctx) + "&code=" + code + "&grant_type=authorization_code"
	body, err := httpGet(ctx, apiurl)
	if err != nil {
		return nil, err
	}
	wxAccessToken := &WxAccessToken{}
	err = json.Unmarshal(body, wxAccessToken)
	return wxAccessToken, err
}

// WxMpWebAuthRefreshAccessToken 刷新认证的AccessToken.认证的accessToken 和API的accessToken不一样,暂时使用同一个struct进行接收
func WxMpWebAuthRefreshAccessToken(ctx context.Context, wxMpConfig IWxMpConfig, refreshToken string) (*WxAccessToken, error) {
	if len(refreshToken) < 1 {
		return nil, errors.New("refreshToken不能为空")
	}
	apiurl := WxMpAPIURL + "/sns/oauth2/refresh_token?appid=" + wxMpConfig.GetAppId(ctx) + "&grant_type=refresh_token&refresh_token=" + refreshToken
	body, err := httpGet(ctx, apiurl)
	if err != nil {
		return nil, err
	}
	wxAccessToken := &WxAccessToken{}
	err = json.Unmarshal(body, wxAccessToken)
	return wxAccessToken, err
}

// WxMpUserInfo 用户信息
type WxMpUserInfo struct {
	OpenId         string   `json:"openid"`          // 用户唯一标识
	Nickname       string   `json:"nickname"`        // 用户的昵称
	Sex            int      `json:"sex"`             // 用户的性别,值为1时是男性,值为2时是女性,值为0时是未知
	Language       string   `json:"language"`        // 用户的语言,简体中文为zh_CN
	Province       string   `json:"province"`        // 用户所在省份
	City           string   `json:"city"`            // 用户所在城市
	Country        string   `json:"country"`         // 用户所在国家
	HeadimgUrl     string   `json:"headimgurl"`      // 用户头像,最后一个数值代表正方形头像大小(有0、46、64、96、132数值可选,0代表640*640正方形头像),用户没有头像时该项为空.若用户更换头像,原有头像URL将失效.
	Privilege      []string `json:"privilege"`       // 用户特权信息
	UnionId        string   `json:"unionid"`         // 只有在用户将公众号绑定到微信开放平台帐号后,才会出现该字段.
	Subscribe      int      `json:"subscribe"`       // 用户是否订阅该公众号标识,值为0时,代表此用户没有关注该公众号,拉取不到其余信息.
	SubscribeTime  int      `json:"subscribe_time"`  // 用户关注时间,为时间戳.如果用户曾多次关注,则取最后关注时间
	SubscribeScene string   `json:"subscribe_scene"` // 返回用户关注的渠道来源,ADD_SCENE_SEARCH 公众号搜索,ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移,ADD_SCENE_PROFILE_CARD 名片分享,ADD_SCENE_QR_CODE 扫描二维码,ADD_SCENEPROFILE LINK 图文页内名称点击,ADD_SCENE_PROFILE_ITEM 图文页右上角菜单,ADD_SCENE_PAID 支付后关注,ADD_SCENE_OTHERS 其他
	Remark         string   `json:"remark"`          // 公众号运营者对粉丝的备注,公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	GroupId        int      `json:"groupid"`         // 用户所在的分组ID(兼容旧的用户分组接口)
	TagidList      []int    `json:"tagid_list"`      // 用户被打上的标签ID列表
	QrScene        int      `json:"qr_scene"`        // 二维码扫码场景(开发者自定义)
	QrSceneStr     string   `json:"qr_scene_str"`    // 二维码扫码场景描述(开发者自定义)
	ErrCode        int      `json:"errcode"`         // 错误码
	ErrMsg         string   `json:"errmsg"`          // 错误信息
}

// WxMpGetUserInfo 获取用户基本信息,包括UnionID机制(授权机制)
// https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html
func WxMpGetUserInfo(ctx context.Context, wxMpConfig IWxMpConfig, openId, lang string) (*WxMpUserInfo, error) {
	if len(lang) <= 0 {
		lang = "zh_CN"
	}
	apiurl := WxMpAPIURL + "/cgi-bin/user/info?access_token=" + wxMpConfig.GetAccessToken(ctx) + "&openid=" + openId + "&lang=" + lang
	body, err := httpGet(ctx, apiurl)
	if err != nil {
		return nil, err
	}
	userInfo := WxMpUserInfo{}
	err = json.Unmarshal(body, &userInfo)
	return &userInfo, err
}

type WxMpUserOpenIds struct {
	Total      int                 `json:"total"`       //	关注该公众账号的总用户数
	Count      int                 `json:"count"`       //	拉取的OPENID个数，最大值为10000
	Data       map[string][]string `json:"data"`        //列表数据，OPENID的列表
	NextOpenid string              `json:"next_openid"` //拉取列表的最后一个用户的OPENID
	ErrCode    int                 `json:"errcode"`     // 错误码
	ErrMsg     string              `json:"errmsg"`      // 错误信息
}

// WxMpGetUserInfos 获取用户基本信息,包括UnionID机制(授权机制)
// https://developers.weixin.qq.com/doc/offiaccount/User_Management/Getting_a_User_List.html
func WxMpGetUserOpenIds(ctx context.Context, wxMpConfig IWxMpConfig, openId string) (*WxMpUserOpenIds, error) {

	apiurl := WxMpAPIURL + "/cgi-bin/user/get?access_token=" + wxMpConfig.GetAccessToken(ctx) + "&next_openid=" + openId
	body, err := httpGet(ctx, apiurl)
	if err != nil {
		return nil, err
	}
	userOpenIds := WxMpUserOpenIds{}
	err = json.Unmarshal(body, &userOpenIds)
	return &userOpenIds, err
}
