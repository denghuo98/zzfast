# 微信支付

ZZFast 基于微信支付官方 APIv3 SDK `github.com/wechatpay-apiv3/wechatpay-go` 封装了常用支付能力。

当前支持：

- Native 支付预下单（PC 扫码）
- JSAPI 支付预下单（公众号/小程序），并生成前端调起支付参数
- 按商户订单号查询支付订单
- 关闭未支付订单
- 发起退款、按商户退款单号查询退款
- 支付/退款回调验签、解密，并转发给业务处理器

> 框架只负责调用微信支付与回调验签；订单创建、金额校验、幂等、订单状态流转等业务逻辑由上层实现。

## 配置

```yaml
pay:
  wechat:
    enabled: true
    appId: "wx0000000000000000"
    mchId: "1900000001"
    mchSerialNo: "商户API证书序列号"
    apiV3Key: "32位APIv3密钥"
    # 二选一：PEM 内容或文件路径
    privateKey: |
      -----BEGIN PRIVATE KEY-----
      ...
      -----END PRIVATE KEY-----
    privateKeyPath: ""
    notifyUrl: "https://example.com/pay/wechat/notify"
    refundNotifyUrl: "https://example.com/pay/wechat/notify"
```

说明：

- 金额统一使用「分」为单位的 `int64`
- `notifyUrl` 必须是外网可访问的 HTTPS 地址
- 客户端懒初始化，首次调用支付接口时会自动下载/更新微信平台证书，需要服务器能访问微信支付接口
- 未启用时调用支付服务会直接返回“微信支付未启用”

## 业务调用

### Native 预下单

```go
out, err := zservice.WechatPay().Native(ctx, &paySchema.WechatPrepayInput{
    OutTradeNo:  "ORDER_20260101_0001",
    Description: "年度会员",
    Total:       9900, // 99.00 元
})
// out.CodeURL 交给前端生成二维码
```

### JSAPI 预下单

```go
out, err := zservice.WechatPay().JSAPI(ctx, &paySchema.WechatPrepayInput{
    OutTradeNo:  "ORDER_20260101_0002",
    Description: "年度会员",
    Total:       9900,
    OpenId:      openId,
})
// out.Params 直接交给前端 WeixinJSBridge / wx.chooseWXPay 调起支付
```

### 查询 / 关单

```go
order, err := zservice.WechatPay().Query(ctx, "ORDER_20260101_0001")
err = zservice.WechatPay().Close(ctx, "ORDER_20260101_0001")
```

### 退款 / 退款查询

```go
refund, err := zservice.WechatPay().Refund(ctx, &paySchema.WechatRefundInput{
    OutTradeNo:  "ORDER_20260101_0001",
    OutRefundNo: "REFUND_20260101_0001",
    Reason:      "用户申请退款",
    Refund:      9900,
    Total:       9900,
})
refund, err = zservice.WechatPay().QueryRefund(ctx, "REFUND_20260101_0001")
```

## 回调处理

框架已内置公开回调接口：

```text
POST /pay/wechat/notify
```

该接口会完成：

1. 校验 `Wechatpay-Signature` 等请求头
2. 使用 APIv3 Key 解密通知内容
3. 转换为 `WechatNotifyInput`
4. 调用业务注册的处理器
5. 按微信要求返回 `{"code":"SUCCESS","message":"成功"}` 或 HTTP 500 + `FAIL`

业务在服务初始化时注册处理器：

```go
func init() {
    zservice.WechatPay().SetNotifyHandler(func(ctx g.Ctx, in *paySchema.WechatNotifyInput) error {
        switch in.EventType {
        case "TRANSACTION.SUCCESS":
            // 根据 in.Payment.OutTradeNo 查询订单、校验金额、幂等更新订单状态
        case "REFUND.SUCCESS":
            // 根据 in.Refund.OutRefundNo 更新退款状态
        }
        return nil
    })
}
```

## 安全建议

- 回调中以 `out_trade_no` 找到本地订单后，必须校验金额与订单状态，不能只相信通知内容
- 支付/退款回调处理器需要实现幂等（微信会重复通知）
- 回调接口不要放入 `AdminAuth`，已在 `/pay` 公开路由组中绑定
- 生产环境私钥、APIv3 Key 通过密钥管理或环境变量注入，不要提交到仓库
