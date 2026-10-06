package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/goods/list", goodsList)
	router.POST("/goods/add", addGoods)
	router.POST("/goods/del", delGoods)
	if err := router.Run("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}

// 以下处理函数只演示路由，不访问数据库。
func goodsList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"action": "list"})
}

func addGoods(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"action": "add"})
}

func delGoods(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"action": "del"})
}
