---
outline: [2, 3]
---

# Eino

## ChatModel 组件

ChatModel 是 Eino 框架中对对话大模型的抽象，它提供了统一的接口来与不同的大模型服务（如 OpenAI、Ollama 等）进行交互。

Model 组件是一个用于与大语言模型交互的组件。它的主要作用是将用户的输入消息发送给语言模型，并获取模型的响应。适合用于以下场景：

- 自然语言对话
- 文本生成
- 生成多模态交互（文本、图片等）
- 工具调用的参数

### 接口定义与三种方法

<<< @/go/codes/eino/iBaseChatModel.go

#### `Generate` 方法

- 功能：生成完整的模型响应
- 参数：
  - `ctx`：上下文对象，用于传递请求级别的信息，同时也用于传递 Callback Manager
  - `input`：输入消息列表
  - `opts`：可选参数，用于配置模型行为
- 返回值：
  - `*schema.Message`：模型生成的响应消息
  - `error`：生成过程中的错误信息

#### `Stream` 方法

- 功能：以流式方式生成模型响应
- 参数：与 `Generate` 方法相同
- 返回值：
  - `*schema.StreamReader[*schema.Message]`：模型响应的流式读取器
  - `error`：生成过程中的错误信息

#### `WithTools` 方法

- 功能：为模型绑定可用的工具
- 参数：
  - `tools`：工具信息列表
- 返回值：
  - `ToolCallingChatModel`: 绑定了 tools 后的 chatmodel
  - `error`：绑定过程中的错误信息

### `schema.Message` 结构体

<<< @/go/codes/eino/sMessage.go

### `Generate` 完整生成

<<< @/go/codes/eino/eGenerate.go

### `Stream` 流式生成

<<< @/go/codes/eino/eStream.go

## ChatTemplate 组件

### 接口定义

<<< @/go/codes/eino/iChatTemplate.go

- `prompt.FromMessages()`
  - 用于把多个 message 变成一个 chat template。
- `schema.Message{}`
  - `schema.Message` 是实现了 `Format` 接口的结构体，因此可直接构建 `schema.Message{}` 作为 template
- `schema.SystemMessage()`
  - 此方法是构建 role 为“system”的 message 快捷方法
- `schema.AssistantMessage()`
  - 此方法是构建 role 为“assistant”的 message 快捷方法
- `schema.UserMessage()`
  - 此方法是构建 role 为“user”的 message 快捷方法
- `schema.ToolMessage()`
  - 此方法是构建 role 为“tool”的 message 快捷方法
- `schema.MessagesPlaceholder()`
  - 可用于把一个 `[]*schema.Message` 插入到 message 列表中，常用于插入历史对话

#### `Format` 方法

- 功能：将变量值填充到消息模板中
- 参数：
  - `ctx`：上下文对象，用于传递请求级别的信息，同时也用于传递 Callback Manager
  - `vs`：变量值映射，用于填充模板中的占位符
  - `opts`：可选参数，用于配置格式化行为
- 返回值：
  - `[]*schema.Message`：格式化后的消息列表
  - `error`：格式化过程中的错误信息

### 基本使用

1. 使用 `prompt.FromMessages` 以及提供的 `schema.FString` 格式化
2. 构造 `params` 参数
3. 使用 `template` 自带的 `Format` 构建即可

### 示例代码

<<< @/go/codes/eino/eChatTemplate.go

## RAG 概念

RAG 的核心思想是在 LLM 生成回复之前，先去外部数据源中查找相关信息，然后将这些信息作为额外上下文提供给 LLM，引导其生成更准确的回答。

### 为什么需要 RAG

1. 减少大模型幻觉现象，提高回答的准确性和相关度
2. 构建专门领域的回答资料库

### RAG 三步走

1. 数据准备（索引）
   - 收集和处理外部数据（提前搜集）：将各种形式的外部数据（如文档、数据库、网页、API 等）收集起来。这些数据可以是企业的内部知识库，也可以是互联网上的公开信息。
   - 切分和嵌入（Transformer）：将这些数据切分成小块（称为“chunk”），并使用嵌入模型将每个小块转换为高维度的数字向量（称为“嵌入”）。这些嵌入能够捕捉文本的语义信息。
   - 存储到向量数据库（Embedding&Indexer）：将这些向量存储在专门的向量数据库中，以便后续进行高效的相似性搜索。
