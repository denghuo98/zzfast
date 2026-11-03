package zsms

import (
	"context"
	"fmt"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gogf/gf/v2/frame/g"

	webSchema "github.com/denghuo98/zzfast/zschema/zweb"
)

type aliyunDriver struct {
	cfg *webSchema.SmsConfig
}

// Send 通过阿里云短信服务发送验证码。
func (d *aliyunDriver) Send(ctx context.Context, msg Message) error {
	if d.cfg.AliyunAccessKeyId == "" || d.cfg.AliyunAccessKeySecret == "" {
		return fmt.Errorf("zsms: 阿里云短信缺少 AccessKeyId/AccessKeySecret")
	}

	client, err := dysmsapi.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(d.cfg.AliyunAccessKeyId),
		AccessKeySecret: tea.String(d.cfg.AliyunAccessKeySecret),
		Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
		Protocol:        tea.String("https"),
	})
	if err != nil {
		return fmt.Errorf("zsms: 创建阿里云短信客户端失败: %w", err)
	}

	timeoutSec := timeout(d.cfg)
	if deadline, ok := ctx.Deadline(); ok {
		remaining := int(time.Until(deadline).Seconds())
		if remaining > 0 && remaining < timeoutSec {
			timeoutSec = remaining
		}
	}
	timeoutMs := timeoutSec * 1000
	req := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(msg.Mobile),
		SignName:      tea.String(signName(d.cfg, DriveAliyun)),
		TemplateCode:  tea.String(msg.Template),
		TemplateParam: tea.String(fmt.Sprintf(`{"code":%q}`, msg.Code)),
	}

	resp, err := client.SendSmsWithOptions(req, &util.RuntimeOptions{
		ConnectTimeout: tea.Int(timeoutMs),
		ReadTimeout:    tea.Int(timeoutMs),
	})
	if err != nil {
		return fmt.Errorf("zsms: 阿里云短信请求失败: %w", err)
	}
	if resp == nil || resp.Body == nil {
		return fmt.Errorf("zsms: 阿里云短信返回为空")
	}
	if tea.StringValue(resp.Body.Code) != "OK" {
		return fmt.Errorf("zsms: 阿里云短信发送失败 code=%s message=%s", tea.StringValue(resp.Body.Code), tea.StringValue(resp.Body.Message))
	}

	g.Log().Infof(ctx, "zsms(aliyun): 短信已发送 mobile=%s template=%s", maskMobile(msg.Mobile), msg.Template)
	return nil
}
