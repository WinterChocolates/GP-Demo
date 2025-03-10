package controllers

import (
	"net/http"
	"strconv"

	"API/models"
	"API/services"
	"API/utils"

	"github.com/gin-gonic/gin"
)

type TrainingController struct {
	trainingService *services.TrainingService
}

func NewTrainingController(s *services.TrainingService) *TrainingController {
	return &TrainingController{trainingService: s}
}

// CreateTraining 创建培训课程
// @Summary 创建培训课程
// @Description 创建一个新的培训课程
// @Tags 培训管理
// @Accept json
// @Produce json
// @Param training body models.Training true "培训课程信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response "无效的请求参数"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/trainings [post]
func (ctl *TrainingController) CreateTraining(c *gin.Context) {
	var training models.Training
	if err := c.ShouldBindJSON(&training); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}
	err := ctl.trainingService.CreateTraining(c.Request.Context(), &training)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "创建培训课程失败")
		return
	}
	utils.RespondSuccess(c, nil)
}

// RegisterTraining 报名培训
// @Summary 报名培训
// @Description 用户报名参加指定的培训课程
// @Tags 培训管理
// @Security Bearer
// @Produce json
// @Param id path int true "培训ID"
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/trainings/{id}/register [post]
func (ctl *TrainingController) RegisterTraining(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	err := ctl.trainingService.RegisterTraining(c.Request.Context(), userID.(uint), uint(id))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, nil)
}

// GetTrainings 获取培训列表
// @Summary 获取培训列表
// @Description 获取所有可用的培训课程列表
// @Tags 培训管理
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Training}
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/trainings [get]
func (ctl *TrainingController) GetTrainings(c *gin.Context) {
	trainings, err := ctl.trainingService.GetTrainings(c.Request.Context())
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取培训列表失败")
		return
	}
	utils.RespondSuccess(c, trainings)
}

// GetMyTrainings 获取我的培训记录
// @Summary 获取我的培训记录
// @Description 获取当前用户的所有培训记录
// @Tags 培训管理
// @Security Bearer
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.TrainingRecord}
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/trainings/my [get]
func (ctl *TrainingController) GetMyTrainings(c *gin.Context) {
	userID, _ := c.Get("userID")
	records, err := ctl.trainingService.GetMyTrainings(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取培训记录失败")
		return
	}
	utils.RespondSuccess(c, records)
}

// UpdateTrainingRecord 更新培训记录
// @Summary 更新培训记录
// @Description 更新指定培训记录的状态和分数
// @Tags 培训管理
// @Accept json
// @Produce json
// @Param id path int true "记录ID"
// @Param request body struct{Status string `json:"status"` Score uint8 `json:"score"`} true "更新信息"
// @Success 200 {object} utils.Response{message=string}
// @Failure 400 {object} utils.Response "无效的记录ID或请求参数"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/training-records/{id} [put]
func (ctl *TrainingController) UpdateTrainingRecord(c *gin.Context) {
	recordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的记录ID")
		return
	}
	var request struct {
		Status string `json:"status"`
		Score  uint8  `json:"score"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}
	if err := ctl.trainingService.UpdateTrainingRecord(c.Request.Context(), uint(recordID), request.Status, request.Score); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "记录更新成功"})
}

// CancelTrainingRegistration 取消培训注册
// @Summary 取消培训注册
// @Description 取消已报名的培训课程
// @Tags 培训管理
// @Security Bearer
// @Produce json
// @Param id path int true "记录ID"
// @Success 200 {object} utils.Response{message=string}
// @Failure 400 {object} utils.Response "无效的记录ID"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/training-records/{id}/cancel [post]
func (ctl *TrainingController) CancelTrainingRegistration(c *gin.Context) {
	recordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的记录ID")
		return
	}
	if err := ctl.trainingService.CancelTrainingRegistration(c.Request.Context(), uint(recordID)); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "取消培训注册失败")
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "培训注册已取消"})
}

// GetTrainingDetail 获取培训详情
// @Summary 获取培训详情
// @Description 获取指定培训课程的详细信息
// @Tags 培训管理
// @Produce json
// @Param id path int true "培训ID"
// @Success 200 {object} utils.Response{data=models.Training}
// @Failure 400 {object} utils.Response "无效的培训ID"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/trainings/{id} [get]
func (ctl *TrainingController) GetTrainingDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "无效的培训ID")
		return
	}
	training, err := ctl.trainingService.GetTrainingByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取培训详情失败")
		return
	}
	utils.RespondSuccess(c, training)
}
