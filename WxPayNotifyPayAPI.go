package gowe

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/beevik/etree"
)

type NotifyPayHandler func(WxPayNotifyPayBody) error

//WxPayNotifyPay 支付结果通知
func WxPayNotifyPay(wxPayConfig IWxPayConfig, handler NotifyPayHandler, requestBody []byte) (rspBody string, err error) {
	// 验证Sign
	if err = wxPayDoVerifySign(wxPayConfig, requestBody, false); err != nil {
		return
	}
	// 解析参数
	var body WxPayNotifyPayBody
	if err = wxPayNotifyParseParams(requestBody, &body); err != nil {
		return
	}
	// 调用外部处理
	if err = handler(body); err != nil {
		return
	}
	// 返回处理结果
	rspModel := WxPayNotifyPayBody{
		ReturnCode: ResponseSuccess,
		ReturnMsg:  ResponseMessageOk,
	}
	rspBody = rspModel.ToXmlString()
	return
}

// 支付结果通知的参数
type WxPayNotifyPayBody struct {
	ReturnCode string `xml:"return_code"` // SUCCESS/FAIL 此字段是通信标识，非交易标识，交易是否成功需要查看result_code来判断
	ReturnMsg  string `xml:"return_msg"`  // 返回信息，如非空，为错误原因：签名失败/参数格式校验错误
	RetMsg     string `xml:"retmsg"`      // 沙盒时返回的错误信息
	// 当return_code为SUCCESS时
	AppId      string `xml:"appid"`        // 微信分配的公众账号ID
	MchId      string `xml:"mch_id"`       // 微信支付分配的商户号
	SubAppId   string `xml:"sub_appid"`    // (服务商模式) 微信分配的子商户公众账号ID
	SubMchId   string `xml:"sub_mch_id"`   // (服务商模式) 微信支付分配的子商户号
	NonceStr   string `xml:"nonce_str"`    // 随机字符串，不长于32位
	Sign       string `xml:"sign"`         // 签名，详见签名生成算法
	ResultCode string `xml:"result_code"`  // SUCCESS/FAIL
	ErrCode    string `xml:"err_code"`     // 详细参见第6节错误列表
	ErrCodeDes string `xml:"err_code_des"` // 错误返回的信息描述

	DeviceInfo         string `xml:"device_info"`          // 微信支付分配的终端设备号
	IsSubscribe        string `xml:"is_subscribe"`         // 用户是否关注公众账号(机构商户不返回)
	SubIsSubscribe     string `xml:"sub_is_subscribe"`     // (服务商模式) 用户是否关注子公众账号(机构商户不返回)
	OpenId             string `xml:"openid"`               // 用户在商户appid下的唯一标识
	SubOpenId          string `xml:"sub_openid"`           // (服务商模式) 用户在子商户appid下的唯一标识
	TradeType          string `xml:"trade_type"`           // 交易类型
	BankType           string `xml:"bank_type"`            // 银行类型，采用字符串类型的银行标识，银行类型见附表
	TotalFee           int    `xml:"total_fee"`            // 订单总金额，单位为分
	FeeType            string `xml:"fee_type"`             // 货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CashFee            int    `xml:"cash_fee"`             // 现金支付金额订单现金支付金额，详见支付金额
	CashFeeType        string `xml:"cash_fee_type"`        // 货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	SettlementTotalFee int    `xml:"settlement_total_fee"` // 应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	CouponFee          int    `xml:"coupon_fee"`           // 代金券或立减优惠金额<=订单总金额，订单总金额-代金券或立减优惠金额=现金支付金额，详见支付金额
	CouponCount        int    `xml:"coupon_count"`         // 代金券或立减优惠使用数量
	TransactionId      string `xml:"transaction_id"`       // 微信支付订单号
	OutTradeNo         string `xml:"out_trade_no"`         // 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	Attach             string `xml:"attach"`               // 商家数据包，原样返回
	TimeEnd            string `xml:"time_end"`             // 支付完成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
	// 使用coupon_count的序号生成的优惠券项
	Coupons []CouponResponseModel `xml:"-"`
}

func (m *WxPayNotifyPayBody) ToXmlString() string {
	buffer := new(bytes.Buffer)
	buffer.WriteString("<xml>")
	buffer.WriteString(fmt.Sprintf("<return_code><![CDATA[%s]]></return_code>", m.ReturnCode))
	buffer.WriteString(fmt.Sprintf("<return_msg><![CDATA[%s]]></return_msg>", m.ReturnMsg))
	buffer.WriteString("</xml>")
	return buffer.String()
}

// 支付结果通知-解析XML参数
func wxPayNotifyParseParams(xmlStr []byte, body *WxPayNotifyPayBody) (err error) {
	if err = xml.Unmarshal(xmlStr, &body); err != nil {
		return
	}
	// 解析CouponCount的对应项
	if body.CouponCount > 0 {
		doc := etree.NewDocument()
		if err = doc.ReadFromBytes(xmlStr); err != nil {
			return
		}
		root := doc.SelectElement("xml")
		for i := 0; i < body.CouponCount; i++ {
			m := NewCouponResponseModel(root, "coupon_id_%d", "coupon_type_%d", "coupon_fee_%d", i)
			body.Coupons = append(body.Coupons, m)
		}
	}
	return
}

// 返回结果中的优惠券条目信息
type CouponResponseModel struct {
	CouponId   string // 代金券或立减优惠ID
	CouponType string // CASH-充值代金券 NO_CASH-非充值优惠券 开通免充值券功能，并且订单使用了优惠券后有返回
	CouponFee  int64  // 单个代金券或立减优惠支付金额
}

// 在XML节点树中，查找labels对应的
func NewCouponResponseModel(
	doc *etree.Element,
	idFormat string,
	typeFormat string,
	feeFormat string,
	numbers ...interface{},
) (m CouponResponseModel) {
	idName := fmt.Sprintf(idFormat, numbers...)
	typeName := fmt.Sprintf(typeFormat, numbers...)
	feeName := fmt.Sprintf(feeFormat, numbers...)
	m.CouponId = doc.SelectElement(idName).Text()
	m.CouponType = doc.SelectElement(typeName).Text()
	m.CouponFee, _ = strconv.ParseInt(doc.SelectElement(feeName).Text(), 10, 64)
	return
}
