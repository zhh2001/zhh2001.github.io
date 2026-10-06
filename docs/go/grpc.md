---
outline: deep
---

# gRPC

gRPC 是一套基于服务定义生成客户端与服务端代码的 RPC 框架。Go 项目通常用 Protocol Buffers 描述消息和服务，调用数据通过 HTTP/2 连接传输。本文先用 `net/rpc` 说明 RPC 的基本过程，再进入 Protobuf、流式调用、元数据、拦截器和错误处理。

示例以 Go 1.26.x、gRPC-Go v1.84.0 和 `google.golang.org/protobuf` v1.36.11 为基准，使用 proto3。生成插件分别固定为 `protoc-gen-go` v1.36.11 和 `protoc-gen-go-grpc` v1.6.2。插件、运行库和 `protoc` 的版本号彼此独立。

## 1. RPC

RPC（Remote Procedure Call，远程过程调用）是一种通信模型，不是某一个固定协议。调用方使用本地接口发起请求，框架负责消息编码、网络传输和结果返回。网络故障、超时和重复执行等问题并不会因此消失，业务代码仍要处理这些边界。

### 1.1 基本原理

- **本地调用**：程序调用一个本地函数，直接执行并返回结果。
- **远程调用**：客户端把方法和参数编码后发送到服务端，服务端执行对应方法，再返回结果或错误。

### 1.2 工作流程

1. **客户端调用**：客户端像调用本地方法一样调用远程方法。
2. **编码**：客户端把方法信息和参数编码为约定的传输格式。
3. **网络传输**：序列化后的数据通过网络发送到服务器。
4. **解码**：服务端把请求恢复为程序可用的数据结构。
5. **执行方法**：服务器根据方法名和参数执行对应的函数。
6. **返回结果**：服务端编码响应，并把响应状态和数据发回客户端。
7. **客户端接收**：客户端解码结果，或者处理超时、取消和远端错误。

### 1.3 `net/rpc` 服务端

Go 标准库的 `net/rpc` 可以通过 TCP 或 HTTP 传输 RPC 消息，默认使用 Gob 编码，因此通常用于 Go 程序之间的通信。它支持自定义编解码器，标准库的 `net/rpc/jsonrpc` 就提供了 JSON-RPC 1.0 实现，可以与其他语言互通。`net/rpc` 已进入冻结状态，不再接受新功能。

可注册的远程方法必须导出，接收两个参数并返回一个 `error`。两个参数的类型必须是导出类型或内建类型，第二个参数必须是指针，用于接收结果。下面用 `RegisterName` 将服务注册为 `HelloService`。

要开发一个 RPC 服务端，通常需要以下几个步骤：

1. **定义服务类型**：结构体的导出方法作为远程方法。
2. **注册服务**：把服务实例注册到 RPC 服务中。
3. **启动服务**：监听端口并处理连接。

<<< @/go/codes/grpc/rpc_server.go

### 1.4 `net/rpc` 客户端

1. **连接到 RPC 服务**：通过 TCP 或 HTTP 连接服务端。
2. **调用远程方法**：通过 RPC 客户端调用服务端提供的方法。

<<< @/go/codes/grpc/rpc_cli.go

分别将服务端和客户端保存为两个独立程序，例如 `rpc_server/main.go` 与 `rpc_client/main.go`。先运行服务端，再运行客户端，输出 `Hello Zhang`。

## 2. Protocol Buffers

Protocol Buffers（Protobuf）同时包含接口描述语言、二进制编码格式和代码生成工具。它可以脱离 gRPC 单独使用。gRPC 默认用 Protobuf 定义服务和消息，也支持其他编解码器。

### 2.1 安装

需要安装 `protoc` 和两个 Go 生成插件。本文用 `protoc` 3.21.12 验证以下 proto3 示例，它支持文中使用的 `optional` 字段。

安装生成插件，并把默认安装目录加入 `PATH`：

<<< @/go/codes/grpc/install.sh

