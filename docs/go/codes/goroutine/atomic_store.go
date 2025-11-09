var value int32 = 0
atomic.StoreInt32(&value, 20)
fmt.Println(value)  // 20