package models

// Tag 标签模型
type Tag struct {
	BaseModel
	Name      string `json:"name" gorm:"size:256;not null;comment:标签名"`
	Color     string `json:"color" gorm:"size:128;comment:颜色"`
	Thumbnail string `json:"thumbnail" gorm:"size:256;comment:缩略图"`
	Slug      string `json:"slug" gorm:"size:128;index;comment:slug"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "m_tag"
}
