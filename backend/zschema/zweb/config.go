package zweb

type SystemConfig struct {
	Mode string `json:"mode" dc:"系统模式"` // dev | test | prod
}

type SuperAdminConfig struct {
	Password string `json:"password"`
}

type AnonymousConfig struct {
	Enabled  bool     `json:"enabled" dc:"是否启用匿名身份"`
	Identity Identity `json:"identity" dc:"匿名用户信息"`
}

type CacheConfig struct {
	Adapter string `json:"adapter" dc:"缓存适配器"`  // redis | file
	FileDir string `json:"fileDir" dc:"文件缓存目录"` // 文件缓存目录，当adapter为file时必填
}

// TokenConfig 登录令牌配置
type TokenConfig struct {
	SecretKey       string `json:"secretKey"`
	Expires         int64  `json:"expires"`
	AutoRefresh     bool   `json:"autoRefresh"`
	RefreshInterval int64  `json:"refreshInterval"`
	MaxRefreshTimes int64  `json:"maxRefreshTimes"`
	MultiLogin      bool   `json:"multiLogin"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	// 通用配置
	Drive string `json:"drive"`
	// 最大文件体积，单位 Mb
	FileSize int64  `json:"fileSize"`
	FileType string `json:"fileType"`
	// 最大图片体积，单位 Mb
	ImageSize int64  `json:"imageSize"`
	ImageType string `json:"imageType"`
	// 本地存储配置
	LocalPath string `json:"localPath"`
	// 腾讯云cos配置
	CosSecretId  string `json:"cosSecretId"`
	CosSecretKey string `json:"cosSecretKey"`
	CosBucketURL string `json:"cosBucketURL"`
	CosPath      string `json:"cosPath"`
	// 阿里云oss配置
	OssSecretId  string `json:"ossSecretId"`
	OssSecretKey string `json:"ossSecretKey"`
	OssEndpoint  string `json:"ossEndpoint"`
	OssBucketURL string `json:"ossBucketURL"`
	OssPath      string `json:"ossPath"`
	OssBucket    string `json:"ossBucket"`
	// minio配置
	MinioAccessKey string `json:"minioAccessKey"`
	MinioSecretKey string `json:"minioSecretKey"`
	MinioEndpoint  string `json:"minioEndpoint"`
	MinioUseSSL    int    `json:"minioUseSSL"`
	MinioPath      string `json:"minioPath"`
	MinioBucket    string `json:"minioBucket"`
	MinioDomain    string `json:"minioDomain"`
}

// MailConfig 邮件发送配置（SMTP）
// 说明：Password 为邮箱的 SMTP 授权码/客户端专用密码，不是登录密码。
type MailConfig struct {
	Enabled  bool     `json:"enabled"  dc:"是否启用邮件通知"`
	Host     string   `json:"host"     dc:"SMTP 服务器地址，如 smtp.163.com"`
	Port     int      `json:"port"     dc:"SMTP 端口，SSL 通常为 465，STARTTLS 为 25/587"`
	Username string   `json:"username" dc:"SMTP 登录账号（通常是完整邮箱地址）"`
	Password string   `json:"password" dc:"SMTP 授权码/客户端专用密码"`
	From     string   `json:"from"     dc:"发件人地址，默认与 username 相同"`
	FromName string   `json:"fromName" dc:"发件人显示名称"`
	SSL      bool     `json:"ssl"      dc:"是否使用隐式 SSL（465 端口开启，25/587 关闭走 STARTTLS）"`
	To       []string `json:"to"       dc:"默认收件人列表"`
	Timeout  int      `json:"timeout"  dc:"发送超时时间（秒），默认 10"`
}
