package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func finalHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "handler executed"})
}

func main() {
	router := gin.Default()
	router.GET("/continue", func(c *gin.Context) {
		return // 只结束当前中间件，finalHandler 仍会执行。
	}, finalHandler)
	router.GET("/abort", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "request denied"})
		return // Abort 不会结束当前函数，仍需按逻辑返回。
	}, finalHandler)
	if err := router.Run("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}
