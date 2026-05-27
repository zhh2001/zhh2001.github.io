---
outline: deep
---

# Goroutine

goroutine 是 Go 运行时自己调度的协程，比 OS 线程轻得多：初始栈只有 2 KB（可按需扩缩），切换不走内核态，单进程拉起几十万个不是问题。`go` 关键字一加就能并发跑，剩下的调度和栈管理全部由 runtime 接管。

下面这一篇覆盖从最小例子到调度器（GMP）、同步原语（Mutex / RWMutex / atomic / WaitGroup）、Channel 及 context 的常用法和坑。

## 起一个协程

<<< @/go/codes/goroutine/hello.go

这段代码末尾的 `time.Sleep(1 * time.Second)` 是关键。主协程跑完 `main` 会直接退出，不等其他协程，去掉 Sleep 的话子协程通常根本来不及打印。Sleep 只是 demo 凑合用，正式代码里要等协程跑完，用下面会讲的 `sync.WaitGroup` 或者 `channel` 做显式同步。

两条 `Println` 是并发执行的，谁先打印不固定，下面两种结果都合法：

```text
Hello from main!
Hello from goroutine!
```

```text
Hello from goroutine!
Hello from main!
```

## GMP 调度器

Go 的调度模型俗称 GMP，是用户态协程能跑得起来的关键。三个角色：

| 缩写 | 实体 | 角色 |
| --- | --- | --- |
| **G** | goroutine | 一段要被执行的协程代码，含自己的栈、PC、状态 |
| **M** | machine | 操作系统线程，真正能被 OS 调度的执行单元 |
| **P** | processor | 逻辑处理器，G 和 M 之间的中转，持有本地运行队列 |

P 的个数等于 `GOMAXPROCS`（默认等于 CPU 核数），决定了同一时刻最多并行跑多少个 G。M 是按需创建的，没事干就睡掉。一个 M 想跑 G，必须先绑定一个 P。

### 调度的几条主线

1. `go f()` 创建出来的 G 默认放到**当前 P 的本地队列**末尾（队列满了就部分迁到全局队列）；
2. M 从绑定的 P 的本地队列头部取一个 G 来跑；
3. 本地队列空了，M 会去**全局队列**捞，再不行就走**工作窃取**，从别的 P 的本地队列后半段偷一半过来；
4. 如果系统里所有队列都空，M 解绑 P 进入休眠，等待被唤醒；
5. G 跑完不会回到运行队列，而是进 **gFree 池**等待下次 `go` 时复用栈结构。

### 阻塞场景：syscall 与 netpoller 走两条路

这是 GMP 里最容易讲错的地方，分开讲清楚：

- **G 进入 syscall**（文件读写、阻塞式系统调用）：当前 M 跟着内核陷入卡住，没法继续跑别的 G。runtime 会**把 P 解绑给另一个空闲 M（必要时新建）**，让 P 上其他 G 继续被调度。syscall 返回后，原 M 想拿回一个 P 继续跑；拿不到就把 G 塞进全局队列，自己去睡。
- **G 阻塞在网络 IO / channel / select / mutex**：G 被 park 挂起，**但 M 不会跟着卡**，它接着在本 P 上调度别的 G。事件就绪时，netpoller 把对应 G 放回运行队列，下次被调度到就接着跑。

理解这一区别之后再看那些"几万并发连接为什么 Go 不会爆"的问题，就豁然开朗。

### 抢占式调度

Go 1.14 起引入了**基于信号的异步抢占**：runtime 给 M 发 `SIGURG`，强制让 G 让出 CPU。在此之前是协作式抢占，靠编译器在函数序言里插入检查点，所以一段不调用任何函数的纯计算死循环会把调度器卡死（经典的"for {} 让其他 goroutine 跑不起来"）。Go 1.14 之后这种情况也能被抢占了。

### 几个相关优化

- **工作窃取**：保证负载在 P 之间动态均衡；
- **局部性**：G 倾向于跟着原 P 跑，提升 L1/L2 命中；
- **netpoller**：所有网络 IO 都过 epoll/kqueue/IOCP，挂起 G 不挂起 M。

## WaitGroup

`sync.WaitGroup` 是最常用的"等一组协程跑完"的同步原语。三步走：

```go
var wg sync.WaitGroup
wg.Add(n)   // 在启动 n 个协程之前
go func() {
    defer wg.Done()  // 协程退出前一定要标记完成
    // ...
}()
wg.Wait()   // 阻塞直到计数器归零
```

最小可运行示例：

1. 声明一个 `WaitGroup`：

<<< @/go/codes/goroutine/wg.go

