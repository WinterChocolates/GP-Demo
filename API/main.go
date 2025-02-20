package main

import (
	"API/config"
	"API/storage/cache"
	"API/storage/database"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	// 初始化MySQL
	if _, err := database.InitMySQL(); err != nil {
		panic("failed to initialize MySQL: " + err.Error())
	}

	// 初始化Redis
	if err := cache.InitRedis(); err != nil {
		panic("failed to initialize Redis: " + err.Error())
	}

}
