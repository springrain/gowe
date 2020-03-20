package gowe

import (
	"encoding/json"
	"errors"
)

//WxMpQrCreateTemporary 生成带参数的临时二维码 API
//https://developers.weixin.qq.com/doc/offiaccount/Account_Management/Generating_a_Parametric_QR_Code.html
//expireSeconds:该二维码有效时间,以秒为单位. 最大不超过2592000(即30天),默认2592000
//sceneStr:场景值ID(字符串形式的ID),字符串类型,长度限制为1到64.这里只使用字符串了,用途更广
func WxMpQrCreateTemporary(wxMpConfig IWxMpConfig, sceneStr string, expireSeconds int) (wxqrr WxMpQrCreateResponse, err error) {
	if len(sceneStr) < 1 {
		return wxqrr, errors.New("sceneStr不能为空")
	}
	if expireSeconds == 0 {
		expireSeconds = 2592000
	}
	apiurl := WxMpAPIURL + "/cgi-bin/qrcode/create?access_token=" + wxMpConfig.GetAccessToken()

	params := make(map[string]interface{})
	params["expire_seconds"] = expireSeconds
	params["action_name"] = "QR_STR_SCENE"
	actionInfo := make(map[string]interface{})
	scene := make(map[string]interface{})
	scene["scene_str"] = sceneStr
	actionInfo["scene"] = scene
	params["action_info"] = actionInfo
	data, err := httpPost(apiurl, params)
	// 发送请求
	if err != nil {
		return wxqrr, err
	}
	// 尝试解码
	_ = json.Unmarshal(data, &wxqrr)

	return wxqrr, nil
}

//WxMpQrCreatePermanent 创建永久的带参数二维码
func WxMpQrCreatePermanent(wxMpConfig IWxMpConfig, sceneStr string) (wxqrr WxMpQrCreateResponse, err error) {
	if len(sceneStr) < 1 {
		return wxqrr, errors.New("sceneStr不能为空")
	}

	apiurl := WxMpAPIURL + "/cgi-bin/qrcode/create?access_token=" + wxMpConfig.GetAccessToken()

	params := make(map[string]interface{})
	params["action_name"] = "QR_LIMIT_STR_SCENE"
	actionInfo := make(map[string]interface{})
	scene := make(map[string]interface{})
	scene["scene_str"] = sceneStr
	actionInfo["scene"] = scene
	params["action_info"] = actionInfo

	data, err := httpPost(apiurl, params)
	// 发送请求
	if err != nil {
		return wxqrr, err
	}
	// 尝试解码
	_ = json.Unmarshal(data, &wxqrr)

	return wxqrr, nil
}

//WxMpQrShowQrCodeUrl 通过ticket换取二维码地址
func WxMpQrShowQrCodeUrl(wxMpConfig IWxMpConfig, ticket string) (string, error) {
	if len(ticket) < 1 {
		return "", errors.New("ticket不能为空")
	}
	return WxMpAPIURL + "/cgi-bin/showqrcode?ticket=" + ticket, nil
}

type WxMpQrCreateResponse struct {
	Ticket        string `json:"ticket"`         // 获取的二维码ticket，凭借此ticket可以在有效时间内换取二维码。
	ExpireSeconds int    `json:"expire_seconds"` // 该二维码有效时间，以秒为单位。 最大不超过2592000（即30天）
	URL           string `json:"url"`            //二维码图片解析后的地址，开发者可根据该地址自行生成需要的二维码图片
}