如果设置了 `GOBIN`，应把该目录加入 `PATH`。用 `protoc-gen-go --version` 和 `protoc-gen-go-grpc --version` 检查插件是否可用。

建立独立的示例模块，后续命令都在 `grpcdemo` 目录执行：

```sh
mkdir -p grpcdemo/proto
cd grpcdemo
go mod init example.com/grpcdemo
go get google.golang.org/grpc@v1.84.0 google.golang.org/protobuf@v1.36.11
```

### 2.2 定义数据结构

创建 `proto/msg.proto`：

<<< @/go/codes/grpc/msg.proto

字段后的数字是线上的字段编号，不是数组下标。同一消息内的编号必须唯一，范围为 1 到 2<sup>29</sup>−1，其中 19000 到 19999 为保留区间。编号一旦发布就不应随意修改或分配给其他字段。`go_package` 填写生成代码的 Go 导入路径，这里的 `proto` 目录对应 `example.com/grpcdemo/proto`。

### 2.3 生成 Go 代码

生成消息类型：

<<< @/go/codes/grpc/protoc.sh

输出为 `proto/msg.pb.go`。`--go_out` 生成消息代码，后面定义服务时还需要 `--go-grpc_out` 生成 gRPC 接口。`paths=source_relative` 按输入文件相对于导入根目录的路径输出文件，因此不会再生成一层 `example.com/grpcdemo` 目录。

如果后面的 `proto/hello.proto` 需要导入其他文件，例如：

<<< @/go/codes/grpc/import.proto

假设第三方文件放在 `../third_party/google/api/`，应将 `../third_party` 作为额外的导入根目录：

<<< @/go/codes/grpc/proto_path.sh

`annotations.proto` 还依赖 `http.proto`，两个文件都需要准备。编译时找到 `.proto` 文件，并不代表相应的 Go 包已经安装，还需要引入生成代码使用的 `google.golang.org/genproto/googleapis/api/annotations` 包。仅导入 HTTP 注解不会自动提供 HTTP/JSON 转码，后面的基础示例不使用该导入。

### 2.4 标量类型

下表列出常用标量类型及其生成代码中的对应类型：

| Proto Type |  Go Type  |  C++ Type  | Python Type |
| :--------: | :-------: | :--------: | :---------: |
|  `double`  | `float64` |  `double`  |   `float`   |
|  `float`   | `float32` |  `float`   |   `float`   |
|  `int32`   |  `int32`  | `int32_t`  |    `int`    |
|  `int64`   |  `int64`  | `int64_t`  |    `int`    |
|  `uint32`  | `uint32`  | `uint32_t` |    `int`    |
|  `uint64`  | `uint64`  | `uint64_t` |    `int`    |
|  `sint32`  |  `int32`  | `int32_t`  |    `int`    |
|  `sint64`  |  `int64`  | `int64_t`  |    `int`    |
| `fixed32`  | `uint32`  | `uint32_t` |    `int`    |
| `fixed64`  | `uint64`  | `uint64_t` |    `int`    |
| `sfixed32` |  `int32`  | `int32_t`  |    `int`    |
| `sfixed64` |  `int64`  | `int64_t`  |    `int`    |
|   `bool`   |  `bool`   |   `bool`   |   `bool`    |
|  `string`  | `string`  |  `string`  |    `str`    |
|  `bytes`   | `[]byte`  |  `string`  |   `bytes`   |

### 2.5 默认值

