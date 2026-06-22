---
outline: deep
---

# P4Runtime

P4 能描述交换机“会怎样处理包”，但把程序编译进设备后，表里仍然是空的。谁来写路由、更新下一跳、读取计数器，主控制器断开后又由谁接手？这些运行期间的问题，就是 P4Runtime 要解决的事。

这份笔记依据 [P4Runtime v1.5.0 规范](https://p4lang.github.io/p4runtime/spec/v1.5.0/P4Runtime-Spec.html) 和同版本的 [Protobuf 定义](https://github.com/p4lang/p4runtime/tree/v1.5.0/proto) 整理。它不是逐段翻译，而是按实际学习顺序重新组织：先弄清 P4Runtime 在系统中的位置，再认识 P4Info、实体和 RPC，最后用 Python 与 BMv2 走通一次完整交互。

::: info 规范用语
规范里的 **MUST / MUST NOT** 表示实现必须遵守，**SHOULD / SHOULD NOT** 表示通常应当遵守，只有充分理由才能偏离，**MAY** 表示可选能力。本文出现“必须”“应该”“可以”时也沿用这一强度；类比和经验性说明会明确写成“可以理解为”或“实践中”。
:::

## 1 先建立直觉：程序、说明书和遥控器

把 P4 交换机想成一台可更换内部机构的自动分拣机：

- **P4 程序**决定机器里有哪些传送带、扫码规则和可填写的表格；
- **设备配置**（device config）是编译器交给具体设备的机器指令；
- **P4Info**是控制器能看懂的操作说明书，记录表、动作、参数和计数器的名字、ID 与类型；
- **P4Runtime**是统一的远程操作协议，负责换配置、填表、读状态和收发控制报文；
- **控制器**才是做决定的人，例如计算路由后把结果写进表里。

因此，P4Runtime 并不是另一门数据平面语言，也不负责计算最短路径。它解决的是：**控制器如何用一套与具体 P4 程序配套、与设备厂商尽量解耦的接口管理 P4 数据平面。**

```text
             P4 源程序
                 │ 编译
       ┌─────────┴──────────┐
       ▼                    ▼
    P4Info             device config
  控制面说明书          设备专用配置
       │                    │
       └─────────┬──────────┘
                 ▼
控制器 ── gRPC / P4Runtime ──► P4Runtime Server ──► 数据平面
```

P4Runtime 定义的是控制器与服务器之间的行为。服务器可能直接运行在交换机上，也可能是机框管理进程或代理；规范并不强制它的内部实现。

## 2 P4Runtime 到底提供了哪些接口？

P4Runtime 使用 Protocol Buffers 描述消息，默认通过 gRPC 传输。`P4Runtime` 服务只有六组核心 RPC：

| RPC | 方向 | 作用 |
| --- | --- | --- |
| `Write` | 一元 | 插入、修改或删除表项等 P4 实体 |
| `Read` | 服务端流 | 按查询条件读取实体；结果可以拆成多个响应返回 |
| `SetForwardingPipelineConfig` | 一元 | 验证、保存或启用新的转发流水线配置 |
| `GetForwardingPipelineConfig` | 一元 | 读取当前 P4Info、设备配置或 cookie |
| `StreamChannel` | 双向流 | 仲裁、Packet I/O、Digest、空闲超时通知和流错误 |
| `Capabilities` | 一元 | 查询服务器支持的 P4Runtime API 版本 |

这里的“一元”指一问一答。  
“服务端流”指一个请求可能对应多个响应。  
“双向流”则允许双方在一条长期连接上独立发送消息。

每台设备由一个 `uint64 device_id` 标识。这个值如何分配属于部署系统的职责，P4Runtime 不提供发现设备或分配 ID 的机制。控制器必须先通过带外配置知道服务器地址、凭据与 `device_id`。

::: warning P4Runtime 不包办所有设备管理
端口创建、机框管理、证书发放、设备发现、P4 编译和拓扑计算都不属于核心 P4Runtime API。厂商可以用 gNMI、gNOI 或自己的管理接口补齐这些能力。
:::

## 3 一次理想的启动流程

从控制器启动到真正转发流量，可以拆成六步：

1. 控制器建立 gRPC 通道并打开 `StreamChannel`。
2. 它首先发送 `MasterArbitrationUpdate`，声明设备、角色和选举 ID。
3. 服务器返回仲裁结果；只有成为该角色的主控制器后，客户端才可执行受主身份保护的写操作。
4. 控制器调用 `Capabilities`，确认 API 版本是否兼容。
5. 控制器用 `SetForwardingPipelineConfig` 下发 P4Info 和设备配置。
6. 控制器依据 P4Info 中的 ID 构造 `Write`，写入表项、动作选择器成员、计数器配置等实体。

```text
Controller                         P4Runtime Server
    │──── StreamChannel ──────────────────►│
    │──── Arbitration(device, role, id) ──►│
    │◄─── Arbitration(status) ─────────────│
    │──── Capabilities ───────────────────►│
    │◄─── api_version ─────────────────────│
    │──── SetPipeline(P4Info + config) ───►│
    │──── Write(TableEntry...) ───────────►│
    │◄════════ PacketIn / Digest ═════════►│
```

这只是最顺的一条路径。实践中控制器可能只读取已有设备，也可能在不改变 P4 程序的情况下重新连接。服务器或控制器重启后，双方都不能凭内存猜测状态：控制器应重新仲裁，并按系统设计读取、核对或恢复所需实体。

规范区分两类状态：

- **P4Runtime 管理的状态**：流水线配置和通过实体 API 暴露的状态；
- **数据平面易失状态（data-plane volatile state）**：例如计数值、计量器运行状态等，重启或重新配置后不一定保留。

是否持久化、何时清空还与目标架构和设备实现有关，控制器不能把所有状态都当成数据库里的永久记录。

## 4 控制器为什么需要 P4Runtime？

规范列出的用例可以归纳为四类：

- **集中式控制**：控制器计算路由、ACL 或负载均衡策略，再批量写入设备。
- **本地代理控制**：设备上的代理把上层意图转换成 P4Runtime 实体。
- **多控制器高可用**：同一角色有主备控制器，主节点故障后通过仲裁切换。
- **混合管理**：不同角色各管一部分实体，例如路由控制器管理三层表，安全控制器管理 ACL 表。

P4Runtime 的“协议无关”不是说控制器完全不知道 P4 程序，而是说 RPC 和消息外壳不随程序改变。某张表有哪些键和动作仍由 P4Info 决定。换句话说，**API 形状稳定，业务对象由 P4 程序描述。**

## 5 主备仲裁、角色与选举 ID

### 5.1 为什么写操作必须先仲裁

如果两个控制器同时认为自己是主节点，一个写“端口 1”，另一个写“端口 2”，交换机状态会不断抖动。P4Runtime 用 `StreamChannel` 上的 `MasterArbitrationUpdate` 建立写权限归属。

仲裁消息包含：

| 字段 | 含义 |
| --- | --- |
| `device_id` | 要控制的设备 |
| `role` | 控制器角色；默认空角色代表默认管理域 |
| `election_id` | 128 位无符号选举 ID，通常随新一轮选主单调增大 |
| `status` | 服务器响应中的仲裁结果 |

同一设备、同一角色下，服务器根据选举 ID 确定主客户端。控制器不应重复使用与其他活动客户端相同的选举 ID。成为主节点后，`Write` 和设置流水线配置等请求还要携带匹配的 `device_id`、角色与 `election_id`，服务器据此拒绝过期或无权请求。

`role_id` 在 v1.4.0 已弃用，v1.5.0 使用字符串 `role`。角色的具体权限通过 `Role.config` 携带的 `google.protobuf.Any` 表示，内容取决于设备实现。角色可以限制可写实体、PacketIn 接收范围等；P4Runtime 核心规范不规定一种通用权限模型。

### 5.2 StreamChannel 的首条消息

客户端为设备打开流后，首条请求必须是仲裁更新。服务器通过它把流与设备和角色关联起来。连接中断就意味着该客户端不再通过这条流存活，备节点应参与新的仲裁，而不是继续沿用旧连接的主身份。

只读操作和主身份之间的关系比“备节点什么都不能做”更细：服务器可以允许备控制器读取，但写请求必须按 RPC 规则校验主身份及角色权限。实际系统仍应避免让大量备节点无节制轮询设备。

::: tip 选举 ID 不是设备版本号
它只比较同一角色下控制器主身份的新旧，不能替代流水线 cookie、配置版本或业务事务号。
:::

## 6 P4Info：连接 P4 程序与控制器的说明书

P4Info 是 `p4.config.v1.P4Info` 消息。编译器从 P4 程序生成它，控制器通常不应手写。它主要包含三类信息：

1. **可控对象**：表、动作、动作配置文件、计数器、计量器、Digest、ValueSet、寄存器和架构相关外部对象；
2. **对象元数据**：数值 ID、完整名称、别名、注解和文档；
3. **类型信息**：结构体、头部、枚举、新类型以及序列化所需的位宽。

### 6.1 Preamble 与对象 ID

大多数顶层对象都有 `Preamble`：

- `id` 是控制面使用的 32 位 ID；
- `name` 是 P4 作用域下的完整名称，如 `MyIngress.ipv4_lpm`；
- `alias` 是能唯一定位对象的最短名称后缀；
- `annotations`、`structured_annotations` 和 `doc` 保留程序中的说明。

对象 ID 的高 8 位表示对象类型，低 24 位标识具体对象。编译器负责稳定生成；若程序用 `@id` 指定 ID，也必须满足对应规则。控制器应从当前 P4Info 查 ID，不能把某次编译得到的数字散落在源码里。

### 6.2 表和动作如何关联

一张表的 P4Info 会列出：

- 匹配字段的 ID、名称、位宽和匹配类型；
- 可用动作的引用；
- 表容量、常量属性、默认动作约束；
- 是否关联直接计数器、直接计量器、空闲超时等资源。

动作则列出参数 ID、名称和位宽。控制器先通过表 ID 指定目标，再用字段 ID 填匹配键，用动作 ID 和参数 ID 填执行内容。

```text
P4: table ipv4_lpm                  P4Info
    key dstAddr: lpm      ───────►  table_id + field_id + bitwidth
    action ipv4_forward   ───────►  action_id + param_id + bitwidth
```

### 6.3 PkgInfo 与 P4TypeInfo

`PkgInfo` 可以记录包名、版本、架构、组织、联系方式和编译信息，便于控制器判断拿到的是哪一份程序。`P4TypeInfo` 描述复杂 P4 类型，使 `P4Data` 能表达结构体、元组、头部、头部栈、联合、枚举和错误类型。

自定义类型如果带 `@p4runtime_translation`，可以声明控制平面看到的 URI 与底层位宽。控制器必须按 P4Info 暴露的控制面类型编码，而不是只看 P4 源码里的表面名称。

## 7 转发流水线配置

`ForwardingPipelineConfig` 包含：

- `p4info`：控制面接口描述；
- `p4_device_config`：不透明的设备专用二进制配置；
- `cookie`：控制器指定的 64 位标识，用于识别配置。

服务器不解释 cookie 的业务含义，控制器可以用它存构建号或配置哈希的截断值。它也不能替代对实际 P4Info 的兼容性检查。

### 7.1 五种设置动作

| Action | 行为 |
| --- | --- |
| `VERIFY` | 只验证配置，不保存也不启用 |
| `VERIFY_AND_SAVE` | 验证并保存，暂不启用 |
| `VERIFY_AND_COMMIT` | 验证后立即启用，是最常见的首次下发方式 |
| `COMMIT` | 启用先前保存的配置 |
| `RECONCILE_AND_COMMIT` | 尝试在保留转发状态的同时切换配置，属于可选能力 |

请求里的 `UNSPECIFIED` 不是有效操作。目标不支持某个可选动作时必须返回错误，而不能悄悄换成别的语义。

流水线切换可能影响现有表项。`VERIFY_AND_COMMIT` 不承诺保留旧状态；`RECONCILE_AND_COMMIT` 的目的才是协调旧状态与新 P4Info，但设备是否支持以及哪些改变可协调都有限制。生产系统应把“不兼容升级如何回滚”作为显式流程，而不是寄希望于设备自动猜对。

### 7.2 读取配置时少拿一点

`GetForwardingPipelineConfig` 的 `response_type` 可以选择：

- `ALL`：返回全部内容；
- `COOKIE_ONLY`：只取 cookie；
- `P4INFO_AND_COOKIE`：不取较大的设备配置；
- `DEVICE_CONFIG_AND_COOKIE`：不取 P4Info。

设备配置可能很大，日常一致性探测优先读 cookie 或 P4Info 与 cookie，不必每次拉回完整二进制。

## 8 消息编码：最容易“看起来对、实际错”的地方

### 8.1 默认值与读写对称

Proto3 不区分“未设置的标量字段”和“显式设置为零”的所有情况，而 P4Runtime 又大量使用零作为通配符或默认语义。构造请求时必须以对应实体章节的规则为准，不能笼统地认为零等于缺省。

所谓 **读写对称（read-write symmetry）**，可以理解为：服务器读出的实体应能在不丢失必要语义的情况下被客户端用于后续写入。只读字段、通配查询字段和某些默认值仍要按规范剔除或处理。

### 8.2 Bytestring

P4 的 `bit<W>`、`int<W>` 没有固定最大位宽，Protobuf 的整数类型不够通用，因此许多 P4 值用 `bytes` 表示。规则要点是：

- 使用大端序；
- 位宽从 P4Info 获得；
- 无符号值的扩展位必须为 0，带符号值必须正确符号扩展；
- 数值必须落在声明位宽范围内；
- 规范形式使用能容纳数值的最短字节串，零至少使用一个字节。

例如 `bit<9>` 的十进制 2 可规范编码为 `02`，510 编码为 `01 fe`。把字符串 `"510"` 的 UTF-8 字节发过去当然不是同一回事。

::: warning 位宽不是字节数
Python 的 `int.to_bytes()` 接收的是字节数。`bit<9>` 最多需要 2 字节，而不是 9 字节。本文示例的 `encode()` 同时检查范围并生成规范的最短表示。
:::

### 8.3 P4Data

`P4Data` 是复杂 P4 值的通用容器，通过 `oneof` 表示位串、整数、布尔值、元组、结构体、头部、头部栈、联合和枚举等。解码时要同时查看 `P4TypeInfo`；只看 Protobuf 外形不足以判断一个字节串在 P4 中的真实类型。

## 9 Entity：P4Runtime 真正管理的对象

`Entity` 是一个 `oneof` 容器，一次只承载一种实体。v1.5.0 的核心实体包括表项、动作配置文件成员与组、计数器、计量器、包复制引擎、ValueSet、寄存器、Digest 和外部对象。

### 9.1 TableEntry

`TableEntry` 是最常用的实体，关键字段有：

| 字段 | 作用 |
| --- | --- |
| `table_id` | P4Info 中的表 ID |
| `match` | 一组字段匹配，字段顺序不应影响语义 |
| `action` | 直接动作、成员 ID、组 ID 或一次性动作集合 |
| `priority` | 需要优先级的表项使用正数；数值越大优先级越高 |
| `is_default_action` | 操作默认动作而不是普通表项 |
| `idle_timeout_ns` | 请求空闲超时通知的时间 |
| `metadata` | 控制器自用的透明字节串 |

匹配类型包括：

- `exact`：精确相等；
- `lpm`：值加前缀长度；
- `ternary`：值加掩码；
- `range`：闭区间下界与上界；
- `optional`：字段存在时精确匹配，不存在时通配；
- `other`：为扩展匹配类型预留的 `Any`。

LPM 的值在前缀以外必须为零；三态匹配的值在掩码为零的位置也必须归零。这种规范化保证同一个逻辑键只有一种表示。对允许通配的字段，省略该字段表达“不关心”，不要塞一个全零掩码冒充通配。

包含 ternary、range 或 optional 等需要优先级语义的普通表项通常必须给出大于零的 `priority`；只含 exact 和 LPM 的表项使用零。默认动作通过 `is_default_action=true` 指定，不携带普通匹配键和优先级，并按表的约束修改。

`INSERT` 要求键不存在，`MODIFY` 要求目标存在，`DELETE` 只需提供能唯一定位实体的键。删除时把旧响应里所有只读字段和动作内容原样带回去，可能导致设备拒绝，所以示例会专门构造 key-only 表项。

### 9.2 ActionProfile：成员、组与一次性选择

带选择器的动作配置文件有两种常见编程方式：

1. **手工成员/组模式**：先写 `ActionProfileMember`，再写引用成员的 `ActionProfileGroup`，表项引用成员或组 ID；
2. **one-shot 模式**：表项直接携带 `ActionProfileActionSet`，由服务器管理内部成员。

一个动作集合可包含动作、权重、监视端口和选择模式。v1.5.0 还描述了 `SUM_OF_WEIGHTS` 与 `SUM_OF_MEMBERS` 等容量语义。控制器不能在同一作用域里随意混用两种模式，也要遵守表与 action profile 在 P4Info 中的绑定关系。

组依赖成员时，控制器应先创建成员再创建组，删除时反过来。一个 `WriteRequest` 内的更新可能被服务器为性能重新排序，不能把有严格依赖的复杂迁移想当然地当作逐行脚本。

### 9.3 Counter 与 Meter

计数器分为：

- `CounterEntry`：独立计数器，通过 `counter_id + index` 定位；
- `DirectCounterEntry`：直接绑定表项，通过表项键定位。

`CounterData` 包含包数和字节数。读取时省略索引可以表示通配查询；写计数数据是否支持、写入后精确效果如何要遵守目标架构。计数器是典型的运行状态，不能假设流水线重载后一定保留。

计量器同样分 `MeterEntry` 和 `DirectMeterEntry`。`MeterConfig` 包含承诺速率/突发量 `cir/cburst`、峰值速率/突发量 `pir/pburst`，v1.5.0 还包含 `eburst`。单位由 P4Info 中的 meter spec 决定，可能是包或字节；控制器必须先看单位再填数字。

### 9.4 Packet Replication Engine

包复制引擎（PRE）实体包括：

- `MulticastGroupEntry`：组 ID 加多个副本，每个副本指定端口和实例；
- `CloneSessionEntry`：克隆会话、副本、服务等级和可选的截断长度。

v1.5.0 的 `Replica` 使用 `port` bytestring，旧的 `egress_port` 字段已经弃用。副本还可带备用副本列表。组 ID、会话 ID 的合法范围及端口编码受目标约束。

### 9.5 ValueSet、Register、Digest 与 Extern

- **ValueSetEntry**：在解析器 ValueSet 中增加或删除成员，成员由一组 `FieldMatch` 表示；支持程度受目标架构约束。
- **RegisterEntry**：用 `register_id + index` 访问寄存器单元，数据放在 `P4Data`；省略索引可用于读取整个数组。
- **DigestEntry**：配置 Digest 的批处理与确认行为，包括最大等待时间、最大列表长度和确认超时；真正的数据通过 `StreamChannel` 上报。
- **ExternEntry**：用外部类型 ID、实例 ID 和 `google.protobuf.Any` 承载非核心架构对象。

复杂实体最重要的共同规律是：**P4Info 告诉你对象是什么，Entity 告诉你这一次要操作哪个实例。**

## 10 Write RPC：批量修改与原子性

`WriteRequest` 包含设备、角色、选举 ID、一组 `Update` 和原子性选项。每个 `Update` 的类型是：

- `INSERT`：创建不存在的实体；
- `MODIFY`：修改已经存在的实体；
- `DELETE`：删除实体。

`UNSPECIFIED` 不是可执行更新。服务器必须校验客户端是否为对应角色的主节点，并逐项校验 ID、字段、位宽、引用和权限。

### 10.1 三种原子性

| Atomicity | 含义 |
| --- | --- |
| `CONTINUE_ON_ERROR` | 必选默认能力；某项失败仍尝试其余更新，数据包可能看到中间状态 |
| `ROLLBACK_ON_ERROR` | 可选；遇错回滚控制面和数据面状态，但处理期间数据包仍可能看到中间状态 |
| `DATAPLANE_ATOMIC` | 可选；数据包只能看到批处理前或完成后的状态 |

后两种不是所有设备都支持。请求可选能力前应了解目标，否则服务器会返回不支持错误。即使选择 `CONTINUE_ON_ERROR`，服务器也可能为了性能重新安排批次内部的处理顺序，所以跨实体依赖最好拆成明确阶段并核对结果。

### 10.2 为什么错误藏在 details 里

批量写可能只有第 3 和第 8 项失败。单个 gRPC 状态码表达不了每项结果，因此服务器在 `grpc-status-details-bin` 中返回 `google.rpc.Status`，其 `details` 按更新位置携带 `p4.v1.Error`。

客户端至少应记录：

- gRPC canonical code；
- 对应更新下标；
- P4Runtime 错误消息；
- 设备扩展的错误空间与代码（如果存在）。

只打印 `RpcError.details()` 往往会丢掉最有用的逐项原因。本文 Python 示例包含一个精简的 details 解码器。

## 11 Read RPC：用模板查询实体

`ReadRequest.entities` 不是待创建的实体，而是一组查询模板。服务器以流式 `ReadResponse` 返回结果，因此数据量大时不会被迫塞进一个巨大响应。

常见查询方式：

- `table_id = 0`：按规范允许的通配语义读取所有表；
- 给 table ID、不写 match：读取该表的所有表项；
- 给完整键：读取一个具体表项；
- 给 counter ID、不写 index：读取整个计数器数组；
- 给 register ID、不写 index：读取整个寄存器。

不同实体允许的通配字段不同，不能把 TableEntry 的规则机械套给其他实体。响应也可能拆包、改变实体顺序；客户端不能依赖“一次响应就是全部”或“返回顺序等于请求顺序”。

读操作通常不携带选举 ID，但仍可能受角色权限与服务端访问控制限制。大量全表读取会消耗设备与控制通道资源，生产控制器应优先做有边界的查询。

## 12 StreamChannel：一条长期存在的双向通道

`StreamMessageRequest` 与 `StreamMessageResponse` 都用 `oneof update` 区分消息类型。客户端方向包括仲裁、PacketOut、Digest ACK 和扩展消息；服务器方向包括仲裁、PacketIn、Digest、空闲超时通知、流错误和扩展消息。

由于响应可以随时到达，客户端必须持续消费流并按类型分发，不能发送 PacketOut 后就假设下一个响应必定是 PacketIn。本文示例用后台线程和多个队列完成最小分发。

### 12.1 PacketIn 与 PacketOut

- `PacketOut`：控制器把原始包载荷和 P4Info 定义的元数据发给设备；
- `PacketIn`：设备把包载荷和元数据送给控制器。

元数据 ID、位宽和含义来自 P4Info 的 `controller_packet_metadata`。它们通常由带 `@controller_header("packet_in")` 或 `@controller_header("packet_out")` 的 P4 头部生成。端口号依然是 bytestring，不能把本机接口名直接写进协议。

服务器收到未知、重复、越界或不合规的 PacketOut 元数据时，应通过流错误报告，而不是关闭整条流来代替每条消息的错误反馈。控制器也要继续监听 `StreamError`。

### 12.2 Digest

P4 数据平面调用 Digest 后，服务器可以把多条样本批成 `DigestList`。控制器处理后发送 `DigestListAck`，其中的 digest ID 和 list ID 用于确认。`DigestEntry.Config` 的三个参数共同决定延迟、批量大小和未确认列表的抑制行为。

Digest 适合上报“数据平面刚发现了一个控制面可能关心的值”，例如未知源地址；它不是可靠消息队列。控制器要考虑重复、延迟和重连后的状态恢复。

### 12.3 IdleTimeoutNotification

表支持空闲超时时，控制器可在表项中设置 `idle_timeout_ns`。服务器检测到表项超过条件后，通过 `IdleTimeoutNotification` 上报一个或多个表项及时间戳。通知本身不会自动等价于删除，控制器应根据策略决定读取、刷新或删除。

设备检测空闲的精度、通知聚合和延迟可能不同，不能把它当作纳秒级定时器。

## 13 Capabilities、版本与兼容性

`Capabilities` 返回 `p4runtime_api_version` 字符串。它说明服务器实现的 P4Runtime API 版本，不代表当前流水线版本，也不代表设备支持所有可选特性。

P4Runtime 使用语义化版本思想：

- 补丁版本用于兼容修正；
- 次版本可加入向后兼容能力；
- 破坏兼容的改变需要主版本变化。

Protobuf 的未知字段机制有助于新旧端点共存，但不等于所有语义都自动兼容。客户端应忽略自己不理解但可安全忽略的响应字段，发送请求时则只使用已确认支持的能力。弃用字段可能仍保留在线格式中，不能因为代码里标了 deprecated 就立即复用其编号。

## 14 可移植性、扩展和非 PSA 架构

P4Runtime 核心消息尽量与架构无关，但端口编号、extern、设备配置、角色配置和部分实体能力仍可能依赖目标。可移植控制器至少应做到：

- 运行时读取 P4Info，不硬编码对象 ID；
- 查询 API 版本并维护能力矩阵；
- 把端口、架构 URI 和设备配置放进目标适配层；
- 对 `Any` 先检查 `type_url`，不认识就保留或明确拒绝；
- 不把 BMv2 的行为误当作所有硬件都必须如此。

非 PSA 架构可以通过 P4Info extern、`ExternEntry`、`Other` 匹配以及流消息中的 `Any` 扩展协议。扩展应使用稳定、全局可辨识的类型 URL，并避免改变核心消息既有语义。

## 15 已知边界

v1.5.0 规范仍明确保留一些边界，学习时尤其要避免过度推断：

- P4Runtime 不定义完整设备生命周期和端口管理；
- 角色权限配置没有跨厂商统一格式；
- 可选原子性和流水线协调升级依赖目标能力；
- 一些架构对象只能通过 extern 扩展表达；
- 协议不能替控制器解决业务层事务、拓扑一致性和故障恢复策略；
- P4Info 与设备配置必须成对匹配，仅凭对象名字相同不能证明二进制兼容。

协议提供的是一套清晰的控制面“语法和合同”，系统是否可靠仍取决于控制器如何管理期望状态、重试、幂等、版本和观测。

## 16 实操：Python 控制 BMv2

下面用一个很小的 v1model 程序验证完整闭环：建立连接、仲裁、查询版本、下发流水线、写表、读表、PacketOut/PacketIn，最后删除表项并关闭资源。

### 16.1 环境与依赖

本文实际验证的工具版本为：

```text
Python                    3.12.3
p4runtime                 1.5.0
protobuf                  3.20.3
googleapis-common-protos  1.56.4
p4c                       1.2.5.10
simple_switch_grpc        1.15.0
```

`p4runtime==1.5.0` 发布包中的 Python 文件由较旧版本 protoc 生成，直接配合新版 protobuf 运行时会出现 `Descriptors cannot be created directly`。因此示例环境显式固定兼容版本：

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install \
  'p4runtime==1.5.0' \
  'protobuf==3.20.3' \
  'googleapis-common-protos==1.56.4'
```

### 16.2 P4 程序

程序包含一张 IPv4 LPM 表，以及 P4Runtime Packet I/O 所需的两个控制器头部。BMv2 要求普通头部总位宽是 8 的倍数，因此 9 位端口字段后补了 7 位保留字段。

<<< @/sdn/codes/p4runtime/basic.p4

编译时同时生成 BMv2 JSON 与 P4Info：

```bash
cd docs/sdn/codes/p4runtime
p4c-bm2-ss --p4v 16 \
  --p4runtime-files basic.p4info.txtpb \
  -o basic.json basic.p4
```

### 16.3 启动交换机

`--no-p4` 让 BMv2 在没有初始流水线的情况下启动，随后由控制器调用 `VERIFY_AND_COMMIT`。CPU 端口必须显式启用，否则 Packet I/O 不可用。

```bash
simple_switch_grpc \
  --device-id 1 \
  --no-p4 \
  --log-console \
  -- \
  --grpc-server-addr 127.0.0.1:50051 \
  --cpu-port 510
```

这里没有绑定真实网卡。为了在单进程实验里验证双向流，控制器把 PacketOut 的出口指定为 CPU 端口，P4 程序再为它生成 PacketIn。真实网络中应把物理端口交给 BMv2，并由表动作或 punt/clone 逻辑把需要的包送往 CPU 端口。

### 16.4 Python 控制器

控制器直接使用官方 protobuf/gRPC 绑定，没有借助更高层封装，所以每个协议步骤都能在源码中找到。

<<< @/sdn/codes/p4runtime/controller.py

在另一个终端运行：

```bash
source .venv/bin/activate
cd docs/sdn/codes/p4runtime
python controller.py
```

一次成功运行的输出如下。BMv2 1.15.0 报告自己实现的 API 为 1.3.0，这也恰好说明：客户端库版本、规范版本和服务器能力版本是三件不同的事。

```text
arbitration: this client is primary
server P4Runtime API: 1.3.0
pipeline: VERIFY_AND_COMMIT succeeded
table: inserted one entry, read back 1 entry/entries
stream: PacketIn received, ingress_port=510, payload=34 bytes
table: entry deleted
```

### 16.5 对照代码再走一遍

1. `P4InfoHelper` 从编译结果按完整名称或别名找对象，拒绝零个或多个匹配；
2. `Client` 建立 gRPC 通道，并在后台持续读取双向流；
3. `become_primary()` 把仲裁作为流的首条消息并检查 `google.rpc.Status`；
4. `set_pipeline()` 同时下发 P4Info、JSON 和 cookie；
5. `make_table_entry()` 从 P4Info 获取表、字段、动作和参数 ID；
6. `Write(INSERT)` 写入 `10.0.0.2/32`，`Read` 再按 table ID 读回；
7. PacketOut 携带出口元数据，后台线程把返回的 PacketIn 分到专用队列；
8. 删除时清掉动作，只保留表项键；
9. `finally` 无论中途是否报错都会关闭流和 channel。

这个例子刻意保持小，但没有省略协议骨架。继续扩展时，可以把连接管理、P4Info 查询、实体构造、期望状态数据库和重试策略拆成独立模块。

## 17 常见误区

### P4Info 是设备配置吗？

不是。P4Info 给控制器看，设备配置给目标设备执行。二者来自同一次编译并一起下发，但用途不同。

### 有了 gRPC stub 就能随便写表吗？

不能。stub 只提供消息类型；对象 ID、匹配字段、动作参数和位宽必须来自当前 P4Info，写操作还要通过主备仲裁与角色权限检查。

### `Write` 成功就表示整批原子切换吗？

不一定。默认 `CONTINUE_ON_ERROR` 既不保证全成全败，也不保证数据包看不到中间状态。需要更强语义时必须请求并确认目标支持相应原子性。

### PacketIn 是抓包接口吗？

它是由 P4 程序和目标共同定义的控制面报文通道。哪些包上送、带哪些元数据、是否截断，都取决于流水线与设备配置，不等同于无条件镜像所有流量。

### API v1.5 客户端能控制 v1.3 服务器吗？

可能可以，但只能使用双方兼容的消息和能力。应先调用 `Capabilities`，再避免发送旧服务器不理解或不支持的可选语义。本文示例就是这种兼容交互，不代表所有 v1.5 功能都能在该 BMv2 版本上验证。

## 18 术语速查

| 术语 | 一句话理解 |
| --- | --- |
| Target | 被控制的 P4 设备或软件交换机 |
| Client / Controller | 发起 P4Runtime 请求的一方 |
| Server | 实现 P4Runtime 服务并操作目标的一方 |
| P4Info | P4 程序暴露给控制面的接口描述 |
| Entity | 表项、计数器等可读写对象的统一容器 |
| Role | 控制器负责的管理域 |
| Election ID | 同一角色下判定主控制器新旧的 128 位值 |
| Preamble | P4Info 对象共有的 ID、名称、别名和注解 |
| Device config | 目标相关的流水线二进制配置 |
| Bytestring | 按 P4 位宽编码的大端 Protobuf `bytes` |
| PRE | 负责组播和克隆的包复制引擎 |
| StreamChannel | 承载仲裁、Packet I/O 和异步通知的双向流 |
