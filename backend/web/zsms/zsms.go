// Package zsms 提供短信发送能力。
//
// 发送动作（驱动选择、签名、模板参数、超时）由框架封装，
// 业务方只需要提供手机号、模板 ID 和验证码。
package zsms

import (
	"context"
	"fmt"
	"strings"

	webSchema "github.com/denghuo98/zzfast/zschema/zweb"
)

const (
	DriveMock    = "mock"
	DriveAliyun  = "aliyun"
	DriveTencent = "tencent"
)

// Message 待发送短信内容。
type Message struct {
	Mobile   string
	Template string
	Code     string
}

// Driver 短信驱动。
type Driver interface {
	Send(ctx context.Context, msg Message) error
}

var config *webSchema.SmsConfig

// SetConfig 由系统配置加载流程注入。
func SetConfig(c *webSchema.SmsConfig) {
	config = c
}

// GetConfig 获取当前短信配置（可能为 nil）。
func GetConfig() *webSchema.SmsConfig {
	return config
}

// Enabled 短信服务是否已启用。
func Enabled() bool {
	return config != nil && config.Enabled
}

// New 按驱动名称创建短信驱动。cfg 为空时回退到 SetConfig 注入的全局配置。
func New(cfg *webSchema.SmsConfig, drive ...string) (Driver, error) {
	if cfg == nil {
		cfg = GetConfig()
	}
	if cfg == nil {
		return nil, fmt.Errorf("zsms: 短信配置未初始化")
	}

	name := cfg.Drive
	if len(drive) > 0 && drive[0] != "" {
		name = drive[0]
	}
	if name == "" {
		name = DriveMock
	}

	switch name {
	case DriveMock:
		return &mockDriver{}, nil
	case DriveAliyun:
		return &aliyunDriver{cfg: cfg}, nil
	case DriveTencent:
		return &tencentDriver{cfg: cfg}, nil
	default:
		return nil, fmt.Errorf("zsms: 不支持的短信驱动 %q", name)
	}
}

func timeout(cfg *webSchema.SmsConfig) int {
	if cfg.Timeout > 0 {
		return cfg.Timeout
	}
	return 10
}

func signName(cfg *webSchema.SmsConfig, provider string) string {
	if provider == DriveAliyun && cfg.AliyunSignName != "" {
		return cfg.AliyunSignName
	}
	if provider == DriveTencent && cfg.TencentSignName != "" {
		return cfg.TencentSignName
	}
	return cfg.SignName
}

func strValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func e164(mobile string) string {
	mobile = strings.TrimSpace(mobile)
	mobile = strings.ReplaceAll(mobile, " ", "")
	mobile = strings.ReplaceAll(mobile, "-", "")
	if strings.HasPrefix(mobile, "00") {
		return "+" + strings.TrimPrefix(mobile, "00")
	}
	if strings.HasPrefix(mobile, "+") {
		return mobile
	}
	if strings.HasPrefix(mobile, "86") && len(mobile) == 13 {
		mobile = strings.TrimPrefix(mobile, "86")
	}
	if len(mobile) == 11 && strings.HasPrefix(mobile, "1") {
		return "+86" + mobile
	}
	return mobile
}

func maskMobile(mobile string) string {
	if len(mobile) <= 7 {
		return strings.Repeat("*", len(mobile))
	}
	if len(mobile) < 11 {
		return mobile[:1] + strings.Repeat("*", len(mobile)-3) + mobile[len(mobile)-2:]
	}
	return mobile[:3] + "****" + mobile[len(mobile)-4:]
}
