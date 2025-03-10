package controllers

import (
	"net/http"
	"strconv"
	"time"

	"API/services"
	"API/utils"

	"github.com/gin-gonic/gin"
)

type SalaryController struct {
	salaryService *services.SalaryService
}

func NewSalaryController(s *services.SalaryService) *SalaryController {
	return &SalaryController{salaryService: s}
}

// GenerateSalary 生成薪资记录
// @Summary 生成薪资记录
// @Description 为指定用户生成指定月份的薪资记录
// @Tags 薪资管理
// @Accept json
// @Produce json
// @Param request body struct{UserID uint `json:"user_id" binding:"required"` Month string `json:"month" binding:"required"`} true "薪资生成请求"
// @Success 200 {object} utils.Response{message=string}
// @Failure 400 {object} utils.Response "无效的请求参数"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/salaries/generate [post]
func (ctl *SalaryController) GenerateSalary(c *gin.Context) {
	var request struct {
		UserID uint   `json:"user_id" binding:"required"`
		Month  string `json:"month" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 验证月份格式
	if _, err := time.Parse("2006-01", request.Month); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "月份格式无效，请使用YYYY-MM格式")
		return
	}

	if err := ctl.salaryService.GenerateSalary(c.Request.Context(), request.UserID, request.Month); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"message": "薪资记录生成成功"})
}

// GetSalaryDetail 获取薪资详情
// @Summary 获取薪资详情
// @Description 获取指定月份的薪资详细信息
// @Tags 薪资管理
// @Produce json
// @Security Bearer
// @Param month path string true "月份(YYYY-MM格式)"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response "无效的月份格式"
// @Failure 401 {object} utils.Response "未授权的请求"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/salaries/{month} [get]
func (ctl *SalaryController) GetSalaryDetail(c *gin.Context) {
	month := c.Param("month")
	if month == "" {
		utils.RespondError(c, http.StatusBadRequest, "月份参数不能为空")
		return
	}

	// 验证月份格式
	if _, err := time.Parse("2006-01", month); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "月份格式无效，请使用YYYY-MM格式")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "未授权的请求")
		return
	}

	salary, err := ctl.salaryService.GetSalaryDetailByMonth(c.Request.Context(), userID.(uint), month)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, salary)
}

// GetSalaryHistory 查看薪资发放记录
// @Summary 查看薪资发放记录
// @Description 获取用户的薪资发放历史记录
// @Tags 薪资管理
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response "未授权的请求"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/salaries/history [get]
func (ctl *SalaryController) GetSalaryHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "未授权的请求")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	history, err := ctl.salaryService.GetSalaryHistory(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, history)
}
