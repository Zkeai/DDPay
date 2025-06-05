package service

import (
	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/Zkeai/DDPay/internal/repo"
	"github.com/Zkeai/DDPay/pkg/email"
	"github.com/Zkeai/DDPay/pkg/jwt"
	"github.com/Zkeai/DDPay/pkg/oauth"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Service 服务层
type Service struct {
	conf           *conf.Conf
	repo           *repo.Repo
	tgService      *tgbotapi.BotAPI
	jwtService     *jwt.JWTService
	emailService   *email.Service
	oauthManager   *oauth.OAuthManager
	subsiteService SubsiteService
	membershipService MembershipService
}

// NewService 创建服务
func NewService(conf *conf.Conf) *Service {
	// 初始化仓储层
	r := repo.NewRepo(conf)

	// 创建服务实例
	s := &Service{
		conf:         conf,
		repo:         r,
		tgService:    conf.Tg,
		jwtService:   jwt.NewJWTService(conf.JWT),
	}
	
	// 设置Email服务
	if conf.Email != nil {
		s.emailService = email.NewService(*conf.Email)
	}

	// 创建OAuth服务管理器
	if conf.OAuth != nil {
		s.oauthManager = oauth.NewOAuthManager(conf.OAuth)
	}
	
	// 初始化分站服务
	s.subsiteService = NewSubsiteService(r.GetSubsiteRepo(), r.GetUserRepo())
	
	// 初始化会员服务
	s.membershipService = NewMembershipService(s)
	
	return s
}

// GetSubsiteService 获取分站服务接口
func (s *Service) GetSubsiteService() SubsiteService {
	return s.subsiteService
}

// GetMembershipService 获取会员服务接口
func (s *Service) GetMembershipService() MembershipService {
	return s.membershipService
}
