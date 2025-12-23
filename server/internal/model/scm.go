package model

import (
	"time"
)

// ScmCategory 商品分类表
type ScmCategory struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(50);not null;comment:分类名"`
	ParentID  uint   `gorm:"not null;default:0"`
	Icon      string `gorm:"type:varchar(255);comment:图标"`
	Sort      int    `gorm:"not null;default:0"`
}

func (ScmCategory) TableName() string {
	return "scm_categories"
}

// ScmProduct 商品标准库 SPU
type ScmProduct struct {
	ID          uint      `gorm:"primarykey"`
	SchoolID    uint      `gorm:"not null;index;comment:监管学校"`
	CategoryID  uint      `gorm:"not null;comment:分类"`
	Name        string    `gorm:"type:varchar(100);not null;comment:商品名称"`
	Image       string    `gorm:"type:varchar(255);comment:标准图"`
	Specs       string    `gorm:"type:varchar(100);not null;comment:固定规格"`
	Unit        string    `gorm:"type:varchar(20);not null;comment:计量单位"`
	SourceType  int8      `gorm:"not null;default:3;comment:1:平台下发 2:学校自建 3:供应商上传"`
	CreatorID   uint      `gorm:"default:0;comment:创建者ID"`
	StaticCerts string    `gorm:"type:json;comment:三证资质"`
	AuditStatus int8      `gorm:"not null;default:0;comment:0:待审 1:通过 2:驳回"`
	IsListed    bool      `gorm:"not null;default:false;comment:上架状态"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (ScmProduct) TableName() string {
	return "scm_products"
}

// ScmProductQuote 供应商报价 SKU
type ScmProductQuote struct {
	ID           uint      `gorm:"primarykey"`
	ProductID    uint      `gorm:"not null;uniqueIndex:uk_prod_sup;comment:关联SPU"`
	SupplierID   uint      `gorm:"not null;uniqueIndex:uk_prod_sup;comment:报价人"`
	Price        float64   `gorm:"type:decimal(10,2);not null;comment:报价"`
	BatchReports string    `gorm:"type:json;comment:批次报告"`
	IsEnabled    bool      `gorm:"not null;default:true;comment:供货开关"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (ScmProductQuote) TableName() string {
	return "scm_product_quotes"
}

// ScmSupplierStaff 供应商员工表
type ScmSupplierStaff struct {
	ID          uint   `gorm:"primarykey"`
	SupplierID  uint   `gorm:"not null;comment:所属供应商"`
	Name        string `gorm:"type:varchar(50);not null;comment:姓名"`
	Mobile      string `gorm:"type:varchar(20);not null"`
	RoleType    int8   `gorm:"not null;comment:1:司机 2:分拣"`
	HealthCert  string `gorm:"type:varchar(255);comment:健康证图"`
}

func (ScmSupplierStaff) TableName() string {
	return "scm_supplier_staffs"
}
