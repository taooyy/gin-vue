package model

// OrgType 定义了组织或用户的类型
type OrgType int8

const (
	OrgTypePlatform OrgType = iota // 0: 平台
	OrgTypeSchool                  // 1: 学校
	OrgTypeSupplier                // 2: 供应商
	OrgTypeCanteen                 // 3: 食堂
	OrgTypeMerchant                // 4: 商户
)

// String 方法返回 OrgType 的字符串表示
func (o OrgType) String() string {
	switch o {
	case OrgTypePlatform:
		return "平台"
	case OrgTypeSchool:
		return "学校"
	case OrgTypeSupplier:
		return "供应商"
	case OrgTypeCanteen:
		return "食堂"
	case OrgTypeMerchant:
		return "商户"
	default:
		return "未知"
	}
}
