package gowe

import (
	"context"
	"errors"
)

// 小程序验证API接口

// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/user-info/auth.getPaidUnionId.html

//WxMaCode2Session 登录凭证校验.通过 wx.login 接口获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func WxMaCode2Session(ctx context.Context, jsCode string) (*APIResult, error) {
	if len(jsCode) < 1 {
		return nil, errors.New("code不能为空")
	}
	wxMaConfig, errWxMaConfig := getWxMaConfig(ctx)
	if errWxMaConfig != nil {
		return nil, errWxMaConfig
	}
	apiurl := WxmpApiUrl + "/sns/jscode2session?appid=" + wxMaConfig.AppId + "&secret=" + wxMaConfig.Secret + "&js_code=" + jsCode + "&grant_type=authorization_code"
	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMaAuthGetPaidUnionId 用户支付完成后，获取该用户的 UnionId,无需用户授权.本接口支持第三方平台代理查询
func WxMaAuthGetPaidUnionId(ctx context.Context, openId string) (*APIResult, error) {
	if len(openId) < 1 {
		return nil, errors.New("openId不能为空")
	}
	wxMaConfig, errWxMaConfig := getWxMaConfig(ctx)
	if errWxMaConfig != nil {
		return nil, errWxMaConfig
	}
	apiurl := WxmpApiUrl + "/wxa/getpaidunionid?access_token=" + wxMaConfig.AccessToken + "&openid=" + openId
	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}
