package controllers

import (
	"net/http"

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

func (ctl *SalaryController) GenerateSalary(c *gin.Context) {
	var req struct {
		UserID uint   `json:"user_id"`
		Month  string `json:"month"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}
	err := ctl.salaryService.GenerateSalary(c.Request.Context(), req.UserID, req.Month)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, nil)
}

func (ctl *SalaryController) GetSalaryDetails(c *gin.Context) {
	userID, _ := c.Get("userID")
	month := c.Param("month")
	roles, _ := c.Get("roles")
	isAdmin := false
	for _, role := range roles.([]string) {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	salary, err := ctl.salaryService.GetSalaryDetails(c.Request.Context(), userID.(uint), month, isAdmin)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, salary)
}

// GetSalaryHistory 查看薪资发放记录
func (ctl *SalaryController) GetSalaryHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "未授权")
		return
	}
	salaries, err := ctl.salaryService.GetSalaryHistory(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取薪资记录失败")
		return
	}
	utils.RespondSuccess(c, salaries)
}
