---
outline: deep
---

# Go

Go 是一门静态类型、编译型语言，带有垃圾回收、结构化并发工具和完整的标准工具链。本文按语言基础、数据结构、接口、模块、测试和代码生成的顺序整理，示例以 Go 1.26 为基线。

## 1 注释

Go 支持行注释和块注释。普通说明优先使用行注释。块注释更适合暂时注释一段代码，或者在表达式中插入说明。

### 1.1 单行注释

单行注释以两个正斜杠（`//`）开头。

编译器将忽略 `//` 到行尾之间的文本（不会执行）。

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

编译器将忽略 `/*` 和 `*/` 之间的文本。

```go
package main

import "fmt"

func main() {
 /* 这是一个
 多行注释 */
 fmt.Println("Hello World!")
}

```

## 2 变量

### 2.1 声明变量

常用的变量声明有两种：

1. 使用 `var` 关键字：  

`var` 关键字后面跟变量名和数据类型：

```go
var age int = 1
```

::: warning 注意
数据类型和值必须指定任意一项，或者像上面那样两个都指定。

```go
var age1 int
var age2 = 1
```

:::

1. 使用 `:=` 符号：

`:=` 后面接变量值：

```go
age := 1
```

::: warning 注意

- 编译器会根据右侧表达式推断变量类型。
- `:=` 只能在函数体内使用，且左侧至少要有一个新变量。
:::

### 2.2 声明多变量

在 Go 中，可以在同一行声明多个变量。

```go
var a, b, c int = 1, 2, 3
```

::: warning 注意
如果指定了数据类型，则每行只能声明同一种数据类型的多个变量。
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

## 3 常量

通过 `const` 关键字声明常量，常量的值是无法修改的。

```go
const PI float64 = 3.1415926
```

::: warning 注意
必须在声明常量时赋值。
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

## 4 基本数据类型

### 4.1 布尔

```go
var b1 bool = true
var b2 = true
var b3 bool
b4 := true
```

### 4.2 整数

整数类型分为两类：

- **有符号整数** - 可以存储正值和负值
- **无符号整数** - 只能存储非负值

#### 4.2.1 有符号整数

| 类型    |                          大小                          |           范围            |
| ------- | :----------------------------------------------------: | :-----------------------: |
| `int`   | 32 bits 或者 64 bits<br />取决于是32位系统还是64位系统 |                           |
| `int8`  |                    8 bits / 1 byte                     |  - (2 ^ 7) ~ (2 ^ 7 - 1)  |
| `int16` |                   16 bits / 2 bytes                    | - (2 ^ 15) ~ (2 ^ 15 - 1) |
| `int32` |                   32 bits / 4 bytes                    | - (2 ^ 31) ~ (2 ^ 31 - 1) |
| `int64` |                   64 bits / 8 bytes                    | - (2 ^ 63) ~ (2 ^ 63 - 1) |

#### 4.2.2 无符号整数

| 类型     |                          大小                          | 范围             |
| -------- | :----------------------------------------------------: | :--------------- |
| `uint`   | 32 bits 或者 64 bits<br />取决于是32位系统还是64位系统 |                  |
| `uint8`  |                    8 bits / 1 byte                     | 0 ~ (2 ^ 8 - 1)  |
| `uint16` |                   16 bits / 2 bytes                    | 0 ~ (2 ^ 16 - 1) |
| `uint32` |                   32 bits / 4 bytes                    | 0 ~ (2 ^ 32 - 1) |
| `uint64` |                   64 bits / 8 bytes                    | 0 ~ (2 ^ 64 - 1) |

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

浮点数类型有两个关键字：`float32`、`float64`。

### 4.4 字符串

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

fmt.Println(s1)  // name: Zhang, age: 23
fmt.Println(s2)  // name: Zhang, age: 23
fmt.Println(s3)  // name: Zhang, age: 23
```

不存在对所有场景都最快的拼接方式。少量固定片段优先使用 `+`；循环拼接或片段较多时再考虑 `strings.Builder`。如果已经有字符串切片，`strings.Join` 往往更直接。

#### 4.4.2 字符串的比较

直接使用比较运算符就行。

```go
s1 := "Zhang"
s2 := "Zhang"

