package gowe

//WxMpAutoreply 获取自动回复
//https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Getting_Rules_for_Auto_Replies.html
func WxMpAutoreply(wxMpConfig IWxMpConfig) (*APIResult, error) {

	apiurl := WxMpAPIURL + "/cgi-bin/get_current_autoreply_info?access_token=" + wxMpConfig.GetAccessToken()

	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}
