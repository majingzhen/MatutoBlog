package models

// ArticleTag 文章标签关联模型
type ArticleTag struct {
	BaseModel
	ArticleId int `json:"article_id" gorm:"not null;comment:文章id"`
	TagId     int `json:"tag_id" gorm:"not null;comment:标签id"`
}

// TableName 指定表名
func (ArticleTag) TableName() string {
	return "m_article_tag"
}
