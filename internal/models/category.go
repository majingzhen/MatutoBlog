package models

import (
	"time"
)

// Category 分类模型
type Category struct {
	BaseModel
	Name            string     `json:"name" gorm:"size:256;not null;comment:分类名"`
	Pid             int        `json:"pid" gorm:"default:-1;comment:父级id"`
	Desc            string     `json:"desc" gorm:"size:512;comment:描述"`
	MetaKeywords    string     `json:"meta_keywords" gorm:"size:256;comment:SEO关键字"`
	Thumbnail       string     `json:"thumbnail" gorm:"size:256;comment:封面图"`
	Slug            string     `json:"slug" gorm:"size:128;index;comment:slug"`
	MetaDescription string     `json:"meta_description" gorm:"size:256;comment:SEO描述内容"`
	Status          int        `json:"status" gorm:"default:0;comment:状态0:正常,1禁用"`
	CreateTime      time.Time  `json:"create_time" gorm:"not null;comment:创建时间"`
	UpdateTime      *time.Time `json:"update_time" gorm:"comment:更新时间"`
	CreateUserId    *uint64    `json:"create_user_id" gorm:"comment:添加人"`
	UpdateUserId    *uint64    `json:"update_user_id" gorm:"comment:更新人"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "p_category"
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