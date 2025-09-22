package gowe

import (
	"context"
)

// 请求结构体
type NativeOrderRequest struct {
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
	SceneInfo struct {
		PayerClientIp string `json:"payer_client_ip"`
		StoreInfo     struct {
			ID string `json:"id"`
		} `json:"store_info"`
	} `json:"scene_info"`
}

// 响应结构体
type NativeOrderResponse struct {
	CodeUrl string `json:"code_url"`
}

// WxPayTransactionsNative 微信支付V3 Native支付（扫码支付）下单接口
// 参数说明:
//   - ctx: 上下文对象，用于传递截止时间、取消信号等跨API信息
//   - wxPayConfig: 微信支付配置接口，必须实现以下方法:
//   - ip: 用户终端IP（IPv4/IPv6格式），用于风控
//   - storeId: 商户门店编号（可选）:
//   - totalFee: 订单总金额（单位：分）:
//   - description: 商品描述（必填）:
//
// 返回值:
//   - NativeOrderResponse: Native支付下单响应结构，包含:
//     - CodeURL string `json:"code_url"`  // 二维码链接（有效期2小时）

func WxPayansactionsNative(ctx context.Context, wxPayConfig IWxPayConfig, ip string, outTradeNo string, storeId string, totalFee int, description string) NativeOrderResponse {
	resp := NativeOrderResponse{}
	codeUrl, err := createNativeOrder(ctx, wxPayConfig, ip, storeId, outTradeNo, totalFee, description)
	if err != nil {
		//fmt.Printf("下单失败: %v\n", err)
		return resp
	}
	//fmt.Printf("下单成功! PrepayID: %s\n", codeUrl)
	resp.CodeUrl = codeUrl
	return resp
}
