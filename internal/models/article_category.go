package models

import (
	"time"
)

// ArticleCategory 文章分类关联模型
type ArticleCategory struct {
	ID           uint64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	ArticleId    uint64     `json:"article_id" gorm:"not null;comment:文章id"`
	CategoryId   uint64     `json:"category_id" gorm:"not null;comment:分类id"`
	CreateTime   time.Time  `json:"create_time" gorm:"not null;comment:创建时间"`
	UpdateTime   *time.Time `json:"update_time" gorm:"comment:更新时间"`
	CreateUserId *uint64    `json:"create_user_id" gorm:"comment:添加人"`
	UpdateUserId *uint64    `json:"update_user_id" gorm:"comment:更新人"`
}

// TableName 指定表名
func (ArticleCategory) TableName() string {
	return "p_article_category"
}