package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Zkeai/DDPay/internal/model"
	"github.com/Zkeai/DDPay/pkg/jwt"
)

// 获取当前北京时间
func getNow() time.Time {
	return time.Now().In(jwt.BeijingLocation)
}

var (
	ErrInvalidCredentials = errors.New("邮箱或密码错误")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrEmailAlreadyExists = errors.New("邮箱已被注册")
	ErrInvalidCode        = errors.New("验证码无效")
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Username  string `json:"username" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Type  string `json:"type" binding:"required"` // register, reset_password
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// OAuthLoginRequest OAuth登录请求
type OAuthLoginRequest struct {
	Provider        string `json:"provider" binding:"required"`
	ProviderUserID  string `json:"provider_user_id" binding:"required"`
	ProviderToken   string `json:"provider_token" binding:"required"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Avatar          string `json:"avatar"`
}

// 用户注册
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*model.UserProfile, *jwt.TokenPair, error) {
	// 验证码验证
	code, err := s.repo.GetVerificationCode(ctx, req.Email, "register")
	if err != nil {
		return nil, nil, err
	}
	
	if code.Code != req.Code {
		return nil, nil, ErrInvalidCode
	}
	
	// 检查邮箱是否已注册
	_, err = s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, nil, ErrEmailAlreadyExists
	}
	
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}
	
	// 创建用户
	now := getNow()
	user := &model.User{
		Email:         req.Email,
		Password:      string(hashedPassword),
		Username:      req.Username,
		Status:        1, // 正常状态
		EmailVerified: 1, // 已验证
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	
	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	
	// 标记验证码已使用
	_ = s.repo.MarkVerificationCodeAsUsed(ctx, code.ID)
	
	// 获取用户信息
	user, err = s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	
	// 生成令牌对
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, nil, err
	}
	
	return user.ToProfile(), tokenPair, nil
}

// 用户登录
func (s *Service) Login(ctx context.Context, req *LoginRequest, ip string, userAgent string) (*model.UserProfile, *jwt.TokenPair, error) {
	// 查找用户
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// 记录失败日志
		_ = s.repo.CreateLoginLog(ctx, &model.LoginLog{
			LoginType:  "email",
			IP:         ip,
			UserAgent:  userAgent,
			Status:     0, // 失败
			FailReason: "用户不存在",
			CreatedAt:  getNow(),
		})
		return nil, nil, ErrInvalidCredentials
	}
	
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// 记录失败日志
		_ = s.repo.CreateLoginLog(ctx, &model.LoginLog{
			UserID:     user.ID,
			LoginType:  "email",
			IP:         ip,
			UserAgent:  userAgent,
			Status:     0, // 失败
			FailReason: "密码错误",
			CreatedAt:  getNow(),
		})
		return nil, nil, ErrInvalidCredentials
	}
	
	// 更新最后登录信息
	_ = s.repo.UpdateLastLogin(ctx, user.ID, ip)
	
	// 记录成功日志
	_ = s.repo.CreateLoginLog(ctx, &model.LoginLog{
		UserID:    user.ID,
		LoginType: "email",
		IP:        ip,
		UserAgent: userAgent,
		Status:    1, // 成功
		CreatedAt: getNow(),
	})
	
	// 生成令牌对
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, nil, err
	}
	
	return user.ToProfile(), tokenPair, nil
}