2. 检索（Reteriever）
   - 用户查询：当用户提出问题或输入查询时，RAG 系统会首先将这个查询也转换为一个向量。
   - 相似性搜索：然后，系统会在向量数据库中进行相似性搜索，找到与用户查询最相关的文档块（或片段）。这个过程就像在海量信息中找到最相关的“证据”。
3. 增强生成
   - 构建增强提示：将用户原始的查询和从向量数据库中检索到的相关信息（作为上下文）一起发送给大型语言模型。
   - LLM 生成回答：LLM 在接收到这个“增强”后的提示后，会根据其自身的训练知识和提供的外部上下文来生成最终的回答。这样，LLM 不再仅仅依赖于其训练时的数据，而是有了“实时”的参考资料。

## Embedding 组件

Embedding 组件是一个用于将文本转换为向量表示的组件。它的主要作用是将文本内容映射到向量空间，使得语义相似的文本在向量空间中的距离较近。常用的场景有：

- 文本相似度计算
- 文本聚类分析
- 语义搜索

### 接口定义

<<< @/go/codes/eino/iEmbedding.go

#### `EmbedStrings` 方法

- 功能：将一组文本转换为向量表示
- 参数：
  - `ctx`：上下文对象，用于传递请求级别的信息，同时也用于传递 Callback Manager
  - `texts`：待转换的文本列表
  - `opts`：转换选项，用于配置转换行为
- 返回值：
  - `[][]float64`：文本对应的向量表示列表，每个向量的维度由具体的实现决定
  - `error`：转换过程中的错误信息

### 示例代码

<<< @/go/codes/eino/eEmbedding.go

## Indexer 组件

Indexer 组件是一个用于存储和索引文档的组件。它的主要作用是将文档及其向量表示存储到后端存储系统中，并提供高效的检索能力。常用的场景有：

- 构建向量数据库，以用于语义关联搜索

### 接口定义

<<< @/go/codes/eino/iIndexer.go

#### `Store` 方法

- 功能：存储文档并建立索引
- 参数：
  - `ctx`：上下文对象，用于传递请求级别的信息，同时也用于传递 Callback Manager
  - `docs`：待存储的文档列表
  - `opts`：存储选项，用于配置存储行为
- 返回值：
  - `ids`：存储成功的文档 ID 列表
  - `error`：存储过程中的错误信息

### 示例代码

<<< @/go/codes/eino/eIndexer.go

## Retriever 组件

Retriever 组件是一个用于从各种数据源检索文档的组件。它的主要作用是根据用户的查询（query）从文档库中检索出最相关的文档。常用的场景有：

- 基于向量相似度的文档检索
- 基于关键词的文档搜索
- 知识库问答系统（rag）

### 接口定义

<<< @/go/codes/eino/iRetriever.go

#### `Retrieve` 方法

- 功能：根据查询检索相关文档
- 参数：
  - `ctx`：上下文对象，用于传递请求级别的信息，同时也用于传递 Callback Manager
  - `query`：查询字符串
  - `opts`：检索选项，用于配置检索行为
- 返回值：
  - `[]*schema.Document`：检索到的文档列表
  - `error`：检索过程中的错误信息

### `Option` 结构体

<<< @/go/codes/eino/sRetrieverOptions.go

### 示例代码

<<< @/go/codes/eino/eRetriever.go

## Document Transformer 组件

Document Transformer 是一个用于文档转换和处理的组件。它的主要作用是对输入的文档进行各种转换操作，如分割、过滤、合并等，从而得到满足特定需求的文档。常用的场景有：

- 将长文档分割成小段落以便于处理
- 提取文档中的特定部分
- 根据特定规则过滤文档内容
- 对文档内容进行结构化转换

### 接口定义

<<< @/go/codes/eino/iTransformer.go

### `Transform` 方法

- 功能：对输入的文档进行转换处理
- 参数：
  - `ctx`：上下文对象，用于传递请求级别的信息，同时也用于传递 Callback Manager
  - `src`：待处理的文档列表
  - `opts`：可选参数，用于配置转换行为
- 返回值：
  - `[]*schema.Document`：转换后的文档列表
  - `error`：转换过程中的错误信息

### 示例代码

<<< @/go/codes/eino/eTransformer.go

## RAG Demo

