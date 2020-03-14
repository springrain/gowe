package gowe

import (
	"errors"
)

//GetUserBaseInfo 获取用户基本信息(包括UnionID机制)
//https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html
func WxMpGetUserBaseInfo(wxMpConfig IWxMpConfig, openId string) (*APIResult, error) {
	if len(openId) < 1 {
		return nil, errors.New("openId不能为空")
	}

	apiurl := WxMpAPIURL + "/cgi-bin/user/info?access_token=" + wxMpConfig.GetAccessToken() + "&openid=" + openId + "&lang=zh_CN"

	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil

}

//GetUserList 获取用户列表
//https://developers.weixin.qq.com/doc/offiaccount/User_Management/Getting_a_User_List.html
func WxMpGetUserList(wxMpConfig IWxMpConfig, nextOpenId string) (*APIResult, error) {

	apiurl := WxMpAPIURL + "/cgi-bin/user/getErrorMsgByCode?access_token=" + wxMpConfig.GetAccessToken()

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

//WxMpUpdateUserRemark 更新用户备注/标识名称,新的备注名,长度必须小于30字符
//https://developers.weixin.qq.com/doc/offiaccount/User_Management/Configuring_user_notes.html
func WxMpUpdateUserRemark(wxMpConfig IWxMpConfig, openId string, remark string) (*APIResult, error) {
	if len(openId) < 1 || len(remark) < 1 {
		return nil, errors.New("openId或者remark不能为空")
	}

	apiurl := WxMpAPIURL + "/cgi-bin/user/info/updateremark?access_token=" + wxMpConfig.GetAccessToken()

	params := make(map[string]interface{})
	params["openid"] = openId
	params["remark"] = remark

	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMpBlackUserOpenIdList 获取公众号的黑名单列表
//当 begin_openid 为空时，默认从开头拉取
//https://developers.weixin.qq.com/doc/offiaccount/User_Management/Manage_blacklist.html
func WxMpBlackUserOpenIdList(wxMpConfig IWxMpConfig, beginOpenid string) (*APIResult, error) {

	apiurl := WxMpAPIURL + "/cgi-bin/tags/members/getblacklist?access_token=" + wxMpConfig.GetAccessToken()

	params := make(map[string]interface{})
	if len(beginOpenid) > 0 {
		params["begin_openid"] = beginOpenid
	}

	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxBatchBlackUserOpenId  批量拉黑用户
//openIds 需要拉黑的用户openId列表
func WxMpBatchBlackUserOpenId(wxMpConfig IWxMpConfig, openIds []string) (*APIResult, error) {
	if len(openIds) < 1 {
		return nil, errors.New("需要拉黑的用户openId列表不能为空")
	}

	apiurl := WxMpAPIURL + "/cgi-bin/tags/members/batchblacklist?access_token=" + wxMpConfig.GetAccessToken()

	params := make(map[string]interface{})
	params["openid_list"] = openIds
	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxBatchUnBlackUserOpenId  批量解封拉黑的用户
//openIds 需要解封的用户openId列表
func WxMpBatchUnBlackUserOpenId(wxMpConfig IWxMpConfig, openIds []string) (*APIResult, error) {
	if len(openIds) < 1 {
		return nil, errors.New("需要拉黑的用户openId列表不能为空")
	}

	apiurl := WxMpAPIURL + "/cgi-bin/tags/members/batchunblacklist?access_token=" + wxMpConfig.GetAccessToken()

	params := make(map[string]interface{})
	params["openid_list"] = openIds
	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}
