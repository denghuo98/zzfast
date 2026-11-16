package pay

import (
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/denghuo98/zzfast/web/zwechatpay"
	paySchema "github.com/denghuo98/zzfast/zschema/pay"
	"github.com/denghuo98/zzfast/zservice"
)

type sWechatPay struct {
	mu      sync.RWMutex
	handler zservice.WechatPayNotifyHandler
}

func init() {
	zservice.RegisterWechatPay(NewWechatPay())
}

func NewWechatPay() *sWechatPay {
	return &sWechatPay{}
}

func (s *sWechatPay) Native(ctx g.Ctx, in *paySchema.WechatPrepayInput) (*paySchema.WechatNativeOutput, error) {
	return zwechatpay.Native(ctx, in)
}

func (s *sWechatPay) JSAPI(ctx g.Ctx, in *paySchema.WechatPrepayInput) (*paySchema.WechatJSAPIOutput, error) {
	return zwechatpay.JSAPI(ctx, in)
}

func (s *sWechatPay) Query(ctx g.Ctx, outTradeNo string) (*paySchema.WechatQueryOutput, error) {
	return zwechatpay.Query(ctx, outTradeNo)
}

func (s *sWechatPay) Close(ctx g.Ctx, outTradeNo string) error {
	return zwechatpay.Close(ctx, outTradeNo)
}

func (s *sWechatPay) Refund(ctx g.Ctx, in *paySchema.WechatRefundInput) (*paySchema.WechatRefundOutput, error) {
	return zwechatpay.Refund(ctx, in)
}

func (s *sWechatPay) QueryRefund(ctx g.Ctx, outRefundNo string) (*paySchema.WechatRefundOutput, error) {
	return zwechatpay.QueryRefund(ctx, outRefundNo)
}

func (s *sWechatPay) ParseNotify(ctx g.Ctx, r *ghttp.Request) (*paySchema.WechatNotifyInput, error) {
	if r == nil {
		return nil, gerror.New("回调请求不能为空")
	}
	return zwechatpay.ParseNotify(ctx, r.Request)
}

func (s *sWechatPay) HandleNotify(ctx g.Ctx, in *paySchema.WechatNotifyInput) error {
	s.mu.RLock()
	handler := s.handler
	s.mu.RUnlock()
	if handler == nil {
		g.Log().Errorf(ctx, "收到微信支付回调但未注册业务处理器 event=%s", in.EventType)
		return gerror.New("未注册微信支付回调处理器")
	}
	return handler(ctx, in)
}

func (s *sWechatPay) SetNotifyHandler(handler zservice.WechatPayNotifyHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handler = handler
}
