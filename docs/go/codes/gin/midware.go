package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("status=%d latency=%s", c.Writer.Status(), time.Since(start))
	}
}

func main() {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.GET("/hello", MyLogger(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "world"})
	})
	if err := router.Run("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}
