package controllers

import (
	"net/http"

	"API/models"
	"API/services"
	"API/utils"

	"github.com/gin-gonic/gin"
)

type ResumeController struct {
	BaseController
	service *services.ResumeService
}

func NewResumeController(s *services.ResumeService) *ResumeController {
	return &ResumeController{service: s}
}

// SubmitResume 提交简历
// @Summary 提交简历信息
// @Tags 用户管理
// @Security Bearer
// @Param resume body models.Resume true "简历信息"
// @Success 200 {object} docs.SwaggerResponse
// @Router /resume [post]
func (ctl *ResumeController) SubmitResume(c *gin.Context) {
	userID, _ := ctl.GetAuthUser(c)

	var resume models.Resume
	if !ctl.BindJSON(c, &resume) {
		return
	}

	if err := ctl.service.SubmitResume(c.Request.Context(), userID, &resume); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "提交简历失败: "+err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"message": "简历提交成功",
	})
}

// GetResume 获取用户简历
func (ctl *ResumeController) GetResume(c *gin.Context) {
	userID, _ := c.Get("userID") // 从认证中间件获取用户ID
	resume, err := ctl.service.GetResumeByUserID(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, resume)
}
