package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"API/controllers"
	"API/middlewares"
	"API/services"
	"API/storage/cache"
	"API/storage/database"
	"API/utils"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
	//_ "API/docs"
)

func composeMiddlewares(middlewares ...gin.HandlerFunc) []gin.HandlerFunc {
	return middlewares
}

var (
	authMiddleware = composeMiddlewares(
		middlewares.JWT(),
		middlewares.AuditLogger(zap.L()),
	)

	adminMiddleware = composeMiddlewares(
		middlewares.JWT(),
		middlewares.AdminOnly(),
		middlewares.OperationLog(zap.L()),
	)
)

func SetupRouter(
	userService *services.UserService,
	jobService *services.JobService,
) *gin.Engine {
	router := gin.New()

	// 初始化所有服务
	attendanceService := services.NewAttendanceService(database.DB)
	trainingService := services.NewTrainingService(database.DB)
	salaryService := services.NewSalaryService(database.DB)
	noticeService := services.NewNoticeService(
		database.DB,
		cache.NewRedisCacheService(cache.RedisClient),
	)
	resumeService := services.NewResumeService(database.DB)
	permissionService := services.NewPermissionService(database.DB)
	applicationService := services.NewApplicationService(database.DB)
	roleService := services.NewRoleService(database.DB)

	ctrls := struct {
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
	}{
		user:        controllers.NewUserController(userService),
		attendance:  controllers.NewAttendanceController(attendanceService),
		training:    controllers.NewTrainingController(trainingService),
		salary:      controllers.NewSalaryController(salaryService),
		notice:      controllers.NewNoticeController(noticeService),
		job:         controllers.NewJobController(jobService),
		resume:      controllers.NewResumeController(resumeService),
		permission:  controllers.NewPermissionController(permissionService),
		application: controllers.NewApplicationController(applicationService),
		role:        controllers.NewRoleController(roleService),
	}

	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 全局中间件（顺序敏感）
	router.Use(
		middlewares.RequestLogger(zap.L()),
		middlewares.SecurityHeaders(),
		middlewares.RateLimiter(100, time.Minute),
		gin.Recovery(),
	)

	apiV1 := router.Group("/api/v1")
	{
		// 公共路由
		apiV1.GET("/health", healthCheckHandler)
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", ctrls.user.Login)
			authGroup.POST("/register", ctrls.user.Register)
		}

		// 认证路由
		authRoutes := apiV1.Group("").Use(authMiddleware...)
		{
			// 考勤路由
			attendance := apiV1.Group("/attendance")
			{
				attendance.POST("/clock-in", ctrls.attendance.ClockIn)
				attendance.POST("/clock-out", ctrls.attendance.ClockOut)
				attendance.GET("/monthly", ctrls.attendance.GetMonthly)
			}

			resumeGroup := apiV1.Group("/resume")
			{
				resumeGroup.POST("", ctrls.resume.SubmitResume)
				resumeGroup.GET("", ctrls.resume.GetResume)
			}

			// 公共认证路由
			authRoutes.GET("/notices", ctrls.notice.GetNotices)

			authRoutes.GET("/users/profile", ctrls.user.GetProfile)
			authRoutes.PUT("/users/profile", ctrls.user.UpdateProfile)

			authRoutes.POST("/upload", ctrls.upload.UploadFile)
			authRoutes.GET("/download/:file_id", ctrls.upload.DownloadFile)
		}

		// 管理员路由
		adminRoutes := apiV1.Group("").Use(adminMiddleware...)
		{
			// 培训管理
			//adminRoutes.POST("/trainings", ctrls.training.CreateTraining)
			training := apiV1.Group("/training")
			{
				training.GET("", ctrls.training.GetTrainingDetail)
				training.POST("/:id/register", ctrls.training.RegisterTraining)
				training.GET("/my", ctrls.training.GetMyTrainings)
				training.DELETE("/records/:id", ctrls.training.CancelTrainingRegistration)
				training.GET("/:id", ctrls.training.GetTrainingDetail)
			}

			// 薪资管理
			adminRoutes.POST("/salaries/generate", ctrls.salary.GenerateSalary)
			authRoutes.GET("/salaries/history", ctrls.salary.GetSalaryHistory)

			// 通知管理
			adminRoutes.POST("/notices", ctrls.notice.CreateNotice)
			adminRoutes.DELETE("/notices/:id", ctrls.notice.DeleteNotice)
			adminRoutes.PUT("/notices/:id", ctrls.notice.UpdateNotice)

			adminRoutes.POST("/roles", ctrls.role.CreateRole)
			adminRoutes.GET("/roles", ctrls.role.GetRoles)

			//更新申请职位接口
			adminRoutes.PUT("/applications/:id", ctrls.application.UpdateApplicationStatus)

			//更新培训记录接口
			adminRoutes.PUT("/training/records/:id", ctrls.training.UpdateTrainingRecord)

			adminRoutes.GET("/attendance/stats", ctrls.attendance.GetAttendanceStats)

			// 岗位管理
			jobs := apiV1.Group("/jobs")
			{
				jobs.GET("", ctrls.job.ListJobs)
				jobs.POST("", ctrls.job.CreateJob)
				jobs.PUT("/:id", ctrls.job.UpdateJob)
				jobs.POST("/:id/apply", ctrls.job.ApplyForJob)
				adminRoutes.DELETE("/:id", ctrls.job.DeleteJob)
			}

			permission := apiV1.Group("/permissions")
			{
				permission.GET("", ctrls.permission.GetPermissions)
				permission.POST("", ctrls.permission.CreatePermission)
			}
		}
	}

	return router
}

