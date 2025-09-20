package gowe

import (
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"
)

// ============ 通知结构 ============
type PaymentNotification struct {
	ID           string `json:"id"`
	CreateTime   string `json:"create_time"`
	EventType    string `json:"event_type"`
	ResourceType string `json:"resource_type"`
	Resource     struct {
		Algorithm      string `json:"algorithm"`
		Ciphertext     string `json:"ciphertext"`
		AssociatedData string `json:"associated_data"`
		Nonce          string `json:"nonce"`
	} `json:"resource"`
}

// 解密后的支付结果
type DecryptedPaymentResult struct {
	AppID         string `json:"appid"`
	MchID         string `json:"mchid"`
	OutTradeNo    string `json:"out_trade_no"`
	TransactionID string `json:"transaction_id"`
	TradeType     string `json:"trade_type"`
	TradeState    string `json:"trade_state"`
	SuccessTime   string `json:"success_time"`
	Amount        struct {
		Total    int    `json:"total"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Payer struct {
		OpenID string `json:"openid"`
	} `json:"payer"`
}

type CallbackResponse struct {
	Code    int                    `json:"code"`
	Massage string                 `json:"massage"`
	Data    DecryptedPaymentResult `json:"data"`
}

// ----------------微信支付回调-----------------
func WechatPayCallback(ctx context.Context, wxPayConfig IWxPayConfig, body []byte) CallbackResponse {
	callbackResponse := CallbackResponse{
		Code:    0,
		Massage: "ok",
	}
	// 3. 解析通知
	var notification PaymentNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		//fmt.Println("解析通知失败")
		callbackResponse.Code = 1
		callbackResponse.Massage = "解析通知失败"
		return callbackResponse
	}

	// 4. 解密
	decryptedData, err := DecryptResource(ctx, wxPayConfig, notification.Resource)
	if err != nil {
		callbackResponse.Code = 1
		callbackResponse.Massage = "解密失败"
		return callbackResponse
	}

	// 5. 解析支付结果
	var paymentResult DecryptedPaymentResult
	if err := json.Unmarshal(decryptedData, &paymentResult); err != nil {
		callbackResponse.Code = 1
		callbackResponse.Massage = "解析支付结果失败"
		return callbackResponse
	}

	// 6. 业务处理
	payment_bool, err := ProcessPayment(ctx, wxPayConfig, paymentResult)
	if err != nil {
		callbackResponse.Code = 1
		callbackResponse.Massage = "业务处理失败"
		return callbackResponse
	}

	if payment_bool {
		callbackResponse.Data = paymentResult
		return callbackResponse
	} else {
		callbackResponse.Code = 1
		callbackResponse.Massage = "error"
		return callbackResponse
	}
}

// ---------------------验证微信签名------------
func VerifyWechatSignature(ctx context.Context, wxPayConfig IWxPayConfig, timestamp string, nonce string, signature string, serial string, body []byte) error {
	if timestamp == "" || nonce == "" || signature == "" || serial == "" {
		return errors.New("缺少签名头")
	}

	// 构造验签串
	signMessage := fmt.Sprintf("%s\n%s\n%s\n", timestamp, nonce, string(body))
	//fmt.Println("验签串: \n", signMessage)
	// 1. 读取商户 API 证书
	keyBytes, err := os.ReadFile(wxPayConfig.GetWechatPayCertificate(ctx))
	if err != nil {
		return fmt.Errorf("读取证书文件失败: %v", err)
	}

	// 解析证书
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return errors.New("证书解析失败")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	// 校验证书序列号
	sprintf := fmt.Sprintf("%X", cert.SerialNumber)
	//fmt.Println("解析出来的 ", sprintf)
	//fmt.Println("携带的 ", serial)
	if sprintf != serial {
		return errors.New("证书序列号不匹配")
	}

	// 提取公钥
	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("无效公钥")
	}

	// Base64 解码签名
	//fmt.Println("签名: ", signature)
	signBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	// 验签
	hashed := sha256.Sum256([]byte(signMessage))
	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signBytes); err != nil {
		return errors.New("签名校验失败")
	}

	// 时间戳校验（可选）
	if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
		if time.Since(t) > 5*time.Minute {
			return errors.New("时间戳过期")
		}
	}

	return nil
}

// 解密资源
func DecryptResource(ctx context.Context, wxPayConfig IWxPayConfig, res struct {
	Algorithm      string `json:"algorithm"`
	Ciphertext     string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	Nonce          string `json:"nonce"`
}) ([]byte, error) {
	if res.Algorithm != "AEAD_AES_256_GCM" {
		return nil, errors.New("不支持的加密算法")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(res.Ciphertext)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher([]byte(wxPayConfig.GetMchAPIv3Key(ctx)))
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesgcm.Open(nil, []byte(res.Nonce), ciphertext, []byte(res.AssociatedData))
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// 业务处理
func ProcessPayment(ctx context.Context, wxPayConfig IWxPayConfig, result DecryptedPaymentResult) (bool, error) {
	if result.MchID != wxPayConfig.GetMchID(ctx) {
		return false, errors.New("商户号不匹配")
	}
	if result.TradeState != "SUCCESS" {
		return false, fmt.Errorf("支付未成功: %s", result.TradeState)
	}
	// 幂等性校验 & 订单更新
	fmt.Printf("支付成功：订单号=%s, 微信单号=%s, 金额=%d分, 用户=%s\n",
		result.OutTradeNo, result.TransactionID, result.Amount.Total, result.Payer.OpenID)
	return true, nil
}
