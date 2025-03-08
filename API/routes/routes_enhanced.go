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

// EnhancedControllers 包含所有控制器实例
type EnhancedControllers struct {
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

func SetupEnhancedRouter(userService *services.UserService, jobService *services.JobService) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由引擎
	router := gin.New()

	// 初始化控制器
	ctrls := EnhancedControllers{
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
		upload:      controllers.NewUploadController(),
	}

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

		// 认证路由
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", ctrls.user.Login)
			authGroup.POST("/register", ctrls.user.Register)
		}

		// 设置需要认证的路由
		setupEnhancedAuthRoutes(apiV1, ctrls)

		// 设置需要管理员权限的路由
		setupEnhancedAdminRoutes(apiV1, ctrls)
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

func setupEnhancedAuthRoutes(apiV1 *gin.RouterGroup, ctrls EnhancedControllers) {
	authRoutes := apiV1.Group("").Use(authMiddleware...)
	{
		// 通知相关
		notices := apiV1.Group("/notices", defaultAuthMiddleware...)
		{
			notices.GET("", ctrls.notice.GetNotices)
			notices.GET("/department/:department", ctrls.notice.GetDepartmentNotices)
			notices.PUT("/:id/read", ctrls.notice.MarkNoticeAsRead)
		}

		// 用户相关
		users := apiV1.Group("/users")
		{
			users.GET("/profile", ctrls.user.GetProfile)
			users.PUT("/profile", ctrls.user.UpdateProfile)
		}

		// 文件上传下载
		authRoutes.POST("/upload", ctrls.upload.UploadFile)
		authRoutes.GET("/download/:file_id", ctrls.upload.DownloadFile)

		// 考勤相关
		attendance := apiV1.Group("/attendance")
		{
			attendance.POST("/clock-in", ctrls.attendance.ClockIn)
			attendance.POST("/clock-out", ctrls.attendance.ClockOut)
			attendance.GET("/monthly", ctrls.attendance.GetMonthly)
		}

		// 简历相关
		resumeGroup := apiV1.Group("/resume")
		{
			resumeGroup.POST("", ctrls.resume.SubmitResume)
			resumeGroup.GET("", ctrls.resume.GetResume)
		}
	}
}

func setupEnhancedAdminRoutes(apiV1 *gin.RouterGroup, ctrls EnhancedControllers) {
	adminRoutes := apiV1.Group("").Use(adminMiddleware...)
	{
		// 薪资相关
		salaries := apiV1.Group("/salaries")
		{
			salaries.POST("/generate", ctrls.salary.GenerateSalary)
			salaries.GET("/history", ctrls.salary.GetSalaryHistory)
		}

		// 通知管理
		notices := apiV1.Group("/notices")
		{
			notices.POST("", ctrls.notice.CreateNotice)
			notices.DELETE("/:id", ctrls.notice.DeleteNotice)
			notices.PUT("/:id", ctrls.notice.UpdateNotice)
		}

		// 角色管理
		roles := apiV1.Group("/roles")
		{
			roles.POST("", ctrls.role.CreateRole)
			roles.GET("", ctrls.role.GetRoles)
		}

		// 申请管理
		adminRoutes.PUT("/applications/:id", ctrls.application.UpdateApplicationStatus)

		// 培训管理
		adminRoutes.PUT("/training/records/:id", ctrls.training.UpdateTrainingRecord)

		// 考勤统计
		adminRoutes.GET("/attendance/stats", ctrls.attendance.GetAttendanceStats)

		// 职位管理
		jobs := apiV1.Group("/jobs")
		{
			jobs.GET("", ctrls.job.ListJobs)
			jobs.POST("", ctrls.job.CreateJob)
			jobs.PUT("/:id", ctrls.job.UpdateJob)
			jobs.POST("/:id/apply", ctrls.job.ApplyForJob)
			jobs.DELETE("/:id", ctrls.job.DeleteJob)
		}

		// 权限管理
		permission := apiV1.Group("/permissions")
		{
			permission.GET("", ctrls.permission.GetPermissions)
			permission.POST("", ctrls.permission.CreatePermission)
		}
	}
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
