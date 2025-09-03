package models

import (
	"time"
)

// Tag 标签模型
type Tag struct {
	BaseModel
	Name         string     `json:"name" gorm:"size:256;not null;comment:标签名"`
	Color        string     `json:"color" gorm:"size:128;comment:颜色"`
	Thumbnail    string     `json:"thumbnail" gorm:"size:256;comment:缩略图"`
	Slug         string     `json:"slug" gorm:"size:128;index;comment:slug"`
	CreateUserId *uint64    `json:"create_user_id" gorm:"comment:添加人"`
	UpdateUserId *uint64    `json:"update_user_id" gorm:"comment:更新人"`
	CreateTime   time.Time  `json:"create_time" gorm:"not null;comment:创建时间"`
	UpdateTime   *time.Time `json:"update_time" gorm:"comment:修改时间"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "p_tag"
}