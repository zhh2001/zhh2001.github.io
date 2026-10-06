package main

import (
	"log"
	"net/http"

	"example.com/gin-notes/pb"
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
	teacher := &pb.Teacher{
		Name:    "zhang",
		Courses: []string{"Gin", "GoLang"},
	}
	c.ProtoBuf(http.StatusOK, teacher)
}
