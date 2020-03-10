package gowe

import (
	"context"
	"errors"
)

//WxMpBlackUserOpenIdList 获取公众号的黑名单列表
//当 begin_openid 为空时，默认从开头拉取
//https://developers.weixin.qq.com/doc/offiaccount/User_Management/Manage_blacklist.html
func WxMpBlackUserOpenIdList(ctx context.Context, wxMpConfig IWxMpConfig, beginOpenid string) (*APIResult, error) {
	apiurl := WxmpApiUrl + "/cgi-bin/tags/members/getblacklist?access_token=" + wxMpConfig.GetAccessToken()

	parm := make(map[string]interface{})
	if len(beginOpenid) > 0 {
		parm["begin_openid"] = beginOpenid
	}

	resultMap, errMap := httpPostResultMap(apiurl, parm)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxBatchBlackUserOpenId  批量拉黑用户
//openIds 需要拉黑的用户openId列表
func WxBatchBlackUserOpenId(ctx context.Context, wxMpConfig IWxMpConfig, openIds []string) (*APIResult, error) {
	apiurl := WxmpApiUrl + "/cgi-bin/tags/members/batchblacklist?access_token=" + wxMpConfig.GetAccessToken()
	if len(openIds) < 1 {
		return nil, errors.New("需要拉黑的用户openId列表不能为空")
	}

	parm := make(map[string]interface{})
	parm["openid_list"] = openIds
	resultMap, errMap := httpPostResultMap(apiurl, parm)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxBatchUnBlackUserOpenId  批量解封拉黑的用户
//openIds 需要解封的用户openId列表
func WxBatchUnBlackUserOpenId(ctx context.Context, wxMpConfig IWxMpConfig, openIds []string) (*APIResult, error) {
	apiurl := WxmpApiUrl + "/cgi-bin/tags/members/batchunblacklist?access_token=" + wxMpConfig.GetAccessToken()
	if len(openIds) < 1 {
		return nil, errors.New("需要拉黑的用户openId列表不能为空")
	}

	parm := make(map[string]interface{})
	parm["openid_list"] = openIds
	resultMap, errMap := httpPostResultMap(apiurl, parm)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}
