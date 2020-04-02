package gowe

import "encoding/xml"

//WxPayReverseMicroOrder 撤销付款码订单 https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_11&index=3
func WxPayReverseMicroOrder(wxPayConfig IWxPayConfig, body *WxPayReverseMicroOrderBody) (*WxPayReverseMicroOrderResponse, error) {
	// 业务逻辑
	bytes, err := wxPayDoWeChatWithCert(wxPayConfig, "/secapi/pay/reverse", body, 0)
	if err != nil {
		return nil, err
	}
	// 结果校验
	if err = wxPayDoVerifySign(wxPayConfig, bytes, true); err != nil {
		return nil, err
	}
	// 解析返回值
	res := &WxPayReverseMicroOrderResponse{}
	err = xml.Unmarshal(bytes, res)
	return res, err
}

//WxPayReverseMicroOrderBody 撤销订单的参数
type WxPayReverseMicroOrderBody struct {
	TransactionId string `json:"transaction_id,omitempty"` // 微信支付订单号
	OutTradeNo    string `json:"out_trade_no"`             // 商户系统内部订单号,要求32个字符内,只能是数字、大小写字母_-|*@ ,且在同一个商户号下唯一.
}

//WxPayReverseMicroOrderResponse 撤销订单的返回值
type WxPayReverseMicroOrderResponse struct {
	WxResponseModel
	WxPayPartnerResponseModel
	Recall string `xml:"recall"` // 是否需要继续调用撤销,Y-需要,N-不需要
}
