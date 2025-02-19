package database

import (
	"API/models"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var DB *gorm.DB

type MySQLConfig struct {
	DSN         string
	MaxIdleConn int
	MaxOpenConn int
	MaxLifetime time.Duration
}

func loadMySQLConfig() MySQLConfig {
	return MySQLConfig{
		DSN:         viper.GetString("database.mysql.dsn"),
		MaxIdleConn: viper.GetInt("database.mysql.max_idle_conn"),
		MaxOpenConn: viper.GetInt("database.mysql.max_open_conn"),
		MaxLifetime: viper.GetDuration("database.mysql.max_lifetime") * time.Second,
	}
}

func InitMySQL() (*gorm.DB, error) {
	config := loadMySQLConfig()

	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	DB = db

	return db, autoMigrate(db)
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Job{},
		&models.Application{},
		&models.Attendance{},
		&models.Notice{},
		&models.Permission{},
		&models.Role{},
		&models.Salary{},
		&models.Training{},
		&models.TrainingRecord{},
		&models.User{},
	)
}
