package controller

import (
	"main/common"
	"main/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RechargeRequest 充值请求结构
type RechargeRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	PayType     string  `json:"pay_type" binding:"required,oneof=alipay wxpay"`
	Description string  `json:"description"`
}

// Recharge 用户充值
func Recharge(c *gin.Context) {
	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userId := c.GetInt("id")
	user, err := model.GetUserById(userId, false)
	if err != nil {
		common.Fail(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 创建充值订单
	order := model.Order{
		UserID:      user.ID,
		Amount:      req.Amount,
		PayType:     req.PayType,
		ProductName: "余额充值",
		Param:       req.Description,
	}

	// 生成支付链接
	payHTML, err := order.CreateRechargeOrder()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "创建充值订单失败: "+err.Error())
		return
	}

	common.Success(c, gin.H{
		"pay_url":  payHTML,
		"order_no": order.OrderNo,
	})
}

// GetBalanceLogs 获取余额变动记录
func GetBalanceLogs(c *gin.Context) {
	userId := c.GetInt("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize
	logs, err := model.GetUserBalanceLogs(userId, offset, pageSize)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "获取余额记录失败: "+err.Error())
		return
	}

	total, err := model.GetUserBalanceLogCount(userId)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "获取余额记录失败: "+err.Error())
		return
	}

	common.Success(c, gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
	})
}

// GetBalance 获取用户余额
func GetBalance(c *gin.Context) {
	userId := c.GetInt("userId")
	user, err := model.GetUserById(userId, true)
	if err != nil {
		common.Fail(c, http.StatusNotFound, "用户不存在")
		return
	}

	common.Success(c, gin.H{
		"balance": user.Balance,
	})
}
