package pay

import "github.com/gogf/gf/v2/frame/g"

// WechatNotifyReq 微信支付/退款回调
type WechatNotifyReq struct {
	g.Meta `path:"/wechat/notify" method:"post" tags:"PAY-01-支付" summary:"微信支付回调"`
}

type WechatNotifyRes struct {
}
