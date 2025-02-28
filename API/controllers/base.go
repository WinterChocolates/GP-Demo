package controllers

import (
	"net/http"
	"strconv"

	"API/utils"
	"github.com/gin-gonic/gin"
)

// BaseController 提供通用控制器方法
type BaseController struct{}

// ParsePagination 解析分页参数
func (bc *BaseController) ParsePagination(c *gin.Context) (page, size int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ = strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	return page, size
}

// BindJSON 统一JSON绑定和校验
func (bc *BaseController) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数: "+err.Error())
		return false
	}
	return true
}

// GetAuthUser 获取认证用户信息
func (bc *BaseController) GetAuthUser(c *gin.Context) (userID uint, roles []string) {
	userID = c.GetUint("userID")
	roles = c.GetStringSlice("roles")
	return
}
