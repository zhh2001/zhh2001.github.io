func main() {
	router := gin.Default()
	router.GET("/hello", hello)
	router.Run("localhost:8000")
}

func hello(context *gin.Context) {
	user := &proto.Teacher{
		Name:    "zhang",
		Courses: []string{"Gin", "GoLang"},
	}
	context.ProtoBuf(http.StatusOK, user)
}