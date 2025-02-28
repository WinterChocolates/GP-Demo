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
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	authMiddleware = []gin.HandlerFunc{
		middlewares.JWT(),
		middlewares.AuditLogger(zap.L()),
	}

	adminMiddleware = []gin.HandlerFunc{
		middlewares.JWT(),
		middlewares.AdminOnly(),
		middlewares.OperationLog(zap.L()),
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

func setupAuthRoutes(apiV1 *gin.RouterGroup, ctrls Controllers) {
	authRoutes := apiV1.Group("").Use(authMiddleware...)
	{
		authRoutes.GET("/notices", ctrls.notice.GetNotices)
		authRoutes.GET("/users/profile", ctrls.user.GetProfile)
		authRoutes.PUT("/users/profile", ctrls.user.UpdateProfile)
		authRoutes.POST("/upload", ctrls.upload.UploadFile)
		authRoutes.GET("/download/:file_id", ctrls.upload.DownloadFile)

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
	}
}

func setupAdminRoutes(apiV1 *gin.RouterGroup, ctrls Controllers) {
	adminRoutes := apiV1.Group("").Use(adminMiddleware...)
	{
		adminRoutes.POST("/salaries/generate", ctrls.salary.GenerateSalary)
		adminRoutes.GET("/salaries/history", ctrls.salary.GetSalaryHistory)

		adminRoutes.POST("/notices", ctrls.notice.CreateNotice)
		adminRoutes.DELETE("/notices/:id", ctrls.notice.DeleteNotice)
		adminRoutes.PUT("/notices/:id", ctrls.notice.UpdateNotice)

		adminRoutes.POST("/roles", ctrls.role.CreateRole)
		adminRoutes.GET("/roles", ctrls.role.GetRoles)

		adminRoutes.PUT("/applications/:id", ctrls.application.UpdateApplicationStatus)
		adminRoutes.PUT("/training/records/:id", ctrls.training.UpdateTrainingRecord)
		adminRoutes.GET("/attendance/stats", ctrls.attendance.GetAttendanceStats)

		jobs := apiV1.Group("/jobs")
		{
			jobs.GET("", ctrls.job.ListJobs)
			jobs.POST("", ctrls.job.CreateJob)
			jobs.PUT("/:id", ctrls.job.UpdateJob)
			jobs.POST("/:id/apply", ctrls.job.ApplyForJob)
			jobs.DELETE("/:id", ctrls.job.DeleteJob)
		}

		permission := apiV1.Group("/permissions")
		{
			permission.GET("", ctrls.permission.GetPermissions)
			permission.POST("", ctrls.permission.CreatePermission)
		}
	}
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
