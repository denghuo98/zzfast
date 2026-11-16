package zwechatpay

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	webSchema "github.com/denghuo98/zzfast/zschema/zweb"
)

func TestValidate(t *testing.T) {
	if err := Validate(nil); err != nil {
		t.Fatalf("nil config should be valid when disabled: %v", err)
	}
	if err := Validate(&webSchema.WechatPayConfig{}); err != nil {
		t.Fatalf("disabled config should be valid: %v", err)
	}

	cfg := &webSchema.WechatPayConfig{Enabled: true}
	if err := Validate(cfg); err == nil {
		t.Fatal("expected missing field error")
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	pemKey := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))

	valid := &webSchema.WechatPayConfig{
		Enabled:     true,
		AppId:       "wx-test",
		MchId:       "1900000001",
		MchSerialNo: "serial",
		ApiV3Key:    "12345678901234567890123456789012",
		PrivateKey:  pemKey,
		NotifyUrl:   "https://example.com/pay/wechat/notify",
	}
	if err = Validate(valid); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}

	valid.ApiV3Key = "short"
	if err = Validate(valid); err == nil {
		t.Fatal("expected apiV3Key length error")
	}
}
