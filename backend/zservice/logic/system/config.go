package system

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gmode"
	"github.com/samber/lo"

	"github.com/denghuo98/zzfast/web/zcache"
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
	if len(mode) == 0 {
		mode = gmode.NOT_SET
	}

	var modes = []string{gmode.DEVELOP, gmode.TESTING, gmode.STAGING, gmode.PRODUCT}

	if lo.Contains(modes, mode) {
		gmode.Set(mode)
	}
}
