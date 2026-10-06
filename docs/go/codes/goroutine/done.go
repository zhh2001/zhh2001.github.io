for i := 0; i < 3; i++ {
	go func(task int) {
		defer wg.Done()
		fmt.Println("Task:", task)
	}(i)
}
