---
outline: [2, 3]
---

# Go 面试题

## 1 基础篇

### 1.1 Go包管理的方式有哪些

发展历史：
- `< Go1.5`：GOPATH
  - 通过统一包存放的路径实现包管理
  - 不支持依赖包的版本控制
- `>= Go1.5`：Go Vendor
  - Go1.5 引入，需要通过环境变量 `GO15VENDOREXPERIMENT=1` 开启
  - Go1.6 Vendor 机制默认开启
  - 把源码拷贝到 `vendor` 目录并维护 `vendor.json` 文件，指定版本
- `>= Go1.11`：Go Modules
  - Go1.11 中，`GO111MODULE` 的默认值是 `auto`：
    - 当项目在 GOPATH 之外，且包含 `go.mod` 文件时，自动启用 Modules；
    - 当项目在 GOPATH 内时，默认不启用（仍使用 GOPATH 模式）。
  - Go1.13 起，`GO111MODULE` 默认值改为 `on`，彻底以 Modules 为主要包管理方式。

### 1.2 `init()` 是什么时候执行的

`init()` 函数的作用：
- 程序执行前包的初始化

`init()` 函数的执行顺序：
- 在同一个 Go 文件中的多个 init 方法，按照代码顺序依次执行
- 同一个包内不同文件中的 `init()` 函数，按照文件名顺序执行
- 不同的包且不相互依赖，按照 `import` 顺序执行
- 存在依赖关系的包，被依赖的包先执行 `init()`

go 文件的初始化顺序：
1. 引入的包
2. 当前包中的常量
3. 当前包中的变量
4. 当前包的 `init()` 函数
5. 若为 `main` 包，最终执行 `main()` 函数

### 1.3 `new` 和 `make` 的区别

- `make` 不仅分配内存，还会初始化。`new` 只会分配零值填充的值（例如，`int` 的零值是 `0`，`*int` 的零值是 `nil`，`[]int` 的零值是 `nil`）
- `make` 只适用于 `slice`、`map`、`channel` 的数据，`new` 没有限制
- `make` 返回原始类型(T)，`new` 返回类型的指针(*T)

### 1.4 内存逃逸

#### 1.4.1 什么是内存逃逸

Go 中，函数内的局部变量默认分配在栈上（栈内存由编译器自动分配和释放，效率极高）。但在某些情况下，变量会被移动到堆上分配，这种现象称为**内存逃逸**。逃逸分析是编译器决定变量分配位置的过程。

核心区别：
- 栈分配：函数退出后，栈内存自动释放，无需垃圾回收（GC）。
- 堆分配：变量生命周期不确定，需由 GC 管理，会增加 GC 压力。

#### 1.4.2 发生内存逃逸的常见场景

编译器进行逃逸分析时，若发现变量的生命周期无法在编译期确定或栈无法容纳，就会将其分配到堆上。常见场景包括：

1. **变量被外部引用（跨函数生命周期）**

若函数返回变量的指针或引用，且该指针被外部持有（变量需在函数退出后继续存在），变量会逃逸到堆。

```go
func create() *int {
	x := 10   // x 会逃逸到堆
	return &x // 返回指针，x 需在函数外存活
}

func main() {
	p := create()
	fmt.Println(*p)
}
```

*原因：函数退出后，栈会被销毁，若变量仍被外部引用，必须放在堆上。*

2. **变量大小超过栈的承载能力**

栈的空间有限，若变量体积过大（如超大数组），编译器会将其分配到堆上。

```go
func bigData() {
	data := [1000000]int64{} // 大小超过栈限制，逃逸到堆
}
```

*原因：避免栈溢出（stack overflow），堆的空间更大且动态分配。*

3. **闭包引用并修改外部变量**

闭包会捕获外部变量的引用，若闭包的生命周期长于变量的原始作用域，变量会逃逸到堆。

```go
func closure() func() {
	x := 10
	return func() {
		x++ // 闭包修改 x，x 需在 closure 退出后存活
		fmt.Println(x)
	}
}

func main() {
	f := closure()
	f() // 调用闭包时，x 仍需存在
}
```

