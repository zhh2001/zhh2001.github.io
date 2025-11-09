---
outline: [2, 3]
---

# gRPC

## 1 RPC

**RPC（Remote Procedure Call，远程过程调用）**​​是一种**网络通信协议**，它允许程序像调用本地方法一样调用远程服务器上的函数或服务。RPC 的核心思想是**隐藏网络通信的复杂性**，让开发者专注于业务逻辑。

### 1.1 基本原理

- **本地调用**：程序调用一个本地函数，直接执行并返回结果。
- **远程调用**：程序调用一个远程函数，RPC 框架会通过网络将调用请求发送到远程服务器，服务器执行函数并将结果返回给客户端。

### 1.2 工作流程

1. **客户端调用**：客户端像调用本地方法一样调用远程方法。
2. **序列化**：客户端将调用的方法名、参数等信息序列化为网络传输格式（如 JSON、Protobuf）。
3. **网络传输**：序列化后的数据通过网络发送到服务器。
4. **反序列化**：服务器接收到数据后，反序列化为可识别的格式。
5. **执行方法**：服务器根据方法名和参数执行对应的函数。
6. **返回结果**：服务器将执行结果序列化并返回给客户端。
7. **客户端接收**：客户端反序列化结果并继续执行后续逻辑。

### 1.3 常见框架

- **gRPC**：Google 开源的 RPC 框架，基于 HTTP/2 和 Protobuf，性能高，支持多种语言。
- **Thrift**：Apache 开源的 RPC 框架，支持多种语言和传输协议。
- **Dubbo**：阿里巴巴开源的 RPC 框架，主要用于 Java 生态。
- **JSON-RPC**：基于 JSON 的轻量级 RPC 协议，适合简单的应用场景。

### 1.4 RPC 服务端

Go 的 `net/rpc` 包提供了RPC的基本功能，支持 TCP 和 HTTP 协议。`net/rpc` 包使用 Go 的编码格式（Gob）进行数据的序列化和反序列化。

要开发一个 RPC 服务端，通常需要以下几个步骤：

1. **​定义服务类型**：定义一个结构体，该结构体的方法将成为 RPC 服务的方法。
2. **​注册服务**：将定义的服务类型注册到RPC服务中。
3. **​启动RPC服务**：监听某个端口，等待客户端的连接和调用。

<<< @/go/codes/grpc/rpc_server.go

### 1.5 RPC 客户端

1. **连接到RPC服务**：通过 TCP 或 HTTP 连接到 RPC 服务端。
2. **调用远程方法**：通过 RPC 客户端调用服务端提供的方法。

<<< @/go/codes/grpc/rpc_cli.go

## 2 Protocol Buffer

Protocol Buffers（protobuf）​​ 是 Google 开发的一种高效、跨语言的数据序列化协议，用于结构化数据的存储和通信。它通过预定义数据结构和自动生成代码的方式，提供了比 JSON、XML 等文本格式更高的性能和更小的数据体积。

### 2.1 安装

1. 在 [https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases) 下载后解压。
2. 将解压出来的 `bin` 目录配置到环境变量。
3. 执行下面的命令安装协议编译器插件：

	<<< @/go/codes/grpc/install.sh

### 2.2 定义数据结构

创建一个 `.proto` 文件（如 `hello.proto`）：

<<< @/go/codes/grpc/msg.proto

### 2.3 生成代码

使用 `protoc` 生成目标语言的代码：

<<< @/go/codes/grpc/protoc.sh

生成的文件（如 `hello.pb.go`）包含数据结构的序列化/反序列化方法。

如果 `hello.proto` 导入了其他的 `.proto` 文件，例如：

<<< @/go/codes/grpc/import.proto

则需要在使用 `protoc` 生成代码时指定这个文件所在的路径，假设它的路径为 `../third_party/google/api/annotations.proto`，那么命令为：

<<< @/go/codes/grpc/proto_path.sh{2,3}

### 2.4 数据类型

