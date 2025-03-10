package controllers

import (
	"net/http"

	"API/models"
	"API/services"
	"API/utils"

	"github.com/gin-gonic/gin"
)

type ResumeController struct {
	resumeService *services.ResumeService
}

func NewResumeController(rs *services.ResumeService) *ResumeController {
	return &ResumeController{resumeService: rs}
}

// SubmitResume 提交简历信息
// @Summary 提交简历信息
// @Description 提交或更新用户的简历信息
// @Tags 简历管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param resume body models.Resume true "简历信息"
// @Success 200 {object} utils.Response{message=string}
// @Failure 400 {object} utils.Response "无效的请求参数"
// @Failure 401 {object} utils.Response "未授权的请求"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/resumes [post]
func (ctl *ResumeController) SubmitResume(c *gin.Context) {
	var resume models.Resume
	if err := c.ShouldBindJSON(&resume); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "未授权的请求")
		return
	}

	resume.UserID = userID.(uint)

	if err := ctl.resumeService.SubmitResume(c.Request.Context(), resume.UserID, &resume); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"message": "简历提交成功"})
}

// GetResume 获取用户简历
// @Summary 获取用户简历
// @Description 获取当前用户的简历信息
// @Tags 简历管理
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=models.Resume}
// @Failure 401 {object} utils.Response "未授权的请求"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/resumes [get]
func (ctl *ResumeController) GetResume(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "未授权的请求")
		return
	}

	resume, err := ctl.resumeService.GetResumeByUserID(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, resume)
}
