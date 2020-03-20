package gowe

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/beevik/etree"
)

//退款查询 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_5

//WxPayQueryRefund 查询退款
func WxPayQueryRefund(wxPayConfig IWxPayConfig, body *WxPayQueryRefundBody) (*WxPayQueryRefundResponse, error) {
	// 业务逻辑
	bytes, err := wxPayDoWeChat(wxPayConfig, WxPayMchAPIURL+"/pay/refundquery", body)
	if err != nil {
		return nil, err
	}
	// 结果校验
	if err = wxPayDoVerifySign(wxPayConfig, bytes, true); err != nil {
		return nil, err
	}
	// 解析返回值
	wxRsp := WxPayQueryRefundResponse{}
	err = wxPayQueryRefundParseResponse(bytes, &wxRsp)
	return &wxRsp, err
}

//WxPayQueryRefundBody 查询退款的参数
type WxPayQueryRefundBody struct {
	SignType      string `json:"sign_type,omitempty"`      // 签名类型,目前支持HMAC-SHA256和MD5,默认为MD5
	TransactionId string `json:"transaction_id,omitempty"` // (非必填,四选一) 微信订单号 查询的优先级是: refund_id > out_refund_no > transaction_id > out_trade_no
	OutTradeNo    string `json:"out_trade_no,omitempty"`   // (非必填,四选一) 商户系统内部订单号,要求32个字符内,只能是数字、大小写字母_-|*@ ,且在同一个商户号下唯一.
	OutRefundNo   string `json:"out_refund_no,omitempty"`  // (非必填,四选一) 商户系统内部的退款单号,商户系统内部唯一,只能是数字、大小写字母_-|*@ ,同一退款单号多次请求只退一笔.
	RefundId      string `json:"refund_id,omitempty"`      // (非必填,四选一) 微信退款单号
	Offset        string `json:"offset,omitempty"`         // (非必填) 偏移量,当部分退款次数超过10次时可使用,表示返回的查询结果从这个偏移量开始取记录
}

//WxPayQueryRefundResponse 查询退款的返回值
type WxPayQueryRefundResponse struct {
	WxResponseModel
	// 当return_code为SUCCESS时
	WxPayServiceResponseModel
	TransactionId      string `xml:"transaction_id"`       // 微信订单号
	OutTradeNo         string `xml:"out_trade_no"`         // 商户系统内部订单号,要求32个字符内,只能是数字、大小写字母_-|*@ ,且在同一个商户号下唯一.
	TotalFee           int    `xml:"total_fee"`            // 订单总金额,单位为分,只能为整数,详见支付金额
	SettlementTotalFee int    `xml:"settlement_total_fee"` // 当订单使用了免充值型优惠券后返回该参数,应结订单金额=订单金额-免充值优惠券金额.
	FeeType            string `xml:"fee_type"`             // 订单金额货币类型,符合ISO 4217标准的三位字母代码,默认人民币:CNY,其他值列表详见货币类型
	CashFee            int    `xml:"cash_fee"`             // 现金支付金额,单位为分,只能为整数,详见支付金额
	RefundCount        int    `xml:"refund_count"`         // 当前返回退款笔数
	TotalRefundCount   int    `xml:"total_refund_count"`   // 订单总共已发生的部分退款次数,当请求参数传入offset后有返回
	// 使用refund_count的序号生成的当前退款项
	CurrentRefunds []WxPayQueryRefundResponseCurrentRefund `xml:"-"`
	// 使用total_refund_count的序号生成的总退款项
	TotalRefunds []WxPayQueryRefundResponseTotalRefund `xml:"-"`
}

//WxPayQueryRefundResponseCurrentRefund 使用refund_count的序号生成的当前退款项
type WxPayQueryRefundResponseCurrentRefund struct {
	OutRefundNo   string // 商户系统内部的退款单号,商户系统内部唯一,只能是数字、大小写字母_-|*@ ,同一退款单号多次请求只退一笔.
	RefundId      string // 微信退款单号
	RefundChannel string // ORIGINAL—原路退款 BALANCE—退回到余额 OTHER_BALANCE—原账户异常退到其他余额账户 OTHER_BANKCARD—原银行卡异常退到其他银行卡
}

