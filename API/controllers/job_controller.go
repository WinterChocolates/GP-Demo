package controllers

import (
	"net/http"
	"strconv"

	"API/models"
	"API/services"
	"API/utils"

	"github.com/gin-gonic/gin"
)

type JobController struct {
	jobService *services.JobService
}

func NewJobController(js *services.JobService) *JobController {
	return &JobController{jobService: js}
}

// CreateJob 创建新职位
func (ctl *JobController) CreateJob(c *gin.Context) {
	var job struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&job); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	err := ctl.jobService.CreateJob(c.Request.Context(), &models.Job{
		Title:       job.Title,
		Description: job.Description,
	})

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, nil)
}

// ListJobs 获取职位列表
func (ctl *JobController) ListJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	jobs, total, err := ctl.jobService.GetOpenJobs(c.Request.Context(), page, pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data":  jobs,
		"total": total,
	})
}

// UpdateJob 更新职位信息
func (ctl *JobController) UpdateJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var job struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&job); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	err = ctl.jobService.UpdateJob(c.Request.Context(), uint(id), &models.Job{
		Title:       job.Title,
		Description: job.Description,
	})

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, nil)
}
