package common

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	commonApi "github.com/denghuo98/zzfast/zapi/common"
	"github.com/denghuo98/zzfast/zservice"
)

type cVerifyCode struct{}

var VerifyCode = cVerifyCode{}

// Send 发送验证码（匿名场景）。
func (c *cVerifyCode) Send(ctx g.Ctx, req *commonApi.VerifyCodeSendReq) (res *commonApi.VerifyCodeSendRes, err error) {
	cfg, err := zservice.SystemConfig().GetVerifyCodeConfig(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "读取验证码配置失败: %v", err)
		return nil, gerror.New("验证码服务暂不可用")
	}
	if !cfg.Enabled {
		return nil, gerror.New("验证码服务未启用")
	}

	sceneCfg := cfg.GetScene(req.Scene)
	if sceneCfg == nil {
		return nil, gerror.New("未知的验证码场景")
	}
	if sceneCfg.AuthOnly {
		return nil, gerror.New("当前场景需要登录后操作")
	}

	out, err := zservice.VerifyCode().Send(ctx, &req.SendVerifyCodeInput)
	if err != nil {
		return nil, err
	}
	res = new(commonApi.VerifyCodeSendRes)
	res.SendVerifyCodeOutput = *out
	return res, nil
}
