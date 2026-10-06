var value int32
atomic.StoreInt32(&value, 20)
fmt.Println(atomic.LoadInt32(&value)) // 20
