package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建 Gin 路由
	r := gin.New()
	r.Use(gin.Recovery())

	// 健康检查端点（用于 Kubernetes/Docker 健康检查）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// 就绪检查端点
	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	// 主业务端点
	r.GET("/", func(c *gin.Context) {
		hostname, _ := os.Hostname()
		c.JSON(http.StatusOK, gin.H{
			"message":   "Hello from Docker!",
			"hostname":  hostname,
			"version":   "1.0.0",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// 环境信息端点（用于调试）
	r.GET("/env", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"go_version":  "1.22",
			"environment": getEnv("ENV", "production"),
		})
	})

	port := getEnv("PORT", "8080")
	fmt.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
