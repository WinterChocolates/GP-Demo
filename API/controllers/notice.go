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
