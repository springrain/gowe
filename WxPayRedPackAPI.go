package gowe

import (
	"encoding/xml"
	"strconv"

	"github.com/beevik/etree"
)

//微信红包APi

//WxPaySendRedPack 发送红包
//https://pay.weixin.qq.com/wiki/doc/api/tools/cash_coupon.php?chapter=13_4&index=3
func WxPaySendRedPack(wxPayConfig IWxPayConfig, body *WxPaySendRedPackBody) (*WxPaySendRedPackResponse, error) {
	apiurl := WxPayMchAPIURL + "/mmpaymkttransfers/sendredpack"
	// 业务逻辑
	bytes, err := wxPayDoWeChatWithCert(wxPayConfig, apiurl, body, 2)
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

//WxPaySendGroupRedPack 发送裂变红包
//https://pay.weixin.qq.com/wiki/doc/api/tools/cash_coupon.php?chapter=13_5&index=4
func WxPaySendGroupRedPack(wxPayConfig IWxPayConfig, body *WxPaySendGroupRedPackBody) (*WxPaySendGroupRedPackResponse, error) {
	apiurl := WxPayMchAPIURL + "/mmpaymkttransfers/sendgroupredpack"
	if len(body.AmtType) < 1 {
		body.AmtType = "ALL_RAND"
	}
	// 业务逻辑
	bytes, err := wxPayDoWeChatWithCert(wxPayConfig, apiurl, body, 2)
	if err != nil {
		return nil, err
	}
	// 结果校验
	if err = wxPayDoVerifySign(wxPayConfig, bytes, true); err != nil {
		return nil, err
	}
	// 解析返回值
	res := &WxPaySendGroupRedPackResponse{}
	err = xml.Unmarshal(bytes, res)
	return res, err
}

//WxPayGetHBInfo 查看红包记录,用于商户对已发放的红包进行查询红包的具体信息,可支持普通红包和裂变包.
//https://pay.weixin.qq.com/wiki/doc/api/tools/cash_coupon.php?chapter=13_6&index=5
func WxPayGetHBInfo(wxPayConfig IWxPayConfig, body *WxPayGetHBInfoBody) (*WxPayGetHBInfoResponse, error) {
	apiurl := WxPayMchAPIURL + "/mmpaymkttransfers/gethbinfo"
	// 业务逻辑
	bytes, err := wxPayDoWeChatWithCert(wxPayConfig, apiurl, body, 0)
	if err != nil {
		return nil, err
	}
	// 结果校验
	if err = wxPayDoVerifySign(wxPayConfig, bytes, true); err != nil {
		return nil, err
	}
	// 解析返回值
	res := &WxPayGetHBInfoResponse{}
	err = wxPayNewHBInfoResponse(bytes, res)
	return res, err
}

//WxPaySendMiniProgramHB 发送小程序红包
//https://pay.weixin.qq.com/wiki/doc/api/tools/cash_coupon.php?chapter=18_2&index=3
func WxPaySendMiniProgramHB(wxPayConfig IWxPayConfig, body *WxPaySendMiniProgramHBBody) (*WxPaySendMiniProgramHBResponse, error) {
	apiurl := WxPayMchAPIURL + "/mmpaymkttransfers/sendminiprogramhb"

	//通过JSAPI方式领取红包,小程序红包固定传MINI_PROGRAM_JSAPI
	if len(body.NotifyWay) < 1 {
		body.NotifyWay = "MINI_PROGRAM_JSAPI"
	}

	// 业务逻辑
	bytes, err := wxPayDoWeChatWithCert(wxPayConfig, apiurl, body, 2)
	if err != nil {
		return nil, err
	}
	// 结果校验
	if err = wxPayDoVerifySign(wxPayConfig, bytes, true); err != nil {
		return nil, err
	}
	// 解析返回值
	res := &WxPaySendMiniProgramHBResponse{}
	err = xml.Unmarshal(bytes, res)
	return res, err
}

//封装返回的裂变红包列表
func wxPayNewHBInfoResponse(xmlStr []byte, rsp *WxPayGetHBInfoResponse) error {
	// 常规解析
	if err := xml.Unmarshal(xmlStr, rsp); err != nil {
		return err
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlStr); err != nil {
		return err
	}
	//rootXML := doc.SelectElement("xml")
	//hblistXML := rootXML.SelectElement("hblist")
	//if hblistXML == nil {
	//	return nil
	//}

	hbList := make([]WxPayHBInfoModel, 0)
	for _, hbinfoXML := range doc.FindElements("./xml/hblist/*") {
		hbInfo := WxPayHBInfoModel{}
		if openid := hbinfoXML.SelectElement("openid"); openid != nil {
			hbInfo.Openid = openid.Text()
		} else if amount := hbinfoXML.SelectElement("amount"); amount != nil {
			am, err := strconv.Atoi(amount.Text())
			if err == nil {
				hbInfo.Amount = am
			}
		} else if rcv := hbinfoXML.SelectElement("rcv_time"); rcv != nil {
			hbInfo.RcvTime = rcv.Text()
		}
		hbList = append(hbList, hbInfo)

	}
	rsp.HBList = hbList
	return nil

}

//WxPayGetHBInfoBody 查看红包记录的请求参数
type WxPayGetHBInfoBody struct {
	MchBillno string `json:"mch_billno"` //商户发放红包的商户订单号
	BillType  string `json:"bill_type"`  //MCHT:通过商户订单号获取红包信息
}

//WxPayGetHBInfoResponse 查看红包记录的返回值
type WxPayGetHBInfoResponse struct {
	ReturnCode string `xml:"return_code"` // SUCCESS/FAIL 此字段是通信标识,非交易标识,交易是否成功需要查看result_code来判断
	ReturnMsg  string `xml:"return_msg"`  // 返回信息,如非空,为错误原因:签名失败/参数格式校验错误

	//以下字段在return_code为SUCCESS的时候有返回
	ResultCode string `xml:"result_code"`  // SUCCESS/FAIL
	ErrCode    string `xml:"err_code"`     // 详细参见第6节错误列表
	ErrCodeDes string `xml:"err_code_des"` // 错误返回的信息描述

	//以下字段在return_code 和result_code都为SUCCESS的时候有返回
	MchBillno    string `xml:"mch_billno"`    // 商户使用查询API填写的商户单号的原路返回
	MchId        string `xml:"mch_id"`        // 微信支付分配的商户号
	DetailId     string `xml:"detail_id"`     // 使用API发放现金红包时返回的红包单号
	Status       string `xml:"status"`        // SENDING:发放中,SENT:已发放待领取,FAILED：发放失败,RECEIVED:已领取,RFUND_ING:退款中,REFUND:已退款
	SendType     string `xml:"send_type"`     // 发送类型 API:通过API接口发放,UPLOAD:通过上传文件方式发放,ACTIVITY:通过活动方式发放
	HBType       string `xml:"hb_type"`       // 红包类型 GROUP:裂变红包,NORMAL:普通红包
	TotalNum     int    `json:"total_num"`    // 红包个数
	TotalAmount  int    `xml:"total_amount"`  // 红包总金额(单位分)
	Reason       string `xml:"reason"`        // 发送失败原因
	SendTime     string `xml:"send_time"`     // 红包发送时间 2015-04-21 20:00:00
	RefundTime   string `xml:"refund_time"`   // 红包的退款时间(如果其未领取的退款) 2015-04-21 23:03:00
	RefundAmount int    `xml:"refund_amount"` // 红包退款金额(单位分)
	Wishing      string `xml:"wishing"`       // 红包祝福语
	ActName      string `xml:"act_name"`      // 活动名称
	Remark       string `xml:"remark"`        // 活动描述，低版本微信可见

	Openid  string             `xml:"openid"`   // 领取红包的openid
	Amount  int                `xml:"amount"`   // 领取金额(单位分)
	RcvTime string             `xml:"rcv_time"` // 领取红包的时间 2015-04-21 20:00:00
	HBList  []WxPayHBInfoModel `xml:"-"`        // 裂变红包的列表
}

//WxPayHBInfoModel 返回的微信裂变红包信息
type WxPayHBInfoModel struct {
	Openid  string // 领取红包的openid
	Amount  int    // 领取金额(单位分)
	RcvTime string // 领取红包的时间 2015-04-21 20:00:00

}

//WxPaySendGroupRedPackBody 微信裂变红包参数
type WxPaySendGroupRedPackBody struct {
	WxPaySendRedPackBody
	AmtType string `json:"amt_type"` //红包金额设置方式 ALL_RAND—全部随机,商户指定总金额和红包发放总人数，由微信支付随机计算出各红包金额
}

//WxPaySendGroupRedPackResponse 微信裂变红包返回值
type WxPaySendGroupRedPackResponse struct {
	WxPaySendRedPackResponse
}

//WxPaySendRedPackBody 微信发送红包参数
type WxPaySendRedPackBody struct {
	MchBillno   string `json:"mch_billno"`          // 商户订单号(每个订单号必须唯一.取值范围:0~9,a~z,A~Z)	接口根据商户订单号支持重入,如出现超时可再调用.
	SendName    string `json:"send_name"`           // 红包发送者名称 注意:敏感词会被转义成字符*
	ReOpenid    string `json:"re_openid"`           // 接受红包的用户openid openid为用户在wxappid下的唯一标识(获取openid参见微信公众平台开发者文档:网页授权获取用户基本信息)
	TotalAmount int    `json:"total_amount"`        //付款金额,单位分
	TotalNum    int    `json:"total_num"`           //红包发放总人数 total_num=1
	Wishing     string `json:"wishing"`             //红包祝福语 注意:敏感词会被转义成字符*
	ClientIp    string `json:"client_ip,omitempty"` //调用接口的机器Ip地址
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

//WxPaySendMiniProgramHBBody 微信发送小程序红包参数
type WxPaySendMiniProgramHBBody struct {
	WxPaySendRedPackBody
	NotifyWay string `json:"notify_way"` //通知用户形式	 通过JSAPI方式领取红包,小程序红包固定传MINI_PROGRAM_JSAPI
}

//WxPaySendMiniProgramHBResponse 微信发送小程序红包返回值
type WxPaySendMiniProgramHBResponse struct {
	WxPaySendRedPackResponse
	JSAPIPackage string `xml:"package"` // 返回jaspi的入参package的值
}
