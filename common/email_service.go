package common

import (
	"fmt"
	"net/smtp"
	"os"
)

// EmailService handles email sending operations
type EmailService struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

// NewEmailService creates a new email service instance from environment variables
func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     os.Getenv("SMTP_FROM"),
	}
}

// SendOTP sends an OTP code to the specified email address
func (s *EmailService) SendOTP(to, otp string) error {
	// Email subject and body
	subject := "Mã OTP đặt lại mật khẩu - Plantheon"
	body := s.buildOTPEmailBody(otp)

	// Use only email address for SMTP From (Gmail requirement)
	fromEmail := s.SMTPUsername

	// Compose email message
	message := fmt.Sprintf("From: %s\r\n", fromEmail)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	// SMTP authentication
	auth := smtp.PlainAuth("", s.SMTPUsername, s.SMTPPassword, s.SMTPHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)
	
	// Debug logging
	fmt.Printf("📧 Attempting to send email to: %s\n", to)
	fmt.Printf("📧 SMTP Host: %s\n", addr)
	fmt.Printf("📧 SMTP Username: %s\n", s.SMTPUsername)
	fmt.Printf("📧 SMTP From: %s\n", fromEmail)
	
	err := smtp.SendMail(addr, auth, fromEmail, []string{to}, []byte(message))
	if err != nil {
		fmt.Printf("❌ Email send error: %v\n", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	fmt.Printf("✅ Email sent successfully to: %s\n", to)
	return nil
}

// buildOTPEmailBody creates the HTML email body for OTP
func (s *EmailService) buildOTPEmailBody(otp string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f9f9f9;
        }
        .header {
            background-color: #4CAF50;
            color: white;
            padding: 20px;
            text-align: center;
            border-radius: 5px 5px 0 0;
        }
        .content {
            background-color: white;
            padding: 30px;
            border-radius: 0 0 5px 5px;
        }
        .otp-code {
            font-size: 32px;
            font-weight: bold;
            color: #4CAF50;
            text-align: center;
            letter-spacing: 5px;
            padding: 20px;
            background-color: #f0f0f0;
            border-radius: 5px;
            margin: 20px 0;
        }
        .warning {
            color: #ff6b6b;
            font-size: 14px;
            margin-top: 20px;
        }
        .footer {
            text-align: center;
            margin-top: 20px;
            font-size: 12px;
            color: #777;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🌱 Plantheon</h1>
        </div>
        <div class="content">
            <h2>Đặt lại mật khẩu</h2>
            <p>Xin chào,</p>
            <p>Bạn đã yêu cầu đặt lại mật khẩu cho tài khoản Plantheon của mình. Vui lòng sử dụng mã OTP dưới đây để tiếp tục:</p>
            
            <div class="otp-code">%s</div>
            
            <p><strong>Lưu ý quan trọng:</strong></p>
            <ul>
                <li>Mã OTP này có hiệu lực trong <strong>5 phút</strong></li>
                <li>Bạn có tối đa <strong>3 lần</strong> nhập mã OTP</li>
                <li>Không chia sẻ mã này với bất kỳ ai</li>
            </ul>
            
            <p class="warning">⚠️ Nếu bạn không yêu cầu đặt lại mật khẩu, vui lòng bỏ qua email này và đảm bảo tài khoản của bạn an toàn.</p>
            
            <div class="footer">
                <p>Email này được gửi tự động, vui lòng không trả lời.</p>
                <p>&copy; 2024 Plantheon. All rights reserved.</p>
            </div>
        </div>
    </div>
</body>
</html>
`, otp)
}
