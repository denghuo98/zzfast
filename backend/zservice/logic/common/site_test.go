package common

import (
	"testing"

	"github.com/gogf/gf/v2/encoding/gbase64"

	"github.com/denghuo98/zzfast/web/zencrypt"
)

// testKey 测试用 AES 密钥（16 字节）
var testKey = []byte("1234567890abcdef")

// encryptPassword 将明文密码加密为前端传输格式
// 先 AES-GCM 加密，再 base64 编码
func encryptPassword(plainText string) string {
	encrypted, err := zencrypt.AesGCMEncrypt([]byte(plainText), testKey)
	if err != nil {
		panic(err)
	}
	return gbase64.EncodeToString(encrypted)
}

func TestVerifyPassword(t *testing.T) {
	s := NewCommonSite()
	plainPassword := "testPassword123"

	correctHash, err := zencrypt.HashPassword(plainPassword)
	if err != nil {
		t.Fatalf("生成 bcrypt 哈希失败: %v", err)
	}

	tests := []struct {
		name      string
		input     string
		hash      string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "正确的密码验证",
			input:     encryptPassword(plainPassword),
			hash:      correctHash,
			wantError: false,
		},
		{
			name:      "错误的密码",
			input:     encryptPassword("wrongPassword"),
			hash:      correctHash,
			wantError: true,
			errorMsg:  "用户密码错误",
		},
		{
			name:      "无效的 base64 编码",
			input:     "!!!invalid-base64!!!",
			hash:      correctHash,
			wantError: true,
		},
		{
			name:      "无效的 AES 加密数据",
			input:     gbase64.EncodeToString([]byte("invalid aes data")),
			hash:      correctHash,
			wantError: true,
		},
		{
			name:      "空密码",
			input:     encryptPassword(""),
			hash:      correctHash,
			wantError: true,
			errorMsg:  "用户密码错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.verifyPassword(tt.input, testKey, tt.hash)
			if tt.wantError {
				if err == nil {
					t.Errorf("期望返回错误，但返回了 nil")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("错误信息不匹配，期望: %s，实际: %s", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("不期望返回错误，但返回了: %v", err)
				}
			}
		})
	}
}
