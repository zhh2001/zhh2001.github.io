func main() {
	router := gin.Default()
	router.GET("/hello", hello)
	router.Run("localhost:8000")
}

func hello(context *gin.Context) {
	// 如果参数没取到，则使用空字符串""
	lang := context.Query("lang")
	// 如果参数没取到，则使用默认值"Gin"
	framework := context.DefaultQuery("framework", "Gin")
	context.JSON(http.StatusOK, gin.H{
		"lang":      lang,
		"framework": framework,
	})
}