// Package zwechatpay 封装微信支付 APIv3 的常用能力。
//
// 当前提供 Native / JSAPI 预下单、订单查询、关单、退款、退款查询与回调验签解密。
// 业务订单、金额校验、幂等与回调业务处理仍由上层业务负责。
package zwechatpay

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	jsapi "github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	native "github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	refund "github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"

	paySchema "github.com/denghuo98/zzfast/zschema/pay"
	webSchema "github.com/denghuo98/zzfast/zschema/zweb"
)

var (
	config    *webSchema.WechatPayConfig
	configMu  sync.RWMutex
	client    *core.Client
	clientMu  sync.Mutex
	clientErr error
)

// SetConfig 由系统配置加载流程注入。
func SetConfig(c *webSchema.WechatPayConfig) {
	configMu.Lock()
	config = c
	configMu.Unlock()

	clientMu.Lock()
	client = nil
	clientErr = nil
	clientMu.Unlock()
}

// GetConfig 获取当前微信支付配置（可能为 nil）。
func GetConfig() *webSchema.WechatPayConfig {
	configMu.RLock()
	defer configMu.RUnlock()
	return config
}

// Enabled 微信支付是否已启用。
func Enabled() bool {
	cfg := GetConfig()
	return cfg != nil && cfg.Enabled
}

// Validate 校验微信支付配置是否完整可用。
func Validate(cfg *webSchema.WechatPayConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.AppId == "" || cfg.MchId == "" || cfg.MchSerialNo == "" {
		return gerror.New("微信支付配置不完整：appId / mchId / mchSerialNo 不能为空")
	}
	if len(cfg.ApiV3Key) != 32 {
		return gerror.Newf("微信支付 apiV3Key 必须为 32 位，当前为 %d 位", len(cfg.ApiV3Key))
	}
	if cfg.NotifyUrl == "" {
		return gerror.New("微信支付配置不完整：notifyUrl 不能为空")
	}
	if _, err := loadPrivateKey(cfg); err != nil {
		return gerror.Wrap(err, "微信支付商户私钥无效")
	}
	return nil
}

// Client 获取微信支付客户端（懒初始化，首次调用会下载平台证书）。
func Client(ctx context.Context) (*core.Client, error) {
	if !Enabled() {
		return nil, gerror.New("微信支付未启用")
	}
	cfg := GetConfig()
	clientMu.Lock()
	defer clientMu.Unlock()
	if client != nil {
		return client, nil
	}
	key, err := loadPrivateKey(cfg)
	if err != nil {
		clientErr = gerror.Wrap(err, "加载微信支付商户私钥失败")
		return nil, clientErr
	}
	client, clientErr = core.NewClient(ctx, option.WithWechatPayAutoAuthCipher(
		cfg.MchId, cfg.MchSerialNo, key, cfg.ApiV3Key,
	))
	return client, clientErr
}

// Native 创建 Native 支付预下单，返回二维码链接。
func Native(ctx context.Context, in *paySchema.WechatPrepayInput) (*paySchema.WechatNativeOutput, error) {
	if err := validatePrepay(in); err != nil {
		return nil, err
	}
	c, err := Client(ctx)
	if err != nil {
		return nil, err
	}
	cfg := GetConfig()
	svc := native.NativeApiService{Client: c}
	req := native.PrepayRequest{
		Appid:       ptr(cfg.AppId),
		Mchid:       ptr(cfg.MchId),
		Description: ptr(in.Description),
		OutTradeNo:  ptr(in.OutTradeNo),
		NotifyUrl:   ptr(cfg.NotifyUrl),
		Amount:      &native.Amount{Total: ptr(in.Total), Currency: ptr("CNY")},
	}
	if in.Attach != "" {
		req.Attach = ptr(in.Attach)
	}
	if in.ExpireAt != nil {
		req.TimeExpire = in.ExpireAt
	}

	resp, _, err := svc.Prepay(ctx, req)
	if err != nil {
		return nil, gerror.Wrap(err, "微信 Native 预下单失败")
	}
	if resp == nil || resp.CodeUrl == nil {
		return nil, gerror.New("微信 Native 预下单返回为空")
	}
	return &paySchema.WechatNativeOutput{CodeURL: *resp.CodeUrl}, nil
}

