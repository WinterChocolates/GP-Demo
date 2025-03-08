package controllers

import (
	"net/http"
	"time"

	"API/services"
	"API/utils"

	"github.com/gin-gonic/gin"
)

type AttendanceController struct {
	BaseController
	service *services.AttendanceService
}

func NewAttendanceController(s *services.AttendanceService) *AttendanceController {
	return &AttendanceController{service: s}
}

// ClockIn 上班打卡
// @Summary 上班打卡
// @Tags 考勤管理
// @Security Bearer
// @Success 200 {object} docs.SwaggerResponse
// @Router /attendance/clock-in [post]
func (ctl *AttendanceController) ClockIn(c *gin.Context) {
	userID, _ := ctl.GetAuthUser(c)

	if err := ctl.service.ClockIn(c.Request.Context(), userID); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "打卡失败: "+err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"clock_time": time.Now().Format(time.RFC3339),
		"message":    "打卡成功",
	})
}

// ClockOut 下班打卡
// @Summary 下班打卡
// @Tags 考勤管理
// @Security Bearer
// @Success 200 {object} docs.SwaggerResponse
// @Router /attendance/clock-out [post]
func (ctl *AttendanceController) ClockOut(c *gin.Context) {
	userID, _ := c.Get("userID")
	err := ctl.service.ClockOut(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "打卡失败")
		return
	}
	utils.RespondSuccess(c, nil)
}

// GetMonthly 获取月度考勤
// @Summary 获取月度考勤
// @Tags 考勤管理
// @Security Bearer
// @Param month query string true "月份格式YYYY-MM"
// @Success 200 {object} docs.SwaggerResponse
// @Router /attendance/monthly [get]
func (ctl *AttendanceController) GetMonthly(c *gin.Context) {
	userID, roles := ctl.GetAuthUser(c)
	month := c.Query("month")

	if _, err := time.Parse("2006-01", month); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "日期格式错误，请使用YYYY-MM格式")
		return
	}

	isAdmin := false
	for _, role := range roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	records, err := ctl.service.GetMonthlyAttendance(c.Request.Context(), userID, month, isAdmin)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取记录失败: "+err.Error())
		return
	}

	utils.RespondSuccess(c, records)
}

// GetAttendanceStats 获取考勤统计
// @Summary 获取考勤统计
// @Tags 考勤管理
// @Security Bearer
// @Success 200 {object} docs.SwaggerResponse
// @Router /attendance/stats [get]
func (ctl *AttendanceController) GetAttendanceStats(c *gin.Context) {
	stats, err := ctl.service.GetAttendanceStats(c.Request.Context())
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取考勤统计失败")
		return
	}
	utils.RespondSuccess(c, stats)
}
