package zverify

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gcache"
)

type fakeSender struct {
	code  string
	calls int
}

func (s *fakeSender) Send(_ context.Context, _ string, _ string, code string) error {
	s.code = code
	s.calls++
	return nil
}

func newTestService(sender Sender) (*Service, *time.Time) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	svc := New(gcache.New(), sender)
	svc.now = func() time.Time { return now }
	return svc, &now
}

func testPolicy() Policy {
	return Policy{
		Scene:          "login",
		Template:       "SMS_TEST",
		HashKey:        "test-secret",
		Length:         6,
		TTL:            time.Minute * 5,
		Cooldown:       time.Second * 60,
		MaxAttempts:    3,
		MaxIpDaily:     10,
		MaxMobileDaily: 10,
	}
}

func TestSendAndVerify(t *testing.T) {
	sender := &fakeSender{}
	svc, _ := newTestService(sender)
	policy := testPolicy()

	expires, err := svc.Send(context.Background(), policy, "13800000001", "127.0.0.1")
	if err != nil {
		t.Fatalf("send failed: %v", err)
	}
	if expires != 300 {
		t.Fatalf("unexpected expires: %d", expires)
	}
	if len(sender.code) != 6 {
		t.Fatalf("unexpected code length: %d", len(sender.code))
	}

	if err = svc.Verify(context.Background(), policy, "13800000001", sender.code); err != nil {
		t.Fatalf("verify failed: %v", err)
	}

	// 验证码一次性使用
	if err = svc.Verify(context.Background(), policy, "13800000001", sender.code); err == nil {
		t.Fatal("expected verify to fail after code consumed")
	}
}

func TestCooldown(t *testing.T) {
	sender := &fakeSender{}
	svc, now := newTestService(sender)
	policy := testPolicy()

	if _, err := svc.Send(context.Background(), policy, "13800000002", "127.0.0.1"); err != nil {
		t.Fatalf("send failed: %v", err)
	}
	if _, err := svc.Send(context.Background(), policy, "13800000002", "127.0.0.1"); err == nil || !strings.Contains(err.Error(), "频繁") {
		t.Fatalf("expected cooldown error, got: %v", err)
	}

	*now = now.Add(61 * time.Second)
	if _, err := svc.Send(context.Background(), policy, "13800000002", "127.0.0.1"); err != nil {
		t.Fatalf("send after cooldown failed: %v", err)
	}
}

func TestMaxAttempts(t *testing.T) {
	sender := &fakeSender{}
	svc, _ := newTestService(sender)
	policy := testPolicy()
	ctx := context.Background()

	if _, err := svc.Send(ctx, policy, "13800000003", "127.0.0.1"); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := svc.Verify(ctx, policy, "13800000003", "000000"); err == nil {
			t.Fatal("expected wrong code error")
		}
	}
	if err := svc.Verify(ctx, policy, "13800000003", sender.code); err == nil || !strings.Contains(err.Error(), "次数过多") {
		t.Fatalf("expected max attempts error, got: %v", err)
	}
}

func TestExpire(t *testing.T) {
	sender := &fakeSender{}
	svc, now := newTestService(sender)
	policy := testPolicy()
	ctx := context.Background()

	if _, err := svc.Send(ctx, policy, "13800000004", "127.0.0.1"); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	*now = now.Add(6 * time.Minute)
	if err := svc.Verify(ctx, policy, "13800000004", sender.code); err == nil || !strings.Contains(err.Error(), "过期") {
		t.Fatalf("expected expired error, got: %v", err)
	}
}

func TestDailyLimit(t *testing.T) {
	sender := &fakeSender{}
	svc, now := newTestService(sender)
	policy := testPolicy()
	policy.MaxMobileDaily = 1
	ctx := context.Background()

	if _, err := svc.Send(ctx, policy, "13800000005", "127.0.0.1"); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	*now = now.Add(61 * time.Second)
	if _, err := svc.Send(ctx, policy, "13800000005", "127.0.0.1"); err == nil || !strings.Contains(err.Error(), "次数过多") {
		t.Fatalf("expected daily limit error, got: %v", err)
	}
}

func TestVerifyTrimsCode(t *testing.T) {
	sender := &fakeSender{}
	svc, _ := newTestService(sender)
	policy := testPolicy()
	ctx := context.Background()

	if _, err := svc.Send(ctx, policy, "13800000006", "127.0.0.1"); err != nil {
		t.Fatalf("send failed: %v", err)
	}
	if err := svc.Verify(ctx, policy, "13800000006", " "+sender.code+" "); err != nil {
		t.Fatalf("verify with padded code failed: %v", err)
	}
}

func TestInvalidMobile(t *testing.T) {
	svc, _ := newTestService(&fakeSender{})
	if _, err := svc.Send(context.Background(), testPolicy(), "abc", "127.0.0.1"); err == nil {
		t.Fatal("expected invalid mobile error")
	}
}
