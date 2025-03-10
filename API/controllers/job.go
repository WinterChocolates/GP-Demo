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
// @Summary 创建新职位
// @Description 创建一个新的职位信息
// @Tags 职位管理
// @Accept json
// @Produce json
// @Param job body struct{Title string `json:"title" binding:"required"` Description string `json:"description"`} true "职位信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response "无效的请求参数"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/jobs [post]
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
// @Summary 获取职位列表
// @Description 获取所有开放的职位列表
// @Tags 职位管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=[]models.Job,total=int}
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/jobs [get]
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
// @Summary 更新职位信息
// @Description 更新指定职位的信息
// @Tags 职位管理
// @Accept json
// @Produce json
// @Param id path int true "职位ID"
// @Param job body struct{Title string `json:"title"` Description string `json:"description"`} true "职位信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response "无效的请求参数"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/jobs/{id} [put]
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

// ApplyForJob 申请职位
// @Summary 申请职位
// @Description 用户申请指定职位
// @Tags 职位管理
// @Security Bearer
// @Produce json
// @Param id path int true "职位ID"
// @Success 200 {object} utils.Response{message=string}
// @Failure 400 {object} utils.Response "无效的职位ID"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/jobs/{id}/apply [post]
func (ctl *JobController) ApplyForJob(c *gin.Context) {
	userID, _ := c.Get("userID") // 从认证中间件获取用户ID
	jobID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的职位ID")
		return
	}
	if err := ctl.jobService.ApplyForJob(c.Request.Context(), userID.(uint), uint(jobID)); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "申请成功"})
}

// DeleteJob 删除职位
// @Summary 删除职位
// @Description 删除指定职位
// @Tags 职位管理
// @Produce json
// @Param id path int true "职位ID"
// @Success 200 {object} utils.Response{message=string}
// @Failure 400 {object} utils.Response "无效的职位ID"
// @Failure 500 {object} utils.Response "删除职位失败"
// @Router /api/v1/jobs/{id} [delete]
func (ctl *JobController) DeleteJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的职位ID")
		return
	}
	if err := ctl.jobService.DeleteJob(c.Request.Context(), uint(id)); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "删除职位失败")
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "职位删除成功"})
}
