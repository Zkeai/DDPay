package model

import "time"

// User 用户模型
type User struct {
	ID            int64     `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	Password      string    `json:"-" db:"password"` // 不暴露密码
	Username      string    `json:"username" db:"username"`
	Avatar        string    `json:"avatar" db:"avatar"`
	Role          string    `json:"role" db:"role"`
	Status        int       `json:"status" db:"status"`
	EmailVerified int       `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt   time.Time `json:"last_login_at" db:"last_login_at"`
	LastLoginIP   string    `json:"last_login_ip" db:"last_login_ip"`
}

// UserSession 用户会话模型
type UserSession struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	Token          string    `json:"token" db:"token"`
	IP             string    `json:"ip" db:"ip"`
	UserAgent      string    `json:"user_agent" db:"user_agent"`
	ExpiresAt      time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	LastActivityAt time.Time `json:"last_activity_at" db:"last_activity_at"`
}

// OAuthAccount OAuth账号关联模型
type OAuthAccount struct {
	ID              int64     `json:"id" db:"id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	Provider        string    `json:"provider" db:"provider"`
	ProviderUserID  string    `json:"provider_user_id" db:"provider_user_id"`
	ProviderUsername string   `json:"provider_username" db:"provider_username"`
	ProviderEmail   string    `json:"provider_email" db:"provider_email"`
	ProviderAvatar  string    `json:"provider_avatar" db:"provider_avatar"`
	AccessToken     string    `json:"-" db:"access_token"` // 不暴露令牌
	RefreshToken    string    `json:"-" db:"refresh_token"` // 不暴露令牌
	TokenExpiresAt  time.Time `json:"token_expires_at" db:"token_expires_at"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// VerificationCode 验证码模型
type VerificationCode struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Code      string    `json:"code" db:"code"`
	Type      string    `json:"type" db:"type"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Used      int       `json:"used" db:"used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// LoginLog 登录日志模型
type LoginLog struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	LoginType  string    `json:"login_type" db:"login_type"`
	IP         string    `json:"ip" db:"ip"`
	UserAgent  string    `json:"user_agent" db:"user_agent"`
	Status     int       `json:"status" db:"status"`
	FailReason string    `json:"fail_reason" db:"fail_reason"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// UserProfile 用户配置文件（返回给前端的安全用户信息）
type UserProfile struct {
	ID            int64     `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Avatar        string    `json:"avatar"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	LastLoginAt   time.Time `json:"last_login_at"`
}

// ToProfile 将用户模型转换为配置文件
func (u *User) ToProfile() *UserProfile {
	return &UserProfile{
		ID:            u.ID,
		Email:         u.Email,
		Username:      u.Username,
		Avatar:        u.Avatar,
		Role:          u.Role,
		EmailVerified: u.EmailVerified == 1,
		CreatedAt:     u.CreatedAt,
		LastLoginAt:   u.LastLoginAt,
	}
}