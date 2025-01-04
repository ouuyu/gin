package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"

	"main/common"
	"main/model"
)

type CreateOrderRequest struct {
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	PayType string  `json:"pay_type" binding:"required,oneof=alipay wxpay"`
	Param   string  `json:"param"`
}

func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	random := rand.Intn(900000) + 100000
	time := time.Now().Unix()
	tradeNo := fmt.Sprintf("%d%d", random, time)

	user, err := model.GetUserById(c.GetInt("id"), false)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "获取用户信息失败: "+err.Error())
		return
	}
	var order model.Order
	order.OrderNo = tradeNo
	order.UserID = user.ID
	order.Amount = req.Amount
	order.PayType = req.PayType
	order.Param = req.Param
	order.ProductName = "充值"

	err = order.Create()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "创建订单失败: "+err.Error())
		return
	}

	payHTML, err := common.GeneratePayURL(req.Amount, req.PayType, req.Param, tradeNo)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "创建订单失败: "+err.Error())
		return
	}

	common.Success(c, payHTML)
}

func GetOrderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	orders, err := model.GetOrderList(page, pageSize)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "获取订单列表失败: "+err.Error())
		return
	}

	common.Success(c, orders)
}

func QueryOrder(c *gin.Context) {
	tradeNo := c.Param("trade_no")
	order, err := model.GetOrderByOrderNo(tradeNo)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "订单不存在: "+err.Error())
		return
	}

	apiUrl, _ := url.Parse(common.EasyPayURL)
	apiUrl = apiUrl.JoinPath("api.php")
	client := resty.New()

	resp, err := client.R().SetQueryParams(map[string]string{
		"act":          "order",
		"pid":          common.EasyPayPid,
		"key":          common.EasyPayKey,
		"out_trade_no": tradeNo,
	}).Get(apiUrl.String())

	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "查询订单状态失败: "+err.Error())
		return
	}

	var Result struct {
		Status string `json:"status"`
	}
	json.Unmarshal(resp.Body(), &Result)

	if Result.Status == "1" {
		order.UpdatePayStatus(1, tradeNo)
	}

	common.Success(c, gin.H{
		"status": Result.Status,
	})
}
