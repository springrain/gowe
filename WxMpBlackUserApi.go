package gowe

import "context"

//WxMpBlackUserList 获取公众号的黑名单列表
//当 begin_openid 为空时，默认从开头拉取
//https://developers.weixin.qq.com/doc/offiaccount/User_Management/Manage_blacklist.html
func WxMpBlackUserList(ctx context.Context, wxMpConfig IWxMpConfig, beginOpenid string) (*APIResult, error) {
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
