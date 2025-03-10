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
// @Summary 上传文件
// @Description 上传文件到服务器
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "要上传的文件"
// @Success 200 {object} utils.Response{data=map[string]string{filename=string}} "文件上传成功"
// @Failure 400 {object} utils.Response "文件上传失败"
// @Failure 500 {object} utils.Response "文件保存失败"
// @Router /api/v1/upload [post]
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
// @Summary 下载文件
// @Description 下载指定的文件
// @Tags 文件管理
// @Produce octet-stream
// @Param file_id path string true "文件ID"
// @Success 200 {file} binary "文件内容"
// @Router /api/v1/download/{file_id} [get]
func (ctl *UploadController) DownloadFile(c *gin.Context) {
	fileID := c.Param("file_id")
	c.FileAttachment("uploads/"+fileID, fileID)
}
