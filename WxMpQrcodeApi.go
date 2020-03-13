package gowe

import (
	"context"
	"errors"
)

//WxMpCreateTemporary 生成带参数的二维码 API
//https://developers.weixin.qq.com/doc/offiaccount/Account_Management/Generating_a_Parametric_QR_Code.html
//expireSeconds:该二维码有效时间,以秒为单位. 最大不超过2592000(即30天),默认2592000
//sceneStr:场景值ID(字符串形式的ID),字符串类型,长度限制为1到64.这里只使用字符串了,用途更广法
func WxMpCreateTemporary(ctx context.Context, sceneStr string, expireSeconds int) (*APIResult, error) {
	if len(sceneStr) < 1 {
		return nil, errors.New("expireSecond或sceneStr不能为空")
	}
	if expireSeconds == 0 {
		expireSeconds = 2592000
	}
	wxMpConfig, errWxMpConfig := getWxMpConfig(ctx)
	if errWxMpConfig != nil {
		return nil, errWxMpConfig
	}
	apiurl := WxmpApiUrl + "/cgi-bin/qrcode/create?access_token=" + wxMpConfig.AccessToken

	params := make(map[string]interface{})
	params["expire_seconds"] = expireSeconds
	params["action_name"] = "QR_STR_SCENE"
	actionInfo := make(map[string]interface{})
	scene := make(map[string]interface{})
	scene["scene_str"] = sceneStr
	actionInfo["scene"] = scene
	params["action_info"] = actionInfo
	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}
