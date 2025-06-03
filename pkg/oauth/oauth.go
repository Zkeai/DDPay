package oauth

import (
	"context"
	"errors"

	"github.com/Zkeai/DDPay/internal/conf"
)

// 支持的OAuth提供商类型
const (
	ProviderGithub = "github"
	ProviderGoogle = "google"
)

var (
	ErrUnsupportedProvider = errors.New("不支持的OAuth提供商")
	ErrInvalidConfig       = errors.New("无效的OAuth配置")
)

// OAuthManager OAuth服务管理器
type OAuthManager struct {
	config        *conf.OAuthConfig
	githubService *GitHubService
	googleService *GoogleService
}

// NewOAuthManager 创建OAuth服务管理器
func NewOAuthManager(config *conf.OAuthConfig) *OAuthManager {
	manager := &OAuthManager{
		config: config,
	}

	// 初始化GitHub服务
	if config != nil && config.Github.ClientID != "" && config.Github.ClientSecret != "" {
		manager.githubService = NewGitHubService(&config.Github)
	}

	// 初始化Google服务
	if config != nil && config.Google.ClientID != "" && config.Google.ClientSecret != "" {
		manager.googleService = NewGoogleService(&config.Google)
	}

	return manager
}

// GetOAuthURL 获取OAuth授权URL
func (m *OAuthManager) GetOAuthURL(provider, state string) (string, error) {
	switch provider {
	case ProviderGithub:
		if m.githubService == nil {
			return "", ErrInvalidConfig
		}
		return m.githubService.GetAuthURL(state), nil
	case ProviderGoogle:
		if m.googleService == nil {
			return "", ErrInvalidConfig
		}
		return m.googleService.GetAuthURL(state), nil
	default:
		return "", ErrUnsupportedProvider
	}
}

// ExchangeCode 使用授权码交换访问令牌
func (m *OAuthManager) ExchangeCode(ctx context.Context, provider, code string) (string, error) {
	switch provider {
	case ProviderGithub:
		if m.githubService == nil {
			return "", ErrInvalidConfig
		}
		return m.githubService.ExchangeCode(ctx, code)
	case ProviderGoogle:
		if m.googleService == nil {
			return "", ErrInvalidConfig
		}
		return m.googleService.ExchangeCode(ctx, code)
	default:
		return "", ErrUnsupportedProvider
	}
}

// GetUserInfo 获取用户信息
func (m *OAuthManager) GetUserInfo(ctx context.Context, provider, accessToken string) (map[string]interface{}, error) {
	switch provider {
	case ProviderGithub:
		if m.githubService == nil {
			return nil, ErrInvalidConfig
		}

		user, err := m.githubService.GetUserInfo(ctx, accessToken)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"id":       user.ID,
			"username": user.Login,
			"name":     user.Name,
			"email":    user.Email,
			"avatar":   user.AvatarURL,
		}, nil

	case ProviderGoogle:
		if m.googleService == nil {
			return nil, ErrInvalidConfig
		}

		user, err := m.googleService.GetUserInfo(ctx, accessToken)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"id":       user.ID,
			"username": user.Name,
			"name":     user.Name,
			"email":    user.Email,
			"avatar":   user.Picture,
		}, nil

	default:
		return nil, ErrUnsupportedProvider
	}
} 