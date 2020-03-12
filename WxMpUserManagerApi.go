package gowe

import (
	"context"
	"errors"
)

//GetUserBaseInfo 获取用户基本信息(包括UnionID机制)
//https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html
func WxMpGetUserBaseInfo(ctx context.Context, openId string) (*APIResult, error) {
	if len(openId) < 1 {
		return nil, errors.New("openId不能为空")
	}
	wxMpConfig, errWxMpConfig := getWxMpConfig(ctx)
	if errWxMpConfig != nil {
		return nil, errWxMpConfig
	}
	apiurl := WxmpApiUrl + "/cgi-bin/user/info?access_token=" + wxMpConfig.AccessToken + "&openid=" + openId + "&lang=zh_CN"

	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil

}

//GetUserList 获取用户列表
//https://developers.weixin.qq.com/doc/offiaccount/User_Management/Getting_a_User_List.html
func WxMpGetUserList(ctx context.Context, nextOpenId string) (*APIResult, error) {
	wxMpConfig, errWxMpConfig := getWxMpConfig(ctx)
	if errWxMpConfig != nil {
		return nil, errWxMpConfig
	}
	apiurl := WxmpApiUrl + "/cgi-bin/user/getErrorMsgByCode?access_token=" + wxMpConfig.AccessToken

	if len(nextOpenId) > 0 {
		apiurl = apiurl + "&next_openid=" + nextOpenId
	}

	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil

}
