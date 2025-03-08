package controllers

import (
	"net/http"
	"strconv"

	"API/models"
	"API/services"
	"API/utils"

	"github.com/gin-gonic/gin"
)

type NoticeController struct {
	noticeService *services.NoticeService
}

func NewNoticeController(s *services.NoticeService) *NoticeController {
	return &NoticeController{noticeService: s}
}

func (ctl *NoticeController) CreateNotice(c *gin.Context) {
	var notice models.Notice
	if err := c.ShouldBindJSON(&notice); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}
	err := ctl.noticeService.CreateNotice(c.Request.Context(), &notice)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "发布通知失败")
		return
	}
	utils.RespondSuccess(c, nil)
}

func (ctl *NoticeController) GetNotices(c *gin.Context) {
	notices, err := ctl.noticeService.GetActiveNotices(c.Request.Context())
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取通知列表失败")
		return
	}
	utils.RespondSuccess(c, notices)
}

// GetDepartmentNotices 获取指定部门的通知
func (ctl *NoticeController) GetDepartmentNotices(c *gin.Context) {
	department := c.Param("department")
	if department == "" {
		utils.RespondError(c, http.StatusBadRequest, "部门参数不能为空")
		return
	}
	notices, err := ctl.noticeService.GetDepartmentNotices(c.Request.Context(), department)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取部门通知列表失败")
		return
	}
	utils.RespondSuccess(c, notices)
}

// DeleteNotice 删除通知
func (ctl *NoticeController) DeleteNotice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的通知ID")
		return
	}
	if err := ctl.noticeService.DeleteNotice(c.Request.Context(), uint(id)); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "删除通知失败")
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "通知删除成功"})
}

// UpdateNotice 更新通知
func (ctl *NoticeController) UpdateNotice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的通知ID")
		return
	}
	var notice models.Notice
	if err := c.ShouldBindJSON(&notice); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}
	if err := ctl.noticeService.UpdateNotice(c.Request.Context(), uint(id), &notice); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "更新通知失败")
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "通知更新成功"})
}

// MarkNoticeAsRead 标记通知为已读
func (ctl *NoticeController) MarkNoticeAsRead(c *gin.Context) {
	noticeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的通知ID")
		return
	}
	
	// 从JWT中获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "未授权的请求")
		return
	}
	
	if err := ctl.noticeService.MarkNoticeAsRead(c.Request.Context(), uint(userID.(float64)), uint(noticeID)); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "标记通知已读失败")
		return
	}
	
	utils.RespondSuccess(c, gin.H{"message": "通知已标记为已读"})
}
