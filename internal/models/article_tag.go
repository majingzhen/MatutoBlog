package models

// ArticleTag 文章标签关联模型
type ArticleTag struct {
	BaseModel
	ArticleId uint64 `json:"article_id" gorm:"not null;comment:文章id"`
	TagId     uint64 `json:"tag_id" gorm:"not null;comment:标签id"`
}

// TableName 指定表名
func (ArticleTag) TableName() string {
	return "m_article_tag"
}