fmt.Println(s1 == s2)  // true
```

## 5 格式化输出

### 5.1 常规格式

以下格式可用于所有数据类型：

| 格式  | 含义               |
| ----- | ------------------ |
| `%v`  | 以默认格式输出     |
| `%+v` | 在打印结构体时，会显示字段名和字段值 |
| `%#v` | 以Go语法的格式输出 |
| `%T`  | 输出值的类型       |
| `%%`  | 输出一个百分号     |

```go
f := 12.3
s := "Zhang"
a := []int{1, 2, 3}

fmt.Printf("%v\n", f)   // 12.3
fmt.Printf("%#v\n", f)  // 12.3
fmt.Printf("%T\n", f)   // float64

fmt.Printf("%v\n", s)   // Zhang
fmt.Printf("%#v\n", s)  // "Zhang"
fmt.Printf("%T\n", s)   // string

fmt.Printf("%v\n", a)   // [1 2 3]
fmt.Printf("%#v\n", a)  // []int{1, 2, 3}
fmt.Printf("%T\n", a)   // []int

fmt.Printf("%%\n")      // %
```

### 5.2 整数格式

以下输出格式要和整数类型一起使用：

| 格式   | 含义                                    |
| ------ | --------------------------------------- |
| `%b`   | 以二进制格式输出                        |
| `%o`   | 以八进制格式输出                        |
| `%O`   | 以八进制格式输出并且显示前缀 `0o`       |
| `%d`   | 以十进制格式输出                        |
| `%+d`  | 以十进制格式输出并且显示符号            |
| `%x`   | 以十六进制格式小写输出                  |
| `%X`   | 以十六进制格式大写输出                  |
| `%#x`  | 以十六进制格式小写输出并且显示前缀 `0x` |
| `%#X`  | 以十六进制格式大写输出并且显示前缀 `0X` |
| `%4d`  | 以宽度为 `4` 的格式输出，左侧填充空格   |
| `%-4d` | 以宽度为 `4` 的格式输出，右侧填充空格   |
| `%04d` | 以宽度为 `4` 的格式输出，左侧填充 `0`   |

```go
var i = 123

fmt.Printf("%b\n", i)    // 1111011
fmt.Printf("%o\n", i)    // 173
fmt.Printf("%O\n", i)    // 0o173
fmt.Printf("%d\n", i)    // 123
fmt.Printf("%+d\n", i)   // +123
fmt.Printf("%x\n", i)    // 7b
fmt.Printf("%X\n", i)    // 7B
fmt.Printf("%#x\n", i)   // 0x7b
fmt.Printf("%#X\n", i)   // 0X7B
fmt.Printf("%4d\n", i)   //  123
fmt.Printf("%-4d\n", i)  // 123 
fmt.Printf("%04d\n", i)  // 0123
```

### 5.3 字符串格式

以下输出格式要和字符串类型一起使用：

| 格式   | 含义                                         |
| ------ | -------------------------------------------- |
| `%s`   | 纯字符串输出                                 |
| `%q`   | 用双引号包裹输出                             |
| `%8s`  | 以宽度为 `8` 的格式输出，左侧填充空格        |
| `%-8s` | 以宽度为 `8` 的格式输出，右侧填充空格        |
| `%x`   | 以十六进制输出字符串的每个字节               |
| `% x`  | 以十六进制输出字符串的每个字节并用空格分隔   |

```go
var s = "Zhang"

fmt.Printf("%s\n", s)    // Zhang
fmt.Printf("%q\n", s)    // "Zhang"
fmt.Printf("%8s\n", s)   //    Zhang
fmt.Printf("%-8s\n", s)  // Zhang   
fmt.Printf("%x\n", s)    // 5a68616e67
fmt.Printf("% x\n", s)   // 5a 68 61 6e 67
```

### 5.4 布尔格式

以下输出格式要和布尔类型一起使用：

| 格式 | 含义       |
| ---- | ---------- |
| `%t` | 输出布尔值 |

```go
fmt.Printf("%t\n", true)  // true
```

### 5.5 浮点数格式

以下输出格式要和浮点数类型一起使用：

