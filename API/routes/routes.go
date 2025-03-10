package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"API/controllers"
	"API/docs"
	"API/middlewares"
	"API/services"
	"API/storage/cache"
	"API/storage/database"
	"API/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

var (
	defaultAuthMiddleware = []gin.HandlerFunc{
		middlewares.JWT(),
		middlewares.EnhancedAuditLogger(zap.L()),
	}

	adminAuthMiddleware = []gin.HandlerFunc{
		middlewares.JWT(),
		middlewares.AdminOnly(),
		middlewares.EnhancedAuditLogger(zap.L()),
	}
)

type Controllers struct {
	user        *controllers.UserController
	attendance  *controllers.AttendanceController
	training    *controllers.TrainingController
	salary      *controllers.SalaryController
	notice      *controllers.NoticeController
	job         *controllers.JobController
	resume      *controllers.ResumeController
	permission  *controllers.PermissionController
	application *controllers.ApplicationController
	role        *controllers.RoleController
	upload      *controllers.UploadController
}

// initSwagger 初始化Swagger文档
func initSwagger() {
	docs.SwaggerInfo.Title = "招聘系统API文档"
	docs.SwaggerInfo.Description = "招聘系统后端API接口文档"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func SetupRouter(userService *services.UserService, jobService *services.JobService) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化Swagger文档
	initSwagger()

	// 创建路由引擎
	router := gin.New()

	// 初始化控制器
	ctrls := Controllers{
		user:        controllers.NewUserController(userService),
		attendance:  controllers.NewAttendanceController(services.NewAttendanceService(database.DB)),
		training:    controllers.NewTrainingController(services.NewTrainingService(database.DB)),
		salary:      controllers.NewSalaryController(services.NewSalaryService(database.DB)),
		notice:      controllers.NewNoticeController(services.NewNoticeService(database.DB, cache.NewRedisCacheService(cache.RedisClient))),
		job:         controllers.NewJobController(jobService),
		resume:      controllers.NewResumeController(services.NewResumeService(database.DB, cache.NewRedisCacheService(cache.RedisClient))),
		permission:  controllers.NewPermissionController(services.NewPermissionService(database.DB)),
		application: controllers.NewApplicationController(services.NewApplicationService(database.DB)),
		role:        controllers.NewRoleController(services.NewRoleService(database.DB)),
		upload:      controllers.NewUploadController(),
	}

	// 配置Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 全局中间件
	router.Use(
		middlewares.RequestContext(),
		middlewares.EnhancedLogger(zap.L()),
		middlewares.ErrorLogger(zap.L()),
		middlewares.EnhancedSecurityHeaders(),
		middlewares.EnhancedCORSMiddleware(),
		middlewares.EnhancedRateLimiter(zap.L()),
		middlewares.PerformanceMonitor(zap.L()),
		gin.Recovery(),
	)

	// API路由组
	apiV1 := router.Group("/api/v1")
	{
		// 健康检查
		apiV1.GET("/health", enhancedHealthCheckHandler)

		setupAuthRoutes(apiV1, ctrls)
		setupAdminRoutes(apiV1, ctrls)
	}

	// 404处理
	router.NoRoute(func(c *gin.Context) {
		utils.RespondWithJSON(c, http.StatusNotFound, 404, "请求的资源不存在", nil)
	})

	// 405处理
	router.NoMethod(func(c *gin.Context) {
		utils.RespondWithJSON(c, http.StatusMethodNotAllowed, 405, "不支持的请求方法", nil)
	})

	return router
}

func enhancedHealthCheckHandler(c *gin.Context) {
	type healthCheck func(context.Context) error

	checks := map[string]healthCheck{
		"database": database.CheckMySQLHealth,
		"redis":    cache.CheckRedisHealth,
	}

	results := make(map[string]string)
	overallStatus := http.StatusOK

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	for name, check := range checks {
		if err := check(ctx); err != nil {
			results[name] = fmt.Sprintf("异常: %v", err)
			overallStatus = http.StatusServiceUnavailable
		} else {
			results[name] = "正常"
		}
	}

	c.JSON(overallStatus, gin.H{
		"status":  utils.MapStatusToString(overallStatus),
		"checks":  results,
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
	})
}