1. 启动协程前 `Add`（必须先调用）：

<<< @/go/codes/goroutine/add.go

1. 协程末尾 `Done`，搭配 `defer` 防遗漏：

<<< @/go/codes/goroutine/done.go

1. 主协程 `Wait` 等齐：

<<< @/go/codes/goroutine/wait.go

::: warning 三个真坑

- **先 `go` 后 `Add` 是错的**。`Add` 必须在 `go func()` 之前调用，否则 `Wait` 可能在 `Add` 之前就看到 0，提前返回。
- **`Done` 多调一次会 panic**（计数器走负）。
- **`Add` 的参数可以是负数**，但实际上几乎没人这么用，正确姿势是循环里 `wg.Add(1)`。

:::

## 互斥锁 Mutex

`sync.Mutex` 守护任何不允许并发写入的共享状态。零值即可使用，不需要构造函数：

```go
var mu sync.Mutex
```

加锁、临界区、解锁三件套：

<<< @/go/codes/goroutine/sync.go

<<< @/go/codes/goroutine/mutex.go

<<< @/go/codes/goroutine/lock.go{2,3}

完整示例：

<<< @/go/codes/goroutine/mu.go

几个易踩的点：

- **不要拷贝带锁的结构体**。`Mutex` 拷一份会得到两个相互独立的锁，等于没锁。`go vet` 会报 `copylocks` 警告。
- **不要在已锁的临界区里再次 `Lock` 同一把锁**。Go 的 Mutex 不可重入，会直接死锁。
- **`Unlock` 必须由持锁的 goroutine 调用**。把 `Unlock` 用 `defer` 包好是最稳的写法。

## 原子操作 sync/atomic

`sync/atomic` 提供针对 `int32 / int64 / uint32 / uint64 / uintptr / unsafe.Pointer` 的原子操作，避免锁的开销。Go 1.19 之后还多了 `atomic.Int32` / `atomic.Int64` / `atomic.Bool` / `atomic.Pointer[T]` 等类型化版本，写起来更清爽。

常用函数：

- `Add`：原子地把 `delta` 加到 `*addr`，返回新值。

<<< @/go/codes/goroutine/atomic_add.go

- `CompareAndSwap`：CAS。比较 `*addr` 是否等于 `old`，是就改成 `new` 并返回 `true`，否则不改返回 `false`。无锁数据结构的基石。

<<< @/go/codes/goroutine/atomic_swap.go

- `Load`：原子读。

<<< @/go/codes/goroutine/atomic_load.go

- `Store`：原子写。

<<< @/go/codes/goroutine/atomic_store.go

- `Swap`：原子地把 `val` 写入 `*addr` 并返回旧值。

<<< @/go/codes/goroutine/atomic_swap_int.go

适用范围有限：原子操作只能针对单个变量做一次"读 → 改 → 写"。一旦临界区里要改两个变量、或者要先判断再改，就得回到 `Mutex`。原子操作之间没有任何顺序保证，需要的话要配合 memory barrier，但 Go 已经在 `atomic` 包里保证了 sequentially consistent 语义，业务代码里通常不用考虑。

## 读写锁 RWMutex

`sync.RWMutex` 区分了**读锁**和**写锁**：

- 多个 goroutine 可以同时持有读锁；
- 写锁是独占的，写锁持有时其他任何读 / 写都拿不到锁；
- 读锁优先级低于已经在等的写锁（避免写饥饿）。

| 方法 | 说明 |
| --- | --- |
| `RLock()` | 拿读锁 |
| `RUnlock()` | 释放读锁 |
| `Lock()` | 拿写锁 |
| `Unlock()` | 释放写锁 |
| `TryLock()` / `TryRLock()` | Go 1.18+，拿不到立刻返回 false |

适用场景：读远多于写的状态（缓存、配置、路由表）。读写比 1:1 上下的话 `RWMutex` 反而比 `Mutex` 慢，因为内部簿记开销更大。

## Channel

channel 是 Go 最具特色的并发原语，等于一根类型化的管道，传值同时也传"happens-before"语义。Go 圈子里那句口号"Don't communicate by sharing memory; share memory by communicating"指的就是它。

### 创建

<<< @/go/codes/goroutine/chan.go

`make(chan T)` 不带容量是**无缓冲**的，`make(chan T, n)` 带容量是**有缓冲**的，两者语义差别巨大。

### 无缓冲

发送和接收必须同时配对才能继续，否则任一方阻塞。常用来做严格的同步点。

<<< @/go/codes/goroutine/chan2.go

### 有缓冲

