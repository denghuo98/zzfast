package admin

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/denghuo98/zzfast/web/zcontext"
	adminApi "github.com/denghuo98/zzfast/zapi/admin"
	"github.com/denghuo98/zzfast/zservice"
)

type cVerifyCode struct{}

var VerifyCode = cVerifyCode{}

// Send 发送验证码（登录后场景）。
func (c *cVerifyCode) Send(ctx g.Ctx, req *adminApi.VerifyCodeSendReq) (res *adminApi.VerifyCodeSendRes, err error) {
	cfg, err := zservice.SystemConfig().GetVerifyCodeConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, gerror.New("验证码服务未启用")
	}

	sceneCfg := cfg.GetScene(req.Scene)
	if sceneCfg == nil {
		return nil, gerror.New("未知的验证码场景")
	}
	if !sceneCfg.AuthOnly {
		return nil, gerror.New("当前场景不支持登录后发送")
	}

	user := zcontext.GetUser(ctx)
	if user == nil {
		return nil, gerror.New("用户身份异常，请重新登录")
	}

	mobile := req.Mobile
	if !sceneCfg.AllowTargetInput || mobile == "" {
		mobile = user.Mobile
	}
	if mobile == "" {
		return nil, gerror.New("当前账号未绑定手机号")
	}

	in := req.SendVerifyCodeInput
	in.Mobile = mobile
	out, err := zservice.VerifyCode().Send(ctx, &in)
	if err != nil {
		return nil, err
	}
	res = new(adminApi.VerifyCodeSendRes)
	res.SendVerifyCodeOutput = *out
	return res, nil
}
