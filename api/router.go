package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有 API 路由
func SetupRouter(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API 版本组
	v1 := r.Group("/api/v1")
	{
		// 钱包相关路由
		v1.GET("/wallet/:address/summary", nil)
		v1.GET("/wallet/:address/transactions", nil)
		v1.GET("/wallet/:address/assets", nil)
		v1.POST("/wallet/:address/analyze-ai", nil)
	}
}