| 格式    | 含义                        |
| ------- | --------------------------- |
| `%f`    | 保留 `6` 位小数             |
| `%.2f`  | 保留 `2` 位小数             |
| `%6.2f` | 宽度 `6`，精度 `2`          |
| `%e`    | 科学计数法，输出的 `e` 小写 |
| `%E`    | 科学计数法，输出的 `E` 大写 |

```go
var f float64 = 0.125

fmt.Printf("%f\n", f)     // 0.125000
fmt.Printf("%.2f\n", f)   // 0.12
fmt.Printf("%6.2f\n", f)  //   0.12
fmt.Printf("%e\n", f)     // 1.250000e-01
fmt.Printf("%E\n", f)     // 1.250000E-01
```

## 6 条件语句

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

## 7 循环语句

`for` 循环是 Go 语言中唯一的循环语句。

```go
for i := 0; i < 10; i++ {
    if i == 2 {
        continue
    }
    if i == 6 {
        break
    }
    fmt.Printf("%d ", i)  // 0 1 3 4 5 
}

for index, value := range "Zhang" {
    fmt.Printf("%d-%c ", index, value)  // 0-Z 1-h 2-a 3-n 4-g 
}

for index, value := range []int{66, 88, 99} {
    fmt.Printf("%d-%d ", index, value)  // 0-66 1-88 2-99 
}
```

遍历字符串时，`index` 是 UTF-8 字节偏移量，`value` 是 `rune`。字符串含有非 ASCII 字符时，索引不一定连续。Go 1.22 起也可以对整数执行 `range`，例如 `for i := range 10` 会依次得到 0 到 9。

## 8 `goto` 语句

```go
 i := 1
LOOP:
 fmt.Printf("%v ", i)  // 1 2 3
 i++
 if i <= 3 {
  goto LOOP
 }
```

## 9 `switch` 语句

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

## 10 数组

```go
arr1 := [3]int{1, 2, 3}
arr2 := [...]int{4, 5}
fmt.Printf("%T %T", arr1, arr2)  // [3]int [2]int
```

::: warning 注意
在 Go 语言中，数组长度固定，要么给数组指定长度，要么使用 `...` 让编译器推断数组的长度。
:::

只初始化数组特定位置：

```go
arr1 := [8]int{3: 33, 5: 55}
arr2 := [...]int{6: 66}
fmt.Println(arr1)  // [0 0 0 33 0 55 0 0]
fmt.Println(arr2)  // [0 0 0 0 0 0 66]
```

可以使用 `==`、`!=` 直接比较两个数组是否相等：

```go
arr1 := [3]int{1, 2, 3}
arr2 := [...]int{1, 2, 3}
fmt.Println(arr1 == arr2)  // true
fmt.Println(arr1 != arr2)  // false
```

只有数组元素类型可比较时，数组才能使用 `==` 和 `!=`。切片不能直接比较，唯一允许的比较是与 `nil` 比较。

