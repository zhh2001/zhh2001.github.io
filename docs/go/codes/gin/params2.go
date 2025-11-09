router.GET("/goods/:id/*action", func(context *gin.Context) {
	id := context.Param("id")
	action := context.Param("action")
	context.JSON(http.StatusOK, gin.H{
		"id":     id,
		"action": action,
	})
})