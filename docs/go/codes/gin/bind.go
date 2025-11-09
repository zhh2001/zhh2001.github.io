type Goods struct {
	Id   int    `uri:"id"   binding:"required"`
	Name string `uri:"name" binding:"required"`
}

func main() {
	router := gin.Default()
	router.GET("/goods/:id/:name", func(context *gin.Context) {
		var goods Goods
		if err := context.ShouldBindUri(&goods); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"id":   goods.Id,
			"name": goods.Name,
		})
	})
	router.Run("localhost:8000")
}