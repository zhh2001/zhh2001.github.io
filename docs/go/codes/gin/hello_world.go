package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "world"})
}

func main() {
	router := gin.Default()
	router.GET("/hello", helloWorld)
	if err := router.Run("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}
