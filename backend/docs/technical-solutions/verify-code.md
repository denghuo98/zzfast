# 验证码

ZZFast 内置短信验证码能力，用于登录二次校验、高危操作确认、换绑手机号等场景。

验证码状态存放在 `zcache`（memory/file/redis），校验成功后立即失效；短信发送通过 `web/zsms` 驱动完成，当前支持 `mock`、阿里云、腾讯云。

## 设计分层

| 层 | 位置 | 职责 |
| --- | --- | --- |
| 状态机 | `web/zverify` | 生成、缓存、冷却、每日限额、错误次数、一次性消费 |
| 短信驱动 | `web/zsms` | 阿里云 / 腾讯云 / mock 发送实现 |
| 业务服务 | `zservice/logic/common/verify_code.go` | 读取场景配置、组装策略、写审计日志 |
| HTTP 接口 | `zcontroller/common`、`zcontroller/admin` | 匿名发码、登录后发码 |

## 配置

```yaml
system:
  verifyCode:
    enabled: true
    secret: ""          # 留空复用 system.token.secretKey
    length: 6
    ttl: 300            # 有效期（秒）
    cooldown: 60        # 同场景同手机号重发间隔（秒）
    maxAttempts: 5      # 最大错误次数
    maxIpDaily: 20      # 同 IP 单场景每日发送上限
    maxMobileDaily: 10  # 同手机号单场景每日发送上限
    # trustedProxies:   # 可信代理 IP/CIDR；默认仅信任本机回环，nginx/SLB 后必须显式配置
    #   - "10.0.0.0/8"
    scenes:
      - scene: "login"
        name: "登录验证"
        template: "SMS_000000000"
        authOnly: false
        requireCaptcha: true
        allowTargetInput: false
      - scene: "high_risk"
        name: "高危操作验证"
        template: "SMS_000000001"
        authOnly: true
        requireCaptcha: false
        allowTargetInput: false

sms:
  enabled: true
  drive: "aliyun"       # mock 仅允许在 develop 模式使用；生产误配会启动失败
  timeout: 10
  signName: "ZZFast"
  aliyunAccessKeyId: ""
  aliyunAccessKeySecret: ""
  # 腾讯云
  tencentSecretId: ""
  tencentSecretKey: ""
  tencentSdkAppId: ""
  tencentRegion: "ap-guangzhou"
```

场景字段说明：

- `authOnly`：`true` 表示仅允许 `/admin/member/verify-code/send`，用于登录后场景
- `requireCaptcha`：发送前是否校验图形验证码，匿名场景建议开启
- `allowTargetInput`：登录后是否允许前端指定手机号（如换绑手机），否则固定使用当前账号手机号
- 场景级 `length/ttl/cooldown/maxAttempts/maxIpDaily/maxMobileDaily` 可覆盖全局值

## HTTP 接口

### 匿名场景（登录）

```http
POST /admin/site/verify-code/send
Content-Type: application/json

{
  "scene": "login",
  "mobile": "13800000000",
  "captchaId": "xxx",
  "captchaCode": "12"
}
```

### 登录后场景（高危操作）

```http
POST /admin/member/verify-code/send
Content-Type: application/json

{
  "scene": "high_risk",
  "mobile": ""
}
```

`allowTargetInput=false` 时忽略 `mobile`，使用当前登录账号绑定的手机号。

> 登录后发码接口位于 `/admin/member/*`，默认受 Casbin 路由权限控制。需要给对应角色分配该路由权限；如果希望所有登录用户都可调用，可将 `/admin/member/verify-code/send` 加入 `system.exceptAuth`（仍要求登录态）。

响应：

```json
{
  "expiresIn": 300,
  "cooldown": 60,
  "mobile": "138****0000"
}
```

## 业务层调用

发送：

```go
_, err := zservice.VerifyCode().Send(ctx, &commonSchema.SendVerifyCodeInput{
    Scene:  zconsts.VerifySceneHighRisk,
    Mobile: mobile,
})
```

校验（成功后立即失效）：

```go
if err = zservice.VerifyCode().Verify(ctx, &commonSchema.VerifyVerifyCodeInput{
    Scene:  zconsts.VerifySceneHighRisk,
    Mobile: mobile,
    Code:   code,
}); err != nil {
    return err
}
```

建议在高危操作中先完成业务参数校验，再调用 `Verify`，避免验证码被无效请求消耗。

## 安全说明

- 验证码使用 HMAC-SHA256（密钥复用 `system.token.secretKey` 或 `verifyCode.secret`）存储，不存明文
- 图形验证码在服务层统一校验，`/admin/site/*` 和 `/admin/member/*` 入口共用同一套场景策略
- 部署在 nginx/SLB 后时必须配置 `trustedProxies`，否则无法取得真实客户端 IP，IP 日限额会跳过
- 校验成功立即删除，天然一次性
- 发送与校验均有进程/多实例共享的短锁，避免并发绕过
- 默认限制同手机号冷却、每日手机号/ IP 上限、错误次数
- 审计日志写入 `zz_sys_verify_code_log`，只记录发送/校验结果，不记录验证码
