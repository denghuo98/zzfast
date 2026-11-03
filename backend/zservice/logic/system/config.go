package system

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gmode"

	"github.com/denghuo98/zzfast/web/zcache"
	"github.com/denghuo98/zzfast/web/zmail"
	"github.com/denghuo98/zzfast/web/zsms"
	"github.com/denghuo98/zzfast/web/zstorager"
	"github.com/denghuo98/zzfast/web/ztoken"
	"github.com/denghuo98/zzfast/zconsts"
	webSchema "github.com/denghuo98/zzfast/zschema/zweb"
	"github.com/denghuo98/zzfast/zservice"
)

type sSystemConfig struct {
}

func init() {
	zservice.RegisterSystemConfig(NewSystemConfig())
}

func NewSystemConfig() *sSystemConfig {
	return &sSystemConfig{}
}

func (s *sSystemConfig) LoadConfig(ctx g.Ctx) (err error) {
	systemConfig := &webSchema.SystemConfig{}
	v, err := g.Cfg().Get(ctx, "system")
	if err != nil {
		return err
	}
	if v != nil {
		if err = v.Struct(systemConfig); err != nil {
			return err
		}
	}

	// 默认生产环境
	if systemConfig.Mode == "" {
		systemConfig.Mode = "prod"
	}
	s.SetGFMode(ctx, systemConfig.Mode)

	// 缓存配置
	cacheCfg, err := s.GetCacheConfig(ctx)
	if err != nil {
		return err
	}
	zcache.SetAdapter(ctx, cacheCfg)

	// token 认证配置
	tokenCfg, err := s.GetTokenConfig(ctx)
	if err != nil {
		return err
	}
	if err := s.validateTokenConfig(ctx, tokenCfg); err != nil {
		return err
	}
	ztoken.SetConfig(tokenCfg)

	// 上传附件配置
	uploadCfg, err := s.GetUploadConfig(ctx)
	if err != nil {
		return err
	}
	zstorager.SetConfig(uploadCfg)

	// 邮件通知配置（未启用时允许为空）
	mailCfg, err := s.GetMailConfig(ctx)
	if err != nil {
		return err
	}
	zmail.SetConfig(mailCfg)

	// 短信发送配置（验证码通道）
	smsCfg, err := s.GetSmsConfig(ctx)
	if err != nil {
		return err
	}
	zsms.SetConfig(smsCfg)

	// 验证码策略
	verifyCfg, err := s.GetVerifyCodeConfig(ctx)
	if err != nil {
		return err
	}
	if verifyCfg.Enabled && !smsCfg.Enabled {
		return gerror.New("启用验证码服务时必须同时启用短信服务（sms.enabled = true）")
	}
	if verifyCfg.Enabled && len(verifyCfg.TrustedProxies) == 0 {
		g.Log().Warning(ctx, "system.verifyCode.trustedProxies 未配置：若部署在 nginx/SLB 后，将无法获取真实客户端 IP，IP 日限额会被跳过")
	}

	// 请求加密密钥（强校验，缺失或长度不合法直接失败）
	if _, err := s.GetRequestEncryptKey(ctx); err != nil {
		return err
	}

	// 初始化超级管理员
	if err := zservice.CommonSite().InitSuperAdmin(ctx); err != nil {
		return err
	}

	return nil
}

func (s *sSystemConfig) GetSuperAdmin(ctx g.Ctx) (conf *webSchema.SuperAdminConfig, err error) {
	conf = &webSchema.SuperAdminConfig{}
	v, _ := g.Cfg().Get(ctx, "system.superAdmin")
	if v != nil {
		err = v.Struct(conf)
	}
	return
}

func (s *sSystemConfig) GetAnonymousConfig(ctx g.Ctx) (conf *webSchema.AnonymousConfig, err error) {
	conf = &webSchema.AnonymousConfig{}
	v, _ := g.Cfg().Get(ctx, "system.anonymous")
	if v != nil {
		err = v.Struct(conf)
	}
	return
}

