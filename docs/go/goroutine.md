---
title: Go Goroutine 调度、同步与并发控制
description: 介绍 Go goroutine 与 GMP 调度器，结合示例说明 WaitGroup、互斥锁、原子操作、channel 和 context 的用法，讨论并发同步、取消和竞态检查。
outline: deep
---

# Goroutine

goroutine 是由 Go 运行时调度的轻量执行单元。多个 goroutine 可以复用操作系统线程，各自的栈按需增长和收缩。栈的初始大小属于运行时实现细节，并不固定为 2 KiB。大量 goroutine 仍会占用内存和其他资源，并发数量需要结合实际负载控制。

本文以 Go 1.26.x 为基准，介绍调度器、同步原语、channel 和 context。

## 1. 启动 goroutine

`go f()` 在新的 goroutine 中调用 `f`。函数值和实参由当前 goroutine 先求值，新 goroutine 随后执行调用。并发不一定意味着并行，实际并行度还受到 `GOMAXPROCS` 和可用 CPU 的限制。

<<< @/go/codes/goroutine/hello.go

`main` 返回后，程序会退出，不会自动等待其他 goroutine 完成。上例通过关闭 `done` 通道通知任务完成，主 goroutine 在 `<-done` 处等待。`time.Sleep` 只保证至少暂停指定时长，不能保证另一个 goroutine 已经执行完毕。

两条打印语句之间没有规定先后顺序，下面两种结果都合法：

```text
Hello from main!
Hello from goroutine!
```

```text
Hello from goroutine!
Hello from main!
```

## 2. GMP 调度器

GMP 是 Go 运行时调度器的常用概括。以下按 Go 1.26.x 的实现说明，具体队列和抢占策略可能随版本变化。

| 缩写 | 实体      | 作用                                                 |
| ---- | --------- | ---------------------------------------------------- |
| G    | goroutine | 保存 goroutine 的栈、执行位置和状态等信息            |
| M    | machine   | 操作系统线程，实际执行代码                           |
| P    | processor | 调度 Go 代码所需的运行时资源，包含本地运行队列等状态 |

M 执行普通 Go 代码时需要持有 P，P 的数量由 `GOMAXPROCS` 决定。`GOMAXPROCS` 限制同时执行 Go 代码的线程数量，阻塞在系统调用中的线程不计入这个限制，所以 M 的总数可以多于 P。`runtime.GOMAXPROCS(0)` 可以查询当前设置。

Go 1.25 起，默认值会综合考虑逻辑 CPU 数、CPU 亲和性及 Linux cgroup 的 CPU 配额，并周期性检查变化。手动设置 `GOMAXPROCS` 环境变量或调用 `runtime.GOMAXPROCS(n)` 设置正值会关闭自动更新。主模块的 `go` 指令为 1.24 或更低版本时，容器感知和自动更新默认关闭，也可通过 `GODEBUG` 中的 `containermaxprocs`、`updatemaxprocs` 设置调整。

### 2.1 查找可运行的 G

