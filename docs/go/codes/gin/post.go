package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/hello", hello)
	if err := router.Run("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}

func hello(c *gin.Context) {
	lang := c.PostForm("lang")
	framework := c.DefaultPostForm("framework", "Gin")
	c.JSON(http.StatusOK, gin.H{"lang": lang, "framework": framework})
}
