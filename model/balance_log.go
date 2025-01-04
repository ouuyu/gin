package model

import (
	"time"
)

// BalanceLogType 余额变动类型
type BalanceLogType int

const (
	BalanceLogTypeRecharge BalanceLogType = iota + 1 // 充值
	BalanceLogTypeConsume                            // 消费
	BalanceLogTypeRefund                             // 退款
)

// BalanceLog 余额变动记录
type BalanceLog struct {
	ID          int            `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	UserId      int            `json:"user_id" gorm:"index"`
	Type        BalanceLogType `json:"type"`
	Amount      float64        `json:"amount" gorm:"type:decimal(10,2)"`
	Balance     float64        `json:"balance" gorm:"type:decimal(10,2)"` // 变动后的余额
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
}

// CreateBalanceLog 创建余额变动记录
func CreateBalanceLog(userId int, logType BalanceLogType, amount, balance float64, description string) error {
	log := &BalanceLog{
		UserId:      userId,
		Type:        logType,
		Amount:      amount,
		Balance:     balance,
		Description: description,
		CreatedAt:   time.Now(),
	}
	return DB.Create(log).Error
}

// GetUserBalanceLogs 获取用户的余额变动记录
func GetUserBalanceLogs(userId int, offset, limit int) ([]BalanceLog, error) {
	var logs []BalanceLog
	err := DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetUserBalanceLogCount 获取用户的余额变动记录总数
func GetUserBalanceLogCount(userId int) (int64, error) {
	var count int64
	err := DB.Model(&BalanceLog{}).Where("user_id = ?", userId).Count(&count).Error
	return count, err
}
