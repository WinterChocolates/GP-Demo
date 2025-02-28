package controllers

import (
	"net/http"
	"strconv"

	"API/services"
	"API/utils"
	"github.com/gin-gonic/gin"
)

type ApplicationController struct {
	applicationService *services.ApplicationService
}

func NewApplicationController(as *services.ApplicationService) *ApplicationController {
	return &ApplicationController{applicationService: as}
}

// UpdateApplicationStatus 更新申请状态
func (ctl *ApplicationController) UpdateApplicationStatus(c *gin.Context) {
	applicationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的申请ID")
		return
	}
	var request struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}
	if err := ctl.applicationService.UpdateApplicationStatus(c.Request.Context(), uint(applicationID), request.Status); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "状态更新成功"})
}
