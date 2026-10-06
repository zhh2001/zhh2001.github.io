---
outline: deep
---

# Go

Go 是一门静态类型、编译型语言，提供垃圾回收、goroutine、channel 和标准工具链。本文按语言基础、数据结构、接口、模块、测试和代码生成的顺序整理。

## 1. 注释

Go 支持行注释和块注释。普通说明优先使用行注释。块注释更适合暂时注释一段代码，或者在表达式中插入说明。

### 1.1 单行注释

单行注释以两个正斜杠（`//`）开头。

普通行注释不会作为语句执行。`//go:generate`、`//go:build` 等指令也使用注释形式，但会由相应工具处理。

```go
// 这是一个单行注释
package main

import "fmt"

func main() {
	// 这是一个单行注释
	fmt.Println("Hello World!") // 这是一个单行注释
}
```

### 1.2 多行注释

多行注释以 `/*` 开头和 `*/` 结尾。

普通块注释不会作为语句执行，块注释不能嵌套。

```go
package main

import "fmt"

func main() {
	/* 这是一个
	多行注释 */
	fmt.Println("Hello World!")
}
```

## 2. 变量

### 2.1 声明变量

常用的变量声明有两种：

使用 `var` 关键字，后面跟变量名和数据类型：

```go
var age int = 1
```

::: warning 注意
数据类型和初始值至少指定一项，也可以同时指定。未提供初始值时，变量取该类型的零值，例如 `int` 为 `0`，`bool` 为 `false`，`string` 为空字符串。

```go
var age1 int
var age2 = 1
```

:::

使用 `:=` 符号，后面接变量值：

```go
age := 1
```

::: warning 注意

- 编译器会根据右侧表达式推断变量类型。
- `:=` 只能在函数体内使用，左侧至少要有一个在当前块中尚未声明的非空白变量，`_` 不算新变量。
- 同一块中的已有变量可以与新变量一起使用，但类型必须保持一致。内层块中的同名声明会遮蔽外层变量。
:::

### 2.2 声明多变量

在 Go 中，可以在同一行声明多个变量。

```go
var a, b, c int = 1, 2, 3
```

::: warning 注意
一个 `var` 声明指定类型时，该声明中的所有变量使用同一种类型。
:::

如果不指定数据类型，就能在同一行中声明不同数据类型的变量：

```go
var name1, age1 = "Zhang", 23
name2, age2 := "Zhang", 23
```

多个变量声明也可以组合成一个块，可读性更高：

```go
var (
	m    int
	n           = 1
	age  int    = 23
	name string = "Zhang"
)
```

## 3. 常量

通过 `const` 关键字声明常量，常量的值是无法修改的。

```go
const PI float64 = 3.1415926
```

::: warning 注意
常量的值必须由常量表达式确定，不能在运行时计算。分组声明中，除第一项外，可以省略表达式，此时沿用前一项的类型和表达式。
:::

Go 不要求常量名全部大写。名称仍按普通 Go 标识符使用 MixedCaps，例如 `maxRetries` 或 `DefaultTimeout`。

```go
const (
	A = 23
	B
	C = "Zhang"
	D
)
```

在上面的声明中，常量 `B` 的值与 `A` 一样，常量 `D` 的值与 `C` 一样。

```go
const (
	A = iota
	B
	C
)
```

在上面的声明中，常量 `A` 的值为 `0`，`B` 为 `1`，`C` 为 `2`。

`iota` 是预声明标识符，在每个 `const` 声明块中从 0 开始，并随每个 `ConstSpec` 递增。

Go 常量只能表示布尔值、字符、整数、浮点数、复数和字符串。数组、切片、map、结构体、函数以及依赖运行时计算的结果都不能声明为常量。

## 4. 基本数据类型

### 4.1 布尔

```go
var b1 bool = true
var b2 = true
var b3 bool
b4 := true
```

### 4.2 整数

整数类型分为两类：

- **有符号整数**：可以存储正值、零和负值。
- **无符号整数**：只能存储非负值。

`int` 和 `uint` 的宽度取决于编译目标，均为 32 位或 64 位，两者宽度相同。不能只根据宿主操作系统判断，交叉编译时还要考虑目标架构。`strconv.IntSize` 给出当前目标的 `int` 位数。下表的指数是数学记法，Go 中的 `^` 表示按位异或，并不是乘方运算。

#### 4.2.1 有符号整数

| 类型    | 大小           | 范围                       |
| ------- | -------------- | -------------------------- |
| `int`   | 32 位或 64 位  | 与 `int32` 或 `int64` 相同 |
| `int8`  | 8 位 / 1 字节  | $-2^7$ 到 $2^7 - 1$        |
| `int16` | 16 位 / 2 字节 | $-2^{15}$ 到 $2^{15} - 1$  |
| `int32` | 32 位 / 4 字节 | $-2^{31}$ 到 $2^{31} - 1$  |
| `int64` | 64 位 / 8 字节 | $-2^{63}$ 到 $2^{63} - 1$  |

#### 4.2.2 无符号整数

| 类型     | 大小           | 范围                         |
| -------- | -------------- | ---------------------------- |
| `uint`   | 32 位或 64 位  | 与 `uint32` 或 `uint64` 相同 |
| `uint8`  | 8 位 / 1 字节  | $0$ 到 $2^8 - 1$             |
| `uint16` | 16 位 / 2 字节 | $0$ 到 $2^{16} - 1$          |
| `uint32` | 32 位 / 4 字节 | $0$ 到 $2^{32} - 1$          |
| `uint64` | 64 位 / 8 字节 | $0$ 到 $2^{64} - 1$          |

#### 4.2.3 其他

`byte` 是 `uint8` 的别名，用于表达字节值，定义如下：

```go
type byte = uint8
```

`rune` 是 `int32` 的别名，通常表示一个 Unicode 码点，定义如下：

```go
type rune = int32
```

### 4.3 浮点数

浮点数类型有 `float32` 和 `float64`，分别采用 IEEE 754 的 32 位和 64 位二进制浮点表示。它们不能精确表示所有十进制小数，例如 `0.1`，计算结果可能需要舍入。通过 `f := 1.2` 声明变量时，默认类型为 `float64`。

### 4.4 字符串

字符串是不可变的字节序列，不保证内容一定是有效的 UTF-8。`len(s)` 返回字节数，`s[i]` 返回第 `i` 个字节。需要逐个 Unicode 码点处理有效的 UTF-8 文本时，可以使用 `range` 或转换为 `[]rune`。

```go
var s string = "Hello World"
```

#### 4.4.1 字符串的拼接

