package model

import (
	"time"
)

// FinBill 资金流水表
type FinBill struct {
	ID          uint      `gorm:"primarykey"`
	BillNo      string    `gorm:"type:varchar(32);not null;comment:流水号"`
	OrderID     uint      `gorm:"not null;comment:关联订单"`
	SchoolID    uint      `gorm:"not null"`
	SupplierID  uint      `gorm:"not null"`
	Amount      float64   `gorm:"type:decimal(12,2);not null"`
	BillType    int8      `gorm:"not null;comment:1:订单收款 2:售后退款"`
	IoDirection int8      `gorm:"not null;comment:1:收入 2:支出"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (FinBill) TableName() string {
	return "fin_bills"
}

// FinStatement 月度对账单
type FinStatement struct {
	ID           uint      `gorm:"primarykey"`
	Period       string    `gorm:"type:varchar(7);not null;comment:账期"`
	SupplierID   uint      `gorm:"not null"`
	SchoolID     uint      `gorm:"not null"`
	SettleAmount float64   `gorm:"type:decimal(12,2);not null;comment:应结金额"`
	Status       int8      `gorm:"not null;default:10;comment:10:待核 20:已核 30:已打款"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (FinStatement) TableName() string {
	return "fin_statements"
}
