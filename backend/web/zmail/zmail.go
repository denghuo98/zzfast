// Package zmail 提供 SMTP 邮件发送能力。
//
// 设计目标：发送动作（连接/认证/编码/重试日志）由框架封装，
// 业务方只需要提供收件人、主题和 HTML 内容（通常来自业务自己的模板）。
//
// 配置示例：
//
//	mail:
//	  enabled: true
//	  host: smtp.163.com
//	  port: 465
//	  ssl: true
//	  username: your_account@163.com
//	  password: SMTP授权码
//	  from: your_account@163.com
//	  fromName: 系统通知
//	  to:
//	    - receive@example.com
package zmail

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"html/template"
	"mime"
	"net"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	webSchema "github.com/denghuo98/zzfast/zschema/zweb"
)

// Message 待发送邮件内容。
// To 为空时使用配置中的 mail.to；Subject/HTML 为必填，Text 为空时自动从 HTML 提取纯文本。
type Message struct {
	To      []string
	Subject string
	HTML    string
	Text    string
}

var config *webSchema.MailConfig

// SetConfig 由系统配置加载流程注入。
func SetConfig(c *webSchema.MailConfig) {
	config = c
}

// GetConfig 获取当前邮件配置（可能为 nil）。
func GetConfig() *webSchema.MailConfig {
	return config
}

// Enabled 邮件通知是否已启用且配置完整。
func Enabled() bool {
	return config != nil && config.Enabled
}

// Send 同步发送邮件，失败返回错误。
func Send(ctx context.Context, msg Message) error {
	if !Enabled() {
		return nil
	}
	cfg := GetConfig()
	to := resolveRecipients(cfg, msg)
	if len(to) == 0 {
		return errors.New("zmail: 未配置收件人（mail.to 为空）")
	}
	raw, err := buildMessage(cfg, msg, to)
	if err != nil {
		return err
	}
	if err = sendSMTP(ctx, cfg, to, raw); err != nil {
		return err
	}
	g.Log().Infof(ctx, "zmail: 邮件已发送 %s", strings.Join(to, ","))
	return nil
}

// SendAsync 异步发送邮件，仅记录错误，不阻塞调用方。
func SendAsync(ctx context.Context, msg Message) {
	if !Enabled() {
		return
	}
	cfg := GetConfig()
	to := resolveRecipients(cfg, msg)
	if len(to) == 0 {
		g.Log().Errorf(ctx, "zmail: 未配置收件人（mail.to 为空），跳过发送")
		return
	}
	raw, err := buildMessage(cfg, msg, to)
	if err != nil {
		g.Log().Errorf(ctx, "zmail: 构建邮件失败: %v", err)
		return
	}
	go func(cfg *webSchema.MailConfig, to []string, raw []byte) {
		sendCtx := gctx.New()
		if err := sendSMTP(sendCtx, cfg, to, raw); err != nil {
			g.Log().Errorf(sendCtx, "zmail: 邮件发送失败 to=%s err=%v", strings.Join(to, ","), err)
			return
		}
		g.Log().Infof(sendCtx, "zmail: 邮件已发送 %s", strings.Join(to, ","))
	}(cfg, to, raw)
}

// SendTemplate 使用 html/template 渲染模板后同步发送。
func SendTemplate(ctx context.Context, tpl string, data any, msg Message) error {
	rendered, err := render(tpl, data)
	if err != nil {
		return err
	}
	msg.HTML = rendered
	return Send(ctx, msg)
}

// SendTemplateAsync 使用 html/template 渲染模板后异步发送。
func SendTemplateAsync(ctx context.Context, tpl string, data any, msg Message) {
	rendered, err := render(tpl, data)
	if err != nil {
		g.Log().Errorf(ctx, "zmail: 渲染邮件模板失败: %v", err)
		return
	}
	msg.HTML = rendered
	SendAsync(ctx, msg)
}

