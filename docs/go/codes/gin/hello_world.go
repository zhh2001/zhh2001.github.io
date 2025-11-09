package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func helloWorld(context *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "world",
	})
}

func main() {
	router := gin.Default()
	router.GET("/hello", helloWorld)
	router.Run("localhost:8000")
}
