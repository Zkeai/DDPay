package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Zkeai/DDPay/internal/conf"
)

// GitHubService GitHub OAuth服务
type GitHubService struct {
	Config *conf.GithubConfig
}

// NewGitHubService 创建GitHub OAuth服务
func NewGitHubService(config *conf.GithubConfig) *GitHubService {
	return &GitHubService{
		Config: config,
	}
}

// GitHubUser GitHub用户信息
type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GetAuthURL 获取GitHub授权URL
func (s *GitHubService) GetAuthURL(state string) string {
	baseURL := "https://github.com/login/oauth/authorize"
	params := url.Values{}
	params.Add("client_id", s.Config.ClientID)
	params.Add("redirect_uri", s.Config.RedirectURI)
	params.Add("scope", s.Config.Scopes)
	params.Add("state", state)

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// ExchangeCode 使用授权码交换访问令牌
func (s *GitHubService) ExchangeCode(ctx context.Context, code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", s.Config.ClientID)
	data.Set("client_secret", s.Config.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", s.Config.RedirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub响应错误，状态码: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		Error       string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Error != "" {
		return "", fmt.Errorf("GitHub错误: %s", result.Error)
	}

	return result.AccessToken, nil
}

// GetUserInfo 获取GitHub用户信息
func (s *GitHubService) GetUserInfo(ctx context.Context, accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API错误，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	// 如果用户没有设置公开邮箱，则尝试获取邮箱列表
	if user.Email == "" {
		userEmail, err := s.GetUserEmails(ctx, accessToken)
		if err == nil && userEmail != "" {
			user.Email = userEmail
		}
	}

	return &user, nil
}

// GetUserEmails 获取GitHub用户邮箱列表
func (s *GitHubService) GetUserEmails(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API错误，状态码: %d", resp.StatusCode)
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("解析邮箱信息失败: %w", err)
	}

	// 优先返回已验证的主邮箱
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	// 返回任何已验证的邮箱
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	// 没有找到已验证的邮箱
	if len(emails) > 0 {
		return emails[0].Email, nil
	}

	return "", nil
} 