```go
name := "Zhang"
age := 23

// 少量固定字符串直接使用 +，最清楚
s1 := "name: " + name + ", age: " + strconv.Itoa(age)

// 需要格式控制时使用 Sprintf
s2 := fmt.Sprintf("name: %s, age: %d", name, age)

// 循环或分段构造较长字符串时使用 Builder
builder := strings.Builder{}
builder.WriteString("name: ")
builder.WriteString(name)
builder.WriteString(", age: ")
builder.WriteString(strconv.Itoa(age))
s3 := builder.String()

fmt.Println(s1) // name: Zhang, age: 23
fmt.Println(s2) // name: Zhang, age: 23
fmt.Println(s3) // name: Zhang, age: 23
```

不存在对所有场景都最快的拼接方式。少量固定片段优先使用 `+`，循环拼接或片段较多时再考虑 `strings.Builder`。如果已经有字符串切片，`strings.Join` 往往更直接。`Builder` 的零值可直接使用，开始写入后不要复制它。

#### 4.4.2 字符串的比较

可以直接使用 `==`、`!=`、`<` 等比较运算符。字符串按字节的字典序比较，不会自动按自然语言的排序规则处理。

```go
s1 := "Zhang"
s2 := "Zhang"

fmt.Println(s1 == s2) // true
```

## 5. 格式化输出

### 5.1 常规格式

以下是常用的通用格式。类型实现了 `fmt.Formatter`、`fmt.Stringer` 等接口时，输出还可能由自定义方法决定。

| 格式  | 含义                                 |
| ----- | ------------------------------------ |
| `%v`  | 以默认格式输出                       |
| `%+v` | 在打印结构体时，会显示字段名和字段值 |
| `%#v` | 以Go语法的格式输出                   |
| `%T`  | 输出值的类型                         |
| `%%`  | 输出一个百分号，不消耗参数           |

```go
f := 12.3
s := "Zhang"
a := []int{1, 2, 3}

fmt.Printf("%v\n", f)  // 12.3
fmt.Printf("%#v\n", f) // 12.3
fmt.Printf("%T\n", f)  // float64

fmt.Printf("%v\n", s)  // Zhang
fmt.Printf("%#v\n", s) // "Zhang"
fmt.Printf("%T\n", s)  // string

fmt.Printf("%v\n", a)  // [1 2 3]
fmt.Printf("%#v\n", a) // []int{1, 2, 3}
fmt.Printf("%T\n", a)  // []int

fmt.Printf("%%\n") // %
```

### 5.2 整数格式

以下输出格式要和整数类型一起使用：

格式中的宽度是最小宽度，数值较长时不会被截断。例如 `%4d` 输出 `12345` 时仍保留全部五位数字。

| 格式   | 含义                                      |
| ------ | ----------------------------------------- |
| `%b`   | 以二进制格式输出                          |
| `%o`   | 以八进制格式输出                          |
| `%O`   | 以八进制格式输出并且显示前缀 `0o`         |
| `%d`   | 以十进制格式输出                          |
| `%+d`  | 以十进制格式输出并且显示符号              |
| `%x`   | 以十六进制格式小写输出                    |
| `%X`   | 以十六进制格式大写输出                    |
| `%#x`  | 以十六进制格式小写输出并且显示前缀 `0x`   |
| `%#X`  | 以十六进制格式大写输出并且显示前缀 `0X`   |
| `%4d`  | 以最小宽度为 `4` 的格式输出，左侧填充空格 |
| `%-4d` | 以最小宽度为 `4` 的格式输出，右侧填充空格 |
| `%04d` | 以最小宽度为 `4` 的格式输出，左侧填充 `0` |

```go
var i = 123

fmt.Printf("%b\n", i)   // 1111011
fmt.Printf("%o\n", i)   // 173
fmt.Printf("%O\n", i)   // 0o173
fmt.Printf("%d\n", i)   // 123
fmt.Printf("%+d\n", i)  // +123
fmt.Printf("%x\n", i)   // 7b
fmt.Printf("%X\n", i)   // 7B
fmt.Printf("%#x\n", i)  // 0x7b
fmt.Printf("%#X\n", i)  // 0X7B
fmt.Printf("%4d\n", i)  //  123
fmt.Printf("%-4d\n", i) // 123
fmt.Printf("%04d\n", i) // 0123
```

### 5.3 字符串格式

以下输出格式要和字符串类型一起使用：

| 格式   | 含义                                       |
| ------ | ------------------------------------------ |
| `%s`   | 纯字符串输出                               |
| `%q`   | 用双引号包裹，并按 Go 字符串字面量规则转义 |
| `%8s`  | 以最小宽度为 `8` 的格式输出，左侧填充空格  |
| `%-8s` | 以最小宽度为 `8` 的格式输出，右侧填充空格  |
| `%x`   | 以十六进制输出字符串的每个字节             |
| `% x`  | 以十六进制输出字符串的每个字节并用空格分隔 |

字符串的宽度通常按 Unicode 码点数计算，不等同于字节数或终端显示列数。

```go
var s = "Zhang"

fmt.Printf("%s\n", s)   // Zhang
fmt.Printf("%q\n", s)   // "Zhang"
fmt.Printf("%8s\n", s)  //    Zhang
fmt.Printf("%-8s\n", s) // Zhang
fmt.Printf("%x\n", s)   // 5a68616e67
fmt.Printf("% x\n", s)  // 5a 68 61 6e 67
```

### 5.4 布尔格式

以下输出格式要和布尔类型一起使用：

| 格式 | 含义       |
| ---- | ---------- |
| `%t` | 输出布尔值 |

```go
fmt.Printf("%t\n", true) // true
```

### 5.5 浮点数格式

以下输出格式要和浮点数类型一起使用：

| 格式    | 含义                        |
| ------- | --------------------------- |
| `%f`    | 保留 `6` 位小数             |
| `%.2f`  | 保留 `2` 位小数             |
| `%6.2f` | 最小宽度 `6`，精度 `2`      |
| `%e`    | 科学计数法，输出的 `e` 小写 |
| `%E`    | 科学计数法，输出的 `E` 大写 |

```go
var f float64 = 0.125

fmt.Printf("%f\n", f)    // 0.125000
fmt.Printf("%.2f\n", f)  // 0.12
fmt.Printf("%6.2f\n", f) //   0.12
fmt.Printf("%e\n", f)    // 1.250000e-01
fmt.Printf("%E\n", f)    // 1.250000E-01
```

上例的 `0.125` 能被二进制浮点精确表示。保留两位小数时恰好处于两种结果的中点，`fmt` 按最接近的偶数舍入，因此结果为 `0.12`。

## 6. 条件语句

```go
age := 23
if age == 18 {
	fmt.Println("刚成年")
} else if age > 18 {
	fmt.Println("已成年")
} else {
	fmt.Println("未成年")
}
```

## 7. 循环语句

`for` 循环是 Go 语言中唯一的循环语句。

