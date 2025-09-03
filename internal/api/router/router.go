package router

import (
	"context"
	"matuto-blog/config"
	middleware2 "matuto-blog/internal/api/middleware"
	api2 "matuto-blog/internal/api/router/api"
	database2 "matuto-blog/internal/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(config.GetString("server.mode"))

	r := gin.New()

	// 添加中间件
	r.Use(middleware2.TraceID())      // 追踪ID中间件
	r.Use(middleware2.Logger())       // 日志中间件
	r.Use(middleware2.ErrorHandler()) // 错误处理中间件
	r.Use(middleware2.CORS())         // 跨域中间件

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		health := gin.H{
			"status":    "ok",
			"message":   "Blog system is running",
			"timestamp": time.Now().Format(time.RFC3339),
		}

		// 检查数据库连接状态
		if db := database2.GetDB(); db != nil {
			sqlDB, err := db.DB()
			if err == nil {
				if err := sqlDB.Ping(); err == nil {
					health["database"] = "connected"
				} else {
					health["database"] = "disconnected"
				}
			} else {
				health["database"] = "unavailable"
			}
		} else {
			health["database"] = "not_initialized"
		}

		c.JSON(http.StatusOK, health)
	})

	// API路由组
	api := r.Group("/api")
	{
		// 注册各模块路由
		api2.RegisterAuthRoutes(api)   // 认证模块路由
		api2.RegisterUserRoutes(api)   // 用户模块路由
		api2.RegisterCommonRoutes(api) // 通用模块路由
	}

	// 前端页面路由
	r.LoadHTMLGlob("web/templates/*")
	r.Static("/static", "web/static")
	
	// 首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "博客首页",
		})
	})

	return r
}