*原因：闭包可能在函数退出后被调用，变量需脱离原函数栈存活。*

4. **变量类型为接口（动态类型不确定）**

当变量被赋值给接口类型，且编译器无法在编译期确定其具体类型（动态类型），变量会逃逸到堆。

```go
type animal interface {
	run()
}

type dog struct{}

func (d dog) run() {}

func main() {
	// 声明的同时赋值了，不会发生逃逸
	var a1 animal = dog{}
	a1.run()

	// 接口类型声明时没赋值，动态类型不确定，发生内存逃逸
	var a2 animal
	a2 = dog{}
	a2.run()
}
```

*原因：接口的动态类型处理需要 runtime 支持，堆分配更灵活。*

5. **切片 / 映射的动态扩容或长度不确定**

若切片的长度是动态计算的（非编译期常量），或可能发生扩容（底层数组需更换），其底层数组可能逃逸到堆。

```go
func main() {
	n, _ := strconv.Atoi(os.Args[1])
	_ = make([]int, n)
}
```

### 1.5 如何手动修改容量和长度

正常情况下无法直接修改，必须通过反射的 `SetLen` 和 `SetCap` 方法操作。这两个方法本质是**直接修改切片头中的长度和容量字段**，而非改变底层数组本身。

- `SetLen(n)` 的限制：`n` 必须满足 `0 ≤ n ≤ cap(slice)`
- `SetCap(n)` 的限制：`n` 必须满足 `len(slice) ≤ n ≤ 原 cap(slice)`

```go
func main() {
	slice := make([]int, 3, 5)
	fmt.Println(len(slice), cap(slice), slice) // 3 5 [0 0 0]
	reflect.ValueOf(&slice).Elem().SetLen(2)
	fmt.Println(len(slice), cap(slice), slice) // 2 5 [0 0]
	reflect.ValueOf(&slice).Elem().SetCap(4)
	fmt.Println(len(slice), cap(slice), slice) // 2 4 [0 0]
}
```

### 1.6 切片和浮点数能作为 `map` 的键吗

只有可比较类型才能作为 `map` 的键。所以切片不行，浮点数可以。

但是如果用浮点数做 `map` 的键会存在问题。例如：

```go
func main() {
	m := make(map[float64]int)
	m[0.1] = 1
	m[0.2] = 2
	m[0.3] = 5
	m[0.30000000000000001] = 6
	fmt.Printf("%v\n", m) // map[0.1:1 0.2:2 0.3:6]
}
```

出现上面问题，是由于浮点数自身的二进制表示精度限制，导致不同十进制数对应同一二进制值，进而引发键冲突。当浮点型作为 `map` 的 `key` 的时候会做一些特别的处理，它会先通过 `math.Float64bit` 函数转为 `uint64` 类型，再作为 `key`：

```go
fmt.Println(math.Float64bits(0.3))                 // 4599075939470750515
fmt.Println(math.Float64bits(0.30000000000000001)) // 4599075939470750515
```

除了精度导致的冲突，浮点数键还有一个常见问题：`NaN` 与任何值都不相等，导致存入 `NaN` 后无法取出：

```go
m := make(map[float64]int)
m[math.NaN()] = 100
fmt.Println(m[math.NaN()] == 100)     // false
fmt.Println(math.NaN() == math.NaN()) // false
```

### 1.7 判断对象是否有某个方法

```go
type Person1 struct {
	Name string
}

func (p *Person1) SayHi() string {
	return "Hi"
}

type Person2 struct {
	Name string
}

// 方法一：类型断言，性能好（无反射开销）
func HasMethodSay(v interface{}) bool {
	_, has := v.(interface {
		SayHi() string
	})
	return has
}

// 方法二：通过反射，灵活性高，可以通过字符串动态指定方法名
func HasMethod(v interface{}, methodName string) bool {
	obj := reflect.ValueOf(v)
	method := obj.MethodByName(methodName)
	if !method.IsValid() {
		return false
	}
	return true
}

func main() {
	p1 := &Person1{"Zhang"}
	p2 := &Person2{"Zhang"}

	// 方法一
	fmt.Println(HasMethodSay(p1)) // true
	fmt.Println(HasMethodSay(p2)) // false

	// 方法二：通过反射
	fmt.Println(HasMethod(p1, "SayHi")) // true
	fmt.Println(HasMethod(p2, "SayHi")) // false
}
```

