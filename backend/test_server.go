package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"employee-management/config"
	"employee-management/models"
	"employee-management/routers"
)

func startServer() {
	// 初始化配置
	if err := config.InitConfig("config.yaml"); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 初始化数据库
	if err := models.InitDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 执行数据库迁移
	if err := models.Migrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 设置路由
	router := routers.SetupRouter()

	// 配置服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.AppConfig.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器
	log.Printf("服务器启动在端口 %d", config.AppConfig.Server.Port)
	log.Printf("访问地址: http://localhost:%d", config.AppConfig.Server.Port)

	// 添加优雅关闭
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	// 这里可以使用信号处理，简化版本直接阻塞
	select {}
}
