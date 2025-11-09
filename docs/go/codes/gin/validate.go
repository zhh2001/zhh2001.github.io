type SignUpInfo struct {
	Username   string `json:"username" binding:"required,min=3,max=20"`        // 必传字段，要求 3 <= length <= 20
	Password   string `json:"password" binding:"required,min=8,max=20"`        // 必传字段，要求 3 <= length <= 20
	RePassword string `json:"rePassword" binding:"required,eqfield=Password"`  // 必传字段，要求和 password 字段一致
	Email      string `json:"email" binding:"required,email"`                  // 必传字段，要求符合邮箱格式
	Age        uint8  `json:"age" binding:"gte=0,lte=120"`                     // 要求 0 <= age <= 120
}

func main() {
	router := gin.Default()
	router.POST("/signUp", signUp)
	router.Run("localhost:8000")
}

func signUp(context *gin.Context) {
	var info SignUpInfo
	if err := context.ShouldBindJSON(&info); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
}