下表为 Proto 的数据类型，以及自动生成的对应语言的数据类型：

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
|   `sint32` |   `0`    | 使用变长编码。比 `int32` 编码负数更快                              |
|   `sint64` |   `0`    | 使用变长编码。比 `int64` 编码负数更快                              |
|  `fixed32` |   `0`    | 始终为 4 字节，如果值大于 2<sup>28</sup>，该类型比 `uint32` 高效   |
|  `fixed64` |   `0`    | 始终为 8 字节，如果值大于 2<sup>56</sup>，该类型比 `uint64` 高效   |
| `sfixed32` |   `0`    | 始终为 4 字节                                                      |
| `sfixed64` |   `0`    | 始终为 8 字节                                                      |
|     `bool` | `false`  | 布尔型                                                             |
|   `string` |   `""`   | 必须是 UTF-8 编码或者 7-bit ASCII 编码的文本                       |
|    `bytes` | `[]byte` | 可以包含任意字节数据的序列                                         |

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

字段 `SearchRequest.corpus` 的默认值是 `CORPUS_UNSPECIFIED`，因为这是枚举中定义的第一个值。枚举定义的第一个值**必须**为 0。

枚举字段的默认值可以显式覆盖，如下所示：

<<< @/go/codes/grpc/default.proto{5}

### 2.9 Map 类型

<<< @/go/codes/grpc/map.proto

`key` 的类型可以是任何整数或字符串类型。`value` 的类型可以是任意类型。

### 2.10 时间戳类型

<<< @/go/codes/grpc/ts.proto

`google/protobuf/timestamp.proto` 的源码如下：

<<< @/go/codes/grpc/timestamp.proto

在使用 Go 实现时，需要导入 `google.golang.org/protobuf/types/known/timestamppb` 包来生成时间戳。

## 3 gRPC

gRPC 是一个开源高性能 RPC 框架，可以在任何环境中运行。

### 3.1 定义服务

<<< @/go/codes/grpc/hello.proto

### 3.2 生成代码

运行 `protoc` 命令生成 Go 代码。会生成两个文件：

- `hello.pb.go`：消息定义
- `hello_grpc.pb.go`：gRPC 服务端/客户端接口

### 3.3 实现服务端

<<< @/go/codes/grpc/grpc_server.go

### 3.3 实现客户端

<<< @/go/codes/grpc/grpc_cli.go

## 4 数据流模式

在 gRPC 中，支持 **​4 种数据流模式**。

### 4.1 单一请求-响应

客户端发送一个请求，服务端返回一个响应。

proto 定义：

<<< @/go/codes/grpc/service.proto

Go 实现：

<<< @/go/codes/grpc/service.go

### 4.2 服务端流式响应

客户端发送**一个请求**，服务端返回**多个响应**​（流式传输）。

proto 定义：

<<< @/go/codes/grpc/stream_server.proto

Go 实现：

<<< @/go/codes/grpc/stream_server.go

### 4.3 客户端流式请求

客户端发送**多个请求**，服务端返回**一个响应**。

proto 定义：

<<< @/go/codes/grpc/stream_cli.proto

Go 实现：

<<< @/go/codes/grpc/stream_cli.go

### 4.4 双向流式

客户端和服务端**同时发送多个请求/响应**，彼此独立。

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

## 5 元数据

gRPC 元数据是使用 HTTP/2 标头实现的。

### 5.1 构建元数据

可以使用包 `google.golang.org/grpc/metadata` 创建元数据。类型 `MD` 实际上是 `map`：

<<< @/go/codes/grpc/md.go

元数据可以像 `map` 一样使用，其中值类型为 `[]string`。

#### 5.1.1 创建新元数据

可以使用函数 `New` 从 `map[string]string` 创建元数据：

<<< @/go/codes/grpc/mdNew.go

另一种方法是使用 `Pairs`。具有相同键的值将被合并到一个切片中：

<<< @/go/codes/grpc/mdPairs.go

::: tip 注意
所有的键都会自动转换为小写， 因此 `key1` 和 `kEy1` 是相同的键，它们的值将被合并到同一个切片中。无论是 `New` 还是 `Pairs` 都是这样的。
:::

#### 5.1.2 在元数据中存储二进制数据

在元数据中，键始终是字符串。但值可以是字符串或二进制数据。 要在元数据中存储二进制数据值，只需在键中添加 `-bin` 后缀即可。 创建元数据时，将对带有 `-bin` 后缀键的值进行编码：

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

流式调用：

<<< @/go/codes/grpc/mdSomeStreamingRPC.go

### 5.3 服务端发送和接收

#### 5.3.1 接收元数据