// JSAPI 创建 JSAPI 支付预下单，并返回前端调起支付所需参数。
func JSAPI(ctx context.Context, in *paySchema.WechatPrepayInput) (*paySchema.WechatJSAPIOutput, error) {
	if err := validatePrepay(in); err != nil {
		return nil, err
	}
	if in.OpenId == "" {
		return nil, gerror.New("JSAPI 支付必须提供 openId")
	}
	c, err := Client(ctx)
	if err != nil {
		return nil, err
	}
	cfg := GetConfig()
	svc := jsapi.JsapiApiService{Client: c}
	req := jsapi.PrepayRequest{
		Appid:       ptr(cfg.AppId),
		Mchid:       ptr(cfg.MchId),
		Description: ptr(in.Description),
		OutTradeNo:  ptr(in.OutTradeNo),
		NotifyUrl:   ptr(cfg.NotifyUrl),
		Amount:      &jsapi.Amount{Total: ptr(in.Total), Currency: ptr("CNY")},
		Payer:       &jsapi.Payer{Openid: ptr(in.OpenId)},
	}
	if in.Attach != "" {
		req.Attach = ptr(in.Attach)
	}
	if in.ExpireAt != nil {
		req.TimeExpire = in.ExpireAt
	}

	resp, _, err := svc.Prepay(ctx, req)
	if err != nil {
		return nil, gerror.Wrap(err, "微信 JSAPI 预下单失败")
	}
	if resp == nil || resp.PrepayId == nil {
		return nil, gerror.New("微信 JSAPI 预下单返回为空")
	}
	params, err := BuildJSAPIParams(ctx, *resp.PrepayId)
	if err != nil {
		return nil, err
	}
	return &paySchema.WechatJSAPIOutput{PrepayId: *resp.PrepayId, Params: params}, nil
}

// BuildJSAPIParams 生成微信 JSAPI 调起支付参数。
func BuildJSAPIParams(ctx context.Context, prepayId string) (map[string]string, error) {
	if prepayId == "" {
		return nil, gerror.New("prepayId 不能为空")
	}
	c, err := Client(ctx)
	if err != nil {
		return nil, err
	}
	cfg := GetConfig()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := randomString(16)
	pkg := "prepay_id=" + prepayId
	message := cfg.AppId + "\n" + timestamp + "\n" + nonce + "\n" + pkg + "\n"
	signature, err := c.Sign(ctx, message)
	if err != nil {
		return nil, gerror.Wrap(err, "微信 JSAPI 支付参数签名失败")
	}
	return map[string]string{
		"appId":     cfg.AppId,
		"timeStamp": timestamp,
		"nonceStr":  nonce,
		"package":   pkg,
		"signType":  "RSA",
		"paySign":   signature.Signature,
	}, nil
}

// Query 按商户订单号查询支付订单。
func Query(ctx context.Context, outTradeNo string) (*paySchema.WechatQueryOutput, error) {
	if outTradeNo == "" {
		return nil, gerror.New("商户订单号不能为空")
	}
	c, err := Client(ctx)
	if err != nil {
		return nil, err
	}
	cfg := GetConfig()
	svc := native.NativeApiService{Client: c}
	resp, _, err := svc.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: ptr(outTradeNo),
		Mchid:      ptr(cfg.MchId),
	})
	if err != nil {
		return nil, gerror.Wrap(err, "微信支付订单查询失败")
	}
	return mapTransaction(resp), nil
}

// Close 关闭未支付订单。
func Close(ctx context.Context, outTradeNo string) error {
	if outTradeNo == "" {
		return gerror.New("商户订单号不能为空")
	}
	c, err := Client(ctx)
	if err != nil {
		return err
	}
	cfg := GetConfig()
	svc := native.NativeApiService{Client: c}
	if _, err = svc.CloseOrder(ctx, native.CloseOrderRequest{
		OutTradeNo: ptr(outTradeNo),
		Mchid:      ptr(cfg.MchId),
	}); err != nil {
		return gerror.Wrap(err, "微信支付关单失败")
	}
	return nil
}

