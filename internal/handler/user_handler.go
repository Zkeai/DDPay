package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateUserLevelReq 更新用户等级请求
type UpdateUserLevelReq struct {
	UserID int64 `json:"user_id" binding:"required"`
	Level  int   `json:"level" binding:"required,min=1,max=3"`
}

// UpdateUserStatusReq 更新用户状态请求
type UpdateUserStatusReq struct {
	UserID int64 `json:"user_id" binding:"required"`
	Status int   `json:"status" binding:"required,oneof=0 1"` // 0-禁用 1-启用
}

// HandleUpdateUserLevel 更新用户等级
func (h *Handler) HandleUpdateUserLevel(c *gin.Context) {
	var req UpdateUserLevelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	// 打印请求参数
	fmt.Printf("收到更新用户等级请求: 用户ID=%d, 新等级=%d\n", req.UserID, req.Level)

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

	// 获取目标用户
	targetUser, err := h.userService.GetUserByID(c, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取用户信息失败"})
		return
	}

	if targetUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
		return
	}

	// 打印更新前的用户等级
	fmt.Printf("更新前用户等级: 用户ID=%d, 当前等级=%d\n", targetUser.ID, targetUser.Level)

	// 更新用户等级
	targetUser.Level = req.Level
	err = h.userService.UpdateUser(c, targetUser)
	if err != nil {
		fmt.Printf("更新用户等级失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新用户等级失败: " + err.Error()})
		return
	}

	// 重新获取用户信息，验证更新是否成功
	updatedUser, err := h.userService.GetUserByID(c, req.UserID)
	if err != nil {
		fmt.Printf("获取更新后用户信息失败: %v\n", err)
	} else {
		fmt.Printf("更新后用户等级: 用户ID=%d, 新等级=%d\n", updatedUser.ID, updatedUser.Level)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "用户等级更新成功",
		"data": gin.H{
			"user_id": targetUser.ID,
			"level":   targetUser.Level,
		},
	})
}

// HandleUpdateUserStatus 更新用户状态
func (h *Handler) HandleUpdateUserStatus(c *gin.Context) {
	var req UpdateUserStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
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

	// 获取目标用户
	targetUser, err := h.userService.GetUserByID(c, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取用户信息失败"})
		return
	}

	if targetUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
		return
	}

	// 不允许禁用管理员账号
	if targetUser.Role == "admin" && req.Status == 0 {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "不允许禁用管理员账号"})
		return
	}

	// 更新用户状态
	targetUser.Status = req.Status
	err = h.userService.UpdateUser(c, targetUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新用户状态失败: " + err.Error()})
		return
	}

	statusText := "启用"
	if req.Status == 0 {
		statusText = "禁用"
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "用户已" + statusText,
		"data": gin.H{
			"user_id": targetUser.ID,
			"status":  targetUser.Status,
		},
	})
} 