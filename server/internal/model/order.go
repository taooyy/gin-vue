package model

import (
	"time"
)

// OrdCart 购物车
type OrdCart struct {
	ID         uint `gorm:"primarykey"`
	MerchantID uint `gorm:"not null;comment:买家"`
	QuoteID    uint `gorm:"not null;comment:选中报价"`
	Quantity   int  `gorm:"not null;default:1"`
}

func (OrdCart) TableName() string {
	return "ord_carts"
}

// OrdOrder 订单主表
type OrdOrder struct {
	ID              uint       `gorm:"primarykey"`
	OrderNo         string     `gorm:"type:varchar(32);not null;uniqueIndex;comment:订单号"`
	MerchantID      uint       `gorm:"not null;comment:买家"`
	SupplierID      uint       `gorm:"not null;comment:卖家"`
	Status          int8       `gorm:"not null;default:10;comment:10:待接 30:配送 40:完成"`
	ReceiptVouchers string     `gorm:"type:json;comment:收货凭证"`
	TotalAmount     float64    `gorm:"type:decimal(12,2);not null;default:0.00"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	DeliveryTime    *time.Time `gorm:"comment:配送时间"`
	ArrivalTime     *time.Time `gorm:"comment:送达时间"`
}

func (OrdOrder) TableName() string {
	return "ord_orders"
}

// OrdOrderItem 订单明细 - 交易快照
type OrdOrderItem struct {
	ID         uint    `gorm:"primarykey"`
	OrderID    uint    `gorm:"not null;index;comment:订单ID"`
	ProductID  uint    `gorm:"not null;comment:商品ID"`
	QuoteID    uint    `gorm:"not null;comment:报价ID"`
	SnapName   string  `gorm:"type:varchar(100);not null;comment:快照:品名"`
	SnapSpecs  string  `gorm:"type:varchar(100);not null;comment:快照:规格"`
	SnapPrice  float64 `gorm:"type:decimal(10,2);not null;comment:快照:单价"`
	Quantity   int     `gorm:"not null;comment:数量"`
	Amount     float64 `gorm:"type:decimal(12,2);not null;comment:该项总价"`
}

func (OrdOrderItem) TableName() string {
	return "ord_order_items"
}

// OrdAfterSale 售后表
type OrdAfterSale struct {
	ID          uint    `gorm:"primarykey"`
	OrderItemID uint    `gorm:"not null;comment:关联明细"`
	Type        int8    `gorm:"not null;comment:1:仅退款 2:退货退款"`
	Reason      string  `gorm:"type:varchar(255);not null;comment:原因"`
	ApplyAmount float64 `gorm:"type:decimal(10,2);not null;comment:金额"`
	Status      int8    `gorm:"not null;default:10;comment:10:待审 20:通过 30:驳回"`
}

func (OrdAfterSale) TableName() string {
	return "ord_after_sales"
}

// OrdItemTrace 全链路溯源表
type OrdItemTrace struct {
	ID            uint      `gorm:"primarykey"`
	OrderItemID   uint      `gorm:"not null;unique;comment:明细ID"`
	TraceCode     string    `gorm:"type:varchar(64);not null;unique;comment:溯源码"`
	MerchantInfo  string    `gorm:"type:json;comment:商户快照"`
	SupplierInfo  string    `gorm:"type:json;comment:供应商快照"`
	CertSnapshot  string    `gorm:"type:json;comment:证书快照"`
	TimeLine      string    `gorm:"type:json;comment:时间轴"`
	DriverID      uint      `gorm:"default:0;comment:司机"`
	QcCertImage   string    `gorm:"type:varchar(255);comment:实物检测图"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (OrdItemTrace) TableName() string {
	return "ord_item_traces"
}
