package handler

import (
	"net/http"

	"github.com/Zkeai/DDPay/common/conf"
	chttp "github.com/Zkeai/DDPay/common/net/cttp"
	_ "github.com/Zkeai/DDPay/docs"
	"github.com/Zkeai/DDPay/internal/service"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

var (
	svc *service.Service
)

func InitRouter(s *chttp.Server, service *service.Service) {
	svc = service

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
}




