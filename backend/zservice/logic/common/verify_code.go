package common

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/denghuo98/zzfast/internal/dao"
	"github.com/denghuo98/zzfast/web/zcache"
	"github.com/denghuo98/zzfast/web/zcaptcha"
	"github.com/denghuo98/zzfast/web/zsms"
	"github.com/denghuo98/zzfast/web/zverify"
	"github.com/denghuo98/zzfast/zconsts"
	commonSchema "github.com/denghuo98/zzfast/zschema/common"
	webSchema "github.com/denghuo98/zzfast/zschema/zweb"
	"github.com/denghuo98/zzfast/zservice"
)

type sVerifyCode struct{}

func init() {
	zservice.RegisterVerifyCode(NewVerifyCode())
}

func NewVerifyCode() *sVerifyCode {
	return &sVerifyCode{}
}

// Send 发送短信验证码。
func (s *sVerifyCode) Send(ctx g.Ctx, in *commonSchema.SendVerifyCodeInput) (out *commonSchema.SendVerifyCodeOutput, err error) {
	if in == nil {
		return nil, gerror.New("请求参数不能为空")
	}
	if in.Mobile == "" {
		return nil, gerror.New("手机号不能为空")
	}
	mobile := zverify.NormalizeMobile(in.Mobile)

	cfg, err := zservice.SystemConfig().GetVerifyCodeConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, gerror.New("验证码服务未启用")
	}
	sceneCfg := cfg.GetScene(in.Scene)
	if sceneCfg == nil {
		return nil, gerror.Newf("未知的验证码场景: %s", in.Scene)
	}
	if sceneCfg.RequireCaptcha && !zcaptcha.Verify(in.CaptchaId, in.CaptchaCode) {
		return nil, gerror.New("图形验证码错误")
	}

	smsCfg, err := zservice.SystemConfig().GetSmsConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !smsCfg.Enabled {
		return nil, gerror.New("短信服务未启用")
	}

	driver, err := zsms.New(smsCfg, smsCfg.Drive)
	if err != nil {
		g.Log().Errorf(ctx, "初始化短信驱动失败: %v", err)
		return nil, gerror.New("短信服务配置异常，请联系管理员")
	}

	ip := clientIP(ctx, cfg.TrustedProxies)
	policy := buildPolicy(cfg, sceneCfg)
	svc := zverify.New(zcache.Instance(), &smsSender{driver: driver})
	expiresIn, err := svc.Send(ctx, policy, mobile, ip)
	if err != nil {
		var sendErr *zverify.SendError
		if errors.As(err, &sendErr) {
			s.writeLog(ctx, in.Scene, mobile, ip, zconsts.VerifyLogStatusSendFailed, sendErr.Err.Error())
			g.Log().Warningf(ctx, "验证码发送失败 scene=%s mobile=%s err=%v", in.Scene, maskMobile(mobile), sendErr.Err)
			return nil, gerror.New("验证码发送失败，请稍后重试")
		}
		var rejectErr *zverify.RejectError
		if errors.As(err, &rejectErr) {
			g.Log().Warningf(ctx, "验证码发送被拒绝 scene=%s mobile=%s err=%v", in.Scene, maskMobile(mobile), err)
			return nil, err
		}
		g.Log().Warningf(ctx, "验证码发送异常 scene=%s mobile=%s err=%v", in.Scene, maskMobile(mobile), err)
		return nil, gerror.New("验证码服务异常，请稍后重试")
	}
	s.writeLog(ctx, in.Scene, mobile, ip, zconsts.VerifyLogStatusSent, "")

	return &commonSchema.SendVerifyCodeOutput{
		ExpiresIn: expiresIn,
		Cooldown:  int(policy.Cooldown.Seconds()),
		Mobile:    maskMobile(mobile),
	}, nil
}

// Verify 校验短信验证码，成功后立即失效。
func (s *sVerifyCode) Verify(ctx g.Ctx, in *commonSchema.VerifyVerifyCodeInput) (err error) {
	if in == nil {
		return gerror.New("请求参数不能为空")
	}

	cfg, err := zservice.SystemConfig().GetVerifyCodeConfig(ctx)
	if err != nil {
		return err
	}
	if !cfg.Enabled {
		return gerror.New("验证码服务未启用")
	}
	sceneCfg := cfg.GetScene(in.Scene)
	if sceneCfg == nil {
		return gerror.Newf("未知的验证码场景: %s", in.Scene)
	}

	mobile := zverify.NormalizeMobile(in.Mobile)
	ip := clientIP(ctx, cfg.TrustedProxies)
	svc := zverify.New(zcache.Instance(), nil)
	if err = svc.Verify(ctx, buildPolicy(cfg, sceneCfg), mobile, in.Code); err != nil {
		s.writeLog(ctx, in.Scene, mobile, ip, zconsts.VerifyLogStatusVerifyFailed, err.Error())
		var rejectErr *zverify.RejectError
		if errors.As(err, &rejectErr) {
			return err
		}
		g.Log().Warningf(ctx, "验证码校验异常 scene=%s mobile=%s err=%v", in.Scene, maskMobile(mobile), err)
		return gerror.New("验证码服务异常，请稍后重试")
	}
	s.writeLog(ctx, in.Scene, mobile, ip, zconsts.VerifyLogStatusVerified, "")
	return nil
}

