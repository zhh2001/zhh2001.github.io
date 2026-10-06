md := metadata.New(map[string]string{"Key1": "value1", "key2": "value2"})
fmt.Println(md.Get("key1")) // [value1]
