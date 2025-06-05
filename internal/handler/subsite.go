package handler

import (
	"net/http"
	"strconv"

	"github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/internal/model"
	"github.com/gin-gonic/gin"
)

// 分站相关处理函数

func createSubsite(c *gin.Context) {
	var req model.CreateSubsiteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}
	
	// 从当前用户获取OwnerID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{
			Code: http.StatusUnauthorized,
			Msg:  "未登录",
		})
		return
	}
	
	ownerID := userID.(int64)
	
	// 获取分站服务
	subsiteService := GetSubsiteService()
	if subsiteService == nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  "服务未初始化",
		})
		return
	}
	
	// 创建分站对象
	subsite := &model.Subsite{
		OwnerID:        ownerID,
		Name:           req.Name,
		Subdomain:      req.Subdomain,
		CommissionRate: req.CommissionRate,
		Status:         1, // 默认启用
		Theme:          "default", // 默认主题
	}
	
	// 调用服务层创建分站
	createdSubsite, err := subsiteService.CreateSubsiteObject(c, subsite)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "创建分站失败: " + err.Error(),
		})
		return
	}
	
	// 返回创建成功的分站信息
	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: gin.H{
			"subsite": createdSubsite,
		},
	})
}

func getSubsiteInfo(c *gin.Context) {
	// 获取分站服务
	subsiteService := GetSubsiteService()
	if subsiteService == nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  "服务未初始化",
		})
		return
	}

	// 获取请求中的分站ID
	subsiteIDStr := c.Query("id")
	var subsiteID int64
	var err error

	if subsiteIDStr == "" {
		// 如果没有提供ID，尝试获取当前用户的分站
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, conf.Response{
				Code: http.StatusUnauthorized,
				Msg:  "未登录",
			})
			return
		}
		ownerID := userID.(int64)
		
		subsite, err := subsiteService.GetSubsiteByOwnerID(c, ownerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, conf.Response{
				Code: http.StatusBadRequest,
				Msg:  "获取分站失败: " + err.Error(),
			})
			return
		}
		if subsite == nil {
			c.JSON(http.StatusBadRequest, conf.Response{
				Code: http.StatusBadRequest,
				Msg:  "未找到分站",
			})
			return
		}
		subsiteID = subsite.ID
	} else {
		// 将ID字符串转换为int64
		subsiteID, err = strconv.ParseInt(subsiteIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, conf.Response{
				Code: http.StatusBadRequest,
				Msg:  "无效的分站ID",
			})
			return
		}
	}

	// 获取分站详细信息
	subsiteInfo, err := subsiteService.GetSubsiteInfo(c, subsiteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "获取分站信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: gin.H{
			"subsite_info": subsiteInfo,
		},
	})
}

func updateSubsite(c *gin.Context) {
	var req model.UpdateSubsiteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 获取分站服务
	subsiteService := GetSubsiteService()
	if subsiteService == nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  "服务未初始化",
		})
		return
	}

	// 检查权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{
			Code: http.StatusUnauthorized,
			Msg:  "未登录",
		})
		return
	}
	ownerID := userID.(int64)

	// 获取分站信息
	subsite, err := subsiteService.GetSubsiteByID(c, req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "获取分站失败: " + err.Error(),
		})
		return
	}
	if subsite == nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "分站不存在",
		})
		return
	}

	// 检查是否为分站所有者或管理员
	userRole, _ := c.Get("user_role")
	role := userRole.(string)
	if subsite.OwnerID != ownerID && role != "admin" {
		c.JSON(http.StatusForbidden, conf.Response{
			Code: http.StatusForbidden,
			Msg:  "无权操作",
		})
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

	err = subsiteService.UpdateSubsite(c, subsite)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "更新分站失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: gin.H{
			"subsite": subsite,
		},
	})
}

func listSubsites(c *gin.Context) {
	// 检查权限
	userRole, exists := c.Get("user_role")
	if !exists || userRole.(string) != "admin" {
		c.JSON(http.StatusOK, conf.Response{
			Code: 403,
			Msg:  "需要管理员权限",
			Data: gin.H{
				"subsites": []interface{}{},
				"total":    0,
			},
		})
		return
	}

	// 获取分站服务
	subsiteService := GetSubsiteService()
	if subsiteService == nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  "服务未初始化",
		})
		return
	}

	// 获取请求参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	statusStr := c.DefaultQuery("status", "-1")
	
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	status, _ := strconv.Atoi(statusStr)
	
	// 获取分站列表
	subsites, total, err := subsiteService.ListSubsites(c, page, pageSize, status)
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
	
	// 确保不返回nil
	if subsites == nil {
		subsites = []*model.Subsite{}
	}
	
	// 精简响应结构体
	type BasicSubsiteInfo struct {
		ID             int64   `json:"id"`
		Name           string  `json:"name"`
		Subdomain      string  `json:"subdomain"`
		Domain         string  `json:"domain"`
		Status         int     `json:"status"`
		CommissionRate float64 `json:"commission_rate"`
	}
	
	var subsiteList []BasicSubsiteInfo
	
	// 处理数据
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

	c.JSON(http.StatusOK, conf.Response{
		Code: 200,
		Msg:  "success",
		Data: gin.H{
			"subsites": subsiteList,
			"total":    total,
		},
	})
}

// handleCreateSubsite 创建分站
// @Summary 创建分站
// @Description 创建一个新的分站
// @Tags subsite
// @Accept json
// @Produce json
// @Param request body model.CreateSubsiteReq true "创建分站请求"
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.Response
// @Router /api/v1/subsite/create [post]
func handleCreateSubsite(c *gin.Context) {
	createSubsite(c)
}

