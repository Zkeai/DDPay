package email

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

// Config 邮件配置
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

// Service 邮件服务
type Service struct {
	config Config
}

// NewService 创建邮件服务
func NewService(config Config) *Service {
	return &Service{
		config: config,
	}
}

// SendVerificationCode 发送验证码邮件
func (s *Service) SendVerificationCode(to, code, codeType string) error {
	var subject, body string

	switch codeType {
	case "register":
		subject = "DDPay - 注册验证码"
		body = fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 5px;">
			<h2 style="color: #333;">DDPay - 注册验证码</h2>
			<p>您好，</p>
			<p>感谢您注册DDPay账户。请使用以下验证码完成注册：</p>
			<p style="font-size: 24px; font-weight: bold; background-color: #f5f5f5; padding: 10px; text-align: center; letter-spacing: 5px;">%s</p>
			<p>该验证码将在10分钟后过期。</p>
			<p>如果您没有请求此验证码，请忽略此邮件。</p>
			<p>祝好，<br>DDPay团队</p>
		</div>
		`, code)
	case "reset_password":
		subject = "DDPay - 重置密码验证码"
		body = fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 5px;">
			<h2 style="color: #333;">DDPay - 重置密码验证码</h2>
			<p>您好，</p>
			<p>您正在重置DDPay账户密码。请使用以下验证码完成密码重置：</p>
			<p style="font-size: 24px; font-weight: bold; background-color: #f5f5f5; padding: 10px; text-align: center; letter-spacing: 5px;">%s</p>
			<p>该验证码将在10分钟后过期。</p>
			<p>如果您没有请求此验证码，请忽略此邮件，但您可能需要检查您的账户安全。</p>
			<p>祝好，<br>DDPay团队</p>
		</div>
		`, code)
	default:
		subject = "DDPay - 验证码"
		body = fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 5px;">
			<h2 style="color: #333;">DDPay - 验证码</h2>
			<p>您好，</p>
			<p>您的验证码是：</p>
			<p style="font-size: 24px; font-weight: bold; background-color: #f5f5f5; padding: 10px; text-align: center; letter-spacing: 5px;">%s</p>
			<p>该验证码将在10分钟后过期。</p>
			<p>如果您没有请求此验证码，请忽略此邮件。</p>
			<p>祝好，<br>DDPay团队</p>
		</div>
		`, code)
	}

	return s.SendEmail(to, subject, body)
}

// SendEmail 使用gomail库发送邮件
func (s *Service) SendEmail(to, subject, body string) error {
	log.Printf("准备发送邮件给 %s，主题：%s", to, subject)
	log.Printf("使用SMTP服务器: %s:%d", s.config.Host, s.config.Port)
	
	// 创建新邮件
	m := gomail.NewMessage()
	
	// 设置发件人
	if s.config.FromName != "" {
		m.SetHeader("From", m.FormatAddress(s.config.From, s.config.FromName))
	} else {
		m.SetHeader("From", s.config.From)
	}
	
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	
	// 创建SMTP发送器
	dialer := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	
	// QQ邮箱需要SSL
	dialer.SSL = true
	
	log.Printf("开始发送邮件...")
	
	// 发送邮件
	if err := dialer.DialAndSend(m); err != nil {
		log.Printf("邮件发送失败: %v", err)
		return fmt.Errorf("发送邮件失败: %v", err)
	}
	
	log.Printf("邮件发送成功!")
	return nil
} 