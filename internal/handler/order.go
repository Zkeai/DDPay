package handler

import (
	"math"
	"net/http"

	"github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/common/logger"
	"github.com/Zkeai/DDPay/internal/model"
	"github.com/gin-gonic/gin"
)

// createTransaction 创建订单
// @Tags order
// @Summary 创建订单
// @Accept json
// @Produce json
// @Param request body model.OrderReq true "创建订单"
// @Router /order/create-transaction [post]
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.Response
// @Failure 500 {object} conf.Response
func createTransaction(c *gin.Context) {
	var req model.OrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("bind error: %v", err)
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "参数错误", Data: err.Error()})
		return
	}

	res, err := svc.CreateOrder(c.Request.Context(), req)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: res})
}

// getOrderStatus 获取当前订单状态
// @Tags order
// @Summary 获取当前订单状态
// @Accept json
// @Produce json
// @Param order query string true "订单key"
// @Router /pay/status [get]
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.Response
// @Failure 500 {object} conf.Response
func getOrderStatus(c *gin.Context) {

	orderKey := c.Query("order")
	if orderKey == "" {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "invalid order", Data: "订单不存在"})
		return
	}

	info := svc.GetOrderStatus(orderKey)
	infoNew := model.RedisWallet{
		Amount:     info.Amount,
		Address:    info.Address,
		Chain:      info.Chain,
		Status:     info.Status,
		MerchantID: info.MerchantID,
	}
	if infoNew.Address == "" {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "invalid order", Data: "订单不存在"})
		return
	}
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: infoNew})

}



func RoundFloat(x float64, prec int) float64 {
	pow := math.Pow(10, float64(prec))
	return math.Round(x*pow) / pow
}
