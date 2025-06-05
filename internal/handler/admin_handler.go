package handler

import (
	"net/http"

	"github.com/Zkeai/DDPay/internal/model"
	"github.com/gin-gonic/gin"
)

// ListUsersReq 获取用户列表请求
type ListUsersReq struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

// ListUsersResp 获取用户列表响应
type ListUsersResp struct {
	Users      []*model.UserProfile `json:"users"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// HandleListUsers 获取用户列表
func (h *Handler) HandleListUsers(c *gin.Context) {
	var req ListUsersReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	// 检查当前用户是否为管理员
	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
		return
	}

	if currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "权限不足"})
		return
	}

	// 获取用户列表
	users, total, err := h.userService.ListUsers(c, req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取用户列表失败: " + err.Error()})
		return
	}

	// 计算总页数
	totalPages := total / req.PageSize
	if total%req.PageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取用户列表成功",
		"data": ListUsersResp{
			Users:      users,
			Total:      total,
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalPages: totalPages,
		},
	})
} 