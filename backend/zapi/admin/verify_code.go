package admin

import (
	commonSchema "github.com/denghuo98/zzfast/zschema/common"
	"github.com/gogf/gf/v2/frame/g"
)

// VerifyCodeSendReq 发送验证码（登录后场景，如高危操作、换绑手机）
type VerifyCodeSendReq struct {
	g.Meta `path:"/member/verify-code/send" method:"post" tags:"SYS-02-用户管理" summary:"发送验证码"`
	commonSchema.SendVerifyCodeInput
}

type VerifyCodeSendRes struct {
	commonSchema.SendVerifyCodeOutput
}
