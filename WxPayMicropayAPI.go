package gowe

import "encoding/xml"

//WxPayMicropay 提交付款码支付 https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_10&index=1
func WxPayMicropay(wxPayConfig IWxPayConfig, body *WxPayMicropayBody) (*WxPayMicropayResponse, error) {
	// 处理参数
	if body.SceneInfo != nil {
		body.SceneInfoStr = JsonString(body.SceneInfo)
	}
	// 业务逻辑
	bytes, err := wxPayDoWeChat(wxPayConfig, WxMpPayMchAPIURL+"/pay/micropay", body)
	if err != nil {
		return nil, err
	}
	// 结果校验
	if err = wxPayDoVerifySign(wxPayConfig, bytes, true); err != nil {
		return nil, err
	}
	// 解析返回值
	wxRsp := WxPayMicropayResponse{}
	err = xml.Unmarshal(bytes, &wxRsp)
	return &wxRsp, err
}

//WxPayMicropayBody 付款码支付的参数
type WxPayMicropayBody struct {
	SignType       string `json:"sign_type,omitempty"`   // 签名类型,目前支持HMAC-SHA256和MD5,默认为MD5
	DeviceInfo     string `json:"device_info,omitempty"` // 终端设备号(商户自定义,如门店编号)
	Body           string `json:"body"`                  // 商品或支付单简要描述,格式要求:门店品牌名-城市分店名-实际商品名称
	Detail         string `json:"detail,omitempty"`      // 单品优惠功能字段,需要接入请见详细说明
	Attach         string `json:"attach,omitempty"`      // 附加数据,在查询API和支付通知中原样返回,该字段主要用于商户携带订单的自定义数据
	OutTradeNo     string `json:"out_trade_no"`          // 商户系统内部订单号,要求32个字符内,只能是数字、大小写字母_-|*且在同一个商户号下唯一.详见商户订单号
	TotalFee       int    `json:"total_fee"`             // 订单总金额,单位为分,只能为整数,详见支付金额
	FeeType        string `json:"fee_type,omitempty"`    // 符合ISO 4217标准的三位字母代码,默认人民币:CNY,其他值列表详见货币类型
	SpbillCreateIP string `json:"spbill_create_ip"`      // 支持IPV4和IPV6两种格式的IP地址.调用微信支付API的机器IP
	GoodsTag       string `json:"goods_tag,omitempty"`   // 订单优惠标记,代金券或立减优惠功能的参数,说明详见代金券或立减优惠
	LimitPay       string `json:"limit_pay,omitempty"`   // no_credit:指定不能使用信用卡支付
	TimeStart      string `json:"time_start,omitempty"`  // 订单生成时间,格式为yyyyMMddHHmmss,如2009年12月25日9点10分10秒表示为20091225091010.其他详见时间规则
	TimeExpire     string `json:"time_expire,omitempty"` // 订单失效时间,格式为yyyyMMddHHmmss,如2009年12月27日9点10分10秒表示为20091227091010.注意:最短失效时间间隔需大于1分钟
	AuthCode       string `json:"auth_code"`             // 扫码支付授权码,设备读取用户微信中的条码或者二维码信息 (注:用户付款码条形码规则:18位纯数字,以10、11、12、13、14、15开头)
	Receipt        string `json:"receipt,omitempty"`     // Y,传入Y时,支付成功消息和支付详情页将出现开票入口.需要在微信支付商户平台或微信公众平台开通电子发票功能,传此字段才可生效
	SceneInfoStr   string `json:"scene_info,omitempty"`  // 该字段用于上报场景信息,目前支持上报实际门店信息.该字段为JSON对象数据,对象格式为{"store_info":{"id": "门店ID","name": "名称","area_code": "编码","address": "地址" }} ,字段详细说明请点击行前的+展开
	// 用于生成SceneInfoStr
	SceneInfo *WxPaySceneInfoModel `json:"-"`
}

//WxPayMicropayResponse 付款码支付的返回值
type WxPayMicropayResponse struct {
	WxResponseModel
	// 当return_code为SUCCESS时
	WxPayServiceResponseModel

	DeviceInfo string `xml:"device_info"` // 调用接口提交的终端设备号
	// 当return_code和result_code都为SUCCESS时
	OpenId             string `xml:"openid"`               // 用户在商户appid下的唯一标识
	IsSubscribe        string `xml:"is_subscribe"`         // 用户是否关注公众账号,仅在公众账号类型支付有效,取值范围:Y或N;Y-关注;N-未关注
	SubOpenId          string `xml:"sub_openid"`           // (服务商模式) 子商户appid下用户唯一标识,如需返回则请求时需要传sub_appid
	SubIsSubscribe     string `xml:"sub_is_subscribe"`     // (服务商模式) 用户是否关注子公众账号,仅在公众账号类型支付有效,取值范围:Y或N;Y-关注;N-未关注
	TradeType          string `xml:"trade_type"`           // 支付类型为MICROPAY(即扫码支付)
	BankType           string `xml:"bank_type"`            // 银行类型,采用字符串类型的银行标识,值列表详见银行类型
	FeeType            string `xml:"fee_type"`             // 符合ISO 4217标准的三位字母代码,默认人民币:CNY,其他值列表详见货币类型
	TotalFee           int    `xml:"total_fee"`            // 订单总金额,单位为分,只能为整数,详见支付金额
	CashFeeType        string `xml:"cash_fee_type"`        // 符合ISO 4217标准的三位字母代码,默认人民币:CNY,其他值列表详见货币类型
	CashFee            int    `xml:"cash_fee"`             // 订单现金支付金额,详见支付金额
	SettlementTotalFee int    `xml:"settlement_total_fee"` // 当订单使用了免充值型优惠券后返回该参数,应结订单金额=订单金额-免充值优惠券金额.
	CouponFee          int    `xml:"coupon_fee"`           // "代金券"金额<=订单金额,订单金额-"代金券"金额=现金支付金额,详见支付金额
	TransactionId      string `xml:"transaction_id"`       // 微信支付订单号
	OutTradeNo         string `xml:"out_trade_no"`         // 商户系统内部订单号,要求32个字符内,只能是数字、大小写字母_-|*且在同一个商户号下唯一.
	Attach             string `xml:"attach"`               // 商家数据包,原样返回
	TimeEnd            string `xml:"time_end"`             // 订单生成时间,格式为yyyyMMddHHmmss,如2009年12月25日9点10分10秒表示为20091225091010.详见时间规则
	PromotionDetail    string `xml:"promotion_detail"`     // TODO 单品优惠详情
}
