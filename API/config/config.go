package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigFile("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	viper.SetDefault("database.mysql.max_idle_conn", 10)
	viper.SetDefault("database.mysql.max_open_conn", 100)
	viper.SetDefault("database.mysql.max_lifetime", 3600)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败 %v", err)
	}

	return nil
}
