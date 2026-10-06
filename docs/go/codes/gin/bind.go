package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Goods struct {
	ID   int    `uri:"id" binding:"required"`
	Name string `uri:"name" binding:"required"`
}

func main() {
	router := gin.Default()
	router.GET("/goods/:id/:name", func(c *gin.Context) {
		var goods Goods
		if err := c.ShouldBindUri(&goods); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": goods.ID, "name": goods.Name})
	})
	if err := router.Run("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}
