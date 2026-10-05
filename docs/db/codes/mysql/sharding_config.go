package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/sharding"
)

type Order struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64
	ProductID int64
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func expectMissingKey(err error) {
	if !errors.Is(err, sharding.ErrMissingShardingKey) {
		log.Fatalf("expected ErrMissingShardingKey, got %v", err)
	}
	fmt.Println("ErrMissingShardingKey")
}

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("set MYSQL_DSN to the demo_sharding database")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	must(err)
	sqlDB, err := db.DB()
	must(err)
	defer sqlDB.Close()

	// 先建立物理表，再注册插件。示例使用 64 张表。
	for i := 0; i < 64; i++ {
		must(db.Table(fmt.Sprintf("orders_%02d", i)).AutoMigrate(&Order{}))
	}
	must(db.Use(sharding.Register(sharding.Config{
		ShardingKey:         "user_id",
		NumberOfShards:      64,
		PrimaryKeyGenerator: sharding.PKSnowflake,
	}, "orders")))

	// 分别路由到 orders_02 和 orders_03。
	must(db.Create(&Order{UserID: 2, ProductID: 100}).Error)
	must(db.Exec("INSERT INTO orders(user_id, product_id) VALUES(?, ?)", 3, 101).Error)

	var orders []Order
	must(db.Where("user_id = ?", 2).Find(&orders).Error)
	fmt.Printf("user 2: %#v\n", orders)
	must(db.Raw("SELECT * FROM orders WHERE user_id = ?", 3).Scan(&orders).Error)
	fmt.Printf("user 3: %#v\n", orders)

	// 缺少路由所需的键，不会自动遍历所有分片。
	expectMissingKey(db.Exec("INSERT INTO orders(product_id) VALUES(?)", 102).Error)
	expectMissingKey(db.Where("product_id = ?", 100).Find(&orders).Error)
	must(db.Exec("UPDATE orders SET product_id = ? WHERE user_id = ?", 103, 3).Error)
	expectMissingKey(db.Exec("DELETE FROM orders WHERE product_id = ?", 103).Error)
}
