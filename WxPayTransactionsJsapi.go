package gowe

import (
	"context"
	"fmt"
)

// 请求结构体
type JsapiOrderRequest struct {
	AppID       string `json:"appid"`
	MchID       string `json:"mchid"`
	Description string `json:"description"`
	OutTradeNo  string `json:"out_trade_no"`
	TimeExpire  string `json:"time_expire,omitempty"`
	NotifyURL   string `json:"notify_url"`
	Amount      struct {
		Total    int    `json:"total"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Payer struct {
		OpenID string `json:"openid"`
	} `json:"payer"`
}

// 响应结构体
type JsapiOrderResponse struct {
	PrepayID string `json:"prepay_id"`
}

// 错误响应
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// 生成JSAPI支付参数
type JsapiPayParams struct {
	AppID     string `json:"appId"`
	Timestamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

// 下单入口
// 参数说明:
//   - ctx: 上下文对象，用于控制请求生命周期（如超时、取消等）
//   - wxPayConfig: 微信支付配置接口，需包含以下字段:
//   - openid: 用户的openid，通过微信授权获取（如: oUpF8uMuAJO_M2pxb1Q9zNjWeS6o）
//   - totalFee: 订单总金额（单位: 分，如100表示1元）
//   - description: 商品描述（如: "腾讯充值中心-QQ会员充值"）
//
// 返回值:
//   - JsapiPayParams: JSAPI调起支付所需参数，包含:
//   - AppID: 公众号/小程序ID
//   - TimeStamp: 时间戳（秒级）
//   - NonceStr: 随机字符串
//   - Package: 预支付交易会话标识（如: prepay_id=wx201410272009395522657a690389285100）
//   - SignType: 签名类型（固定为"RSA"）
//   - PaySign: 签名值
//   - error: 错误信息（成功时为nil）
func WxPayTransactionsJsapi(ctx context.Context, wxPayConfig IWxPayConfig, openid string, outTradeNo string, totalFee int, description string) (*JsapiPayParams, error) {

	prepayID, err := createJsapiOrder(ctx, wxPayConfig, openid, outTradeNo, totalFee, description)
	if err != nil {
		//fmt.Printf("下单失败: %v\n", err)
		return nil, err
	}
	//fmt.Printf("下单成功! PrepayID: %s\n", prepayID)

	payParams, err := generateJsapiPayParams(ctx, wxPayConfig, prepayID)
	if err != nil {
		fmt.Printf("生成支付参数失败: %v\n", err)
		return nil, err
	}

	//fmt.Println("JSAPI支付参数:")
	//fmt.Printf("%+v\n", payParams)
	return payParams, nil
}
