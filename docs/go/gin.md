---
title: Gin 路由、中间件与参数绑定
description: 使用 Gin 编写 Go HTTP 服务，介绍路由分组、路径与查询参数、表单绑定、数据验证、中间件及 Protobuf 响应，并通过示例说明请求处理和优雅退出。
outline: deep
---

# Gin

Gin 是基于 Go 标准库 `net/http` 的 Web 框架，提供路由、中间件、参数绑定和响应渲染等功能。本文使用 Go 1.26.x 和 Gin 1.12.0。该 Gin 版本要求 Go 1.25 或更高版本。

## 1. 安装

先创建模块，再添加指定版本的依赖：

<<< @/go/codes/gin/install.sh

将下一节代码保存为模块目录中的 `main.go`，再运行 `go mod tidy` 和 `go run main.go`。`go get` 在这里用于管理模块依赖，不是安装可执行程序。初次安装无需使用 `-u` 同时升级其他依赖。

## 2. Hello World

<<< @/go/codes/gin/hello_world.go

`gin.Default()` 创建带有 `Logger` 和 `Recovery` 中间件的路由引擎，分别用于请求日志和请求处理链中的 panic 恢复。`gin.New()` 不添加这些默认中间件。`Recovery` 不能替其他 goroutine 恢复 panic。

示例只监听本机的 `127.0.0.1:8000`。`Run` 会阻塞并返回监听或服务错误，应检查返回值。另开终端访问：

```shell
curl http://127.0.0.1:8000/hello
```

响应状态为 `200 OK`，正文为：

```json
{"message":"world"}
```

`c.JSON` 会序列化数据并设置 JSON 的 Content-Type。`gin.H` 的底层类型是 `map[string]any`，适合构造简单响应。

## 3. 路由分组

假设需要提供 `/goods/list`、`/goods/add` 和 `/goods/del`，可以分别注册路由：

<<< @/go/codes/gin/router.go

这些路径都有 `/goods` 前缀，可以将 `main` 中创建路由引擎和注册路由的部分替换为：

<<< @/go/codes/gin/group.go

`goodsList`、`addGoods`、`delGoods` 和最后的 `Run` 保持原样。分组后的完整路径和 HTTP 方法不变。分组也可以接收中间件，统一作用于该组随后注册的路由。

示例的处理函数只返回操作名称，没有访问数据库。GET 用于查询，添加和删除这里用 POST 演示，不应仅因路径名称包含 `delete` 就用 GET 修改资源。

## 4. 路径参数

### 4.1 单个路径段

将 Hello World 的 `main` 中创建路由引擎和注册路由的部分替换为：

<<< @/go/codes/gin/param.go

`/goods/:id` 匹配一个非空路径段。访问 `/goods/123` 时，`Param("id")` 返回字符串 `"123"`，不会自动转换成整数。只替换上例的 `GET` 注册语句，多个参数可以写为：

<<< @/go/codes/gin/params.go

`router` 沿用前例，两个示例分别运行。请求 `/goods/123/delete` 时，响应为：

```json
{"action":"delete","id":"123"}
```

这个路由只回显参数，不执行删除操作。`:action` 只匹配一个非空路径段，不能匹配 `/goods/123/delete/test`。

### 4.2 通配参数

`*action` 匹配剩余路径，必须放在路由末尾，取得的值包含开头的 `/`。沿用已有的 `router`，将 `GET` 注册语句替换为：

<<< @/go/codes/gin/params2.go

在默认配置下，响应如下：

| 请求路径                 | 状态码 | `action` 或重定向位置   |
| ------------------------ | ------ | ----------------------- |
| `/goods/123/delete`      | 200    | `"/delete"`             |
| `/goods/123/delete/test` | 200    | `"/delete/test"`        |
| `/goods/123/`            | 200    | `"/"`                   |
| `/goods/123`             | 301    | `Location: /goods/123/` |

