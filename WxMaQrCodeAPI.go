package gowe

import (
	"context"
	"encoding/json"
)

//WxMaQrCode 微信二维码的实体类
type WxMaQrCode struct {
	qrCodeMap    map[string]interface{}
	lineColorMap map[string]int
}

func (wxMaQrCode *WxMaQrCode) SetScene(value string) {
	if wxMaQrCode.qrCodeMap == nil {
		wxMaQrCode.qrCodeMap = make(map[string]interface{})
	}
	wxMaQrCode.qrCodeMap["scene"] = value
}

func (wxMaQrCode *WxMaQrCode) SetPage(value string) {
	if wxMaQrCode.qrCodeMap == nil {
		wxMaQrCode.qrCodeMap = make(map[string]interface{})
	}
	wxMaQrCode.qrCodeMap["page"] = value
}

func (wxMaQrCode *WxMaQrCode) SetWidth(value int) {
	if wxMaQrCode.qrCodeMap == nil {
		wxMaQrCode.qrCodeMap = make(map[string]interface{})
	}
	wxMaQrCode.qrCodeMap["width"] = value
}
func (wxMaQrCode *WxMaQrCode) SetAutoColor(value bool) {
	if wxMaQrCode.qrCodeMap == nil {
		wxMaQrCode.qrCodeMap = make(map[string]interface{})
	}
	wxMaQrCode.qrCodeMap["auto_color"] = value
}
func (wxMaQrCode *WxMaQrCode) SetHyaline(value bool) {
	if wxMaQrCode.qrCodeMap == nil {
		wxMaQrCode.qrCodeMap = make(map[string]interface{})
	}
	wxMaQrCode.qrCodeMap["is_hyaline"] = value
}
func (wxMaQrCode *WxMaQrCode) SetLineColorR(value int) {
	if wxMaQrCode.lineColorMap == nil {
		wxMaQrCode.lineColorMap = make(map[string]int)
	}
	wxMaQrCode.lineColorMap["r"] = value
}
func (wxMaQrCode *WxMaQrCode) SetLineColorG(value int) {
	if wxMaQrCode.lineColorMap == nil {
		wxMaQrCode.lineColorMap = make(map[string]int)
	}
	wxMaQrCode.lineColorMap["g"] = value
}
func (wxMaQrCode *WxMaQrCode) SetLineColorB(value int) {
	if wxMaQrCode.lineColorMap == nil {
		wxMaQrCode.lineColorMap = make(map[string]int)
	}
	wxMaQrCode.lineColorMap["b"] = value
}
func (wxMaQrCode *WxMaQrCode) getQrCodeMap() map[string]interface{} {
	if wxMaQrCode.lineColorMap != nil {
		wxMaQrCode.qrCodeMap["line_color"] = wxMaQrCode.lineColorMap
	}

	return wxMaQrCode.qrCodeMap

}

//WxMaCodeGetUnlimited 小程序码接口
//https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html
func WxMaCodeGetUnlimited(ctx context.Context, wxMaQrCode *WxMaQrCode) (*APIResult, error) {
	wxMaConfig, errWxMaConfig := getWxMaConfig(ctx)
	if errWxMaConfig != nil {
		return nil, errWxMaConfig
	}
	apiurl := WxmpApiUrl + "/wxa/getwxacodeunlimit?access_token=" + wxMaConfig.AccessToken
	apiResult := APIResult{}
	data, errPost := httpPost(apiurl, wxMaQrCode.getQrCodeMap())
	if errPost == nil { //正常
		apiResult.FileData = data
		return &apiResult, nil
	}

	// 尝试解码
	resultMap := make(map[string]interface{})
	_ = json.Unmarshal(data, &resultMap)
	apiResult.ResultMap = resultMap
	return &apiResult, errPost
}
