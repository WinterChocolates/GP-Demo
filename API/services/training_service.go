package services

import (
	"context"
	"errors"

	"API/models"
	"gorm.io/gorm"
)

type TrainingService struct {
	db *gorm.DB
}

func NewTrainingService(db *gorm.DB) *TrainingService {
	return &TrainingService{db: db}
}

func (s *TrainingService) CreateTraining(ctx context.Context, training *models.Training) error {
	return s.db.WithContext(ctx).Create(training).Error
}

func (s *TrainingService) RegisterTraining(ctx context.Context, userID, trainingID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var training models.Training
		if err := tx.First(&training, trainingID).Error; err != nil {
			return err
		}
		var count int64
		tx.Model(&models.TrainingRecord{}).Where("user_id = ? AND training_id = ?", userID, trainingID).Count(&count)
		if count > 0 {
			return errors.New("已报名该课程")
		}
		return tx.Create(&models.TrainingRecord{
			UserID:     userID,
			TrainingID: trainingID,
			Status:     "registered",
		}).Error
	})
}

func (s *TrainingService) GetTrainings(ctx context.Context) ([]models.Training, error) {
	var trainings []models.Training
	err := s.db.WithContext(ctx).Where("end_time > NOW()").Find(&trainings).Error
	return trainings, err
}

func (s *TrainingService) GetMyTrainings(ctx context.Context, userID uint) ([]models.TrainingRecord, error) {
	var records []models.TrainingRecord
	err := s.db.WithContext(ctx).Preload("Training").Where("user_id = ?", userID).Find(&records).Error
	return records, err
}
