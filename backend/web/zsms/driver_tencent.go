package zsms

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	webSchema "github.com/denghuo98/zzfast/zschema/zweb"
)

type tencentDriver struct {
	cfg *webSchema.SmsConfig
}

// Send 通过腾讯云短信服务发送验证码。
func (d *tencentDriver) Send(ctx context.Context, msg Message) error {
	if d.cfg.TencentSecretId == "" || d.cfg.TencentSecretKey == "" {
		return fmt.Errorf("zsms: 腾讯云短信缺少 SecretId/SecretKey")
	}
	if d.cfg.TencentSdkAppId == "" {
		return fmt.Errorf("zsms: 腾讯云短信缺少 SdkAppId")
	}

	credential := common.NewCredential(d.cfg.TencentSecretId, d.cfg.TencentSecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = timeout(d.cfg)
	if d.cfg.TencentEndpoint != "" {
		cpf.HttpProfile.Endpoint = d.cfg.TencentEndpoint
	}
	region := d.cfg.TencentRegion
	if region == "" {
		region = "ap-guangzhou"
	}

	client, err := sms.NewClient(credential, region, cpf)
	if err != nil {
		return fmt.Errorf("zsms: 创建腾讯云短信客户端失败: %w", err)
	}

	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = common.StringPtr(d.cfg.TencentSdkAppId)
	req.SignName = common.StringPtr(signName(d.cfg, DriveTencent))
	req.TemplateId = common.StringPtr(msg.Template)
	req.TemplateParamSet = common.StringPtrs([]string{msg.Code})
	req.PhoneNumberSet = common.StringPtrs([]string{e164(msg.Mobile)})

	resp, err := client.SendSmsWithContext(ctx, req)
	if err != nil {
		return fmt.Errorf("zsms: 腾讯云短信请求失败: %w", err)
	}
	if resp == nil || resp.Response == nil || len(resp.Response.SendStatusSet) == 0 {
		return fmt.Errorf("zsms: 腾讯云短信返回为空")
	}
	status := resp.Response.SendStatusSet[0]
	if strValue(status.Code) != "Ok" {
		return fmt.Errorf("zsms: 腾讯云短信发送失败 code=%s message=%s", strValue(status.Code), strValue(status.Message))
	}

	g.Log().Infof(ctx, "zsms(tencent): 短信已发送 mobile=%s template=%s", maskMobile(msg.Mobile), msg.Template)
	return nil
}
