md := metadata.Pairs("key1", "value1", "Key1", "value2")
fmt.Println(md.Get("key1")) // [value1 value2]
