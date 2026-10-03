package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"ui-design-backend/internal/config"
)

type Mailer struct {
	cfg *config.Config
}

func NewMailer(cfg *config.Config) *Mailer {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) SendPasswordReset(toEmail, displayName, resetURL string) error {
	subject := "Yêu cầu đặt lại mật khẩu của bạn tại Mộc"

	htmlBody, err := renderResetEmailHTML(displayName, resetURL)
	if err != nil {
		return fmt.Errorf("failed to render reset email template: %w", err)
	}

	headers := make(map[string]string)
	headers["From"] = m.cfg.MailFrom
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message bytes.Buffer
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(htmlBody)

	addr := fmt.Sprintf("%s:%d", m.cfg.SMTPHost, m.cfg.SMTPPort)

	var auth smtp.Auth
	if m.cfg.SMTPUser != "" || m.cfg.SMTPPass != "" {
		auth = smtp.PlainAuth("", m.cfg.SMTPUser, m.cfg.SMTPPass, m.cfg.SMTPHost)
	}

	fromAddr := extractEmail(m.cfg.MailFrom)

	return smtp.SendMail(addr, auth, fromAddr, []string{toEmail}, message.Bytes())
}

func extractEmail(from string) string {
	start := strings.Index(from, "<")
	end := strings.Index(from, ">")
	if start != -1 && end != -1 && end > start {
		return from[start+1 : end]
	}
	return strings.TrimSpace(from)
}

func renderResetEmailHTML(displayName, resetURL string) (string, error) {
	tmpl := `<!DOCTYPE html>
<html lang="vi">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Đặt lại mật khẩu</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; background-color: #f7f9f4; margin: 0; padding: 40px 20px; color: #2d382e; }
    .card { max-width: 520px; margin: 0 auto; background: #ffffff; border-radius: 16px; border: 1px solid #e1e7db; padding: 36px 32px; box-shadow: 0 10px 30px rgba(34, 53, 39, 0.05); }
    .logo { font-size: 24px; font-weight: 700; color: #255036; margin-bottom: 24px; letter-spacing: -0.5px; }
    h1 { font-size: 20px; font-weight: 600; color: #1c261e; margin: 0 0 16px; }
    p { font-size: 14px; line-height: 1.6; color: #526154; margin: 0 0 20px; }
    .btn { display: inline-block; background-color: #2b5e3f; color: #ffffff !important; text-decoration: none; padding: 12px 28px; border-radius: 9999px; font-weight: 500; font-size: 14px; margin: 8px 0 24px; }
    .btn:hover { background-color: #224c32; }
    .note { font-size: 12px; color: #7f8f81; border-top: 1px solid #edf1ea; padding-top: 18px; margin-top: 24px; }
    .link-alt { word-break: break-all; color: #2b5e3f; font-size: 12px; }
  </style>
</head>
<body>
  <div class="card">
    <div class="logo">Mộc</div>
    <h1>Chào {{.DisplayName}},</h1>
    <p>Chúng tôi nhận được yêu cầu đặt lại mật khẩu cho tài khoản của bạn tại <strong>Mộc</strong>.</p>
    <p>Nhấp vào nút bên dưới để tiến hành tạo mật khẩu mới. Liên kết này sẽ hết hạn trong vòng <strong>60 phút</strong>.</p>
    <div>
      <a href="{{.ResetURL}}" class="btn" target="_blank">Đặt lại mật khẩu của tôi</a>
    </div>
    <p>Nếu bạn không gửi yêu cầu này, vui lòng bỏ qua email. Mật khẩu của bạn vẫn an toàn tuyệt đối.</p>
    <div class="note">
      <p>Nếu nút bấm trên không hoạt động, bạn có thể sao chép liên kết này vào trình duyệt:</p>
      <p><a href="{{.ResetURL}}" class="link-alt">{{.ResetURL}}</a></p>
    </div>
  </div>
</body>
</html>`

	t, err := template.New("reset_email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	data := map[string]string{
		"DisplayName": displayName,
		"ResetURL":    resetURL,
	}
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