// OAuth登录
func (s *Service) OAuthLogin(ctx context.Context, req *OAuthLoginRequest, ip string, userAgent string) (*model.UserProfile, string, error) {
	// 查找OAuth账号
	oauthAccount, err := s.repo.GetOAuthAccount(ctx, req.Provider, req.ProviderUserID)
	
	var user *model.User
	
	if err != nil || oauthAccount == nil {
		// 新用户，需要创建
		if req.Email == "" {
			req.Email = fmt.Sprintf("%s_%s@oauth.user", req.Provider, req.ProviderUserID)
		}
		
		if req.Username == "" {
			req.Username = fmt.Sprintf("%s_user_%s", req.Provider, req.ProviderUserID[:6])
		}
		
		// 检查邮箱是否已存在
		existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
		if err == nil {
			// 邮箱已存在，关联到现有账号
			user = existingUser
		} else {
			// 创建新用户
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(generateRandomPassword(12)), bcrypt.DefaultCost)
			
			now := getNow()
			user = &model.User{
				Email:         req.Email,
				Password:      string(hashedPassword),
				Username:      req.Username,
				Avatar:        req.Avatar,
				Status:        1, // 正常状态
				EmailVerified: 1, // OAuth用户视为已验证
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			
			userID, err := s.repo.CreateUser(ctx, user)
			if err != nil {
				return nil, "", err
			}
			
			user.ID = userID
		}
		
		// 创建OAuth账号关联
		now := getNow()
		oauthAccount = &model.OAuthAccount{
			UserID:          user.ID,
			Provider:        req.Provider,
			ProviderUserID:  req.ProviderUserID,
			ProviderUsername: req.Username,
			ProviderEmail:   req.Email,
			ProviderAvatar:  req.Avatar,
			AccessToken:     req.ProviderToken,
			TokenExpiresAt:  now.Add(24 * time.Hour), // 假设令牌有效期为24小时
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		
		err = s.repo.CreateOAuthAccount(ctx, oauthAccount)
		if err != nil {
			return nil, "", err
		}
	} else {
		// 现有用户，更新OAuth信息
		user, err = s.repo.GetUserByID(ctx, oauthAccount.UserID)
		if err != nil {
			return nil, "", err
		}
		
		// 更新OAuth账号信息
		now := getNow()
		oauthAccount.ProviderUsername = req.Username
		oauthAccount.ProviderEmail = req.Email
		oauthAccount.ProviderAvatar = req.Avatar
		oauthAccount.AccessToken = req.ProviderToken
		oauthAccount.TokenExpiresAt = now.Add(24 * time.Hour)
		oauthAccount.UpdatedAt = now
		
		err = s.repo.UpdateOAuthAccount(ctx, oauthAccount)
		if err != nil {
			return nil, "", err
		}
	}
	
	// 更新最后登录信息
	_ = s.repo.UpdateLastLogin(ctx, user.ID, ip)
	
	// 记录成功日志
	_ = s.repo.CreateLoginLog(ctx, &model.LoginLog{
		UserID:    user.ID,
		LoginType: fmt.Sprintf("oauth_%s", req.Provider),
		IP:        ip,
		UserAgent: userAgent,
		Status:    1, // 成功
		CreatedAt: getNow(),
	})
	
	// 生成JWT令牌
	token, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, "", err
	}
	
	return user.ToProfile(), token, nil
}

// 发送验证码
func (s *Service) SendVerificationCode(ctx context.Context, req *SendCodeRequest) error {
	// 生成6位数验证码
	code := generateVerificationCode()
	
	// 创建验证码记录
	now := getNow()
	verificationCode := &model.VerificationCode{
		Email:     req.Email,
		Code:      code,
		Type:      req.Type,
		ExpiresAt: now.Add(10 * time.Minute), // 10分钟有效期
		Used:      0,
		CreatedAt: now,
	}
	
	err := s.repo.CreateVerificationCode(ctx, verificationCode)
	if err != nil {
		return err
	}
	
	// 发送验证码到用户邮箱
	err = s.emailService.SendVerificationCode(req.Email, code, req.Type)
	if err != nil {
		// 记录邮件发送失败，但不影响API返回
		// 可以在这里添加日志记录
		return fmt.Errorf("发送验证码邮件失败: %v", err)
	}
	
	return nil
}

// 重置密码
func (s *Service) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	// 验证码验证
	code, err := s.repo.GetVerificationCode(ctx, req.Email, "reset_password")
	if err != nil {
		return err
	}
	
	if code.Code != req.Code {
		return ErrInvalidCode
	}
	
	// 查找用户
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return ErrUserNotFound
	}
	
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	// 更新密码
	err = s.repo.UpdatePassword(ctx, user.ID, string(hashedPassword))
	if err != nil {
		return err
	}
	
	// 标记验证码已使用
	_ = s.repo.MarkVerificationCodeAsUsed(ctx, code.ID)
	
	return nil
}

