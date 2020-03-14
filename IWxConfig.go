package gowe

// 微信API的服务器域名,方便处理请求代理跳转的情况
var WxMpAPIURL = "https://api.weixin.qq.com"
var WxMpWeiXinURL = "https://mp.weixin.qq.com"
var WxMpOpenURL = "https://open.weixin.qq.com"
var WxMpPayMchAPIURL = "https://api.mch.weixin.qq.com"
var WxqyAPIURL = "https://qyapi.weixin.qq.com"
var WxPayReporMchtURL = "http://report.mch.weixin.qq.com"
var WxPayAppURL = "https://payapp.weixin.qq.com"

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
