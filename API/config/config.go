package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitConfig() error {
	configPath := filepath.Join("config")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	viper.SetDefault("database.mysql.max_idle_conn", 10)
	viper.SetDefault("database.mysql.max_open_conn", 100)
	viper.SetDefault("database.mysql.max_lifetime", 3600)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败 %v", err)
	}

	return nil
}
