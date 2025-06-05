package handler

import (
	"net/http"
	"strings"

	"github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// jwtAuth JWT认证中间件
func jwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头或查询参数中获取令牌
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		// 去除Bearer前缀
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[7:]
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, conf.Response{
				Code: http.StatusUnauthorized,
				Msg:  "未授权：未提供令牌",
			})
			c.Abort()
			return
		}

		// 解析令牌
		claims, err := svc.ParseToken(c, tokenString)
		if err != nil {
			var msg string
			switch err {
			case jwt.ErrTokenExpired:
				msg = "未授权：令牌已过期"
			case jwt.ErrTokenInvalid:
				msg = "未授权：无效的令牌"
			case jwt.ErrTokenRevoked:
				msg = "未授权：令牌已被撤销"
			default:
				msg = "未授权：" + err.Error()
			}

			c.JSON(http.StatusUnauthorized, conf.Response{
				Code: http.StatusUnauthorized,
				Msg:  msg,
			})
			c.Abort()
			return
		}

		// 确保是访问令牌
		if claims.Type != jwt.AccessToken {
			c.JSON(http.StatusUnauthorized, conf.Response{
				Code: http.StatusUnauthorized,
				Msg:  "未授权：需要访问令牌",
			})
			c.Abort()
			return
		}

		// 自动延长令牌有效期
		_ = svc.ExtendAccessToken(c, claims.UserID)

		// 将用户ID和角色存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// 角色验证中间件
func roleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, conf.Response{
				Code: http.StatusUnauthorized,
				Msg:  "未授权：未找到用户角色",
			})
			c.Abort()
			return
		}

		// 检查角色是否在允许的角色列表中
		roleStr := role.(string)
		for _, r := range roles {
			if r == roleStr {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, conf.Response{
			Code: http.StatusForbidden,
			Msg:  "禁止访问：权限不足",
		})
		c.Abort()
	}
}

// 从上下文中获取用户ID
func getUserIDFromContext(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(int64)
}

// API签名验证中间件
func signVerify(c *gin.Context) {
	// TODO: 实现API签名验证
	c.Next()
}

// refreshToken 刷新访问令牌
func refreshToken(c *gin.Context) {
	// 从请求中获取刷新令牌
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 使用刷新令牌获取新的令牌对
	tokenPair, err := svc.RefreshAccessToken(c, req.RefreshToken)
	if err != nil {
		var statusCode int
		var msg string

		switch err {
		case jwt.ErrRefreshTokenInvalid, jwt.ErrTokenExpired, jwt.ErrTokenRevoked:
			statusCode = http.StatusUnauthorized
			msg = "刷新令牌无效或已过期，请重新登录"
		default:
			statusCode = http.StatusInternalServerError
			msg = "刷新令牌失败: " + err.Error()
		}

		c.JSON(statusCode, conf.Response{
			Code: statusCode,
			Msg:  msg,
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "令牌刷新成功",
		Data: tokenPair,
	})
}

// adminCheck 管理员验证中间件
func adminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户角色
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, conf.Response{
				Code: http.StatusUnauthorized,
				Msg:  "未授权：未找到用户角色",
			})
			c.Abort()
			return
		}

		// 检查是否为管理员
		if role.(string) != "admin" {
			c.JSON(http.StatusForbidden, conf.Response{
				Code: http.StatusForbidden,
				Msg:  "禁止访问：仅管理员可操作",
			})
			c.Abort()
			return
		}

		c.Next()
	}
} 