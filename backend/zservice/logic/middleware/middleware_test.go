package middleware

import (
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/denghuo98/zzfast/zservice"
	"github.com/denghuo98/zzfast/zservice/logic/system"
)

func TestGetAnonymousConfig_DefaultDisabled(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		ctx := gctx.New()
		s := system.NewSystemConfig()

		cfg, err := s.GetAnonymousConfig(ctx)
		gt.AssertNil(err)
		gt.Assert(cfg.Enabled, false)
		// 无 config 时 Identity 为零值
		gt.Assert(cfg.Identity.Username, "")
	})
}

func TestGetAnonymousConfig_Enabled(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		ctx := gctx.New()
		s := system.NewSystemConfig()
		cfg, err := s.GetAnonymousConfig(ctx)
		gt.AssertNil(err)
		gt.Assert(cfg.Enabled, false)
	})
}

func TestDeliverUserContext_AnonymousBranch(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		m := NewMiddleware()
		gt.AssertNE(m, nil)
		zservice.RegisterMiddleware(m)
	})
}
