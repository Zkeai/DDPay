package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/Zkeai/DDPay/pkg/redis"
	"github.com/golang-jwt/jwt/v5"
)

// 定义北京时区常量
var (
	// 北京时区（东八区，UTC+8）
	BeijingLocation = time.FixedZone("CST", 8*60*60)
	
	ErrTokenExpired      = errors.New("令牌已过期")
	ErrTokenInvalid      = errors.New("无效的令牌")
	ErrTokenNotProvided  = errors.New("未提供令牌")
	ErrTokenRevoked      = errors.New("令牌已被撤销")
	ErrRefreshTokenInvalid = errors.New("刷新令牌无效")
)

// TokenType 令牌类型
type TokenType string

const (
	AccessToken  TokenType = "access"  // 访问令牌，有效期短
	RefreshToken TokenType = "refresh" // 刷新令牌，有效期长
)

// Config JWT配置
type Config struct {
	Secret           string        `yaml:"secret"`
	Issuer           string        `yaml:"issuer"`
	ExpirationTime   time.Duration `yaml:"expirationTime"`
	RefreshTokenTime time.Duration `yaml:"refreshTokenTime"`
	Expire           int           `yaml:"expire"` // 以秒为单位的过期时间（兼容旧配置）
}

// Claims 自定义JWT声明
type Claims struct {
	UserID int64     `json:"user_id"`
	Role   string    `json:"role"`
	Type   TokenType `json:"type"` // 令牌类型：access 或 refresh
	jwt.RegisteredClaims
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int64  `json:"expires_in"`    // 访问令牌过期时间（秒）
}

// JWTService JWT服务
type JWTService struct {
	config          *Config
	accessExpiration  time.Duration // 访问令牌有效期 (2小时)
	refreshExpiration time.Duration // 刷新令牌有效期 (30天)
}

// NewJWTService 创建JWT服务
func NewJWTService(config *Config) *JWTService {
	// 访问令牌默认有效期为2小时
	accessExpiration := 2 * time.Hour
	
	// 刷新令牌默认有效期为30天
	refreshExpiration := 30 * 24 * time.Hour
	
	// 如果配置中的Expire大于0，则使用配置中的过期时间（兼容旧配置）
	if config.Expire > 0 {
		accessExpiration = time.Duration(config.Expire) * time.Second
	} else if config.ExpirationTime > 0 {
		// 如果配置中的ExpirationTime大于0，则使用ExpirationTime
		accessExpiration = config.ExpirationTime
	}
	
	// 如果配置中的RefreshTokenTime大于0，则使用RefreshTokenTime
	if config.RefreshTokenTime > 0 {
		refreshExpiration = time.Duration(config.RefreshTokenTime) * time.Second
	}
	
	return &JWTService{
		config:          config,
		accessExpiration:  accessExpiration,
		refreshExpiration: refreshExpiration,
	}
}

// GenerateTokenPair 生成访问令牌和刷新令牌对
func (s *JWTService) GenerateTokenPair(userID int64, role string) (*TokenPair, error) {
	// 生成访问令牌
	accessToken, err := s.generateToken(userID, role, AccessToken, s.accessExpiration)
	if err != nil {
		return nil, err
	}
	
	// 生成刷新令牌
	refreshToken, err := s.generateToken(userID, role, RefreshToken, s.refreshExpiration)
	if err != nil {
		return nil, err
	}
	
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessExpiration.Seconds()),
	}, nil
}

// getNow 获取当前北京时间
func getNow() time.Time {
	return time.Now().In(BeijingLocation)
}

// generateToken 生成JWT令牌
func (s *JWTService) generateToken(userID int64, role string, tokenType TokenType, expiration time.Duration) (string, error) {
	// 使用北京时间作为基准
	now := getNow()
	
	// 设置声明
	claims := &Claims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.config.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	// 生成令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", err
	}

	// 将令牌存入Redis，用于后续验证和撤销
	redisKey := fmt.Sprintf("jwt:%s:%d", string(tokenType), userID)
	err = redis.Set(redisKey, tokenString, expiration)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 续约访问令牌（延长有效期）
func (s *JWTService) ExtendAccessToken(userID int64) error {
	// 获取当前存储的访问令牌
	redisKey := fmt.Sprintf("jwt:%s:%d", string(AccessToken), userID)
	token, err := redis.Get(redisKey)
	if err != nil {
		return err
	}
	
	// 解析令牌以确保有效
	claims, err := s.ParseToken(token)
	if err != nil {
		return err
	}
	
	// 只有访问令牌才能续约
	if claims.Type != AccessToken {
		return errors.New("只能续约访问令牌")
	}
	
	// 延长Redis中的令牌有效期
	return redis.Set(redisKey, token, s.accessExpiration)
}

// RefreshAccessToken 使用刷新令牌获取新的访问令牌
func (s *JWTService) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
	// 解析刷新令牌
	claims, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}
	
	// 确保是刷新令牌
	if claims.Type != RefreshToken {
		return nil, ErrRefreshTokenInvalid
	}
	
	// 检查Redis中的刷新令牌是否匹配
	redisKey := fmt.Sprintf("jwt:%s:%d", string(RefreshToken), claims.UserID)
	storedToken, err := redis.Get(redisKey)
	if err != nil || storedToken != refreshToken {
		return nil, ErrRefreshTokenInvalid
	}
	
	// 生成新的令牌对
	return s.GenerateTokenPair(claims.UserID, claims.Role)
}

// ParseToken 解析JWT令牌
func (s *JWTService) ParseToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	// 类型断言
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	// 检查令牌是否被撤销
	redisKey := fmt.Sprintf("jwt:%s:%d", string(claims.Type), claims.UserID)
	storedToken, err := redis.Get(redisKey)
	if err != nil || storedToken != tokenString {
		return nil, ErrTokenRevoked
	}

	return claims, nil
}

// RevokeTokens 撤销用户的所有令牌
func (s *JWTService) RevokeTokens(userID int64) error {
	// 撤销访问令牌
	accessKey := fmt.Sprintf("jwt:%s:%d", string(AccessToken), userID)
	err1 := redis.Del(accessKey)
	
	// 撤销刷新令牌
	refreshKey := fmt.Sprintf("jwt:%s:%d", string(RefreshToken), userID)
	err2 := redis.Del(refreshKey)
	
	if err1 != nil {
		return err1
	}
	return err2
}

// GenerateToken 兼容旧接口，生成访问令牌
func (s *JWTService) GenerateToken(userID int64, role string) (string, error) {
	return s.generateToken(userID, role, AccessToken, s.accessExpiration)
}

// GetUserIDFromToken 从令牌中获取用户ID
func (s *JWTService) GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
} 