func (s *sSystemConfig) GetCacheConfig(ctx g.Ctx) (conf *webSchema.CacheConfig, err error) {
	conf = &webSchema.CacheConfig{}
	v, _ := g.Cfg().Get(ctx, "system.cache")
	if v != nil {
		err = v.Struct(conf)
	} else {
		conf = &webSchema.CacheConfig{
			Adapter: "file",
			FileDir: "tmp/cache",
		}
	}
	return
}

func (s *sSystemConfig) GetUploadConfig(ctx g.Ctx) (conf *webSchema.UploadConfig, err error) {
	conf = &webSchema.UploadConfig{}
	v, _ := g.Cfg().Get(ctx, "system.upload")
	if v != nil {
		err = v.Struct(conf)
	} else {
		conf = &webSchema.UploadConfig{
			Drive:     zconsts.UploadDriveLocal,
			FileSize:  50,
			FileType:  "doc,docx,dot,xls,xlsx,xltx,ppt,pptx,pdf,txt,csv,html,xml,pptm,md",
			ImageSize: 10,
			ImageType: "jpg,png,jpeg,gif,webp",
			LocalPath: "upload-files-",
		}
	}
	return
}

// GetMailConfig 读取邮件发送配置（SMTP）。
// 未启用时只做默认值填充，不做强校验，避免影响未配置邮件的项目启动。
func (s *sSystemConfig) GetMailConfig(ctx g.Ctx) (conf *webSchema.MailConfig, err error) {
	conf = &webSchema.MailConfig{}
	v, _ := g.Cfg().Get(ctx, "mail")
	if v != nil {
		if err = v.Struct(conf); err != nil {
			return nil, gerror.Wrap(err, "解析 mail 配置失败")
		}
	}
	if conf.Port <= 0 {
		conf.Port = 465
	}
	if conf.Timeout <= 0 {
		conf.Timeout = 10
	}
	if conf.From == "" {
		conf.From = conf.Username
	}
	if !conf.Enabled {
		return conf, nil
	}
	if conf.Host == "" {
		return nil, gerror.New("mail 配置不完整：启用邮件通知时必须配置 host")
	}
	if conf.Username != "" && conf.Password == "" {
		return nil, gerror.New("mail 配置不完整：配置了 username 时必须同时配置 password（SMTP 授权码）")
	}
	if len(conf.To) == 0 {
		return nil, gerror.New("mail 配置不完整：启用邮件通知时必须配置至少一个收件人 mail.to")
	}
	return conf, nil
}

// GetSmsConfig 读取短信发送配置。
func (s *sSystemConfig) GetSmsConfig(ctx g.Ctx) (conf *webSchema.SmsConfig, err error) {
	conf = &webSchema.SmsConfig{}
	v, _ := g.Cfg().Get(ctx, "sms")
	if v != nil {
		if err = v.Struct(conf); err != nil {
			return nil, gerror.Wrap(err, "解析 sms 配置失败")
		}
	}
	if conf.Timeout <= 0 {
		conf.Timeout = 10
	}
	if conf.TencentRegion == "" {
		conf.TencentRegion = "ap-guangzhou"
	}
	if conf.TencentEndpoint == "" {
		conf.TencentEndpoint = "sms.tencentcloudapi.com"
	}
	if conf.AliyunSignName == "" {
		conf.AliyunSignName = conf.SignName
	}
	if conf.TencentSignName == "" {
		conf.TencentSignName = conf.SignName
	}

	if !conf.Enabled {
		return conf, nil
	}

	switch conf.Drive {
	case zsms.DriveMock:
		if !gmode.IsDevelop() {
			return nil, gerror.New("sms.drive=mock 仅用于本地联调：生产环境不会真正下发短信且验证码会以明文写入日志，请改用 aliyun / tencent")
		}
	case zsms.DriveAliyun:
		if conf.AliyunAccessKeyId == "" || conf.AliyunAccessKeySecret == "" {
			return nil, gerror.New("sms 配置不完整：阿里云短信需要 aliyunAccessKeyId / aliyunAccessKeySecret")
		}
		if conf.AliyunSignName == "" {
			return nil, gerror.New("sms 配置不完整：阿里云短信需要 aliyunSignName 或 signName")
		}
	case zsms.DriveTencent:
		if conf.TencentSecretId == "" || conf.TencentSecretKey == "" {
			return nil, gerror.New("sms 配置不完整：腾讯云短信需要 tencentSecretId / tencentSecretKey")
		}
		if conf.TencentSdkAppId == "" {
			return nil, gerror.New("sms 配置不完整：腾讯云短信需要 tencentSdkAppId")
		}
		if conf.TencentSignName == "" {
			return nil, gerror.New("sms 配置不完整：腾讯云短信需要 tencentSignName 或 signName")
		}
	default:
		return nil, gerror.Newf("sms 配置错误：不支持的驱动 %q", conf.Drive)
	}

	return conf, nil
}