## 2 高级篇

### 2.1 `switch`

#### 2.1.1 基础特性：自动 `break`（与其他语言的核心区别）

Go 的 `switch` 中，每个 `case` 执行完毕后**自动终止**（无需显式写 `break`），不会像 C/C++ 那样默认穿透到下一个 `case`。

```go
num := 2
switch num {
case 1:
	fmt.Println("1")
case 2:
	fmt.Println("2") // 执行后自动退出switch，不会继续执行case 3
case 3:
	fmt.Println("3")
}
```

#### 2.1.2 `fallthrough`

如果需要像 C++ 那样执行当前 `case` 后继续执行下一个 `case`，需显式使用 `fallthrough`。

**注意事项**：
- `fallthrough` 必须放在 `case` 块的最后一行。
- 会忽略下一个 `case` 的条件，直接执行其代码块。
- 不能用在最后一个 `case` 中（否则编译错误）。


```go
num := 1
switch num {
case 1:
	fmt.Println("1")
	fallthrough // 强制执行下一个case
case 2:
	fmt.Println("2") // 即使num≠2，仍会执行
}
```

#### 2.1.3 无表达式

等价于 `switch true`（类似 if-else 链）。

```go
age := 25
switch { // 等价于 switch true
case age < 18:
	fmt.Println("未成年")
case age >= 18 && age < 30:
	fmt.Println("青年") // 匹配成功
}
```

#### 2.1.4 类型 switch：判断接口动态类型

通过 `switch v.(type)` 语法判断接口变量 `v` 的动态类型（即接口实际存储的值的类型）。

```go
var x interface{} = "hello"
switch t := x.(type) { // t是x的值，类型为case匹配的类型
case int:
	fmt.Printf("int: %d\n", t)
case string:
	fmt.Printf("string: %s\n", t) // 匹配成功，输出 "string: hello"
default:
	fmt.Printf("unknown type: %T\n", t)
}
```

#### 2.1.5 `switch` 内的变量声明

可以在 `switch` 后直接声明变量，其作用域仅限于 `switch` 块内。

```go
switch num := 3; num { // 声明变量num，仅在switch内有效
case 1:
	fmt.Println(1)
case 3:
	fmt.Println(3) // 执行
}
```

### 2.2 `defer`

#### 2.2.1 基础特性

`defer` 用于延迟一个函数（或方法）的执行，延迟到包含它的函数（父函数）返回前执行（无论父函数是正常返回还是因 `panic` 异常退出）。

```go
func main() {
	fmt.Println("start")
	defer fmt.Println("defer 1") // 延迟到 main 返回前执行
	fmt.Println("end")
}
// 输出顺序：start → end → defer 1
```

#### 2.2.2 执行顺序——后进先出（LIFO）

多个 `defer` 按声明顺序入栈，父函数返回时按栈顶到栈底顺序执行（最后声明的 `defer` 最先执行）。

```go
func main() {
	defer fmt.Println("1") // 先入栈
	defer fmt.Println("2") // 后入栈，返回时先执行
	defer fmt.Println("3") // 最后入栈，最先执行
}
// 输出顺序：3 → 2 → 1
```

#### 2.2.3 参数捕获时机——声明时计算

`defer` 后面的函数参数，会在 `defer` 声明的那一刻计算出值并固定（捕获），而非延迟函数执行时计算。

```go
func main() {
	i := 0
	defer fmt.Println("defer:", i) // 声明时 i=0，参数固定为 0
	i = 100                        // 修改 i 不影响已捕获的参数
	fmt.Println("main:", i)
}
// 输出：main: 100 → defer: 0（而非 100）
```

