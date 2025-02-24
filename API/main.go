package main

import (
	"context"
	"log"

	"API/cmd"
	"API/config"
	"API/routes"
	"API/services"
	"API/storage/cache"
	"API/storage/database"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("❌ 配置初始化失败: %v", err)
	}

	// 初始化数据库
	db, err := database.InitMySQL()
	if err != nil {
		log.Fatalf("❌ 数据库初始化失败: %v", err)
	}

	// 初始化Redis
	redisClient, err := cache.InitRedis()
	if err != nil {
		log.Fatalf("❌ Redis初始化失败: %v", err)
	}

	_, stop := context.WithCancel(context.Background())
	defer stop()

	// 使用闭包处理资源释放
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("⚠️ 关闭数据库错误: %v", err)
		}
		if err := cache.Close(); err != nil {
			log.Printf("⚠️ 关闭Redis错误: %v", err)
		}
	}()

	// 初始化服务层
	cacheService := cache.NewRedisCacheService(redisClient)
	userService := services.NewUserService(db, cacheService)
	jobService := services.NewJobService(db)

	// 创建路由
	router := routes.SetupRouter(userService, jobService)

	// 启动服务器
	cmd.StartServer(router)
}