// Refund 发起退款。
func Refund(ctx context.Context, in *paySchema.WechatRefundInput) (*paySchema.WechatRefundOutput, error) {
	if in == nil || in.OutTradeNo == "" || in.OutRefundNo == "" {
		return nil, gerror.New("商户订单号和退款单号不能为空")
	}
	if in.Refund <= 0 || in.Total <= 0 {
		return nil, gerror.New("退款金额和原订单金额必须大于 0")
	}
	c, err := Client(ctx)
	if err != nil {
		return nil, err
	}
	cfg := GetConfig()
	svc := refund.RefundsApiService{Client: c}
	req := refund.CreateRequest{
		OutTradeNo:  ptr(in.OutTradeNo),
		OutRefundNo: ptr(in.OutRefundNo),
		Amount: &refund.AmountReq{
			Refund:   ptr(in.Refund),
			Total:    ptr(in.Total),
			Currency: ptr("CNY"),
		},
	}
	if in.Reason != "" {
		req.Reason = ptr(in.Reason)
	}
	if cfg.RefundNotifyUrl != "" {
		req.NotifyUrl = ptr(cfg.RefundNotifyUrl)
	}

	resp, _, err := svc.Create(ctx, req)
	if err != nil {
		return nil, gerror.Wrap(err, "微信退款申请失败")
	}
	return mapRefund(resp), nil
}

// QueryRefund 按商户退款单号查询退款。
func QueryRefund(ctx context.Context, outRefundNo string) (*paySchema.WechatRefundOutput, error) {
	if outRefundNo == "" {
		return nil, gerror.New("商户退款单号不能为空")
	}
	c, err := Client(ctx)
	if err != nil {
		return nil, err
	}
	svc := refund.RefundsApiService{Client: c}
	resp, _, err := svc.QueryByOutRefundNo(ctx, refund.QueryByOutRefundNoRequest{
		OutRefundNo: ptr(outRefundNo),
	})
	if err != nil {
		return nil, gerror.Wrap(err, "微信退款查询失败")
	}
	return mapRefund(resp), nil
}

// ParseNotify 校验并解密微信支付/退款回调。
func ParseNotify(ctx context.Context, r *http.Request) (*paySchema.WechatNotifyInput, error) {
	if r == nil {
		return nil, gerror.New("回调请求不能为空")
	}
	if _, err := Client(ctx); err != nil {
		return nil, err
	}
	cfg := GetConfig()
	verifier := verifiers.NewSHA256WithRSAVerifier(downloader.MgrInstance().GetCertificateVisitor(cfg.MchId))
	handler, err := notify.NewRSANotifyHandler(cfg.ApiV3Key, verifier)
	if err != nil {
		return nil, gerror.Wrap(err, "创建微信支付回调处理器失败")
	}

	content := notify.ContentMap{}
	req, err := handler.ParseNotifyRequest(ctx, r, &content)
	if err != nil {
		return nil, gerror.Wrap(err, "微信支付回调验签失败")
	}

	input := &paySchema.WechatNotifyInput{
		Id:           req.ID,
		EventType:    req.EventType,
		ResourceType: req.ResourceType,
		Summary:      req.Summary,
	}
	raw, _ := json.Marshal(content)
	switch {
	case strings.HasPrefix(req.EventType, "TRANSACTION"):
		var payment notifyPayment
		if err = json.Unmarshal(raw, &payment); err != nil {
			return nil, gerror.Wrap(err, "解析微信支付通知内容失败")
		}
		input.Payment = &paySchema.WechatNotifyPayment{
			OutTradeNo:     payment.OutTradeNo,
			TransactionId:  payment.TransactionId,
			TradeState:     payment.TradeState,
			TradeStateDesc: payment.TradeStateDesc,
			TradeType:      payment.TradeType,
			SuccessTime:    payment.SuccessTime,
			Total:          payment.Amount.Total,
			PayerTotal:     payment.Amount.PayerTotal,
			OpenId:         payment.Payer.Openid,
			Attach:         payment.Attach,
		}
	case strings.HasPrefix(req.EventType, "REFUND"):
		var refundInfo notifyRefund
		if err = json.Unmarshal(raw, &refundInfo); err != nil {
			return nil, gerror.Wrap(err, "解析微信退款通知内容失败")
		}
		input.Refund = &paySchema.WechatNotifyRefund{
			OutRefundNo:   refundInfo.OutRefundNo,
			RefundId:      refundInfo.RefundId,
			OutTradeNo:    refundInfo.OutTradeNo,
			TransactionId: refundInfo.TransactionId,
			RefundStatus:  refundInfo.RefundStatus,
			SuccessTime:   refundInfo.SuccessTime,
			Total:         refundInfo.Amount.Total,
			Refund:        refundInfo.Amount.Refund,
			PayerTotal:    refundInfo.Amount.PayerTotal,
			PayerRefund:   refundInfo.Amount.PayerRefund,
		}
	}
	return input, nil
}