```go
for i := 0; i < 10; i++ {
	if i == 2 {
		continue
	}
	if i == 6 {
		break
	}
	fmt.Printf("%d ", i) // 0 1 3 4 5
}

for index, value := range "Zhang" {
	fmt.Printf("%d-%c ", index, value) // 0-Z 1-h 2-a 3-n 4-g
}

for index, value := range []int{66, 88, 99} {
	fmt.Printf("%d-%d ", index, value) // 0-66 1-88 2-99
}
```

遍历字符串时，`index` 是 UTF-8 字节偏移量，`value` 是 `rune`。字符串含有非 ASCII 字符时，索引不一定连续。Go 1.22 起也可以对整数执行 `range`，例如 `for i := range 10` 会依次得到 0 到 9。

当模块使用 Go 1.22 或更高的语言版本时，循环通过 `:=` 声明的变量在每次迭代中都是独立变量，闭包可以捕获各自的值。通过 `=` 给已有变量赋值时，仍然复用同一变量。

```go
var funcs []func() int
for i := range 3 {
	funcs = append(funcs, func() int { return i })
}
for _, f := range funcs {
	fmt.Println(f()) // 依次输出 0、1、2
}
```

## 8. `goto` 语句

`goto` 只能跳转到同一函数中的标签，不能跳入更内层的块，也不能跨过变量声明而让变量尚未初始化就进入作用域。

```go
func main() {
	i := 1
LOOP:
	fmt.Printf("%v ", i) // 1 2 3
	i++
	if i <= 3 {
		goto LOOP
	}
}
```

## 9. `switch` 语句

表达式 `switch` 匹配一个分支后默认结束，不需要写 `break`。需要继续执行紧接着的下一分支时，可在分支末尾使用 `fallthrough`，下一分支的条件不会重新判断。

```go
status := 200
switch status {
case 200:
	fmt.Println("OK")
case 403:
	fmt.Println("Permission Denied")
case 404:
	fmt.Println("Not Found")
default:
	fmt.Println("Unknown status")
}
```

## 10. 数组

```go
arr1 := [3]int{1, 2, 3}
arr2 := [...]int{4, 5}
fmt.Printf("%T %T", arr1, arr2) // [3]int [2]int
```

::: warning 注意
在 Go 语言中，数组长度固定，要么给数组指定长度，要么使用 `...` 让编译器推断数组的长度。
:::

只初始化数组特定位置：

```go
arr1 := [8]int{3: 33, 5: 55}
arr2 := [...]int{6: 66}
fmt.Println(arr1) // [0 0 0 33 0 55 0 0]
fmt.Println(arr2) // [0 0 0 0 0 0 66]
```

可以使用 `==`、`!=` 直接比较两个数组是否相等：

```go
arr1 := [3]int{1, 2, 3}
arr2 := [...]int{1, 2, 3}
fmt.Println(arr1 == arr2) // true
fmt.Println(arr1 != arr2) // false
```

数组长度也是类型的一部分，`[2]int` 和 `[3]int` 是不同类型。只有元素类型可比较且两个数组的类型满足比较规则时，才能使用 `==` 和 `!=`。数组赋值和传参会复制整个数组值。切片不能互相使用 `==` 或 `!=`，但可以与 `nil` 比较。

## 11. 切片

### 11.1 创建

切片是对底层数组某个连续区间的描述。切片本身包含长度和容量。重新切片可以改变可见范围，`append` 则返回更新后的切片，并可能分配新的底层数组。

在 Go 中，有几种方法可以创建切片：

- 直接声明
- 从数组创建
- 使用 `make()` 函数

```go
arr := [...]int{1, 2, 3, 4, 5}
slice1 := []int{1, 2, 3}
slice2 := arr[1:4]
slice3 := make([]int, 3, 6)
fmt.Println(slice1, len(slice1), cap(slice1)) // [1 2 3] 3 3
fmt.Println(slice2, len(slice2), cap(slice2)) // [2 3 4] 3 4
fmt.Println(slice3, len(slice3), cap(slice3)) // [0 0 0] 3 6
```

`cap` 返回切片的容量。对数组使用 `arr[low:high]` 时，容量为数组长度减去 `low`。完整切片表达式 `arr[low:high:max]` 可以限制容量，此时容量为 `max - low`。

`make()` 的第二个参数为长度，第三个参数为容量，如果容量未指定，则默认等于长度。

`var s []int` 声明的是 nil 切片，其长度和容量均为 0。`[]int{}` 是非 nil 的空切片。二者都能用于 `len`、`range` 和 `append`，但与 `nil` 比较的结果不同。

Go 1.21 引入的内置函数 `clear` 会把切片现有长度范围内的元素重置为零值，但不会改变长度和容量。

### 11.2 修改

```go
slice1 := []int{1, 2, 3}
slice2 := append(slice1, 4, 5)      // 追加两个元素
slice3 := append(slice2, slice1...) // 追加一个切片的所有元素
fmt.Println(slice1)                 // [1 2 3]
fmt.Println(slice2)                 // [1 2 3 4 5]
fmt.Println(slice3)                 // [1 2 3 4 5 1 2 3]
```

::: warning 注意
将另一个切片的所有元素追加到一个切片上时，需要在另一个切片后面写上 `...`。`append` 返回新的切片值，需要接收这个返回值，原变量的长度不会自行改变。
:::

```go
// 删除切片中索引为 2 的元素
slice1 := []int{1, 2, 3, 4, 5}
slice2 := append(slice1[:2], slice1[3:]...)
fmt.Println(slice1) // [1 2 4 5 5]
fmt.Println(slice2) // [1 2 4 5]
```

::: warning 注意
`slice1` 变成了 `[1 2 4 5 5]`，这是因为切片是对底层数组的引用，当执行 `append(slice1[:2], slice1[3:]...)` 时，`slice1` 的底层数组会被修改。具体来说，`slice1` 的前 `2` 个元素保持不变，但后面的元素会被覆盖为 `[4, 5]`。
:::

```go
slice1 := []int{1, 2, 3, 4, 5, 6}
slice2 := []int{}
slice3 := make([]int, 3)
copy(slice2, slice1[1:])
copy(slice3, slice1[1:])
slice1[1] = 100
fmt.Println(slice1) // [1 100 3 4 5 6]
fmt.Println(slice2) // []，因为长度为 0，所以复制了 0 个元素
fmt.Println(slice3) // [2 3 4]，因为长度为 3，所以复制了 3 个元素
```

`copy(dst, src)` 复制的元素个数为 `min(len(dst), len(src))`，并返回实际复制数量。它不会扩展目标切片，即使目标容量足够，长度为 0 时也不会复制任何元素。源和目标可以重叠。上例的 `slice3` 由 `make` 创建，拥有独立的底层数组，因此修改 `slice1` 不会影响它。

