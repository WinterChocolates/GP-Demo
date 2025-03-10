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

// CreatePermission 创建权限
// @Summary 创建权限
// @Description 创建一个新的权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param permission body models.Permission true "权限信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response "无效的请求参数"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/permissions [post]
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

// GetPermissions 获取权限列表
// @Summary 获取权限列表
// @Description 获取所有权限列表
// @Tags 权限管理
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Permission}
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/permissions [get]
func (ctl *PermissionController) GetPermissions(c *gin.Context) {
	permissions, err := ctl.service.GetPermissions(c.Request.Context())
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, permissions)
}
