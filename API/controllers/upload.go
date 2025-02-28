package controllers

import (
	"net/http"

	"API/utils"
	"github.com/gin-gonic/gin"
)

// UploadController 文件上传控制器
type UploadController struct{}

// NewUploadController 初始化文件上传控制器
func NewUploadController() *UploadController {
	return &UploadController{}
}

// UploadFile 上传文件
func (ctl *UploadController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "文件上传失败")
		return
	}
	if err := c.SaveUploadedFile(file, "uploads/"+file.Filename); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "文件保存失败")
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "文件上传成功", "filename": file.Filename})
}

// DownloadFile 下载文件
func (ctl *UploadController) DownloadFile(c *gin.Context) {
	fileID := c.Param("file_id")
	c.FileAttachment("uploads/"+fileID, fileID)
}
