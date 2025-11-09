func MyLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		t := time.Now()
		// context.Abort()  终止请求
		context.Next()   // 处理请求
		latency := time.Since(t)
		log.Print(latency)
	}
}

func main() {
	router := gin.New()
	router.Use(gin.Logger())  // 全局中间件
	router.GET("/signUp", MyLogger(), signUp)  // 路由中间件
	router.Run("localhost:8000")
}