package common

// SendVerifyCodeInput 发送验证码入参。
//
// 匿名场景（如登录）由 /site/verify-code/send 使用，需传 mobile 并按场景配置校验图形验证码；
// 登录后场景（如高危操作、换绑手机）由 /member/verify-code/send 使用，未传 mobile 时默认使用当前账号手机号。
type SendVerifyCodeInput struct {
	Scene       string `json:"scene"       dc:"业务场景，如 login / high_risk"`
	Mobile      string `json:"mobile"      dc:"手机号"`
	CaptchaId   string `json:"captchaId"   dc:"图形验证码 ID"`
	CaptchaCode string `json:"captchaCode" dc:"图形验证码"`
}

// SendVerifyCodeOutput 发送验证码输出。
type SendVerifyCodeOutput struct {
	ExpiresIn int    `json:"expiresIn" dc:"有效期（秒）"`
	Cooldown  int    `json:"cooldown"  dc:"可重发间隔（秒）"`
	Mobile    string `json:"mobile"    dc:"脱敏后的手机号"`
}

// VerifyVerifyCodeInput 校验验证码入参。
type VerifyVerifyCodeInput struct {
	Scene  string `json:"scene"  dc:"业务场景"`
	Mobile string `json:"mobile" dc:"手机号"`
	Code   string `json:"code"   dc:"验证码"`
}
