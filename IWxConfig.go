package gowe

//IWxConfig 微信的基础配置
type IWxConfig interface {
	//GetId 数据库记录的Id
	GetId() (string, error)
	//GetAppId 微信号的appId
	GetAppId() (string, error)
	//GetAccessToken 获取到AccessToken
	GetAccessToken() (string, error)
	//GetAccessToken 设置AccessToken
	SetAccessToken(accessToken string) (string, error)
	//GetAppId 微信号的secret
	GetSecret() (string, error)
}

//IWxMpCongfig 公众号的配置
type IWxMpCongfig interface {
	IWxConfig

	//GetToken 获取token
	GetToken() (string, error)

	//GetAesKey 获取aesKey
	GetAesKey() (string, error)

	//开启oauth2.0认证,是否能够获取openId,0是关闭,1是开启
	GetOauth2() (int, error)
}

//IWxPayConfig 公众号的配置
type IWxPayConfig interface {
	IWxConfig

	//GetCertificateFile 获取商户证路径
	GetCertificateFile() (string, error)

	//GetMchId 获取 Mch ID
	GetMchId() (string, error)

	//GetKey 获取 API 密钥
	GetKey() (string, error)

	//GetNotifyUrl 获取回调地址
	GetNotifyUrl() (string, error)

	//GetSignType 获取加密类型
	GetSignType() (string, error)
}
