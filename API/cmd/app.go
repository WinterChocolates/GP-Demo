package cmd

import (
	"context"
	"log"

	"API/config"
	"API/routes"
	"API/services"
	"API/storage/cache"
	"API/storage/database"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Start 是应用程序的入口点，包含所有初始化逻辑和服务器启动
func Start() {
	// 初始化应用
	_, stop := context.WithCancel(context.Background())
	defer stop()

	// 初始化各组件
	db := initDatabase()
	redisClient := initRedis()
	logger := initLogger()
	defer logger.Sync()

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

	// 创建增强版路由
	router := routes.SetupRouter(userService, jobService)

	// 启动服务器
	log.Println("🚀 启动服务器...")
	StartServer(router)
}

// initConfig 初始化配置
func initConfig() {
	log.Println("🚀 开始初始化配置...")
	if err := config.InitConfig(); err != nil {
		log.Fatalf("❌ 配置初始化失败: %v", err)
	}
}

// initDatabase 初始化数据库
func initDatabase() *gorm.DB {
	// 先初始化配置
	initConfig()

	log.Println("🚀 开始初始化MySQL...")
	db, err := database.InitMySQL()
	if err != nil {
		log.Fatalf("❌ 数据库初始化失败: %v", err)
	}
	return db
}

// initRedis 初始化Redis
func initRedis() *redis.Client {
	log.Println("🚀 开始初始化Redis...")
	redisClient, err := cache.InitRedis()
	if err != nil {
		log.Fatalf("❌ Redis初始化失败: %v", err)
	}
	return redisClient
}

// initLogger 初始化日志
func initLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	return logger
}
