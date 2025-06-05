package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Zkeai/DDPay/internal/model"
	"github.com/gin-gonic/gin"
)

// 会员等级相关请求结构体
type (
	// PurchaseMembershipReq 购买会员请求
	PurchaseMembershipReq struct {
		LevelID       int64  `json:"level_id" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	// UpgradeMembershipReq 升级会员请求
	UpgradeMembershipReq struct {
		LevelID       int64  `json:"level_id" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	// RenewMembershipReq 续费会员请求
	RenewMembershipReq struct {
		DurationDays  int    `json:"duration_days" binding:"required,min=1"`
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	// CreateMembershipLevelReq 创建会员等级请求
	CreateMembershipLevelReq struct {
		Name               string    `json:"name" binding:"required"`
		Level              int       `json:"level" binding:"required,min=1"`
		Icon               string    `json:"icon"`
		Price              float64   `json:"price" binding:"required,min=0"`
		Description        string    `json:"description"`
		DiscountRate       float64   `json:"discount_rate" binding:"required,min=0,max=1"`
		MaxSubsites        int       `json:"max_subsites" binding:"required"`
		CustomServiceAccess bool      `json:"custom_service_access"`
		VIPGroupAccess     bool      `json:"vip_group_access"`
		Priority           int       `json:"priority" binding:"required,min=1"`
		Benefits           []BenefitInfo `json:"benefits"`
	}

	// UpdateMembershipLevelReq 更新会员等级请求
	UpdateMembershipLevelReq struct {
		ID                 int64     `json:"id" binding:"required"`
		Name               string    `json:"name" binding:"required"`
		Level              int       `json:"level" binding:"required,min=1"`
		Icon               string    `json:"icon"`
		Price              float64   `json:"price" binding:"required,min=0"`
		Description        string    `json:"description"`
		DiscountRate       float64   `json:"discount_rate" binding:"required,min=0,max=1"`
		MaxSubsites        int       `json:"max_subsites" binding:"required"`
		CustomServiceAccess bool      `json:"custom_service_access"`
		VIPGroupAccess     bool      `json:"vip_group_access"`
		Priority           int       `json:"priority" binding:"required,min=1"`
		Benefits           []BenefitInfo `json:"benefits"`
	}

	// BenefitInfo 权益信息
	BenefitInfo struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}
)

// HandleGetMembershipLevels 获取会员等级列表
func (h *Handler) HandleGetMembershipLevels(c *gin.Context) {
	levels, err := h.membershipService.GetMembershipLevels(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取会员等级失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取会员等级成功",
		"data": gin.H{
			"levels": levels,
		},
	})
}

// HandleGetUserMembership 获取用户会员信息
func (h *Handler) HandleGetUserMembership(c *gin.Context) {
	// 获取当前用户
	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	// 获取用户会员信息
	membership, level, err := h.membershipService.GetUserMembershipWithLevel(c, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取会员信息失败: " + err.Error(),
		})
		return
	}

	// 检查会员是否过期
	var isExpired bool
	var remainingDays int
	
	if membership == nil {
		// 没有会员记录，使用默认的免费会员
		isExpired = false
		remainingDays = -1 // 永久有效
	} else if membership.EndDate == nil {
		// 没有结束日期，表示永久会员
		isExpired = false
		remainingDays = -1
	} else {
		now := time.Now()
		isExpired = now.After(*membership.EndDate)
		
		if isExpired {
			remainingDays = 0
		} else {
			// 计算剩余天数
			remainingDays = int(membership.EndDate.Sub(now).Hours() / 24)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取会员信息成功",
		"data": gin.H{
			"membership":     membership,
			"level":          level,
			"is_expired":     isExpired,
			"remaining_days": remainingDays,
		},
	})
}

// HandlePurchaseMembership 购买会员
func (h *Handler) HandlePurchaseMembership(c *gin.Context) {
	var req PurchaseMembershipReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户
	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	// 获取会员等级
	level, err := h.membershipService.GetMembershipLevelByID(c, req.LevelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "会员等级不存在",
		})
		return
	}

	// 购买会员
	orderID, err := h.membershipService.PurchaseMembership(c, currentUser.ID, req.LevelID, req.PaymentMethod, level.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "购买会员失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "会员购买请求已提交",
		"data": gin.H{
			"order_id": orderID,
			"level":    level,
			"amount":   level.Price,
		},
	})
}

// HandleUpgradeMembership 升级会员
func (h *Handler) HandleUpgradeMembership(c *gin.Context) {
	var req UpgradeMembershipReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户
	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	// 获取目标会员等级
	level, err := h.membershipService.GetMembershipLevelByID(c, req.LevelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "会员等级不存在",
		})
		return
	}

	// 升级会员
	orderID, err := h.membershipService.UpgradeMembership(c, currentUser.ID, req.LevelID, req.PaymentMethod, level.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "升级会员失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "会员升级请求已提交",
		"data": gin.H{
			"order_id": orderID,
			"level":    level,
			"amount":   level.Price,
		},
	})
}