```go
func handle(nums []int) {
	for index := range nums {
		nums[index] *= 10
		nums = append(nums, index)
	}
}

func main() {
	slice := []int{1, 2, 3}
	fmt.Println(slice) // [1 2 3]
	handle(slice)
	fmt.Println(slice) // [10 2 3]
}
```

::: warning 注意
切片参数按值传递，复制后的切片描述符最初指向同一个底层数组。上例的初始长度和容量都是 3，第一次修改会影响调用方，紧接着的第一次 `append` 必须分配新数组，后续修改便发生在新数组上。`range` 的迭代次数由进入循环时的切片长度决定，因此这里只迭代三次。一般情况下，追加后的长度超过容量才需要分配新数组，具体增长比例不属于语言规范的保证。
:::

### 11.3 底层

可以把切片概念化为下面三个字段，但这不是语言规范公开的结构体定义：

```go
type slice[T any] struct {
	ptr *T  // 指向底层数组的指针
	len int // 切片的长度
	cap int // 切片的容量
}
```

切片本身并不存储数据，而是引用一个底层数组。

## 12. Map

`map` 是可变的键值集合，键的类型必须可比较，切片、map 和函数不能作为键的类型。使用接口类型作为键时，存入的动态值也必须可比较，否则会触发 `panic`。`range` 的遍历顺序没有规定，同一个 `map` 的多次遍历也可能得到不同顺序。没有写操作时可以并发读取，涉及并发写入时需要使用锁、单 goroutine 管理或适合相应场景的 `sync.Map`。

`fmt` 格式化 map 时会按键排序，这与 `range` 的遍历规则不同。下面的字符串键会按字节的字典序输出。

### 12.1 创建

```go
var myMap = map[string]string{
	"key1": "value1",
	"key2": "value2",
}
myMap["key3"] = "value3"
fmt.Println(myMap) // map[key1:value1 key2:value2 key3:value3]
```

下面这种用法将会报错：

```go
var myMap map[string]string
fmt.Println(myMap == nil) // true
myMap["key"] = "value"    // 报错：panic: assignment to entry in nil map
```

nil map 可以读取、遍历，也可以用于 `delete` 和 `clear`，但不能写入。通过 `make` 创建空 map 后就可以写入：

```go
var myMap = make(map[string]string)
fmt.Println(myMap == nil) // false
myMap["key"] = "value"
fmt.Println(myMap) // map[key:value]
```

### 12.2 遍历

```go
myMap := map[string]string{"key1": "value1", "key2": "value2"}
for key, value := range myMap {
	fmt.Println(key, value)
}
for key := range myMap {
	fmt.Println(key)
}
```

### 12.3 删除

```go
var myMap = map[string]string{
	"key1": "value1",
	"key2": "value2",
	"key3": "value3",
}
fmt.Println(myMap) // map[key1:value1 key2:value2 key3:value3]
delete(myMap, "key2")
fmt.Println(myMap)    // map[key1:value1 key3:value3]
delete(myMap, "key4") // 删除不存在的元素也不会报错
fmt.Println(myMap)    // map[key1:value1 key3:value3]
```

Go 1.21 起，需要删除全部键时可以调用 `clear(myMap)`，已初始化的 map 仍可继续写入。对 nil map 调用 `clear` 不会将它初始化。

### 12.4 查询

```go
var myMap = map[string]string{
	"key1": "value1",
	"key2": "value2",
	"key3": "value3",
}
key1, ok1 := myMap["key2"]
key2, ok2 := myMap["key4"]
fmt.Println(ok1, key1) // true value2
fmt.Println(ok2, key2) // false (空字符串)
```

## 13. 函数

### 13.1 声明

```go
func add(m, n int) (sum int, err error) {
	sum = m + n
	return
}

// 效果同上
func add1(m int, n int) (int, error) {
	return m + n, nil
}

// 可变参数
func add2(slice ...int) (sum int, err error) {
	for _, value := range slice {
		sum += value
	}
	return
}

// 返回值为函数
func getFunc() (getN func() int) {
	getN = func() (n int) {
		n = 10
		return
	}
	return
}

func main() {
	fmt.Println(add(1, 2))     // 3 <nil>
	fmt.Println(add2(1, 2, 3)) // 6 <nil>
	fmt.Println(getFunc()())   // 10
}
```

### 13.2 闭包

```go
func autoIncrement() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {
	nextNum := autoIncrement()
	fmt.Println(nextNum()) // 1
	fmt.Println(nextNum()) // 2
	fmt.Println(nextNum()) // 3

	nextNum = autoIncrement()
	fmt.Println(nextNum()) // 1
	fmt.Println(nextNum()) // 2
}
```

### 13.3 `defer`

多个 `defer` 按照 `LIFO` 的顺序执行：

```go
func deferPrint() int {
	defer fmt.Print("1")
	defer fmt.Print("2")
	defer fmt.Print("3")
	return 0
}

func main() {
	fmt.Print(deferPrint()) // 3210
}
```

延迟调用的函数值和实参在执行 `defer` 语句时就已经求值，但调用本身延后执行：

```go
n := 1
defer fmt.Println(n) // 1
n++
fmt.Println(n) // 2
```

闭包直接引用的外层变量并不会在注册 `defer` 时自动复制。下面会读取执行延迟函数时的 `n`：

```go
n := 1
defer func() { fmt.Println(n) }() // 2
n++
```

函数返回时，会先给返回值赋值，再按后进先出的顺序执行延迟调用，最后把结果交给调用方。因此，延迟函数可以修改外层函数的具名返回值：

```go
func outer() (result int) {
	defer func() {
		result *= 2
	}()
	return 10
}

func main() {
	fmt.Println(outer()) // 20
}
```

上述代码执行 `return` 语句时，先给具名返回值赋值 `result = 10`，然后执行延迟函数的内容 `result *= 2`，最后返回 `result`。

### 13.4 `panic` 和 `recover`

`panic` 会停止当前函数的正常执行，并沿调用栈运行已经注册的 `defer`。如果一直没有被 `recover`，当前 goroutine 的 panic 最终会终止程序。

```go
func setAge(age int) {
	if age < 0 {
		panic("negative age")
	}
}

func main() {
	setAge(-12)    // panic: negative age
	fmt.Println(0) // 执行不到这里
}
```

`recover` 必须由延迟函数直接调用，才能捕获同一 goroutine 正在传播的 panic。通过另一层普通函数间接调用 `recover` 无法恢复。它可用于 goroutine 入口或请求处理入口的异常隔离，普通业务错误仍应返回 `error`。

```go
func setAge(age int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r) // negative age
		}
	}()
	if age < 0 {
		panic("negative age")
	}
}

func main() {
	setAge(-12)
	fmt.Println(0) // 0
}
```

