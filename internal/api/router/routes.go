package router

import (
	"matuto-blog/config"
	"matuto-blog/internal/api/controllers"
	"matuto-blog/internal/api/middlewares"
	"matuto-blog/pkg/utils"

	"github.com/gin-gonic/gin"
)

// InitRoutes 初始化路由
func InitRoutes() *gin.Engine {
	r := gin.New()

	// 使用中间件
	r.Use(middlewares.Logger())
	r.Use(middlewares.CORS())

	// 设置模板路径 - 根据主题配置加载模板
	themePath := config.GetString("theme.path")
	templateNames := []string{
		"default",
		"theme2",
	}
	// 设置模板方法
	customFuncs := utils.GenTemplateFuncMap()
	// 加载模板
	tplManager := utils.NewTemplateManager(themePath, templateNames)
	tplManager.LoadTemplates(r, customFuncs)

	// 静态文件
	r.Static("/static", "./web/static")
	r.Static("/uploads", "./web/uploads")

	// 初始化控制器
	authController := &controllers.AuthController{}
	articleController := &controllers.ArticleController{}
	categoryController := &controllers.CategoryController{}
	tagController := &controllers.TagController{}
	commentController := &controllers.CommentController{}
	attachmentController := &controllers.AttachmentController{}

	// 前台路由
	frontend := r.Group("/")
	{
		// 首页
		frontend.GET("", articleController.Index)

		// 文章详情
		frontend.GET("/article/:id", articleController.Show)

		// 分类页面
		frontend.GET("/category/:id", articleController.Index)

		// 分类列表页面
		frontend.GET("/categories", categoryController.CategoryListPage)

		// 标签页面
		frontend.GET("/tag/:id", articleController.Index)

		// 搜索页面
		frontend.GET("/search", articleController.Index)

		// 评论提交
		frontend.POST("/comment/submit", commentController.Submit)
	}

	// API路由 (用于AJAX请求)
	api := r.Group("/api")
	{
		// 认证接口
		api.POST("/login", authController.Login)
		api.POST("/logout", authController.Logout)
		// 需要认证的API
		apiAuth := api.Group("", middlewares.JWTAuth())
		{
			apiAuth.GET("/profile", authController.GetProfile)
			// 文件管理
			attr := apiAuth.Group("/attach")
			{
				attr.POST("/upload", attachmentController.Upload)
				attr.GET("/page", attachmentController.AttachPage)
				attr.DELETE("/:id", attachmentController.DeleteAttach)
			}
			// 文章管理
			articles := apiAuth.Group("/articles")
			{
				articles.GET("/page", articleController.ArticlePage)
				articles.GET("/:id", articleController.GetArticle)
				articles.DELETE("/:id", articleController.DeleteArticle)
				articles.POST("/publish", articleController.PublishArticle)
				articles.PUT("/update", articleController.UpdateArticle)
			}
			// 分类管理
			categories := apiAuth.Group("/categories")
			{
				categories.GET("/page", categoryController.CategoryPage)
				categories.POST("", categoryController.CreateCategory)
				categories.PUT("/:id", categoryController.UpdateCategory)
				categories.DELETE("/:id", categoryController.DeleteCategory)
				categories.GET("/enable-list", categoryController.CategoryEnableList)
			}
			// 标签管理
			tags := apiAuth.Group("/tags")
			{
				tags.GET("/page", tagController.TagPage)
				tags.POST("/create", tagController.CreateTag)
				tags.PUT("/:id", tagController.UpdateTag)
				tags.POST("/:id", tagController.UpdateTag)
				tags.DELETE("/:id", tagController.DeleteTag)
				tags.GET("/enable-list", tagController.TagEnableList)
			}
			// 评论管理
			comments := apiAuth.Group("/comments")
			{
				comments.GET("/page", commentController.CommentPage)
				comments.PUT("/:id/status", commentController.ReviewComment)
				comments.DELETE("/:id", commentController.DestroyComment)
				comments.POST("/batch-review", commentController.BatchReviewComment)
			}
		}

	}

	return r
}
