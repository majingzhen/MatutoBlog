package models

import (
	"time"

	"gorm.io/gorm"
	"matuto-blog/pkg/utils"
)

// User 博客用户模型
type User struct {
	Id        int       `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Account   string    `json:"account" gorm:"uniqueIndex;size:100;comment:账号"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Password  string    `json:"-" gorm:"size:100;not null;comment:密码"`
	Avatar    string    `json:"avatar" gorm:"size:255;comment:头像URL"`
	Email     string    `json:"email" gorm:"uniqueIndex;size:100;comment:邮箱"`
	Status    int       `json:"status" gorm:"default:0;comment:状态:0正常,1禁用"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

// UserStatus 用户状态常量
const (
	UserStatusActive   = 0 // 激活
	UserStatusDisabled = 1 // 禁用
)

// 使用密码工具类的配置
var DefaultPasswordConfig = utils.DefaultPasswordConfig

// TableName 指定表名
func (User) TableName() string {
	return "m_user"
}

// HashPassword 加密密码
func (u *User) HashPassword(password string) error {
	hash, err := utils.HashPassword(password, DefaultPasswordConfig)
	if err != nil {
		return err
	}
	u.Password = hash
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	return utils.CheckPassword(password, u.Password)
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// ScopeByUsername 按用户名查询
func ScopeByUsername(Account string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("account = ?", Account)
	}
}

// ScopeByEmail 按邮箱查询
func ScopeByEmail(email string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	}
}
