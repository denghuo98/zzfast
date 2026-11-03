package zservice

import (
	"github.com/gogf/gf/v2/frame/g"

	commonSchema "github.com/denghuo98/zzfast/zschema/common"
)

// IVerifyCode 短信验证码服务。
//
// 上层业务通过 Send 发送验证码，并在登录、高危操作等动作前调用 Verify 校验。
// 验证码校验成功后立即失效，具体场景、模板、有效期与限流策略由配置决定。
type IVerifyCode interface {
	// Send 发送验证码。
	Send(ctx g.Ctx, in *commonSchema.SendVerifyCodeInput) (out *commonSchema.SendVerifyCodeOutput, err error)

	// Verify 校验验证码，成功后立即失效。
	Verify(ctx g.Ctx, in *commonSchema.VerifyVerifyCodeInput) error
}

var (
	localVerifyCode IVerifyCode
)

func VerifyCode() IVerifyCode {
	if localVerifyCode == nil {
		panic("VerifyCode is not initialized, please register it first")
	}
	return localVerifyCode
}

func RegisterVerifyCode(i IVerifyCode) {
	localVerifyCode = i
}
