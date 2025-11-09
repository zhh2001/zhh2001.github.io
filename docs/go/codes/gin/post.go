func main() {
	router := gin.Default()
	router.GET("/hello", hello)   // [!code --]
	router.POST("/hello", hello)  // [!code ++]
	router.Run("localhost:8000")
}

func hello(context *gin.Context) {
	lang := context.Query("lang")                             // [!code --]
	lang := context.PostForm("lang")                          // [!code ++]
	framework := context.DefaultQuery("framework", "Gin")     // [!code --]
	framework := context.DefaultPostForm("framework", "Gin")  // [!code ++]
	context.JSON(http.StatusOK, gin.H{
		"lang":      lang,
		"framework": framework,
	})
}