package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 添加多级搜索路径
	viper.AddConfigPath(".")                         // 项目根目录
	viper.AddConfigPath("./config")                  // 专用配置目录
	if exePath, err := os.Executable(); err == nil { // 兼容二进制部署
		viper.AddConfigPath(filepath.Dir(exePath))
	}
	viper.AddConfigPath("/etc/hrms/") // 系统级配置

	// 设置智能默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("jwt.secret", "default-insecure-secret")
	viper.SetDefault("jwt.expiration", 720*time.Hour) // 30天

	// 环境变量支持
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HRMS")

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("⚠️ 未找到配置文件，使用环境变量和默认配置")
		} else {
			log.Fatalf("❗ 配置文件解析错误: %v", err)
		}
	} else {
		log.Printf("✅ 成功加载配置文件: %s", viper.ConfigFileUsed())
	}

	// 输出关键配置
	log.Println("📄 有效配置:")
	log.Printf("  Server Port: %s", viper.GetString("server.port"))
	log.Printf("  MySQL DSN: %s", viper.GetString("database.mysql.dsn"))
	log.Printf("  JWT Secret: %s", maskSecret(viper.GetString("jwt.secret")))

	return nil
}

// 敏感信息脱敏显示
func maskSecret(s string) string {
	if len(s) < 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}