// HandleRenewMembership 续费会员
func (h *Handler) HandleRenewMembership(c *gin.Context) {
	var req RenewMembershipReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户
	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	// 计算续费金额（实际项目中应该根据业务逻辑计算）
	// 这里简单示例，按天数比例计算
	// 获取用户当前会员信息
	membership, level, err := h.membershipService.GetUserMembershipWithLevel(c, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取会员信息失败: " + err.Error(),
		})
		return
	}

	if membership == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "您还不是会员，请先购买会员",
		})
		return
	}

	// 计算续费金额，按天计算（实际项目中可能有不同策略）
	baseAmount := level.Price / 30 // 假设30天为一个月
	amount := baseAmount * float64(req.DurationDays)

	// 续费会员
	orderID, err := h.membershipService.RenewMembership(c, currentUser.ID, req.DurationDays, req.PaymentMethod, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "续费会员失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "会员续费请求已提交",
		"data": gin.H{
			"order_id":       orderID,
			"level":          level,
			"amount":         amount,
			"duration_days":  req.DurationDays,
		},
	})
}

// HandleCheckMembershipRequirements 检查会员升级条件
func (h *Handler) HandleCheckMembershipRequirements(c *gin.Context) {
	levelID, err := strconv.ParseInt(c.Query("level_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的会员等级ID",
		})
		return
	}

	// 获取当前用户
	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	// 获取会员等级
	level, err := h.membershipService.GetMembershipLevelByID(c, levelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "会员等级不存在",
		})
		return
	}

	// 检查升级条件
	allMet, results, err := h.membershipService.CheckMembershipRequirements(c, currentUser.ID, levelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "检查升级条件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "检查升级条件成功",
		"data": gin.H{
			"level":       level,
			"all_met":     allMet,
			"requirements": results,
		},
	})
}

// HandleGetMembershipTransactions 获取会员交易记录
func (h *Handler) HandleGetMembershipTransactions(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 获取当前用户
	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	// 获取交易记录
	transactions, total, err := h.membershipService.GetMembershipTransactions(c, currentUser.ID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取交易记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取交易记录成功",
		"data": gin.H{
			"transactions": transactions,
			"total":        total,
			"page":         page,
			"page_size":    pageSize,
			"total_pages":  (total + pageSize - 1) / pageSize,
		},
	})
}

// HandleCreateMembershipLevel 创建会员等级
func (h *Handler) HandleCreateMembershipLevel(c *gin.Context) {
	var req CreateMembershipLevelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查同级别的会员是否已存在
	existingLevel, err := h.membershipService.GetMembershipLevelByLevel(c, req.Level)
	if err == nil && existingLevel != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "相同级别的会员已存在",
		})
		return
	}

	// 创建会员等级
	level := &model.MembershipLevel{
		Name:               req.Name,
		Level:              req.Level,
		Icon:               req.Icon,
		Price:              req.Price,
		Description:        req.Description,
		DiscountRate:       req.DiscountRate,
		MaxSubsites:        req.MaxSubsites,
		CustomServiceAccess: req.CustomServiceAccess,
		VIPGroupAccess:     req.VIPGroupAccess,
		Priority:           req.Priority,
	}

	levelID, err := h.membershipService.CreateMembershipLevel(c, level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建会员等级失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建会员等级成功",
		"data": gin.H{
			"id":    levelID,
			"level": level,
		},
	})
}

// HandleUpdateMembershipLevel 更新会员等级
func (h *Handler) HandleUpdateMembershipLevel(c *gin.Context) {
	var req UpdateMembershipLevelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查会员等级是否存在
	existingLevel, err := h.membershipService.GetMembershipLevelByID(c, req.ID)
	if err != nil || existingLevel == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "会员等级不存在",
		})
		return
	}

	// 检查同级别的会员是否已存在（排除自身）
	if req.Level != existingLevel.Level {
		otherLevel, err := h.membershipService.GetMembershipLevelByLevel(c, req.Level)
		if err == nil && otherLevel != nil && otherLevel.ID != req.ID {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "相同级别的会员已存在",
			})
			return
		}
	}

	// 更新会员等级
	level := &model.MembershipLevel{
		ID:                 req.ID,
		Name:               req.Name,
		Level:              req.Level,
		Icon:               req.Icon,
		Price:              req.Price,
		Description:        req.Description,
		DiscountRate:       req.DiscountRate,
		MaxSubsites:        req.MaxSubsites,
		CustomServiceAccess: req.CustomServiceAccess,
		VIPGroupAccess:     req.VIPGroupAccess,
		Priority:           req.Priority,
	}

	err = h.membershipService.UpdateMembershipLevel(c, level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新会员等级失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新会员等级成功",
		"data": gin.H{
			"level": level,
		},
	})
}

// HandleDeleteMembershipLevel 删除会员等级
func (h *Handler) HandleDeleteMembershipLevel(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "会员等级ID不能为空",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的会员等级ID",
		})
		return
	}

	// 检查会员等级是否存在
	existingLevel, err := h.membershipService.GetMembershipLevelByID(c, id)
	if err != nil || existingLevel == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "会员等级不存在",
		})
		return
	}

	// 删除会员等级
	err = h.membershipService.DeleteMembershipLevel(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除会员等级失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除会员等级成功",
	})
} 