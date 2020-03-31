package gowe

//IWxConfig 微信的基础配置
type IWxConfig interface {
	//Id 数据库记录的Id
	GetId() string
	//AppId 微信号的appId
	GetAppId() string
	//AccessToken 获取到AccessToken
	GetAccessToken() string
	//Secret 微信号的secret
	GetSecret() string
}

//IWxMpConfig 公众号的配置
type IWxMpConfig interface {
	IWxConfig
	//Token 获取token
	GetToken() string
	//AesKey 获取aesKey
	GetAesKey() string
	//开启oauth2.0认证,是否能够获取openId,0是关闭,1是开启
	GetOauth2() bool
}

//IWxMaConfig 微信小程序配置
type IWxMaConfig interface {
	IWxConfig
}

//IWxPayConfig 公众号的配置
type IWxPayConfig interface {
	IWxConfig
	//证数文件路径
	GetCertificateFile() string
	//支付的mchId
	GetMchId() string
	GetSubAppId() string // 微信分配的子商户公众账号ID
	GetSubMchId() string // 微信支付分配的子商户号,开发者模式下必填
	//获取 API 密钥
	GetAPIKey() string
	//支付通知回调的地址
	GetNotifyUrl() string
	//摘要加密类型
	GetSignType() string
	GetServiceType() int // 服务模式
	IsProd() bool        // 是否是生产环境
	MchType() int        //0:普通商户接口, 1:特殊的商户接口(微信找零),2:红包

}
