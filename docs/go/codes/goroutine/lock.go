func example() {
    mu.Lock()
    defer mu.Unlock()

    // 访问共享资源
}