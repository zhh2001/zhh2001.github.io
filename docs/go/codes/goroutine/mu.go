var counter int
var mu sync.Mutex
var wg = sync.WaitGroup{}

func count() {
    mu.Lock()
    defer mu.Unlock()
    defer wg.Done()
    counter++
}

func main() {
    for i := 0; i < 1000000; i++ {
        wg.Add(1)
        go count()
    }
    wg.Wait()
    fmt.Println(counter)
}