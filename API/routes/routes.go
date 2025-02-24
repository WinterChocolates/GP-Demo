package routes

import (
	"net/http"
	"time"

	"API/controllers"
	"API/middlewares"
	"API/utils"

	"API/services"
	"API/storage/cache"
	"API/storage/database"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userService *services.UserService,
	jobService *services.JobService,
) *gin.Engine {
	router := gin.New()

	// 全局中间件（顺序敏感）
	router.Use(
		middlewares.RequestLogger(zap.L()),
		middlewares.SecurityHeaders(),
		middlewares.RateLimiter(100, time.Minute),
		gin.Recovery(),
	)

	// API版本分组
	apiV1 := router.Group("/api/v1")
	{
		// 公共路由
		public := apiV1.Group("")
		{
			userCtrl := controllers.NewUserController(userService)
			public.POST("/auth/login", userCtrl.Login)
			public.POST("/auth/register", userCtrl.Register)

			public.GET("/health", healthCheckHandler)
		}

		// 认证路由
		auth := apiV1.Group("").Use(
			middlewares.JWT(),
			middlewares.AuditLogger(zap.L()),
		)
		{
			userCtrl := controllers.NewUserController(userService)
			auth.GET("/users/me", userCtrl.GetProfile)
			auth.PUT("/users/me", userCtrl.UpdateProfile)
		}

		// 管理员路由
		admin := apiV1.Group("").Use(
			middlewares.AdminOnly(),
			middlewares.OperationLog(zap.L()),
		)
		{
			jobCtrl := controllers.NewJobController(jobService)
			admin.GET("/jobs", jobCtrl.ListJobs)
			admin.POST("/jobs", jobCtrl.CreateJob)
			admin.PUT("/jobs/:id", jobCtrl.UpdateJob)
		}
	}

	return router
}

func healthCheckHandler(c *gin.Context) {
	logger := zap.L()
	ctx := c.Request.Context()

	healthStatus := make(map[string]string)
	status := http.StatusOK

	// 数据库检查
	if err := database.CheckMySQLHealth(ctx); err != nil {
		logger.Error("MySQL健康检查失败", zap.Error(err))
		healthStatus["mysql"] = "unhealthy"
		status = http.StatusServiceUnavailable
	} else {
		healthStatus["mysql"] = "healthy"
	}

	// Redis检查
	if err := cache.CheckRedisHealth(ctx); err != nil {
		logger.Error("Redis健康检查失败", zap.Error(err))
		healthStatus["redis"] = "unhealthy"
		status = http.StatusServiceUnavailable
	} else {
		healthStatus["redis"] = "healthy"
	}

	// 记录检查结果
	logger.Info("健康检查完成",
		zap.Int("status", status),
		zap.Any("services", healthStatus),
	)

	c.JSON(status, gin.H{
		"status":   utils.MapStatus(status),
		"services": healthStatus,
		"time":     time.Now().UTC().Format(time.RFC3339),
	})
}