可以结合运行时的 [proc.go](https://github.com/golang/go/blob/go1.26.8/src/runtime/proc.go) 理解下面的过程：

1. `go` 语句创建的 G 通常先进入当前 P 的 `runnext` 槽位，而不是直接放到本地队列末尾。原有的 `runnext` 任务可以移入普通队列，本地队列满时会转移一批任务到全局队列。
2. 调度器从本地队列、全局队列、网络轮询和工作窃取等来源寻找任务。它也会周期性检查全局队列，避免持续存在的本地任务使全局任务长期得不到调度。
3. 工作窃取会从其他 P 获取一部分待运行任务，以分担负载。本地队列也减少了对全局队列的竞争。
4. 暂时没有任务时，M 可以释放 P，等待任务或事件到来。
5. G 结束后，其运行时对象可以进入空闲缓存供后续复用，栈可能保留或释放。

这些实现细节不保证 goroutine 按创建顺序执行，也不提供固定的调度延迟。

### 2.2 阻塞与唤醒

不同的阻塞原因有不同的处理路径：

- **阻塞式系统调用和部分 cgo 调用**会占用当前 M。运行时可以让其他 M 接管 P，继续执行其他 G，但接管不一定在进入系统调用时立即发生。调用返回后，所在 M 需要重新取得 P，才能继续执行 Go 代码。
- **支持运行时网络轮询的 I/O**在 Unix 上通常使用非阻塞文件描述符。等待就绪时可以挂起 G，让 M 调度其他任务，事件就绪后由 netpoller 将 G 置为可运行状态。具体机制包括 Linux 的 epoll、部分 Unix 系统的 kqueue 和 Windows 的 IOCP。
- **channel、select 和锁的等待**也可能挂起 G，但由通信、关闭通道或释放锁等相应路径唤醒。锁竞争还可能先短暂自旋，这些等待不能全部归到 netpoller。

### 2.3 抢占式调度

Go 1.14 在支持的平台上引入异步抢占。Unix 系统通常通过信号请求抢占，例如 Linux 使用 `SIGURG`，运行时仍需要检查当前位置是否可以安全抢占。

此前的协作式抢占依赖函数调用等检查点。在 `GOMAXPROCS=1` 的情况下，不调用函数的忙循环可能长时间占用唯一的 P。异步抢占改善了这种情况，但不保证任意指令都能立即被抢占，也不能代替程序中的同步和取消机制。

## 3. WaitGroup

`sync.WaitGroup` 用于等待一组任务完成，零值可直接使用。它只管理任务计数，不会自动保护任务访问的共享变量。

### 3.1 Add、Done 与 Wait

下面四段代码按顺序放入同一个函数，并导入 `fmt` 和 `sync`，即可等待三个任务完成。

声明 `WaitGroup`：

<<< @/go/codes/goroutine/wg.go

启动任务前增加计数：

<<< @/go/codes/goroutine/add.go

每个任务结束时调用一次 `Done`，通常使用 `defer`：

<<< @/go/codes/goroutine/done.go

等待计数归零：

<<< @/go/codes/goroutine/wait.go

三个任务的打印顺序不固定，`All Done` 一定在它们全部完成打印后出现。使 `Wait` 返回的 `Done` 调用建立同步关系，任务在此之前的写入对等待方可见。

::: warning 使用要求

- 计数为 0 时，增加计数的 `Add` 必须先于 `Wait`。通常在启动 goroutine 前调用，避免 `Wait` 提前返回或 `Done` 先执行。
- `Done` 等价于 `Add(-1)`。计数变成负数会触发 `panic`，增加和减少的数量应匹配。
- 复用 `WaitGroup` 等待下一组任务时，要等上一组的所有 `Wait` 调用返回后再增加新任务。
- `WaitGroup` 首次使用后不能复制，跨函数传递时通常使用指针。

:::

### 3.2 WaitGroup.Go

Go 1.25 新增 `WaitGroup.Go`，将增加计数、启动 goroutine 和任务返回后减少计数合为一次调用：

<<< @/go/codes/goroutine/wg_go.go

传入的函数不应发生未恢复的 `panic`。使用 `wg.Go` 时，不需要再为同一个任务手动调用 `Add` 或 `Done`。空 `WaitGroup` 的首次 `Go` 调用仍应先于 `Wait`。

## 4. 互斥锁 Mutex

`sync.Mutex` 用于保护共享状态。当多个 goroutine 可能并发修改或读写同一数据时，相关访问需要遵循同一种同步方式。使用锁保护时，读和写都应遵守同一把锁。

零值即可使用：

<<< @/go/codes/goroutine/sync.go

<<< @/go/codes/goroutine/mutex.go

获得锁后进入临界区，通过 `defer` 确保离开函数时解锁：

<<< @/go/codes/goroutine/lock.go{2,3}

完整示例启动 1000 个任务，每个任务将共享计数器加一：

<<< @/go/codes/goroutine/mu.go

`count` 的延迟调用会先解锁，再通知任务完成。`Wait` 返回后，主 goroutine 才读取最终计数。

使用时需要注意：

- `Mutex` 首次使用后不能复制，含锁的结构体通常通过指针传递。复制锁可能破坏同步关系，`go vet` 可以检查许多这类问题。
- `Mutex` 不可重入，已持有锁时再次 `Lock` 同一把锁会阻塞。
- 锁不绑定特定 goroutine，允许由另一个 goroutine 调用 `Unlock`。普通代码中，在获得锁的函数内配对使用 `defer Unlock` 更容易维护。
- 对未锁定的 `Mutex` 调用 `Unlock` 会导致运行时致命错误。

## 5. 原子操作 sync/atomic

`sync/atomic` 提供不可被其他并发访问分割的原子操作。Go 1.19 引入了 `atomic.Int32`、`atomic.Int64`、`atomic.Uint64`、`atomic.Bool` 和 `atomic.Pointer[T]` 等类型化接口，通常比传入裸指针的函数更方便。

### 5.1 类型化接口

下面通过 `atomic.Uint64` 累加共享计数器：

<<< @/go/codes/goroutine/atomic.go

类型化原子值首次使用后不能复制。共享变量的并发访问应统一使用原子操作，不能一边原子写入，一边用普通方式读取。`atomic.Pointer[T]` 只保证指针本身的原子访问，不会自动保护它指向的对象。

### 5.2 函数接口

旧式函数接口操作 `int32`、`int64`、`uint32`、`uint64`、`uintptr` 和 `unsafe.Pointer` 等类型。下面是常用操作：

`Add` 原子地增加指定值，并返回新值：

<<< @/go/codes/goroutine/atomic_add.go

`CompareAndSwap` 比较当前值与 `old`，相等时替换为 `new` 并返回 `true`，否则不修改并返回 `false`：

<<< @/go/codes/goroutine/atomic_swap.go

`Load` 原子读取：

<<< @/go/codes/goroutine/atomic_load.go

`Store` 原子写入：

<<< @/go/codes/goroutine/atomic_store.go

`Swap` 原子替换并返回旧值：

<<< @/go/codes/goroutine/atomic_swap_int.go

在部分 32 位平台上，旧式 64 位原子函数对地址对齐有要求，类型化的 `atomic.Int64` 和 `atomic.Uint64` 会自动满足对齐要求。

### 5.3 顺序与适用范围

Go 的原子操作具有顺序一致性，执行效果相当于这些原子操作按一个与各 goroutine 程序顺序一致的全局顺序发生。如果一个原子操作观察到另一个原子操作的效果，两者之间还建立同步关系。这是 [Go 内存模型](https://go.dev/ref/mem#atomic)规定的语义，使用这些接口时无需额外手写内存屏障。

CAS 可以用于条件更新，但多个原子操作并不会自动组成一个原子事务。需要共同维护多个字段的不变量时，通常用 `Mutex` 更清楚。原子操作也存在竞争和缓存同步开销，是否更快需要结合实际负载测量。

## 6. 读写锁 RWMutex

`sync.RWMutex` 的零值可直接使用，区分读锁和写锁：

- 多个 goroutine 可以同时持有读锁，读锁内只能进行不会与其他读者冲突的访问。
- 写锁独占，获得写锁前需要等待现有读者和写者释放锁。
- 写者等待时，新的读锁请求会被阻塞，直到该写者获得并释放锁。

| 方法                       | 说明                                       |
| -------------------------- | ------------------------------------------ |
| `RLock()`                  | 获得读锁                                   |
| `RUnlock()`                | 释放读锁                                   |
| `Lock()`                   | 获得写锁                                   |
| `Unlock()`                 | 释放写锁                                   |
| `TryLock()` / `TryRLock()` | Go 1.18 引入，尝试获得锁，不等待锁变为可用 |

`RWMutex` 首次使用后不能复制，也不支持递归读锁、读锁直接升级为写锁或写锁直接降级为读锁。尝试加锁失败时没有建立同步关系，不能据此直接读取需要保护的数据。

读操作较多时可以考虑 `RWMutex`，但效果还取决于临界区长度、竞争程度和 CPU 并行度。读写比例本身不足以判断它是否比 `Mutex` 更快。

## 7. Channel

channel 是带元素类型的通信通道，发送和对应接收可以建立同步关系。例如，发送前完成的写入，对对应接收完成后的代码可见。

通道传递的是元素值的副本。若元素是指针、切片或 map，接收方仍可能与发送方共享底层数据，后续访问仍需遵守约定的同步方式。

### 7.1 创建

<<< @/go/codes/goroutine/chan.go

`make(chan T)` 创建无缓冲通道，`make(chan T, n)` 创建容量为 `n` 的缓冲通道，`n` 为 0 时仍是无缓冲通道。未初始化的通道为 `nil`，对它发送或接收会一直阻塞，调用 `close` 会触发 `panic`。

### 7.2 无缓冲

开放的无缓冲通道需要发送方和接收方配对才能完成通信。发送方要等待接收方准备好，接收方也要等待发送方，因而可以用于同步：

<<< @/go/codes/goroutine/chan2.go

### 7.3 有缓冲

开放的缓冲通道有空位时发送可以完成，有数据时接收可以完成。缓冲区满时发送阻塞，缓冲区空时接收阻塞：

<<< @/go/codes/goroutine/chan3.go

一个发送方连续发送、一个接收方连续接收时，接收顺序与发送顺序一致。多个发送方并发发送时，不能依赖它们按 goroutine 创建顺序到达。

### 7.4 关闭

`close(ch)` 表示以后不会再向该通道发送数据：

| 操作            | 关闭后的行为                               |
| --------------- | ------------------------------------------ |
| 发送            | `panic`，错误为 `send on closed channel`   |
| 再次关闭        | `panic`，错误为 `close of closed channel`  |
| 接收            | 先取出已有缓冲值，再立即返回元素类型的零值 |
| `v, ok := <-ch` | 已关闭且缓冲区为空时，`ok` 为 `false`      |
| `range ch`      | 持续接收，已关闭且缓冲区为空时结束         |

<<< @/go/codes/goroutine/chan4.go

通道通常由发送方或协调所有发送者的 goroutine 关闭。多个发送者共享通道时，可以先等待它们全部完成发送，再关闭通道：

<<< @/go/codes/goroutine/chan_multi.go

`sync.Once` 只能防止重复关闭，不能阻止其他 goroutine 同时发送。关闭前仍必须确保后续不会再发送。通道不需要依靠 `close` 释放内存，关闭主要用于通知接收方发送已经结束。

### 7.5 select 多路复用

`select` 从能够进行通信的分支中选择一个执行：

<<< @/go/codes/goroutine/select.go

上例的两个通道都已有数据，输出可能是 `ch1`，也可能是 `ch2`。需要注意：

- 多个通信分支同时就绪时，按均匀伪随机方式选择一个，不保证固定优先级或有限时间内每个分支都被选中。
- 没有通信分支就绪时，有 `default` 就执行它，否则阻塞等待。反复执行带 `default` 的空循环可能忙等。
- nil 通道对应的分支不会就绪。循环中可以把已处理完的通道变量设为 `nil`，停用该分支。
- 已关闭通道的接收始终就绪，循环中需要处理 `ok == false`，否则可能反复读到零值。
- 进入 `select` 时，通道表达式和发送值会先求值。未选中的发送分支也可能执行这些表达式中的副作用。

可以使用定时器限制本次等待时间。下面没有向 `ch` 发送数据，因此会输出 `Timed out`：

<<< @/go/codes/goroutine/timer.go

`time.After(d)` 也能提供一次性定时通道。Go 1.23 引入的新语义允许垃圾回收未被引用的定时器，不能再笼统地说循环中调用 `time.After` 一定会泄漏。频繁创建定时器仍有分配开销，需要时可以复用 `time.NewTimer` 并调用 `Reset`。

新语义下，定时器通道是同步通道，`Stop` 或 `Reset` 返回后不会再收到旧设置产生的时间值。主模块的 `go` 指令低于 1.23 时默认使用旧语义，`GODEBUG=asynctimerchan=1` 也会恢复旧语义，设置为 0 则强制使用新语义。旧语义下重用定时器时，需要正确停止并处理尚未消费的旧通知，不能直接照搬新语义的写法。

### 7.6 单向 channel

`chan<- T` 只允许发送，`<-chan T` 只允许接收，常用于限制函数参数的使用方向：

<<< @/go/codes/goroutine/chan_only.go

生产者和消费者可以分别使用这两种参数类型：

<<< @/go/codes/goroutine/chan_only_example.go

双向通道可以赋给对应的单向通道，单向通道不能赋回双向通道。发送方向的通道可以用于 `close`，接收方向的通道不能关闭。

## 8. 循环变量与 goroutine

下面通过闭包读取循环变量，并用 `WaitGroup` 等待所有打印完成：

<<< @/go/codes/goroutine/loop_var.go

在 Go 1.22 及之后的语言版本中，通过 `:=` 声明的循环变量在每轮迭代中独立，上例会分别打印 `0`、`1`、`2`，顺序不固定。

Go 1.21 及之前的语言语义复用同一个 `i`，循环更新与 goroutine 读取可能产生数据竞争，不能保证输出是什么。可以在循环内使用 `i := i` 创建新变量，或把当前值作为参数传给 goroutine：

```go
for i := 0; i < 3; i++ {
	wg.Add(1)
	go func(value int) {
		defer wg.Done()
		fmt.Println(value)
	}(i)
}
wg.Wait()
```

语言版本通常由模块的 `go.mod` 中的 `go` 指令决定，只升级工具链不一定改变旧模块的循环语义。使用 `=` 给循环外已有变量赋值时，即使采用新语言版本，仍然复用该变量。

## 9. context

`context.Context` 用于传递取消信号、截止时间和请求级数据。需要上下文的函数通常将 `ctx context.Context` 放在第一个参数位置。不要传入 nil，暂时无法确定上下文时可以使用 `context.TODO()`。

Context 的方法可以被多个 goroutine 同时调用，但它不会自动为业务数据提供互斥保护。

### 9.1 核心接口

<<< @/go/codes/goroutine/context.go

`Deadline` 返回截止时间及是否存在截止时间。`Done` 在上下文取消时关闭，不能取消的上下文可以返回 nil。`Err` 在未取消时返回 nil，取消后返回 `context.Canceled` 或 `context.DeadlineExceeded`。`Value` 查询关联值，不存在时返回 nil。

### 9.2 根 context

顶层入口可以使用 `context.Background()`：

<<< @/go/codes/goroutine/bg.go

`Background` 和 `TODO` 都没有截止时间、取消信号或关联值。`Background` 表示明确的顶层上下文，`TODO` 表示当前尚未确定应该传入哪个上下文。

### 9.3 派生 context

下面几种派生方式保留父上下文的值查询，并接收父上下文的取消信号。取消子上下文不会反过来取消父上下文或其他子上下文。

主动取消：

<<< @/go/codes/goroutine/with_cancel.go

设置绝对截止时间：

<<< @/go/codes/goroutine/with_deadline.go

设置相对超时：

<<< @/go/codes/goroutine/with_timeout.go

`WithDeadline` 和 `WithTimeout` 不能延长父上下文已有的更早截止时间。

关联请求级数据：

<<< @/go/codes/goroutine/with_val.go

`WithValue` 的 key 必须非 nil 且可比较，通常使用包内自定义类型，避免与其他包冲突。读取时要使用同一类型的 key，并检查类型断言结果。`Value` 适合请求级元数据，不适合代替普通业务参数，共享可变值仍需另行同步。

Go 1.20 新增 `WithCancelCause` 和 `Cause`，可以在标准取消状态之外保存具体原因：

```go
ctx, cancel := context.WithCancelCause(context.Background())
defer cancel(nil)
cancel(errors.New("client disconnected"))
fmt.Println(ctx.Err())          // context canceled
fmt.Println(context.Cause(ctx)) // client disconnected
```

Go 1.21 新增的 `WithoutCancel` 会保留父上下文的值查询，但不继承取消信号和截止时间，其 `Done` 为 nil，`Err` 为 nil。需要让任务脱离原请求时，应另外明确它的超时和结束条件。

### 9.4 取消与等待退出

取消只是在上下文中传播信号，不会强行终止 goroutine，也不会等待任务退出。任务需要主动检查 `Done` 或 `Err`，并使用支持 context 的阻塞操作。调用方还需要通过通道或 `WaitGroup` 等待任务结束：

<<< @/go/codes/goroutine/context_worker.go

`cancel` 可以重复调用。获得 `WithCancel`、`WithDeadline`、`WithTimeout` 返回的取消函数后，通常应立即安排 `defer cancel()`，任务提前完成时也要释放关联资源。未调用取消函数可能使父上下文保留子上下文及相关资源，但这些函数并不都创建定时器或 goroutine。

一般将 Context 作为参数显式传递，避免存入长期存在的结构体并跨请求复用。如果 `select` 中的工作分支与 `ctx.Done()` 同时就绪，取消分支也没有自动获得优先级，具体退出行为要由任务逻辑保证。

## 10. 并发代码检查

在 Go 模块中，可以结合测试和静态检查检查并发代码：

```shell
go test -race ./...
go vet ./...
```

单文件程序也可以运行竞态检测，例如：

```shell
go run -race hello.go
```

竞态检测只能检查实际执行到的路径。测试未报警仍需要确认其他路径的同步关系，死锁和 goroutine 长期不退出等问题也需要结合测试及退出条件检查。
