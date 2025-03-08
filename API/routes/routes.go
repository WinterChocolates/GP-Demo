package routes

import (
	"time"

	"API/controllers"
	"API/middlewares"
	"API/services"
	"API/storage/cache"
	"API/storage/database"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "API/docs"
)

// 使用base_routes.go中定义的中间件
var (
	authMiddleware  = defaultAuthMiddleware
	adminMiddleware = adminAuthMiddleware
	// 使用base_routes.go中的enhancedHealthCheckHandler
	healthCheckHandler = enhancedHealthCheckHandler
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

func SetupRouter(userService *services.UserService, jobService *services.JobService) *gin.Engine {
	router := gin.New()

	ctrls := Controllers{
		user:        controllers.NewUserController(userService),
		attendance:  controllers.NewAttendanceController(services.NewAttendanceService(database.DB)),
		training:    controllers.NewTrainingController(services.NewTrainingService(database.DB)),
		salary:      controllers.NewSalaryController(services.NewSalaryService(database.DB)),
		notice:      controllers.NewNoticeController(services.NewNoticeService(database.DB, cache.NewRedisCacheService(cache.RedisClient))),
		job:         controllers.NewJobController(jobService),
		resume:      controllers.NewResumeController(services.NewResumeService(database.DB)),
		permission:  controllers.NewPermissionController(services.NewPermissionService(database.DB)),
		application: controllers.NewApplicationController(services.NewApplicationService(database.DB)),
		role:        controllers.NewRoleController(services.NewRoleService(database.DB)),
	}

	// 配置Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 全局中间件
	router.Use(
		middlewares.RequestLogger(zap.L()),
		middlewares.SecurityHeaders(),
		middlewares.RateLimiter(100, time.Minute),
		gin.Recovery(),
	)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/health", healthCheckHandler)

		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", ctrls.user.Login)
			authGroup.POST("/register", ctrls.user.Register)
		}

		setupAuthRoutes(apiV1, ctrls)
		setupAdminRoutes(apiV1, ctrls)
	}

	return router
}
