package handler

import (
	"net/http"
	"strconv"

	"github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/internal/model"
	"github.com/gin-gonic/gin"
)

// GetSubsiteBalance 获取分站余额
func (h *SubsiteHandler) GetSubsiteBalance(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{Code: 401, Msg: "未登录", Data: nil})
		return
	}
	ownerID := userID.(int64)
	
	// 获取余额
	balance, err := h.subsiteService.GetSubsiteBalance(c, ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取余额失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: gin.H{"balance": balance}})
}

// ListSubsiteBalanceLogs 获取分站余额变动记录
func (h *SubsiteHandler) ListSubsiteBalanceLogs(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{Code: 401, Msg: "未登录", Data: nil})
		return
	}
	ownerID := userID.(int64)
	
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	
	logs, total, err := h.subsiteService.ListSubsiteBalanceLogs(c, ownerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取余额变动记录失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{
		Code: 200,
		Msg: "success",
		Data: gin.H{
			"logs":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateSubsiteWithdrawal 创建分站提现申请
func (h *SubsiteHandler) CreateSubsiteWithdrawal(c *gin.Context) {
	var req CreateSubsiteWithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "参数错误", Data: err.Error()})
		return
	}
	
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{Code: 401, Msg: "未登录", Data: nil})
		return
	}
	ownerID := userID.(int64)
	
	// 检查提现申请是否为当前用户
	if req.OwnerID != ownerID {
		c.JSON(http.StatusForbidden, conf.Response{Code: 403, Msg: "只能为自己申请提现", Data: nil})
		return
	}
	
	// 创建提现申请
	withdrawal := &model.SubsiteWithdrawal{
		OwnerID:     req.OwnerID,
		Amount:      req.Amount,
		AccountType: req.AccountType,
		AccountName: req.AccountName,
		AccountNo:   req.AccountNo,
		Remark:      req.Remark,
	}
	
	id, err := h.subsiteService.CreateSubsiteWithdrawal(c, withdrawal)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "创建提现申请失败", Data: err.Error()})
		return
	}
	
	// 获取创建后的提现申请
	createdWithdrawal, err := h.subsiteService.GetSubsiteWithdrawal(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取提现申请失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "提现申请创建成功", Data: gin.H{"withdrawal": createdWithdrawal}})
}

// GetSubsiteWithdrawal 获取分站提现申请
func (h *SubsiteHandler) GetSubsiteWithdrawal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "无效的提现申请ID", Data: err.Error()})
		return
	}
	
	// 获取提现申请
	withdrawal, err := h.subsiteService.GetSubsiteWithdrawal(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取提现申请失败", Data: err.Error()})
		return
	}
	if withdrawal == nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "提现申请不存在", Data: nil})
		return
	}
	
	// 检查权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{Code: 401, Msg: "未登录", Data: nil})
		return
	}
	ownerID := userID.(int64)
	
	// 检查是否为提现申请的所有者或管理员
	userRole, _ := c.Get("user_role")
	role := userRole.(string)
	if withdrawal.OwnerID != ownerID && role != "admin" {
		c.JSON(http.StatusForbidden, conf.Response{Code: 403, Msg: "无权查看该提现申请", Data: nil})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: gin.H{"withdrawal": withdrawal}})
}

// ProcessSubsiteWithdrawal 处理分站提现申请
func (h *SubsiteHandler) ProcessSubsiteWithdrawal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "无效的提现申请ID", Data: err.Error()})
		return
	}
	
	var req ProcessSubsiteWithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "参数错误", Data: err.Error()})
		return
	}
	
	// 检查是否为管理员
	userRole, exists := c.Get("user_role")
	if !exists || userRole.(string) != "admin" {
		c.JSON(http.StatusForbidden, conf.Response{Code: 403, Msg: "需要管理员权限", Data: nil})
		return
	}
	
	// 处理提现申请
	err = h.subsiteService.ProcessSubsiteWithdrawal(c, id, req.Status, req.AdminRemark)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "处理提现申请失败", Data: err.Error()})
		return
	}
	
	// 获取处理后的提现申请
	withdrawal, err := h.subsiteService.GetSubsiteWithdrawal(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取提现申请失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "提现申请处理成功", Data: gin.H{"withdrawal": withdrawal}})
}

// ListSubsiteWithdrawals 获取分站提现申请列表
func (h *SubsiteHandler) ListSubsiteWithdrawals(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{Code: 401, Msg: "未登录", Data: nil})
		return
	}
	ownerID := userID.(int64)
	
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	statusStr := c.DefaultQuery("status", "-1")
	
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	status, _ := strconv.Atoi(statusStr)
	
	withdrawals, total, err := h.subsiteService.ListSubsiteWithdrawals(c, ownerID, page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取提现申请列表失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{
		Code: 200,
		Msg: "success",
		Data: gin.H{
			"withdrawals": withdrawals,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
		},
	})
} 