// 获取用户信息
func (s *Service) GetUserProfile(ctx context.Context, userID int64) (*model.UserProfile, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	return user.ToProfile(), nil
}

// 更新用户信息
func (s *Service) UpdateUserProfile(ctx context.Context, userID int64, username string, avatar string) (*model.UserProfile, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	user.Username = username
	if avatar != "" {
		user.Avatar = avatar
	}
	
	user.UpdatedAt = getNow()
	
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	
	return user.ToProfile(), nil
}

// 注销登录
func (s *Service) Logout(ctx context.Context, userID int64) error {
	return s.jwtService.RevokeTokens(userID)
}

// 生成验证码
func generateVerificationCode() string {
	// 生成6位数字验证码
	code := ""
	for i := 0; i < 6; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += num.String()
	}
	return code
}

// 生成随机密码
func generateRandomPassword(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	password := make([]byte, length)
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		password[i] = chars[num.Int64()]
	}
	return string(password)
}

// ParseToken 解析JWT令牌
func (s *Service) ParseToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	return s.jwtService.ParseToken(tokenString)
}

// ExtendAccessToken 延长访问令牌有效期
func (s *Service) ExtendAccessToken(ctx context.Context, userID int64) error {
	return s.jwtService.ExtendAccessToken(userID)
}

// RefreshAccessToken 使用刷新令牌获取新的令牌对
func (s *Service) RefreshAccessToken(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
	return s.jwtService.RefreshAccessToken(refreshToken)
}

// CheckEmailExists 检查邮箱是否已存在
func (s *Service) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	// 调用数据库方法检查邮箱是否存在
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// 如果是未找到用户的错误，表示邮箱不存在
		if errors.Is(err, ErrUserNotFound) {
			return false, nil
		}
		// 其他错误
		return false, err
	}
	
	// 如果找到用户，表示邮箱已存在
	return user != nil, nil
}

// LoginLogRequest 登录日志查询请求
type LoginLogRequest struct {
	UserID    int64     `json:"user_id"`
	IP        string    `json:"ip"`
	Status    *int      `json:"status"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Page      int       `json:"page" binding:"required,min=1"`
	PageSize  int       `json:"page_size" binding:"required,min=1,max=100"`
}

// LoginLogResponse 登录日志查询响应
type LoginLogResponse struct {
	Logs       []*model.LoginLog `json:"logs"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// GetLoginLogs 获取登录日志
func (s *Service) GetLoginLogs(ctx context.Context, req *LoginLogRequest) (*LoginLogResponse, error) {
	// 转换查询参数
	params := make(map[string]interface{})
	
	if req.UserID > 0 {
		params["user_id"] = req.UserID
	}
	
	if req.IP != "" {
		params["ip"] = req.IP
	}
	
	if req.Status != nil {
		params["status"] = *req.Status
	}
	
	if !req.StartTime.IsZero() {
		params["start_time"] = req.StartTime
	}
	
	if !req.EndTime.IsZero() {
		params["end_time"] = req.EndTime
	}
	
	// 查询日志
	logs, total, err := s.repo.GetLoginLogs(ctx, params, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	
	// 计算总页数
	totalPages := (total + req.PageSize - 1) / req.PageSize
	if totalPages < 1 {
		totalPages = 1
	}
	
	return &LoginLogResponse{
		Logs:       logs,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *Service) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// UpdateUser 更新用户信息
func (s *Service) UpdateUser(ctx context.Context, user *model.User) error {
	user.UpdatedAt = getNow()
	return s.repo.UpdateUser(ctx, user)
}

// ListUsers 获取用户列表
func (s *Service) ListUsers(ctx context.Context, page, pageSize int) ([]*model.UserProfile, int, error) {
	users, total, err := s.repo.ListUsers(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	
	// 转换为UserProfile
	profiles := make([]*model.UserProfile, len(users))
	for i, user := range users {
		profiles[i] = user.ToProfile()
	}
	
	return profiles, total, nil
} 