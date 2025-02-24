package controllers

import (
	"net/http"

	"API/services"
	"API/utils"
	"github.com/gin-gonic/gin"
)

type AttendanceController struct {
	attendanceService *services.AttendanceService
}

func NewAttendanceController(s *services.AttendanceService) *AttendanceController {
	return &AttendanceController{attendanceService: s}
}

func (ctl *AttendanceController) ClockIn(c *gin.Context) {
	userID, _ := c.Get("userID")
	err := ctl.attendanceService.ClockIn(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "打卡失败")
		return
	}
	utils.RespondSuccess(c, nil)
}

func (ctl *AttendanceController) ClockOut(c *gin.Context) {
	userID, _ := c.Get("userID")
	err := ctl.attendanceService.ClockOut(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "打卡失败")
		return
	}
	utils.RespondSuccess(c, nil)
}

func (ctl *AttendanceController) GetMonthly(c *gin.Context) {
	userID, _ := c.Get("userID")
	yearMonth := c.Query("yearMonth")
	roles, _ := c.Get("roles")
	isAdmin := false
	for _, role := range roles.([]string) {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	records, err := ctl.attendanceService.GetMonthlyAttendance(c.Request.Context(), userID.(uint), yearMonth, isAdmin)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取考勤记录失败")
		return
	}
	utils.RespondSuccess(c, records)
}
