md := metadata.Pairs("payload-bin", string([]byte{0, 255}))
fmt.Println([]byte(md.Get("payload-bin")[0])) // [0 255]
