var ops uint64
newValue := atomic.AddUint64(&ops, 1)
fmt.Println("ops:", newValue) // ops: 1
