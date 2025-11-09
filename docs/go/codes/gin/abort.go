func MyLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		return  // 后续逻辑依旧会被执行
		context.Next()
	}
}