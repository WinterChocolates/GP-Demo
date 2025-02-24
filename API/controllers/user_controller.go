package controllers

import (
	"context"
	"net/http"
	"time"

	"API/models"
	"API/services"
	"API/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(us *services.UserService) *UserController {
	return &UserController{userService: us}
}

func (ctl *UserController) Register(c *gin.Context) {
	var user struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 调用服务层
	err := ctl.userService.RegisterUser(c.Request.Context(), &models.User{
		Username:     user.Username,
		PasswordHash: user.Password,
	})

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, nil)
}

func (ctl *UserController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := ctl.userService.GetUserByID(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取用户信息失败")
		return
	}
	utils.RespondSuccess(c, user)
}

func (ctl *UserController) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	var updateData struct {
		Department string `json:"department"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	err := ctl.userService.UpdateProfile(c.Request.Context(), userID.(uint), map[string]interface{}{
		"department": updateData.Department,
	})

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "更新失败")
		return
	}

	utils.RespondSuccess(c, nil)
}

func (ctl *UserController) Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	token, err := ctl.userService.Authenticate(ctx, credentials.Username, credentials.Password)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "认证失败")
		return
	}

	utils.RespondSuccess(c, gin.H{"token": token})
}