最后一种请求由默认开启的 `RedirectTrailingSlash` 重定向，第一次响应不是 JSON。`curl -L` 会跟随重定向，随后得到 `action` 为 `"/"` 的 JSON。关闭这个选项后，该路径返回 404。

`:action` 和 `*action` 是两种替代写法，不要把这两个相冲突的路由同时注册到同一个引擎。

### 4.3 绑定路径参数

`ShouldBindUri` 按 `uri` 标签将路径参数绑定到结构体，并执行 `binding` 验证：

<<< @/go/codes/gin/bind.go

请求 `/goods/123/abc` 时，响应为：

```json
{"id":123,"name":"abc"}
```

`uri:"id"` 与路由中的 `:id` 对应。`ID` 为 `int`，因此非整数参数会绑定失败，示例返回 400。数值字段的 `required` 会拒绝零值，但不会自动要求它为正数。业务需要正数 ID 时，可再增加 `gt=0`。

## 5. 查询参数与表单

### 5.1 查询参数

<<< @/go/codes/gin/get.go

`Query` 和 `DefaultQuery` 读取 URL 的查询字符串，并不限定只能在 GET 请求中使用。

| 请求路径                            | 响应正文                               |
| ----------------------------------- | -------------------------------------- |
| `/hello`                            | `{"framework":"Gin","lang":""}`        |
| `/hello?lang=Java&framework=Spring` | `{"framework":"Spring","lang":"Java"}` |
| `/hello?framework=`                 | `{"framework":"","lang":""}`           |

参数缺失与显式传入空值有所区别。`Query` 在这两种情况下都返回空字符串，`DefaultQuery` 只在参数缺失时使用默认值。需要区分两者时使用 `GetQuery` 返回的布尔值。重复参数可以通过 `QueryArray` 获取。

### 5.2 表单参数

<<< @/go/codes/gin/post.go

发送 URL 编码的表单：

```http
POST http://127.0.0.1:8000/hello
Content-Type: application/x-www-form-urlencoded

lang=Go&framework=Gin
```

响应为：

```json
{"framework":"Gin","lang":"Go"}
```

表单正文不要在字段名和值之间额外加入空格，否则空格可能成为数据的一部分。也可以用下面的命令，由 curl 编码表单字段：

```shell
curl --data-urlencode 'lang=Go' --data-urlencode 'framework=Gin' \
  http://127.0.0.1:8000/hello
```

`PostForm` 和 `DefaultPostForm` 读取 URL 编码或 multipart 表单，不读取 URL 查询参数，也不会解析 JSON 正文。`DefaultPostForm` 同样只在字段缺失时使用默认值，`GetPostForm` 可以区分字段缺失和显式空值。

## 6. Protobuf 渲染

先安装 Protobuf 编译器 `protoc`。将下面的消息定义保存为模块根目录下的 `msg.proto`：

<<< @/go/codes/gin/msg.proto

`package example` 声明 Protobuf 命名空间。`go_package` 中的 `example.com/gin-notes/pb` 是生成代码的 Go 导入路径，`;pb` 指定 Go 包名。这里与第 1 节的模块路径对应，换成其他模块时需要同时调整下面的导入。

将服务代码保存为同一目录的 `main.go`：

<<< @/go/codes/gin/proto.go

安装 Go 代码生成器，将其加入 PATH，再生成 `pb/msg.pb.go`：

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10
export PATH="$(go env GOPATH)/bin:$PATH"
mkdir -p pb
protoc --go_out=pb --go_opt=paths=source_relative msg.proto
go mod tidy
go run main.go
```

`paths=source_relative` 让输出文件按输入文件的相对路径放在 `--go_out` 目录下。生成后的目录包含：

```text
gin-notes/
├── go.mod
├── go.sum
├── main.go
├── msg.proto
└── pb/
    └── msg.pb.go
