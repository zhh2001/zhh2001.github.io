var value int32 = 18
loadedValue := atomic.LoadInt32(&value)
fmt.Println(loadedValue)  // 18