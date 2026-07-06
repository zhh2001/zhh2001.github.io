---
outline: [2, 3]
---

# gRPC

gRPC 是一套基于服务定义生成客户端与服务端代码的 RPC 框架。Go 项目通常用 Protocol Buffers 描述消息和服务，调用数据通过 HTTP/2 连接传输。本文先用 `net/rpc` 说明 RPC 的基本过程，再进入 Protobuf、流式调用、元数据、拦截器和错误处理。

## 1 RPC

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

Go 标准库的 `net/rpc` 可以通过 TCP 或 HTTP 传输 RPC 消息，默认使用 Gob 编码。它只面向 Go 程序，并且已经进入冻结状态，适合用来理解 RPC 或处理简单的内部通信；跨语言接口通常使用 gRPC 等方案。

要开发一个 RPC 服务端，通常需要以下几个步骤：

1. **定义服务类型**：结构体的导出方法作为远程方法。
2. **注册服务**：把服务实例注册到 RPC 服务中。
3. **启动服务**：监听端口并处理连接。

<<< @/go/codes/grpc/rpc_server.go

### 1.4 `net/rpc` 客户端

1. **连接到 RPC 服务**：通过 TCP 或 HTTP 连接服务端。
2. **调用远程方法**：通过 RPC 客户端调用服务端提供的方法。

<<< @/go/codes/grpc/rpc_cli.go

## 2 Protocol Buffers

Protocol Buffers（Protobuf）同时包含接口描述语言、二进制编码格式和代码生成工具。它可以脱离 gRPC 单独使用；gRPC 只是把 Protobuf 作为默认的服务定义与消息格式。

### 2.1 安装

