package gowe

import "encoding/xml"

//微信红包APi

//WxPaySendRedPack 发送红包
//https://pay.weixin.qq.com/wiki/doc/api/tools/cash_coupon.php?chapter=13_4&index=3
func WxPaySendRedPack(wxPayConfig IWxPayConfig, body *WxPaySendRedPackBody) (*WxPaySendRedPackResponse, error) {
	apiurl := WxPayMchAPIURL + "/mmpaymkttransfers/sendredpack"
	// 业务逻辑
	bytes, err := wxPayDoWeChatWithCert(wxPayConfig, apiurl, body)
	if err != nil {
		return nil, err
	}
	// 结果校验
	if err = wxPayDoVerifySign(wxPayConfig, bytes, true); err != nil {
		return nil, err
	}
	// 解析返回值
	res := &WxPaySendRedPackResponse{}
	err = xml.Unmarshal(bytes, res)
	return res, err
}

//WxPaySendRedPackBody 微信发送红包参数
type WxPaySendRedPackBody struct {
	MchBillno   string `json:"mch_billno"`          // 商户订单号(每个订单号必须唯一.取值范围:0~9,a~z,A~Z)	接口根据商户订单号支持重入,如出现超时可再调用.
	SendName    string `json:"send_name"`           // 红包发送者名称 注意:敏感词会被转义成字符*
	ReOpenid    string `json:"re_openid"`           // 接受红包的用户openid openid为用户在wxappid下的唯一标识(获取openid参见微信公众平台开发者文档:网页授权获取用户基本信息)
	TotalAmount int    `json:"total_amount"`        //付款金额,单位分
	TotalNum    int    `json:"total_num"`           //红包发放总人数 total_num=1
	Wishing     string `json:"wishing"`             //红包祝福语 注意:敏感词会被转义成字符*
	ClientIp    string `json:"client_ip"`           //调用接口的机器Ip地址
	ActName     string `json:"act_name"`            //活动名称 注意:敏感词会被转义成字符*
	Remark      string `json:"remark"`              //备注信息
	SceneId     string `json:"scene_id,omitempty"`  //发放红包使用场景,红包金额大于200或者小于1元时必传. PRODUCT_1:商品促, PRODUCT_2:抽奖, PRODUCT_3:虚拟物品兑奖, PRODUCT_4:企业内部福利, PRODUCT_5:渠道分润, PRODUCT_6:保险回馈, PRODUCT_7:彩票派奖, PRODUCT_8:税务刮奖
	RiskInfo    string `json:"risk_info,omitempty"` //活动信息 posttime:用户操作的时间戳 mobile:业务系统账号的手机号,国家代码-手机号.不需要+号 deviceid :mac 地址或者设备唯一标识 clientversion :用户操作的客户端版本 把值为非空的信息用key=value进行拼接,再进行urlencode urlencode(posttime=xx& mobile =xx&deviceid=xx)
}

//WxPaySendRedPackResponse 微信发送红包返回值
type WxPaySendRedPackResponse struct {
	ReturnCode string `xml:"return_code"` // SUCCESS/FAIL 此字段是通信标识,非交易标识,交易是否成功需要查看result_code来判断
	ReturnMsg  string `xml:"return_msg"`  // 返回信息,如非空,为错误原因:签名失败/参数格式校验错误

	//以下字段在return_code为SUCCESS的时候有返回
	ResultCode  string `xml:"result_code"`  // SUCCESS/FAIL
	ErrCode     string `xml:"err_code"`     // 详细参见第6节错误列表
	ErrCodeDes  string `xml:"err_code_des"` // 错误返回的信息描述
	MchBillno   string `xml:"mch_billno"`   // 商户订单号(每个订单号必须唯一.取值范围:0~9,a~z,A~Z)	接口根据商户订单号支持重入,如出现超时可再调用.
	MchId       string `xml:"mch_id"`       // 微信支付分配的商户号
	WxAppId     string `xml:"wxappid"`      // 微信分配的公众账号ID(企业号corpid即为此appId).在微信开放平台(open.weixin.qq.com)申请的移动应用appid无法使用该接口
	ReOpenid    string `xml:"re_openid"`    // 接受红包的用户openid openid为用户在wxappid下的唯一标识(获取openid参见微信公众平台开发者文档:网页授权获取用户基本信息)
	TotalAmount int    `xml:"total_amount"` // 付款金额,单位分
	SendListid  string `xml:"send_listid"`  // 红包订单的微信单号

}
