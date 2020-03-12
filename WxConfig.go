package gowe

import (
	"context"
	"errors"
)

// 微信API的服务器域名,方便处理请求代理跳转的情况
var WxmpApiUrl = "https://api.weixin.qq.com"
var WxmpWeixinUrl = "https://mp.weixin.qq.com"
var WxmpOpenUrl = "https://open.weixin.qq.com"
var Wxmppaybaseurl = "https://api.mch.weixin.qq.com"
var Wxqyapiurl = "https://qyapi.weixin.qq.com"
var Wxreporturl = "http://report.mch.weixin.qq.com"
var Wxpayappbaseurl = "https://payapp.weixin.qq.com"

//全局变量,用于初始化,默认第一次初始化赋值,如果是多个微信号配置,就需要放到context
var wxConfig *WxConfig = nil
var wxMpConfig *WxMpConfig = nil
var wxPayConfig *WxPayConfig = nil

type wrapContextStringKey string

//context WithValue的key,不能是基础类型,例如字符串,包装一下
const contextWxConfigValueKey = wrapContextStringKey("contextWxConfigValueKey")

//WxConfig 微信的基础配置
type WxConfig struct {
	//Id 数据库记录的Id
	Id string
	//AppId 微信号的appId
	AppId string
	//AccessToken 获取到AccessToken
	AccessToken string
	//Secret 微信号的secret
	Secret string
}

//WxMpConfig 公众号的配置
type WxMpConfig struct {
	WxConfig
	//Token 获取token
	Token string
	//AesKey 获取aesKey
	AesKey string
	//开启oauth2.0认证,是否能够获取openId,0是关闭,1是开启
	Oauth2 int
}

//WxPayConfig 公众号的配置
type WxPayConfig struct {
	WxConfig
	//CertificateFile 获取商户证路径
	CertificateFile string
	//MchId 获取 Mch ID
	MchId string
	//Key 获取 API 密钥
	Key string
	//NotifyUrl 获取回调地址
	NotifyUrl string
	//SignType 获取加密类型
	SignType string
}

//BindContextWxConfig 绑定wxConfig到Context,用于多个公众号进行绑定
func BindContextWxConfig(parent context.Context, wxConfig *WxConfig) (context.Context, error) {
	if parent == nil {
		return nil, errors.New("context的parent不能为nil")
	}
	if wxConfig == nil {
		return parent, errors.New("wxConfig不能为空")
	}
	ctx := context.WithValue(parent, contextWxConfigValueKey, wxConfig)
	return ctx, nil
}

//从ctx获取*WxConfig,如果ctx中没有,返回默认值
func getWxConfig(ctx context.Context) (*WxConfig, error) {
	if ctx == nil {
		return nil, errors.New("ctx不能为空")
	}
	value := ctx.Value(contextWxConfigValueKey)
	if value == nil {
		return wxConfig, nil
	}

	config := value.(*WxConfig)
	return config, nil
}

//从ctx获取*WxConfig,如果ctx中没有,返回默认值
func getWxMpConfig(ctx context.Context) (*WxMpConfig, error) {
	if ctx == nil {
		return nil, errors.New("ctx不能为空")
	}
	value := ctx.Value(contextWxConfigValueKey)
	if value == nil {
		return wxMpConfig, nil
	}

	config := value.(*WxMpConfig)
	return config, nil
}

