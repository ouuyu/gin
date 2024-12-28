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
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误", "error": err.Error()})
		return
	}

	random := rand.Intn(900000) + 100000
	time := time.Now().Unix()
	tradeNo := fmt.Sprintf("%d%d", random, time)

	user, err := model.GetUserById(c.GetInt("id"), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取用户信息失败", "error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建订单失败", "error": err.Error()})
		return
	}

	payHTML, err := common.GeneratePayURL(req.Amount, req.PayType, req.Param, tradeNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建订单失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建订单成功",
		"data":    payHTML,
	})
}

func GetOrderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	orders, err := model.GetOrderList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取订单列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取订单列表成功",
		"data":    orders,
	})
}

func QueryOrder(c *gin.Context) {
	tradeNo := c.Query("trade_no")
	order, err := model.GetOrderByOrderNo(tradeNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "订单不存在", "error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询订单状态失败 " + err.Error(),
			"success": false,
		})
		return
	}

	var Result struct {
		Status string `json:"status"`
	}
	json.Unmarshal(resp.Body(), &Result)

	if Result.Status == "1" {
		order.UpdatePayStatus(1, tradeNo)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": common.Ternary(Result.Status == "1", "已支付", "未支付"),
		"data": gin.H{
			"status": Result.Status,
		},
	})
}
