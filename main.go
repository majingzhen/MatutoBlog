package main

import (
	"matuto-blog/config"
	"matuto-blog/internal/api/router"
	database2 "matuto-blog/internal/database"
	"matuto-blog/pkg/logger"
	"matuto-blog/pkg/storage"
	"matuto-blog/pkg/utils"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化日志
	logger.Init()

	// 初始化数据库
	if err := database2.Init(); err != nil {
		logger.Error("Warning: Failed to initialize database: %v", err)
		logger.Error("Continuing without database connection...")
		return
	} else {
		// 初始化数据库表
		if err := database2.InitTables(database2.GetDB()); err != nil {
			logger.Error("Warning: Failed to initialize database tables: %v", err)
		}

		// 创建数据库索引
		if err := database2.CreateIndexes(database2.GetDB()); err != nil {
			logger.Error("Warning: Failed to create database indexes: %v", err)
		}
	}

	// 初始化存储系统
	if err := storage.InitStorage(); err != nil {
		logger.Error("Warning: Failed to initialize storage: %v", err)
	}

	// 初始化路由
	r := router.InitRoutes()

	// 启动服务器
	port := config.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting on port " + port)

	password, err := utils.HashDefaultConfigPassword("123456")
	if err != nil {
		logger.Fatal("Failed to hash password:", err)
	}
	logger.Info("Password:", password)

	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server:", err)
	}

}
