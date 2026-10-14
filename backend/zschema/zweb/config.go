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
