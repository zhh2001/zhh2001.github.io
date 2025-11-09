var ops uint64 = 0
atomic.AddUint64(&ops, 1)
fmt.Println("ops:", ops)