[https://github.com/zhh2001/RAG_Learning](https://github.com/zhh2001/RAG_Learning)

## Tool 组件

ToolsNode 组件是一个用于扩展模型能力的组件，它允许模型调用外部工具来完成特定的任务。这个组件可用于以下场景中：

- 让模型能够获取实时信息（如搜索引擎、天气查询等）
- 使模型能够执行特定的操作（如数据库操作、API 调用等）
- 扩展模型的能力范围（如数学计算、代码执行等）
- 与外部系统集成（如知识库查询、插件系统等）

### 接口定义

<<< @/go/codes/eino/iTool.go

#### `Info` 方法

- 功能：获取工具的描述信息
- 参数：
  - `ctx`：上下文对象
- 返回值：
  - `*schema.ToolInfo`：工具的描述信息，用于提供给大模型
  - `error`：获取信息过程中的错误

#### `InvokableRun` 方法

- 功能：同步执行工具
- 参数：
  - `ctx`：上下文对象，用于传递请求级别的信息，同时也用于传递 Callback Manager
  - `argumentsInJSON`：JSON 格式的参数字符串
  - `opts`：工具执行的选项
- 返回值：
  - `string`：执行结果
  - `error`：执行过程中的错误

#### `StreamableRun` 方法

- 功能：以流式方式执行工具
- 参数：
  - `ctx`：上下文对象，用于传递请求级别的信息，同时也用于传递 Callback Manager
  - `argumentsInJSON`：JSON 格式的参数字符串
  - `opts`：工具执行的选项
- 返回值：
  - `*schema.StreamReader[string]`：流式执行结果
  - `error`：执行过程中的错误

### `ToolInfo` 结构体

<<< @/go/codes/eino/sToolInfo.go

### 示例代码

<<< @/go/codes/eino/eTool.go

创建一个 tool：

<<< @/go/codes/eino/eNewTool.go

## 编排

### Chain

<img width="550 " src="/eino/chain_simple_llm.png" alt="chain" />

<<< @/go/codes/eino/eChain.go

### Graph

<img width="680 " src="/eino/eino_take_first_toolcall_output.png" alt="graph" />

<<< @/go/codes/eino/eGraph.go

### GraphWithModel

<<< @/go/codes/eino/eGraphWithModel.go

### GraphWithState

#### State 定义

Graph 可以有自身的全局状态，在创建 Graph 时传入 `compose.WithGenLocalState` Option 开启此功能。这个请求维度的全局状态在一次请求的各环节可读写使用。

Eino 推荐用 `StatePreHandler` 和 `StatePostHandler`，功能定位是：

- `StatePreHandler`：在每个节点执行前读写 State，以及按需替换节点的 Input。输入需对齐节点的非流式输入类型。
- `StatePostHandler`：在每个节点执行后读写 State，以及按需替换节点的 Output。输入需对齐节点的非流式输出类型。

针对流式场景，使用对应的 `StreamStatePreHandler` 和 `StreamStatePostHandler`，输入需分别对齐节点的流式输入和流式输出类型。

这些 state handlers 位于节点外部，通过对 Input 或 Output 的修改影响节点，从而保证了节点的**状态无关**特性。

如果需要在节点内部读写 State，Eino 提供了对应函数：

<<< @/go/codes/eino/fProcessState.go

Eino 框架会在所有读写 State 的位置加锁。

#### 示例代码

<<< @/go/codes/eino/eGraphWithState.go{36,41-46,81-84,104}

### GraphWithCallback

核心概念串起来，就是：Eino 中的 Component 和 Graph 等**实体**，在固定的**时机** (Callback Timing)，回调用户提供的 **function** (Callback Handler)，并把**自己是谁** (`RunInfo`)，以及**当时发生了什么** (Callback Input & Output) 传出去。

#### `RunInfo` 结构体

<<< @/go/codes/eino/sRunInfo.go

#### `CallbackInput` & `CallbackOutput` 类型

本质是任意类型，因为不同的 Component 的输入输出、内部状态完全不同。

<<< @/go/codes/eino/tCallbackInOut.go

#### 示例代码

<<< @/go/codes/eino/eGraphWithCallback.go

### GraphWithGraph（图嵌套）

图编排产物 `Runnable` 与 Lambda 的接口形式非常相似。因此编译好的图可以简单的封装为 Lambda，并以 Lambda 节点的形式嵌套进其他图中。

其实就是将图当作一个特殊的 Lambda 节点嵌入。

注意依旧要保持**上下游类型对齐**。

#### 示例代码

<<< @/go/codes/eino/eGraphWithGraph.go

## CozeLoop

字节官方为 CozeLoop 实现了 Trace 回调，可以与 Eino 的应用无缝集成增强可观测能力。

<<< @/go/codes/eino/eCozeLoop.go
