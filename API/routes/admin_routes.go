package routes

import (
	"github.com/gin-gonic/gin"
)

func setupAdminRoutes(apiV1 *gin.RouterGroup, ctrls Controllers) {
	apiV1.Group("").Use(adminAuthMiddleware...)
	{
		// 薪资管理
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
