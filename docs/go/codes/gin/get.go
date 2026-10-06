package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello", hello)
	if err := router.Run("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}

func hello(c *gin.Context) {
	// 参数缺失或值为空时，Query 都返回空字符串。
	lang := c.Query("lang")
	// 参数缺失时使用默认值，显式空值不会被替换。
	framework := c.DefaultQuery("framework", "Gin")
	c.JSON(http.StatusOK, gin.H{"lang": lang, "framework": framework})
}
