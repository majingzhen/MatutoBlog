package main

import (
	"log"
	"matuto-blog/config"
	"matuto-blog/internal/api/router"
	database2 "matuto-blog/internal/database"
	"matuto-blog/pkg/logger"
	"matuto-blog/pkg/storage"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化日志
	logger.Init()

	// 初始化数据库
	if err := database2.Init(); err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		log.Println("Continuing without database connection...")
	} else {
		// 初始化数据库表
		if err := database2.InitTables(database2.GetDB()); err != nil {
			log.Printf("Warning: Failed to initialize database tables: %v", err)
		}

		// 创建数据库索引
		if err := database2.CreateIndexes(database2.GetDB()); err != nil {
			log.Printf("Warning: Failed to create database indexes: %v", err)
		}
	}

	// 初始化存储系统
	if err := storage.InitStorage(); err != nil {
		log.Printf("Warning: Failed to initialize storage: %v", err)
	}

	// 初始化路由
	r := router.Init()

	// 启动服务器
	port := config.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting on port " + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
