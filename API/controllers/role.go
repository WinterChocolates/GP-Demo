package controllers

import (
	"net/http"

	"API/models"
	"API/services"
	"API/utils"
	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type RoleController struct {
	roleService *services.RoleService
}

// NewRoleController 初始化角色控制器
func NewRoleController(rs *services.RoleService) *RoleController {
	return &RoleController{roleService: rs}
}

// CreateRole 创建角色
func (ctl *RoleController) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}
	if err := ctl.roleService.CreateRole(c.Request.Context(), &role); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "创建角色失败")
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "角色创建成功"})
}

// GetRoles 获取角色列表
func (ctl *RoleController) GetRoles(c *gin.Context) {
	roles, err := ctl.roleService.GetRoles(c.Request.Context())
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取角色列表失败")
		return
	}
	utils.RespondSuccess(c, roles)
}
