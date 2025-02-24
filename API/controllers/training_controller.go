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
