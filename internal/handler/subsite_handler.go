package handler

import (
	"net/http"
	"strconv"

	"github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/internal/model"
	"github.com/Zkeai/DDPay/internal/service"
	"github.com/gin-gonic/gin"
)

// SubsiteHandler 分站处理程序
type SubsiteHandler struct {
	subsiteService service.SubsiteService
}

// NewSubsiteHandler 创建分站处理程序
func NewSubsiteHandler(subsiteService service.SubsiteService) *SubsiteHandler {
	return &SubsiteHandler{
		subsiteService: subsiteService,
	}
}

// RegisterRoutes 注册路由
func (h *SubsiteHandler) RegisterRoutes(router *gin.RouterGroup) {
	subsiteRouter := router.Group("/subsite")
	{
		// 分站管理
		subsiteRouter.POST("/create", h.CreateSubsite)
		subsiteRouter.GET("/info", h.GetSubsiteInfo)
		subsiteRouter.PUT("/update", h.UpdateSubsite)
		subsiteRouter.GET("/list", h.ListSubsites)
		
		// 分站JSON配置
		subsiteRouter.POST("/config", h.SaveSubsiteJsonConfig)
		subsiteRouter.GET("/config", h.GetSubsiteJsonConfig)
		
		// 分站商品
		subsiteRouter.POST("/product", h.CreateSubsiteProduct)
		subsiteRouter.GET("/product/:id", h.GetSubsiteProduct)
		subsiteRouter.PUT("/product/:id", h.UpdateSubsiteProduct)
		subsiteRouter.GET("/products", h.ListSubsiteProducts)
		
		// 分站订单
		subsiteRouter.POST("/order", h.CreateSubsiteOrder)
		subsiteRouter.GET("/order/:id", h.GetSubsiteOrder)
		subsiteRouter.PUT("/order/:id/status", h.UpdateSubsiteOrderStatus)
		subsiteRouter.GET("/orders", h.ListSubsiteOrders)
		
		// 分站余额
		subsiteRouter.GET("/balance", h.GetSubsiteBalance)
		subsiteRouter.GET("/balance/logs", h.ListSubsiteBalanceLogs)
		
		// 分站提现
		subsiteRouter.POST("/withdrawal", h.CreateSubsiteWithdrawal)
		subsiteRouter.GET("/withdrawal/:id", h.GetSubsiteWithdrawal)
		subsiteRouter.PUT("/withdrawal/:id", h.ProcessSubsiteWithdrawal)
		subsiteRouter.GET("/withdrawals", h.ListSubsiteWithdrawals)
	}
}

