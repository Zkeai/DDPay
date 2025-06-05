package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/internal/service"
	"github.com/gin-gonic/gin"
)

// 注册
// @Summary 用户注册
// @Description 通过邮箱验证码注册用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body service.RegisterRequest true "注册信息"
// @Success 200 {object} conf.Response
// @Router /api/v1/user/register [post]
func register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	profile, tokenPair, err := svc.Register(c, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "注册成功",
		Data: gin.H{
			"user":          profile,
			"access_token":  tokenPair.AccessToken,
			"refresh_token": tokenPair.RefreshToken,
			"expires_in":    tokenPair.ExpiresIn,
		},
	})
}

// 登录
// @Summary 用户登录
// @Description 通过邮箱密码登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body service.LoginRequest true "登录信息"
// @Success 200 {object} conf.Response
// @Router /api/v1/user/login [post]
func login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	profile, tokenPair, err := svc.Login(c, &req, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, conf.Response{
			Code: http.StatusUnauthorized,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "登录成功",
		Data: gin.H{
			"user":          profile,
			"access_token":  tokenPair.AccessToken,
			"refresh_token": tokenPair.RefreshToken,
			"expires_in":    tokenPair.ExpiresIn,
		},
	})
}

// OAuth登录
// @Summary OAuth登录
// @Description 通过OAuth提供商登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body service.OAuthLoginRequest true "OAuth登录信息"
// @Success 200 {object} conf.Response
// @Router /api/v1/user/oauth/login [post]
func oauthLogin(c *gin.Context) {
	var req service.OAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// TODO: 更新OAuthLogin方法返回TokenPair
	profile, token, err := svc.OAuthLogin(c, &req, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, conf.Response{
			Code: http.StatusUnauthorized,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "登录成功",
		Data: gin.H{
			"user":  profile,
			"token": token,
		},
	})
}

// 发送验证码
// @Summary 发送验证码
// @Description 发送邮箱验证码
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body service.SendCodeRequest true "发送验证码请求"
// @Success 200 {object} conf.Response
// @Router /api/v1/user/send-code [post]
func sendCode(c *gin.Context) {
	var req service.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	err := svc.SendVerificationCode(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  "发送验证码失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "验证码已发送",
	})
}

// 重置密码
// @Summary 重置密码
// @Description 通过验证码重置密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body service.ResetPasswordRequest true "重置密码请求"
// @Success 200 {object} conf.Response
// @Router /api/v1/user/reset-password [post]
func resetPassword(c *gin.Context) {
	var req service.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	err := svc.ResetPassword(c, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "密码重置成功",
	})
}

// 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} conf.Response
// @Router /api/v1/user/profile [get]
func getUserProfile(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, conf.Response{
			Code: http.StatusUnauthorized,
			Msg:  "未授权",
		})
		return
	}

	profile, err := svc.GetUserProfile(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  "获取用户信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: profile,
	})
}

// 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username formData string false "用户名"
// @Param avatar formData string false "头像"
// @Success 200 {object} conf.Response
// @Router /api/v1/user/profile [put]
func updateUserProfile(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, conf.Response{
			Code: http.StatusUnauthorized,
			Msg:  "未授权",
		})
		return
	}

	username := c.PostForm("username")
	avatar := c.PostForm("avatar")

	profile, err := svc.UpdateUserProfile(c, userID, username, avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  "更新用户信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "更新成功",
		Data: profile,
	})
}

// 注销登录
// @Summary 注销登录
// @Description 注销当前用户的登录状态
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} conf.Response
// @Router /api/v1/user/logout [post]
func logout(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, conf.Response{
			Code: http.StatusUnauthorized,
			Msg:  "未授权",
		})
		return
	}

	err := svc.Logout(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  "注销失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "注销成功",
	})
}

// @Summary 检查邮箱是否已存在
// @Description 检查提供的邮箱是否已经注册
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param email query string true "邮箱地址"
// @Success 200 {object} conf.Response
// @Router /api/v1/user/check-email [get]
func checkEmail(c *gin.Context) {
	// 获取邮箱参数
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "邮箱不能为空",
		})
		return
	}

	// 调用service层检查邮箱
	exists, err := svc.CheckEmailExists(c, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("检查邮箱失败: %v", err),
		})
		return
	}

	// 返回检查结果
	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "查询成功",
		Data: gin.H{
			"exists": exists,
		},
	})
}

// @Summary 获取登录日志
// @Description 获取用户登录日志，支持分页和筛选
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user_id query int false "用户ID"
// @Param ip query string false "IP地址"
// @Param status query int false "状态(0:失败,1:成功)"
// @Param start_time query string false "开始时间(格式:2006-01-02T15:04:05Z)"
// @Param end_time query string false "结束时间(格式:2006-01-02T15:04:05Z)"
// @Param page query int true "页码"
// @Param page_size query int true "每页大小"
// @Success 200 {object} conf.Response
// @Router /api/v1/user/login-logs [get]
func getLoginLogs(c *gin.Context) {
	// 解析请求参数
	var req service.LoginLogRequest
	
	// 获取分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	req.Page = page
	
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	req.PageSize = pageSize
	
	// 获取筛选参数
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err == nil && userID > 0 {
			req.UserID = userID
		}
	}
	
	req.IP = c.Query("ip")
	
	if statusStr := c.Query("status"); statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err == nil && (status == 0 || status == 1) {
			req.Status = &status
		}
	}
	
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err == nil {
			req.StartTime = startTime
		}
	}
	
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err == nil {
			req.EndTime = endTime
		}
	}
	
	// 调用service层获取登录日志
	resp, err := svc.GetLoginLogs(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("获取登录日志失败: %v", err),
		})
		return
	}
	
	// 返回查询结果
	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "查询成功",
		Data: resp,
	})
} 