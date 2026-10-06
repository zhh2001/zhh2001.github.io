---
outline: deep
---

# Eino

Eino 是 CloudWeGo 的 Go 大模型应用开发框架。它把模型、提示模板、检索和工具等能力抽象为组件，再通过 Chain、Graph 等编排方式连接这些组件。本文整理组件接口和基础编排，示例使用 `schema.Message`，不展开 ADK 和 `AgenticMessage` 的用法。

本文以 Go 1.26.x 和 Eino v0.9.21 为基准，接口声明按该版本整理。`eino-ext` 的各组件是独立 Go 模块，版本号与 Eino 核心库不同。具体 API 可核对 [Eino v0.9.21 文档](https://pkg.go.dev/github.com/cloudwego/eino@v0.9.21)。

## 1. 环境与依赖

创建独立的示例目录，在该目录初始化模块并安装依赖：

<<< @/go/codes/eino/install.sh

其中 Ark ChatModel 使用 v0.1.71，Ark Embedding 使用 v0.1.2，CozeLoop 回调扩展使用 v0.3.1。Redis 索引、检索、Markdown 分割和 HTTP 工具模块使用同一个 Go 伪版本号，它们尚未提供可用于本例的语义版本标签。

每个包含 `package main` 的示例都是独立程序。例如，将生成示例保存为 `cmd/generate/main.go`，然后在模块根目录执行 `go run ./cmd/generate`。不同示例不能全部放进同一个包，否则 `main` 等名称会重复。接口和结构体摘录用于说明框架 API，其中的 `Option` 等类型属于相应组件包，不需要在业务代码里重新定义。

模型相关示例从进程环境读取配置，不会自动加载 `.env`。运行前设置相应变量：

| 变量                | 含义                                                                                                  |
| ------------------- | ----------------------------------------------------------------------------------------------------- |
| `ARK_API_KEY`       | Ark 服务的 API Key                                                                                    |
| `MODEL`             | 对话模型或推理接入点的 ID                                                                             |
| `EMBEDDER`          | 嵌入模型或推理接入点的 ID                                                                             |
| `ARK_BASE_URL`      | 可选，所用地域的 API 基础地址，留空使用 SDK 默认的北京地域地址                                        |
| `EMBEDDER_API_TYPE` | 可选，`text_api` 对应 `/embeddings`，`multi_modal_api` 对应 `/embeddings/multimodal`，默认 `text_api` |
| `REDIS_ADDR`        | Redis 的 `host:port`，例如 `127.0.0.1:6379`                                                           |
| `REDIS_PASSWORD`    | 可选，Redis 密码                                                                                      |

`MODEL` 和 `EMBEDDER` 是这里约定的环境变量名，不是 Eino 自动识别的配置。应填写账号中可用的模型或接入点 ID，并选择该模型支持的接口类型。模板、文档分割、自定义工具和不含模型的编排示例无需 Ark 凭据。

## 2. ChatModel 组件

### 2.1 接口

ChatModel 抽象模型调用。`BaseChatModel` 是 `BaseModel[*schema.Message]` 的类型别名，包含 `Generate` 和 `Stream`。`WithTools` 属于扩展接口 `ToolCallingChatModel`，并非每个基础模型实现都支持它。

<<< @/go/codes/eino/iBaseChatModel.go

`Generate` 返回一条完整响应消息。`Stream` 返回流读取器，后续响应或错误从 `Recv` 获取。两者都接收消息列表、上下文和模型选项。上下文用于取消、超时和传递回调信息，选项的支持情况还取决于模型实现。

接口可以容纳多模态消息，但模型和服务适配器是否支持图片、音频等输入输出，需要分别确认。

### 2.2 `schema.Message`

<<< @/go/codes/eino/sMessage.go

`Role` 区分系统指令、用户输入、模型回复和工具结果。用户问题应放在 `user` 消息中，不应因经过模板处理就改成 `system` 消息。`ToolCalls` 描述模型请求的工具调用，工具结果通过 `ToolCallID` 与对应调用关联。

`MultiContent` 已弃用。用户的多模态输入使用 `UserInputMultiContent`，模型的多模态输出使用 `AssistantGenMultiContent`。`ReasoningContent` 只在服务实际返回相应内容时有值，不能据此认为可以获得模型完整的内部推理过程。

### 2.3 完整生成

<<< @/go/codes/eino/eGenerate.go

示例使用 30 秒调用超时，模型响应的具体内容不固定。返回错误时应先处理错误，再读取响应。

### 2.4 流式生成

<<< @/go/codes/eino/eStream.go

流块通常是增量内容，不是截至当前时刻的完整答案，也不保证按词或句子划分。只打印 `Content` 适合展示文本，工具调用参数、推理内容和多模态结果还可能位于其他字段中。

`Recv` 返回 `io.EOF` 表示流正常结束，其他错误应返回给调用方。无论是否读到末尾，都要关闭读取器。需要合并消息时使用 `schema.ConcatMessages`，不能只拼接文本就丢掉其他字段。一个读取器只有一个消费进度，多处消费应先复制流，并分别关闭各自的读取器。

### 2.5 绑定工具

下面的辅助函数接收支持工具调用的模型和工具实例，取得工具描述后，调用 `WithTools` 派生一个绑定工具的新模型。它复用第 9 节创建的工具，需要导入 `model`、`tool` 和 `schema` 包。

<<< @/go/codes/eino/eWithTools.go

`WithTools` 不修改原实例。绑定工具后，模型可能返回 `ToolCalls`，也可能直接回答。模型生成调用名称和参数，并不执行 Go 函数，执行过程由应用代码或 ToolsNode 完成。

## 3. ChatTemplate 组件

### 3.1 接口与格式化

<<< @/go/codes/eino/iChatTemplate.go

`prompt.FromMessages` 把实现 `schema.MessagesTemplate` 的对象组合为模板。`*schema.Message` 实现了相应格式化方法，消息构造函数返回的指针可以直接使用。`Format` 接收变量映射，输出 `[]*schema.Message`，它本身不会调用模型。

`SystemMessage`、`UserMessage`、`AssistantMessage` 和 `ToolMessage` 用于构造不同角色的消息。`AssistantMessage` 还接收工具调用列表，`ToolMessage` 需要对应的工具调用 ID。

### 3.2 示例

<<< @/go/codes/eino/eChatTemplate.go

这里使用 `schema.FString`，`{role}` 和 `{topic}` 从变量映射取值。需要输出字面量花括号时，将对应的花括号写两次。缺少必需的格式化变量会返回错误。

`MessagesPlaceholder("history", true)` 把历史消息列表插入模板，第二个参数表示该变量可以缺省。传入历史时，值应为 `[]*schema.Message`。模板不会自动保存历史，需要由应用准备并传入。

## 4. RAG

RAG（Retrieval-Augmented Generation，检索增强生成）在生成前检索外部资料，再把相关片段作为上下文交给模型。它可以补充训练知识之外的信息，也有助于减少缺乏依据的回答，但不能保证消除幻觉。结果仍取决于资料质量、检索覆盖和模型是否正确使用资料。

典型流程分为两个阶段：

1. 建立索引。收集资料，用 Document Transformer 清理或切分文档。采用向量检索时，用 Embedding 生成向量，再由 Indexer 写入后端。
2. 检索与生成。Retriever 根据用户问题取回文档，应用整理来源和上下文，通过 ChatTemplate 或消息列表交给 ChatModel。

RAG 不要求必须使用向量数据库，也可以采用关键词检索、稀疏检索或混合检索。向量只是其中一种方式。检索结果是否及时更新，取决于资料和索引的更新流程，并不天然具有实时性。

下面的 Redis 示例演示索引和检索两个步骤，还没有把检索结果传给模型生成最终答案。另有 [RAG_Learning 练习项目](https://github.com/zhh2001/RAG_Learning)，使用时应核对其自身依赖。

## 5. Embedding 组件

### 5.1 接口

<<< @/go/codes/eino/iEmbedding.go

`EmbedStrings` 将文本列表转换为向量列表，`[][]float64` 的外层对应输入文本，内层是各维的值。向量维度、适用语言和检索效果由所选模型决定。

### 5.2 示例

<<< @/go/codes/eino/eEmbedding.go

文本索引和查询应使用同一个嵌入模型，并保持接口类型、维度及其他相关配置一致。两个模型即使输出维度相同，向量空间也不一定兼容。更换模型后通常需要重新生成文档向量。

语义接近的文本在合适的向量空间中通常较接近，但具体排序还取决于余弦距离、内积或欧氏距离等度量，不能只凭向量维度判断检索质量。

## 6. Indexer 组件

### 6.1 接口

<<< @/go/codes/eino/iIndexer.go

`Store` 接收文档并返回写入后的 ID。Indexer 负责与索引后端交互，并不局限于向量数据库。是否生成向量、如何映射字段以及写入后的检索可见性，由具体实现决定。

### 6.2 Redis 示例的条件

本例需要支持 `FT.CREATE` 和 `FT.SEARCH` 的 Redis Search/Query Engine。可使用包含这些功能的 Redis 8.2.x 发行包，并确认服务已加载查询引擎。只有普通数据命令可用的 Redis 服务不足以运行这个示例。

Redis 扩展把嵌入结果转换为 FLOAT32 二进制向量，因此索引的向量字段也使用 `FLOAT32`。以下配置必须对应：

| 项目        | 本例配置                   |
| ----------- | -------------------------- |
| 索引名称    | `eino_notes`               |
| Hash 键前缀 | `eino:note:`               |
| 文本字段    | `content`                  |
| 向量字段    | `vector_content`           |
| 向量维度    | 从所选模型的实际响应中读取 |
| 距离度量    | `COSINE`                   |

### 6.3 写入文档

<<< @/go/codes/eino/eIndexer.go

程序先生成一个探测向量，用其长度创建索引，再写入三条文档。索引已经存在时，创建步骤会返回错误。复用已有索引前，应核对字段、维度、数据类型和距离度量，不能仅凭索引同名就判断配置相同。

这里的 `Store` 返回文档 ID，例如 `p4`，实际 Hash 键为 `eino:note:p4`。同一前缀下重复写入相同 ID 会更新对应 Hash 中的字段。示例使用 `Protocol: 2` 让查询结果按 RESP2 解析，不需要同时启用 `UnstableResp3`。

## 7. Retriever 组件

### 7.1 接口与选项

<<< @/go/codes/eino/iRetriever.go

`Retrieve` 根据查询返回 `[]*schema.Document`。它可以使用向量、关键词或其他检索方式，公共选项如下：

<<< @/go/codes/eino/sRetrieverOptions.go

`TopK` 控制返回数量。`Index`、`SubIndex`、`DSLInfo` 等选项的含义和支持情况由后端实现决定。`ScoreThreshold` 也不是所有检索器通用的“相似度下限”，不能把某个后端的数值解释直接套到另一个后端。

### 7.2 Redis 检索

先运行上一节的索引程序，再运行检索程序：

<<< @/go/codes/eino/eRetriever.go

示例用 KNN 查询返回最多两条文档，并只取回 `content` 字段。该扩展返回的文档 ID 是实际 Redis 键，包含 `eino:note:` 前缀。检索顺序取决于模型的向量响应和距离，正文不固定具体命中的文档。

本文固定的 Redis 扩展使用 `RetrieverConfig.DistanceThreshold` 配置向量范围查询，距离越小表示越接近。它不会把公共 `WithScoreThreshold` 选项作为这个配置的替代。需要限制距离时，应配置 `DistanceThreshold`，而不是按“分数必须大于 0.5”理解阈值。

## 8. Document Transformer 组件

### 8.1 接口

<<< @/go/codes/eino/iTransformer.go

Transformer 接收并返回文档列表，用于切分、清理、过滤或其他转换。它不等同于嵌入模型，也不保证每个实现都会生成向量。

### 8.2 按 Markdown 标题切分

<<< @/go/codes/eino/eTransformer.go

`Headers` 把标题级别映射为元数据字段。`TrimHeaders: true` 从输出正文中移除匹配的标题行，标题仍可保存在 `h1`、`h2`、`h3` 元数据中。

标题切分并不保证每块都小于模型的 token 上限，较长章节还需要进一步按长度切分。分块后写入索引时，也需要为每个块分配合适的文档 ID。

## 9. Tool 与 ToolsNode

### 9.1 工具接口

Tool 描述一项可执行能力。ToolsNode 是编排中的执行节点，它根据 assistant 消息里的工具名称和参数调用已注册的工具，再生成工具结果消息，两者职责不同。

<<< @/go/codes/eino/iTool.go

`BaseTool.Info` 提供名称、描述和参数定义。`InvokableTool.InvokableRun` 接收 JSON 参数字符串并返回完整结果，`StreamableTool.StreamableRun` 则返回字符串流。这里展示的是文本结果接口，该版本还提供支持结构化多模态结果的增强接口。

### 9.2 工具描述

<<< @/go/codes/eino/sToolInfo.go

参数可通过 `NewParamsOneOfByParams` 或 `NewParamsOneOfByJSONSchema` 描述。参数声明用于指导模型生成调用，不应代替工具自身的输入校验。

### 9.3 HTTP GET 工具

<<< @/go/codes/eino/eTool.go

示例直接调用工具获取本站的 `sitemap.xml`，不经过模型。工具使用 HTTP 客户端超时，调用本身也带上下文超时。

### 9.4 创建本地工具

<<< @/go/codes/eino/eNewTool.go

`utils.NewTool` 根据显式的 `ToolInfo` 封装 Go 函数，负责 JSON 参数解码和结果转换。`GetNote` 检查空名称和未知名称，正常运行时输出 P4 笔记链接。也可以用 `utils.InferTool` 从函数参数类型推导参数定义。

### 9.5 通过 ToolsNode 执行

下面复用上一节的 `InputParams`、`GetNote` 和 `CreateTool`。将辅助函数与这些定义放在同一个包，补充 `compose` 导入，在已有 `main` 中调用 `executeToolCall` 并处理返回值：

<<< @/go/codes/eino/eToolsNode.go

这里手工构造一条工具调用消息，因此无需模型服务。返回的工具消息通过 `ToolCallID` 关联到 `call_1`。实际对话中，还要把 assistant 的工具调用消息和 tool 结果追加到消息历史，再交给模型继续回答。ToolsNode 本身不会完成这一整轮对话。

## 10. 编排

### 10.1 Chain

Chain 适合顺序执行的流程。编译时检查节点连接和类型是否匹配，得到 `Runnable` 后再调用 `Invoke` 或 `Stream`。本例由模板生成消息列表，再调用模型：

<img width="550" src="/eino/chain_simple_llm.png" alt="Chain：ChatTemplate 生成消息后交给 ChatModel" />

<<< @/go/codes/eino/eChain.go

输入是 `map[string]any`，模板输出是 `[]*schema.Message`，模型输出是 `*schema.Message`。模板允许缺省历史，所以此处只传入角色和主题。

### 10.2 Graph 与分支

Graph 通过节点、边和分支描述执行路径，可以支持循环，也可以配置为 DAG。`AddBranch` 的条件函数根据上游输出选择目标，可能到达的目标需要在分支定义中声明。

本例的执行路径为：

```text
START → classify → cat   → END
                 → dog   → END
                 → other → END
```

<<< @/go/codes/eino/eGraph.go

输入 `1`、`2` 和其他值分别选择三个分支，示例依次输出 `喵喵喵`、`汪汪汪` 和 `你好`。分支由普通 Go 逻辑决定，不需要模型判断。

### 10.3 Graph 中的模型节点

<<< @/go/codes/eino/eGraphWithModel.go

输入包含 `style` 和 `content`。分支选择简要或详细的系统提示，两条路径都把用户问题保留为 `user` 消息，再交给同一个模型节点。节点之间的输入输出类型必须匹配，添加节点、边和编译时的错误都需要处理。

### 10.4 调用内状态

`WithGenLocalState` 为每次图调用创建状态。它在本次调用的节点间共享，不会自动变成跨请求的会话历史或持久化数据。状态生成函数应创建新对象，重复返回同一个共享指针会破坏请求隔离。

`StatePreHandler` 和 `StatePostHandler` 在节点前后读取或修改状态，也可以改变节点的输入输出。流式场景有对应的 `StreamStatePreHandler` 和 `StreamStatePostHandler`。在流上使用非流式处理器可能触发聚合，不能据此保证增量输出不受影响。

节点内部访问状态可以使用 `compose.ProcessState`，公开签名如下：

<<< @/go/codes/eino/fProcessState.go

框架在状态处理器和 `ProcessState` 的回调执行期间加锁。不要把状态指针带到回调外直接修改，也不要在持有同一状态锁的回调里再次调用 `ProcessState`，否则可能产生数据竞争或死锁。流式处理器返回后，另起 goroutine 的访问也不自动受这把锁保护。

<<< @/go/codes/eino/eGraphWithState.go

程序连续调用同一个 `Runnable` 两次，每次都输出 `request: count=1`。`count` 节点通过 `ProcessState` 更新状态，`format` 节点的前处理器读取它，展示的是调用内共享和调用间隔离。

### 10.5 回调

回调用于记录组件和编排执行过程。`RunInfo` 描述触发回调的实体，`Name` 是展示名称，`Type` 是实现类型，`Component` 是组件类别：

<<< @/go/codes/eino/sRunInfo.go

回调输入输出按组件区分，下面的类型定义不代表它们具有统一的结构：

<<< @/go/codes/eino/tCallbackInOut.go

<<< @/go/codes/eino/eGraphWithCallback.go

示例记录开始、结束和错误，并通过 `compose.WithCallbacks` 应用于本次调用。需要读取特定组件的数据时，应使用该组件的转换函数，例如 `model.ConvCallbackInput`，并检查转换结果。

回调不应修改共享的输入输出对象。多个处理器之间没有可依赖的执行顺序，并发节点的日志也可能交错。流式回调拿到的是流副本，使用后同样需要关闭，不能把副本遗留到调用结束以后。

### 10.6 嵌套图

`AddGraphNode` 可以把图或 Chain 直接加入外层图，外层 `Compile` 会处理内部编排，不要求先把内层图编译后再包装为 Lambda。嵌套前后仍需保持输入输出类型匹配。

<<< @/go/codes/eino/eGraphWithGraph.go

内层图去除首尾空白并转换为大写，外层图追加前缀，最终输出 `result: EINO`。

## 11. CozeLoop

CozeLoop 回调扩展把 Eino 的执行信息转换为 trace。运行示例前，还需要设置 `COZELOOP_WORKSPACE_ID` 和 `COZELOOP_API_TOKEN`，对应账号中的工作空间与访问凭据。

<<< @/go/codes/eino/eCozeLoop.go

`AppendGlobalHandlers` 影响整个进程，应在初始化阶段注册一次，不要在每次请求时重复追加，也不要与正在执行的调用并发修改全局回调列表。只想给某次调用添加追踪时，可以使用 `compose.WithCallbacks(handler)`。

客户端关闭时需要留出上报时间。示例另建关闭上下文，防止使用已经取消的业务上下文。程序得到 `EINO` 只表示本地编排成功，trace 是否出现在平台上还取决于凭据、网络和服务端接收情况。
