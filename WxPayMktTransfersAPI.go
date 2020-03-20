package gowe

import "encoding/xml"

// 企业付款到零钱
//https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_1

//WxPayPromotionMktTransfers 企业付款到零钱(前提用户必须关注公众号)
func WxPayPromotionMktTransfers(wxPayConfig IWxPayConfig, body WxPayPromotionMktTransfersBody) (wxRsp WxPayPromotionMktTransfersResponse, err error) {
	// 业务逻辑
	bytes, err := wxPayDoWeChatWithCert(wxPayConfig, WxMpPayMchAPIURL+"/mmpaymkttransfers/promotion/transfers", body)
	if err != nil {
		return
	}
	// 不返回sign不需要校验
	// 解析返回值
	err = xml.Unmarshal(bytes, &wxRsp)
	return
}

//WxPayQueryMktTransfer 企业付款到零钱的查询
func WxPayQueryMktTransfer(wxPayConfig IWxPayConfig, body WxPayQueryMktTransferBody) (wxRsp WxPayQueryMktTransferResponse, err error) {
	// 业务逻辑
	bytes, err := wxPayDoWeChatWithCert(wxPayConfig, WxMpPayMchAPIURL+"/mmpaymkttransfers/gettransferinfo", body)
	if err != nil {
		return
	}
	// 不返回sign不需要校验
	// 解析返回值
	err = xml.Unmarshal(bytes, &wxRsp)
	return
}

// 微信找零的参数
type WxPayPromotionMktTransfersBody struct {
	DeviceInfo     string `json:"device_info,omitempty"`  // 终端设备号
	PartnerTradeNo string `json:"partner_trade_no"`       // 商户系统内部订单号,要求32个字符内,只能是数字、大小写字母_-|*且在同一个商户号下唯一.详见商户订单号
	OpenId         string `json:"openid"`                 // 商品或支付单简要描述,格式要求：门店品牌名-城市分店名-实际商品名称
	CheckName      string `json:"check_name"`             // 校验用户姓名选项(见constant定义)
	ReUserName     string `json:"re_user_name,omitempty"` // 收款用户姓名
	Amount         int    `json:"amount"`                 // 企业找零金额,单位为分
	Desc           string `json:"desc"`                   // 企业付款备注
	SpbillCreateIP string `json:"spbill_create_ip"`       // IP可传用户端或者服务端的IP
}

//WxPayPromotionMktTransfersResponse 微信找零的返回值
type WxPayPromotionMktTransfersResponse struct {
	WxResponseModel
	WxPayMchServiceResponseModel
	DeviceInfo     string `xml:"device_info"`      // 终端设备号
	PartnerTradeNo string `xml:"partner_trade_no"` // 商户系统内部订单号,要求32个字符内,只能是数字、大小写字母_-|*且在同一个商户号下唯一.详见商户订单号
	PaymentNo      string `xml:"payment_no"`       // 企业付款成功,返回的微信付款单号
	PaymentTime    string `xml:"payment_time"`     // 企业付款成功时间
}

// 微信找零查询的参数
type WxPayQueryMktTransferBody struct {
	PartnerTradeNo string `json:"partner_trade_no"` // 商户系统内部订单号
}

// 微信找零查询的返回值
type WxPayQueryMktTransferResponse struct {
	WxResponseModel
	WxPayMchServiceResponseModel
	PartnerTradeNo string `xml:"partner_trade_no"` // 商户系统内部订单号,要求32个字符内,只能是数字、大小写字母_-|*且在同一个商户号下唯一.详见商户订单号
	DetailId       string `xml:"detail_id"`        // 调用企业付款API时,微信系统内部产生的单号
	Status         string `xml:"status"`           // 转账状态
	Reason         string `xml:"reason"`           // 失败原因
	OpenId         string `xml:"openid"`           // 转账的openid
	TransferName   string `xml:"transfer_name"`    // 收款用户姓名
	PaymentAmount  int64  `xml:"payment_amount"`   // 付款金额单位为“分”
	TransferTime   string `xml:"transfer_time"`    // 发起转账的时间
	PaymentTime    string `xml:"payment_time"`     // 企业付款成功时间
	Desc           string `xml:"desc"`             // 企业付款备注
}
