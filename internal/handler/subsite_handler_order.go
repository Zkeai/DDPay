package handler

import (
	"net/http"
	"strconv"

	"github.com/Zkeai/DDPay/common/conf"
	"github.com/gin-gonic/gin"
)

// CreateSubsiteOrder 创建分站订单
func (h *SubsiteHandler) CreateSubsiteOrder(c *gin.Context) {
	var req CreateSubsiteOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "参数错误", Data: err.Error()})
		return
	}
	
	// 如果用户ID为0，则使用当前用户ID
	if req.UserID == 0 {
		userID, exists := c.Get("user_id")
		if exists {
			req.UserID = userID.(int64)
		}
	}
	
	order, err := h.subsiteService.CreateSubsiteOrder(c, req.SubsiteID, req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "创建订单失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "订单创建成功", Data: gin.H{"order": order}})
}

// GetSubsiteOrder 获取分站订单
func (h *SubsiteHandler) GetSubsiteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "无效的订单ID", Data: err.Error()})
		return
	}
	
	order, err := h.subsiteService.GetSubsiteOrder(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取订单失败", Data: err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "订单不存在", Data: nil})
		return
	}
	
	// 检查权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{Code: 401, Msg: "未登录", Data: nil})
		return
	}
	ownerID := userID.(int64)
	
	// 获取分站信息
	subsite, err := h.subsiteService.GetSubsiteByID(c, order.SubsiteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取分站失败", Data: err.Error()})
		return
	}
	if subsite == nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "分站不存在", Data: nil})
		return
	}
	
	// 检查是否为分站所有者、管理员或订单所属用户
	userRole, _ := c.Get("user_role")
	role := userRole.(string)
	if subsite.OwnerID != ownerID && role != "admin" && order.UserID != ownerID {
		c.JSON(http.StatusForbidden, conf.Response{Code: 403, Msg: "无权查看该订单", Data: nil})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: gin.H{"order": order}})
}

// UpdateSubsiteOrderStatus 更新分站订单状态
func (h *SubsiteHandler) UpdateSubsiteOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "无效的订单ID", Data: err.Error()})
		return
	}
	
	var req UpdateSubsiteOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "参数错误", Data: err.Error()})
		return
	}
	
	// 获取订单信息
	order, err := h.subsiteService.GetSubsiteOrder(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取订单失败", Data: err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "订单不存在", Data: nil})
		return
	}
	
	// 检查权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{Code: 401, Msg: "未登录", Data: nil})
		return
	}
	ownerID := userID.(int64)
	
	// 获取分站信息
	subsite, err := h.subsiteService.GetSubsiteByID(c, order.SubsiteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取分站失败", Data: err.Error()})
		return
	}
	if subsite == nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "分站不存在", Data: nil})
		return
	}
	
	// 检查是否为分站所有者或管理员
	userRole, _ := c.Get("user_role")
	role := userRole.(string)
	if subsite.OwnerID != ownerID && role != "admin" {
		c.JSON(http.StatusForbidden, conf.Response{Code: 403, Msg: "无权操作", Data: nil})
		return
	}
	
	// 更新订单状态
	err = h.subsiteService.UpdateSubsiteOrderStatus(c, id, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "更新订单状态失败", Data: err.Error()})
		return
	}
	
	// 获取更新后的订单
	updatedOrder, err := h.subsiteService.GetSubsiteOrder(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取更新后的订单失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "订单状态更新成功", Data: gin.H{"order": updatedOrder}})
}

// ListSubsiteOrders 获取分站订单列表
func (h *SubsiteHandler) ListSubsiteOrders(c *gin.Context) {
	subsiteIDStr := c.Query("subsite_id")
	if subsiteIDStr == "" {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "缺少分站ID", Data: nil})
		return
	}
	
	subsiteID, err := strconv.ParseInt(subsiteIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "无效的分站ID", Data: err.Error()})
		return
	}
	
	// 检查权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{Code: 401, Msg: "未登录", Data: nil})
		return
	}
	ownerID := userID.(int64)
	
	// 获取分站信息
	subsite, err := h.subsiteService.GetSubsiteByID(c, subsiteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取分站失败", Data: err.Error()})
		return
	}
	if subsite == nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "分站不存在", Data: nil})
		return
	}
	
	// 检查是否为分站所有者或管理员
	userRole, _ := c.Get("user_role")
	role := userRole.(string)
	if subsite.OwnerID != ownerID && role != "admin" {
		c.JSON(http.StatusForbidden, conf.Response{Code: 403, Msg: "无权操作", Data: nil})
		return
	}
	
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	statusStr := c.DefaultQuery("status", "-1")
	
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	status, _ := strconv.Atoi(statusStr)
	
	orders, total, err := h.subsiteService.ListSubsiteOrders(c, subsiteID, page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取订单列表失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{
		Code: 200,
		Msg: "success",
		Data: gin.H{
			"orders":    orders,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
} 