package models

import (
	"time"
)

// ArticleTag 文章标签关联模型
type ArticleTag struct {
	ID           uint64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	ArticleId    uint64     `json:"article_id" gorm:"not null;comment:文章id"`
	TagId        uint64     `json:"tag_id" gorm:"not null;comment:标签id"`
	CreateTime   *time.Time `json:"create_time" gorm:"comment:创建时间"`
	UpdateTime   *time.Time `json:"update_time" gorm:"comment:更新时间"`
	CreateUserId *uint64    `json:"create_user_id" gorm:"comment:添加人"`
	UpdateUserId *uint64    `json:"update_user_id" gorm:"comment:更新人"`
}

// TableName 指定表名
func (ArticleTag) TableName() string {
	return "p_article_tag"
}