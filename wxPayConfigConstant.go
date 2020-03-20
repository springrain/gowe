package gowe

// 微信API的服务器域名,方便处理请求代理跳转的情况
var WxMpAPIURL = "https://api.weixin.qq.com"
var WxMpWeiXinURL = "https://mp.weixin.qq.com"
var WxMpOpenURL = "https://open.weixin.qq.com"
var WxMpPayMchAPIURL = "https://api.mch.weixin.qq.com"
var WxMpPaySanBoxAPIURL = "https://api.mch.weixin.qq.com/sandboxnew"
var WxqyAPIURL = "https://qyapi.weixin.qq.com"
var WxPayReporMchtURL = "http://report.mch.weixin.qq.com"
var WxPayAppURL = "https://payapp.weixin.qq.com"

const (

	//baseUrl        = "https://api.mch.weixin.qq.com/"            // (生产环境) 微信支付的基地址
	//baseUrlSandbox = "https://api.mch.weixin.qq.com/sandboxnew/" // (沙盒环境) 微信支付的基地址

	// wxURL_DownloadFundFlow  = wxBaseUrl + "pay/downloadfundflow"            // 下载资金账单
	// wxURL_BatchQueryComment = wxBaseUrl + "billcommentsp/batchquerycomment" //
	// wxURL_SanBox_DownloadFundFlow  = wxBaseUrlSandbox + "pay/downloadfundflow"
	// wxURL_SanBox_BatchQueryComment = wxBaseUrlSandbox + "billcommentsp/batchquerycomment"

	// 服务模式
	serviceTypeNormalDomestic      = 1 // 境内普通商户
	serviceTypeNormalAbroad        = 2 // 境外普通商户
	serviceTypeFacilitatorDomestic = 3 // 境内服务商
	serviceTypeFacilitatorAbroad   = 4 // 境外服务商
	serviceTypeBankServiceProvidor = 5 // 银行服务商

	// 支付类型
	tradeTypeApplet   = "JSAPI"    // 小程序支付
	tradeTypeJsApi    = "JSAPI"    // JSAPI支付
	tradeTypeApp      = "APP"      // APP支付
	tradeTypeH5       = "MWEB"     // H5支付
	tradeTypeNative   = "NATIVE"   // Native支付
	tradeTypeMicropay = "MICROPAY" // 付款码支付

	// 交易状态
	tradeStateSuccess    = "SUCCESS"    // 支付成功
	tradeStateRefund     = "REFUND"     // 转入退款
	tradeStateNotPay     = "NOTPAY"     // 未支付
	tradeStateClosed     = "CLOSED"     // 已关闭
	tradeStateRevoked    = "REVOKED"    // 已撤销(刷卡支付)
	tradeStateUserPaying = "USERPAYING" // 用户支付中
	tradeStatePayError   = "PAYERROR"   // 支付失败(其他原因,如银行返回失败)

	// 交易保障(MICROPAY)上报数据包的交易状态
	reportMicropayTradeStateOk     = "OK"     // 成功
	reportMicropayTradeStateFail   = "FAIL"   // 失败
	reportMicropayTradeStateCancel = "CANCLE" // 取消

	// 签名方式
	signTypeMD5        = "MD5" // 默认
	signTypeHmacSHA256 = "HMAC-SHA256"

	// 货币类型
	feeTypeCNY = "CNY" // 人民币

	// 指定支付方式
	limitPayNoCredit = "no_credit" // 指定不能使用信用卡支付

	// 压缩账单
	tarTypeGzip = "GZIP"

	// 电子发票
	receiptEnable = "Y" // 支付成功消息和支付详情页将出现开票入口

	// 代金券类型
	couponTypeCash   = "CASH"    // 充值代金券
	couponTypeNoCash = "NO_CASH" // 非充值优惠券

	// 账单类型
	billTypeAll            = "ALL"             // 返回当日所有订单信息,默认值
	billTypeSuccess        = "SUCCESS"         // 返回当日成功支付的订单
	billTypeRefund         = "REFUND"          // 返回当日退款订单
	billTypeRechargeRefund = "RECHARGE_REFUND" // 返回当日充值退款订单

	// 退款渠道
	refundChannelOriginal      = "ORIGINAL"       // 原路退款
	refundChannelBalance       = "BALANCE"        // 退回到余额
	refundChannelOtherBalance  = "OTHER_BALANCE"  // 原账户异常退到其他余额账户
	refundChannelOtherBankCard = "OTHER_BANKCARD" // 原银行卡异常退到其他银行卡

	// 退款状态
	refundStatusSuccess    = "SUCCESS"     // 退款成功
	refundStatusClose      = "REFUNDCLOSE" // 退款关闭
	refundStatusProcessing = "PROCESSING"  // 退款处理中
	refundStatusChange     = "CHANGE"      // 退款异常,退款到银行发现用户的卡作废或者冻结了,导致原路退款银行卡失败,可前往商户平台(pay.weixin.qq.com)-交易中心,手动处理此笔退款

	// 退款资金来源
	refundAccountRechargeFunds  = "REFUND_SOURCE_RECHARGE_FUNDS"  // 可用余额退款/基本账户
	refundAccountUnsettledFunds = "REFUND_SOURCE_UNSETTLED_FUNDS" // 未结算资金退款

	// 退款发起来源
	refundRequestSourceApi            = "API"             // API接口
	refundRequestSourceVendorPlatform = "VENDOR_PLATFORM" // 商户平台

	// 找零校验用户姓名选项
	checkNameTypeNoCheck    = "NO_CHECK"    //不校验真实姓名
	checkNameTypeForceCheck = "FORCE_CHECK" //强校验真实姓名

	// 返回结果
	responseSuccess = "SUCCESS" // 成功,通信标识或业务结果
	responseFail    = "FAIL"    // 失败,通信标识或业务结果

	// 返回消息
	responseMessageOk = "OK" // 返回成功信息

	// 错误代码,包括描述、支付状态、原因、解决方案
	errCodeAppIdMchIdNotMatch   = "APPID_MCHID_NOT_MATCH" // appid和mch_id不匹配 支付确认失败 appid和mch_id不匹配 请确认appid和mch_id是否匹配
	errCodeAppIdNotExist        = "APPID_NOT_EXIST"       // APPID不存在 支付确认失败 参数中缺少APPID 请检查APPID是否正确
	errCodeAuthCodeError        = "AUTH_CODE_ERROR"       // 授权码参数错误 支付确认失败 请求参数未按指引进行填写 每个二维码仅限使用一次,请刷新再试
	errCodeAuthCodeExpire       = "AUTHCODEEXPIRE"        // 二维码已过期,请用户在微信上刷新后再试 支付确认失败 用户的条码已经过期 请收银员提示用户,请用户在微信上刷新条码,然后请收银员重新扫码. 直接将错误展示给收银员
	errCodeAuthCodeInvalid      = "AUTH_CODE_INVALID"     // 授权码检验错误 支付确认失败 收银员扫描的不是微信支付的条码 请扫描微信支付被扫条码/二维码
	errCodeBankError            = "BANKERROR"             // 银行系统异常 支付结果未知 银行端超时 请立即调用被扫订单结果查询API,查询当前订单的不同状态,决定下一步的操作.
	errCodeBuyerMismatch        = "BUYER_MISMATCH"        // 支付帐号错误 支付确认失败 暂不支持同一笔订单更换支付方 请确认支付方是否相同
	errCodeInvalidRequest       = "INVALID_REQUEST"       // 无效请求 支付确认失败 商户系统异常导致,商户权限异常、重复请求支付、证书错误、频率限制等 请确认商户系统是否正常,是否具有相应支付权限,确认证书是否正确,控制频率
	errCodeInvalidTransactionId = "INVALID_TRANSACTIONID" // 无效transaction_id 请求参数未按指引进行填写 请求参数错误,检查原交易号是否存在或发起支付交易接口返回失败
	errCodeLackParams           = "LACK_PARAMS"           // 缺少参数 支付确认失败 缺少必要的请求参数 请检查参数是否齐全
	errCodeMchIdNotExist        = "MCHID_NOT_EXIST"       // MCHID不存在 支付确认失败 参数中缺少MCHID 请检查MCHID是否正确
	errCodeNoAuth               = "NOAUTH"                // 商户无权限 支付确认失败 商户没有开通被扫支付权限 请开通商户号权限.请联系产品或商务申请
	errCodeNotEnough            = "NOTENOUGH"             // 余额不足 支付确认失败 用户的零钱余额不足 请收银员提示用户更换当前支付的卡,然后请收银员重新扫码.建议:商户系统返回给收银台的提示为“用户余额不足.提示用户换卡支付”
	errCodeNotSuportCard        = "NOTSUPORTCARD"         // 不支持卡类型 支付确认失败 用户使用卡种不支持当前支付形式 请用户重新选择卡种 建议:商户系统返回给收银台的提示为“该卡不支持当前支付,提示用户换卡支付或绑新卡支付”
	errCodeNotUtf8              = "NOT_UTF8"              // 编码格式错误 支付确认失败 未使用指定编码格式 请使用UTF-8编码格式
	errCodeOrderClosed          = "ORDERCLOSED"           // 订单已关闭 支付确认失败 该订单已关 商户订单号异常,请重新下单支付
	errCodeOrderPaid            = "ORDERPAID"             // 订单已支付 支付确认失败 订单号重复 请确认该订单号是否重复支付,如果是新单,请使用新订单号提交
	errCodeOrderReversed        = "ORDERREVERSED"         // 订单已撤销 支付确认失败 当前订单已经被撤销 当前订单状态为“订单已撤销”,请提示用户重新支付
	errCodeOutTradeNoUsed       = "OUT_TRADE_NO_USED"     // 商户订单号重复 支付确认失败 同一笔交易不能多次提交 请核实商户订单号是否重复提交
	errCodeParamError           = "PARAM_ERROR"           // 参数错误 支付确认失败 请求参数未按指引进行填写 请根据接口返回的详细信息检查您的程序
	errCodePostDataEmpty        = "POST_DATA_EMPTY"       // post数据为空 post数据不能为空 请检查post数据是否为空
	errCodeRefundNotExist       = "REFUNDNOTEXIST"        // 退款订单查询失败 订单号错误或订单状态不正确 请检查订单号是否有误以及订单状态是否正确,如:未支付、已支付未退款
	errCodeRequirePostMethod    = "REQUIRE_POST_METHOD"   // 请使用post方法 支付确认失败 未使用post传递参数 请检查请求参数是否通过post方法提交
	errCodeSignError            = "SIGNERROR"             // 签名错误 支付确认失败 参数签名结果不正确 请检查签名参数和方法是否都符合签名算法要求
	errCodeSystemError          = "SYSTEMERROR"           // 接口返回错误 支付结果未知 系统超时 请立即调用被扫订单结果查询API,查询当前订单状态,并根据订单的状态决定下一步的操作.
	errCodeTradeError           = "TRADE_ERROR"           // 交易错误 支付确认失败 业务错误导致交易失败、用户账号异常、风控、规则限制等 请确认帐号是否存在异常
	errCodeUserPaying           = "USERPAYING"            // 用户支付中,需要输入密码 支付结果未知 该笔交易因为业务规则要求,需要用户输入支付密码. 等待5秒,然后调用被扫订单结果查询API,查询当前订单的不同状态,决定下一步的操作.
	errCodeXmlFormatError       = "XML_FORMAT_ERROR"      // XML格式错误 支付确认失败 XML格式错误 请检查XML参数格式是否正确

	// 是否关注公众账号
	isSubscribeYes = "Y" // 关注
	isSubscribeNo  = "N" // 未关注

	// 银行类型
	bankTypeIcbcDebit    = "ICBC_DEBIT"    // 工商银行(借记卡)
	bankTypeIcbcCredit   = "ICBC_CREDIT"   // 工商银行(信用卡)
	bankTypeAbcDebit     = "ABC_DEBIT"     // 农业银行(借记卡)
	bankTypeAbcCredit    = "ABC_CREDIT"    // 农业银行(信用卡)
	bankTypePsbcDebit    = "PSBC_DEBIT"    // 邮政储蓄银行(借记卡)
	bankTypePsbcCredit   = "PSBC_CREDIT"   // 邮政储蓄银行(信用卡)
	bankTypeCcbDebit     = "CCB_DEBIT"     // 建设银行(借记卡)
	bankTypeCcbCredit    = "CCB_CREDIT"    // 建设银行(信用卡)
	bankTypeCmbDebit     = "CMB_DEBIT"     // 招商银行(借记卡)
	bankTypeCmbCredit    = "CMB_CREDIT"    // 招商银行(信用卡)
	bankTypeBocDebit     = "BOC_DEBIT"     // 中国银行(借记卡)
	bankTypeBocCredit    = "BOC_CREDIT"    // 中国银行(信用卡)
	bankTypeCommDebit    = "COMM_DEBIT"    // 交通银行(借记卡)
	bankTypeCommCredit   = "COMM_CREDIT"   // 交通银行(信用卡)
	bankTypeSpdbDebit    = "SPDB_DEBIT"    // 浦发银行(借记卡)
	bankTypeSpdbCredit   = "SPDB_CREDIT"   // 浦发银行(信用卡)
	bankTypeGdbDebit     = "GDB_DEBIT"     // 广发银行(借记卡)
	bankTypeGdbCredit    = "GDB_CREDIT"    // 广发银行(信用卡)
	bankTypeCmbcDebit    = "CMBC_DEBIT"    // 民生银行(借记卡)
	bankTypeCmbcCredit   = "CMBC_CREDIT"   // 民生银行(信用卡)
	bankTypePabDebit     = "PAB_DEBIT"     // 平安银行(借记卡)
	bankTypePabCredit    = "PAB_CREDIT"    // 平安银行(信用卡)
	bankTypeCebDebit     = "CEB_DEBIT"     // 光大银行(借记卡)
	bankTypeCebCredit    = "CEB_CREDIT"    // 光大银行(信用卡)
	bankTypeCibDebit     = "CIB_DEBIT"     // 兴业银行(借记卡)
	bankTypeCibCredit    = "CIB_CREDIT"    // 兴业银行(信用卡)
	bankTypeCiticDebit   = "CITIC_DEBIT"   // 中信银行(借记卡)
	bankTypeCiticCredit  = "CITIC_CREDIT"  // 中信银行(信用卡)
	bankTypeBoshDebit    = "BOSH_DEBIT"    // 上海银行(借记卡)
	bankTypeBoshCredit   = "BOSH_CREDIT"   // 上海银行(信用卡)
	bankTypeCrbDebit     = "CRB_DEBIT"     // 华润银行(借记卡)
	bankTypeHzbDebit     = "HZB_DEBIT"     // 杭州银行(借记卡)
	bankTypeHzbCredit    = "HZB_CREDIT"    // 杭州银行(信用卡)
	bankTypeBsbDebit     = "BSB_DEBIT"     // 包商银行(借记卡)
	bankTypeBsbCredit    = "BSB_CREDIT"    // 包商银行(信用卡)
	bankTypeCqbDebit     = "CQB_DEBIT"     // 重庆银行(借记卡)
	bankTypeSdebDebit    = "SDEB_DEBIT"    // 顺德农商行(借记卡)
	bankTypeSzrcbDebit   = "SZRCB_DEBIT"   // 深圳农商银行(借记卡)
	bankTypeSzrcbCredit  = "SZRCB_CREDIT"  // 深圳农商银行(信用卡)
	bankTypeHrbbDebit    = "HRBB_DEBIT"    // 哈尔滨银行(借记卡)
	bankTypeBocdDebit    = "BOCD_DEBIT"    // 成都银行(借记卡)
	bankTypeGdnybDebit   = "GDNYB_DEBIT"   // 南粤银行(借记卡)
	bankTypeGdnybCredit  = "GDNYB_CREDIT"  // 南粤银行(信用卡)
	bankTypeGzcbDebit    = "GZCB_DEBIT"    // 广州银行(借记卡)
	bankTypeGzcbCredit   = "GZCB_CREDIT"   // 广州银行(信用卡)
	bankTypeJsbDebit     = "JSB_DEBIT"     // 江苏银行(借记卡)
	bankTypeJsbCredit    = "JSB_CREDIT"    // 江苏银行(信用卡)
	bankTypeNbcbDebit    = "NBCB_DEBIT"    // 宁波银行(借记卡)
	bankTypeNbcbCredit   = "NBCB_CREDIT"   // 宁波银行(信用卡)
	bankTypeNjcbDebit    = "NJCB_DEBIT"    // 南京银行(借记卡)
	bankTypeQhnxDebit    = "QHNX_DEBIT"    // 青海农信(借记卡)
	bankTypeOrdosbCredit = "ORDOSB_CREDIT" // 鄂尔多斯银行(信用卡)
	bankTypeOrdosbDebit  = "ORDOSB_DEBIT"  // 鄂尔多斯银行(借记卡)
	bankTypeBjrcbCredit  = "BJRCB_CREDIT"  // 北京农商(信用卡)
	bankTypeBhbDebit     = "BHB_DEBIT"     // 河北银行(借记卡)
	bankTypeBgzbDebit    = "BGZB_DEBIT"    // 贵州银行(借记卡)
	bankTypeBeebDebit    = "BEEB_DEBIT"    // 鄞州银行(借记卡)
	bankTypePzhccbDebit  = "PZHCCB_DEBIT"  // 攀枝花银行(借记卡)
	bankTypeQdccbCredit  = "QDCCB_CREDIT"  // 青岛银行(信用卡)
	bankTypeQdccbDebit   = "QDCCB_DEBIT"   // 青岛银行(借记卡)
	bankTypeShinhanDebit = "SHINHAN_DEBIT" // 新韩银行(借记卡)
	bankTypeQlbDebit     = "QLB_DEBIT"     // 齐鲁银行(借记卡)
	bankTypeQsbDebit     = "QSB_DEBIT"     // 齐商银行(借记卡)
	bankTypeZzbDebit     = "ZZB_DEBIT"     // 郑州银行(借记卡)
	bankTypeCcabDebit    = "CCAB_DEBIT"    // 长安银行(借记卡)
	bankTypeRzbDebit     = "RZB_DEBIT"     // 日照银行(借记卡)
	bankTypeScnxDebit    = "SCNX_DEBIT"    // 四川农信(借记卡)
	bankTypeBeebCredit   = "BEEB_CREDIT"   // 鄞州银行(信用卡)
	bankTypeSdrcuDebit   = "SDRCU_DEBIT"   // 山东农信(借记卡)
	bankTypeBczDebit     = "BCZ_DEBIT"     // 沧州银行(借记卡)
	bankTypeSjbDebit     = "SJB_DEBIT"     // 盛京银行(借记卡)
	bankTypeLnnxDebit    = "LNNX_DEBIT"    // 辽宁农信(借记卡)
	bankTypeJufengbDebit = "JUFENGB_DEBIT" // 临朐聚丰村镇银行(借记卡)
	bankTypeZzbCredit    = "ZZB_CREDIT"    // 郑州银行(信用卡)
	bankTypeJxnxbDebit   = "JXNXB_DEBIT"   // 江西农信(借记卡)
	bankTypeJzbDebit     = "JZB_DEBIT"     // 晋中银行(借记卡)
	bankTypeJzcbCredit   = "JZCB_CREDIT"   // 锦州银行(信用卡)
	// bankType                 = "JZCB_DEBIT"        // 锦州银行(借记卡)
	// bankType                 = "KLB_DEBIT"         // 昆仑银行(借记卡)
	// bankType                 = "KRCB_DEBIT"        // 昆山农商(借记卡)
	// bankType                 = "KUERLECB_DEBIT"    // 库尔勒市商业银行(借记卡)
	// bankType                 = "LJB_DEBIT"         // 龙江银行(借记卡)
	// bankType                 = "NYCCB_DEBIT"       // 南阳村镇银行(借记卡)
	// bankType                 = "LSCCB_DEBIT"       // 乐山市商业银行(借记卡)
	// bankType                 = "LUZB_DEBIT"        // 柳州银行(借记卡)
	// bankType                 = "LWB_DEBIT"         // 莱商银行(借记卡)
	// bankType                 = "LYYHB_DEBIT"       // 辽阳银行(借记卡)
	// bankType                 = "LZB_DEBIT"         // 兰州银行(借记卡)
	// bankType                 = "MINTAIB_CREDIT"    // 民泰银行(信用卡)
	// bankType                 = "MINTAIB_DEBIT"     // 民泰银行(借记卡)
	// bankType                 = "NCB_DEBIT"         // 宁波通商银行(借记卡)
	// bankType                 = "NMGNX_DEBIT"       // 内蒙古农信(借记卡)
	// bankType                 = "XAB_DEBIT"         // 西安银行(借记卡)
	// bankType                 = "WFB_CREDIT"        // 潍坊银行(信用卡)
	// bankType                 = "WFB_DEBIT"         // 潍坊银行(借记卡)
	// bankType                 = "WHB_CREDIT"        // 威海商业银行(信用卡)
	// bankType                 = "WHB_DEBIT"         // 威海市商业银行(借记卡)
	// bankType                 = "WHRC_CREDIT"       // 武汉农商(信用卡)
	// bankType                 = "WHRC_DEBIT"        // 武汉农商行(借记卡)
	// bankType                 = "WJRCB_DEBIT"       // 吴江农商行(借记卡)
	// bankType                 = "WLMQB_DEBIT"       // 乌鲁木齐银行(借记卡)
	// bankType                 = "WRCB_DEBIT"        // 无锡农商(借记卡)
	// bankType                 = "WZB_DEBIT"         // 温州银行(借记卡)
	// bankType                 = "XAB_CREDIT"        // 西安银行(信用卡)
	// bankType                 = "WEB_DEBIT"         // 微众银行(借记卡)
	// bankType                 = "XIB_DEBIT"         // 厦门国际银行(借记卡)
	// bankType                 = "XJRCCB_DEBIT"      // 新疆农信银行(借记卡)
	// bankType                 = "XMCCB_DEBIT"       // 厦门银行(借记卡)
	// bankType                 = "YNRCCB_DEBIT"      // 云南农信(借记卡)
	// bankType                 = "YRRCB_CREDIT"      // 黄河农商银行(信用卡)
	// bankType                 = "YRRCB_DEBIT"       // 黄河农商银行(借记卡)
	// bankType                 = "YTB_DEBIT"         // 烟台银行(借记卡)
	// bankType                 = "ZJB_DEBIT"         // 紫金农商银行(借记卡)
	// bankType                 = "ZJLXRB_DEBIT"      // 兰溪越商银行(借记卡)
	// bankType                 = "ZJRCUB_CREDIT"     // 浙江农信(信用卡)
	// bankType                 = "AHRCUB_DEBIT"      // 安徽省农村信用社联合社(借记卡)
	// bankType                 = "BCZ_CREDIT"        // 沧州银行(信用卡)
	// bankType                 = "SRB_DEBIT"         // 上饶银行(借记卡)
	// bankType                 = "ZYB_DEBIT"         // 中原银行(借记卡)
	// bankType                 = "ZRCB_DEBIT"        // 张家港农商行(借记卡)
	// bankType                 = "SRCB_CREDIT"       // 上海农商银行(信用卡)
	// bankType                 = "SRCB_DEBIT"        // 上海农商银行(借记卡)
	// bankType                 = "ZJTLCB_DEBIT"      // 浙江泰隆银行(借记卡)
	// bankType                 = "SUZB_DEBIT"        // 苏州银行(借记卡)
	// bankType                 = "SXNX_DEBIT"        // 山西农信(借记卡)
	// bankType                 = "SXXH_DEBIT"        // 陕西信合(借记卡)
	// bankType                 = "ZJRCUB_DEBIT"      // 浙江农信(借记卡)
	// bankType                 = "AE_CREDIT"         // AE(信用卡)
	// bankType                 = "TACCB_CREDIT"      // 泰安银行(信用卡)
	// bankType                 = "TACCB_DEBIT"       // 泰安银行(借记卡)
	// bankType                 = "TCRCB_DEBIT"       // 太仓农商行(借记卡)
	// bankType                 = "TJBHB_CREDIT"      // 天津滨海农商行(信用卡)
	// bankType                 = "TJBHB_DEBIT"       // 天津滨海农商行(借记卡)
	// bankType                 = "TJB_DEBIT"         // 天津银行(借记卡)
	// bankType                 = "TRCB_DEBIT"        // 天津农商(借记卡)
	// bankType                 = "TZB_DEBIT"         // 台州银行(借记卡)
	// bankType                 = "URB_DEBIT"         // 联合村镇银行(借记卡)
	// bankType                 = "DYB_CREDIT"        // 东营银行(信用卡)
	// bankType                 = "CSRCB_DEBIT"       // 常熟农商银行(借记卡)
	// bankType                 = "CZB_CREDIT"        // 浙商银行(信用卡)
	// bankType                 = "CZB_DEBIT"         // 浙商银行(借记卡)
	// bankType                 = "CZCB_CREDIT"       // 稠州银行(信用卡)
	// bankType                 = "CZCB_DEBIT"        // 稠州银行(借记卡)
	// bankType                 = "DANDONGB_CREDIT"   // 丹东银行(信用卡)
	// bankType                 = "DANDONGB_DEBIT"    // 丹东银行(借记卡)
	// bankType                 = "DLB_CREDIT"        // 大连银行(信用卡)
	// bankType                 = "DLB_DEBIT"         // 大连银行(借记卡)
	// bankType                 = "DRCB_CREDIT"       // 东莞农商银行(信用卡)
	// bankType                 = "DRCB_DEBIT"        // 东莞农商银行(借记卡)
	// bankType                 = "CSRCB_CREDIT"      // 常熟农商银行(信用卡)
	// bankType                 = "DYB_DEBIT"         // 东营银行(借记卡)
	// bankType                 = "DYCCB_DEBIT"       // 德阳银行(借记卡)
	// bankType                 = "FBB_DEBIT"         // 富邦华一银行(借记卡)
	// bankType                 = "FDB_DEBIT"         // 富滇银行(借记卡)
	// bankType                 = "FJHXB_CREDIT"      // 福建海峡银行(信用卡)
	// bankType                 = "FJHXB_DEBIT"       // 福建海峡银行(借记卡)
	// bankType                 = "FJNX_DEBIT"        // 福建农信银行(借记卡)
	// bankType                 = "FUXINB_DEBIT"      // 阜新银行(借记卡)
	// bankType                 = "BOCDB_DEBIT"       // 承德银行(借记卡)
	// bankType                 = "JSNX_DEBIT"        // 江苏农商行(借记卡)
	// bankType                 = "BOLFB_DEBIT"       // 廊坊银行(借记卡)
	// bankType                 = "CCAB_CREDIT"       // 长安银行(信用卡)
	// bankType                 = "CBHB_DEBIT"        // 渤海银行(借记卡)
	// bankType                 = "CDRCB_DEBIT"       // 成都农商银行(借记卡)
	// bankType                 = "BYK_DEBIT"         // 营口银行(借记卡)
	// bankType                 = "BOZ_DEBIT"         // 张家口市商业银行(借记卡)
	// bankType                 = "CFT"               // 零钱
	// bankType                 = "BOTSB_DEBIT"       // 唐山银行(借记卡)
	// bankType                 = "BOSZS_DEBIT"       // 石嘴山银行(借记卡)
	// bankType                 = "BOSXB_DEBIT"       // 绍兴银行(借记卡)
	// bankType                 = "BONX_DEBIT"        // 宁夏银行(借记卡)
	// bankType                 = "BONX_CREDIT"       // 宁夏银行(信用卡)
	// bankType                 = "GDHX_DEBIT"        // 广东华兴银行(借记卡)
	// bankType                 = "BOLB_DEBIT"        // 洛阳银行(借记卡)
	// bankType                 = "BOJX_DEBIT"        // 嘉兴银行(借记卡)
	// bankType                 = "BOIMCB_DEBIT"      // 内蒙古银行(借记卡)
	// bankType                 = "BOHN_DEBIT"        // 海南银行(借记卡)
	// bankType                 = "BOD_DEBIT"         // 东莞银行(借记卡)
	// bankType                 = "CQRCB_CREDIT"      // 重庆农商银行(信用卡)
	// bankType                 = "CQRCB_DEBIT"       // 重庆农商银行(借记卡)
	// bankType                 = "CQTGB_DEBIT"       // 重庆三峡银行(借记卡)
	// bankType                 = "BOD_CREDIT"        // 东莞银行(信用卡)
	// bankType                 = "CSCB_DEBIT"        // 长沙银行(借记卡)
	// bankType                 = "BOB_CREDIT"        // 北京银行(信用卡)
	// bankType                 = "GDRCU_DEBIT"       // 广东农信银行(借记卡)
	// bankType                 = "BOB_DEBIT"         // 北京银行(借记卡)
	// bankType                 = "HRXJB_DEBIT"       // 华融湘江银行(借记卡)
	// bankType                 = "HSBC_DEBIT"        // 恒生银行(借记卡)
	// bankType                 = "HSB_CREDIT"        // 徽商银行(信用卡)
	// bankType                 = "HSB_DEBIT"         // 徽商银行(借记卡)
	// bankType                 = "HUNNX_DEBIT"       // 湖南农信(借记卡)
	// bankType                 = "HUSRB_DEBIT"       // 湖商村镇银行(借记卡)
	// bankType                 = "HXB_CREDIT"        // 华夏银行(信用卡)
	// bankType                 = "HXB_DEBIT"         // 华夏银行(借记卡)
	// bankType                 = "HNNX_DEBIT"        // 河南农信(借记卡)
	// bankType                 = "BNC_DEBIT"         // 江西银行(借记卡)
	// bankType                 = "BNC_CREDIT"        // 江西银行(信用卡)
	// bankType                 = "BJRCB_DEBIT"       // 北京农商行(借记卡)
	// bankType                 = "JCB_DEBIT"         // 晋城银行(借记卡)
	// bankType                 = "JJCCB_DEBIT"       // 九江银行(借记卡)
	// bankType                 = "JLB_DEBIT"         // 吉林银行(借记卡)
	// bankType                 = "JLNX_DEBIT"        // 吉林农信(借记卡)
	// bankType                 = "JNRCB_DEBIT"       // 江南农商(借记卡)
	// bankType                 = "JRCB_DEBIT"        // 江阴农商行(借记卡)
	// bankType                 = "JSHB_DEBIT"        // 晋商银行(借记卡)
	// bankType                 = "HAINNX_DEBIT"      // 海南农信(借记卡)
	// bankType                 = "GLB_DEBIT"         // 桂林银行(借记卡)
	// bankType                 = "GRCB_CREDIT"       // 广州农商银行(信用卡)
	// bankType                 = "GRCB_DEBIT"        // 广州农商银行(借记卡)
	// bankType                 = "GSB_DEBIT"         // 甘肃银行(借记卡)
	// bankType                 = "GSNX_DEBIT"        // 甘肃农信(借记卡)
	// bankType                 = "GXNX_DEBIT"        // 广西农信(借记卡)
	bankTypeGycbCredit       = "GYCB_CREDIT"       // 贵阳银行(信用卡)
	bankTypeGycbDebit        = "GYCB_DEBIT"        // 贵阳银行(借记卡)
	bankTypeGznxDebit        = "GZNX_DEBIT"        // 贵州农信(借记卡)
	bankTypeHainnxCredit     = "HAINNX_CREDIT"     // 海南农信(信用卡)
	bankTypeHkbDebit         = "HKB_DEBIT"         // 汉口银行(借记卡)
	bankTypeHanabDebit       = "HANAB_DEBIT"       // 韩亚银行(借记卡)
	bankTypeHbcbCredit       = "HBCB_CREDIT"       // 湖北银行(信用卡)
	bankTypeHbcbDebit        = "HBCB_DEBIT"        // 湖北银行(借记卡)
	bankTypeHbnxCredit       = "HBNX_CREDIT"       // 湖北农信(信用卡)
	bankTypeHbnxDebit        = "HBNX_DEBIT"        // 湖北农信(借记卡)
	bankTypeHdcbDebit        = "HDCB_DEBIT"        // 邯郸银行(借记卡)
	bankTypeHebnxDebit       = "HEBNX_DEBIT"       // 河北农信(借记卡)
	bankTypeHfbDebit         = "HFB_DEBIT"         // 恒丰银行(借记卡)
	bankTypeHkbeaDebit       = "HKBEA_DEBIT"       // 东亚银行(借记卡)
	bankTypeJcbCredit        = "JCB_CREDIT"        // JCB(信用卡)
	bankTypeMasterCardCredit = "MASTERCARD_CREDIT" // MASTERCARD(信用卡)
	bankTypeVisaCredit       = "VISA_CREDIT"       // VISA(信用卡)
	bankTypeLqt              = "LQT"               // 零钱通
)
