package common

import (
	commonSchema "github.com/denghuo98/zzfast/zschema/common"
	"github.com/gogf/gf/v2/frame/g"
)

// VerifyCodeSendReq 发送验证码（匿名场景，如登录）
type VerifyCodeSendReq struct {
	g.Meta `path:"/site/verify-code/send" method:"post" tags:"SYS-00-系统管理" summary:"发送验证码"`
	commonSchema.SendVerifyCodeInput
}

type VerifyCodeSendRes struct {
	commonSchema.SendVerifyCodeOutput
}
