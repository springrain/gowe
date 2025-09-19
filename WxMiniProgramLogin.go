package gowe

import (
	"context"
	"encoding/json"
	"fmt"
)

type PhoneResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    PhoneInfo `json:"data"`
}

// 手机号信息结构
type PhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}

func WxGetPhoneInfo(ctx context.Context, SessionKey string, Iv string, EncryptedData string) *PhoneResponse {
	var err error
	resp := &PhoneResponse{
		Code:    0,
		Message: "success",
	}
	// 2. 解密手机号数据
	decryptedData, err := AesDecrypt(EncryptedData, SessionKey, Iv)
	if err != nil {
		resp.Code = 1
		resp.Message = fmt.Sprintf("解密失败: %v", err)
		return resp
	}
	// 3. 解析手机号信息
	var phoneInfo PhoneInfo
	if err := json.Unmarshal(decryptedData, &phoneInfo); err != nil {
		resp.Code = 1
		resp.Message = fmt.Sprintf("解析手机号信息失败: %v", err)
		return resp
	}

	// 4. 验证手机号格式
	if !ValidatePhoneNumber(phoneInfo.PurePhoneNumber) {
		resp.Code = 1
		resp.Message = "无效的手机号格式"
		return resp
	}

	resp.Data = phoneInfo
	return resp
}