## 11 切片

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
fmt.Println(slice1, len(slice1), cap(slice1))  // [1 2 3] 3 3
fmt.Println(slice2, len(slice2), cap(slice2))  // [2 3 4] 3 4
fmt.Println(slice3, len(slice3), cap(slice3))  // [0 0 0] 3 6
```

`cap()` 返回切片的容量，如果切片从数组创建的，那么切片的容量为选中的起始位置到数组末尾。

`make()` 的第二个参数为长度，第三个参数为容量，如果容量未指定，则默认等于长度。

内置函数 `clear` 会把切片现有长度范围内的元素重置为零值，但不会改变长度和容量。

### 11.2 修改

```go
slice1 := []int{1, 2, 3}
slice2 := append(slice1, 4, 5)       // 追加一个元素
slice3 := append(slice2, slice1...)  // 追加一个切片的所有元素
fmt.Println(slice1)  // [1 2 3]
fmt.Println(slice2)  // [1 2 3 4 5]
fmt.Println(slice3)  // [1 2 3 4 5 1 2 3]
```

::: warning 注意
将另一个切片的所有元素追加到一个切片上时，需要在另一个切片后面写上 `...`。
:::

```go
// 删除切片中索引为 2 的元素
slice1 := []int{1, 2, 3, 4, 5}
slice2 := append(slice1[:2], slice1[3:]...)
fmt.Println(slice1)  // [1 2 4 5 5]
fmt.Println(slice2)  // [1 2 4 5]
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
fmt.Println(slice1)  // [1 100 3 4 5 6]
fmt.Println(slice2)  // []，因为容量为 0，所以只 copy 了 0 个元素
fmt.Println(slice3)  // [2 3 4]，因为容量为 3，所以只 copy 了 3 个元素
```

只有目标切片拥有独立的底层数组时，`copy` 才能解除共享关系。上例的 `slice3` 由 `make` 创建，因此修改 `slice1` 不会影响它。

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
切片参数仍然按值传递，但复制后的切片描述符最初指向同一个底层数组。上例第一次修改会影响调用方；随后 `append` 触发扩容，局部切片改为指向新数组，后续修改不再影响原切片。是否扩容取决于容量，代码不能依赖某一次具体的扩容时机。
:::

### 11.3 底层

可以把切片概念化为下面三个字段，但这不是语言规范公开的结构体定义：

```go
type slice struct {
    ptr *T   // 指向底层数组的指针
    len int  // 切片的长度
    cap int  // 切片的容量
}
```

切片本身并不存储数据，而是引用一个底层数组。

## 12 Map

`map` 是可变的键值集合。遍历顺序没有规定，同一个 `map` 的多次遍历也可能得到不同顺序。普通 `map` 不支持无同步的并发读写，需要使用锁、单 goroutine 管理或 `sync.Map` 等方案。

### 12.1 创建

```go
var myMap = map[string]string{
    "key1": "value1",
    "key2": "value2",
}
myMap["key3"] = "value3"
fmt.Println(myMap)  // 输出顺序不固定
```

下面这种用法将会报错：

```go
var myMap map[string]string
fmt.Println(myMap == nil)  // true
myMap["key"] = "value"  // 报错：panic: assignment to entry in nil map
```

通过 `make()` 创建空 Map，可以防止这种错误：

```go
var myMap = make(map[string]string)
fmt.Println(myMap == nil)  // false
myMap["key"] = "value"
fmt.Println(myMap)  // map[key:value]
```

### 12.2 遍历

```go
for key, value := range myMap {}
for key := range myMap {}
```

### 12.3 删除

```go
var myMap = map[string]string{
    "key1": "value1",
    "key2": "value2",
    "key3": "value3",
}
fmt.Println(myMap)  // 输出顺序不固定
delete(myMap, "key2")
fmt.Println(myMap)
delete(myMap, "key4")  // 删除不存在的元素也不会报错
fmt.Println(myMap)
```

需要删除全部键时可以调用 `clear(myMap)`，`map` 本身仍可继续使用。

### 12.4 查询

```go
var myMap = map[string]string{
    "key1": "value1",
    "key2": "value2",
    "key3": "value3",
}
key1, ok1 := myMap["key2"]
key2, ok2 := myMap["key4"]
fmt.Println(ok1, key1)  // true value2
fmt.Println(ok2, key2)  // false (空字符串)
```

## 13 函数

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
 fmt.Println(add(1, 2))  // 3 <nil>
    fmt.Println(add2(1, 2, 3))  // 6 <nil>
    fmt.Println(getFunc()())  // 10
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
 fmt.Println(nextNum())  // 1
 fmt.Println(nextNum())  // 2
 fmt.Println(nextNum())  // 3

 nextNum = autoIncrement()
 fmt.Println(nextNum())  // 1
 fmt.Println(nextNum())  // 2
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
 fmt.Print(deferPrint())  // 3210
}
```

延迟函数的参数在 `defer` 语句出现时就已经确定了：

```go
n := 1
defer fmt.Println(n)  // 1
n++
fmt.Println(n)  // 2
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
 fmt.Println(outer())  // 20
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
 setAge(-12)  // panic: negative age
 fmt.Println(0)  // 执行不到这里
}
```

`recover` 只能在延迟函数中捕获同一 goroutine 正在传播的 panic。它适合在进程边界做隔离，普通业务错误仍应返回 `error`。

