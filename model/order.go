package model

import (
	"time"
)

type Order struct {
	ID          int       `json:"id" gorm:"primarykey"`
	OrderNo     string    `json:"order_no" gorm:"uniqueIndex;not null"` // 商户订单号
	TradeNo     string    `json:"trade_no"`                             // 易支付订单号
	UserID      int       `json:"user_id" gorm:"not null"`
	ProductName string    `json:"product_name" gorm:"not null"`
	Amount      float64   `json:"amount" gorm:"not null"`
	PayType     string    `json:"pay_type"`                    // alipay or wxpay
	PayStatus   int       `json:"pay_status" gorm:"default:0"` // 0:未支付 1:已支付
	Param       string    `json:"param"`                       // 附加参数
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (o *Order) Create() error {
	return DB.Create(o).Error
}

func (o *Order) Update() error {
	return DB.Save(o).Error
}

func GetOrderByOrderNo(orderNo string) (*Order, error) {
	var order Order
	err := DB.Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

func GetOrderList(page int, pageSize int) ([]Order, error) {
	var orders []Order
	err := DB.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error
	return orders, err
}

func (o *Order) UpdatePayStatus(status int, tradeNo string) error {
	o.PayStatus = status
	o.TradeNo = tradeNo
	return o.Update()
}