// GetVerifyCodeConfig 读取验证码策略配置。
func (s *sSystemConfig) GetVerifyCodeConfig(ctx g.Ctx) (conf *webSchema.VerifyCodeConfig, err error) {
	conf = &webSchema.VerifyCodeConfig{}
	v, _ := g.Cfg().Get(ctx, "system.verifyCode")
	if v != nil {
		if err = v.Struct(conf); err != nil {
			return nil, gerror.Wrap(err, "解析 system.verifyCode 配置失败")
		}
	}

	if conf.Length <= 0 {
		conf.Length = zconsts.DefaultVerifyCodeLength
	}
	if conf.TTL <= 0 {
		conf.TTL = zconsts.DefaultVerifyCodeTTL
	}
	if conf.Cooldown <= 0 {
		conf.Cooldown = zconsts.DefaultVerifyCodeCooldown
	}
	if conf.MaxAttempts <= 0 {
		conf.MaxAttempts = zconsts.DefaultVerifyCodeMaxAttempts
	}
	if conf.MaxIpDaily == 0 {
		conf.MaxIpDaily = zconsts.DefaultVerifyCodeMaxIpDaily
	}
	if conf.MaxMobileDaily == 0 {
		conf.MaxMobileDaily = zconsts.DefaultVerifyCodeMaxMobDaily
	}
	if conf.Secret == "" {
		tokenSecret, cfgErr := g.Cfg().Get(ctx, "system.token.secretKey")
		if cfgErr != nil {
			return nil, gerror.Wrap(cfgErr, "解析 system.token.secretKey 失败")
		}
		conf.Secret = tokenSecret.String()
	}

	if !conf.Enabled {
		return conf, nil
	}

	if conf.Length < 4 || conf.Length > 10 {
		return nil, gerror.Newf("system.verifyCode.length 必须在 4~10 之间，当前为 %d", conf.Length)
	}
	if conf.Cooldown > conf.TTL {
		return nil, gerror.New("system.verifyCode.cooldown 不能大于 ttl")
	}
	if conf.Secret == "" {
		return nil, gerror.New("system.verifyCode.secret 不能为空（可复用 system.token.secretKey）")
	}
	if conf.MaxIpDaily < -1 {
		return nil, gerror.Newf("system.verifyCode.maxIpDaily 只能为 -1（不限）或正数，当前为 %d", conf.MaxIpDaily)
	}
	if conf.MaxMobileDaily < -1 {
		return nil, gerror.Newf("system.verifyCode.maxMobileDaily 只能为 -1（不限）或正数，当前为 %d", conf.MaxMobileDaily)
	}
	if len(conf.Scenes) == 0 {
		return nil, gerror.New("system.verifyCode 配置不完整：启用验证码服务时必须配置至少一个 scenes")
	}

	seen := make(map[string]struct{}, len(conf.Scenes))
	for _, scene := range conf.Scenes {
		if scene == nil || scene.Scene == "" {
			return nil, gerror.New("system.verifyCode 配置错误：scene 不能为空")
		}
		if scene.Template == "" {
			return nil, gerror.Newf("system.verifyCode 配置错误：场景 %s 未配置 template", scene.Scene)
		}
		if scene.Length > 0 && (scene.Length < 4 || scene.Length > 10) {
			return nil, gerror.Newf("system.verifyCode 配置错误：场景 %s 的 length 必须在 4~10 之间", scene.Scene)
		}
		if scene.MaxIpDaily < -1 {
			return nil, gerror.Newf("system.verifyCode 配置错误：场景 %s 的 maxIpDaily 只能为 -1（不限）或正数", scene.Scene)
		}
		if scene.MaxMobileDaily < -1 {
			return nil, gerror.Newf("system.verifyCode 配置错误：场景 %s 的 maxMobileDaily 只能为 -1（不限）或正数", scene.Scene)
		}
		ttl := scene.TTL
		if ttl <= 0 {
			ttl = conf.TTL
		}
		cooldown := scene.Cooldown
		if cooldown <= 0 {
			cooldown = conf.Cooldown
		}
		if cooldown > ttl {
			return nil, gerror.Newf("system.verifyCode 配置错误：场景 %s 的 cooldown 不能大于 ttl", scene.Scene)
		}
		if _, ok := seen[scene.Scene]; ok {
			return nil, gerror.Newf("system.verifyCode 配置错误：场景 %s 重复", scene.Scene)
		}
		seen[scene.Scene] = struct{}{}
		if scene.Name == "" {
			scene.Name = scene.Scene
		}
	}

	return conf, nil
}

