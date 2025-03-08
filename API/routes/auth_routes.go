package routes

import (
	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(apiV1 *gin.RouterGroup, ctrls Controllers) {
	authGroup := apiV1.Group("/auth")
	{
		authGroup.POST("/login", ctrls.user.Login)
		authGroup.POST("/register", ctrls.user.Register)
	}

	authRoutes := apiV1.Group("").Use(defaultAuthMiddleware...)
	{
		authRoutes.GET("/notices", ctrls.notice.GetNotices)
		authRoutes.GET("/notices/department/:department", ctrls.notice.GetDepartmentNotices)
		authRoutes.PUT("/notices/:id/read", ctrls.notice.MarkNoticeAsRead)

		users := apiV1.Group("/users")
		{
			users.GET("/profile", ctrls.user.GetProfile)
			users.PUT("/profile", ctrls.user.UpdateProfile)
		}

		authRoutes.POST("/upload", ctrls.upload.UploadFile)
		authRoutes.GET("/download/:file_id", ctrls.upload.DownloadFile)
	}
}
