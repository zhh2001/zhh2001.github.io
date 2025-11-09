---
outline: deep
---

# Goroutine

Go 协程是由 Go 运行时管理的轻量级线程。它们允许在单个进程中同时执行多个函数或方法，而无需显式地管理线程的生命周期。Go 协程比传统的操作系统线程更轻量，启动和切换的开销更小，因此可以轻松地创建成千上万个协程。

## 1 创建协程

<<< @/go/codes/goroutine/hello.go

输出结果可能是：

```text
Hello from main!
Hello from goroutine!
```

也可能是：

```text
Hello from goroutine!
Hello from main!
```

因为两个打印操作是并发执行的，顺序不确定。

## 2 GMP 调度器

GMP 调度器是 Go 语言实现高效并发的关键组件。它通过将 goroutine 与操作系统线程解耦，并使用工作窃取、局部性原理、抢占式调度和网络轮询器等技术，实现了高效的并发调度和执行。

### 2.1 组成部分

1. <strong>G（Goroutine）</strong>​：表示一个 goroutine，它是 Go 语言中的轻量级线程。每个 goroutine 都有自己的栈、程序计数器和局部变量等。
2. <strong>​M（Machine）</strong>​：表示一个操作系统线程。M 负责执行 G，并与操作系统进行交互。一个 M 可以执行多个 G，但是在任意时刻，一个 M 只能执行一个 G。
3. <strong>P（Processor）</strong>​：表示一个逻辑处理器，它是 G 和 M 之间的中间层。P 负责管理和调度分配给它的 G。每个 P 都有自己的本地队列，用于存储待执行的 G。P 的数量决定了可以同时运行的 G 的最大数量。

### 2.2 工作原理

1. 当创建一个新的 goroutine 时，它会被添加到当前 P 的本地队列中。
2. M 从与其关联的 P 的本地队列中获取一个 G 来执行。如果本地队列为空，M 会尝试从其他 P 的本地队列中窃取 G，或者从全局队列中获取 G。
3. 如果所有队列都为空，M 会进入休眠状态，等待新的 G 被创建或者被唤醒。
4. 当一个 G 阻塞时（例如等待 I/O 操作完成），M 会将其与当前的 P 解除关联，并尝试从其他 P 的本地队列中获取一个新的 G 来执行。
5. 当一个 G 完成执行时，它会返回到其所属的 P 的本地队列中，等待下一次调度。

### 2.3 优化

1. <strong>工作窃取（Work Stealing）</strong>​：当一个 P 的本地队列为空时，它会尝试从其他 P 的本地队列中窃取 G。这种机制可以平衡各个 P 之间的负载，提高整体调度效率。
2. <strong>局部性原理（Locality Principle）</strong>​：GMP 调度器会尽量将相关的 G 调度到同一个 P 上执行，以提高缓存命中率和减少内存访问延迟。
3. <strong>抢占式调度（Preemptive Scheduling）</strong>​：Go 1.14 引入了基于信号的抢占式调度，可以在长时间运行的 G 上强制插入调度点，以避免某个 G 长时间占用 M，导致其他 G 无法得到执行。
4. <strong>网络轮询器（Network Poller）</strong>​：Go 运行时内置了一个网络轮询器，用于高效地处理 I/O 多路复用。当 G 阻塞在 I/O 操作上时，M 可以与网络轮询器协作，继续执行其他 G，从而提高整体性能。

## 3 WaitGroup

`WaitGroup` 是 Go 语言中用于等待一组 goroutine 完成的同步原语。它属于 `sync` 包，提供了一种简单而有效的方式来协调多个并发执行的 goroutine，确保主程序在所有子任务完成后再继续执行。

1. **创建一个 `WaitGroup`**

<<< @/go/codes/goroutine/wg.go

2. **添加等待的 goroutine 数量**

使用 `Add` 方法来设置需要等待的 goroutine 数量。这通常在启动 goroutine 之前调用。

<<< @/go/codes/goroutine/add.go

3. **在每个 goroutine 完成时调用 `Done`**

在每个 goroutine 的逻辑结束时，调用 `Done` 方法来通知 `WaitGroup` 该 goroutine 已完成。通常使用 `defer` 语句确保 `Done` 被调用。

<<< @/go/codes/goroutine/done.go

4. **等待所有 goroutine 完成**

在主 goroutine 中调用 `Wait` 方法，它会阻塞直到所有的 goroutine 调用了 `Done`。

<<< @/go/codes/goroutine/wait.go

## 4 互斥锁

Go 语言中的互斥锁（Mutex）是一种同步原语，用于在多个 goroutine 之间保护共享资源的访问，防止数据竞争和不一致性。互斥锁确保在同一时间只有一个 goroutine 能够访问被保护的资源。

Go 的 `sync` 包提供了 `Mutex` 类型，用于实现互斥锁。以下是互斥锁的基本使用方法：

首先，需要导入 `sync` 包：

<<< @/go/codes/goroutine/sync.go

定义一个 `sync.Mutex` 变量：

<<< @/go/codes/goroutine/mutex.go

使用 `Lock()` 方法对互斥锁进行加锁，使用 `Unlock()` 方法进行解锁。通常，`Unlock()` 会在 `defer` 语句中调用，以确保在函数退出时释放锁。

<<< @/go/codes/goroutine/lock.go{2,3}

完整示例：

<<< @/go/codes/goroutine/mu.go

## 5 原子操作包

`sync/atomic` 包提供了一系列原子操作函数，这些函数可以对基本数据类型（如 `int32`、`int64`、`uint32`、`uint64`、`uintptr` 和 `unsafe.Pointer`）进行原子性的读取、写入、比较和交换等操作。

