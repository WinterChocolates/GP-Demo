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
	log.Println("🚀 开始初始化配置...")
	if err := config.InitConfig(); err != nil {
		log.Fatalf("❌ 配置初始化失败: %v", err)
	}

	// 初始化数据库
	log.Println("🚀 开始初始化MySQL...")
	db, err := database.InitMySQL()
	if err != nil {
		log.Fatalf("❌ 数据库初始化失败: %v", err)
	}

	// 初始化Redis
	log.Println("🚀 开始初始化Redis...")
	redisClient, err := cache.InitRedis()
	if err != nil {
		log.Fatalf("❌ Redis初始化失败: %v", err)
	}

	_, stop := context.WithCancel(context.Background())
	defer stop()

	// 资源释放
	defer func() {
		log.Println("🛑 开始关闭资源...")
		if err := database.Close(); err != nil {
			log.Printf("⚠️ 关闭数据库错误: %v", err)
		}
		if err := cache.Close(); err != nil {
			log.Printf("⚠️ 关闭Redis错误: %v", err)
		}
		log.Println("✅ 所有资源已关闭")
	}()

	// 初始化服务层
	cacheService := cache.NewRedisCacheService(redisClient)
	userService := services.NewUserService(db, cacheService)
	jobService := services.NewJobService(db)

	// 创建路由
	router := routes.SetupRouter(userService, jobService)

	// 启动服务器
	log.Println("🚀 启动服务器...")
	cmd.StartServer(router)
}
