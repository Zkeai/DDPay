package handler

import (
	"encoding/json"
	"github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/common/logger"
	"github.com/Zkeai/DDPay/common/utils"
	config "github.com/Zkeai/DDPay/internal/conf"
	"github.com/Zkeai/DDPay/internal/model"
	"github.com/gin-gonic/gin"

	"net/http"
)

// createTransaction 创建订单
// @Tags order
// @Summary 创建订单
// @Param a req body model.OrderReq true "创建订单"
// @Router /order/create-transaction [post]
// @Success 200 {object} conf.Response
// @Failure 400 {object} string "参数错误"
// @Failure 500 {object} string "内部错误"
// @Produce JSON
// @Accept JSON
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
// @Param order query string true "订单key"
// @Router /pay/status [get]
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.ResponseError
// @Failure 500 {object} conf.ResponseError
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

func signVerify(c *gin.Context) {
	rawData, err := c.GetRawData()

	if err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "读取请求失败", Data: err.Error()})
		c.Abort()
		return
	}

	var m = make(map[string]any)
	if err = json.Unmarshal(rawData, &m); err != nil {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "JSON解析失败", Data: err.Error()})
		c.Abort()
		return
	}

	sign, ok := m["signature"]
	if !ok {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "签名缺失"})
		c.Abort()
		return
	}
	token := config.Get().Config.SignKey
	if utils.Sign(m, token) != sign {
		c.JSON(http.StatusUnauthorized, conf.Response{Code: 401, Msg: "签名验证失败"})
		c.Abort()
		return
	}

	c.Set("data", m)
	c.Request.Body = utils.RestoreBody(rawData)
}