//WxPayQueryRefundResponseTotalRefund 使用total_refund_count的序号生成的总退款项
type WxPayQueryRefundResponseTotalRefund struct {
	RefundFee           int64  // 退款总金额,单位为分,可以做部分退款
	SettlementRefundFee int64  // 退款金额=申请退款金额-非充值代金券退款金额,退款金额<=申请退款金额
	CouponRefundFee     int64  // 代金券退款金额<=退款金额,退款金额-代金券或立减优惠退款金额为现金,说明详见代金券或立减优惠
	CouponRefundCount   int64  // 退款代金券使用数量 ,$n为下标,从0开始编号
	RefundStatus        string // 退款状态:SUCCESS—退款成功 REFUNDCLOSE—退款关闭 PROCESSING—退款处理中 CHANGE—退款异常,退款到银行发现用户的卡作废或者冻结了,导致原路退款银行卡失败,可前往商户平台(pay.weixin.qq.com)-交易中心,手动处理此笔退款.$n为下标,从0开始编号.
	RefundAccount       string // REFUND_SOURCE_RECHARGE_FUNDS---可用余额退款/基本账户 REFUND_SOURCE_UNSETTLED_FUNDS---未结算资金退款 $n为下标,从0开始编号.
	RefundRecvAccout    string // 取当前退款单的退款入账方 1)退回银行卡:{银行名称}{卡类型}{卡尾号} 2)退回支付用户零钱: 支付用户零钱 3)退还商户: 商户基本账户 商户结算银行账户 4)退回支付用户零钱通: 支付用户零钱通
	RefundSuccessTime   string // 退款成功时间,当退款状态为退款成功时有返回.$n为下标,从0开始编号.
	// 使用coupon_refund_count的序号生成的代金券项
	Coupons []WxPayCouponResponseModel
}

//wxPayQueryRefundParseResponse 查询退款-解析返回值
func wxPayQueryRefundParseResponse(xmlStr []byte, rsp *WxPayQueryRefundResponse) (err error) {
	// 常规解析
	if err = xml.Unmarshal(xmlStr, rsp); err != nil {
		return
	}
	// 解析RefundCount和TotalRefundCount的对应项
	if rsp.RefundCount > 0 || rsp.TotalRefundCount > 0 {
		doc := etree.NewDocument()
		if err = doc.ReadFromBytes(xmlStr); err != nil {
			return
		}
		root := doc.SelectElement("xml")
		// 解析RefundCount的对应项
		if rsp.RefundCount > 0 {
			rsp.CurrentRefunds = make([]WxPayQueryRefundResponseCurrentRefund, rsp.RefundCount)
			for i := 0; i < rsp.RefundCount; i++ {
				rsp.CurrentRefunds[i].OutRefundNo = root.SelectElement(fmt.Sprintf("out_refund_no_%d", i)).Text()
				rsp.CurrentRefunds[i].RefundId = root.SelectElement(fmt.Sprintf("refund_id_%d", i)).Text()
				rsp.CurrentRefunds[i].RefundChannel = root.SelectElement(fmt.Sprintf("refund_channel_%d", i)).Text()
			}
		}
		// 解析TotalRefundCount的对应项
		if rsp.TotalRefundCount > 0 {
			rsp.TotalRefunds = make([]WxPayQueryRefundResponseTotalRefund, rsp.TotalRefundCount)
			for i := 0; i < rsp.TotalRefundCount; i++ {
				rsp.TotalRefunds[i].RefundFee, _ = strconv.ParseInt(root.SelectElement(fmt.Sprintf("refund_fee_%d", i)).Text(), 10, 64)
				rsp.TotalRefunds[i].SettlementRefundFee, _ = strconv.ParseInt(root.SelectElement(fmt.Sprintf("settlement_refund_fee_%d", i)).Text(), 10, 64)
				rsp.TotalRefunds[i].CouponRefundFee, _ = strconv.ParseInt(root.SelectElement(fmt.Sprintf("coupon_refund_fee_%d", i)).Text(), 10, 64)
				rsp.TotalRefunds[i].CouponRefundCount, _ = strconv.ParseInt(root.SelectElement(fmt.Sprintf("coupon_refund_count_%d", i)).Text(), 10, 64)
				rsp.TotalRefunds[i].RefundStatus = root.SelectElement(fmt.Sprintf("refund_status_%d", i)).Text()
				rsp.TotalRefunds[i].RefundAccount = root.SelectElement(fmt.Sprintf("refund_account_%d", i)).Text()
				rsp.TotalRefunds[i].RefundRecvAccout = root.SelectElement(fmt.Sprintf("refund_recv_accout_%d", i)).Text()
				rsp.TotalRefunds[i].RefundSuccessTime = root.SelectElement(fmt.Sprintf("refund_success_time_%d", i)).Text()
				if rsp.TotalRefunds[i].CouponRefundCount > 0 {
					for j := int64(0); j < rsp.TotalRefunds[i].CouponRefundCount; j++ {
						m := WxPayNewCouponResponseModel(root, "coupon_refund_id_%d_%d", "coupon_type_%d_%d", "coupon_refund_fee_%d_%d", i, j)
						rsp.TotalRefunds[i].Coupons = append(rsp.TotalRefunds[i].Coupons, m)
					}
				}
			}
		}
	}
	return
}
