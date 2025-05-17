package handler

import (
	"github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/common/logger"
	"github.com/Zkeai/DDPay/internal/model"

	"github.com/gin-gonic/gin"

	"net/http"
)

// test 接口测试
// @Tags  test
// @Summary 接口测试
// @Param msg query string true "测试消息"
// @Router /test/test [get]
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.ResponseError
// @Failure 500 {object} string "内部错误"
func test(c *gin.Context) {
	r := new(model.TestReq)

	if err := c.Bind(r); err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, conf.Response{Code: 500, Msg: "err", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: r.Msg})
}