容量没满时发送不阻塞，容量没空时接收不阻塞。容量满了写、容量空了读，才会阻塞。

<<< @/go/codes/goroutine/chan3.go

### 关闭

`close(ch)` 表示生产者宣告再也不会写新值。关闭后的行为要记牢：

| 操作 | 行为 |
| --- | --- |
| 再次发送 | panic：`send on closed channel` |
| 再次 close | panic：`close of closed channel` |
| 接收 | 先把缓冲里剩的值取干净，之后返回零值；`v, ok := <-ch` 中 `ok == false` 表示已关闭且空 |
| `range ch` | 取干净后自动结束循环 |

<<< @/go/codes/goroutine/chan4.go

经验法则：**只让发送方关闭，不要让接收方关闭**，并且只在能确定不再有发送时关闭。多生产者场景下，需要额外用 `sync.Once` 或者一个独立的"关闭信号 channel"。

`for range` 是消费 channel 直到关闭的最常见写法：

```go
for v := range ch {
    fmt.Println(v)
}
// 退出循环 = ch 已被 close 且取空
```

### select 多路复用

`select` 同时盯多个 channel，谁先就绪就跑哪个 case：

<<< @/go/codes/goroutine/select.go

要点：

- 多个 case 同时就绪时，`select` **随机**挑一个，避免某些 case 长期被饿死；
- 加 `default` 分支可以让 select 立刻返回，做**非阻塞**的尝试发送/接收；
- 加 `time.After` 做**超时**：

<<< @/go/codes/goroutine/timer.go

注意 `time.After` 每次都会 `make` 一个新的定时器，长循环里频繁用会泄漏定时器，循环内推荐改用 `time.NewTimer` + `Reset`，或者让 `ctx.Done()` 当超时入口。

### 单向 channel

只发送（`chan<- T`）或只接收（`<-chan T`）。常作为函数参数类型，约束函数对 channel 的使用方向：

<<< @/go/codes/goroutine/chan_only.go

<<< @/go/codes/goroutine/chan_only_example.go

双向 channel 可以隐式赋值给单向版本，反向不行。

## 循环变量与协程的经典坑

```go
for i := 0; i < 3; i++ {
    go func() { fmt.Println(i) }()
}
```

**Go 1.21 及之前**：三个协程很可能都打印 `3`，因为闭包捕获的是同一个 `i`，循环结束时它的值是 3。修法是在循环里 `i := i` 创建新变量。

**Go 1.22 起**：循环变量改为每轮迭代独立作用域，上面这段会按预期打印 `0`、`1`、`2`（顺序仍然不固定）。所以维护老项目升到 1.22+ 后还要顺便复查一遍，行为是真的会变。

## context

`context.Context` 用来在调用链里传**取消信号**、**截止时间**和**请求级数据**。Go 里凡是涉及"可能要中途取消"的 IO、RPC、SQL，第一个参数基本都是 `ctx`。

### 核心接口

<<< @/go/codes/goroutine/context.go

四个方法：`Deadline()` 看有没有截止时间、`Done()` 拿到一个会在取消时被 close 的 channel、`Err()` 看取消原因（`Canceled` 或 `DeadlineExceeded`）、`Value(key)` 取关联值。

### 根 context

整个调用链的源头通常是 `context.Background()`：

<<< @/go/codes/goroutine/bg.go

`context.TODO()` 与 `Background()` 等价，专门用在"暂时不知道用什么，先占位"的地方。两者在静态分析工具里会被区别对待，所以语义上有差别。

### 派生 context

所有 `WithXxx` 都基于父 context 派一个子 context，子 context 自动继承父的截止时间和取消事件。

- 主动取消：

<<< @/go/codes/goroutine/with_cancel.go

- 绝对截止时间：

<<< @/go/codes/goroutine/with_deadline.go

- 相对超时（语法糖）：

<<< @/go/codes/goroutine/with_timeout.go

- 带值：

<<< @/go/codes/goroutine/with_val.go

::: warning context 的几条规矩

- `context.WithCancel` / `WithTimeout` / `WithDeadline` 返回的 `cancel()` **必须调用**，哪怕 ctx 是因超时被取消的也要 `defer cancel()`，否则会泄漏内部 timer 和 goroutine。
- 不要把 `Context` 塞进结构体字段，永远当函数参数显式传递，并放在第一个位置；
- `Value` 只用来传请求级元数据（trace id、用户身份），别拿它传业务参数；
- 自定义 key 要用**自定义类型**而不是 `string`，防止跨包冲突：`type ctxKey int; const userKey ctxKey = 1`。

:::