```

`c.ProtoBuf` 返回 Protobuf 二进制数据，Content-Type 为 `application/x-protobuf`，不会自动转换成 JSON。另开终端，在模块根目录获取并解码响应：

```shell
curl --output response.bin http://127.0.0.1:8000/hello
protoc --decode=example.Teacher msg.proto < response.bin
```

解码结果为：

```text
name: "zhang"
courses: "Gin"
courses: "GoLang"
```

## 7. 参数绑定与验证

### 7.1 JSON 绑定

Gin 使用 `go-playground/validator/v10` 执行 `binding` 标签中的验证规则。Gin 1.12.0 默认依赖 validator 10.30.1，以下错误示例按这个版本展示。

<<< @/go/codes/gin/validate.go

`ShouldBindJSON` 明确按 JSON 解析请求正文，绑定或验证失败时返回错误，不会自动设置错误响应。示例选择返回 400，并结束当前处理函数。`BindJSON` 则在失败时自动中止处理链并写入 400，不适合在这之后再改成其他响应状态。

| 标签               | 含义                                     |
| ------------------ | ---------------------------------------- |
| `required`         | 字段不能为对应类型的零值                 |
| `min=3,max=20`     | 字符串长度在 3 到 20 之间                |
| `min=8,max=20`     | 字符串长度在 8 到 20 之间                |
| `eqfield=Password` | 与同一结构体中的 Go 字段 `Password` 相等 |
| `email`            | 符合该验证器支持的邮箱格式               |
| `lte=120`          | 数值不大于 120                           |

字符串长度按 Unicode 码点数量检查，不是 UTF-8 字节数。`Age` 使用无符号整数，负数或不匹配的 JSON 类型会在解析阶段失败。这里没有给年龄设置 `required`，省略它时保留零值 0，仍能通过验证。如果需要区分未传字段与显式传入零值，可以使用指针字段。

发送以下请求，用户名长度、密码确认、邮箱格式和年龄会分别验证失败：

```http
POST http://127.0.0.1:8000/signUp
Content-Type: application/json