## 14. `type`

```go
type myInt1 = int // 类型别名，与 int 是同一个类型
type myInt2 int   // 定义一个新类型，底层类型为 int

func main() {
	var a myInt1
	var b myInt2
	fmt.Printf("%T\n", a) // int
	fmt.Printf("%T\n", b) // main.myInt2
}
```

类型别名不会创建新类型。类型定义会创建不同的类型，赋值给底层类型相同的其他定义类型时通常需要显式转换，也不会自动继承原类型的方法。

## 15. 结构体

### 15.1 声明

```go
type Student struct {
	name string
	age  int
}
```

### 15.2 访问

```go
stu1 := Student{"Zhang", 20}
stu2 := Student{name: "Klose"}
fmt.Println(stu1.name, stu1.age) // Zhang 20
fmt.Println(stu2.name, stu2.age) // Klose 0
```

### 15.3 嵌套

```go{7,13}
type Info struct {
	name string
	age  int
}

type Student struct {
	info  Info
	score float64
}

func main() {
	stu := Student{Info{"Zhang", 20}, 95.5}
	fmt.Println(stu.info.name) // Zhang
}
```

### 15.4 匿名嵌套

嵌入字段可以提升其字段和方法，使它们能够通过外层结构体直接访问。结构体字面量中不能直接用被提升的字段名作为键，仍然要初始化对应的嵌入字段。

```go{7,13}
type Info struct {
	name string
	age  int
}

type Student struct {
	Info
	score float64
}

func main() {
	stu := Student{Info{"Zhang", 20}, 95.5}
	fmt.Println(stu.name) // Zhang
}
```

同名字段优先选择嵌入层次更浅的一项，因此外层字段会遮蔽内层字段。如果多个候选字段处于同一最浅层次，直接访问会因歧义而无法编译。

```go{2,8}
type Info struct {
	name string
	age  int
}

type Student struct {
	Info
	name  string
	score float64
}

func main() {
	stu := Student{Info{"Zhang", 20}, "Klose", 95.5}
	fmt.Println(stu.name) // Klose
}
```

### 15.5 方法

```go{6-8}
type Student struct {
	name string
	age  int
}

func (stu Student) print() {
	fmt.Printf("name: %s, age: %d\n", stu.name, stu.age)
}

func main() {
	stu := Student{"zhang", 20}
	stu.print() // name: zhang, age: 20
}
```

上例使用值接收者，调用时会复制 `Student` 值，修改接收者本身不会改变原结构体。需要修改原结构体时可以使用 `*Student` 指针接收者。若结构体字段本身是切片、map 或指针，复制结构体后仍可能共享这些字段引用的数据。

## 16. 指针

```go
func increase(n *int) {
	(*n)++
}

func main() {
	n := 15
	increase(&n)
	fmt.Println(n) // 16
}
```

未初始化的指针为 `nil`，可以传递和比较，但不能解引用。`new(T)` 为类型 `T` 的零值分配存储并返回 `*T`，并不保证存储一定分配在堆上：

```go
var n1 *int
var n2 = new(int)
fmt.Printf("type: %T, n1==nil: %t\n", n1, n1 == nil) // type: *int, n1==nil: true
fmt.Printf("type: %T, n2==nil: %t\n", n2, n2 == nil) // type: *int, n2==nil: false
```

Go 1.26 起，`new` 还可以接收表达式，将其值保存到新分配的存储中：

```go
p := new(42)
fmt.Println(*p) // 42
```

## 17. 接口

普通接口描述一组方法。Go 中没有关键字显式声明某个类型实现了接口，只要该类型的方法集包含接口要求的全部方法，就实现了该接口。泛型约束还可以使用类型集合，这类非基本接口只能用于约束，不能直接作为普通变量的类型。

### 17.1 接口定义

```go
type Duck interface {
	walk()
	eat()
	sleep()
}
```

### 17.2 接口实现

```go
type PskDuck struct {
	age uint8
}

func (p *PskDuck) walk() {
	fmt.Println("pskDuck walk")
}
func (p *PskDuck) eat() {
	fmt.Println("pskDuck eat")
}
func (p *PskDuck) sleep() {
	fmt.Println("pskDuck sleep")
}

func main() {
	var pskDuck Duck = &PskDuck{age: 1}
	pskDuck.walk()
	pskDuck.eat()
	pskDuck.sleep()
}
```

上例的方法全部使用指针接收者，因此实现 `Duck` 的是 `*PskDuck`，`PskDuck` 本身没有实现。可寻址的结构体变量调用指针接收者方法时，编译器可以隐式取地址，但接口赋值不会这样转换。

### 17.3 空接口

空接口 `interface{}` 不要求任何方法，因此任意类型都满足它。Go 1.18 起通常写成等价的别名 `any`。

### 17.4 类型断言

类型断言 `x.(T)` 用于检查接口值的动态类型，并取得对应的值。`T` 是具体类型时，动态类型必须与它相同。`T` 是接口时，动态类型必须实现该接口。单返回值形式在检查失败时会触发 `panic`。

```go
var a any = "zhang"
name := a.(string)
fmt.Println(name) // zhang
age := a.(int)    // 报错：panic: interface conversion: interface {} is string, not int
fmt.Println(age)
```

为了避免 `panic`，可以使用带检查的类型断言。失败时得到目标类型的零值和 `false`：

```go
var a any = "zhang"
age, isInt := a.(int)
if isInt {
	fmt.Println("age:", age)
} else {
	fmt.Println("a:", a) // a: zhang
}
```

### 17.5 类型选择

```go
var i any = "zhang"
switch i.(type) {
case nil:
	fmt.Println("nil")
case int:
	fmt.Println("int")
case string:
	fmt.Println("string") // string
default:
	fmt.Println("unknown")
}
```

### 17.6 接口遇到切片的常见错误

```go
func printSlice(slice ...any) {
	for _, v := range slice {
		fmt.Println(v)
	}
}

func main() {
	data := []string{"zhang", "heng", "hua"}
	printSlice(data...) // 编译错误：[]string 不能作为 []any 展开传递
}
```

单个 `string` 值可以赋给 `any`，但 `[]string` 与 `[]any` 是不同的切片类型，不能整体赋值或直接展开传递。需要新建 `[]any`，逐项转换后再调用：

```go
data := []string{"zhang", "heng", "hua"}
values := make([]any, len(data))
for i, value := range data {
	values[i] = value
}
printSlice(values...) // 依次输出 zhang、heng、hua
```

### 17.7 `error` 接口

`error` 内置接口类型的源码如下：

```go
type error interface {
	Error() string
}
```

自定义错误：

```go
type newError struct{}

func (e *newError) Error() string {
	return "新错误"
}

func main() {
	err := &newError{}
	fmt.Println(err) // 新错误
}
```

