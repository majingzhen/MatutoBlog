package models

// Link 友情链接模型
type Link struct {
	BaseModel
	Name    string `json:"name" gorm:"size:256;not null;comment:网站名"`
	Logo    string `json:"logo" gorm:"size:256;comment:网站logo"`
	Desc    string `json:"desc" gorm:"size:512;comment:网站描述"`
	Address string `json:"address" gorm:"size:256;not null;comment:网站地址"`
}

// TableName 指定表名
func (Link) TableName() string {
	return "m_link"
}
