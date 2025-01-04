package model

import (
	"fmt"
	"main/common"
	"math/rand"
	"time"
)

type Order struct {
	ID          int       `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	OrderNo     string    `json:"order_no" gorm:"unique;index"`
	UserID      int       `json:"user_id" gorm:"index"`
	Amount      float64   `json:"amount"`
	PayType     string    `json:"pay_type"`
	Status      int       `json:"status" gorm:"default:0"` // 0:未支付 1:已支付
	ProductName string    `json:"product_name"`
	Param       string    `json:"param"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Create 创建订单
func (o *Order) Create() error {
	random := rand.Intn(900000) + 100000
	timestamp := time.Now().Unix()
	o.OrderNo = fmt.Sprintf("%d%d", random, timestamp)
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return DB.Create(o).Error
}

// CreateRechargeOrder 创建充值订单并生成支付链接
func (o *Order) CreateRechargeOrder() (string, error) {
	if err := o.Create(); err != nil {
		return "", err
	}

	// 生成支付链接
	payHTML, err := common.GeneratePayURL(o.Amount, o.PayType, o.Param, o.OrderNo)
	if err != nil {
		return "", err
	}

	return payHTML, nil
}

// UpdatePayStatus 更新支付状态
func (o *Order) UpdatePayStatus(status int, orderNo string) error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新订单状态
	if err := tx.Model(o).Where("order_no = ?", orderNo).Update("status", status).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果是充值订单且支付成功，更新用户余额
	if status == 1 && o.ProductName == "余额充值" {
		user, err := GetUserById(o.UserID, false)
		if err != nil {
			tx.Rollback()
			return err
		}

		// 增加用户余额
		if err := user.AddBalance(o.Amount); err != nil {
			tx.Rollback()
			return err
		}

		// 记录余额变动
		if err := CreateBalanceLog(o.UserID, BalanceLogTypeRecharge, o.Amount, user.Balance, o.Param); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetOrderByOrderNo 根据订单号获取订单
func GetOrderByOrderNo(orderNo string) (*Order, error) {
	var order Order
	err := DB.Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrderList 获取订单列表
func GetOrderList(page, pageSize int) ([]Order, error) {
	var orders []Order
	offset := (page - 1) * pageSize
	err := DB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, err
}
