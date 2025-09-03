package models

import (
	"time"
)

// Link 友情链接模型
type Link struct {
	BaseModel
	Name         string     `json:"name" gorm:"size:256;not null;comment:网站名"`
	Logo         string     `json:"logo" gorm:"size:256;comment:网站logo"`
	Desc         string     `json:"desc" gorm:"size:512;comment:网站描述"`
	Address      string     `json:"address" gorm:"size:256;not null;comment:网站地址"`
	CreateTime   time.Time  `json:"create_time" gorm:"not null;comment:创建时间"`
	UpdateTime   *time.Time `json:"update_time" gorm:"comment:更新时间"`
	CreateUserId *uint64    `json:"create_user_id" gorm:"comment:添加人"`
	UpdateUserId *uint64    `json:"update_user_id" gorm:"comment:更新人"`
}

// TableName 指定表名
func (Link) TableName() string {
	return "p_link"
}