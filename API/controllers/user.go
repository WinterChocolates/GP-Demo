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
	BaseController
	userService *services.UserService
}

func NewUserController(us *services.UserService) *UserController {
	return &UserController{userService: us}
}

// Register 用户注册
// @Summary 用户注册
// @Tags 用户管理
// @Param request body object true "注册信息"
// @Success 200 {object} docs.SwaggerResponse
// @Router /users/register [post]
func (ctl *UserController) Register(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}

	if !ctl.BindJSON(c, &request) {
		return
	}

	// 验证至少填写一种联系方式
	if request.Phone == "" && request.Email == "" {
		utils.RespondError(c, http.StatusBadRequest, "必须填写手机号或邮箱")
		return
	}

	// 创建用户对象
	newUser := &models.User{
		Username:     request.Username,
		PasswordHash: request.Password,
		Phone:        request.Phone,
		Email:        request.Email,
		Usertype:     "candidate",
	}

	if err := ctl.userService.RegisterUser(c.Request.Context(), newUser); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"user_id": newUser.ID,
		"message": "注册成功，请完善简历信息",
	})
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Tags 用户管理
// @Security Bearer
// @Success 200 {object} docs.SwaggerResponse{data=models.User}
// @Router /users/profile [get]
func (ctl *UserController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID") // 从认证中间件获取用户ID
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "未授权")
		return
	}
	user, err := ctl.userService.GetUserByID(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取用户信息失败")
		return
	}
	utils.RespondSuccess(c, user)
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Tags 用户管理
// @Security Bearer
// @Param request body object true "更新信息"
// @Success 200 {object} docs.SwaggerResponse
// @Router /users/profile [put]
func (ctl *UserController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "未授权")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}
	if err := ctl.userService.UpdateProfile(c.Request.Context(), userID.(uint), updates); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "更新用户信息失败")
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "用户信息更新成功"})
}

// Login 用户登录
// @Summary 用户登录
// @Tags 用户管理
// @Param request body object true "登录信息"
// @Success 200 {object} docs.SwaggerResponse{data=object{token=string}}
// @Router /users/login [post]
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
