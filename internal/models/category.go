package models

// Category 分类模型
type Category struct {
	BaseModel
	Name            string `json:"name" gorm:"size:256;not null;comment:分类名"`
	Pid             int    `json:"pid" gorm:"default:-1;comment:父级id"`
	Desc            string `json:"desc" gorm:"size:512;comment:描述"`
	MetaKeywords    string `json:"meta_keywords" gorm:"size:256;comment:SEO关键字"`
	Thumbnail       string `json:"thumbnail" gorm:"size:256;comment:封面图"`
	Slug            string `json:"slug" gorm:"size:128;index;comment:slug"`
	MetaDescription string `json:"meta_description" gorm:"size:256;comment:SEO描述内容"`
	Status          int    `json:"status" gorm:"default:0;comment:状态0:正常,1禁用"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "m_category"
}

// CategoryStatus 分类状态常量
const (
	CategoryStatusActive   = 0 // 正常
	CategoryStatusDisabled = 1 // 禁用
)

// IsActive 检查分类是否激活
func (c *Category) IsActive() bool {
	return c.Status == CategoryStatusActive
}

// IsRoot 检查是否为根分类
func (c *Category) IsRoot() bool {
	return c.Pid == -1
}