接口只有在动态类型和动态值都为空时才等于 `nil`。包含 nil 指针的接口仍然不是 nil，这也是返回 `error` 时常见的问题：

```go
var p *newError
var err error = p
fmt.Println(p == nil)   // true
fmt.Println(err == nil) // false
```

函数成功时应返回 `nil` 接口，而不是将某个 nil 错误指针转换为 `error` 后返回。

## 18. `package`

假设项目目录为 `proj`，模块路径为 `example.com/proj`，`go.mod` 中包含：

```text
module example.com/proj

go 1.26.0
```

`user/person.go` 定义如下。目录名决定导入路径，`package` 声明决定默认导入名称，两者不要求相同：

```go
package person

type Person struct {
	Name string
}
```

名称以 Unicode 大写字母开头的包级标识符和字段可以被其他包访问，例如 `Person` 和 `Name`。下面几种导入方式分别演示，不要同时放在一个文件中。

### 18.1 导入

```go
import (
	"fmt"

	"example.com/proj/user"
)

func main() {
	p := person.Person{Name: "Zhang"}
	fmt.Println(p.Name) // Zhang
}
```

也可以给包起个别名：

```go{4}
import (
	"fmt"

	u "example.com/proj/user"
)

func main() {
	p := u.Person{Name: "Zhang"}
	fmt.Println(p.Name) // Zhang
}
```

点导入会把目标包的导出标识符直接引入当前文件，容易造成来源不明和命名冲突，普通业务代码不建议使用：

```go{4}
import (
	"fmt"

	. "example.com/proj/user"
)

func main() {
	p := Person{Name: "Zhang"}
	fmt.Println(p.Name) // Zhang
}
```

空白导入只执行目标包的初始化副作用，常用于注册数据库驱动、图片解码器等实现：

```go{1}
import _ "example.com/proj/user"

func main() {}
```

### 18.2 `init`

被导入的包先于导入它的包初始化，每个包在一个程序中只初始化一次。包内先按依赖关系初始化包级变量，再执行 `init` 函数，全部初始化完成后才进入 `main.main`。同一包中 `init` 的执行顺序取决于声明顺序及文件提交给编译器的顺序，标准 Go 工具通常按文件名排序。需要明确顺序或返回错误的复杂初始化，更适合放进显式函数。

例如，给 `user/person.go` 加上 `init`：

```go{10-12}
package person

import "fmt"

type Person struct {
	Name string
}

// 在当前包初始化时自动执行
func init() {
	fmt.Println("init")
}
```

## 19. `go modules`

### 19.1 初始化模块

在新项目的根目录初始化模块，创建 `go.mod`。模块路径是包导入路径的前缀，可以与本地目录名不同。如果前面已经创建了 `go.mod`，不需要再次执行 `go mod init`。

```shell
go mod init example.com/proj
go get go@1.26.0
```