```go
func setAge(age int) {
 defer func() {
  if r := recover(); r != nil {
   fmt.Println(r)  // negative age
  }
 }()
 if age < 0 {
  panic("negative age")
 }
}

func main() {
 setAge(-12)
 fmt.Println(0)  // 0
}
```

## 14 `type`

```go
type myInt1 = int  // 类型别名，在编译时会直接替换为 int
type myInt2 int    // 自定义类型

func main() {
 var a myInt1
 var b myInt2
 fmt.Printf("%T\n", a)  // int
 fmt.Printf("%T\n", b)  // main.myInt2
}
```

## 15 结构体

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
fmt.Println(stu1.name, stu1.age)  // Zhang 20
fmt.Println(stu2.name, stu2.age)  // Klose 0
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
 fmt.Println(stu.info.name)  // Zhang
}
```

### 15.4 匿名嵌套

匿名嵌套可直接用 `.` 访问的被嵌入内部的结构体的成员，但是声明变量时，依旧要写完整被嵌入的结构体。

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
 fmt.Println(stu.name)  // Zhang
}
```

如果结构体中与被嵌入内部的结构体中含有相同的成员，会优先使用外层的。

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
 fmt.Println(stu.name)  // Klose
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
 stu.print()  // name: zhang, age: 20
}
```

## 16 指针

```go
func increase(n *int) {
 *n++
}

func main() {
 n := 15
 increase(&n)
 fmt.Println(n)  // 16
}
```

未初始化的指针的值为 `nil`，不可直接使用，通过 `new()` 可以避免该问题：

```go
var n1 *int
var n2 = new(int)
fmt.Printf("type: %T, n1==nil: %t\n", n1, n1 == nil)  // type: *int, n1==nil: true
fmt.Printf("type: %T, n2==nil: %t\n", n2, n2 == nil)  // type: *int, n2==nil: false
```

## 17 接口

Go 中没有关键字显式声明某个类型实现了某个接口。只要一个类型实现了接口要求的所有方法，该类型就自动被认为实现了该接口。

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

### 17.3 空接口

空接口 `interface{}` 不要求任何方法，因此任意类型都满足它。Go 1.18 起通常写成等价的别名 `any`。

### 17.4 类型断言

类型断言用于从接口类型中提取其底层值。如果类型不匹配，会触发 `panic`。

```go
var a any = "zhang"
name := a.(string)
fmt.Println(name)  // zhang
age := a.(int)  // 报错：panic: interface conversion: interface {} is string, not int
fmt.Println(age)
```

为了避免 `panic`，可以使用带检查的类型断言：

```go
var a any = "zhang"
age, isInt := a.(int)
if isInt {
 fmt.Println("age:", age)
} else {
 fmt.Println("a:", a)  // zhang
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
 fmt.Println("string")  // string
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
 printSlice(data...)  // 报错：cannot use data (variable of type []string) as []interface{} value in argument to printSlice
}
```

`[]string` 与 `[]any` 的底层元素表示不同，不能直接展开传递。需要新建 `[]any`，逐项转换后再调用。

### 17.7 `error` 接口

`error` 内置接口类型的源码如下：

```go
type error interface {
 Error() string
}
```

自定义错误：

```go
type newError struct {}
func (e *newError) Error() string {
 return "新错误"
}

func main() {
 err := &newError{}
 fmt.Println(err)  // 新错误
}
```

## 18 `package`

假设有个包在目录 `proj/user` 下：

```go
package person

type Person struct {
 Name string
}
```

### 18.1 导入

```go
import (
 "fmt"
 "proj/user"
)

func main() {
 p := person.Person{Name: "Zhang"}
 fmt.Println(p.Name)  // Zhang
}
```

也可以给包起个别名：

```go{3}
import (
 "fmt"
 u "proj/user"
)

func main() {
 p := u.Person{Name: "Zhang"}
 fmt.Println(p.Name)  // Zhang
}
```

点导入会把目标包的导出标识符直接引入当前文件，容易造成来源不明和命名冲突，普通业务代码不建议使用：

```go{3}
import (
 "fmt"
 . "proj/user"
)

