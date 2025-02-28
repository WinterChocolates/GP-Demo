package controllers

import (
	"net/http"

	"API/models"
	"API/services"
	"API/utils"
	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	service *services.PermissionService
}

func NewPermissionController(s *services.PermissionService) *PermissionController {
	return &PermissionController{service: s}
}

func (ctl *PermissionController) CreatePermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}
	if err := ctl.service.CreatePermission(c.Request.Context(), &permission); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, nil)
}

func (ctl *PermissionController) GetPermissions(c *gin.Context) {
	permissions, err := ctl.service.GetPermissions(c.Request.Context())
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, permissions)
}
