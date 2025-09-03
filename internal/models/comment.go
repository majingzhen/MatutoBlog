package models

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	BaseModel
	ArticleId    uint64     `json:"article_id" gorm:"not null;comment:文章id"`
	Pid          int        `json:"pid" gorm:"default:-1;comment:父级id"`
	TopPid       int        `json:"top_pid" gorm:"default:-1;comment:顶层父级id"`
	UserId       *uint64    `json:"user_id" gorm:"comment:用户ID"`
	Content      string     `json:"content" gorm:"size:2048;comment:评论内容"`
	Status       int        `json:"status" gorm:"default:0;comment:状态:0正常,1:待审核"`
	Avatar       string     `json:"avatar" gorm:"size:256;comment:头像"`
	Website      string     `json:"website" gorm:"size:256;comment:网站地址"`
	Email        string     `json:"email" gorm:"size:256;comment:邮箱"`
	UserName     string     `json:"user_name" gorm:"size:256;comment:评论人"`
	Ip           string     `json:"ip" gorm:"size:256;comment:ip"`
	Device       string     `json:"device" gorm:"size:256;comment:设备类型"`
	CreateTime   time.Time  `json:"create_time" gorm:"not null;comment:创建时间"`
	UpdateTime   *time.Time `json:"update_time" gorm:"comment:更新时间"`
	CreateUserId *uint64    `json:"create_user_id" gorm:"comment:添加人"`
	UpdateUserId *uint64    `json:"update_user_id" gorm:"comment:更新人"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "p_comment"
}

// CommentStatus 评论状态常量
const (
	CommentStatusActive  = 0 // 正常
	CommentStatusPending = 1 // 待审核
)

// IsActive 检查评论是否激活
func (c *Comment) IsActive() bool {
	return c.Status == CommentStatusActive
}

// IsRoot 检查是否为根评论
func (c *Comment) IsRoot() bool {
	return c.Pid == -1
}

// IsTopLevel 检查是否为顶层评论
func (c *Comment) IsTopLevel() bool {
	return c.TopPid == -1
}