func main() {
 p := Person{Name: "Zhang"}
 fmt.Println(p.Name)  // Zhang
}
```

空白导入只执行目标包的初始化副作用，常用于注册数据库驱动、图片解码器等实现：

```go{3}
import (
 "fmt"
 _ "proj/user"
)

func main() {}
```

### 18.2 `init`

包初始化时会先计算包级变量，再按文件顺序执行 `init`。程序不应依赖同一包内不同文件的文件名顺序；复杂初始化更适合放进显式函数并返回错误。

```go{7-9}
package person

type Person struct {
 Name string
}

// 在当前包被导入时自动执行
func init() {
 fmt.Println("init")
}
```

## 19 `go modules`

### 19.1 初始化模块

初始化一个新的模块，并创建 `go.mod` 文件。

```shell
go mod init github.com/username/project
```

### 19.2 添加依赖

添加或更新依赖包到 `go.mod` 文件中。

```shell
go get github.com/gin-gonic/gin
```

### 19.3 移除未使用的依赖

清理 `go.mod` 和 `go.sum` 文件，移除未使用的依赖，并添加缺失的依赖。

```shell
go mod tidy
```

### 19.4 查看依赖关系

列出当前模块及其所有依赖。

```shell
go list -m all
```

显示模块依赖图。

```shell
go mod graph
```

### 19.5 下载依赖

下载 `go.mod` 中指定的所有模块。

```shell
go mod download
```

### 19.6 校验依赖

```shell
go mod verify
```

该命令检查本地模块缓存中的依赖内容是否与 `go.sum` 记录一致。查看全部模块命令可运行 `go help mod`。

## 20 单元测试

Go 的单元测试主要依赖于 `testing` 包，并且通过 `go test` 命令来执行测试。

### 20.1 基本结构

测试文件必须以 `_test.go` 结尾。测试可以与被测代码使用同一包，也可以使用以 `_test` 结尾的外部测试包；前者能访问未导出标识符，后者更接近真实调用方。

测试函数的名称必须以 `Test` 开头，并接受一个 `*testing.T` 类型的参数。例如：

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

在包含测试文件的目录下，运行以下命令即可执行所有测试：

```shell
go test
```

使用 `-v` 标志可以查看详细的测试输出：

```shell
go test -v
```

可以通过指定测试函数名称来运行特定的测试：

```shell
go test -run TestAdd
```

### 20.4 性能测试

Go 支持基准测试（Benchmark Tests），用于衡量代码的性能。基准测试函数的名称必须以 `Benchmark` 开头，并接受一个 `*testing.B` 类型的参数。

```go
// math/math_test.go
func BenchmarkAdd(b *testing.B) {
 for b.Loop() {
  Add(1, 2)
 }
}
```

使用 `-bench` 标志运行基准测试：

```shell
go test -bench=. -benchmem
```

`B.Loop` 由测试框架决定迭代次数，并把循环控制排除在计时之外。也可以指定特定的基准测试：

```shell
go test -bench='^BenchmarkAdd$' -benchmem
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

使用 `-short` 标志跳过该测试用例：

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

## 21 泛型

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

### 21.2 Map上使用

```go
type MyMap[
 K int | string,
 V float32 | float64,
] map[K]V

func main() {
 m := MyMap[string, float64]{}
}
```

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