::: warning 注意
如果 `defer` 后面是匿名函数（闭包），且函数内部引用了外部变量，则变量的值会在延迟函数执行时读取（而非声明时），因为此时变量不是函数的参数，而是闭包的外部引用。
:::

```go
func main() {
	i := 0
	defer func() {
		// 这里的i是闭包引用的外部变量，不是函数参数
		// 延迟函数执行时（main返回前）才会读取i的当前值
		fmt.Println("defer:", i)
	}()
	i = 100 // 延迟函数执行前修改i，会影响输出结果
	fmt.Println("main:", i)
}
// 输出：main: 100 → defer: 100
```

#### 2.2.4 与函数返回值的交互（命名返回值 vs 匿名返回值）

**仅当返回值是命名返回值时**，`defer` 能影响函数返回值。

```go
func f1() int {
	i := 10
	defer func() { i += 5 }() // 修改的是局部变量 i，而非返回值
	return i                  // 返回的是 i 的副本（匿名返回值），此时 i=10
}

func f2() (i int) { // 命名返回值 i，函数开始时已初始化（默认 0）
	i = 10
	defer func() { i += 5 }() // 直接修改命名返回值 i
	return                    // 等价于 return i，此时 i 已被修改为 15
}

func main() {
	fmt.Println(f1()) // 输出 10（而非 15）
	fmt.Println(f2()) // 输出 15
}
```

#### 2.2.5 常见使用场景

1. **资源释放**：确保文件、网络连接、锁等资源在函数结束后释放，避免泄漏（无论函数正常还是异常退出，资源都能释放）。
	```go
	// 示例：文件关闭
	func readFile(path string) {
		file, err := os.Open(path)
		if err != nil {
			return
		}
		defer file.Close() // 函数返回前自动关闭文件，即使后续代码报错
		// 读取文件逻辑...
	}

	// 示例：互斥锁释放
	var mu sync.Mutex
	func safeAdd() {
		mu.Lock()
		defer mu.Unlock() // 函数结束后自动解锁，避免死锁
		// 并发安全操作...
	}
	```

2. **捕获 `panic`**：在 `defer` 中用 `recover()` 捕获父函数的 `panic`，实现异常恢复（避免程序崩溃）。
	```go
	func risky() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("捕获异常：%v\n", err) // 捕获 panic，程序不会崩溃
			}
		}()
		panic("致命错误") // 触发 panic
	}

	func main() {
		risky()
		fmt.Println("程序继续执行") // 正常输出
	}
	```

### 2.3 channel

#### 2.3.1 基础特性

channel 是 Go 语言中用于 **goroutine 间通信** 的核心机制，通过传递数据实现 goroutine 同步，而非共享内存。

底层是一个带锁的队列（或环形缓冲区），保证并发安全（多个 goroutine 操作同一个 channel 无需额外加锁）。

#### 2.3.2 核心分类：无缓冲 channel vs 有缓冲 channel

两者的核心区别在于 **是否有数据缓冲区**，直接影响发送 / 接收操作的阻塞行为。

1. 无缓冲 channel（`make(chan Type)`）
    - **特点**：没有缓冲区，发送和接收操作 **必须同步**，发送方会阻塞直到有接收方接收数据，接收方会阻塞直到有发送方发送数据。
    - **示例**：有容量为 `n` 的缓冲区，发送和接收操作的阻塞条件与缓冲区状态相关：
		```go
		ch := make(chan int)
		go func() {
			ch <- 100 // 发送方：阻塞等待接收方
		}()
		fmt.Println(<-ch) // 接收方：阻塞等待发送方，输出 100
		```

2. 有缓冲 channel（`make(chan Type, n)`）
	- **特点**：有容量为 `n` 的缓冲区，发送和接收操作的阻塞条件与缓冲区状态相关：
    	- 发送操作：缓冲区未满时，数据存入缓冲区，不阻塞；缓冲区满时，发送方阻塞直到有数据被接收。
    	- 接收操作：缓冲区非空时，直接取数据，不阻塞；缓冲区空时，接收方阻塞直到有数据被发送。
  	- **示例**：
		```go
		ch := make(chan int, 2) // 容量为2的有缓冲channel
		ch <- 1                 // 缓冲区未满（1/2），不阻塞
		ch <- 2                 // 缓冲区满（2/2），不阻塞
		//ch <- 3                 // 缓冲区满，发送方阻塞（若无人接收）

		fmt.Println(<-ch) // 取走1，缓冲区剩1（1/2），不阻塞
		fmt.Println(<-ch) // 取走2，缓冲区空，若再接收会阻塞
		```

