package zservice

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	paySchema "github.com/denghuo98/zzfast/zschema/pay"
)

// WechatPayNotifyHandler 微信支付/退款回调业务处理器。
type WechatPayNotifyHandler func(ctx g.Ctx, in *paySchema.WechatNotifyInput) error

// IWechatPay 微信支付能力。
type IWechatPay interface {
	Native(ctx g.Ctx, in *paySchema.WechatPrepayInput) (*paySchema.WechatNativeOutput, error)
	JSAPI(ctx g.Ctx, in *paySchema.WechatPrepayInput) (*paySchema.WechatJSAPIOutput, error)
	Query(ctx g.Ctx, outTradeNo string) (*paySchema.WechatQueryOutput, error)
	Close(ctx g.Ctx, outTradeNo string) error
	Refund(ctx g.Ctx, in *paySchema.WechatRefundInput) (*paySchema.WechatRefundOutput, error)
	QueryRefund(ctx g.Ctx, outRefundNo string) (*paySchema.WechatRefundOutput, error)

	ParseNotify(ctx g.Ctx, r *ghttp.Request) (*paySchema.WechatNotifyInput, error)
	HandleNotify(ctx g.Ctx, in *paySchema.WechatNotifyInput) error
	SetNotifyHandler(handler WechatPayNotifyHandler)
}

var localWechatPay IWechatPay

func WechatPay() IWechatPay {
	if localWechatPay == nil {
		panic("WechatPay is not initialized, please register it first")
	}
	return localWechatPay
}

func RegisterWechatPay(i IWechatPay) {
	localWechatPay = i
}