func healthCheckHandler(c *gin.Context) {
	type healthCheck func(context.Context) error

	checks := map[string]healthCheck{
		"mysql": database.CheckMySQLHealth,
		"redis": cache.CheckRedisHealth,
	}

	result := make(map[string]string)
	status := http.StatusOK
	logger := zap.L()

	for name, check := range checks {
		if err := check(c.Request.Context()); err != nil {
			logger.Error(fmt.Sprintf("%s健康检查失败", name), zap.Error(err))
			result[name] = "unhealthy"
			status = http.StatusServiceUnavailable
		} else {
			result[name] = "healthy"
		}
	}

	logger.Info("健康检查完成",
		zap.Int("status", status),
		zap.Any("services", result),
	)

	c.JSON(status, gin.H{
		"status":   utils.MapStatus(status),
		"services": result,
		"time":     time.Now().UTC().Format(time.RFC3339),
	})
}

// API版本分组
//apiV1 := router.Group("/api/v1")
//{
//	// 公共路由
//	public := apiV1.Group("")
//	{
//		userCtrl := controllers.NewUserController(userService)
//		auth := apiV1.Group("/auth")
//		{
//			// 登录路由
//			auth.POST("/login", userCtrl.Login)
//			// 注册路由
//			auth.POST("/register", userCtrl.Register)
//		}
//		// 数据库检查
//		public.GET("/health", healthCheckHandler)
//	}
//
//	// 认证路由
//	auth := apiV1.Group("").Use(
//		middlewares.JWT(),
//		middlewares.AuditLogger(zap.L()),
//	)
//	{
//		// 考勤路由组
//		attendanceCtrl := controllers.NewAttendanceController(attendanceService)
//		attendance := apiV1.Group("/attendance")
//		{
//			attendance.POST("/clock-in", attendanceCtrl.ClockIn)
//			attendance.POST("/clock-out", attendanceCtrl.ClockOut)
//			attendance.GET("/monthly", attendanceCtrl.GetMonthly)
//		}
//
//		// 培训路由组
//		trainingCtrl := controllers.NewTrainingController(trainingService)
//		training := apiV1.Group("/training")
//		{
//			training.GET("", trainingCtrl.GetTrainings)
//			training.POST("/:id/register", trainingCtrl.RegisterTraining)
//			training.GET("/my", trainingCtrl.GetMyTrainings)
//		}
//
//		// 薪资路由组
//		salaryCtrl := controllers.NewSalaryController(salaryService)
//		salary := apiV1.Group("/salary")
//		{
//			salary.GET("/:month", salaryCtrl.GetSalaryDetails)
//		}
//
//		// 通知路由组（对所有认证用户开放）
//		noticeCtrl := controllers.NewNoticeController(noticeService)
//		auth.GET("/notices", noticeCtrl.GetNotices)
//	}
//
//	// 管理员路由
//	admin := apiV1.Group("").Use(
//		middlewares.JWT(),
//		middlewares.AdminOnly(),
//		middlewares.OperationLog(zap.L()),
//	)
//	{
//		// 培训管理
//		trainingCtrl := controllers.NewTrainingController(trainingService)
//		admin.POST("/trainings", trainingCtrl.CreateTraining)
//
//		// 薪资管理
//		salaryCtrl := controllers.NewSalaryController(salaryService)
//		admin.POST("/salaries/generate", salaryCtrl.GenerateSalary)
//
//		// 通知管理
//		noticeCtrl := controllers.NewNoticeController(noticeService)
//		admin.POST("/notices", noticeCtrl.CreateNotice)
//
//		// 岗位管理
//		jobCtrl := controllers.NewJobController(jobService)
//		admin.GET("/jobs", jobCtrl.ListJobs)
//		admin.POST("/jobs", jobCtrl.CreateJob)
//		admin.PUT("/jobs/:id", jobCtrl.UpdateJob)
//	}
//}