#### 2.3.3 与 `select` 配合

`select` 用于同时监听多个 channel 的发送 / 接收操作，是 channel 最常用的高级特性。

1. **基本行为**
	- `select` 会阻塞等待 **任意一个 `case` 可执行**，然后执行该 `case`。
	- 若 **多个 `case` 同时可执行**：随机选择一个执行（避免某一 channel 饥饿）。
	- 若 **所有 `case` 都不可执行**：
    	- 有 `default` 分支：执行 `default`（非阻塞模式）。
    	- 无 `default` 分支：`select` 阻塞，直到某个 `case` 可执行。

2. **常见用法**
	- **示例 1：非阻塞接收**（避免在空 channel 上阻塞）：
		```go
		ch := make(chan int)
		select {
		case val := <-ch:
			fmt.Println("接收:", val)
		default:
			fmt.Println("channel 为空，不阻塞") // 执行此分支
		}
		```
	- **示例 2：超时控制**（避免永久阻塞）：
		```go
		ch := make(chan int)
		select {
		case val := <-ch:
			fmt.Println("接收:", val)
		case <-time.After(1 * time.Second): // 1秒后超时
			fmt.Println("超时") // 1秒内无数据则执行
		}
		```

### 2.4 鸭子类型

鸭子类型是一种以行为定义类型的设计思想，核心原则是 **如果一个东西看起来像鸭子、走起来像鸭子、叫起来像鸭子，那它就可以被当作鸭子**。它不关注对象的具体类型，只关注对象是否具备所需的 “行为”（即方法或属性）。

Go 没有类和继承，而是通过**接口（interface）** 实现鸭子类型，且是 隐式接口（无需显式声明 实现了某接口），这是 Go 的核心设计特点之一。

在 Go 中，**只要一个结构体（或其他类型）拥有某接口的 全部方法**，就会被视为 实现了该接口，无需像 Java 那样用 `implements` 显式声明。此时，该结构体的实例就可以被当作接口类型使用。

```go
// 1. 定义一个接口（代表“能叫”的行为）
type Speaker interface {
	Speak() string // 所需行为：能发出声音（Speak方法）
}

// 2. 定义两个不同类型的结构体（无继承关系）
type Dog struct {
	Name string
}
type Cat struct {
	Age int
}

// 3. 给两个结构体分别实现Speak方法（具备“叫”的行为）
func (d *Dog) Speak() string {
	return d.Name + ": 汪汪汪"
}

func (c *Cat) Speak() string {
	return fmt.Sprintf("猫咪（%d岁）: 喵喵喵", c.Age)
}

// 4. 定义一个函数：接收“能叫的对象”（即Speaker接口）
func WakeUp(s Speaker) {
	fmt.Println("唤醒声音：", s.Speak())
}

func main() {
	// Dog和Cat虽类型不同，但都有Speak方法（具备“叫”的行为）
	// 因此都能传给WakeUp（接收Speaker接口）
	dog := &Dog{Name: "阿黄"}
	cat := &Cat{Age: 2}

	WakeUp(dog) // 输出：唤醒声音：阿黄: 汪汪汪
	WakeUp(cat) // 输出：唤醒声音：猫咪（2岁）: 喵喵喵
}
```

### 2.5 重载（Overload）

同一作用域内，**函数名相同但参数列表（类型、个数、顺序）不同**的函数，编译器通过参数匹配调用对应实现（如 Java/C++）。

Go 明确不支持重载。Go 在同一包、同一作用域内，不允许定义同名函数或同名方法（即使参数列表不同），编译会直接报错。