以下是一些常用的 `atomic` 操作函数：

1. `Add`：原子地将 `delta` 添加到 `*addr` 并返回新值。

<<< @/go/codes/goroutine/atomic_add.go

2. `CompareAndSwap`：原子地比较 `addr` 的旧值和 `old`。如果相等，则将 `addr` 的值设为新值并返回 `true`；否则返回 `false`。

<<< @/go/codes/goroutine/atomic_swap.go

3. `Load`：原子地加载 `*addr`。

<<< @/go/codes/goroutine/atomic_load.go

4. `Store`：原子地存储 `val` 到 `*addr`。

<<< @/go/codes/goroutine/atomic_store.go

5. `Swap`：原子地将 `val` 存储到 `*addr` 并返回 `*addr` 的旧值。

<<< @/go/codes/goroutine/atomic_swap_int.go

原子操作在并发编程中非常有用，因为它们可以避免使用锁，从而减少性能开销和死锁的风险。然而，需要注意的是，并非所有操作都可以通过原子操作来实现。对于复杂的操作，仍然需要使用互斥锁（`sync.Mutex`）或其他同步机制来确保数据的一致性。

## 6 读写互斥锁

`RWMutex`（读写互斥锁）是 Go 语言中用于管理共享资源访问的一种同步机制，位于 `sync` 包下。与普通的互斥锁（`Mutex`）不同，`RWMutex` 允许多个读操作同时进行，但在写操作时会阻塞所有其他读写操作。

### 6.1 主要特点

1. <b>读锁（RLock）</b>​：
    - 多个 goroutine 可以同时持有读锁。
    - 当有线程持有读锁时，其他 goroutine 仍然可以获取读锁，但不能获取写锁。
2. <b>写锁（Lock）</b>​：
    - 写锁是独占的，当一个 goroutine 持有写锁时，其他任何 goroutine 都不能获取读锁或写锁。
    - 写锁会阻塞所有其他读锁和写锁的获取，直到写锁被释放。

### 6.2 常用方法

|        方法 | 说明         |
| ----------: | ------------ |
|   `RLock()` | 获取一个读锁 |
| `RUnlock()` | 释放一个读锁 |
|    `Lock()` | 获取一个写锁 |
|  `Unlock()` | 释放一个写锁 |

## 7 Channel

- <b>Channel</b>：一种类型化的管道，可以通过它发送和接收特定类型的值。
- <b>发送（Send）</b>：将一个值发送到 Channel 中。
- <b>接收（Receive）</b>：从 Channel 中接收一个值。
- <b>阻塞（Blocking）</b>：如果发送方尝试向满的 Channel 发送数据，或者接收方尝试从空的 Channel 接收数据，操作将会阻塞，直到另一端准备好。

### 7.1 创建

<<< @/go/codes/goroutine/chan.go

### 7.2 无缓冲

无缓冲 Channel 在发送和接收操作完成之前会阻塞。

<<< @/go/codes/goroutine/chan2.go

### 7.3 有缓冲

有缓冲 Channel 允许在阻塞之前存储一定数量的元素。

<<< @/go/codes/goroutine/chan3.go

### 7.4 关闭

关闭 Channel 表示不会再有更多的值发送到该 Channel。接收方可以通过检测 Channel 是否关闭来决定如何处理接收到的数据。

<<< @/go/codes/goroutine/chan4.go

### 7.5 多路复用（Select）

`select` 语句允许同时等待多个 Channel 操作，类似于 `switch` 语句，但用于 Channel。

<<< @/go/codes/goroutine/select.go

当多个 `case` 同时就绪时，`select` 会随机选择一个执行。这有助于避免某些 Goroutine 被饿死（即一直得不到执行机会）。

通常还会设置一个超时时间，避免阻塞等待的时间过长：

<<< @/go/codes/goroutine/timer.go

### 7.6 单向 Channel

单向Channel分为两种类型：

1. <b>只发送通道（Send-only Channel）</b>​：只能用于发送数据，不能用于接收数据。
2. <b>只接收通道（Receive-only Channel）</b>​：只能用于接收数据，不能用于发送数据。

首先，需要创建一个双向 Channel，然后可以将其转换为单向 Channel：

<<< @/go/codes/goroutine/chan_only.go

单向 Channel 常用于函数参数，以确保函数只能以特定的方式使用通道：

<<< @/go/codes/goroutine/chan_only_example.go

## 8 context

`context` 包提供了一种机制，用于在不同的 Goroutine 之间传递请求范围内的数据、取消信号以及截止时间（deadline）。

主要用途：

1. <b>​传递请求范围内的数据</b>：如请求 ID、用户认证信息等。
2. <b>​取消信号</b>：通知相关的 Goroutine 停止当前操作。
3. <b>​设置截止时间（Deadline）</b>​：为操作设置一个超时时间，防止操作无限期地阻塞。

### 8.1 核心接口

<<< @/go/codes/goroutine/context.go

### 8.2 创建根上下文

根上下文是一个没有任何值、不会被取消且没有截止时间的上下文，通常作为其他上下文的父级使用。

<<< @/go/codes/goroutine/bg.go

### 8.3 派生上下文

可以从根上下文或其他派生上下文中创建新的上下文，添加取消功能、截止时间或键值对。

带取消功能的上下文：

<<< @/go/codes/goroutine/with_cancel.go

带截止时间的上下文：

<<< @/go/codes/goroutine/with_deadline.go

或者使用 `WithTimeout` 简化：

<<< @/go/codes/goroutine/with_timeout.go

带值的上下文：

<<< @/go/codes/goroutine/with_val.go
