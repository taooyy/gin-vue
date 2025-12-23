package model

import (
	"time"
)

// SysOrganization 组织架构表
type SysOrganization struct {
	ID           uint      `gorm:"primarykey"`
	Name         string    `gorm:"type:varchar(100);not null;comment:组织名称"`
	OrgType      int8      `gorm:"not null;comment:1:平台 2:学校 3:供应商 4:食堂 5:商户"`
	ParentID     uint      `gorm:"not null;default:0;comment:父级ID"`
	ContactName  string    `gorm:"type:varchar(50);comment:负责人"`
	ContactPhone string    `gorm:"type:varchar(20);comment:电话"`
	Address      string    `gorm:"type:varchar(255);comment:地址"`
	Longitude    float64   `gorm:"type:decimal(10,6);comment:经度"`
	Latitude     float64   `gorm:"type:decimal(10,6);comment:纬度"`
	IsEnabled    bool      `gorm:"not null;default:true;comment:状态 1:启用 0:禁用"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (SysOrganization) TableName() string {
	return "sys_organizations"
}

// SysUser 用户/账号表
type SysUser struct {
	ID        uint      `gorm:"primarykey"`
	OrgID     uint      `gorm:"not null;comment:所属组织"`
	Username  string    `gorm:"type:varchar(50);not null;uniqueIndex;comment:账号"`
	Password  string    `gorm:"type:varchar(100);not null;comment:密码"`
	RealName  string    `gorm:"type:varchar(50);not null;comment:真实姓名"`
	Mobile    string    `gorm:"type:varchar(20);comment:手机号"`
	Avatar    string    `gorm:"type:varchar(255);comment:头像"`
	RoleID    uint      `gorm:"not null;default:0;comment:角色ID"`
	Status    int8      `gorm:"not null;default:1;comment:1:正常 2:锁定"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (SysUser) TableName() string {
	return "sys_users"
}

// SysDictionary 数据字典表
type SysDictionary struct {
	ID        uint   `gorm:"primarykey"`
	DictCode  string `gorm:"type:varchar(50);not null;comment:编码"`
	ItemLabel string `gorm:"type:varchar(50);not null;comment:展示名"`
	ItemValue string `gorm:"type:varchar(50);not null;comment:存储值"`
	Sort      int    `gorm:"not null;default:0"`
}

func (SysDictionary) TableName() string {
	return "sys_dictionaries"
}

// SysOpLog 操作日志表
type SysOpLog struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null;comment:操作人"`
	Module    string    `gorm:"type:varchar(50);not null;comment:模块"`
	Action    string    `gorm:"type:varchar(50);not null;comment:动作"`
	Params    string    `gorm:"type:json;comment:参数"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (SysOpLog) TableName() string {
	return "sys_op_logs"
}

// SysBanner 运营广告表
type SysBanner struct {
	ID        uint   `gorm:"primarykey"`
	SchoolID  uint   `gorm:"not null;default:0;comment:定向投放"`
	ImageURL  string `gorm:"type:varchar(255);not null;comment:图片"`
	LinkURL   string `gorm:"type:varchar(255);comment:跳转链"`
	IsEnabled bool   `gorm:"not null;default:true"`
}

func (SysBanner) TableName() string {
	return "sys_banners"
}
