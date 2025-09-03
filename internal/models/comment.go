package models

// Comment 评论模型
type Comment struct {
	BaseModel
	ArticleId int    `json:"articleId" gorm:"not null;comment:文章id"`
	Pid       int    `json:"pId" gorm:"default:-1;comment:父级id"`
	TopPid    int    `json:"topPId" gorm:"default:-1;comment:顶层父级id"`
	UserId    *int   `json:"userId" gorm:"comment:用户ID"`
	Content   string `json:"content" gorm:"size:2048;comment:评论内容"`
	Status    int    `json:"status" gorm:"default:0;comment:状态:0正常,1:待审核"`
	Avatar    string `json:"avatar" gorm:"size:256;comment:头像"`
	Website   string `json:"website" gorm:"size:256;comment:网站地址"`
	Email     string `json:"email" gorm:"size:256;comment:邮箱"`
	Username  string `json:"username" gorm:"size:256;comment:评论人"`
	Ip        string `json:"ip" gorm:"size:256;comment:ip"`
	Device    string `json:"device" gorm:"size:256;comment:设备类型"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "m_comment"
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
