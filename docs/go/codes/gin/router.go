router := gin.Default()
router.GET("/goods/list", goodsList)
router.POST("/goods/add", addGoods)
router.POST("/goods/del", delGoods)