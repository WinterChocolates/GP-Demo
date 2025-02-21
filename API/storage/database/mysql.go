package database

import (
	"fmt"
	"time"

	"API/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type MySQLConfig struct {
	DSN           string
	MaxIdleConn   int
	MaxOpenConn   int
	MaxLifetime   time.Duration
	MaxRetries    int
	RetryInterval time.Duration
}

func loadMySQLConfig() MySQLConfig {
	return MySQLConfig{
		DSN:           viper.GetString("database.mysql.dsn"),
		MaxIdleConn:   viper.GetInt("database.mysql.max_idle_conn"),
		MaxOpenConn:   viper.GetInt("database.mysql.max_open_conn"),
		MaxLifetime:   viper.GetDuration("database.mysql.max_lifetime") * time.Second,
		MaxRetries:    viper.GetInt("database.mysql.max_retries"),
		RetryInterval: viper.GetDuration("database.mysql.retry_interval") * time.Second,
	}
}

func InitMySQL() (*gorm.DB, error) {
	config := loadMySQLConfig()

	var db *gorm.DB
	var err error

	for i := 0; i <= config.MaxRetries; i++ {
		db, err = gorm.Open(mysql.Open(config.DSN), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			break
		}

		if i < config.MaxRetries {
			time.Sleep(config.RetryInterval)
		}
	}

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

func CheckMySQLHealth() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接异常: %v", err)
	}
	return nil
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
