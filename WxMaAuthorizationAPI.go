package gowe

import (
	"encoding/json"
	"errors"
)

// 小程序验证API接口

// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/user-info/auth.getPaidUnionId.html

//SessionKey jsCode换取用户信息,获得sessionKey
type SessionKey struct {
	OpenId     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionId    string `json:"unionid"`     // 只有在用户将公众号绑定到微信开放平台帐号后,才会出现该字段.
	ErrCode    int    `json:"errcode"`     // 错误码
	ErrMsg     string `json:"errmsg"`      // 错误信息
}

//WxMaCode2Session 登录凭证校验.通过 wx.login 接口获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func WxMaCode2Session(wxMaConfig IWxMaConfig, jsCode string) (*SessionKey, error) {
	if len(jsCode) < 1 {
		return nil, errors.New("jsCode不能为空")
	}
	apiurl := WxMpAPIURL + "/sns/jscode2session?appid=" + wxMaConfig.GetAppId() + "&secret=" + wxMaConfig.GetSecret() + "&js_code=" + jsCode + "&&grant_type=authorization_code"
	body, err := httpGet(apiurl)
	if err != nil {
		return nil, err
	}
	sessionKey := SessionKey{}
	err = json.Unmarshal(body, &sessionKey)
	return &sessionKey, err
}

//WxMaAuthGetPaidUnionId 用户支付完成后,获取该用户的 UnionId,无需用户授权.本接口支持第三方平台代理查询
func WxMaAuthGetPaidUnionId(wxMaConfig IWxMaConfig, openId string) (*WxMaAuthGetPaidUnionIdResponse, error) {
	if len(openId) < 1 {
		return nil, errors.New("openId不能为空")
	}

	apiurl := WxMpAPIURL + "/wxa/getpaidunionid?access_token=" + wxMaConfig.GetAccessToken() + "&openid=" + openId
	data, err := httpGet(apiurl)
	// 发送请求
	if err != nil {
		return nil, err
	}
	// 尝试解码
	res := WxMaAuthGetPaidUnionIdResponse{}
	err = json.Unmarshal(data, &res)

	return &res, err
}

//WxMaAuthGetPaidUnionIdResponse 支付后获取用户unionid
type WxMaAuthGetPaidUnionIdResponse struct {
	UnionId string `json:"unionid"` // 用户唯一标识,调用成功后返回
	ErrCode int    `json:"errcode"` // 错误码
	ErrMsg  string `json:"errmsg"`  // 错误信息
}