func render(tpl string, data any) (string, error) {
	t, err := template.New("mail").Option("missingkey=zero").Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func resolveRecipients(cfg *webSchema.MailConfig, msg Message) []string {
	raw := msg.To
	if len(raw) == 0 {
		raw = cfg.To
	}
	seen := make(map[string]struct{}, len(raw))
	to := make([]string, 0, len(raw))
	for _, addr := range raw {
		addr = strings.TrimSpace(addr)
		if addr == "" {
			continue
		}
		if _, ok := seen[addr]; ok {
			continue
		}
		seen[addr] = struct{}{}
		to = append(to, addr)
	}
	return to
}

func fromAddress(cfg *webSchema.MailConfig) string {
	if cfg.From != "" {
		return cfg.From
	}
	return cfg.Username
}

func buildMessage(cfg *webSchema.MailConfig, msg Message, to []string) ([]byte, error) {
	from := fromAddress(cfg)
	if from == "" {
		return nil, errors.New("zmail: 未配置发件人（mail.from / mail.username）")
	}
	if strings.TrimSpace(msg.Subject) == "" {
		return nil, errors.New("zmail: 邮件主题不能为空")
	}

	fromHeader := from
	if cfg.FromName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", mime.QEncoding.Encode("UTF-8", cfg.FromName), from)
	}
	subject := mime.QEncoding.Encode("UTF-8", msg.Subject)
	boundary := randomHex(16)
	text := msg.Text
	if strings.TrimSpace(text) == "" {
		text = htmlToText(msg.HTML)
	}

	var b bytes.Buffer
	b.WriteString("From: " + fromHeader + "\r\n")
	b.WriteString("To: " + strings.Join(to, ", ") + "\r\n")
	b.WriteString("Subject: " + subject + "\r\n")
	b.WriteString("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n")
	b.WriteString("Message-ID: <" + randomHex(12) + "@" + messageDomain(cfg.Host) + ">\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString(`Content-Type: multipart/alternative; boundary="` + boundary + `"` + "\r\n")
	b.WriteString("\r\n")

	writePart := func(contentType, body string) {
		b.WriteString("--" + boundary + "\r\n")
		b.WriteString("Content-Type: " + contentType + `; charset="UTF-8"` + "\r\n")
		b.WriteString("Content-Transfer-Encoding: base64\r\n")
		b.WriteString("\r\n")
		b.WriteString(wrapBase64(body))
		b.WriteString("\r\n")
	}
	writePart("text/plain", text)
	writePart("text/html", msg.HTML)
	b.WriteString("--" + boundary + "--\r\n")

	return b.Bytes(), nil
}

func sendSMTP(ctx context.Context, cfg *webSchema.MailConfig, to []string, raw []byte) (err error) {
	if cfg.Host == "" {
		return errors.New("zmail: 未配置 SMTP 服务器")
	}
	port := cfg.Port
	if port <= 0 {
		port = 465
	}
	timeout := time.Duration(cfg.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(port))
	dialer := &net.Dialer{Timeout: timeout}

	var conn net.Conn
	if cfg.SSL {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{ServerName: cfg.Host})
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("zmail: 连接 %s 失败: %w", addr, err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(timeout))

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return fmt.Errorf("zmail: 创建 SMTP 客户端失败: %w", err)
	}
	defer client.Close()

	if !cfg.SSL {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err = client.StartTLS(&tls.Config{ServerName: cfg.Host}); err != nil {
				return fmt.Errorf("zmail: STARTTLS 失败: %w", err)
			}
		}
	}

	if cfg.Username != "" {
		if ok, _ := client.Extension("AUTH"); ok {
			auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
			if err = client.Auth(auth); err != nil {
				return fmt.Errorf("zmail: SMTP 认证失败: %w", err)
			}
		}
	}

	if err = client.Mail(fromAddress(cfg)); err != nil {
		return fmt.Errorf("zmail: MAIL FROM 失败: %w", err)
	}
	for _, rcpt := range to {
		if err = client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("zmail: RCPT TO %s 失败: %w", rcpt, err)
		}
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("zmail: DATA 失败: %w", err)
	}
	if _, err = writer.Write(raw); err != nil {
		return fmt.Errorf("zmail: 写入邮件内容失败: %w", err)
	}
	if err = writer.Close(); err != nil {
		return fmt.Errorf("zmail: 提交邮件内容失败: %w", err)
	}
	if err = client.Quit(); err != nil {
		// 部分服务端在 QUIT 前已关闭连接，此时邮件通常已投递成功。
		g.Log().Debugf(ctx, "zmail: QUIT 失败（可忽略）: %v", err)
	}
	return nil
}

func wrapBase64(s string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	const width = 76
	var out strings.Builder
	for i := 0; i < len(encoded); i += width {
		end := i + width
		if end > len(encoded) {
			end = len(encoded)
		}
		out.WriteString(encoded[i:end])
		out.WriteString("\r\n")
	}
	return strings.TrimSuffix(out.String(), "\r\n")
}

var (
	htmlTagRe   = regexp.MustCompile(`(?s)<[^>]*>`)
	htmlBreakRe = regexp.MustCompile(`(?i)<(br|/p|/div|/tr|/td|/h[1-6])[^>]*>`)
)

func htmlToText(s string) string {
	s = htmlBreakRe.ReplaceAllString(s, "\n")
	s = htmlTagRe.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.Join(strings.Fields(line), " ")
		if line != "" {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

func randomHex(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	return hex.EncodeToString(buf)
}

func messageDomain(host string) string {
	if host == "" {
		return "localhost"
	}
	return host
}
