package main

import (
	"fmt"
	"log"

	"employee-management/config"
	"employee-management/models"
	"employee-management/routers"

	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	if err := config.InitConfig("config.yaml"); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	if err := models.InitDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	if err := models.Migrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	r := routers.SetupRouter()
	r.Run(fmt.Sprintf(":%d", config.AppConfig.Server.Port))
}
