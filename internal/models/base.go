package models

import (
	"time"
)

// BaseModel 基础模型，包含通用字段
type BaseModel struct {
	Id        int       `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	CreatedBy int       `json:"created_by" gorm:"comment:创建人"`
	UpdatedBy int       `json:"updated_by" gorm:"comment:更新人"`
}

// 状态常量定义
const (
	StatusDisabled = 0 // 禁用
	StatusActive   = 1 // 正常
)