// 请求和响应结构体
type (
	// 创建分站请求
	CreateSubsiteRequest struct {
		Name           string  `json:"name" binding:"required"`
		Subdomain      string  `json:"subdomain" binding:"required"`
		Description    string  `json:"description"`
		Logo           string  `json:"logo"`
		Domain         string  `json:"domain"`
		Theme          string  `json:"theme"`
		Status         int     `json:"status"`
		CommissionRate float64 `json:"commission_rate" binding:"required,gte=0,lte=100"`
	}
	
	// 更新分站请求
	UpdateSubsiteRequest struct {
		ID             int64   `json:"id" binding:"required"`
		Name           string  `json:"name" binding:"required"`
		Domain         string  `json:"domain"`
		Subdomain      string  `json:"subdomain" binding:"required"`
		Logo           string  `json:"logo"`
		Description    string  `json:"description"`
		Theme          string  `json:"theme"`
		Status         int     `json:"status" binding:"oneof=0 1"`
		CommissionRate float64 `json:"commission_rate" binding:"required,gte=0,lte=100"`
	}
	
	// 分站JSON配置请求
	SubsiteJsonConfigRequest struct {
		SubsiteID int64                  `json:"subsite_id" binding:"required"`
		Config    map[string]interface{} `json:"config" binding:"required"`
	}
	
	// 创建分站商品请求
	CreateSubsiteProductRequest struct {
		SubsiteID      int64   `json:"subsite_id" binding:"required"`
		MainProductID  int64   `json:"main_product_id"`
		Name           string  `json:"name" binding:"required"`
		Description    string  `json:"description"`
		Price          float64 `json:"price" binding:"required,gt=0"`
		OriginalPrice  float64 `json:"original_price"`
		Stock          int     `json:"stock" binding:"gte=0"`
		Image          string  `json:"image"`
		Status         int     `json:"status" binding:"required,oneof=0 1"`
		IsTimeLimited  int     `json:"is_time_limited" binding:"oneof=0 1"`
		StartTime      string  `json:"start_time"`
		EndTime        string  `json:"end_time"`
		SortOrder      int     `json:"sort_order"`
	}
	
	// 更新分站商品请求
	UpdateSubsiteProductRequest struct {
		ID             int64   `json:"id" binding:"required"`
		Name           string  `json:"name" binding:"required"`
		Description    string  `json:"description"`
		Price          float64 `json:"price" binding:"required,gt=0"`
		OriginalPrice  float64 `json:"original_price"`
		Stock          int     `json:"stock" binding:"gte=0"`
		Image          string  `json:"image"`
		Status         int     `json:"status" binding:"required,oneof=0 1"`
		IsTimeLimited  int     `json:"is_time_limited" binding:"oneof=0 1"`
		StartTime      string  `json:"start_time"`
		EndTime        string  `json:"end_time"`
		SortOrder      int     `json:"sort_order"`
	}
	
	// 创建分站订单请求
	CreateSubsiteOrderRequest struct {
		SubsiteID int64 `json:"subsite_id" binding:"required"`
		UserID    int64 `json:"user_id"`
		ProductID int64 `json:"product_id" binding:"required"`
		Quantity  int   `json:"quantity" binding:"required,gt=0"`
	}
	
	// 更新分站订单状态请求
	UpdateSubsiteOrderStatusRequest struct {
		Status int `json:"status" binding:"required,oneof=0 1 2 3"`
	}
	
	// 创建分站提现申请请求
	CreateSubsiteWithdrawalRequest struct {
		OwnerID     int64   `json:"owner_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		AccountType string  `json:"account_type" binding:"required"`
		AccountName string  `json:"account_name" binding:"required"`
		AccountNo   string  `json:"account_no" binding:"required"`
		Remark      string  `json:"remark"`
	}
	
	// 处理分站提现申请请求
	ProcessSubsiteWithdrawalRequest struct {
		Status      int    `json:"status" binding:"required,oneof=1 2"`
		AdminRemark string `json:"admin_remark"`
	}
)

// CreateSubsite 创建分站
func (h *SubsiteHandler) CreateSubsite(c *gin.Context) {
	var req CreateSubsiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "参数错误", Data: err.Error()})
		return
	}
	
	// 从当前用户获取OwnerID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{Code: 401, Msg: "未登录", Data: nil})
		return
	}
	ownerID := userID.(int64)
	
	// 创建完整的分站对象
	subsite := &model.Subsite{
		OwnerID:        ownerID,
		Name:           req.Name,
		Domain:         req.Domain,
		Subdomain:      req.Subdomain,
		Logo:           req.Logo,
		Description:    req.Description,
		Theme:          req.Theme,
		Status:         req.Status,
		CommissionRate: req.CommissionRate,
	}
	
	// 如果主题为空，设置默认主题
	if subsite.Theme == "" {
		subsite.Theme = "default"
	}
	
	// 如果状态为0，设置为默认启用状态
	if subsite.Status == 0 {
		subsite.Status = 1 // 默认启用
	}
	
	// 创建分站
	createdSubsite, err := h.subsiteService.CreateSubsiteObject(c, subsite)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "创建分站失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: gin.H{"subsite": createdSubsite}})
}

// GetSubsiteInfo 获取分站信息
func (h *SubsiteHandler) GetSubsiteInfo(c *gin.Context) {
	subsiteIDStr := c.Query("id")
	if subsiteIDStr == "" {
		// 如果没有提供ID，尝试获取当前用户的分站
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, conf.Response{Code: 401, Msg: "未登录", Data: nil})
			return
		}
		ownerID := userID.(int64)
		
		subsite, err := h.subsiteService.GetSubsiteByOwnerID(c, ownerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取分站失败", Data: err.Error()})
			return
		}
		if subsite == nil {
			c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "未找到分站", Data: nil})
			return
		}
		subsiteIDStr = strconv.FormatInt(subsite.ID, 10)
	}
	
	subsiteID, err := strconv.ParseInt(subsiteIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "无效的分站ID", Data: err.Error()})
		return
	}
	
	subsiteInfo, err := h.subsiteService.GetSubsiteInfo(c, subsiteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取分站信息失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: gin.H{"subsite_info": subsiteInfo}})
}

// UpdateSubsite 更新分站
func (h *SubsiteHandler) UpdateSubsite(c *gin.Context) {
	var req UpdateSubsiteRequest
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
	
	subsite, err := h.subsiteService.GetSubsiteByID(c, req.ID)
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
	
	// 更新分站信息
	subsite.Name = req.Name
	subsite.Domain = req.Domain
	subsite.Subdomain = req.Subdomain
	subsite.Logo = req.Logo
	subsite.Description = req.Description
	subsite.Theme = req.Theme
	subsite.Status = req.Status
	subsite.CommissionRate = req.CommissionRate
	
	err = h.subsiteService.UpdateSubsite(c, subsite)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "更新分站失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "更新成功", Data: gin.H{"subsite": subsite}})
}

// ListSubsites 获取分站列表
func (h *SubsiteHandler) ListSubsites(c *gin.Context) {
	// 检查权限
	userRole, exists := c.Get("user_role")
	if !exists || userRole.(string) != "admin" {
		c.JSON(http.StatusOK, conf.Response{
			Code: 403,
			Msg: "需要管理员权限",
			Data: map[string]interface{}{
				"subsites": []interface{}{},
				"total": 0,
			},
		})
		return
	}
	
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	statusStr := c.DefaultQuery("status", "-1")
	
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	status, _ := strconv.Atoi(statusStr)
	
	// 获取分站列表
	subsites, total, err := h.subsiteService.ListSubsites(c, page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusOK, conf.Response{
			Code: 400,
			Msg: err.Error(),
			Data: map[string]interface{}{
				"subsites": []interface{}{},
				"total": 0,
			},
		})
		return
	}
	
	// 即使没有数据，也返回空数组而不是nil
	if subsites == nil {
		subsites = []*model.Subsite{}
	}
	
	// 极简的响应结构体
	type BasicSubsiteInfo struct {
		ID             int64   `json:"id"`
		Name           string  `json:"name"`
		Subdomain      string  `json:"subdomain"`
		Domain         string  `json:"domain"`
		Status         int     `json:"status"`
		CommissionRate float64 `json:"commission_rate"`
	}
	
	var subsiteList []BasicSubsiteInfo
	
	// 极简处理
	for _, subsite := range subsites {
		subsiteList = append(subsiteList, BasicSubsiteInfo{
			ID:             subsite.ID,
			Name:           subsite.Name,
			Subdomain:      subsite.Subdomain,
			Domain:         subsite.Domain,
			Status:         subsite.Status,
			CommissionRate: subsite.CommissionRate,
		})
	}
	
	// 使用最简单的结构
	c.JSON(http.StatusOK, conf.Response{
		Code: 200,
		Msg: "success",
		Data: map[string]interface{}{
			"subsites": subsiteList,
			"total":    total,
		},
	})
}

// SaveSubsiteJsonConfig 保存分站JSON配置
func (h *SubsiteHandler) SaveSubsiteJsonConfig(c *gin.Context) {
	var req SubsiteJsonConfigRequest
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
	
	// 保存分站JSON配置
	err = h.subsiteService.SaveSubsiteJsonConfig(c, req.SubsiteID, req.Config)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "保存配置失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: gin.H{"message": "配置保存成功"}})
}

// GetSubsiteJsonConfig 获取分站JSON配置
func (h *SubsiteHandler) GetSubsiteJsonConfig(c *gin.Context) {
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
	
	config, err := h.subsiteService.GetSubsiteJsonConfig(c, subsiteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "获取分站配置失败", Data: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: gin.H{"config": config}})
} 