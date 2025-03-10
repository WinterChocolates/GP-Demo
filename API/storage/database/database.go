package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"API/models"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type MySQLConfig struct {
	Host          string
	Port          int
	Username      string
	Password      string
	DBName        string
	Charset       string
	ParseTime     bool
	MaxIdleConn   int
	MaxOpenConn   int
	MaxLifetime   time.Duration
	MaxRetries    int
	RetryInterval time.Duration
}

func loadMySQLConfig() MySQLConfig {
	return MySQLConfig{
		Host:          viper.GetString("database.mysql.host"),
		Port:          viper.GetInt("database.mysql.port"),
		Username:      viper.GetString("database.mysql.username"),
		Password:      viper.GetString("database.mysql.password"),
		DBName:        viper.GetString("database.mysql.dbname"),
		Charset:       viper.GetString("database.mysql.charset"),
		ParseTime:     viper.GetBool("database.mysql.parse_time"),
		MaxIdleConn:   viper.GetInt("database.mysql.max_idle_conn"),
		MaxOpenConn:   viper.GetInt("database.mysql.max_open_conn"),
		MaxLifetime:   viper.GetDuration("database.mysql.max_lifetime") * time.Second,
		MaxRetries:    viper.GetInt("database.mysql.max_retries"),
		RetryInterval: viper.GetDuration("database.mysql.retry_interval") * time.Second,
	}
}

func (c MySQLConfig) BuildDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
		c.ParseTime,
	)
}

func InitMySQL() (*gorm.DB, error) {
	config := loadMySQLConfig()
	log.Println("正在尝试连接数据库...")

	var db *gorm.DB
	var err error

	for attempt := 1; attempt <= config.MaxRetries+1; attempt++ {
		db, err = gorm.Open(mysql.Open(config.BuildDSN()), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			log.Printf("✅ 数据库连接成功（第%d次尝试）", attempt)
			break
		}

		log.Printf("❌ 数据库连接失败（第%d次尝试）: %v", attempt, err)
		if attempt < config.MaxRetries+1 {
			log.Printf("将在%d秒后重试...", config.RetryInterval/time.Second)
			time.Sleep(config.RetryInterval)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("数据库连接失败（重试%d次）: %v", config.MaxRetries, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	DB = db
	return db, autoMigrate(db)
}

func CheckMySQLHealth(ctx context.Context) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
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
		&models.Resume{},
	)
}

func Close() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接失败: %w", err)
	}
	log.Println("✅ MySQL连接已关闭")
	return nil
}
