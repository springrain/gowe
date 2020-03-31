package gowe

import "encoding/xml"

//WxPayReportJsApi  交易保障(JSAPI) https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_8&index=9
func WxPayReportJsApi(wxPayConfig IWxPayConfig, body *WxPayReportJsAPIBody) (*WxResponseModel, error) {
	var err error
	// 处理参数
	if body.InterfaceUrl, err = encodePath(body.InterfaceUrl); err != nil {
		return nil, err
	}
	// 业务逻辑
	bytes, err := wxPayDoWeChat(wxPayConfig, WxPayMchAPIURL+"/payitil/report", body, 0)
	if err != nil {
		return nil, err
	}
	// 解析返回值
	res := &WxResponseModel{}
	err = xml.Unmarshal(bytes, res)
	return res, err
}

//WxPayReportJsAPIBody 交易保障(JSAPI)的参数
type WxPayReportJsAPIBody struct {
	SignType     string `json:"sign_type,omitempty"`    // 签名类型,目前支持HMAC-SHA256和MD5,默认为MD5
	DeviceInfo   string `json:"device_info,omitempty"`  // (非必填) 微信支付分配的终端设备号,商户自定义
	InterfaceUrl string `json:"interface_url"`          // 上报对应的接口的完整URL,类似:https://api.mch.weixin.qq.com/pay/unifiedorder 对于刷卡支付,为更好的和商户共同分析一次业务行为的整体耗时情况,对于两种接入模式,请都在门店侧对一次刷卡行为进行一次单独的整体上报,上报URL指定为:https://api.mch.weixin.qq.com/pay/micropay/total 关于两种接入模式具体可参考本文档章节:刷卡支付商户接入模式 其它接口调用仍然按照调用一次,上报一次来进行.
	ExecuteTime  int64  `json:"execute_time"`           // 接口耗时情况,单位为毫秒
	ReturnCode   string `json:"return_code"`            // SUCCESS/FAIL 此字段是通信标识,非交易标识,交易是否成功需要查看trade_state来判断
	ReturnMsg    string `json:"return_msg,omitempty"`   // (非必填) 返回信息,如非空,为错误原因 签名失败 参数格式校验错误
	ResultCode   string `json:"result_code"`            // SUCCESS/FAIL
	ErrCode      string `json:"err_code,omitempty"`     // (非必填) ORDERNOTEXIST—订单不存在 SYSTEMERROR—系统错误
	ErrCodeDes   string `json:"err_code_des,omitempty"` // (非必填) 结果信息描述
	OutTradeNo   string `json:"out_trade_no,omitempty"` // (非必填) 商户系统内部的订单号,商户可以在上报时提供相关商户订单号方便微信支付更好的提高服务质量.
	UserIp       string `json:"user_ip"`                // 发起接口调用时的机器IP
	Time         string `json:"time,omitempty"`         // (非必填) 系统时间,格式为yyyyMMddHHmmss,如2009年12月27日9点10分10秒表示为20091227091010.其他详见时间规则
}
