package service

import (
	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/Zkeai/DDPay/internal/repo"
	"github.com/Zkeai/DDPay/pkg/email"
	"github.com/Zkeai/DDPay/pkg/jwt"
	"github.com/Zkeai/DDPay/pkg/oauth"
	"github.com/Zkeai/DDPay/pkg/telegram"
)

// Service 服务层
type Service struct {
	conf         *conf.Conf
	repo         *repo.Repo
	tgService    *telegram.TelegramService
	jwtService   *jwt.JWTService
	emailService *email.Service
	oauthManager *oauth.OAuthManager
}

// NewService 创建服务
func NewService(conf *conf.Conf) *Service {
	// 初始化仓储层
	r := repo.NewRepo(conf)

	// 初始化Telegram服务
	var tgService *telegram.TelegramService
	if conf.Tg != nil {
		tgService = &telegram.TelegramService{
			Bot: conf.Tg,
		}
	}
	
	// 创建邮件服务
	var emailService *email.Service
	if conf.Email != nil {
		emailService = email.NewService(*conf.Email)
	}
	
	// 创建OAuth服务管理器
	var oauthManager *oauth.OAuthManager
	if conf.OAuth != nil {
		oauthManager = oauth.NewOAuthManager(conf.OAuth)
	}

	return &Service{
		conf:         conf,
		repo:         r,
		tgService:    tgService,
		jwtService:   jwt.NewJWTService(conf.JWT),
		emailService: emailService,
		oauthManager: oauthManager,
	}
}
