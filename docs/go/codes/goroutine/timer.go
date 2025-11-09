// 会在 3s 后向 Channel timer.C 写入时间
timer := time.NewTimer(3 * time.Second)

select {
case msg1 := <-ch1:
    fmt.Println(msg1)
case msg2 := <-ch2:
    fmt.Println(msg2)
case <-timer.C:
    fmt.Println("Timed out")
}