func (s *sSystemConfig) GetTokenConfig(ctx g.Ctx) (conf *webSchema.TokenConfig, err error) {
	conf = &webSchema.TokenConfig{}
	v, _ := g.Cfg().Get(ctx, "system.token")
	if v != nil {
		err = v.Struct(conf)
	} else {
		conf = &webSchema.TokenConfig{
			Expires:         86400,
			RefreshInterval: 3600,
			MaxRefreshTimes: 10,
			SecretKey:       "",
			MultiLogin:      true,
		}
	}
	return
}

// validateTokenConfig 校验 token 配置的密钥强度
func (s *sSystemConfig) validateTokenConfig(ctx g.Ctx, conf *webSchema.TokenConfig) error {
	if conf == nil || conf.SecretKey == "" {
		return gerror.New("system.token.secretKey 未配置，请参考 config.example.yaml 进行配置")
	}
	if len(conf.SecretKey) < 32 {
		if gmode.IsDevelop() {
			g.Log().Warningf(ctx, "system.token.secretKey 长度不足 32 位，生产环境请务必修改")
			return nil
		}
		return gerror.New("system.token.secretKey 长度必须至少 32 位")
	}
	return nil
}

// GetRequestEncryptKey 获取请求加密密钥（强校验）
func (s *sSystemConfig) GetRequestEncryptKey(ctx g.Ctx) (key []byte, err error) {
	v := g.Cfg().MustGet(ctx, "system.requestEncryptKey").String()
	if v == "" {
		return nil, gerror.New("system.requestEncryptKey 未配置，请参考 config.example.yaml 进行配置")
	}
	key = []byte(v)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, gerror.Newf("system.requestEncryptKey 长度必须为 16/24/32 字节，当前为 %d 字节", len(key))
	}
	return key, nil
}

func (s *sSystemConfig) SetGFMode(ctx g.Ctx, mode string) {
	switch mode {
	case "dev", gmode.DEVELOP:
		gmode.Set(gmode.DEVELOP)
	case "test", gmode.TESTING:
		gmode.Set(gmode.TESTING)
	case gmode.STAGING:
		gmode.Set(gmode.STAGING)
	case "prod", "production", gmode.PRODUCT:
		gmode.Set(gmode.PRODUCT)
	}
}
