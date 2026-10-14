package zconsts

const (
	DemoTips                  = "演示系统已隐藏" // 演示系统敏感数据打码
	NilJsonToString           = "{}"          // 空json初始化值
	RegionSpilt               = " / "         // 地区分隔符
	Unknown                   = "Unknown"     // Unknown
	SuperRoleKey              = "super"       // 超管角色唯一标识符，通过角色验证超管
	DefaultSuperAdminUsername = "superAdmin"  // 默认超管用户名，生产环境请通过配置覆盖
	DefaultSuperAdminPassword = "zzfast@admin123" // 默认超管密码，生产环境请通过 system.superAdmin.password 配置覆盖
	DefaultUserPassword       = "123456"      // 默认用户密码，生产环境请覆盖
	MaxServeLogContentLen     = 2048          // 最大保留服务日志内容大小
	EventServerClose          = "server.close" // 服务关闭事件
)
