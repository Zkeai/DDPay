package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/internal/model"
	"github.com/gin-gonic/gin"
)

// CreateSubsiteProduct 创建分站商品
func (h *SubsiteHandler) CreateSubsiteProduct(c *gin.Context) {
	var req CreateSubsiteProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "参数错误", Data: err.Error()})
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
	subsite, err := h.subsiteService.GetSubsiteByID(c, req.SubsiteID)
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
	
	// 创建商品对象
	product := &model.SubsiteProduct{
		SubsiteID:     req.SubsiteID,
		MainProductID: req.MainProductID,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         req.Stock,
		Image:         req.Image,
		Status:        req.Status,
		IsTimeLimited: req.IsTimeLimited,
		SortOrder:     req.SortOrder,
	}
	
	// 处理时间
	if req.IsTimeLimited == 1 {
		if req.StartTime == "" || req.EndTime == "" {
			c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "限时商品必须设置开始和结束时间", Data: nil})
			return
		}
		
		startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "开始时间格式错误，应为 YYYY-MM-DD HH:MM:SS", Data: err.Error()})
			return
		}
		
		endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "结束时间格式错误，应为 YYYY-MM-DD HH:MM:SS", Data: err.Error()})
			return
		}
		
		if endTime.Before(startTime) {
			c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "结束时间不能早于开始时间", Data: nil})
			return
		}
		
		product.StartTime = startTime
		product.EndTime = endTime
	}
	
	id, err := h.subsiteService.CreateSubsiteProduct(c, product)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "创建商品失败", Data: err.Error()})
		return
	}
	
	product.ID = id
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "商品创建成功", Data: gin.H{"product": product}})
}

// GetSubsiteProduct 获取分站商品
func (h *SubsiteHandler) GetSubsiteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "无效的商品ID", Data: err.Error()})
		return
	}
	
	product, err := h.subsiteService.GetSubsiteProduct(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取商品失败", Data: err.Error()})
		return
	}
	if product == nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "商品不存在", Data: nil})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: gin.H{"product": product}})
}

// UpdateSubsiteProduct 更新分站商品
func (h *SubsiteHandler) UpdateSubsiteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "无效的商品ID", Data: err.Error()})
		return
	}
	
	var req UpdateSubsiteProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "参数错误", Data: err.Error()})
		return
	}
	
	// 检查ID是否一致
	if id != req.ID {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "URL中的ID与请求体中的ID不一致", Data: nil})
		return
	}
	
	// 获取商品信息
	product, err := h.subsiteService.GetSubsiteProduct(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取商品失败", Data: err.Error()})
		return
	}
	if product == nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "商品不存在", Data: nil})
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
	subsite, err := h.subsiteService.GetSubsiteByID(c, product.SubsiteID)
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
	
	// 更新商品信息
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.OriginalPrice = req.OriginalPrice
	product.Stock = req.Stock
	product.Image = req.Image
	product.Status = req.Status
	product.IsTimeLimited = req.IsTimeLimited
	product.SortOrder = req.SortOrder
	
	// 处理时间
	if req.IsTimeLimited == 1 {
		if req.StartTime == "" || req.EndTime == "" {
			c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "限时商品必须设置开始和结束时间", Data: nil})
			return
		}
		
		startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "开始时间格式错误，应为 YYYY-MM-DD HH:MM:SS", Data: err.Error()})
			return
		}
		
		endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "结束时间格式错误，应为 YYYY-MM-DD HH:MM:SS", Data: err.Error()})
			return
		}
		
		if endTime.Before(startTime) {
			c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "结束时间不能早于开始时间", Data: nil})
			return
		}
		
		product.StartTime = startTime
		product.EndTime = endTime
	} else {
		// 清空时间
		product.StartTime = time.Time{}
		product.EndTime = time.Time{}
	}
	
	err = h.subsiteService.UpdateSubsiteProduct(c, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, conf.Response{Code: 500, Msg: "更新商品失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "商品更新成功", Data: gin.H{"product": product}})
}

// ListSubsiteProducts 获取分站商品列表
func (h *SubsiteHandler) ListSubsiteProducts(c *gin.Context) {
	subsiteIDStr := c.Query("subsite_id")
	if subsiteIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少分站ID"})
		return
	}
	
	subsiteID, err := strconv.ParseInt(subsiteIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分站ID"})
		return
	}
	
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	statusStr := c.DefaultQuery("status", "-1")
	
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	status, _ := strconv.Atoi(statusStr)
	
	products, total, err := h.subsiteService.ListSubsiteProducts(c, subsiteID, page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"products":  products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
} 