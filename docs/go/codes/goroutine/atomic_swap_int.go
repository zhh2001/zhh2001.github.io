var value int32 = 4
oldValue := atomic.SwapInt32(&value, 5)
fmt.Println(oldValue, value)  // 4, 5