1. 从 [Protobuf Releases](https://github.com/protocolbuffers/protobuf/releases) 安装 `protoc`，并确认它位于 `PATH`：

	```sh
	protoc --version
	```

2. 安装 Go 消息代码和 gRPC 服务代码的生成插件：

	<<< @/go/codes/grpc/install.sh

3. 确保 Go 的可执行文件目录位于 `PATH`：

	```sh
	export PATH="$PATH:$(go env GOPATH)/bin"
	```

实际项目应在构建脚本中固定工具版本，避免 `@latest` 更新后生成结果发生意外变化。

### 2.2 定义数据结构

创建一个 `.proto` 文件，例如 `msg.proto`：

<<< @/go/codes/grpc/msg.proto

字段后的数字是线上的字段编号，不是数组下标。编号一旦发布就不应随意修改或分配给其他字段。`go_package` 应填写生成代码的完整 Go 导入路径，并与项目的模块路径保持一致。

### 2.3 生成 Go 代码

使用 `protoc` 生成目标语言的代码：

<<< @/go/codes/grpc/protoc.sh

`--go_out` 生成消息类型，`--go-grpc_out` 生成 gRPC 客户端和服务端接口。`paths=source_relative` 表示输出文件与输入的 `.proto` 文件放在同一相对目录下。

如果 `hello.proto` 导入了其他的 `.proto` 文件，例如：

<<< @/go/codes/grpc/import.proto

则需要通过 `--proto_path` 指定导入根目录。假设文件位于 `../third_party/google/api/annotations.proto`，应把 `../third_party` 作为搜索路径：

<<< @/go/codes/grpc/proto_path.sh{2,3}

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

|   数据类型 |  默认值  | 说明                                                               |
| ---------: | :------: | :----------------------------------------------------------------- |
|   `double` |   `0`    | 双精度浮点型                                                       |
|    `float` |   `0`    | 浮点型                                                             |
|    `int32` |   `0`    | 使用变长编码。负数编码效率低下。如果字段可能为负，用 `sint32` 代替 |
|    `int64` |   `0`    | 使用变长编码。负数编码效率低下。如果字段可能为负，用 `sint64` 代替 |
|   `uint32` |   `0`    | 使用变长编码                                                       |
|   `uint64` |   `0`    | 使用变长编码                                                       |
|   `sint32` |   `0`    | 使用 ZigZag 编码，负数通常比 `int32` 更紧凑                        |
|   `sint64` |   `0`    | 使用 ZigZag 编码，负数通常比 `int64` 更紧凑                        |
|  `fixed32` |   `0`    | 始终为 4 字节，如果值大于 2<sup>28</sup>，该类型比 `uint32` 高效   |
|  `fixed64` |   `0`    | 始终为 8 字节，如果值大于 2<sup>56</sup>，该类型比 `uint64` 高效   |
| `sfixed32` |   `0`    | 始终为 4 字节                                                      |
| `sfixed64` |   `0`    | 始终为 8 字节                                                      |
|     `bool` | `false`  | 布尔型                                                             |
|   `string` |   `""`   | 必须是 UTF-8 编码的文本                                            |
|    `bytes` |  `nil`   | Go 中为 `[]byte`，读取语义为空字节序列                             |

消息类型字段的默认值为 `nil`，枚举字段取编号为 0 的值。Proto3 不允许使用 `[default = ...]` 声明自定义默认值；业务默认值应在应用代码中处理。需要区分“字段未设置”和“字段显式设置为零值”时，可以使用 `optional`：

```proto
optional int32 page_size = 1;
```

### 2.6 消息嵌套

可以使用其他消息类型作为字段类型。例如，要在每条 `SearchResponse` 消息中包含 `Result` 消息，可以直接在同一个 `.proto` 文件中定义一个 `Result` 消息类型，然后在 `SearchResponse` 中指定一个字段类型为 `Result`：

<<< @/go/codes/grpc/res.proto

也可以直接在 `SearchResponse` 消息类型中直接定义和使用 `Result` 消息类型：

<<< @/go/codes/grpc/res2.proto

如果要在父消息类型之外重用此消息类型，需要使用 `_Parent_._Type_` 进行调用：

<<< @/go/codes/grpc/parent.proto

也可以根据需要嵌套多层消息，下面的例子中，两个名为 `Inner` 的嵌套类型是完全独立的，因为它们定义在不同的消息中：

<<< @/go/codes/grpc/inner.proto

### 2.7 导入定义

可以通过导入其他 `.proto` 文件来使用里面的定义：

```proto
import "myproject/other_protos.proto";
```

默认情况下，只能使用直接导入的 `.proto` 文件中的定义。但是通过 `import public` 导入的 `.proto` 文件可以传递到下一个文件。例如：

```proto
// new.proto
// 这里有很多定义
```

```proto
// old.proto
import public "new.proto";
import "other.proto";

// 这里有很多定义
```

```proto
// client.proto
import "old.proto";

// 可以使用 old.proto 和 new.proto 中的定义，但是不能使用 other.proto 的定义
```

### 2.8 枚举类型

<<< @/go/codes/grpc/enum.proto

字段 `SearchRequest.corpus` 的默认值是 `CORPUS_UNSPECIFIED`，因为它的编号为 0。Proto3 枚举定义的第一个值必须为 0，通常用 `*_UNSPECIFIED` 表示调用方没有明确选择。

### 2.9 Map 类型

<<< @/go/codes/grpc/map.proto

`key` 可以使用整数、`bool` 或字符串类型，不能使用浮点数、`bytes` 和枚举。`value` 可以使用除 `map` 以外的任意字段类型。

### 2.10 时间戳类型

<<< @/go/codes/grpc/ts.proto

`google/protobuf/timestamp.proto` 的源码如下：

<<< @/go/codes/grpc/timestamp.proto

在使用 Go 实现时，需要导入 `google.golang.org/protobuf/types/known/timestamppb` 包来生成时间戳。

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

不要修改已经发布字段的编号，也不要把同一编号改成含义不同的字段。消息兼容性依赖编号，而不是字段名称。

## 3 gRPC

gRPC 根据 `.proto` 中的 `service` 生成强类型客户端和服务端接口。一次调用通常包含方法路径、消息、元数据、状态码和可选的流数据；Protobuf 负责消息编码，gRPC 负责调用语义与传输。

### 3.1 定义服务

<<< @/go/codes/grpc/hello.proto

### 3.2 生成代码

运行 `protoc` 命令生成 Go 代码。会生成两个文件：

- `hello.pb.go`：消息定义
- `hello_grpc.pb.go`：gRPC 服务端/客户端接口

示例模块可以这样初始化依赖：

```sh
go mod init example.com/grpcdemo
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

### 3.3 实现服务端

<<< @/go/codes/grpc/grpc_server.go

服务结构体嵌入 `UnimplementedGreeterServer`，可以在后续新增 RPC 时保持向前兼容。正式代码还应处理监听和 `Serve` 返回的错误。

### 3.4 实现客户端

<<< @/go/codes/grpc/grpc_cli.go

`grpc.NewClient` 只创建逻辑上的客户端连接，不会立即发起网络 I/O；通常在第一次 RPC 时才开始解析地址并建立连接。示例使用 `insecure.NewCredentials()` 便于本机测试，它不会启用 TLS，不能直接照搬到生产环境。

## 4 数据流模式

gRPC 有四种调用模式。单个流内的消息保持发送顺序，但不同 RPC 之间没有统一的到达顺序。当前版本的 `protoc-gen-go-grpc` 会为流式方法生成带泛型参数的接口；旧版生成代码中的接口名称可能不同。

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

客户端和服务端都可以连续发送消息。gRPC-Go 允许一个 goroutine 负责 `Send`、另一个 goroutine 负责 `Recv`，但不能让多个 goroutine 同时调用同一方向的方法。

proto 定义：

<<< @/go/codes/grpc/stream.proto

Go 实现：

<<< @/go/codes/grpc/stream.go

### 4.5 示例

`.proto` 文件定义：

<<< @/go/codes/grpc/stream_example.proto

服务端：

<<< @/go/codes/grpc/stream_server_example.go

客户端：

<<< @/go/codes/grpc/stream_client_example.go

流结束时要区分 `io.EOF` 与真正的传输错误。服务端处理函数应返回错误，不要在处理请求的 goroutine 中直接 `panic`，否则一次异常调用可能终止整个进程。

## 5 元数据

元数据用于传递认证信息、追踪 ID 等调用上下文，在线上对应 HTTP/2 的 header 和 trailer。键不区分大小写，不能以保留前缀 `grpc-` 开头；普通值应使用 ASCII，二进制键以 `-bin` 结尾。

### 5.1 构建元数据

可以使用包 `google.golang.org/grpc/metadata` 创建元数据。类型 `MD` 实际上是 `map`：

<<< @/go/codes/grpc/md.go

元数据可以像 `map` 一样读取，其中每个键对应一组值。它不适合承载大块业务数据，过大的 header 可能被客户端、代理或服务端拒绝。

#### 5.1.1 创建新元数据

可以使用函数 `New` 从 `map[string]string` 创建元数据：

<<< @/go/codes/grpc/mdNew.go

另一种方法是使用 `Pairs`。具有相同键的值将被合并到一个切片中：

<<< @/go/codes/grpc/mdPairs.go

::: tip 注意
所有的键都会自动转换为小写， 因此 `key1` 和 `kEy1` 是相同的键，它们的值将被合并到同一个切片中。无论是 `New` 还是 `Pairs` 都是这样的。
:::

#### 5.1.2 在元数据中存储二进制数据

在 Go API 中，元数据的键和值都是字符串。传递原始字节时，把字节转换为字符串，并给键添加 `-bin` 后缀；gRPC 会处理线上的二进制编码：

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

流式调用中，`Header` 会等待响应头到达；`Trailer` 应在读取到 `io.EOF` 后获取：

<<< @/go/codes/grpc/mdSomeStreamingRPC.go

### 5.3 服务端发送和接收

#### 5.3.1 接收元数据

如果是一元调用，则可以使用 RPC 处理程序的上下文。对于流式调用，服务端需要从流中获取上下文。

一元调用：

```go
func (s *server) SomeRPC(ctx context.Context, in *pb.SomeRequest) (*pb.SomeResponse, error) {
    md, ok := metadata.FromIncomingContext(ctx)
}
```

流式调用：

```go
func (s *server) SomeStreamingRPC(stream pb.Service_SomeStreamingRPCServer) error {
    md, ok := metadata.FromIncomingContext(stream.Context())
}
```

#### 5.3.2 发送元数据

一元调用可以使用 `grpc.SetHeader` 和 `grpc.SetTrailer`。它们接收当前 RPC 的上下文，返回的错误不应忽略：

<<< @/go/codes/grpc/someRPC.go

对于流式调用，可以使用接口 `ServerStream` 中的函数 `SetHeader` 和 `SetTrailer` 发送 `header` 和 `trailer`：

<<< @/go/codes/grpc/someStreamingRPC.go

::: tip 重要
服务端不能通过 `metadata.NewOutgoingContext` 发送响应元数据，应使用 `grpc.SetHeader`、`grpc.SetTrailer` 或流对象上的对应方法。
:::

## 6 拦截器

拦截器位于应用代码与 gRPC 调用之间，常用于日志、认证、指标和追踪。客户端与服务端分别提供一元和流式拦截器，共四种类型。拦截器必须继续调用传入的 `invoker`、`streamer` 或 `handler`，除非它明确要提前拒绝这次调用。

### 6.1 客户端拦截器

#### 6.1.1 一元拦截器

客户端一元拦截器的类型为 `UnaryClientInterceptor`。它本质上是一个带有签名的函数类型：

<<< @/go/codes/grpc/unaryClientInterceptor.go

一元拦截器通常分为三步：预处理、调用 RPC 方法和后处理。

对于预处理，可以通过检查传入的参数来获取有关当前 RPC 调用的信息。参数包括 RPC 上下文、方法字符串、要发送的请求和配置的 `CallOptions`。有了这些信息甚至可以修改 RPC 调用。

预处理后，用户可以通过调用 `invoker` 来调用 RPC 调用。

一旦调用程序返回，就可以对 RPC 调用进行后处理。这通常涉及处理返回的回复和错误。

要在 `ClientConn` 上安装一元拦截器，可以给 `NewClient` 传入 `WithUnaryInterceptor`。需要组合多个拦截器时使用 `WithChainUnaryInterceptor`，执行顺序与注册顺序一致。

#### 6.1.2 流拦截器

客户端流拦截器的类型是 `StreamClientInterceptor`。它是一个带有签名的函数类型：

<<< @/go/codes/grpc/streamClientInterceptor.go

流拦截器的实现通常包括预处理和流操作拦截。

预处理类似于一元拦截器。

然而，流拦截器不是调用 RPC 方法然后进行后处理，而是拦截对流的操作。拦截器首先调用传入的 `streamer` 以获取 `ClientStream`，然后在用拦截逻辑重载其方法的同时封装 `ClientStream`。最后，拦截器将包装好的 `ClientStream` 返回。

流拦截器通过 `WithStreamInterceptor` 安装，多个流拦截器可以使用 `WithChainStreamInterceptor` 组合。

#### 6.1.3 示例

<<< @/go/codes/grpc/streamInterceptor.go{1-7,13}

### 6.2 服务端拦截器

服务端拦截器与客户端拦截器相似，但提供的参数信息略有不同。

#### 6.2.1 一元拦截器

服务端一元拦截器的类型是 `UnaryServerInterceptor`。它是一个带有签名的函数类型：

<<< @/go/codes/grpc/unaryServerInterceptor.go

服务端一元拦截器通过 `UnaryInterceptor` 安装，多个拦截器使用 `ChainUnaryInterceptor`。

#### 6.2.2 流拦截器

服务端流拦截器的类型是 `StreamServerInterceptor`。它是一个具有签名的函数类型：

<<< @/go/codes/grpc/streamServerInterceptor.go

服务端流拦截器通过 `StreamInterceptor` 安装，多个拦截器使用 `ChainStreamInterceptor`。

#### 6.2.3 示例

<<< @/go/codes/grpc/serverInterceptor.go

## 7 错误处理

gRPC 错误由状态码、说明文本和可选详情组成。服务端应返回稳定、可判断的状态码，不要把数据库错误或内部堆栈直接暴露给客户端。客户端先判断状态码，再决定提示用户、重试还是终止调用。

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
|       `ABORTED`       | `10`  | 操作被中止                                           |
|    `OUT_OF_RANGE`     | `11`  | 操作已超出有效范围                                   |
|    `UNIMPLEMENTED`    | `12`  | 此服务中未实现或不支持/未启用该操作                  |
|      `INTERNAL`       | `13`  | 内部错误                                             |
|     `UNAVAILABLE`     | `14`  | 该服务当前不可用                                     |
|      `DATA_LOSS`      | `15`  | 无法恢复的数据丢失或损坏                             |
|   `UNAUTHENTICATED`   | `16`  | 请求没有操作的有效身份验证凭据                       |

### 7.2 服务端

<<< @/go/codes/grpc/err_server.go{2}

### 7.3 客户端

<<< @/go/codes/grpc/err_cli.go

## 8 超时与取消

gRPC 默认不会自动给调用设置 deadline。客户端应根据业务耗时设置超时，并在结束后调用 `cancel` 释放计时器：

<<< @/go/codes/grpc/timeout.go

客户端取消上下文或超过 deadline 后，服务端上下文也会被取消。耗时操作应定期检查 `ctx.Err()` 或 `ctx.Done()`，否则业务 goroutine 可能在调用已经失败后继续占用资源。

## 9 上线前检查

- **传输安全**：本机示例可以使用 `insecure.NewCredentials()`，跨主机部署应配置 TLS，并校验服务端身份。
- **优雅退出**：停止接收新请求后调用 `GracefulStop`，给正在执行的 RPC 留出收尾时间；超时后再决定是否强制 `Stop`。
- **健康检查**：使用 gRPC 标准健康检查服务，不要用一个业务 RPC 代替实例健康状态。
- **反射服务**：`grpcurl` 等工具依赖服务反射。生产环境是否开启，应根据接口暴露范围决定。
- **请求边界**：为消息大小、并发数和执行时间设置合理限制，不能只依赖框架默认值。