Go 1.26.0 的 `go mod init` 默认写入 `go 1.25.0`，Go 1.26.1 起已恢复为当前工具链版本，例如 Go 1.26.8 写入 `go 1.26.8`。这项补丁变更记录在官方 [#77860](https://github.com/golang/go/issues/77860) 中。上面的第二条命令显式将模块的最低 Go 版本设为 `1.26.0`，使这些示例采用一致的语言版本。最低版本与本机工具链补丁版本是两个概念，可以使用 Go 1.26.8 工具链构建这个模块。

### 19.2 添加依赖

`go get` 解析指定包所需的模块版本，并更新 `go.mod` 和相应的校验记录。显式指定版本可以避免示例随最新版本变化：

```shell
go get github.com/gin-gonic/gin@v1.12.0
```

### 19.3 移除未使用的依赖

`go mod tidy` 根据源码和测试中的导入补充缺失的依赖，并移除不再需要的依赖及校验记录。它会考虑大多数构建标签下的源码。如果刚用 `go get` 添加 Gin，却还没有任何源码导入它，执行 `tidy` 会将该依赖移除。

```shell
go mod tidy
```

### 19.4 查看依赖关系

列出当前模块和构建列表中选定的依赖模块版本。这不是列出每个依赖的全部历史版本。

```shell
go list -m all
```

显示模块依赖图中的版本需求关系。依赖边表示某个模块要求的版本，最终选用的版本仍以构建列表为准。

```shell
go mod graph
```

### 19.5 下载依赖

提前把依赖下载到本地模块缓存。对本文使用的模块版本，不带参数时下载 `go.mod` 中显式要求的模块，包括标为 `indirect` 的条目。平时执行 `go build`、`go test` 等命令也会按需下载依赖。

```shell
go mod download
```

### 19.6 校验依赖

```shell
go mod verify
```

该命令比较模块缓存中的压缩包、解压目录与下载时保存在缓存中的哈希，检查下载后内容是否被修改。它不是直接拿依赖源码与 `go.sum` 比较。构造模块图时仍可能下载 `go.mod` 并按 `go.sum` 校验或补充记录。

## 20. 单元测试

Go 的单元测试主要依赖于 `testing` 包，并且通过 `go test` 命令来执行测试。

### 20.1 基本结构

测试文件必须以 `_test.go` 结尾。测试可以与被测代码使用同一包，也可以使用以 `_test` 结尾的外部测试包。前者能访问未导出标识符，后者更接近真实调用方。

测试函数不带类型参数和返回值，名称以 `Test` 开头，后面的首个字符不能是小写字母，并接受一个 `*testing.T` 参数。例如：

```go
func TestAdd(t *testing.T) {}
```

### 20.2 编写测试用例

假设有一个简单的 `math` 包，包含一个 `Add` 函数：

```go
// math/math.go
package math

func Add(a, b int) int {
	return a + b
}
```

对应的测试文件 `math_test.go` 如下：

```go
// math/math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5
	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}
```

### 20.3 运行测试

在包含测试文件的目录下，运行以下命令执行当前包的测试：

```shell
go test
```

在模块根目录运行 `go test ./...` 可以测试当前模块下的所有包。使用包列表测试时可能复用成功的测试缓存，需要实际重跑时可以加 `-count=1`。

使用 `-v` 标志可以查看详细的测试输出：

```shell
go test -v
```

`-run` 接收正则表达式，下面只选择名称为 `TestAdd` 的测试及其子测试：

```shell
go test -run '^TestAdd$'
```

### 20.4 性能测试

Go 支持基准测试，用于衡量代码的性能。基准测试函数不带类型参数和返回值，以 `Benchmark` 开头，后面的首个字符不能是小写字母，并接受一个 `*testing.B` 参数。下面采用 Go 1.24 引入的 `B.Loop`：

```go
// math/math_test.go
func BenchmarkAdd(b *testing.B) {
	for b.Loop() {
		Add(1, 2)
	}
}
```

使用 `-bench` 选择基准测试，`-benchmem` 报告每次操作的内存分配，`-run='^$'` 排除普通测试：

```shell
go test -run='^$' -bench=. -benchmem
```

`B.Loop` 由测试框架决定迭代次数，在第一次调用时重置计时器，返回 `false` 时停止计时。因此，循环前的准备和循环后的清理不计入结果，循环控制本身的开销仍计入结果。按 `for b.Loop() { ... }` 的形式编写时，编译器会保留循环体内调用的参数、结果及赋值变量，避免被测代码整体消除。也可以指定特定的基准测试：

```shell
go test -run='^$' -bench='^BenchmarkAdd$' -benchmem
```

### 20.5 跳过用例

```go{2-4}
func TestAdd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	result := Add(2, 3)
	expected := 5
	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}
```

`-short` 只使 `testing.Short()` 返回 `true`，不会自动跳过测试。上例主动调用了 `t.Skip`，因此下面的命令会跳过它：

```shell
go test -short
```

### 20.6 表格驱动测试

```go
func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"positive", 6, 2, 8},
		{"zero", 5, 0, 5},
		{"negative", -6, 2, -4},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Add(test.a, test.b)
			if result != test.expected {
				t.Errorf("Add(%d, %d) = %d, want %d", test.a, test.b, result, test.expected)
			}
		})
	}
}
```

## 21. 泛型

### 21.1 函数上使用

在函数名后声明类型参数及其约束。下面的 `T` 只能取约束中列出的三种类型：

```go
func Add[T int | float64 | string](a, b T) T {
	return a + b
}

func main() {
	fmt.Println(Add[int](1, 2))            // 3
	fmt.Println(Add[float64](1.1, 2.1))    // 3.2
	fmt.Println(Add[string]("1.1", "2.2")) // 1.12.2
}
```

调用时通常可以推断类型参数，例如 `Add(1, 2)`。当前约束只包含列出的类型，不包含 `type MyInt int` 这样的定义类型。如果希望接受底层类型相同的定义类型，可以将约束写为 `~int | ~float64 | ~string`。

用 `any` 和类型断言也能实现相似效果，但会失去编译期类型检查，而且调用时可能因参数类型不一致而 panic：

```go
func IAdd(a, b interface{}) interface{} {
	switch a.(type) {
	case int:
		return a.(int) + b.(int)
	case float64:
		return a.(float64) + b.(float64)
	case string:
		return a.(string) + b.(string)
	}
	return nil
}

func main() {
	fmt.Println(IAdd(1, 2))         // 3
	fmt.Println(IAdd(1.1, 2.1))     // 3.2
	fmt.Println(IAdd("1.1", "2.2")) // 1.12.2
}
```

### 21.2 Map 上使用

```go
type MyMap[
	K int | string,
	V float32 | float64,
] map[K]V

func main() {
	m := MyMap[string, float64]{}
	fmt.Println(m) // map[]
}
```

这里的 `int` 和 `string` 都可比较，因此能作为 map 的键。需要接受任意可比较键类型时，可以使用约束 `K comparable`。

### 21.3 结构体上使用

```go
type S[
	T1 string | float64,
	T2 int | uint,
] struct {
	A T1
	B T2
}

func main() {
	s := S[string, uint]{"Hello", 2025}
	fmt.Println(s) // {Hello 2025}
}
```

## 22. 函数选项模式

函数选项模式适合可选参数较多、默认值稳定且未来可能继续扩展的构造函数。参数很少时，直接使用普通参数或配置结构体通常更清楚。

例如，有如下结构体表示数据库连接配置：

```go
type DBOptions struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}
```

定义一个函数类型 `Option`，它接收一个指向配置结构体的指针，目的是修改该结构体的字段：

```go
type Option func(*DBOptions)
```

声明选项函数，即对配置结构体进行某项定制化设置的函数。例如：

```go
func WithHost(host string) Option {
	return func(o *DBOptions) {
		o.Host = host
	}
}

func WithPort(port int) Option {
	return func(o *DBOptions) {
		o.Port = port
	}
}
```

构造函数中，先设定默认值，然后遍历用户提供的选项函数，依次修改配置：

```go
func NewDBOptions(options ...Option) *DBOptions {
	dbOptions := &DBOptions{
		Host: "127.0.0.1",
		Port: 3306,
	}
	for _, option := range options {
		option(dbOptions)
	}
	return dbOptions
}

func main() {
	// 只修改了 Host 字段，其它使用默认值
	dbOptions := NewDBOptions(WithHost("192.168.0.1"))
	fmt.Println(dbOptions.Host) // 192.168.0.1
	fmt.Println(dbOptions.Port) // 3306
}
```

这个示例只创建配置，没有建立数据库连接。选项按传入顺序执行，同一字段多次设置时，后一次覆盖前一次。真实客户端的构造函数通常还要校验配置，并在配置无效时返回 `error`。

## 23. 错误处理（`error`）

### 23.1 什么是 `error`

在 Go 中，`error` 是一种内建接口，用于表示函数执行中的错误状态。

```go
type error interface {
	Error() string
}
```

任何实现了 `Error() string` 方法的类型都可以被视为一个 `error`。

### 23.2 返回 `error` 的基本用法

可能失败的函数通常把 `error` 作为最后一个返回值：

```go
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}
```

调用时进行错误判断：

```go
result, err := divide(10, 0)
if err != nil {
	fmt.Println("Error:", err)
} else {
	fmt.Println("Result:", result)
}
```

### 23.3 创建错误的几种方式

标准库的 `errors.New` 创建一个只包含消息的错误值，它不是内置函数：

```go
err := errors.New("something went wrong")
fmt.Println(err)
```

`fmt.Errorf` 可以将参数格式化为错误消息：

```go
input := -1
err := fmt.Errorf("invalid input: %d", input)
fmt.Println(err) // invalid input: -1
```

Go 1.13 起，`fmt.Errorf` 的 `%w` 可以在添加上下文的同时包装原错误。调用方仍能通过 `errors.Is`、`errors.As` 检查原错误，而 `%v` 只保留格式化后的文字：

```go
baseErr := errors.New("disk error")
err := fmt.Errorf("upload failed: %w", baseErr)
fmt.Println(errors.Is(err, baseErr)) // true
```

### 23.4 判断和提取底层错误

`errors.Is(err, target)` 检查错误及其包装的错误是否与目标匹配，通常用于判断哨兵错误，而不是判断类型。匹配可以通过值相等或错误类型提供的 `Is(error) bool` 方法实现。两次 `errors.New` 即使消息相同，也会产生不同的错误值。

`errors.As` 查找可赋值给指定类型的错误，并通过指针参数提取它。错误类型也可以通过 `As(any) bool` 自定义匹配行为。

```go
_, err := os.ReadFile("missing.txt")
if errors.Is(err, os.ErrNotExist) {
	fmt.Println("File does not exist")
}

var pathErr *os.PathError
if errors.As(err, &pathErr) {
	fmt.Println("Path error:", pathErr.Path)
}
```

Go 1.26 新增 `errors.AsType`，可以直接返回匹配的错误值，减少单独声明目标变量的步骤：

```go
if pathErr, ok := errors.AsType[*os.PathError](err); ok {
	fmt.Println("Path error:", pathErr.Path)
}
```

这些函数也支持 Go 1.20 引入的 `errors.Join` 和包含多个 `%w` 的错误包装，此时错误关系可以形成树，而不只是单条链。

### 23.5 错误链与调用栈

标准库的错误包装及检查函数不会自动保存错误发生时的调用栈。`runtime/debug.Stack` 获取的是调用它的 goroutine 此刻的栈，不能还原已经返回的函数。下面的示例先打印错误及其匹配结果，再打印当前 goroutine 的栈，其中不会包含已经返回的 `NewStudent` 和 `SetName`。如果需要保存错误发生位置，应在该位置采集栈信息。

<<< @/go/codes/golang/error_stack.go

## 24. AST 代码生成

### 24.1 使用 stringer 生成方法

假设我们要维护状态码的相关代码，可能会这样来写：

```go
type Code int64

const (
	OK            Code = 0 // OK
	InvalidParams Code = 1 // 参数错误
	Timeout       Code = 2 // 超时
)

var mapCodeDesc = map[Code]string{
	OK:            "OK",
	InvalidParams: "参数错误",
	Timeout:       "超时",
}

func GetCodeDesc(code Code) string {
	if desc, ok := mapCodeDesc[code]; ok {
		return desc
	}
	return "未知错误"
}
```

手动维护常量和描述映射容易遗漏。`stringer` 可以为整数类型生成 `String() string` 方法，使用 `-linecomment` 时以常量后的行注释作为描述。这里使用 `golang.org/x/tools` v0.51.0，该版本要求 Go 1.26 或更高版本：

```shell
go install golang.org/x/tools/cmd/stringer@v0.51.0
```

可执行文件会安装到 `GOBIN`，未设置时默认使用 `$(go env GOPATH)/bin`。需要将对应目录加入 `PATH`，以便 `go generate` 找到 `stringer`。

在前面示例的模块中新建 `code` 目录，将原来的手动映射替换为下面的 `code/code.go`：

<<< @/go/codes/golang/code.go

在模块根目录执行：

```shell
go generate ./code
```

`go generate` 执行指令注释中的命令，生成 `code/code_string.go`。它不会由 `go build` 或 `go test` 自动执行，修改状态码后需要再次生成。生成文件应由工具维护，不要手工编辑：

<<< @/go/codes/golang/code_string.go

调用方可以直接使用 `code.InvalidParams.String()`，输出 `参数错误`。`fmt.Println(code.Timeout)` 也会调用这个方法，输出 `超时`。未定义的值，例如 `code.Code(99)`，会得到 `Code(99)`，与原来手写函数的 `未知错误` 返回值不同。

### 24.2 什么是 AST

AST（Abstract Syntax Tree，抽象语法树）以树状节点表示源代码结构。Go 的 `go/ast` 定义节点，`go/token` 管理词法符号和源码位置。生成器可以构造这些节点，再输出为 `.go` 文件，之后仍需通过正常的 Go 工具链编译。

### 24.3 常用标准库

| 包名         | 用途                     |
| ------------ | ------------------------ |
| `go/ast`     | 抽象语法树节点结构和操作 |
| `go/token`   | 词法符号和源码位置       |
| `go/parser`  | 将源代码解析为 AST       |
| `go/printer` | 将 AST 输出为源代码      |
| `go/format`  | 按 Go 格式规范输出源码   |

### 24.4 基本步骤

下面生成一个 `Hello` 函数。将四个步骤的代码按顺序放入模块中的 `hello/gen/main.go`。生成器使用 `package main`，生成的库文件使用 `package hello`，二者放在不同目录。

#### 24.4.1 步骤一：导入包

```go
package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"os"
)
```

#### 24.4.2 步骤二：构造代码结构

目标函数如下，它是生成结果，不要复制到生成器文件中：

```go
func Hello(name string) string {
	return "Hello, " + name
}
```

用 AST 构造等效结构：

```go
func buildFunc() ast.Decl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("Hello"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("name")},
						Type:  ast.NewIdent("string"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BinaryExpr{
							X:  &ast.BasicLit{Kind: token.STRING, Value: `"Hello, "`},
							Op: token.ADD,
							Y:  ast.NewIdent("name"),
						},
					},
				},
			},
		},
	}
}
```

#### 24.4.3 步骤三：组合为完整文件

```go
func buildFile() *ast.File {
	return &ast.File{
		Name:  ast.NewIdent("hello"),
		Decls: []ast.Decl{buildFunc()},
	}
}
```

#### 24.4.4 步骤四：输出生成的代码到 `.go` 文件

```go
func main() {
	fset := token.NewFileSet()
	f := buildFile()

	var output bytes.Buffer
	if err := format.Node(&output, fset, f); err != nil {
		panic(err)
	}
	if err := os.WriteFile("hello_gen.go", output.Bytes(), 0o644); err != nil {
		panic(err)
	}
}
```

从模块根目录进入 `hello` 目录，再运行生成器。`os.WriteFile` 的相对路径以命令的工作目录为准：

```shell
cd hello
go run ./gen
go test .
```

生成的 `hello/hello_gen.go` 内容如下。`go test .` 可以检查这个库包能否编译，此时还没有测试用例。调用方可导入 `example.com/proj/hello` 并使用 `hello.Hello("Go")`，结果为 `Hello, Go`。

<<< @/go/codes/golang/hello_gen.go

### 24.5 AST 节点常见结构

| AST 类型          | 描述                             | 示例            |
| ----------------- | -------------------------------- | --------------- |
| `*ast.File`       | 表示一个 Go 文件                 | 整体文件结构    |
| `*ast.FuncDecl`   | 表示函数声明                     | `func Foo() {}` |
| `*ast.GenDecl`    | 常量/变量/类型声明               | `var x int`     |
| `*ast.AssignStmt` | 赋值语句                         | `x := 1`        |
| `*ast.ReturnStmt` | 返回语句                         | `return x`      |
| `*ast.BinaryExpr` | 二元运算符                       | `x + y`         |
| `*ast.BasicLit`   | 字面量                           | `"hello"`, `1`  |
| `*ast.Ident`      | 标识符，常用 `ast.NewIdent` 创建 | `x`, `Foo`      |
