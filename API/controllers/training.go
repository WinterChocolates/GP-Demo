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

func (ctl *TrainingController) GetTrainings(c *gin.Context) {
	trainings, err := ctl.trainingService.GetTrainings(c.Request.Context())
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "获取培训列表失败")
		return
	}
	utils.RespondSuccess(c, trainings)
}

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
