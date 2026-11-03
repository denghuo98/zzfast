package zconsts

// 验证码场景（具体可用场景由配置 system.verifyCode.scenes 决定）
const (
	VerifySceneLogin         = "login"
	VerifySceneRegister      = "register"
	VerifySceneResetPassword = "reset_password"
	VerifySceneBindMobile    = "bind_mobile"
	VerifySceneHighRisk      = "high_risk"
)

// 验证码发送通道
const (
	VerifyChannelSms = "sms"
)

// 验证码审计日志状态
const (
	VerifyLogStatusSent         = 1 // 发送成功
	VerifyLogStatusSendFailed   = 2 // 发送失败
	VerifyLogStatusVerified     = 3 // 校验成功
	VerifyLogStatusVerifyFailed = 4 // 校验失败
)

// 验证码默认策略
const (
	DefaultVerifyCodeLength      = 6
	DefaultVerifyCodeTTL         = 300 // 秒
	DefaultVerifyCodeCooldown    = 60  // 秒
	DefaultVerifyCodeMaxAttempts = 5
	DefaultVerifyCodeMaxIpDaily  = 20
	DefaultVerifyCodeMaxMobDaily = 10
)
