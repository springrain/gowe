package gowe

import "context"

//IWxConfig 微信的基础配置
type IWxConfig interface {
	//Id 数据库记录的Id
	GetId(ctx context.Context) string
	//AppId 微信号的appId
	GetAppId(ctx context.Context) string
	//AccessToken 获取到AccessToken
	GetAccessToken(ctx context.Context) string
	//Secret 微信号的secret
	GetSecret(ctx context.Context) string
}

//IWxMpConfig 公众号的配置
type IWxMpConfig interface {
	IWxConfig
	//Token 获取token
	GetToken(ctx context.Context) string
	//AesKey 获取aesKey
	GetAesKey(ctx context.Context) string
	//开启oauth2.0认证,是否能够获取openId,0是关闭,1是开启
	GetOauth2(ctx context.Context) bool
}

//IWxMaConfig 微信小程序配置
type IWxMaConfig interface {
	IWxConfig
}

//IWxPayConfig 公众号的配置
type IWxPayConfig interface {
	IWxConfig
	//证数文件路径
	GetCertificateFile(ctx context.Context) string
	//支付的mchId
	GetMchId(ctx context.Context) string
	GetSubAppId(ctx context.Context) string // 微信分配的子商户公众账号ID
	GetSubMchId(ctx context.Context) string // 微信支付分配的子商户号,开发者模式下必填
	//获取 API 密钥
	GetAPIKey(ctx context.Context) string
	//支付通知回调的地址
	GetNotifyUrl(ctx context.Context) string
	//摘要加密类型
	GetSignType(ctx context.Context) string
	GetServiceType(ctx context.Context) int // 服务模式
	IsProd(ctx context.Context) bool        // 是否是生产环境

}
