package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/sharding"
)

func main() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "postgres://localhost:5432/sharding-db?sslmode=disable",
	}))

	_ = db.Use(sharding.Register(sharding.Config{
		ShardingKey:         "user_id",
		NumberOfShards:      64,
		PrimaryKeyGenerator: sharding.PKSnowflake,
	}, "orders", Notification{}, AuditLog{}))

	// 像平常一样使用数据库会话
	// 在对表进行分片操作时，查询语句中必须包含 ShardingKey

	// GORM 创建示例，会插入到 orders_02 表
	db.Create(&Order{UserID: 2})

	// 使用原生 SQL
	db.Exec("INSERT INTO orders(user_id) VALUES(?)", int64(3))

	// 抛出 ErrMissingShardingKey 错误，缺少 ShardingKey
	db.Create(&Order{Amount: 10, ProductID: 100})

	// 查找，这将把查询重定向到 orders_02 表
	var orders []Order
	db.Model(&Order{}).Where("user_id", int64(2)).Find(&orders)
	fmt.Printf("%#v\n", orders)

	// 原生 SQL 也支持
	db.Raw("SELECT * FROM orders WHERE user_id = ?", int64(3)).Scan(&orders)
	fmt.Printf("%#v\n", orders)

	// 这将抛出 ErrMissingShardingKey 错误，因为 WHERE 条件中缺少 ShardingKey
	err = db.Model(&Order{}).Where("product_id", "1").Find(&orders).Error
	fmt.Println(err)

	// 更新和删除操作与创建和查询类似
	db.Exec("UPDATE orders SET product_id = ? WHERE user_id = ?", 2, int64(3))
	err = db.Exec("DELETE FROM orders WHERE product_id = 3").Error
	fmt.Println(err) // ErrMissingShardingKey
}