```go
type Calculator struct{}

// 编译错误：method Calculator.Add already declared
func (c *Calculator) Add(x, y int) int {
	return x + y
}
func (c *Calculator) Add(x, y, z int) int {
	return x + y + z
}
```

### 2.6 重写（Override）

```go
type Parent struct{}
func (p Parent) F() { fmt.Println("Parent.F") }

type Child struct{ Parent } // 嵌套Parent
func (c Child) F() { fmt.Println("Child.F") } // 隐藏Parent.F

func main() {
	c := Child{}
	c.F()        // 输出：Child.F（调用外层方法）
	c.Parent.F() // 需显式调用内层方法，输出：Parent.F

	var p Parent = c.Parent
	p.F() // 输出：Parent.F
}
```

## 3 Runtime 篇

Runtime 是 Go 程序的运行时系统，主要功能包括：内存管理、垃圾回收、协程调度、反射、系统交互。

### 3.1 垃圾回收（GC）

垃圾回收是 Go runtime 自动管理内存的机制，负责识别并回收不再被引用的对象所占用的内存。

#### 3.1.1 发展历程

1. Go 1.0: 标记-清除算法（Mark-Sweep），全量STW（Stop The World），延迟高
2. Go 1.5: 三色标记法（Tri-color Marking）+ 写屏障（Write Barrier）
3. Go 1.8: 混合写屏障（Hybrid Write Barrier）

::: tip 标记清除算法的流程
1. 暂停程序业务逻辑，对可达和不可达的对象进行分类，做上标记
2. 找出程序所有可达对象并做上标记
3. 解除暂停继续执行代码，重复上述过程直到程序生命周期结束
:::

#### 3.1.2 三色标记法

当前 Go GC 的基础算法，通过颜色标记对象的可达性：

三色状态定义：
- 白色：未被标记的对象（初始状态），GC 结束后会被回收
- 灰色：自身已被标记，但引用的子对象未被标记（待处理）
- 黑色：自身及引用的子对象均已标记（存活对象）

执行流程：
1. 初始标记（Initial Mark）：
    - STW（极短），标记根对象（全局变量、栈上变量等直接可达对象）
	- 根对象从白色变为灰色，放入标记队列
2. 并发标记（Concurrent Mark）：
	- 恢复程序运行，GC 后台线程并发处理灰色对象
	- 将灰色对象变为黑色，同时将其引用的白色对象变为灰色
	- 借助写屏障跟踪并发过程中对象引用的变化
3. 并发清理（Concurrent Sweep）：
	- 程序运行的同时，回收所有白色对象的内存

#### 3.1.3 三色不变式

满足下面一种三色不变式就能确保对象不丢失：

1. 强三色不变式：
    - 不允许黑色对象直接引用白色对象
2. 弱三色不变式：
	- 允许黑色对象引用白色对象，但需满足一定条件
	- 这个白色对象必须存在其他灰色对象对它的引用
	- 或者这个白色对象的链路上游存在灰色对象

#### 3.1.3 屏障

屏障本质上是在程序执行过程中增加的额外判断机制。

如果满足一定条件就使用类似回调或者钩子（Hook）。

1. 插入屏障

又称为增量更新屏障。插入到黑色对象后的对象会被保守的标记为灰色对象，以满足强三色不变性。

2. 删除屏障

又称为基于起始快照的屏障，满足弱三色不变式。当一个白色或灰色对象的引用被移除时，该对象被标记为灰色。

#### 3.1.4 混合写屏障

1. GC 开始时扫描栈上的对象并将可达对象全部标记为黑色
2. 栈上新增的对象也直接标记为黑色，避免 STW 重复扫描
3. 栈上的引用新增和删除时都不启用写屏障机制和 STW
4. 当对象删除时触发删除写屏障，即讲删除的对象标记为灰色
5. 当对象新增的时候触发插入写屏障，将其标记为灰色

### 3.2 内存泄漏

内存泄漏指已不再使用的内存未被GC回收，导致程序占用内存持续增长。

常见内存泄漏场景及原因：