type notifyPayment struct {
	OutTradeNo     string `json:"out_trade_no"`
	TransactionId  string `json:"transaction_id"`
	TradeState     string `json:"trade_state"`
	TradeStateDesc string `json:"trade_state_desc"`
	TradeType      string `json:"trade_type"`
	SuccessTime    string `json:"success_time"`
	Attach         string `json:"attach"`
	Amount         struct {
		Total      int64 `json:"total"`
		PayerTotal int64 `json:"payer_total"`
	} `json:"amount"`
	Payer struct {
		Openid string `json:"openid"`
	} `json:"payer"`
}

type notifyRefund struct {
	OutRefundNo   string `json:"out_refund_no"`
	RefundId      string `json:"refund_id"`
	OutTradeNo    string `json:"out_trade_no"`
	TransactionId string `json:"transaction_id"`
	RefundStatus  string `json:"refund_status"`
	SuccessTime   string `json:"success_time"`
	Amount        struct {
		Total       int64 `json:"total"`
		Refund      int64 `json:"refund"`
		PayerTotal  int64 `json:"payer_total"`
		PayerRefund int64 `json:"payer_refund"`
	} `json:"amount"`
}

func validatePrepay(in *paySchema.WechatPrepayInput) error {
	if in == nil {
		return gerror.New("预下单参数不能为空")
	}
	if in.OutTradeNo == "" {
		return gerror.New("商户订单号不能为空")
	}
	if in.Description == "" {
		return gerror.New("商品描述不能为空")
	}
	if in.Total <= 0 {
		return gerror.New("订单金额必须大于 0")
	}
	return nil
}

func mapTransaction(resp *payments.Transaction) *paySchema.WechatQueryOutput {
	out := &paySchema.WechatQueryOutput{}
	if resp == nil {
		return out
	}
	out.OutTradeNo = value(resp.OutTradeNo)
	out.TransactionId = value(resp.TransactionId)
	out.TradeState = value(resp.TradeState)
	out.TradeStateDesc = value(resp.TradeStateDesc)
	out.TradeType = value(resp.TradeType)
	out.SuccessTime = value(resp.SuccessTime)
	if resp.Amount != nil {
		out.Total = valueInt64(resp.Amount.Total)
		out.PayerTotal = valueInt64(resp.Amount.PayerTotal)
	}
	if resp.Payer != nil {
		out.OpenId = value(resp.Payer.Openid)
	}
	return out
}

func mapRefund(resp *refund.Refund) *paySchema.WechatRefundOutput {
	out := &paySchema.WechatRefundOutput{}
	if resp == nil {
		return out
	}
	out.RefundId = value(resp.RefundId)
	out.OutRefundNo = value(resp.OutRefundNo)
	if resp.Status != nil {
		out.Status = string(*resp.Status)
	}
	if resp.Amount != nil {
		out.Total = valueInt64(resp.Amount.Total)
		out.Refund = valueInt64(resp.Amount.Refund)
	}
	return out
}

func loadPrivateKey(cfg *webSchema.WechatPayConfig) (*rsa.PrivateKey, error) {
	if cfg.PrivateKey != "" {
		return utils.LoadPrivateKey(cfg.PrivateKey)
	}
	if cfg.PrivateKeyPath != "" {
		return utils.LoadPrivateKeyWithPath(cfg.PrivateKeyPath)
	}
	return nil, fmt.Errorf("未配置 wechat.privateKey 或 wechat.privateKeyPath")
}

func ptr[T any](v T) *T {
	return &v
}

func value(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func valueInt64(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

func randomString(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(buf)
}