type smsSender struct {
	driver zsms.Driver
}

func (s *smsSender) Send(ctx g.Ctx, mobile, template, code string) error {
	return s.driver.Send(ctx, zsms.Message{Mobile: mobile, Template: template, Code: code})
}

func buildPolicy(cfg *webSchema.VerifyCodeConfig, scene *webSchema.VerifyCodeSceneConfig) zverify.Policy {
	length := pickInt(scene.Length, cfg.Length, zconsts.DefaultVerifyCodeLength)
	ttl := pickInt(scene.TTL, cfg.TTL, zconsts.DefaultVerifyCodeTTL)
	cooldown := pickInt(scene.Cooldown, cfg.Cooldown, zconsts.DefaultVerifyCodeCooldown)
	maxAttempts := pickInt(scene.MaxAttempts, cfg.MaxAttempts, zconsts.DefaultVerifyCodeMaxAttempts)
	maxIpDaily := pickLimit(scene.MaxIpDaily, cfg.MaxIpDaily, zconsts.DefaultVerifyCodeMaxIpDaily)
	maxMobileDaily := pickLimit(scene.MaxMobileDaily, cfg.MaxMobileDaily, zconsts.DefaultVerifyCodeMaxMobDaily)

	return zverify.Policy{
		Scene:          scene.Scene,
		Template:       scene.Template,
		HashKey:        cfg.Secret,
		Length:         length,
		TTL:            time.Duration(ttl) * time.Second,
		Cooldown:       time.Duration(cooldown) * time.Second,
		MaxAttempts:    maxAttempts,
		MaxIpDaily:     maxIpDaily,
		MaxMobileDaily: maxMobileDaily,
	}
}

// pickLimit 解析每日发送上限：场景 > 全局 > 默认；负数表示不限。
func pickLimit(scene, global, fallback int) int {
	if scene < 0 {
		return -1
	}
	if scene > 0 {
		return scene
	}
	if global < 0 {
		return -1
	}
	if global > 0 {
		return global
	}
	return fallback
}

func pickInt(values ...int) int {
	for _, v := range values {
		if v > 0 {
			return v
		}
	}
	return 0
}

func clientIP(ctx g.Ctx, trustedProxies []string) string {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	remote := r.GetRemoteIp()
	if !isTrustedProxy(remote, trustedProxies) {
		if isPrivateAddress(remote) {
			// 内网地址但不是已配置的可信代理，无法确认真实客户端 IP，跳过 IP 限流。
			return ""
		}
		return remote
	}
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		for i := len(parts) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(parts[i])
			if net.ParseIP(ip) == nil || isTrustedProxy(ip, trustedProxies) {
				continue
			}
			return ip
		}
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); net.ParseIP(realIP) != nil && !isTrustedProxy(realIP, trustedProxies) {
		return realIP
	}
	// 无法确定真实客户端 IP 时不做 IP 限流，避免所有用户共享代理 IP 的每日额度。
	return ""
}

// isTrustedProxy 仅当请求来自本机/内网代理时才信任 X-Forwarded-For，避免公网请求伪造 IP 绕过限流。
var cgnatRange = func() *net.IPNet {
	_, network, _ := net.ParseCIDR("100.64.0.0/10")
	return network
}()

func isPrivateAddress(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	return parsed.IsLoopback() || parsed.IsPrivate() || cgnatRange.Contains(parsed)
}

func isTrustedProxy(ip string, extra []string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	if parsed.IsLoopback() {
		return true
	}
	for _, item := range extra {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if strings.Contains(item, "/") {
			if _, network, err := net.ParseCIDR(item); err == nil && network.Contains(parsed) {
				return true
			}
			continue
		}
		if trusted := net.ParseIP(item); trusted != nil && trusted.Equal(parsed) {
			return true
		}
	}
	return false
}

func (s *sVerifyCode) writeLog(ctx g.Ctx, scene, mobile, ip string, status int, errMsg string) {
	mobile = truncateRunes(mobile, 32)
	scene = truncateRunes(scene, 64)
	ip = truncateRunes(ip, 64)
	errMsg = truncateRunes(errMsg, 1000)
	cols := dao.SysVerifyCodeLog.Columns()
	_, err := dao.SysVerifyCodeLog.Ctx(ctx).Data(g.Map{
		cols.Scene:     scene,
		cols.Channel:   zconsts.VerifyChannelSms,
		cols.Mobile:    mobile,
		cols.Ip:        ip,
		cols.Status:    status,
		cols.ErrMsg:    errMsg,
		cols.CreatedAt: gtime.Now(),
	}).Insert()
	if err != nil {
		g.Log().Warningf(ctx, "写入验证码日志失败: %v", err)
	}
}

func truncateRunes(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit])
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
