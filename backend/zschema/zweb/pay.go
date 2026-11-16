package zweb

// PayConfig 支付配置。
type PayConfig struct {
	Wechat *WechatPayConfig `json:"wechat"`
}

// WechatPayConfig 微信支付 APIv3 配置。
//
// 证书私钥支持直接填写 PEM 内容（privateKey），或填写文件路径（privateKeyPath）。
// APIv3 Key 必须为 32 字节，用于回调通知解密。
type WechatPayConfig struct {
	Enabled         bool   `json:"enabled"          dc:"是否启用微信支付"`
	AppId           string `json:"appId"            dc:"公众号/小程序 AppId"`
	MchId           string `json:"mchId"            dc:"商户号"`
	MchSerialNo     string `json:"mchSerialNo"      dc:"商户 API 证书序列号"`
	ApiV3Key        string `json:"apiV3Key"         dc:"APIv3 密钥（32 字节）"`
	PrivateKey      string `json:"privateKey"       dc:"商户私钥 PEM 内容"`
	PrivateKeyPath  string `json:"privateKeyPath"   dc:"商户私钥文件路径，privateKey 为空时使用"`
	NotifyUrl       string `json:"notifyUrl"        dc:"支付结果回调地址，必须为外网可访问的 HTTPS URL"`
	RefundNotifyUrl string `json:"refundNotifyUrl"  dc:"退款结果回调地址，留空时不单独指定"`
}
