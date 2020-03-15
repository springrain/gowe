package gowe

import (
	"encoding/json"
	"errors"
)

//WxMpUserInfo 用户信息
type WxMpUserInfo struct {
	OpenId         string   `json:"openid"`          // 用户唯一标识
	Nickname       string   `json:"nickname"`        // 用户的昵称
	Sex            int      `json:"sex"`             // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Language       string   `json:"language"`        // 用户的语言，简体中文为zh_CN
	Province       string   `json:"province"`        // 用户所在省份
	City           string   `json:"city"`            // 用户所在城市
	Country        string   `json:"country"`         // 用户所在国家
	HeadimgUrl     string   `json:"headimgurl"`      // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	Privilege      []string `json:"privilege"`       // 用户特权信息
	UnionId        string   `json:"unionid"`         // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	Subscribe      int      `json:"subscribe"`       // 用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。
	SubscribeTime  int      `json:"subscribe_time"`  // 用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
	SubscribeScene string   `json:"subscribe_scene"` // 返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENEPROFILE LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_OTHERS 其他
	Remark         string   `json:"remark"`          // 公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	GroupId        int      `json:"groupid"`         // 用户所在的分组ID（兼容旧的用户分组接口）
	TagidList      []int    `json:"tagid_list"`      // 用户被打上的标签ID列表
	QrScene        int      `json:"qr_scene"`        // 二维码扫码场景（开发者自定义）
	QrSceneStr     string   `json:"qr_scene_str"`    // 二维码扫码场景描述（开发者自定义）
	ErrCode        int      `json:"errcode"`         // 错误码
	ErrMsg         string   `json:"errmsg"`          // 错误信息
}

//WxMpGetUserInfo 获取用户基本信息,包括UnionID机制(授权机制)
//https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html
func WxMpGetUserInfo(wxMpConfig IWxMpConfig, openId, lang string) (*WxMpUserInfo, error) {
	if len(lang) <= 0 {
		lang = "zh_CN"
	}
	apiurl := WxMpAPIURL + "/sns/userinfo?access_token=" + wxMpConfig.GetAccessToken() + "&openid=" + openId + "s&lang=" + lang
	body, err := httpGet(apiurl)
	if err != nil {
		return nil, err
	}
	userInfo := WxMpUserInfo{}
	err = json.Unmarshal(body, &userInfo)
	return &userInfo, err
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