如果是一元调用，则可以使用 RPC 处理程序的上下文。对于流式调用，服务端需要从流中获取上下文。

一元调用：

```go
func (s *server) SomeRPC(ctx context.Context, in *pb.someRequest) (*pb.someResponse, error) {
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

如果采用一元调用方式，服务端可以调用 `grpc` 模块中的 `SetHeader` 和 `SetTrailer` 函数。这两个函数将上下文作为第一个参数，它是 RPC 处理程序的上下文或从它派生的上下文：

<<< @/go/codes/grpc/someRPC.go

对于流式调用，可以使用接口 `ServerStream` 中的函数 `SetHeader` 和 `SetTrailer` 发送 `header` 和 `trailer`：

<<< @/go/codes/grpc/someStreamingRPC.go

::: tip 重要
不要在服务端使用 `FromOutgoingContext` 写入元数据。`FromOutgoingContext` 仅供客户端使用。
:::

## 6 拦截器

gRPC 提供了简单的 API，可以在每个 ClientConn/Server 的基础上实现和安装拦截器。拦截器充当应用程序和 gRPC 的中间层，可用于观察或控制 gRPC 的行为。拦截器可用于日志记录、身份验证/授权、指标收集以及跨 RPC 共享的其他功能。客户端和服务端都有自己的一元拦截器和流拦截器。因此，总共有四种不同类型的拦截器。

### 6.1 客户端拦截器

#### 6.1.1 一元拦截器

客户端一元拦截器的类型为 `UnaryClientInterceptor`。它本质上是一个带有签名的函数类型：

<<< @/go/codes/grpc/unaryClientInterceptor.go

一元拦截器实现通常可分为三个部分：预处理、调用RPC方法和后处理。

对于预处理，可以通过检查传入的参数来获取有关当前 RPC 调用的信息。参数包括 RPC 上下文、方法字符串、要发送的请求和配置的 `CallOptions`。有了这些信息甚至可以修改 RPC 调用。

预处理后，用户可以通过调用 `invoker` 来调用 RPC 调用。

一旦调用程序返回，就可以对 RPC 调用进行后处理。这通常涉及处理返回的回复和错误。

要在 ClientConn 上安装一元拦截器，需要使用 `WithUnaryInterceptor` 配置 `NewClient(target string, opts ...DialOption)` 的 `DialOption`。

#### 6.1.2 流拦截器

客户端流拦截器的类型是 `StreamClientInterceptor`。它是一个带有签名的函数类型：

<<< @/go/codes/grpc/streamClientInterceptor.go

流拦截器的实现通常包括预处理和流操作拦截。

预处理类似于一元拦截器。

然而，流拦截器不是调用 RPC 方法然后进行后处理，而是拦截对流的操作。拦截器首先调用传入的 `streamer` 以获取 `ClientStream`，然后在用拦截逻辑重载其方法的同时封装 `ClientStream`。最后，拦截器将包装好的 `ClientStream` 返回。

要为 ClientConn 安装流拦截器，需要使用 `WithStreamInterceptor` 配置 `NewClient(target string, opts ...DialOption)` 的 `DialOption`。

#### 6.1.3 示例

<<< @/go/codes/grpc/streamInterceptor.go{1-7,13}

### 6.2 服务端拦截器

服务端拦截器与客户端拦截器相似，但提供的参数信息略有不同。

#### 6.2.1 一元拦截器

服务一元拦截器的类型是 `UnaryServerInterceptor`。它是一个带有签名的函数类型：

<<< @/go/codes/grpc/unaryServerInterceptor.go

要在服务端安装一元拦截器，需要使用 `UnaryInterceptor` 配置 `NewServer(opt ...ServerOption)` 的 `ServerOption`。

#### 6.2.2 流拦截器

服务端流拦截器的类型是 `StreamServerInterceptor`。它是一个具有签名的函数类型：

<<< @/go/codes/grpc/streamServerInterceptor.go

要在服务端安装流拦截器，需要使用 `StreamInterceptor` 配置 `NewServer(opt ...ServerOption)` 的 `ServerOption`。

#### 6.2.3 示例

<<< @/go/codes/grpc/serverInterceptor.go

## 7 错误处理

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

## 8 超时

客户端设置一个 3 秒超时：

<<< @/go/codes/grpc/timeout.go
