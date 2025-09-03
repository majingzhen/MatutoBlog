package database

import (
	"matuto-blog/internal/models"
	"matuto-blog/pkg/logger"

	"gorm.io/gorm"
)

// InitTables 初始化数据库表（使用GORM AutoMigrate）
func InitTables(db *gorm.DB) error {
	logger.Info("Initializing database tables using GORM AutoMigrate...")

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Category{},
		&models.Tag{},
		&models.Comment{},
		&models.Attach{},
		&models.Link{},
		&models.ArticleCategory{},
		&models.ArticleTag{},
	)

	if err != nil {
		logger.Error("Failed to migrate database tables:", err)
		return err
	}

	logger.Info("Database tables initialized successfully")
	return nil
}

// CreateIndexes 创建必要的索引
func CreateIndexes(db *gorm.DB) error {
	logger.Info("Creating database indexes...")

	// 文章表索引
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_article_status ON p_article(status)",
		"CREATE INDEX IF NOT EXISTS idx_article_type ON p_article(type)",
		"CREATE INDEX IF NOT EXISTS idx_article_is_top ON p_article(is_top)",
		"CREATE INDEX IF NOT EXISTS idx_article_slug ON p_article(slug)",
		"CREATE INDEX IF NOT EXISTS idx_article_create_user_id ON p_article(create_user_id)",

		// 分类表索引
		"CREATE INDEX IF NOT EXISTS idx_category_status ON p_category(status)",
		"CREATE INDEX IF NOT EXISTS idx_category_slug ON p_category(slug)",
		"CREATE INDEX IF NOT EXISTS idx_category_pid ON p_category(pid)",

		// 标签表索引
		"CREATE INDEX IF NOT EXISTS idx_tag_slug ON p_tag(slug)",

		// 评论表索引
		"CREATE INDEX IF NOT EXISTS idx_comment_article_id ON p_comment(article_id)",
		"CREATE INDEX IF NOT EXISTS idx_comment_status ON p_comment(status)",
		"CREATE INDEX IF NOT EXISTS idx_comment_pid ON p_comment(pid)",

		// 附件表索引
		"CREATE INDEX IF NOT EXISTS idx_attach_type ON p_attach(type)",

		// 文章分类关联表索引
		"CREATE INDEX IF NOT EXISTS idx_article_category_article_id ON p_article_category(article_id)",
		"CREATE INDEX IF NOT EXISTS idx_article_category_category_id ON p_article_category(category_id)",

		// 文章标签关联表索引
		"CREATE INDEX IF NOT EXISTS idx_article_tag_article_id ON p_article_tag(article_id)",
		"CREATE INDEX IF NOT EXISTS idx_article_tag_tag_id ON p_article_tag(tag_id)",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			logger.Warn("Failed to create index:", err)
			// 继续执行其他索引，不中断流程
		}
	}

	logger.Info("Database indexes created successfully")
	return nil
}