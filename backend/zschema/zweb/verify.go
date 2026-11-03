package zweb

// SmsConfig 短信发送配置。
//
// 当前支持 mock / aliyun / tencent 三种驱动：
//   - mock: 不真正发送，仅记录日志，便于本地联调
//   - aliyun: 阿里云短信（dysmsapi）
//   - tencent: 腾讯云短信（sms v20210111）
type SmsConfig struct {
	Enabled  bool   `json:"enabled"  dc:"是否启用短信服务"`
	Drive    string `json:"drive"    dc:"短信驱动：mock | aliyun | tencent"`
	Timeout  int    `json:"timeout"  dc:"发送超时时间（秒），默认 10"`
	SignName string `json:"signName" dc:"短信签名，阿里云/腾讯云可共用"`

	// 阿里云
	AliyunAccessKeyId     string `json:"aliyunAccessKeyId"     dc:"阿里云 AccessKeyId"`
	AliyunAccessKeySecret string `json:"aliyunAccessKeySecret" dc:"阿里云 AccessKeySecret"`
	AliyunSignName        string `json:"aliyunSignName"        dc:"阿里云短信签名，留空时使用 signName"`

	// 腾讯云
	TencentSecretId  string `json:"tencentSecretId"  dc:"腾讯云 SecretId"`
	TencentSecretKey string `json:"tencentSecretKey" dc:"腾讯云 SecretKey"`
	TencentSdkAppId  string `json:"tencentSdkAppId"  dc:"腾讯云短信 SdkAppId"`
	TencentSignName  string `json:"tencentSignName"  dc:"腾讯云短信签名，留空时使用 signName"`
	TencentRegion    string `json:"tencentRegion"    dc:"腾讯云地域，默认 ap-guangzhou"`
	TencentEndpoint  string `json:"tencentEndpoint"  dc:"腾讯云短信接口域名，默认 sms.tencentcloudapi.com"`
}

// VerifyCodeConfig 验证码全局策略。
type VerifyCodeConfig struct {
	Enabled        bool                     `json:"enabled"         dc:"是否启用验证码服务"`
	Secret         string                   `json:"secret"          dc:"验证码哈希密钥，留空时使用 system.token.secretKey"`
	Length         int                      `json:"length"          dc:"验证码长度，默认 6"`
	TTL            int                      `json:"ttl"             dc:"验证码有效期（秒），默认 300"`
	Cooldown       int                      `json:"cooldown"        dc:"同场景同手机号重发间隔（秒），默认 60"`
	MaxAttempts    int                      `json:"maxAttempts"     dc:"最大错误次数，默认 5"`
	MaxIpDaily     int                      `json:"maxIpDaily"      dc:"同 IP 单场景每日发送上限，默认 20；-1 表示不限"`
	MaxMobileDaily int                      `json:"maxMobileDaily"  dc:"同手机号单场景每日发送上限，默认 10；-1 表示不限"`
	TrustedProxies []string                 `json:"trustedProxies"   dc:"可信代理 IP/CIDR；默认仅信任本机回环，代理部署需显式配置，否则跳过 IP 限流"`
	Scenes         []*VerifyCodeSceneConfig `json:"scenes"          dc:"业务场景列表"`
}

// VerifyCodeSceneConfig 单个验证码场景配置。
type VerifyCodeSceneConfig struct {
	Scene            string `json:"scene"         dc:"场景标识，如 login / high_risk"`
	Name             string `json:"name"          dc:"场景名称"`
	Template         string `json:"template"      dc:"短信模板 ID（阿里云 TemplateCode / 腾讯云 TemplateId）"`
	AuthOnly         bool   `json:"authOnly"      dc:"仅允许登录后调用"`
	RequireCaptcha   bool   `json:"requireCaptcha" dc:"发送前是否校验图形验证码"`
	AllowTargetInput bool   `json:"allowTargetInput" dc:"登录后是否允许指定手机号（如换绑手机）"`

	// 以下字段为 0 时继承全局配置；每日上限设为 -1 表示不限
	Length         int `json:"length"         dc:"验证码长度"`
	TTL            int `json:"ttl"            dc:"验证码有效期（秒）"`
	Cooldown       int `json:"cooldown"       dc:"重发间隔（秒）"`
	MaxAttempts    int `json:"maxAttempts"    dc:"最大错误次数"`
	MaxIpDaily     int `json:"maxIpDaily"     dc:"同 IP 每日发送上限"`
	MaxMobileDaily int `json:"maxMobileDaily" dc:"同手机号每日发送上限"`
}

// GetScene 按场景标识查找配置，未找到返回 nil。
func (c *VerifyCodeConfig) GetScene(scene string) *VerifyCodeSceneConfig {
	if c == nil {
		return nil
	}
	for _, item := range c.Scenes {
		if item != nil && item.Scene == scene {
			return item
		}
	}
	return nil
}
