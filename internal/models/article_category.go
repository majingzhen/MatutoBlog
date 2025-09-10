package models

// ArticleCategory 文章分类关联模型
type ArticleCategory struct {
	BaseModel
	ArticleId  int `json:"article_id" gorm:"not null;comment:文章id"`
	CategoryId int `json:"category_id" gorm:"not null;comment:分类id"`
}

// TableName 指定表名
func (ArticleCategory) TableName() string {
	return "m_article_category"
}