// handleGetSubsiteInfo 获取分站信息
// @Summary 获取分站信息
// @Description 获取分站详细信息
// @Tags subsite
// @Accept json
// @Produce json
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.Response
// @Router /api/v1/subsite/info [get]
func handleGetSubsiteInfo(c *gin.Context) {
	getSubsiteInfo(c)
}

// handleUpdateSubsite 更新分站
// @Summary 更新分站
// @Description 更新分站信息
// @Tags subsite
// @Accept json
// @Produce json
// @Param request body model.UpdateSubsiteReq true "更新分站请求"
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.Response
// @Router /api/v1/subsite/update [put]
func handleUpdateSubsite(c *gin.Context) {
	updateSubsite(c)
}

// handleListSubsites 获取分站列表
// @Summary 获取分站列表
// @Description 获取所有分站列表
// @Tags subsite
// @Accept json
// @Produce json
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.Response
// @Router /api/v1/subsite/list [get]
func handleListSubsites(c *gin.Context) {
	listSubsites(c)
}

// deleteSubsite 删除分站
func deleteSubsite(c *gin.Context) {
	// 获取分站ID
	subsiteIDStr := c.Query("id")
	if subsiteIDStr == "" {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "缺少分站ID",
		})
		return
	}

	subsiteID, err := strconv.ParseInt(subsiteIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "无效的分站ID",
		})
		return
	}

	// 获取分站服务
	subsiteService := GetSubsiteService()
	if subsiteService == nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  "服务未初始化",
		})
		return
	}

	// 检查权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{
			Code: http.StatusUnauthorized,
			Msg:  "未登录",
		})
		return
	}
	ownerID := userID.(int64)

	// 获取分站信息
	subsite, err := subsiteService.GetSubsiteByID(c, subsiteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "获取分站失败: " + err.Error(),
		})
		return
	}
	if subsite == nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "分站不存在",
		})
		return
	}

	// 检查是否为分站所有者或管理员
	userRole, _ := c.Get("user_role")
	role := userRole.(string)
	if subsite.OwnerID != ownerID && role != "admin" {
		c.JSON(http.StatusForbidden, conf.Response{
			Code: http.StatusForbidden,
			Msg:  "无权操作",
		})
		return
	}

	// 删除分站
	err = subsiteService.DeleteSubsite(c, subsiteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "删除分站失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: gin.H{
			"message": "分站已成功删除",
		},
	})
}

// handleDeleteSubsite 删除分站
// @Summary 删除分站
// @Description 删除指定的分站
// @Tags subsite
// @Accept json
// @Produce json
// @Param id query string true "分站ID"
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.Response
// @Router /api/v1/subsite/delete [delete]
func handleDeleteSubsite(c *gin.Context) {
	deleteSubsite(c)
}

// saveSubsiteJsonConfig 保存分站JSON配置
func saveSubsiteJsonConfig(c *gin.Context) {
	var req struct {
		SubsiteID int64                  `json:"subsite_id" binding:"required"`
		Config    map[string]interface{} `json:"config" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 获取分站服务
	subsiteService := GetSubsiteService()
	if subsiteService == nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  "服务未初始化",
		})
		return
	}

	// 检查权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, conf.Response{
			Code: http.StatusUnauthorized,
			Msg:  "未登录",
		})
		return
	}
	ownerID := userID.(int64)

	// 获取分站信息
	subsite, err := subsiteService.GetSubsiteByID(c, req.SubsiteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "获取分站失败: " + err.Error(),
		})
		return
	}
	if subsite == nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "分站不存在",
		})
		return
	}

	// 检查是否为分站所有者或管理员
	userRole, _ := c.Get("user_role")
	role := userRole.(string)
	if subsite.OwnerID != ownerID && role != "admin" {
		c.JSON(http.StatusForbidden, conf.Response{
			Code: http.StatusForbidden,
			Msg:  "无权操作",
		})
		return
	}

	// 保存分站JSON配置
	err = subsiteService.SaveSubsiteJsonConfig(c, req.SubsiteID, req.Config)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: gin.H{
			"message": "配置保存成功",
		},
	})
}

// getSubsiteJsonConfig 获取分站JSON配置
func getSubsiteJsonConfig(c *gin.Context) {
	// 获取分站ID
	subsiteIDStr := c.Query("subsite_id")
	if subsiteIDStr == "" {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "缺少分站ID",
		})
		return
	}

	subsiteID, err := strconv.ParseInt(subsiteIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "无效的分站ID",
		})
		return
	}

	// 获取分站服务
	subsiteService := GetSubsiteService()
	if subsiteService == nil {
		c.JSON(http.StatusInternalServerError, conf.Response{
			Code: http.StatusInternalServerError,
			Msg:  "服务未初始化",
		})
		return
	}

	// 获取分站JSON配置
	config, err := subsiteService.GetSubsiteJsonConfig(c, subsiteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{
			Code: http.StatusBadRequest,
			Msg:  "获取分站配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conf.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: gin.H{
			"config": config,
		},
	})
}

// handleSaveSubsiteJsonConfig 保存分站JSON配置
// @Summary 保存分站JSON配置
// @Description 保存分站JSON格式配置信息
// @Tags subsite
// @Accept json
// @Produce json
// @Param request body object true "保存分站JSON配置请求"
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.Response
// @Router /api/v1/subsite/config [post]
func handleSaveSubsiteJsonConfig(c *gin.Context) {
	saveSubsiteJsonConfig(c)
}

// handleGetSubsiteJsonConfig 获取分站JSON配置
// @Summary 获取分站JSON配置
// @Description 获取分站JSON格式配置
// @Tags subsite
// @Accept json
// @Produce json
// @Param subsite_id query string true "分站ID"
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.Response
// @Router /api/v1/subsite/config [get]
func handleGetSubsiteJsonConfig(c *gin.Context) {
	getSubsiteJsonConfig(c)
}
