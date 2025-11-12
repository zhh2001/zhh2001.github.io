ctx := context.Background()

pipe := rdb.Pipeline()
incr := pipe.Incr(ctx, "pipeline_counter")
pipe.Expire(ctx, "pipeline_counter", time.Hour)

cmds, err := pipe.Exec(ctx)
if err != nil {
	panic(err)
}

// 结果需要在调用 Exec 后才能使用
fmt.Println(incr.Val())