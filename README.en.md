# gowe

#### 介绍
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

### 初始化

```go
type WxConfig struct {
	Id     string
	AppId  string
	Secret string
}

var wxConfig = &WxConfig{
	Id:     "test",
	AppId:  "XXXXXXXXXXXXXXxxx",
	Secret: "XXXXXXXXXXXXXXX",
}

func (wxConfig *WxConfig) GetId() string {
	return wxConfig.Id
}

func (wxConfig *WxConfig) GetAppId() string {
	return wxConfig.AppId
}

func (wxConfig *WxConfig) GetAccessToken() string {
	//从缓存中获取wxAccessToken,这里只是演示
	wxAccessToken, err := GetAccessToken(wxConfig)
	if err == nil && wxAccessToken.ErrCode == 0 {
		return wxAccessToken.AccessToken
	}
	return ""
}

func (wxConfig *WxConfig) GetSecret() string {
	return wxConfig.Secret
}

```

### 使用

以下是通用的接口，WxConfig 设置为全局变量，可以直接使用`gowe.XXX`调用

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

#### 微信支付

* 提交付款码支付 `WxPayMicropay`
* 统一下单：`WxPayUnifiedOrder`
* 查询订单：`WxPayQueryOrder`
* 关闭订单：`WxPayCloseOrder`
* 撤销订单：`WxPayReverse`
* 申请退款：`WxPayRefund`
* 查询退款：`WxPayQueryRefund`
* 下载对账单：`WxPayDownloadBill`
* 交易保障(JSAPI)：`WxPayReportJsApi`
* 交易保障(MICROPAY)：`WxPayReportMicropay` 

#### 微信支付回调

* 支付回调：`WxPayNotifyPay`
* 退款回调：`WxPayNotifyRefund`

#### 微信公众号

* 获取基础支持的AccessToken：`WxMpWebAuthAccessToken`
* 获取用户基本信息(UnionId机制)：`WxMpGetUserInfo`
* 获取H5支付签名：`WxPayH5Sign`

#### 微信小程序

* 获取小程序支付签名：WxPayAppSign
* 获取小程序码：`WxMpQrCreateTemporary //生成带参数的临时二维码 WxMpQrCreatePermanent //创建永久的带参数二维码` 

### 文档

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