|   数据类型 | 默认值  | 说明                                                         |
| ---------: | :-----: | :----------------------------------------------------------- |
|   `double` |   `0`   | 双精度浮点型                                                 |
|    `float` |   `0`   | 浮点型                                                       |
|    `int32` |   `0`   | 使用 varint，负数的值编码占 10 字节，含负数时可考虑 `sint32` |
|    `int64` |   `0`   | 使用 varint，负数的值编码占 10 字节，含负数时可考虑 `sint64` |
|   `uint32` |   `0`   | 使用变长编码                                                 |
|   `uint64` |   `0`   | 使用变长编码                                                 |
|   `sint32` |   `0`   | 使用 ZigZag 编码，负数通常比 `int32` 更紧凑                  |
|   `sint64` |   `0`   | 使用 ZigZag 编码，负数通常比 `int64` 更紧凑                  |
|  `fixed32` |   `0`   | 值编码固定 4 字节，值 ≥ 2<sup>28</sup> 时比 `uint32` 更短    |
|  `fixed64` |   `0`   | 值编码固定 8 字节，值 ≥ 2<sup>56</sup> 时比 `uint64` 更短    |
| `sfixed32` |   `0`   | 始终为 4 字节                                                |
| `sfixed64` |   `0`   | 始终为 8 字节                                                |
|     `bool` | `false` | 布尔型                                                       |
|   `string` |  `""`   | 必须是 UTF-8 编码的文本                                      |
|    `bytes` |  `nil`  | Go 中为 `[]byte`，读取语义为空字节序列                       |

上表的默认值按 Go 表示，字节数只计算字段值，不含字段标签。选择固定长度或变长编码时，应看数值分布，而非只看最大值。

Go 中消息字段默认是 `nil`，枚举字段默认取编号为 0 的值，`repeated` 和 `map` 字段的零值为空集合，底层切片或 map 可以是 `nil`。Proto3 不允许用 `[default = ...]` 声明自定义默认值，业务默认值应在应用代码中处理。需要区分“字段未设置”和“字段显式设置为零值”时，可以使用 `optional`：

```proto
message PageRequest {
  optional int32 page_size = 1;
}
```

普通 proto3 标量字段不记录是否设置，未设置和显式零值具有相同的读取语义。`optional int32` 在本文生成方式下对应 `*int32`，可以检查指针是否为 `nil`，但 getter 仍会为未设置的字段返回零值。

### 2.6 消息嵌套

可以使用其他消息类型作为字段类型。例如，要在每条 `SearchResponse` 消息中包含 `Result` 消息，可以直接在同一个 `.proto` 文件中定义一个 `Result` 消息类型，然后在 `SearchResponse` 中指定一个字段类型为 `Result`：

<<< @/go/codes/grpc/res.proto

也可以在 `SearchResponse` 内部定义和使用 `Result` 消息类型：

<<< @/go/codes/grpc/res2.proto

如果要在父消息类型之外重用此消息类型，需要使用 `Parent.Type` 引用。下面的片段接在上一段嵌套定义后：

<<< @/go/codes/grpc/parent.proto

也可以根据需要嵌套多层消息，下面的例子中，两个名为 `Inner` 的嵌套类型是完全独立的，因为它们定义在不同的消息中：

<<< @/go/codes/grpc/inner.proto

### 2.7 导入定义

可以通过导入其他 `.proto` 文件来使用里面的定义：

```proto
import "myproject/other_protos.proto";
```

通常只能引用直接导入文件的定义。`import public` 可以转导出定义，Go 代码生成支持这种用法，但并非所有语言的生成器都支持。以下片段省略各文件的语法和包声明，用于说明导入关系：

```proto
// new.proto
message NewMessage {
  string value = 1;
}
```

```proto
// old.proto
import public "new.proto";
import "other.proto";

message OldMessage {
  NewMessage item = 1;
}
```

```proto
// client.proto
import "old.proto";

message ClientMessage {
  NewMessage item = 1;
  OldMessage old = 2;
}
// 不能直接引用 other.proto 中的类型，使用前需要显式导入该文件。
```

### 2.8 枚举类型

<<< @/go/codes/grpc/enum.proto

字段 `SearchRequest.corpus` 的默认值是 `CORPUS_UNSPECIFIED`，因为它的编号为 0。Proto3 枚举定义的第一个值必须为 0，通常用 `*_UNSPECIFIED` 表示调用方没有明确选择。

### 2.9 Map 类型

<<< @/go/codes/grpc/map.proto

`key` 可以使用整数、`bool` 或字符串类型，不能使用浮点数、`bytes` 和枚举。`value` 可以使用除 `map` 以外的字段类型。map 不保证序列化顺序或遍历顺序，解析遇到重复键时保留最后一个值。

