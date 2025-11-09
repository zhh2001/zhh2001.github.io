var value int32 = 3
swapped := atomic.CompareAndSwapInt32(&value, 3, 5)
fmt.Println(swapped, value)  // true, 5