//错误代码
var errCodeToErrMsgMap = map[int]string{
	-1:      "系统繁忙",
	0:       "请求成功",
	40001:   "获取access_token时AppSecret错误，或者access_token无效",
	40002:   "不合法的凭证类型",
	40003:   "不合法的OpenID",
	40004:   "不合法的媒体文件类型",
	40005:   "不合法的文件类型",
	40006:   "不合法的文件大小",
	40007:   "不合法的媒体文件id",
	40008:   "不合法的消息类型",
	40009:   "不合法的图片文件大小",
	40010:   "不合法的语音文件大小",
	40011:   "不合法的视频文件大小",
	40012:   "不合法的缩略图文件大小",
	40013:   "不合法的APPID",
	40014:   "不合法的access_token",
	40015:   "不合法的菜单类型",
	40016:   "不合法的按钮个数",
	40017:   "不合法的按钮类型",
	40018:   "不合法的按钮名字长度",
	40019:   "不合法的按钮KEY长度",
	40020:   "不合法的按钮URL长度",
	40021:   "不合法的菜单版本号",
	40022:   "不合法的子菜单级数",
	40023:   "不合法的子菜单按钮个数",
	40024:   "不合法的子菜单按钮类型",
	40025:   "不合法的子菜单按钮名字长度",
	40026:   "不合法的子菜单按钮KEY长度",
	40027:   "不合法的子菜单按钮URL长度",
	40028:   "不合法的自定义菜单使用用户",
	40029:   "不合法的oauth_code",
	40030:   "不合法的refresh_token",
	40031:   "不合法的openid列表",
	40032:   "不合法的openid列表长度,一次只能拉黑20个用户",
	40033:   "不合法的请求字符，不能包含\\uxxxx格式的字符",
	40035:   "不合法的参数",
	40037:   "不合法的模板id",
	40038:   "不合法的请求格式",
	40039:   "不合法的URL长度",
	40050:   "不合法的分组id",
	40051:   "分组名字不合法",
	40053:   "不合法的actioninfo，请开发者确认参数正确",
	40056:   "不合法的Code码",
	40059:   "不合法的消息id",
	40071:   "不合法的卡券类型",
	40072:   "不合法的编码方式",
	40078:   "card_id未授权",
	40079:   "不合法的时间",
	40080:   "不合法的CardExt",
	40097:   "参数不正确，请参考字段要求检查json字段",
	40099:   "卡券已被核销",
	40100:   "不合法的时间区间",
	40116:   "不合法的Code码",
	40122:   "不合法的库存数量",
	40124:   "会员卡设置查过限制的 custom_field字段",
	40127:   "卡券被用户删除或转赠中",
	40130:   "不合法的openid列表长度, 长度至少大于2个", //invalid openid list size, at least two openid
	41001:   "缺少access_token参数",
	41002:   "缺少appid参数",
	41003:   "缺少refresh_token参数",
	41004:   "缺少secret参数",
	41005:   "缺少多媒体文件数据",
	41006:   "缺少media_id参数",
	41007:   "缺少子菜单数据",
	41008:   "缺少oauth code",
	41009:   "缺少openid",
	41011:   "缺少必填字段",
	41012:   "缺少cardid参数",
	42001:   "access_token超时",
	42002:   "refresh_token超时",
	42003:   "oauth_code超时",
	43001:   "需要GET请求",
	43002:   "需要POST请求",
	43003:   "需要HTTPS请求",
	43004:   "需要接收者关注",
	43005:   "需要好友关系",
	43009:   "自定义SN权限，请前往公众平台申请",
	43010:   "无储值权限，请前往公众平台申请",
	43100:   "修改模板所属行业太频繁",
	44001:   "多媒体文件为空",
	44002:   "POST的数据包为空",
	44003:   "图文消息内容为空",
	44004:   "文本消息内容为空",
	45001:   "多媒体文件大小超过限制",
	45002:   "消息内容超过限制",
	45003:   "标题字段超过限制",
	45004:   "描述字段超过限制",
	45005:   "链接字段超过限制",
	45006:   "图片链接字段超过限制",
	45007:   "语音播放时间超过限制",
	45008:   "图文消息超过限制",
	45009:   "接口调用超过限制",
	45010:   "创建菜单个数超过限制",
	45015:   "回复时间超过限制",
	45016:   "系统分组，不允许修改",
	45017:   "分组名字过长",
	45018:   "分组数量超过上限",
	45027:   "模板与所选行业不符", //template conflict with industry
	45028:   "没有群发配额",    //has no masssend quota
	45030:   "该cardid无接口权限",
	45031:   "库存为0",
	45033:   "用户领取次数超过限制get_limit",
	45056:   "创建的标签数过多，请注意不能超过100个",
	45057:   "该标签下粉丝数超过10w，不允许直接删除",
	45058:   "不能修改0/1/2这三个系统默认保留的标签",
	45059:   "有粉丝身上的标签数已经超过限制",
	45157:   "标签名非法，请注意不能和其他标签重名",
	45158:   "标签名长度超过30个字节",
	45159:   "非法的tag_id",
	46001:   "不存在媒体数据",
	46002:   "不存在的菜单版本",
	46003:   "不存在的菜单数据",
	46004:   "不存在的用户",
	46005:   "不存在的门店",
	47001:   "解析JSON/XML内容错误",
	48001:   "api功能未授权",
	48004:   "api接口被封禁，请登录mp.weixin.qq.com查看详情",
	49003:   "传入的openid不属于此AppID",
	50001:   "用户未授权该api",
	50002:   "用户受限，可能是违规后接口被封禁",
	61451:   "参数错误(invalid parameter)",
	61452:   "无效客服账号(invalid kf_account)",
	61453:   "客服帐号已存在(kf_account exsited)",
	61454:   "客服帐号名长度超过限制(仅允许10个英文字符，不包括@及@后的公众号的微信号)(invalid kf_acount length)",
	61455:   "客服帐号名包含非法字符(仅允许英文+数字)(illegal character in kf_account)",
	61456:   "客服帐号个数超过限制(10个客服账号)(kf_account count exceeded)",
	61457:   "无效头像文件类型(invalid file type)",
	61450:   "系统错误(system error)",
	61500:   "日期格式错误",
	65104:   "门店的类型不合法，必须严格按照附表的分类填写",
	65105:   "图片url 不合法，必须使用接口1 的图片上传接口所获取的url",
	65106:   "门店状态必须未审核通过",
	65107:   "扩展字段为不允许修改的状态",
	65109:   "门店名为空",
	65110:   "门店所在详细街道地址为空",
	65111:   "门店的电话为空",
	65112:   "门店所在的城市为空",
	65113:   "门店所在的省份为空",
	65114:   "图片列表为空",
	65115:   "poi_id 不正确",
	65301:   "不存在此menuid对应的个性化菜单",
	65302:   "没有相应的用户",
	65303:   "没有默认菜单，不能创建个性化菜单",
	65304:   "MatchRule信息为空",
	65305:   "个性化菜单数量受限",
	65306:   "不支持个性化菜单的帐号",
	65307:   "个性化菜单信息为空",
	65308:   "包含没有响应类型的button",
	65309:   "个性化菜单开关处于关闭状态",
	65310:   "填写了省份或城市信息，国家信息不能为空",
	65311:   "填写了城市信息，省份信息不能为空",
	65312:   "不合法的国家信息",
	65313:   "不合法的省份信息",
	65314:   "不合法的城市信息",
	65316:   "该公众号的菜单设置了过多的域名外跳（最多跳转到3个域名的链接）",
	65317:   "不合法的URL",
	65400:   "API不可用，即没有开通/升级到新客服功能",
	65401:   "无效客服帐号",
	65402:   "帐号尚未绑定微信号，不能投入使用",
	65403:   "客服昵称不合法",
	65404:   "客服帐号不合法",
	65405:   "帐号数目已达到上限，不能继续添加",
	65406:   "已经存在的客服帐号",
	65407:   "邀请对象已经是本公众号客服",
	65408:   "本公众号已发送邀请给该微信号",
	65409:   "无效的微信号",
	65410:   "邀请对象绑定公众号客服数量达到上限（目前每个微信号最多可以绑定5个公众号客服帐号）",
	65411:   "该帐号已经有一个等待确认的邀请，不能重复邀请",
	65412:   "该帐号已经绑定微信号，不能进行邀请",
	65413:   "不存在对应用户的会话信息",
	65414:   "客户正在被其他客服接待",
	65415:   "客服不在线",
	65416:   "查询参数不合法",
	65417:   "查询时间段超出限制",
	72015:   "没有操作发票的权限，请检查是否已开通相应权限。",
	72017:   "发票抬头不一致",
	72023:   "发票已被其他公众号锁定",
	72024:   "发票状态错误",
	72025:   "wx_invoice_token无效",
	72028:   "未设置微信支付商户信息",
	72029:   "未设置授权字段",
	72030:   "mchid无效",
	72031:   "参数错误。可能为请求中包括无效的参数名称或包含不通过后台校验的参数值",
	72035:   "发票已经被拒绝开票",
	72036:   "发票正在被修改状态，请稍后再试",
	72038:   "订单没有授权，可能是开票平台appid、商户appid、订单order_id不匹配",
	72039:   "订单未被锁定",
	72040:   "Pdf无效，请提供真实有效的pdf",
	72042:   "发票号码和发票代码重复",
	72043:   "发票号码和发票代码错误",
	72044:   "发票抬头二维码超时",
	88000:   "没有留言权限",
	9001001: "POST数据参数不合法",
	9001002: "远端服务不可用",
	9001003: "Ticket不合法",
	9001004: "获取摇周边用户信息失败",
	9001005: "获取商户信息失败",
	9001006: "获取OpenID失败",
	9001007: "上传文件缺失",
	9001008: "上传素材的文件类型不合法",
	9001009: "上传素材的文件尺寸不合法",
	9001010: "上传失败",
	9001020: "帐号不合法",
	9001021: "已有设备激活率低于50%，不能新增设备",
	9001022: "设备申请数不合法，必须为大于0的数字",
	9001023: "已存在审核中的设备ID申请",
	9001024: "一次查询设备ID数量不能超过50",
	9001025: "设备ID不合法",
	9001026: "页面ID不合法",
	9001027: "页面参数不合法",
	9001028: "一次删除页面ID数量不能超过10",
	9001029: "页面已应用在设备中，请先解除应用关系再删除",
	9001030: "一次查询页面ID数量不能超过50",
	9001031: "时间区间不合法",
	9001032: "保存设备与页面的绑定关系参数错误",
	9001033: "门店ID不合法",
	9001034: "设备备注信息过长",
	9001035: "设备申请参数不合法",
	9001036: "查询起始值begin不合法",
}