{
  "username": "ab",
  "password": "password123",
  "rePassword": "other123",
  "email": "invalid",
  "age": 130
}
```

响应状态为 400，JSON 正文中的 `error` 是字符串，包含以下各行。序列化时换行会编码为 `\n`：

```text
Key: 'SignUpInfo.Username' Error:Field validation for 'Username' failed on the 'min' tag
Key: 'SignUpInfo.RePassword' Error:Field validation for 'RePassword' failed on the 'eqfield' tag
Key: 'SignUpInfo.Email' Error:Field validation for 'Email' failed on the 'email' tag
Key: 'SignUpInfo.Age' Error:Field validation for 'Age' failed on the 'lte' tag
```

示例只检查这些字段约束，没有执行完整的用户注册流程。

### 7.2 中文错误信息

可以在启动服务前注册中文翻译，并使用 JSON 字段名作为响应中的错误键：

<<< @/go/codes/gin/translation.go

同一个无效请求会返回 400，正文为：

<<< @/go/codes/gin/error.json

`eqfield` 的参数仍使用 Go 字段名，因此默认翻译中的 `Password` 对应结构体里的同名字段。示例是平铺结构体，错误键使用 `Field()`。有嵌套结构或切片时，需要设计完整的字段路径，避免同名字段相互覆盖。

只有 `validator.ValidationErrors` 才能按字段翻译。JSON 语法错误、空正文和字段类型不匹配属于绑定错误，示例为它们返回独立的提示。初始化失败会停止启动，避免继续使用未初始化的翻译器。字段名称规则和翻译应在处理并发请求前注册完成。

## 8. 中间件

### 8.1 执行顺序

中间件与最终处理函数使用同一种 `gin.HandlerFunc` 类型，按注册顺序组成请求处理链：

<<< @/go/codes/gin/midware.go

`c.Next()` 在当前中间件内部执行尚未执行的处理函数，返回后再执行计时和日志逻辑。多个中间件都使用这个结构时，前置逻辑按注册顺序执行，后置逻辑按相反顺序返回。

上例的 `Logger` 和 `Recovery` 用于全局路由，`MyLogger` 仅用于 `/hello`。也可以在 `Group` 中配置中间件。`Use` 应放在相应路由注册之前，后续添加的中间件不会改变已注册路由的处理链。子组创建时也会复制已有的上级处理链。

### 8.2 return 与 Abort

直接 `return` 只结束当前中间件函数，框架仍会继续调用后面的处理函数：

<<< @/go/codes/gin/abort.go

访问 `/continue` 时，最终处理函数仍执行并返回 200。访问 `/abort` 时，`AbortWithStatusJSON` 中止后续处理链并返回 403，最终处理函数不会执行。

`Abort` 不会终止当前函数，也不会撤销已经执行的逻辑。调用后是否继续执行当前函数，由普通 Go 控制流决定。仅调用 `Abort()` 不会自动写入错误状态或正文，所以通常搭配 `AbortWithStatus` 或 `AbortWithStatusJSON`，并按需要 `return`。外层已经进入 `Next` 的中间件仍可能继续执行自己的后置逻辑。

### 8.3 处理链源码

`RouterGroup.Use` 追加中间件：

<<< @/go/codes/gin/use.go

注册路由时，`combineHandlers` 将组中间件和路由处理函数合成新的切片：

<<< @/go/codes/gin/combineHandlers.go

引擎匹配路由后调用 `Next`，它按索引执行处理链：

<<< @/go/codes/gin/next.go

`Abort` 将索引设置为中止值，后续调用不再进入循环：

<<< @/go/codes/gin/abortIndex.go

处理链源码见该版本的 [context.go](https://github.com/gin-gonic/gin/blob/v1.12.0/context.go)。

### 8.4 goroutine 与请求上下文

`*gin.Context` 会被复用，不应直接交给请求结束后仍在运行的 goroutine。优先在启动 goroutine 前取出需要的数据，确需传递 Gin 上下文时使用 `c.Copy()` 的副本，只用于读取，不能通过副本写入响应。

`Copy` 不会让共享数据自动获得并发保护，也不会延长 `c.Request.Context()` 的生命周期。需要取消信号和请求截止时间的下游调用，应使用请求的标准库 Context。

## 9. 优雅退出

优雅退出是先停止接受新连接，再为正在处理的 HTTP 请求留出完成时间。Gin 路由引擎实现了 `http.Handler`，可以交给显式创建的 `http.Server`，通过 `Shutdown` 控制关闭过程：

<<< @/go/codes/gin/quit.go

示例使用 `signal.NotifyContext` 接收 `SIGINT` 和 `SIGTERM`，适用于这里的 Linux/Unix 信号场景。`SIGKILL` 不能捕获，也不能用于触发优雅退出。

可以先启动耗时请求，再在服务器终端按 Ctrl+C。服务器停止接受新连接，但该请求在 2 秒后仍可得到响应，随后程序退出。这里使用独立的 `context.Background()` 创建 10 秒关闭期限，不能直接把已经取消的信号上下文传给 `Shutdown`。

需要区分几个行为：

- `ListenAndServe` 在关闭时返回 `http.ErrServerClosed`，这属于正常退出。主 goroutine 仍应等待 `Shutdown` 返回，不能只因监听结束就退出程序。
- `Shutdown` 的期限限制等待时间，不会在超时时自动强制关闭仍在处理的连接。示例在失败时调用 `Close`，并返回错误。
- `Shutdown` 不会自动等待自行创建的后台任务，也不会关闭或等待已经被接管的连接，例如 WebSocket。应用需要另外通知这些任务或连接结束，并等待完成。
- 正常请求不会仅因调用 `Shutdown` 就被强制中断。请求提前取消时，上例通过 `c.Request.Context().Done()` 停止模拟等待。

这些关闭行为由标准库的 [http.Server.Shutdown](https://github.com/golang/go/blob/go1.26.8/src/net/http/server.go) 定义。使用 `http.Server` 也便于按实际服务需要配置请求读取和响应写入超时。
