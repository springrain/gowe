package gowe

import (
	"encoding/xml"

	"github.com/beevik/etree"
)

//查询支付订单 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2

//WxPayQueryOrder 查询订单
func WxPayQueryOrder(wxPayConfig IWxPayConfig, body *WxPayQueryOrderBody) (*WxPayQueryOrderResponse, error) {
	// 业务逻辑
	bytes, err := wxPayDoWeChat(wxPayConfig, WxPayMchAPIURL+"/pay/orderquery", body, 0)
	if err != nil {
		return nil, err
	}
	// 结果校验
	if err = wxPayDoVerifySign(wxPayConfig, bytes, true); err != nil {
		return nil, err
	}
	// 解析返回值
	res := &WxPayQueryOrderResponse{}
	err = wxPayQueryOrderParseResponse(bytes, res)
	return res, err
}

//WxPayQueryOrderBody 查询订单的参数
type WxPayQueryOrderBody struct {
	SignType      string `json:"sign_type,omitempty"`      // 签名类型,目前支持HMAC-SHA256和MD5,默认为MD5
	TransactionId string `json:"transaction_id,omitempty"` // (非必填,二选一) 微信的订单号,优先使用
	OutTradeNo    string `json:"out_trade_no,omitempty"`   // (非必填,二选一) 商户系统内部订单号,要求32个字符内,只能是数字、大小写字母_-|*且在同一个商户号下唯一.详见商户订单号
}

//WxPayQueryOrderResponse 查询订单的返回值
type WxPayQueryOrderResponse struct {
	WxResponseModel
	// 当return_code为SUCCESS时
	WxPayPartnerResponseModel
	// 当return_code、result_code、trade_state都为SUCCESS时有返回,如trade_state不为SUCCESS,则只返回out_trade_no(必传)和attach(选传).
	DeviceInfo         string `xml:"device_info"`          // 微信支付分配的终端设备号
	OpenId             string `xml:"openid"`               // 用户在商户appid下的唯一标识
	IsSubscribe        string `xml:"is_subscribe"`         // 用户是否关注公众账号,Y-关注,N-未关注(机构商户不返回)
	SubOpenId          string `xml:"sub_openid"`           // (服务商模式) 子商户appid下用户唯一标识,如需返回则请求时需要传sub_appid
	SubIsSubscribe     string `xml:"sub_is_subscribe"`     // (服务商模式) 用户是否关注子公众账号,仅在公众账号类型支付有效,取值范围:Y或N;Y-关注;N-未关注
	TradeType          string `xml:"trade_type"`           // 调用接口提交的交易类型,取值如下:JSAPI,NATIVE,APP,MICROPAY,详细说明见参数规定
	TradeState         string `xml:"trade_state"`          // SUCCESS—支付成功 REFUND—转入退款 NOTPAY—未支付 CLOSED—已关闭 REVOKED—已撤销(刷卡支付) USERPAYING--用户支付中 PAYERROR--支付失败(其他原因,如银行返回失败)
	BankType           string `xml:"bank_type"`            // 银行类型,采用字符串类型的银行标识
	Detail             string `xml:"detail"`               // 商品详细列表,使用Json格式,传输签名前请务必使用CDATA标签将JSON文本串保护起来.如果使用了单品优惠,会有单品优惠信息返回
	TotalFee           int    `xml:"total_fee"`            // 订单总金额,单位为分
	FeeType            string `xml:"fee_type"`             // 货币类型,符合ISO 4217标准的三位字母代码,默认人民币:CNY,其他值列表详见货币类型
	SettlementTotalFee int    `xml:"settlement_total_fee"` // 当订单使用了免充值型优惠券后返回该参数,应结订单金额=订单金额-免充值优惠券金额.
	CashFee            int    `xml:"cash_fee"`             // 现金支付金额订单现金支付金额,详见支付金额
	CashFeeType        string `xml:"cash_fee_type"`        // 货币类型,符合ISO 4217标准的三位字母代码,默认人民币:CNY,其他值列表详见货币类型
	CouponFee          int    `xml:"coupon_fee"`           // "代金券或立减优惠"金额<=订单总金额,订单总金额-"代金券或立减优惠"金额=现金支付金额,详见支付金额
	CouponCount        int    `xml:"coupon_count"`         // 代金券或立减优惠使用数量
	TransactionId      string `xml:"transaction_id"`       // 微信支付订单号
	OutTradeNo         string `xml:"out_trade_no"`         // 商户系统内部订单号,要求32个字符内,只能是数字、大小写字母_-|*@ ,且在同一个商户号下唯一.
	Attach             string `xml:"attach"`               // 商家数据包,原样返回
	TimeEnd            string `xml:"time_end"`             // 订单支付时间,格式为yyyyMMddHHmmss,如2009年12月25日9点10分10秒表示为20091225091010.其他详见时间规则
	TradeStateDesc     string `xml:"trade_state_desc"`     // 对当前查询订单状态的描述和下一步操作的指引
	// 使用coupon_count的序号生成的优惠券项
	Coupons []WxPayCouponResponseModel `xml:"-"`
}

//wxPayQueryOrderParseResponse 查询订单-解析XML返回值
func wxPayQueryOrderParseResponse(xmlStr []byte, rsp *WxPayQueryOrderResponse) (err error) {
	// 常规解析
	if err = xml.Unmarshal(xmlStr, rsp); err != nil {
		return
	}
	// 解析CouponCount的对应项
	if rsp.CouponCount > 0 {
		doc := etree.NewDocument()
		if err = doc.ReadFromBytes(xmlStr); err != nil {
			return
		}
		root := doc.SelectElement("xml")
		for i := 0; i < rsp.CouponCount; i++ {
			m := wxPayNewCouponResponseModel(root, "coupon_id_%d", "coupon_type_%d", "coupon_fee_%d", i)
			rsp.Coupons = append(rsp.Coupons, m)
		}
	}
	return
}
