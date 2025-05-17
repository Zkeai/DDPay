package handler

import (
	ws2 "github.com/Zkeai/go_template/internal/ws"
	"log"
	"net/http"
	"os"

	"github.com/Zkeai/go_template/common/conf"
	"github.com/Zkeai/go_template/common/middleware"
	chttp "github.com/Zkeai/go_template/common/net/cttp"
	_ "github.com/Zkeai/go_template/docs"
	"github.com/Zkeai/go_template/internal/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

var (
	svc   *service.Service
	wsHub *ws2.Hub
)

func InitRouter(s *chttp.Server, service *service.Service, hub *ws2.Hub) {
	wsHub = hub
	svc = service
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// 初始化 Casbin Enforcer
	e, err := casbin.NewEnforcer(dir+"/common/conf/rbac_model.conf", dir+"/common/conf/rbac_policy.csv")
	if err != nil {
		log.Fatalf("Failed to initialize Casbin enforcer: %v", err)
	}

	g := s.Group("/api/v1")
	g.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	g.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: "AIDog"})
	})

	// WebSocket
	// 注册 WebSocket 路由（支持模块路径）
	wg := s.Group("/ws/v1")
	wg.GET("/:module", func(c *gin.Context) {
		module := c.Param("module")
		ws2.ServeWs(wsHub, c.Writer, c.Request, module, "")
	})

	wg.GET("/:module/:target", func(c *gin.Context) {
		module := c.Param("module")
		target := c.Param("target")
		ws2.ServeWs(wsHub, c.Writer, c.Request, module, target)
	})

	// 用户相关路由
	ug := g.Group("/user")
	ugpub := ug.Group("/public")
	{
		ugpub.POST("/login", userLogin)
		ugpub.POST("/logout", userLogout)
	}
	ugpro := ug.Group("/protected")

	ugpro.Use(middleware.Middleware())
	ugpro.Use(middleware.CasbinMiddleware(e))
	{
		ugpro.GET("/query", userQuery)
	}

	//会员相关路由
	membership := g.Group("/membership")
	membershippub := membership.Group("/protected")
	membershippub.Use(middleware.Middleware())
	membershippub.Use(middleware.CasbinMiddleware(e))
	{
		membershippub.POST("/add", AddMembership)
		membershippub.GET("/current", GetCurrentMembership)
		membershippub.GET("/user/:user_id", GetMembershipsByUserID)
		membershippub.GET("/by_order", GetMembershipByOrderID)

		// 管理员
		membershippub.GET("/admin/all", GetAllMemberships)
	}

	// 订单相关路由
	orderGroup := g.Group("/order")
	orderProtected := orderGroup.Group("/protected")
	orderProtected.Use(middleware.Middleware())
	orderProtected.Use(middleware.CasbinMiddleware(e))
	{
		orderProtected.POST("/create", createOrder)
		orderProtected.GET("/query_by_id", queryOrderByID)
		orderProtected.GET("/query_by_txhash", queryOrderByTxHash)
		orderProtected.POST("/query", queryOrdersWithPagination)
		orderProtected.POST("/count", countOrders)
		orderProtected.POST("/update_status", updateOrderStatus)
	}

	//通道相关路由
	channelGroup := g.Group("/channel")
	channelGroup.Use(middleware.Middleware())
	{
		channelGroup.POST("/upsert", UpsertChannel)
		channelGroup.GET("/list", GetChannelsByUserID)
		channelGroup.POST("/status", UpdateChannelStatus)
		channelGroup.POST("/deleted", DeleteChannel)
		channelGroup.GET("/create/telegram/personal", CreateChannelByTg)
		channelGroup.GET("/create/telegram/personal/getByKey", GetChannelRe)

		channelGroup.GET("/test/telegram", TestTelegram)
	}

	//监控相关
	monitorGroup := g.Group("/monitor")
	monitorGroup.Use(middleware.Middleware())
	orderProtected.Use(middleware.CasbinMiddleware(e))
	{
		monitorGroup.POST("/monitorAccountsWithSolana", monitorAccountsWithSolana)
	}

}
