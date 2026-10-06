package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignUpInfo struct {
	Username   string `json:"username" binding:"required,min=3,max=20"`
	Password   string `json:"password" binding:"required,min=8,max=20"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password"`
	Email      string `json:"email" binding:"required,email"`
	Age        uint   `json:"age" binding:"lte=120"`
}

func main() {
	router := gin.Default()
	router.POST("/signUp", signUp)
	if err := router.Run("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}

func signUp(c *gin.Context) {
	var info SignUpInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "验证通过"})
}
