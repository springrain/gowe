package gowe

import (
	"encoding/json"
)

//https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html

//WxMaQrCode 微信二维码的实体类
type WxMaQrCode struct {
	qrCodeMap    map[string]interface{}
	lineColorMap map[string]int
}

//SetScene 最大32个可见字符,只支持数字,大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~,其它字符请自行编码为合法字符(因不支持%,中文无法使用 urlencode 处理,请使用其他编码方式)
func (wxMaQrCode *WxMaQrCode) SetScene(value string) {
	if wxMaQrCode.qrCodeMap == nil {
		wxMaQrCode.qrCodeMap = make(map[string]interface{})
	}
	wxMaQrCode.qrCodeMap["scene"] = value
}

//SetPage 必须是已经发布的小程序存在的页面(否则报错),例如 pages/index/index, 根路径前不要填加 /,不能携带参数(参数请放在scene字段里),如果不填写这个字段,默认跳主页面
func (wxMaQrCode *WxMaQrCode) SetPage(value string) {
	if wxMaQrCode.qrCodeMap == nil {
		wxMaQrCode.qrCodeMap = make(map[string]interface{})
	}
	wxMaQrCode.qrCodeMap["page"] = value
}

//SetWidth 二维码的宽度,单位 px,最小 280px,最大 1280px
func (wxMaQrCode *WxMaQrCode) SetWidth(value int) {
	if wxMaQrCode.qrCodeMap == nil {
		wxMaQrCode.qrCodeMap = make(map[string]interface{})
	}
	wxMaQrCode.qrCodeMap["width"] = value
}

//SetAutoColor 自动配置线条颜色,如果颜色依然是黑色,则说明不建议配置主色调,默认 false
func (wxMaQrCode *WxMaQrCode) SetAutoColor(value bool) {
	if wxMaQrCode.qrCodeMap == nil {
		wxMaQrCode.qrCodeMap = make(map[string]interface{})
	}
	wxMaQrCode.qrCodeMap["auto_color"] = value
}

//SetHyaline 是否需要透明底色,为 true 时,生成透明底色的小程序
func (wxMaQrCode *WxMaQrCode) SetHyaline(value bool) {
	if wxMaQrCode.qrCodeMap == nil {
		wxMaQrCode.qrCodeMap = make(map[string]interface{})
	}
	wxMaQrCode.qrCodeMap["is_hyaline"] = value
}

//SetLineColorR auto_color 为 false 时生效,使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"} 十进制表示
func (wxMaQrCode *WxMaQrCode) SetLineColorR(value int) {
	if wxMaQrCode.lineColorMap == nil {
		wxMaQrCode.lineColorMap = make(map[string]int)
	}
	wxMaQrCode.lineColorMap["r"] = value
}

//SetLineColorG auto_color 为 false 时生效,使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"} 十进制表示
func (wxMaQrCode *WxMaQrCode) SetLineColorG(value int) {
	if wxMaQrCode.lineColorMap == nil {
		wxMaQrCode.lineColorMap = make(map[string]int)
	}
	wxMaQrCode.lineColorMap["g"] = value
}

//SetLineColorB auto_color 为 false 时生效,使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"} 十进制表示
func (wxMaQrCode *WxMaQrCode) SetLineColorB(value int) {
	if wxMaQrCode.lineColorMap == nil {
		wxMaQrCode.lineColorMap = make(map[string]int)
	}
	wxMaQrCode.lineColorMap["b"] = value
}

//getQrCodeMap 获取到小程序码的map信息,作为请求参数
func (wxMaQrCode *WxMaQrCode) getQrCodeMap() map[string]interface{} {
	if wxMaQrCode.lineColorMap != nil {
		wxMaQrCode.qrCodeMap["line_color"] = wxMaQrCode.lineColorMap
	}

	return wxMaQrCode.qrCodeMap

}

//WxMaCodeGetUnlimited 小程序码接口
//https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html
func WxMaCodeGetUnlimited(wxMaConfig IWxMaConfig, wxMaQrCode *WxMaQrCode) (*APIResult, error) {

	apiurl := WxMpAPIURL + "/wxa/getwxacodeunlimit?access_token=" + wxMaConfig.GetAccessToken()
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