### 2.10 时间戳类型

<<< @/go/codes/grpc/ts.proto

`Timestamp` 的主要字段如下。这里仅用于说明结构，项目应导入编译器提供的标准定义，不要另建同名类型：

<<< @/go/codes/grpc/timestamp.proto

`seconds` 表示从 Unix epoch 起的秒数，`nanos` 表示非负的纳秒部分，范围为 0 到 999999999。有效时间范围是公元 0001 年至 9999 年，采用 UTC 表示。

Go 中使用 `google.golang.org/protobuf/types/known/timestamppb`。`Now` 和 `New` 分别从当前时间和 `time.Time` 构造消息。来自外部的时间戳应先用 `CheckValid` 校验，再转换为 `time.Time`，因为 `AsTime` 会对越界的纳秒值做规范化，不会返回校验错误：

```go
ts := timestamppb.Now()
if err := ts.CheckValid(); err != nil {
    return err
}
t := ts.AsTime()
fmt.Println(t)
```

### 2.11 字段演进

删除字段后应保留原编号和名称，防止后续代码误用旧数据：

```proto
message User {
  reserved 2, 4 to 6;
  reserved "nickname";

  string name = 1;
  string email = 3;
}
```

不要修改已经发布字段的编号，也不要把同一编号改成含义不同的字段。Protobuf 二进制兼容性还涉及字段类型及其线上的编码，不能仅凭编号相同就判断兼容。字段改名可能保持二进制兼容，但会改变生成的 API 和默认 JSON 字段名。以上规则基于 [proto3 语言规范](https://protobuf.dev/programming-guides/proto3/)。

## 3. gRPC

gRPC 根据 `.proto` 中的 `service` 生成强类型客户端和服务端接口。一次调用通常包含方法路径、消息、元数据、状态码和可选的流数据。Protobuf 负责消息编码，gRPC 负责调用语义与传输。

### 3.1 定义服务

<<< @/go/codes/grpc/hello.proto

将上述定义保存为 `proto/hello.proto`。

### 3.2 生成代码

在模块根目录运行：

```sh
protoc --proto_path=. \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/hello.proto
```

生成 `proto/hello.pb.go` 和 `proto/hello_grpc.pb.go`，前者包含消息，后者包含服务接口。项目目录如下：

```text
grpcdemo/
├── go.mod
├── go.sum
├── proto/
│   ├── hello.proto
│   ├── hello.pb.go
│   └── hello_grpc.pb.go
├── server/
│   └── main.go
└── client/
    └── main.go
```

代码用 `pb` 作为 `example.com/grpcdemo/proto` 的导入别名。

### 3.3 实现服务端

<<< @/go/codes/grpc/grpc_server.go

保存为 `server/main.go`。服务结构体按值嵌入 `UnimplementedGreeterServer`，新增 RPC 后，尚未实现的方法会返回 `Unimplemented`。它用于兼容服务接口的扩展，不会自动实现业务逻辑。

### 3.4 实现客户端

<<< @/go/codes/grpc/grpc_cli.go

保存为 `client/main.go`。在两个终端分别执行 `go run ./server` 和 `go run ./client`，客户端输出 `Hello Zhang`。

`grpc.NewClient` 自 v1.63.0 起提供，创建时不会进行网络 I/O，因此返回成功不代表服务器可达。通常在第一次 RPC 时解析地址并建立连接，调用的 deadline 也涵盖等待连接的时间。`ClientConn` 可以复用并由多个 goroutine 发起 RPC，无需每次调用都重新创建。

示例只监听本机地址，`insecure.NewCredentials()` 不启用 TLS。跨主机部署时应配置相应的传输凭据。上述 API 的行为以 [gRPC-Go v1.84.0 文档](https://pkg.go.dev/google.golang.org/grpc@v1.84.0) 为准。

## 4. 数据流模式

gRPC 有四种调用模式。同一流的每个发送方向保持消息顺序，不同 RPC 之间没有统一的到达顺序。本文固定的生成插件默认使用泛型流接口。旧版生成代码可能使用服务专属的流接口名称。

第 4.1 至 4.4 节分别展示服务定义和处理方法，不是完整的入口程序。生成方式与第 3 节相同，Go 片段按使用情况导入 `context`、`fmt`、`io`、`grpc` 和 `pb`。第 4.5 节给出可以直接运行的三种流式调用。

### 4.1 单一请求-响应

客户端发送一个请求，服务端返回一个响应。

proto 定义：

<<< @/go/codes/grpc/service.proto

Go 实现：

<<< @/go/codes/grpc/service.go

### 4.2 服务端流式响应

客户端发送一个请求，服务端按顺序返回多个响应。

proto 定义：

<<< @/go/codes/grpc/stream_server.proto

Go 实现：

<<< @/go/codes/grpc/stream_server.go

### 4.3 客户端流式请求

客户端连续发送多个请求，发送结束后由服务端返回一个汇总响应。客户端应调用 `CloseAndRecv`，服务端则在读到 `io.EOF` 后调用 `SendAndClose`。

proto 定义：

<<< @/go/codes/grpc/stream_cli.proto

Go 实现：

<<< @/go/codes/grpc/stream_cli.go

### 4.4 双向流式

客户端和服务端都可以连续发送消息。gRPC-Go 允许一个 goroutine 负责 `Send`、另一个 goroutine 负责 `Recv`，但不能让多个 goroutine 同时调用同一方向的方法。`CloseSend` 也不能与 `Send` 并发执行。发送后不要再修改该消息，因为库仍可能读取它。

`Send` 可能因流量控制而阻塞，成功返回只表示消息已交给传输层，不代表对端已经处理或持久化。双向流应及时接收消息，两端都持续发送而不接收可能造成阻塞。

proto 定义：

<<< @/go/codes/grpc/stream.proto

Go 实现：

<<< @/go/codes/grpc/stream.go

### 4.5 示例

保存为 `proto/stream_example.proto`：

<<< @/go/codes/grpc/stream_example.proto

在模块根目录生成代码：

```sh
protoc --proto_path=. \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/stream_example.proto
```

服务端保存为 `stream_server/main.go`：

<<< @/go/codes/grpc/stream_server_example.go

客户端保存为 `stream_client/main.go`：

<<< @/go/codes/grpc/stream_client_example.go

先停止第 3 节的服务端，再分别运行 `go run ./stream_server` 和 `go run ./stream_client`。客户端先接收三条日志，再得到 `received 3 messages`，最后接收三条聊天回显。

服务端 `Recv` 返回 `io.EOF` 表示客户端已结束发送，此时服务端仍可以返回响应。客户端 `Recv` 返回 `io.EOF` 表示 RPC 以 `OK` 结束，其他错误表示调用失败。客户端 `Send` 返回 `io.EOF` 则不能视为成功，应继续用 `Recv` 或 `CloseAndRecv` 读取最终状态。

`CloseSend` 只关闭客户端的发送方向，不取消调用，也不等待最终状态。应继续接收到流结束，或取消该调用的上下文，释放相关资源。服务端处理函数应返回错误，避免因请求内的 `panic` 终止整个进程。

## 5. 元数据

元数据用于传递认证信息、追踪 ID 等调用上下文，在线上对应 HTTP/2 的 header 和 trailer。键不区分大小写，由字母、数字、`-`、`_` 和 `.` 组成，不能使用保留前缀 `grpc-`。普通值应使用可打印 ASCII，即 0x20 到 0x7E，二进制键以 `-bin` 结尾。

### 5.1 构建元数据

可以使用包 `google.golang.org/grpc/metadata` 创建元数据。类型 `MD` 实际上是 `map`：

<<< @/go/codes/grpc/md.go

元数据可以像 `map` 一样读取，其中每个键对应一组值。它不适合承载大块业务数据，过大的 header 可能被客户端、代理或服务端拒绝。`MD` 是普通 map，不支持并发修改。将它关联到上下文后不要继续修改，需要调整时先调用 `Copy` 创建副本。

本节的 Go 代码是基于第 3、4 节服务的片段，使用 `metadata` 包，并按需补充其他导入。调用辅助函数时传入已有客户端和带 deadline 的上下文。

#### 5.1.1 创建新元数据

可以使用函数 `New` 从 `map[string]string` 创建元数据：

<<< @/go/codes/grpc/mdNew.go

另一种方法是使用 `Pairs`。具有相同键的值将被合并到一个切片中：

<<< @/go/codes/grpc/mdPairs.go

`New` 和 `Pairs` 都会将键转换为小写。`New` 的输入是 map，若大小写不同的多个键合并到一起，值的顺序不应依赖 map 的遍历顺序。

#### 5.1.2 在元数据中存储二进制数据

在 Go API 中，元数据的键和值都是字符串。传递原始字节时，把字节转换为字符串，并给键添加 `-bin` 后缀。gRPC 会在传输时处理 Base64 编解码，无需自行编码：

<<< @/go/codes/grpc/mdPairsBin.go

### 5.2 客户端发送和接收

#### 5.2.1 发送元数据

有两种方法可以将元数据发送到服务端。推荐的方法是使用 `AppendToOutgoingContext` 将键值对追加到上下文中。当不存在元数据时，则添加元数据；当上下文中已存在元数据时，将合并键值对。

<<< @/go/codes/grpc/mdSend.go

也可以使用 `NewOutgoingContext`，但是这会替换上下文中的现有元数据。

<<< @/go/codes/grpc/mdSendNew.go

#### 5.2.2 接收元数据

一元调用：

<<< @/go/codes/grpc/mdSomeRPC.go

流式调用中，`Header` 会等待响应头或调用结束。`Trailer` 应在 `Recv` 返回 `io.EOF` 或其他终止错误后获取：

<<< @/go/codes/grpc/mdSomeStreamingRPC.go

### 5.3 服务端发送和接收

#### 5.3.1 接收元数据

如果是一元调用，则可以使用 RPC 处理程序的上下文。对于流式调用，服务端需要从流中获取上下文。

下面两个实现中的 `FromIncomingContext` 分别演示这两种读取方式。

#### 5.3.2 发送元数据

用下面的方法替换第 3 节的 `SayHello`，演示接收请求元数据和设置响应元数据。一元调用使用 `grpc.SetHeader` 和 `grpc.SetTrailer`，两者都接收当前 RPC 的上下文并返回 `error`：

<<< @/go/codes/grpc/someRPC.go

用下面的方法替换第 4.5 节的 `GetStream`。流对象的 `SetHeader` 返回 `error`，但 `SetTrailer` 没有返回值：

<<< @/go/codes/grpc/someStreamingRPC.go

`SetHeader` 先保存响应头，通常在首条响应消息发送时发出。需要提前发送时使用 `SendHeader`，响应头发送后不能再追加。Trailer 随 RPC 的最终状态发送。`metadata.NewOutgoingContext` 用于发起另一个 RPC，不能用来设置当前调用的响应元数据。

## 6. 拦截器

拦截器位于应用代码与 gRPC 调用之间，常用于日志、认证、指标和追踪。客户端与服务端分别提供一元和流式拦截器，共四种类型。一元拦截器不处理流式 RPC，流式拦截器也不处理一元 RPC。以下类型声明摘自 `grpc` 包，实际项目使用 `grpc.` 前缀引用这些类型。拦截器必须继续调用传入的 `invoker`、`streamer` 或 `handler`，除非它明确要提前拒绝这次调用。

### 6.1 客户端拦截器

#### 6.1.1 一元拦截器

客户端一元拦截器的类型为 `UnaryClientInterceptor`。它本质上是一个带有签名的函数类型：

<<< @/go/codes/grpc/unaryClientInterceptor.go

一元拦截器通常分为三步：预处理、调用 RPC 方法和后处理。

对于预处理，可以通过检查传入的参数来获取有关当前 RPC 调用的信息。参数包括 RPC 上下文、方法字符串、要发送的请求和配置的 `CallOptions`。有了这些信息甚至可以修改 RPC 调用。

预处理后，用户可以通过调用 `invoker` 来调用 RPC 调用。

一旦调用程序返回，就可以对 RPC 调用进行后处理。这通常涉及处理返回的回复和错误。

要在 `ClientConn` 上安装一元拦截器，可以给 `NewClient` 传入 `WithUnaryInterceptor`。需要组合多个拦截器时使用 `WithChainUnaryInterceptor`，注册顺序决定嵌套顺序。例如注册 A、B 后，执行顺序为 A 前处理、B 前处理、RPC、B 后处理、A 后处理。

#### 6.1.2 流拦截器

客户端流拦截器的类型是 `StreamClientInterceptor`。它是一个带有签名的函数类型：

<<< @/go/codes/grpc/streamClientInterceptor.go

流拦截器的实现通常包括预处理和流操作拦截。

预处理类似于一元拦截器。

客户端流拦截器调用 `streamer` 获取 `ClientStream`。此时返回并不表示整个 RPC 已结束。要记录消息收发或最终接收状态，需要包装 `ClientStream`，覆盖 `SendMsg`、`RecvMsg` 等方法，再返回包装后的流。

流拦截器通过 `WithStreamInterceptor` 安装，多个流拦截器可以使用 `WithChainStreamInterceptor` 组合。

#### 6.1.3 示例

下面是客户端一元拦截器的完整程序，可作为第 3 节客户端的替代实现。调用完成后记录方法、耗时和状态码，并将原错误返回。

<<< @/go/codes/grpc/clientInterceptor.go

### 6.2 服务端拦截器

服务端拦截器与客户端拦截器相似，但提供的参数信息略有不同。

#### 6.2.1 一元拦截器

服务端一元拦截器的类型是 `UnaryServerInterceptor`。它是一个带有签名的函数类型：

<<< @/go/codes/grpc/unaryServerInterceptor.go

服务端一元拦截器通过 `UnaryInterceptor` 安装，多个拦截器使用 `ChainUnaryInterceptor`。

#### 6.2.2 流拦截器

服务端流拦截器的类型是 `StreamServerInterceptor`。它是一个具有签名的函数类型：

<<< @/go/codes/grpc/streamServerInterceptor.go

服务端流拦截器通过 `StreamInterceptor` 安装，多个拦截器使用 `ChainStreamInterceptor`。服务端流拦截器的 `handler` 会持续执行到流处理方法返回。多个服务端拦截器也按外层到内层调用，返回时顺序相反。

#### 6.2.3 示例

下面是第 3 节服务端的替代实现。RPC 处理函数可能并发执行，拦截器中共享的可变状态也需要保护。

<<< @/go/codes/grpc/serverInterceptor.go

## 7. 错误处理

gRPC 错误由状态码、说明文本和可选详情组成。HTTP 状态为 200 不代表 RPC 成功，还要读取 gRPC 的最终状态，`OK` 才表示成功。服务端应返回稳定、可判断的状态码，不要把数据库错误或内部堆栈直接暴露给客户端。客户端先判断状态码，再决定提示用户、重试还是终止调用。

### 7.1 状态码

|        状态码         |  ID   | 描述                                                 |
| :-------------------: | :---: | ---------------------------------------------------- |
|         `OK`          |  `0`  | 成功时返回                                           |
|      `CANCELLED`      |  `1`  | 操作被取消，通常是由调用者取消的                     |
|       `UNKNOWN`       |  `2`  | 未知错误                                             |
|  `INVALID_ARGUMENT`   |  `3`  | 客户端指定的参数无效                                 |
|  `DEADLINE_EXCEEDED`  |  `4`  | 操作超时                                             |
|      `NOT_FOUND`      |  `5`  | 找不到某些请求的实体（例如文件或目录）               |
|   `ALREADY_EXISTS`    |  `6`  | 客户端尝试创建的实体（例如，文件或目录）已存在       |
|  `PERMISSION_DENIED`  |  `7`  | 调用者没有执行指定操作的权限                         |
| `RESOURCE_EXHAUSTED`  |  `8`  | 某些资源已耗尽，可能是用户的配额或者文件系统空间不足 |
| `FAILED_PRECONDITION` |  `9`  | 操作被拒绝，因为系统未处于执行操作所需的状态         |
|       `ABORTED`       | `10`  | 通常因并发冲突而中止，例如事务冲突                   |
|    `OUT_OF_RANGE`     | `11`  | 操作已超出有效范围                                   |
|    `UNIMPLEMENTED`    | `12`  | 此服务中未实现或不支持/未启用该操作                  |
|      `INTERNAL`       | `13`  | 内部错误                                             |
|     `UNAVAILABLE`     | `14`  | 该服务当前不可用                                     |
|      `DATA_LOSS`      | `15`  | 无法恢复的数据丢失或损坏                             |
|   `UNAUTHENTICATED`   | `16`  | 请求没有操作的有效身份验证凭据                       |

### 7.2 服务端

下面用固定的 `NotFound` 状态演示返回业务错误，可替换第 3 节的 `SayHello`。需要导入 `codes` 和 `status` 包。

<<< @/go/codes/grpc/err_server.go

普通的非状态错误通常映射为 `Unknown`。取消和超时应使用对应状态，而非全部转换成 `Internal`。

### 7.3 客户端

辅助函数使用已有客户端，调用时传入带 deadline 的上下文：

<<< @/go/codes/grpc/err_cli.go

状态码适合程序判断，错误说明文本不应作为稳定的业务协议。`FailedPrecondition` 通常要求先改变系统状态，`Aborted` 常需要重做更高层的操作，`Unavailable` 才可能适合稍后重试，仍需结合幂等性和重试策略。即使收到 `DeadlineExceeded`，服务端也可能已经完成业务写入，不能据此判断操作没有发生。

## 8. 超时与取消

gRPC 默认不会自动给调用设置 deadline。客户端应根据业务耗时设置超时，并在结束后调用 `cancel` 释放计时器。为了演示超时，先用下面的实现替换第 3 节服务端的 `SayHello`，它等待 5 秒，但会响应上下文取消：

<<< @/go/codes/grpc/timeout_server.go

服务端方法需补充 `time` 和 `status` 导入。下面的客户端片段使用第 3 节创建的 `client`，设置 3 秒超时，调用该延迟服务端时会得到 `DeadlineExceeded`：

<<< @/go/codes/grpc/timeout.go

客户端取消上下文或超过 deadline 后，服务端上下文也会被取消。耗时操作应定期检查 `ctx.Err()` 或 `ctx.Done()`，否则业务 goroutine 可能在调用已经失败后继续占用资源。取消上下文不会强制终止 goroutine，也不会自动撤销已经发生的业务写入。服务端继续调用其他服务时应传递当前上下文，子上下文的 deadline 不能延长父上下文的有效期限。

## 9. 上线前检查

- **传输安全**：本机示例可以使用 `insecure.NewCredentials()`，跨主机部署应配置 TLS，并校验服务端身份。
- **优雅退出**：`GracefulStop` 会停止接收新连接和新 RPC，并等待在途调用结束，它本身没有超时参数。可在 goroutine 中执行，主流程用计时器控制等待，超时后调用 `Stop`。长期流式调用可能使等待一直持续，`Stop` 也不会强制杀死忽略取消的业务 goroutine。
- **健康检查**：标准健康检查服务需要显式注册，并根据实例状态更新 `SERVING` 等值。仅能建立 TCP 连接不代表服务已经就绪。
- **反射服务**：服务反射需要显式注册。`grpcurl` 可以使用反射，也可以通过 `-proto` 或 `-protoset` 提供定义。是否开启反射，应根据接口暴露范围决定。
- **请求边界**：为消息大小、并发数和执行时间设置合理限制，不能只依赖框架默认值。
