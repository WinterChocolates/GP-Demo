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

// CreateNotice 创建通知
// @Summary 创建通知
// @Description 创建一个新的通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param notice body models.Notice true "通知信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response "无效的请求参数"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/notices [post]
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

// GetNotices 获取通知列表
// @Summary 获取通知列表
// @Description 获取所有活动的通知列表
// @Tags 通知管理
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Notice}
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/notices [get]
func (ctl *NoticeController) GetNotices(c *gin.Context) {
	notices, err := ctl.noticeService.GetActiveNotices(c.Request.Context())
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取通知列表失败")
		return
	}
	utils.RespondSuccess(c, notices)
}

// GetDepartmentNotices 获取指定部门的通知
// @Summary 获取部门通知
// @Description 获取指定部门的所有通知
// @Tags 通知管理
// @Produce json
// @Param department path string true "部门名称"
// @Success 200 {object} utils.Response{data=[]models.Notice}
// @Failure 400 {object} utils.Response "部门参数不能为空"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/notices/department/{department} [get]
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
// @Summary 删除通知
// @Description 删除指定的通知
// @Tags 通知管理
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} utils.Response{message=string}
// @Failure 400 {object} utils.Response "无效的通知ID"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/notices/{id} [delete]
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
// @Summary 更新通知
// @Description 更新指定通知的信息
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Param notice body models.Notice true "通知信息"
// @Success 200 {object} utils.Response{message=string}
// @Failure 400 {object} utils.Response "无效的请求参数"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/notices/{id} [put]
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
// @Summary 标记通知已读
// @Description 将指定通知标记为已读状态
// @Tags 通知管理
// @Security Bearer
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} utils.Response{message=string}
// @Failure 400 {object} utils.Response "无效的通知ID"
// @Failure 401 {object} utils.Response "未授权的请求"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/notices/{id}/read [post]
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