## 22 函数选项模式

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
func NewDBClient(options ...Option) *DBOptions {
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
 dbOptions := NewDBClient(WithHost("192.168.0.1"))
 fmt.Println(dbOptions.Host) // 192.168.0.1
 fmt.Println(dbOptions.Port) // 3306
}
```

真实的构造函数通常还要校验端口、地址和凭据，并在配置无效时返回 `error`。

## 23 错误处理（`error`）

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

1. 使用内建 `errors.New`
 ```go
 import "errors"

 err := errors.New("something went wrong")
 ```

2. 使用 `fmt.Errorf`
 ```go
 err := fmt.Errorf("invalid input: %v", input)
 ```

3. 使用错误包装
 ```go
 baseErr := errors.New("disk error")
 err := fmt.Errorf("upload failed: %w", baseErr) // %w 表示包装
 ```

### 23.4 判断和提取底层错误

使用 `errors.Is` 和 `errors.As` 判断错误类型或提取特定错误。

1. `errors.Is` 判断是否是某种错误
 ```go
 if errors.Is(err, os.ErrNotExist) {
  fmt.Println("File does not exist")
 }
 ```

2. `errors.As` 判断并提取特定错误类型
 ```go
 var pathErr *os.PathError
 if errors.As(err, &pathErr) {
  fmt.Println("Path error:", pathErr.Path)
 }
 ```

### 23.5 错误链与调用栈

标准库的 `%w`、`errors.Is` 和 `errors.As` 维护错误链，但不会自动记录每一层调用栈。需要排查崩溃位置时，可以在日志边界使用 `runtime/debug.Stack`；生产服务通常把错误链、请求信息和追踪数据一起交给日志或可观测系统。

<<< @/go/codes/golang/error_stack.go

## 24 AST 代码生成

假设我们要维护状态码的相关代码，可能会这样来写：

```go
const (
 OK            Code = 0 // OK
 InvalidParams Code = 1 // 参数错误
 Timeout       Code = 2 // 超时
)

var mapCodeDesc = map[int]string{
 OK:            "OK",
 InvalidParams: "参数错误",
 Timeout:       "超时",
}

func GetCodeDesc(code int) string {
 if desc, ok := mapCodeDesc[code]; ok {
  return desc
 }
 return "未知错误"
}
```

但是，有没有一种办法可以只维护状态码，而状态码的描述映射和获取函数都自动生成呢。

通过 `stringer` 工具能够做到这点。先安装：

```shell
go install golang.org/x/tools/cmd/stringer@latest
```

此时原本的代码文件就可以写为：

```go
//go:generate stringer -type Code -linecomment

package code

type Code int64

const (
 OK            Code = 0 // OK
 InvalidParams Code = 1 // 参数错误
 Timeout       Code = 2 // 超时
)

```

然后在代码文件同级目录下执行：

```go
go generate .
```

此时，会多出一个文件 `code_string.go`：

<<< @/go/codes/golang/code_string.go

### 24.1 什么是 AST

- AST（Abstract Syntax Tree）是编程语言源代码的树状结构表示。
- 在 Go 中，AST 表示形式由 `go/ast` 和 `go/token` 标准库提供。
- 借助 AST，我们可以**动态构建代码**，然后导出为 `.go` 文件。

### 24.2 常用标准库

| 包名 | 用途 |
|------|------|
| `go/ast` | 抽象语法树节点结构和操作 |
| `go/token` | 源码位置（token位置管理） |
| `go/parser` | 源代码 -> AST |
| `go/printer` | AST -> 源代码输出 |
| `go/format` | 美化输出代码 |

### 24.3 基本步骤

以构造函数生成器为例。

#### 24.3.1 步骤一：导入包

```go
import (
 "bytes"
 "go/ast"
 "go/format"
 "go/token"
 "os"
)
```

#### 24.3.2 步骤二：构造代码结构

我们创建一个函数：

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

#### 24.3.3 步骤三：组合为完整文件

```go
func buildFile() *ast.File {
 return &ast.File{
  Name:  ast.NewIdent("main"),
  Decls: []ast.Decl{buildFunc()},
 }
}
```

#### 24.3.4 步骤四：输出生成的代码到 `.go` 文件

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

运行后，会生成 `hello_gen.go`，内容类似于：

<<< @/go/codes/golang/hello_gen.go

### 24.4 AST 节点常见结构

| AST 类型             | 描述          | 示例            |
| -------------------- | ------------- | --------------- |
| `*ast.File`          | 表示一个 Go 文件    | 整体文件结构          |
| `*ast.FuncDecl`      | 表示函数声明        | `func Foo() {}` |
| `*ast.GenDecl`       | 常量/变量/类型声明    | `var x int`     |
| `*ast.AssignStmt`    | 赋值语句          | `x := 1`        |
| `*ast.ReturnStmt`    | 返回语句          | `return x`      |
| `*ast.BinaryExpr`    | 二元运算符         | `x + y`         |
| `*ast.BasicLit`      | 字面量           | `"hello"`, `1`  |
| `ast.NewIdent(name)` | 标识符（变量名/函数名等） | `x`, `Foo`      |