1. 未正确关闭的资源（文件句柄、网络连接等）
	```go
	func readFile() {
		f, _ := os.Open("data.txt")
		// 遗漏 `defer f.Close()`，导致文件描述符和关联内存泄漏
		data := make([]byte, 1024)
		f.Read(data)
	}
	```
2. Goroutine 无法退出导致泄漏
	1. 错误的通道操作导致的永久阻塞。
	2. `select` 分支未覆盖所有情况，且所有 `case` 都无法触发，导致 goroutine 永久阻塞。
	3. Goroutine 内部死循环。
	4. 未正确使用 `context` 控制生命周期。
3. 使用 cgo 时，没有手动释放 C 分配的内存。

### 3.3 GMP（面试必问）

#### 3.3.1 基本概念

G（Goroutine）：
  - 含义：Go 协程
  - 本质：包含执行栈、程序计数器、状态等信息的结构体（`runtime.g`）
  - 特点：创建成本低（KB级内存）、启动快，由 Go runtime 调度（非 OS 内核）

M（Machine）：
  - 含义：工作线程
  - 本质：封装了 OS 线程的结构体（`runtime.m`），与内核线程一一对应
  - 特点：由 OS 调度，创建成本高（MB 级内存），数量受 OS 限制

P（Processor）：
  - 含义：逻辑处理器，连接 G 和 M 的中间层
  - 本质：包含本地 goroutine 队列、P 状态、与 M 的绑定关系等（`runtime.p`）
  - 核心作用：
	1. 作为 G 的容器（本地队列存储待执行的 G）
	2. 限制并发执行的 M 数量（与 GOMAXPROCS 值一致）
	3. 缓存线程本地资源（如内存分配缓存）

#### 3.3.2 三者关系

M 必须绑定 P 才能执行 G，才能从 P 的本地队列或全局队列获取 G 执行。

P 管理 G 的调度：每个 P 维护一个本地队列，存放待执行的 G；全局还有一个全局队列。

数量：
- P 的数量由 `GOMAXPROCS` 控制（默认等于 CPU 核心数）
- M 的数量动态调整（默认最大10000），通常略多于 P
- G 的数量理论无上限（受内存限制）

#### 3.3.3 调度流程

1. **G 的创建与入队**

	当使用 `go func()` 创建 G 时，runtime 会：
   - 初始化 G 的栈、程序计数器等信息
   - 优先将 G 放入当前 P 的本地队列，本地队列满时放入全局队列
   - 若本地队列和全局队列都满，会触发扩容或转移逻辑

2. **G 的调度执行**
	1. M 绑定 P：M 启动后，会从 P 列表中获取一个空闲的 P 并绑定
	2. 获取 G：
        1. 优先从 P 的本地队列取 G。为了防止全局队列中的 G 饿死，P 每调度 61 次，就会优先获取全局队列中的 G。
        2. 本地队列为空时，从全局队列取一批 G（通常取 `min(全局队列长度/P数量, 128)` 个）
        3. 若全局队列也为空，触发**工作窃取**：从其他 P 的本地队列窃取一半 G
	3. 执行 G：M 执行 G 的代码，直到 G 阻塞、结束或被抢占
	4. G 状态切换：G 执行过程中可能切换状态（如 `running` → `waiting` → `runnable`），重新进入队列等待调度

#### 3.3.4 调度时机

1. **G 阻塞于系统调用（如文件I/O）**
   1. M 与 P 解绑，P 可被其他 M 获取
   2. 阻塞结束后，M 尝试绑定 P，若无空闲 P，G 进入全局队列
2. **G 阻塞于 Go 原生操作（如 channel、mutex）**
   1. G 状态置为 `waiting`，M 立即执行其他 G
   2. 当阻塞解除（如 Channel 数据到达），G 被重新加入队列
3. **抢占式调度**
   - 早期 Go 调度是协作式的，若 G 无 IO 操作、不主动让出 CPU（如死循环），会导致其他 G 饥饿
   - 当本地队列中有等待的 G，且当前 G 运行超过 10ms，触发抢占
