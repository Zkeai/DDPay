package handler

import (
	"errors"
	"net/http"

	"github.com/Zkeai/DDPay/common/conf"
	chttp "github.com/Zkeai/DDPay/common/net/cttp"
	_ "github.com/Zkeai/DDPay/docs"
	"github.com/Zkeai/DDPay/internal/model"
	"github.com/Zkeai/DDPay/internal/service"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

var (
	svc *service.Service
)

// Handler 处理器结构体
type Handler struct {
	userService      *service.Service
	subsiteService   service.SubsiteService
	membershipService service.MembershipService
}

// NewHandler 创建处理器
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		userService:      service,
		subsiteService:   service.GetSubsiteService(),
		membershipService: service.GetMembershipService(),
	}
}

// getCurrentUser 获取当前用户
func (h *Handler) getCurrentUser(c *gin.Context) (*model.User, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, errors.New("未找到用户ID")
	}

	user, err := h.userService.GetUserByID(c, userID.(int64))
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("用户不存在")
	}

	return user, nil
}

func InitRouter(s *chttp.Server, service *service.Service) {
	svc = service
	h := NewHandler(service)

	g := s.Group("/api/v1")
	g.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	g.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: "AIDog"})
	})

	// 用户相关路由
	ug := g.Group("/user")
	{
		// 无需认证的路由
		ug.POST("/register", register)
		ug.POST("/login", login)
		ug.POST("/oauth/login", oauthLogin)
		ug.POST("/send-code", sendCode)
		ug.POST("/reset-password", resetPassword)
		ug.GET("/check-email", checkEmail)
		
		// 刷新访问令牌的路由
		ug.POST("/refresh-token", refreshToken)
		
		// 需要认证的路由
		auth := ug.Group("")
		auth.Use(jwtAuth())
		{
			auth.GET("/profile", getUserProfile)
			auth.PUT("/profile", updateUserProfile)
			auth.POST("/logout", logout)
			auth.GET("/login-logs", getLoginLogs)
			
			// 管理员路由
			admin := auth.Group("/admin")
			{
				admin.PUT("/update-level", h.HandleUpdateUserLevel)
				admin.PUT("/update-status", h.HandleUpdateUserStatus)
				admin.GET("/list", h.HandleListUsers)
			}
		}
	}

	// 订单相关路由
	wg := g.Group("/order")
	wg.Use(signVerify)
	{
		wg.POST("/create-transaction", createTransaction)
	}

	pg := g.Group("/pay")
	{
		pg.GET("/status", getOrderStatus)
	}
	
	// 分站相关路由
	// 注册分站相关路由，需要JWT认证
	subsiteGroup := g.Group("/subsite")
	subsiteGroup.Use(jwtAuth())
	{
		// 分站管理
		subsiteGroup.POST("/create", handleCreateSubsite)
		subsiteGroup.GET("/info", handleGetSubsiteInfo)
		subsiteGroup.PUT("/update", handleUpdateSubsite)
		subsiteGroup.GET("/list", handleListSubsites)
		subsiteGroup.DELETE("/delete", handleDeleteSubsite)
		
		// 分站JSON配置
		subsiteGroup.POST("/config", handleSaveSubsiteJsonConfig)
		subsiteGroup.GET("/config", handleGetSubsiteJsonConfig)
	}

	// 会员等级相关路由
	membershipGroup := g.Group("/membership")
	{
		// 无需认证的路由
		membershipGroup.GET("/levels", h.HandleGetMembershipLevels)
		
		// 需要认证的路由
		memberAuth := membershipGroup.Group("")
		memberAuth.Use(jwtAuth())
		{
			memberAuth.GET("/info", h.HandleGetUserMembership)
			memberAuth.POST("/purchase", h.HandlePurchaseMembership)
			memberAuth.POST("/upgrade", h.HandleUpgradeMembership)
			memberAuth.POST("/renew", h.HandleRenewMembership)
			memberAuth.GET("/check-requirements", h.HandleCheckMembershipRequirements)
			memberAuth.GET("/transactions", h.HandleGetMembershipTransactions)
			
			// 管理员路由 - 使用adminCheck中间件验证是否为管理员
			memberAuth.POST("/admin/level", adminCheck(), h.HandleCreateMembershipLevel)
			memberAuth.PUT("/admin/level", adminCheck(), h.HandleUpdateMembershipLevel)
			memberAuth.DELETE("/admin/level", adminCheck(), h.HandleDeleteMembershipLevel)
		}
	}
}

// GetSubsiteService 获取分站服务接口
func GetSubsiteService() service.SubsiteService {
	if svc == nil {
		return nil
	}
	
	// 调用Service的GetSubsiteService方法
	return svc.GetSubsiteService()
}







