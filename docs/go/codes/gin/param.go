router := gin.Default()
router.GET("/goods/:id", func(context *gin.Context) {
	id := context.Param("id")
	context.JSON(http.StatusOK, gin.H{
		"id": id,
	})
})