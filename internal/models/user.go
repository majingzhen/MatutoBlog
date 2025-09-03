package models

import (
	"time"

	"gorm.io/gorm"
	"matuto-blog/pkg/utils"
)

// User 博客用户模型
type User struct {
	BaseModel
	Username string     `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Password string     `json:"-" gorm:"size:100;not null;comment:密码"`
	Nickname string     `json:"nickname" gorm:"size:50;comment:昵称"`
	Avatar   string     `json:"avatar" gorm:"size:255;comment:头像URL"`
	Phone    string     `json:"phone" gorm:"uniqueIndex;size:20;comment:手机号"`
	Email    string     `json:"email" gorm:"uniqueIndex;size:100;comment:邮箱"`
	Status   int        `json:"status" gorm:"default:1;comment:状态:1正常,2禁用"`
	LastIP   string     `json:"last_ip" gorm:"size:45;comment:最后登录IP"`
	LastTime *time.Time `json:"last_time" gorm:"comment:最后登录时间"`
}

// UserStatus 用户状态常量
const (
	UserStatusActive   = 1 // 激活
	UserStatusDisabled = 2 // 禁用
)

// 使用密码工具类的配置
var DefaultPasswordConfig = utils.DefaultPasswordConfig

// TableName 指定表名
func (User) TableName() string {
	return "p_user"
}

// BeforeCreate GORM钩子：创建前处理
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的BeforeCreate
	if err := u.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 设置默认昵称
	if u.Nickname == "" {
		u.Nickname = u.Username
	}

	return nil
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

// UpdateLastLogin 更新最后登录信息
func (u *User) UpdateLastLogin(ip string) {
	now := time.Now()
	u.LastIP = ip
	u.LastTime = &now
}

// ScopeByUsername 按用户名查询
func ScopeByUsername(username string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("username = ?", username)
	}
}

// ScopeByPhone 按手机号查询
func ScopeByPhone(phone string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("phone = ?", phone)
	}
}

// ScopeByEmail 按邮箱查询
func ScopeByEmail(email string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	}
}