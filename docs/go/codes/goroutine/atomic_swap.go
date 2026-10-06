var value int32 = 3
swapped := atomic.CompareAndSwapInt32(&value, 3, 5)
fmt.Println(swapped, atomic.LoadInt32(&value)) // true 5
