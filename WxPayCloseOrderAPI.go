package gowe

import (
	"context"
	"encoding/xml"
)

// WxPayCloseOrder 关闭订单 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_3
func WxPayCloseOrder(ctx context.Context, wxPayConfig IWxPayConfig, body *WxPayCloseOrderBody) (*WxPayCloseOrderResponse, error) {
	// 业务逻辑
	bytes, err := wxPayDoWeChat(ctx, wxPayConfig, "/pay/closeorder", body, 0)
	if err != nil {
		return nil, err
	}
	// 结果校验
	if err = wxPayDoVerifySign(ctx, wxPayConfig, bytes, true); err != nil {
		return nil, err
	}
	// 解析返回值
	res := &WxPayCloseOrderResponse{}
	err = xml.Unmarshal(bytes, res)
	return res, err
}

// WxPayCloseOrderBody 关闭订单的参数
type WxPayCloseOrderBody struct {
	SignType   string `json:"sign_type,omitempty"` // 签名类型,目前支持HMAC-SHA256和MD5,默认为MD5
	OutTradeNo string `json:"out_trade_no"`        // 商户系统内部订单号,要求32个字符内,只能是数字、大小写字母_-|*且在同一个商户号下唯一.详见商户订单号
}

// WxPayCloseOrderResponse 关闭订单的返回值
type WxPayCloseOrderResponse struct {
	WxResponseModel
	// 当return_code为SUCCESS时
	WxPayPartnerResponseModel
	ResultMsg string `xml:"result_msg"` // 对业务结果的补充说明
}
