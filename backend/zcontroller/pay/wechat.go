package pay

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	payApi "github.com/denghuo98/zzfast/zapi/pay"
	"github.com/denghuo98/zzfast/zservice"
)

type cWechat struct{}

var WechatNotify = cWechat{}

// Notify 接收微信支付/退款回调，按微信要求返回 SUCCESS/FAIL。
func (c *cWechat) Notify(ctx g.Ctx, req *payApi.WechatNotifyReq) (res *payApi.WechatNotifyRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return nil, nil
	}

	in, parseErr := zservice.WechatPay().ParseNotify(ctx, r)
	if parseErr != nil {
		writeNotifyResult(r, false, parseErr)
		return nil, nil
	}
	if handleErr := zservice.WechatPay().HandleNotify(ctx, in); handleErr != nil {
		writeNotifyResult(r, false, handleErr)
		return nil, nil
	}
	writeNotifyResult(r, true, nil)
	return nil, nil
}

func writeNotifyResult(r *ghttp.Request, success bool, err error) {
	r.Response.Header().Set("Content-Type", "application/json")
	if success {
		r.Response.WriteJson(g.Map{"code": "SUCCESS", "message": "成功"})
		return
	}
	g.Log().Errorf(r.Context(), "微信支付回调处理失败: %v", err)
	r.Response.Status = 500
	r.Response.WriteJson(g.Map{"code": "FAIL", "message": "处理失败"})
}
