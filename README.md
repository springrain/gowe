
## 介绍
golang 多微信号SDK,[readygo](https://gitee.com/chunanyong/readygo)子项目  [API文档](https://pkg.go.dev/gitee.com/chunanyong/gowe?tab=doc)  

感谢 [https://gitee.com/xiaochengtech/wechat](https://gitee.com/xiaochengtech/wechat) 提供的基础代码

``` 
go get gitee.com/chunanyong/gowe 
```
* 支持境内普通商户和境内服务商(境外和银行服务商没有条件测试)
* 全部参数和返回值均使用`struct`类型传递
* 缓存前置,使用项目现有的缓存体系  
* 原生支持多微信号  
* 支持跳板请求微信API服务.例如内网服务器没有网络出口权限,可以使用Nginx跳板请求微信API服务 
* 原生支持集群部署  
 
## 初始化
WxPayConfig初始化
```go
type WxPayConfig struct {
    Id                   string
    AppId                string
    Secret               string
    MchID                string // 商户号
    MchAPIv3Key          string // 商户API v3密钥 (32字节)
    WechatPayCertificate string // 微信支付平台证书（示例，实际应从微信平台下载并定期更新）
    CertSerialNo         string // 商户证书序列号
    PrivateKey           string // 商户API私钥
    NotifyURL            string // 支付结果通知地址
}

var Wx = &WxPayConfig{
    Id:                   "xxx",
    AppId:                "xxx",               // 小程序/公众号AppID
    Secret:               "xxx", // 小程序/公众号密钥
    MchID:                "xxx",                       // 商户号
    MchAPIv3Key:          "xxx",  // 商户API v3密钥 (32字节)
    WechatPayCertificate: "xxx",//商户证书wechatpay_开头的
    CertSerialNo:         "xxx",// 商户API证书的序列号
    PrivateKey:           "xxx",// 商户API私钥
    NotifyURL:            "https://xxx/xxx/xxx",//回调地址
}

func (w WxPayConfig) GetId(ctx context.Context) string {
    return w.Id
}

func (w WxPayConfig) GetAppId(ctx context.Context) string {
    return w.AppId
}

func (w WxPayConfig) GetAccessToken(ctx context.Context) string {
    panic("implement me")
}

func (w WxPayConfig) GetSecret(ctx context.Context) string {
    return w.Secret
}

func (w WxPayConfig) GetCertificateFile(ctx context.Context) string {
    return "../cert/apiclient_cert.pem"
}

func (w WxPayConfig) GetSubAppId(ctx context.Context) string {
    panic("implement me")
}
func (w WxPayConfig) GetAPIKey(ctx context.Context) string {
    return w.MchAPIv3Key
}

func (w WxPayConfig) GetSignType(ctx context.Context) string {
    return "MD5"
}

func (w WxPayConfig) GetServiceType(ctx context.Context) int {
    return 1
}

func (w WxPayConfig) IsProd(ctx context.Context) bool {
    return true
}

func (w WxPayConfig) GetMchID(ctx context.Context) string {
    return w.MchID
}
func (w WxPayConfig) GetMchAPIv3Key(ctx context.Context) string {
    return w.MchAPIv3Key
}
func (w WxPayConfig) GetWechatPayCertificate(ctx context.Context) string {
    return w.WechatPayCertificate
}
func (w WxPayConfig) GetCertSerialNo(ctx context.Context) string {
    return w.CertSerialNo
}
func (w WxPayConfig) GetPrivateKey(ctx context.Context) string {
    return w.PrivateKey
}
func (w WxPayConfig) GetNotifyURL(ctx context.Context) string {
    return w.NotifyURL
}

```


## 使用

以下是通用的接口，WxConfig 设置为全局变量，使用`gowe.XXX`调用

使用样例：

```go

func TestGetAccessToken(t *testing.T)  {
	token, err := gowe.GetAccessToken(wx)
	if err != nil {
		t.Log("error:" ,err)
	}
	t.Log("token:",token)
}

```
## 微信支付v3

* Native支付(扫码支付) `WxPayTransactionsNative`
* 统一下单：`WxPayTransactionsJsapi`
## 微信支付v3回调

* 验签：`VerifyWechatSignature`
* 回调：`WechatPayCallback`

##基于hertz框架(示例)

```go
//微信支付统一下单,返回一个prepay_id,参数解释请看方法注释
jsapi, err := gowe.WxPayTransactionsJsapi(ctx, Wx, OpenId, 1, "测试微信支付")
//扫码支付(Native方式),返回一个code_url,参数解释请看方法注释
native, err := gowe.WxPayTransactionsNative(ctx, Wx, ip, "0001", 1, "测试微信支付")
//回调
func HandleWechatPayCallback(ctx context.Context, c *app.RequestContext) {
	// 1. 验证签名
	fmt.Println("微信回调请求到达----------")

	// 从请求头获取签名相关字段
	timestamp := string(c.Request.Header.Peek("Wechatpay-Timestamp"))
	nonce := string(c.Request.Header.Peek("Wechatpay-Nonce"))
	signature := string(c.Request.Header.Peek("Wechatpay-Signature"))
	serial := string(c.Request.Header.Peek("Wechatpay-Serial"))

	if timestamp == "" || nonce == "" || signature == "" || serial == "" {
		fmt.Println("缺少签名头----------")
		c.JSON(400, map[string]string{"message": "缺少签名头"})
		return
	}

	body, err := c.Body()
	if err != nil {
		fmt.Println("读取请求体失败----------", err.Error())
		c.JSON(400, map[string]string{"message": "读取请求体失败"})
		return
	}
	// 通常不需要再手动调用 SetBody，但如果你后续处理需要，可以设置
	c.Request.SetBody(body)
	//验签
	err = gowe.VerifyWechatSignature(ctx, Wx, timestamp, nonce, signature, serial, body)
	if err != nil {
		fmt.Printf("签名验证失败: %v", err)
		c.JSON(400, map[string]string{"message": fmt.Sprintf("签名验证失败: %v", err)})
		return
	}

	//回调
	callback := gowe.WechatPayCallback(ctx, Wx, body)
	if callback.Code == 0 {
		fmt.Println("微信回调成功----------")
		// 7. 返回成功
		c.JSON(200, map[string]string{"message": "Success"})
	} else {
		fmt.Println("微信回调失败----------")
		// 7. 返回成功
		c.JSON(400, map[string]string{"message": "error"})
	}

}
```



## 微信支付v2

* 提交付款码支付 `WxPayMicropay`
* 查询订单：`WxPayQueryOrder`
* 关闭订单：`WxPayCloseOrder`
* 撤销订单：`WxPayReverse`
* 申请退款：`WxPayRefund`
* 查询退款：`WxPayQueryRefund`
* 下载对账单：`WxPayDownloadBill`
* 交易保障(JSAPI)：`WxPayReportJsApi`
* 交易保障(MICROPAY)：`WxPayReportMicropay` 

## 微信红包

* 发送现金红包 `WxPaySendRedPack`
* 发送裂变红包 `WxPaySendGroupRedPack`
* 发送小程序红包 `WxPaySendMiniProgramHB`
* 查询红包记录  `WxPayGetHBInfo`

## 企业付款

* 企业付款到零钱 `WxPayPromotionMktTransfers`
* 查询企业付款 `WxPayQueryMktTransfer`

## 微信支付回调
* 退款回调：`WxPayNotifyRefund`

## 微信公众号

* 获取基础支持的AccessToken：`WxMpWebAuthAccessToken`
* 获取用户基本信息(UnionId机制)：`WxMpGetUserInfo`
* 获取H5支付签名：`WxPayH5Sign`
* 临时二维码：`WxMpQrCreateTemporary` 
* 永久二维码：`WxMpQrCreatePermanent`
* 发送模板消息：`WxMpTemplateMsgSend` 
* 发送订阅消息: `WxMpSubscribeMsgSend`

## 微信小程序

* 获取小程序支付签名：`WxPayMaSign`
* 获取小程序码：`WxMaCodeGetUnlimited`
* 发送订阅消息：`WxMaSubscribeMessageSend`

## 文档

* 微信支付文档:[https://pay.weixin.qq.com/wiki/doc/api/index.html](https://pay.weixin.qq.com/wiki/doc/api/index.html)
* 随机数生成算法:[https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_3](https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_3)
* 签名生成算法:[https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_3](https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_3)
* 交易金额:[https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2](https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2)
* 交易类型:[https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2](https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2)
* 货币类型:[https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2](https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2)
* 时间规则:[https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2](https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2)
* 时间戳:[https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2](https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2)
* 商户订单号:[https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2](https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2)
* 银行类型:[https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2](https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=4_2)
* 单品优惠功能字段:[https://pay.weixin.qq.com/wiki/doc/api/danpin.php?chapter=9_101&index=1](https://pay.weixin.qq.com/wiki/doc/api/danpin.php?chapter=9_101&index=1)
* 代金券或立减优惠:[https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=12_1](https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=12_1)
* 最新县及县以上行政区划代码:[https://pay.weixin.qq.com/wiki/doc/api/download/store_adress.csv](https://pay.weixin.qq.com/wiki/doc/api/download/store_adress.csv)   
