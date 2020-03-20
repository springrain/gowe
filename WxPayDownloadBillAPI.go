package gowe

import (
	"encoding/xml"
	"errors"
)

//WxPayDownloadBill 下载对账单 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_6
func WxPayDownloadBill(wxPayConfig IWxPayConfig, body *WxPayDownloadBillBody) (*WxPayDownloadBillResponse, error) {
	// 业务逻辑
	bytes, err := wxPayDoWeChat(wxPayConfig, WxPayMchAPIURL+"/pay/downloadbill", body)
	if err != nil {
		return nil, err
	}
	// 解析返回值
	failRsp := new(WxPayDownloadBillResponse)
	err = xml.Unmarshal(bytes, failRsp)
	if err != nil {
		return nil, errors.New(string(bytes))
	} else {
		err = errors.New(failRsp.ReturnMsg)
		return failRsp, err
	}
}

//WxPayDownloadBillBody 下载对账单的参数
type WxPayDownloadBillBody struct {
	SignType string `json:"sign_type,omitempty"` // 签名类型,目前支持HMAC-SHA256和MD5,默认为MD5
	BillDate string `json:"bill_date"`           // 下载对账单的日期,格式:20140603
	BillType string `json:"bill_type,omitempty"` // ALL,返回当日所有订单信息,默认值 SUCCESS,返回当日成功支付的订单 REFUND,返回当日退款订单 RECHARGE_REFUND,返回当日充值退款订单
	TarType  string `json:"tar_type,omitempty"`  // 非必传参数,固定值:GZIP,返回格式为.gzip的压缩包账单.不传则默认为数据流形式.
}

//WxPayDownloadBillResponse 下载对账单的返回值
type WxPayDownloadBillResponse struct {
	WxResponseModel
	ErrCode string `xml:"err_code"` // 失败错误码,详见错误码列表 TODO
}
