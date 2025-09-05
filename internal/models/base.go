package models

import (
	"time"
)

// BaseModel 基础模型，包含通用字段
type BaseModel struct {
	Id        int       `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:主键ID"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
	CreatedBy int       `json:"createdBy" gorm:"column:created_by;comment:创建人"`
	UpdatedBy int       `json:"updatedBy" gorm:"column:updated_by;comment:更新人"`
}

// 状态常量定义
const (
	StatusDisabled = 0 // 禁用
	StatusActive   = 1 // 正常
)
