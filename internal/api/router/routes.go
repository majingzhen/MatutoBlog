package router

import (
	"github.com/gin-gonic/gin"
	"matuto-blog/internal/api/controllers"
	"matuto-blog/internal/api/middleware"
	"matuto-blog/pkg/utils"
)

// InitRoutes 初始化路由
func InitRoutes() *gin.Engine {
	r := gin.New()

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 设置模板路径 - 加载所有HTML文件
	//r.LoadHTMLGlob("web/templates/*")
	// 如果有子目录，也可以用这种方式
	r.HTMLRender = utils.LoadTemplateFiles("web", ".html")
	// 静态文件
	r.Static("/static", "web/static")
	r.Static("/uploads", "./uploads")

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

		// 标签页面
		frontend.GET("/tag/:id", articleController.Index)

		// 搜索页面
		frontend.GET("/search", articleController.Index)

		// 评论提交
		frontend.POST("/comment/submit", commentController.Submit)
	}

	// 管理员登录路由
	admin := r.Group("/admin")
	{
		admin.GET("/login", authController.AdminLoginPage)
		admin.POST("/login", authController.AdminLogin)
		admin.GET("/logout", authController.AdminLogout)
	}

	// 需要认证的管理后台路由
	adminAuth := r.Group("/admin", middleware.SessionAuth())
	{
		// 后台首页
		adminAuth.GET("", func(c *gin.Context) {
			c.HTML(200, "admin/index.html", gin.H{
				"title": "管理后台",
			})
		})
		adminAuth.GET("/", func(c *gin.Context) {
			c.HTML(200, "admin/index.html", gin.H{
				"title": "管理后台",
			})
		})

		// 文章管理
		articles := adminAuth.Group("/articles")
		{
			articles.GET("", articleController.AdminIndex)
			articles.GET("/", articleController.AdminIndex)
			articles.GET("/create", articleController.AdminCreate)
			articles.POST("/create", articleController.AdminStore)
			articles.GET("/:id/edit", articleController.AdminEdit)
			articles.PUT("/:id", articleController.AdminUpdate)
			articles.POST("/:id", articleController.AdminUpdate) // 兼容表单提交
			articles.DELETE("/:id", articleController.AdminDestroy)
			articles.POST("/:id/delete", articleController.AdminDestroy) // 兼容表单提交
		}

		// 分类管理
		categories := adminAuth.Group("/categories")
		{
			categories.GET("", categoryController.AdminIndex)
			categories.GET("/", categoryController.AdminIndex)
			categories.GET("/create", categoryController.AdminCreate)
			categories.POST("/create", categoryController.AdminStore)
			categories.GET("/:id/edit", categoryController.AdminEdit)
			categories.PUT("/:id", categoryController.AdminUpdate)
			categories.POST("/:id", categoryController.AdminUpdate)
			categories.DELETE("/:id", categoryController.AdminDestroy)
			categories.POST("/:id/delete", categoryController.AdminDestroy)
		}

		// 标签管理
		tags := adminAuth.Group("/tags")
		{
			tags.GET("", tagController.AdminIndex)
			tags.GET("/", tagController.AdminIndex)
			tags.GET("/create", tagController.AdminCreate)
			tags.POST("/create", tagController.AdminStore)
			tags.GET("/:id/edit", tagController.AdminEdit)
			tags.PUT("/:id", tagController.AdminUpdate)
			tags.POST("/:id", tagController.AdminUpdate)
			tags.DELETE("/:id", tagController.AdminDestroy)
			tags.POST("/:id/delete", tagController.AdminDestroy)
		}

		// 评论管理
		comments := adminAuth.Group("/comments")
		{
			comments.GET("", commentController.AdminIndex)
			comments.GET("/", commentController.AdminIndex)
			comments.POST("/:id/review", commentController.AdminReview)
			comments.DELETE("/:id", commentController.AdminDestroy)
			comments.POST("/:id/delete", commentController.AdminDestroy)
			comments.POST("/batch/review", commentController.AdminBatchReview)
		}

		// 附件管理
		attachments := adminAuth.Group("/attachments")
		{
			attachments.GET("", attachmentController.AdminIndex)
			attachments.GET("/", attachmentController.AdminIndex)
			attachments.DELETE("/:id", attachmentController.AdminDestroy)
			attachments.POST("/:id/delete", attachmentController.AdminDestroy)
			attachments.POST("/batch/delete", attachmentController.AdminBatchDelete)
		}

		// 文件上传
		adminAuth.POST("/upload", attachmentController.Upload)
		adminAuth.GET("/upload/token", attachmentController.GetUploadToken)
	}

	// API路由 (用于AJAX请求)
	api := r.Group("/api")
	{
		// 认证接口
		api.POST("/login", authController.Login)

		// 需要认证的API
		apiAuth := api.Group("", middleware.JWTAuth())
		{
			apiAuth.GET("/profile", authController.GetProfile)
		}
	}

	return r
}
