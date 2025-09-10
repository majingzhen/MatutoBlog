package models

// Article 文章模型
type Article struct {
	BaseModel
	Title           string `json:"title" gorm:"size:256;comment:文章标题"`
	Content         string `json:"content" gorm:"type:longtext;not null;comment:文章内容"`
	ParseContent    string `json:"parseContent" gorm:"type:longtext;not null;comment:解析后的文章内容"`
	ContentModel    string `json:"contentModel" gorm:"size:32;comment:文章内容类型:html/markdown"`
	Type            string `json:"type" gorm:"size:32;comment:文章类型:article文章,page页面"`
	Summary         string `json:"summary" gorm:"size:1024;comment:文章摘要"`
	MetaKeywords    string `json:"metaKeywords" gorm:"size:512;comment:SEO关键字"`
	MetaDescription string `json:"metaDescription" gorm:"size:512;comment:SEO描述"`
	Thumbnail       string `json:"thumbnail" gorm:"size:256;comment:缩略图"`
	Slug            string `json:"slug" gorm:"size:128;index;comment:slug"`
	IsTop           int8   `json:"isTop" gorm:"default:0;comment:是否置顶0:否,1:是"`
	Status          int8   `json:"status" gorm:"default:0;comment:状态0:已发布,1:草稿"`
	ViewCount       int    `json:"viewCount" gorm:"default:0;comment:访问量"`
	GreatCount      int    `json:"greatCount" gorm:"default:0;comment:点赞量"`
	IsComment       int8   `json:"isComment" gorm:"default:1;comment:是否允许评论0:否,1是"`
	Flag            string `json:"flag" gorm:"size:256;comment:标识"`
	Template        string `json:"template" gorm:"size:256;comment:模板"`
	Visibility      int8   `json:"visibility" gorm:"default:0;comment:是否可见, 0是, 1否"`
}

// TableName 指定表名
func (Article) TableName() string {
	return "m_article"
}

// ArticleType 文章类型常量
const (
	ArticleTypeArticle = "article" // 文章
	ArticleTypePage    = "page"    // 页面
)

// ArticleStatus 文章状态常量
const (
	ArticleStatusPublished = 0 // 已发布
	ArticleStatusDraft     = 1 // 草稿
)

// ContentModel 内容模型常量
const (
	ContentModelHTML     = "html"     // HTML
	ContentModelMarkdown = "markdown" // Markdown
)

// IsPublished 检查文章是否已发布
func (a *Article) IsPublished() bool {
	return a.Status == ArticleStatusPublished
}

// IsDraft 检查文章是否为草稿
func (a *Article) IsDraft() bool {
	return a.Status == ArticleStatusDraft
}

// AllowComment 检查文章是否允许评论
func (a *Article) AllowComment() bool {
	return a.IsComment == 1
}

// IsVisible 检查文章是否可见
func (a *Article) IsVisible() bool {
	return a.Visibility == 0
}

// IsTopArticle 检查文章是否置顶
func (a *Article) IsTopArticle() bool {
	return a.IsTop == 1
}
