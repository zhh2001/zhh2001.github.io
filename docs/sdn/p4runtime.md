---
outline: deep
---

# P4Runtime

P4 程序定义数据包的解析和处理逻辑。程序装入设备后，控制器还需要写入路由、更新下一跳、读取计数器，并协调多个控制器的访问权限。P4Runtime 为这些运行期间的操作定义了统一接口。表的默认动作和程序声明的初始表项可以在安装流水线时就已存在，并非所有表都从空表开始。

这份笔记依据 [P4Runtime v1.5.0 规范](https://p4lang.github.io/p4runtime/spec/v1.5.0/P4Runtime-Spec.html) 和同版本的 [Protobuf 定义](https://github.com/p4lang/p4runtime/tree/v1.5.0/proto) 整理。先介绍 P4Runtime 在系统中的位置，再讨论 P4Info、实体和 RPC，最后用 Python 与 BMv2 验证基本交互。示例使用 v1model 架构，API 版本为 1.3.0，文中涉及的 v1.5.0 新特性需要另行确认目标支持。

::: info 规范用语
规范里的 **MUST / MUST NOT** 表示实现必须遵守，**SHOULD / SHOULD NOT** 表示通常应当遵守，只有充分理由才能偏离，**MAY** 表示可选能力。本文引用协议要求时沿用这些强度，部署建议和示例的选择会单独说明。
:::

## 1 P4Runtime 在系统中的位置

编译和控制流程中有五个主要部分：

- **P4 程序**定义解析器、匹配动作表和数据包处理逻辑。
- **设备配置**（device config）是目标相关的流水线配置，例如 BMv2 JSON 或硬件编译产物。
- **P4Info**记录控制面可访问对象的名字、ID、类型和约束。
- **P4Runtime**定义安装配置、操作实体和收发控制报文的协议。
- **控制器**计算路由或策略，并将结果转换成 P4Runtime 请求。

控制器通过 P4Info 了解当前程序的控制接口，再通过 P4Runtime 操作对应对象。路由算法、拓扑计算和业务策略由控制器实现。

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

P4Runtime 定义控制器与服务器之间的行为。服务器可能直接运行在交换机上，也可能是机框管理进程或代理，规范不强制它的内部实现。固定功能设备也可以提供 P4Info 和 P4Runtime 服务，控制器未必能取得对应的 P4 源程序。

## 2 RPC 接口

P4Runtime 使用 Protocol Buffers 描述消息，并通过 gRPC 定义 RPC 服务。`P4Runtime` 服务包含六个 RPC：

| RPC                           | 方向     | 作用                                           |
| ----------------------------- | -------- | ---------------------------------------------- |
| `Write`                       | 一元     | 插入、修改或删除表项等 P4 实体                 |
| `Read`                        | 服务端流 | 按查询条件读取实体，结果可以拆成多个响应返回   |
| `SetForwardingPipelineConfig` | 一元     | 验证、保存或启用新的转发流水线配置             |
| `GetForwardingPipelineConfig` | 一元     | 读取当前 P4Info、设备配置或 cookie             |
| `StreamChannel`               | 双向流   | 仲裁、Packet I/O、Digest、空闲超时通知和流错误 |
| `Capabilities`                | 一元     | 查询服务器支持的 P4Runtime API 版本            |

这里的“一元”指一问一答。  
“服务端流”指一个请求可能对应多个响应。  
“双向流”则允许双方在一条长期连接上独立发送消息。

每台设备由一个 `uint64 device_id` 标识。v1.5.0 要求设备 ID 非零，但 `CapabilitiesRequest.device_id=0` 表示查询服务器共有的能力。ID 的分配属于部署系统的职责，P4Runtime 不提供设备发现或 ID 分配机制。控制器通过带外配置取得服务器地址、凭据与 `device_id`。

::: warning P4Runtime 不包办所有设备管理
端口创建、机框管理、证书发放、设备发现、P4 编译和拓扑计算都不属于核心 P4Runtime API。厂商可以用 gNMI、gNOI 或自己的管理接口补齐这些能力。
:::

## 3 控制器启动流程

从控制器启动到真正转发流量，可以拆成六步：

1. 控制器建立 gRPC 通道并打开 `StreamChannel`。
2. 它首先发送 `MasterArbitrationUpdate`，声明设备、角色和选举 ID。
3. 服务器返回仲裁结果。成为该角色的主控制器后，客户端才可执行受主身份保护的写操作。
4. 控制器调用 `Capabilities`，确认 API 版本是否兼容。
5. 控制器用 `SetForwardingPipelineConfig` 下发 P4Info 和设备配置。
6. 控制器依据 P4Info 中的 ID 构造 `Write`，写入表项、动作选择器成员等实体，或修改计数器和计量器。

```text
Controller                         P4Runtime Server
    │──── StreamChannel ──────────────────►│
    │──── Arbitration(device, role, id) ──►│
    │◄─── Arbitration(status) ─────────────│
    │──── Capabilities ───────────────────►│
    │◄─── api_version ─────────────────────│
    │──── SetPipeline(P4Info + config) ───►│
    │──── Write(TableEntry...) ───────────►│
    │◄─── PacketIn / Digest ───────────────│
    │──── PacketOut / Digest ACK ─────────►│
```

这个流程适用于首次安装程序。只读客户端可以直接查询已有设备，重新连接的控制器也不必重新安装流水线。需要写入的控制器应重新仲裁，再读取、核对或恢复所需实体。

规范将由 P4Runtime 管理的状态称为 P4Runtime 状态，其中部分字段具有数据平面易变性：

- **P4Runtime 状态**：流水线配置和通过实体 API 暴露的状态。
- **数据平面易变状态（data-plane volatile state）**：无需控制器写入，仅由数据包处理就能改变的状态，例如计数器数值、寄存器数据和部分空闲超时字段。

这里的 volatile 描述状态会自行变化，与重启后是否持久化是两个问题。完整服务器重启会清空转发状态，支持 ISSU 的升级流程则可能保留状态。流水线配置操作另有明确的状态保留规则，见第 7 节。

## 4 控制器为什么需要 P4Runtime？

规范列出的用例可以归纳为四类：

- **集中式控制**：控制器计算路由、ACL 或负载均衡策略，再批量写入设备。
- **本地代理控制**：设备上的代理把上层意图转换成 P4Runtime 实体。
- **多控制器高可用**：同一角色有主备控制器，主节点故障后通过仲裁切换。
- **混合管理**：不同角色各管一部分实体，例如路由控制器管理三层表，安全控制器管理 ACL 表。

P4Runtime 的“协议无关”指 RPC 和消息结构不随数据包协议改变。某张表有哪些键和动作仍由 P4Info 决定，因此控制器需要了解当前流水线的接口。

## 5 主备仲裁、角色与选举 ID

### 5.1 为什么写操作必须先仲裁

如果两个控制器同时认为自己是主节点，一个写“端口 1”，另一个写“端口 2”，交换机状态会不断抖动。P4Runtime 用 `StreamChannel` 上的 `MasterArbitrationUpdate` 建立写权限归属。

仲裁消息包含：

| 字段          | 含义                                                   |
| ------------- | ------------------------------------------------------ |
| `device_id`   | 要控制的设备                                           |
| `role`        | 控制器角色，省略时使用具有完整流水线访问权限的默认角色 |
| `election_id` | 128 位无符号选举 ID，通常随新一轮选主单调增大          |
| `status`      | 服务器响应中的仲裁结果                                 |

服务器为每个 `(device_id, role)` 记录已见过的最大选举 ID，即 `election_id_past`。一次合法仲裁更新的 ID 大于或等于这个值时，客户端才能成为主节点。同一设备、同一角色的活动客户端必须使用不同的选举 ID，重复 ID 会使相应流以 `INVALID_ARGUMENT` 结束。

主节点断开后，较小 ID 的备节点不会自动接任。它需要发送新的仲裁更新，使用不小于服务器记录值、且未被其他活动客户端占用的 ID。仲裁响应中的选举 ID 表示服务器已见过的最大值，并不总是当前仍在线的客户端 ID。具体规则见规范的[仲裁更新](https://p4lang.github.io/p4runtime/spec/v1.5.0/P4Runtime-Spec.html#sec-arbitration-updates)与[仲裁通知](https://p4lang.github.io/p4runtime/spec/v1.5.0/P4Runtime-Spec.html#sec-arbitration-notification)。

响应的 `status` 为 `OK` 时，该客户端是主节点。有主节点时，备节点收到 `ALREADY_EXISTS`。当前没有主节点时，备节点收到 `NOT_FOUND`。请求中不填写 `status`。

成为主节点后，`Write` 和设置流水线配置的请求还要携带匹配的 `device_id`、角色与 `election_id`，服务器据此拒绝过期或无权请求。客户端省略 `election_id` 时只能参与备节点仲裁，显式设置的全零 128 位 ID 则仍可参与主节点选举。

数值角色 ID 在 v1.4.0 已弃用。仲裁使用 `Role.name`，`Write` 和设置流水线配置使用字符串 `role`。角色的具体权限通过 `Role.config` 中的 `google.protobuf.Any` 表示，其格式由控制器和服务器带外约定，可以限制可写实体、PacketIn 接收范围等。未设置 `Role.config` 表示完整流水线访问范围，核心规范不提供通用权限配置格式。

### 5.2 StreamChannel 的首条消息

客户端为设备打开流后，首条请求必须是仲裁更新。服务器通过它把流与设备和角色关联起来。已建立的流不能更换设备或角色，需要关闭后另开一条流。连接中断后，原有流上的主身份不再有效。

`Read` 不要求主身份，也不要求预先建立 `StreamChannel`，但仍可能受到身份认证和角色访问范围的限制。备节点可以读取状态，实际系统应控制轮询频率。

::: tip 选举 ID 不是设备版本号
它只比较同一角色下控制器主身份的新旧，不能替代流水线 cookie、配置版本或业务事务号。
:::

## 6 P4Info：连接 P4 程序与控制器的说明书

P4Info 是 `p4.config.v1.P4Info` 消息。可编程目标通常由编译器生成它，固定功能目标可以由设备供应方提供。它主要包含三类信息：

1. **可控对象**：表、动作、动作配置文件、计数器、计量器、Digest、ValueSet、寄存器和架构相关外部对象。
2. **对象元数据**：数值 ID、完整名称、别名、注解和文档。
3. **类型信息**：结构体、头部、枚举、新类型以及序列化所需的位宽。

### 6.1 Preamble 与对象 ID

大多数顶层对象都有 `Preamble`：

- `id` 是控制面使用的 32 位 ID。
- `name` 是 P4 作用域下的完整名称，如 `MyIngress.ipv4_lpm`。
- `alias` 通常是能唯一定位对象的最短名称后缀。
- `annotations`、`structured_annotations` 和 `doc` 保留程序中的说明。

顶层对象 ID 的高 8 位表示对象类型，低 24 位标识具体对象，完整 ID 在一份 P4Info 内必须唯一。编译器生成 ID，但规范不保证跨程序修改或不同编译器时保持不变。`@id` 可以指定后缀，类型前缀必须正确。匹配字段、动作参数和 Packet I/O 元数据的 ID 没有类型前缀，分别在所属表、动作或控制器头部内唯一。控制器应从当前 P4Info 查询 ID。

### 6.2 表和动作如何关联

一张表的 P4Info 会列出：

- 匹配字段的 ID、名称、位宽和匹配类型。
- 可用动作的引用及其适用范围。
- 表容量、常量属性、初始表项标记和默认动作约束。
- 是否关联直接计数器、直接计量器、空闲超时等资源。

动作则列出参数 ID、名称和位宽。控制器先通过表 ID 指定目标，再用字段 ID 填匹配键，用动作 ID 和参数 ID 填执行内容。

```text
P4: table ipv4_lpm                 P4Info
    key dstAddr: lpm     ───────►  table_id + field_id + bitwidth
    action ipv4_forward  ───────►  action_id + param_id + bitwidth
```

### 6.3 PkgInfo 与 P4TypeInfo

`PkgInfo` 可以记录包名、版本、架构、组织、联系方式和编译信息，便于控制器判断拿到的是哪一份程序。`P4TypeInfo` 描述复杂 P4 类型，使 `P4Data` 能表达结构体、元组、头部、头部栈、联合、枚举和错误类型。

底层为 `bit<W>` 的新类型如果带 `@p4runtime_translation`，可以声明标识转换规则的 URI，以及控制面使用的类型。控制面类型可以是另一种位宽的无符号位串，也可以是字符串，分别由 `translated_type.sdn_bitwidth` 或 `sdn_string` 描述。控制器必须按 P4Info 暴露的控制面类型编码。

## 7 转发流水线配置

`ForwardingPipelineConfig` 包含：

- `p4info`：控制面接口描述。
- `p4_device_config`：不透明的设备专用配置字节串。
- `cookie`：控制器可选填的 64 位标识，用于识别配置。

服务器不解释 cookie 的业务含义，控制器可以用它存构建号或配置哈希的截断值。它也不能替代对实际 P4Info 的兼容性检查。

### 7.1 五种设置动作

| Action                 | 行为                                                                             |
| ---------------------- | -------------------------------------------------------------------------------- |
| `VERIFY`               | 验证配置，目标的转发状态不变                                                     |
| `VERIFY_AND_SAVE`      | 验证并保存配置，目标转发状态不变，但后续 `Read` / `Write` 必须使用新配置中的对象 |
| `VERIFY_AND_COMMIT`    | 验证、保存并启用配置，清空目标中的转发状态                                       |
| `COMMIT`               | 启用最近保存且尚未提交的配置，重放保存之后的写请求                               |
| `RECONCILE_AND_COMMIT` | 验证、保存并启用配置，同时保留转发状态，属于可选能力                             |

请求里的 `UNSPECIFIED` 不是有效操作。`COMMIT` 请求必须省略 `config`，携带配置时返回 `INVALID_ARGUMENT`。此前没有保存配置时返回 `NOT_FOUND`。其余四种动作需要提供配置。

`VERIFY_AND_COMMIT` 清空原有转发状态，新流水线仍会建立自身的默认动作和初始表项。`RECONCILE_AND_COMMIT` 则要求保留状态，目标无法为新配置保留状态时返回 `INVALID_ARGUMENT`，未实现该能力时可以返回 `UNIMPLEMENTED`。规范不保证切换期间零丢包，也不限制丢包持续时间。生产系统需要单独设计升级和恢复流程。完整语义见规范的[设置流水线配置](https://p4lang.github.io/p4runtime/spec/v1.5.0/P4Runtime-Spec.html#_setforwardingpipelineconfig_rpc)。

### 7.2 读取配置

`GetForwardingPipelineConfig` 的 `response_type` 可以选择：

- `ALL`：默认值，请求全部内容。
- `COOKIE_ONLY`：只取 cookie。
- `P4INFO_AND_COOKIE`：不取设备配置。
- `DEVICE_CONFIG_AND_COOKIE`：不取 P4Info。

尚未设置流水线时，响应的 `config` 字段不设置。已有配置时，服务器必须能够返回 P4Info，但可以不支持取回 `p4_device_config`。未配置 cookie 时，响应也不含 cookie。设备配置可能很大，日常核对通常只需请求 cookie 或 P4Info 与 cookie。

## 8 消息编码

### 8.1 默认值与读写对称

Proto3 的普通标量字段通常不区分未设置与显式默认值，消息字段和 `oneof` 成员则具有存在性。P4Runtime 中零值的含义取决于字段，例如读取时 `table_id=0` 表示通配，而 `Index{index: 0}` 表示数组第一个元素，省略 `Index` 消息才表示通配。构造请求时需要查对应实体的规则。

**读写对称（read-write symmetry）** 要求成功写入的实体随后读出时与写入内容相符，数据平面易变字段和规范明确列出的例外除外。重复字段的顺序通常不影响对称性，数值可以规范化为最短字节串。它不表示任意读响应都能直接作为写请求，只读字段和各操作的约束仍需处理。

### 8.2 Bytestring

P4 的定宽整数可能超过 Protobuf 整数的位宽，因此许多 P4 值用 `bytes` 表示。规则要点是：

- 使用大端序。
- 位宽从 P4Info 获得。
- 无符号值使用零扩展，带符号值使用二进制补码和符号扩展。
- 数值必须落在声明位宽范围内。
- 规范形式使用能容纳数值的最短字节串，零至少使用一个字节。

例如 `bit<9>` 的十进制 2 规范编码为 `02`，510 编码为 `01 fe`。接收方也必须接受正确扩展且数值在范围内的非最短表示，例如用 `00 02` 表示 2。空字节串不是整数零的有效编码。完整规则见规范的[Bytestrings](https://p4lang.github.io/p4runtime/spec/v1.5.0/P4Runtime-Spec.html#sec-bytestrings)。

v1.x 的表匹配键、动作参数和 Packet I/O 元数据只支持无符号整数及规定的类型转换，不能因 bytestring 支持补码就直接使用 `int<W>`。带符号整数主要通过 `P4Data` 表达。若类型被转换为字符串，则按对应字符串编码规则处理。

::: warning 位宽不是字节数
Python 的 `int.to_bytes()` 接收的是字节数。表示 `bit<9>` 范围内的值至多需要 2 字节，合法的非最短编码可以更长。本文示例的 `encode()` 检查数值范围并生成最短表示。
:::

### 8.3 P4Data

`P4Data` 是 P4 值的通用容器，通过 `oneof` 表示定宽位串、变长位串、布尔值、元组、结构体、头部、头部栈、联合、枚举和错误等。`bit<W>` 与 `int<W>` 共用 `bitstring` 字段，具体解释由类型描述决定。变长位串还需要携带当前有效位宽。

元组和结构体的成员按类型声明顺序序列化。头部需要表示有效性，无效头部不携带字段值，联合则标识当前有效成员。解码时需要同时查看对象的类型描述及 `P4TypeInfo`，仅凭 Protobuf 字段不足以还原 P4 类型。

## 9 Entity：P4Runtime 真正管理的对象

`Entity` 是一个 `oneof` 容器，一次只承载一种实体。v1.5.0 的核心实体包括表项、动作配置文件成员与组、计数器、计量器、包复制引擎、ValueSet、寄存器、Digest 和外部对象。

### 9.1 TableEntry

`TableEntry` 是最常用的实体，关键字段有：

| 字段                | 作用                                         |
| ------------------- | -------------------------------------------- |
| `table_id`          | P4Info 中的表 ID                             |
| `match`             | 一组字段匹配，字段顺序不应影响语义           |
| `action`            | 直接动作、成员 ID、组 ID 或一次性动作集合    |
| `priority`          | 需要优先级的表项使用正数，数值越大优先级越高 |
| `is_default_action` | 操作默认动作而不是普通表项                   |
| `idle_timeout_ns`   | 请求空闲超时通知的时间                       |
| `metadata`          | 控制器自用的透明字节串                       |

匹配类型包括：

- `exact`：精确相等。
- `lpm`：值加前缀长度。
- `ternary`：值加掩码。
- `range`：闭区间下界与上界。
- `optional`：提供该匹配字段时精确匹配，省略时通配。这里的“存在”指请求中的字段，与 P4 头部有效性无关。
- `other`：为扩展匹配类型预留的 `Any`。

LPM 的值在前缀以外必须为零，三元匹配的值在掩码为零的位置也必须归零。提供 LPM 字段时，前缀长度必须在 1 到字段位宽之间。`range` 的边界必须在字段范围内，且下界不大于上界。

`exact` 字段必须提供。其余标准匹配类型的全通配需要省略整个字段，不能使用前缀长度 0、全零掩码或覆盖全部值的区间代替。数值的 bytestring 长度仍可以有不同的合法表示。

表定义只要包含 `ternary`、`range` 或 `optional`，普通表项就必须给出正的 `priority`，即使该表项省略了这些字段。只含 `exact` 和 LPM 的表使用零。相同优先级的多个表项若同时匹配一个包，选择结果没有确定性保证。

普通表项的键由 `table_id`、`match` 和 `priority` 组成。`INSERT` 要求键不存在，`MODIFY` 和 `DELETE` 要求目标存在。普通表项的 `MODIFY` 若省略 `action`，会保留当前动作。`DELETE` 只依据键和 `is_default_action` 定位表项，服务器必须忽略其他字段，包括动作。示例单独构造只含键的删除请求，便于看清操作意图。

默认动作通过 `is_default_action=true` 指定，`match` 必须为空，`priority` 必须为零。默认表项只允许 `MODIFY`，不允许 `INSERT` 或 `DELETE`。省略 `action` 的默认表项修改会恢复 P4 程序声明的默认动作，这与普通表项不同。常量默认动作不可修改。

`is_const_table`、`has_initial_entries` 和读响应中的 `is_const` 分别描述常量表、初始表项和常量表项。常量条目的匹配键和动作不可改，但直接资源仍可按规则更新，未声明为 `const` 的默认动作也可修改。写请求中的 `is_const` 必须为 false。完整约束见规范的[TableEntry](https://p4lang.github.io/p4runtime/spec/v1.5.0/P4Runtime-Spec.html#sec-table-entry)。

### 9.2 ActionProfile：成员、组与一次性选择

带选择器的动作配置文件有两种常见编程方式：

1. **手工成员/组模式**：先写 `ActionProfileMember`，再写引用成员的 `ActionProfileGroup`，表项引用成员或组 ID。
2. **one-shot 模式**：表项直接携带 `ActionProfileActionSet`，由服务器管理内部成员。

动作集合包含动作、权重和监视端口。v1.5.0 的 `action_selection_mode` 可以选择默认方式、`HASH` 或 `RANDOM`，`size_semantics` 则决定按权重和（`SUM_OF_WEIGHTS`）还是成员数（`SUM_OF_MEMBERS`）计算容量。这是选择算法和资源计量两个独立维度。若 P4Info 的 `weights_disallowed=true`，客户端必须将权重留为零，否则权重需要为正。

同一个 action selector 关联的所有表必须统一使用一种编程方式，不能在已有条目中混用 one-shot 和手工成员/组。规范要求支持 one-shot，手工成员/组方式是可选能力。没有选择器的 action profile 则通过成员 ID 编程，不使用组或 one-shot。

手工方式需要分批创建成员、组、表项，删除时按相反的依赖顺序操作。一个 `WriteRequest` 内的更新可能被服务器重新排序，不能靠同一批次中的排列顺序保证引用已存在。

### 9.3 Counter 与 Meter

计数器分为：

- `CounterEntry`：独立计数器，通过 `counter_id + index` 定位。
- `DirectCounterEntry`：直接绑定表项，通过表项键定位。

`CounterData` 包含包数和字节数。两类计数器只接受 `MODIFY` 更新，`INSERT` 和 `DELETE` 返回 `INVALID_ARGUMENT`。省略 `data` 不改变数值，提供全零数据则将计数清零。独立计数器省略 `index` 时，可以读取或修改整个数组。数据平面同时处理包时，数值还会继续变化。

计量器分 `MeterEntry` 和 `DirectMeterEntry`，同样只接受 `MODIFY`。与计数器不同，省略 `config` 会恢复默认的全 GREEN 行为。`MeterConfig` 包含承诺速率和突发量 `cir/cburst`、峰值速率和突发量 `pir/pburst`，以及超额突发量 `eburst`。

单位和计量器类型由 P4Info 的 meter spec 决定。双速率三色计量器不使用 `eburst`，单速率三色计量器要求 `pir=cir`、`pburst=cburst`，再用 `eburst` 指定超额突发量。单速率双色计量器也要求前述两组值相等，并不使用 `eburst`。支持按颜色计数的目标还可通过 `counter_data` 暴露各颜色的计数值。

直接资源也可以随 `TableEntry` 一起读写。读取时，需要在查询模板中设置空的 `counter_data` 或 `meter_config` 消息，服务器才返回对应资源。修改表项时，省略 `counter_data` 保留计数值，省略 `meter_config` 却会重置计量器。要保留原计量器配置，需要在 `MODIFY` 中再次提供它。默认的全 GREEN 配置在读响应中以未设置的 `meter_config` 表示。

### 9.4 Packet Replication Engine

包复制引擎（PRE）实体包括：

- `MulticastGroupEntry`：组 ID 加多个副本，每个副本指定端口和实例。
- `CloneSessionEntry`：克隆会话、副本、服务等级和可选的截断长度。

`Replica` 从 v1.4.0 起使用 `port` bytestring，旧的 `egress_port` 字段已经弃用。v1.5.0 增加备用副本列表，目标对此的支持是可选的，同一个组内的主副本和备用副本都需要满足 `(port, instance)` 唯一性。

P4Runtime 写入的组 ID 和克隆会话 ID 必须非零，零保留给通配读取。PSA 数据平面的克隆会话 ID 可以为零，控制面可通过类型转换映射到这个值。其余合法范围及端口编码仍受目标约束。

### 9.5 ValueSet、Register、Digest 与 Extern

- **ValueSetEntry**：用 `MODIFY` 替换指定解析器 ValueSet 的全部成员，空成员列表表示清空。`INSERT` 和 `DELETE` 不合法。成员由一组 `FieldMatch` 表示，列表不能重复或超过 P4Info 声明的容量。
- **RegisterEntry**：用 `register_id + index` 访问寄存器单元，数据放在 `P4Data`。省略索引可读取整个数组，也可将同一个值写入整个数组。数组索引为有符号整数，负数无效，越界会返回错误。
- **DigestEntry**：配置最大等待时间、最大列表长度和确认超时，通过 `INSERT` 启用、`MODIFY` 调整、`DELETE` 停止生成 Digest。数据通过 `StreamChannel` 上报。
- **ExternEntry**：用外部类型 ID、实例 ID 和 `google.protobuf.Any` 承载非核心架构对象。

P4Info 描述对象的类型和约束，Entity 指定这次访问的实例。具体支持哪些实体还取决于目标架构。

## 10 Write RPC：批量修改与原子性

`WriteRequest` 包含设备、角色、选举 ID、一组 `Update` 和原子性选项。每个 `Update` 的类型是：

- `INSERT`：创建不存在的实体。
- `MODIFY`：修改已经存在的实体。
- `DELETE`：删除实体。

这三个操作不是每种实体都支持，例如计数器、计量器和 ValueSet 只接受 `MODIFY`。`UNSPECIFIED` 不是可执行更新。服务器必须校验客户端是否为对应角色的主节点，并逐项校验 ID、字段、位宽、引用和权限。

### 10.1 三种原子性

| Atomicity           | 含义                                                                 |
| ------------------- | -------------------------------------------------------------------- |
| `CONTINUE_ON_ERROR` | 必选默认能力，某项失败仍尝试其余更新，数据包可能看到批次中的中间状态 |
| `ROLLBACK_ON_ERROR` | 可选，遇错恢复控制面和数据面状态，但处理期间数据包仍可能看到中间状态 |
| `DATAPLANE_ATOMIC`  | 可选，整批作为事务处理，每个数据包只能看到批处理前或成功完成后的状态 |

即使选择 `CONTINUE_ON_ERROR`，每个单独更新也必须对数据包处理保持原子性。它允许看到多个更新之间的状态，不允许看到一个表项更新到一半的状态。后两种模式依赖目标支持，未实现时返回 `UNIMPLEMENTED`。跨实体依赖应拆成阶段并核对结果。

### 10.2 批量错误报告

批量写可能只有部分更新失败。处理批次后，只要存在失败项，RPC 的整体状态就使用 `UNKNOWN`，服务器在 `grpc-status-details-bin` 中返回 `google.rpc.Status`。其 `details` 按更新位置逐项携带 `p4.v1.Error`，成功项也需要以 `OK` 占位，客户端不能把失败项下标重新编号。整批处理前发生的请求级错误，例如无主身份，则直接使用相应的 gRPC 状态码，不属于这种逐项结果。

客户端至少应记录：

- gRPC canonical code。
- 对应更新下标。
- P4Runtime 错误消息。
- 设备扩展的错误空间与代码（如果存在）。

只打印 `RpcError.details()` 往往会丢掉逐项原因。本文示例的解码器保留更新下标，并输出设备扩展的 `space` 和 `code`。

## 11 Read RPC：用模板查询实体

`ReadRequest.entities` 不是待创建的实体，而是一组查询模板。服务器以流式 `ReadResponse` 返回结果，因此数据量大时不会被迫塞进一个巨大响应。

常见查询方式：

- `table_id=0` 且其余过滤字段缺省：读取所有表的普通表项。
- 给 table ID、不写 match 且其余过滤字段缺省：读取该表的普通表项。
- `is_default_action=true`：读取默认表项，可以结合 table ID 指定表。
- 给完整的非通配键：读取匹配的具体表项。全通配条目的空 `match` 会与全表查询重合，需要从结果中筛选。
- 给 counter ID、不写 index：读取整个计数器数组。
- 给 register ID、不写 index：读取整个寄存器。

查询全部表项也要将 `Entity` 的 `oneof` 设为 `table_entry`，空 `Entity` 不能表示读取所有实体。不同实体允许的通配字段不同，不能机械套用 TableEntry 的规则。

响应可以拆分并改变实体顺序，多个重叠查询模板也可能返回重复实体。流可能先返回部分结果，再以错误结束，客户端要读取到最终状态才能确认查询成功。

`ReadRequest` 没有选举 ID 字段，也不要求主身份或仲裁流。并发 `Read` 与 `Write` 的结果必须与整批写入开始前或完成后的控制面状态一致，但数据平面易变字段不构成静止快照。大量全表读取会消耗设备与控制通道资源，生产控制器应控制查询范围。

## 12 StreamChannel：一条长期存在的双向通道

`StreamMessageRequest` 与 `StreamMessageResponse` 都用 `oneof update` 区分消息类型。客户端方向包括仲裁、PacketOut、Digest ACK 和扩展消息。服务器方向包括仲裁、PacketIn、Digest、空闲超时通知、流错误和扩展消息。

由于响应可以随时到达，客户端必须持续消费流并按类型分发，不能发送 PacketOut 后就假设下一个响应必定是 PacketIn。本文示例用后台线程和多个队列完成最小分发。

### 12.1 PacketIn 与 PacketOut

- `PacketOut`：控制器把原始包载荷和 P4Info 定义的元数据发给设备。
- `PacketIn`：设备把包载荷和元数据送给控制器。

元数据 ID、位宽和含义来自 P4Info 的 `controller_packet_metadata`，通常由带 `@controller_header("packet_in")` 或 `@controller_header("packet_out")` 的 P4 头部生成。声明的每个元数据字段都必须出现一次，不能省略、重复或增加未知 ID。PacketIn 的元数据顺序没有保证，客户端应按 ID 查找字段。

本例的端口使用无符号 bytestring。若 P4Info 声明了类型转换，则按转换后的位宽或字符串类型编码。`@controller_header` 是可选的，不使用它时，控制器需要自行在 payload 中序列化或解析控制器头部。

服务器收到不合规的 PacketOut 元数据时必须丢弃该包，可以通过 `StreamError` 报告，但错误报告是可选能力。没有收到流错误不能证明包已发送，PacketOut 也没有逐包成功确认。规则见规范的[Packet I/O](https://p4lang.github.io/p4runtime/spec/v1.5.0/P4Runtime-Spec.html#sec-packet-i_o)。

### 12.2 Digest

P4 数据平面调用 Digest 后，服务器可以把多条样本批成 `DigestList`。控制器处理后发送 `DigestListAck`，其中的 digest ID 和 list ID 用于确认。`DigestEntry.Config` 的 `max_timeout_ns` 和 `max_list_size` 控制组批，`ack_timeout_ns` 控制未确认列表在重复抑制缓存中的保留时间。

Digest 适合上报未知源地址等控制面可能关心的值。确认用于管理重复抑制缓存，不提供可靠消息队列的重传保证。控制器需要考虑丢失、重复、延迟和重连后的状态恢复。

### 12.3 IdleTimeoutNotification

表支持空闲超时时，控制器可在普通表项中设置非零 `idle_timeout_ns`，零表示不超时。默认表项不支持空闲超时。服务器尽力检测空闲条件，并通过 `IdleTimeoutNotification` 上报表项及时间戳。查询距离上次命中的时间，需要在读请求中设置空的 `time_since_last_hit` 消息，该字段不能用于写入。

PSA 的通知模式由控制器决定是否删除，发送通知后服务器重置通知计时器，因此持续空闲的表项可以再次被上报。通知允许聚合和丢弃，不保证可靠交付。PNA 另有数据平面自动删除模式，不能将 PSA 的行为推广到所有架构。

设备检测空闲的精度、通知聚合和延迟可能不同，不能把它当作纳秒级定时器。

### 12.4 流错误与连接结束

`StreamError` 报告某条异步请求处理失败，可以携带错误码、设备扩展错误和相关请求。服务器不保证发送所有流错误，负载较高时也可以丢弃错误报告。客户端既要分发这种消息，也要处理 gRPC 流本身的结束或异常。

非法仲裁更新属于连接级错误，服务器需要以相应的 gRPC 状态结束流，不能用 `StreamError` 代替。示例在等待仲裁或 PacketIn 时检查错误队列，避免将连接失败误报为普通超时。

## 13 Capabilities、版本与兼容性

`Capabilities` 返回 `p4runtime_api_version` 字符串。它说明服务器实现的 P4Runtime API 版本，不代表当前流水线版本，也不代表设备支持所有可选特性。

v1.5.0 增加 `CapabilitiesRequest.device_id` 和响应的 `experimental`（`Any`）字段。前者可以限定设备能力范围，后者用于双方约定的实验功能。核心响应仍没有一份完整的可选能力清单，支持哪些原子性模式等信息需要通过目标文档或带外约定确认。

P4Runtime 使用语义化版本思想：

- 补丁版本用于兼容修订。
- 次版本可加入向后兼容能力。
- 破坏兼容的改变需要主版本变化。

Protobuf 的未知字段机制有助于新旧端点共存，但不等于所有语义都自动兼容。客户端应忽略自己不理解但可安全忽略的响应字段，发送请求时则只使用已确认支持的能力。弃用字段可能仍保留在线格式中，不能因为代码里标了 deprecated 就立即复用其编号。

## 14 可移植性、扩展和非 PSA 架构

P4Runtime 核心消息尽量与架构无关，但端口编号、extern、设备配置、角色配置和部分实体能力仍可能依赖目标。可移植控制器至少应做到：

- 运行时读取 P4Info，按名称解析对象 ID。
- 查询 API 版本，并记录目标支持的可选能力。
- 把端口、架构 URI 和设备配置放进目标适配层。
- 对 `Any` 先检查 `type_url`，再按对应扩展规则处理。
- 不把 BMv2 的行为误当作所有硬件都必须如此。

非 PSA 架构可以通过 P4Info extern、`ExternEntry`、`Other` 匹配以及流消息中的 `Any` 扩展协议。扩展应使用稳定、全局可辨识的类型 URL，并避免改变核心消息既有语义。

## 15 已知边界

v1.5.0 规范仍明确保留一些边界，学习时尤其要避免过度推断：

- P4Runtime 不定义完整设备生命周期和端口管理。
- 角色权限配置没有跨厂商统一格式。
- 可选原子性和流水线协调升级依赖目标能力。
- 一些架构对象只能通过 extern 扩展表达。
- 业务层事务、拓扑一致性和故障恢复策略需要由控制系统设计。
- P4Info 与设备配置必须成对匹配，仅凭对象名字相同不能证明二进制兼容。

控制器需要持续核对期望状态与设备状态，并根据具体操作设计重试和恢复逻辑。例如，在响应丢失后重试 `INSERT` 可能得到 `ALREADY_EXISTS`，此时需要读取内容，而不能仅凭错误码判断前一次操作未生效。

## 16 实操：Python 控制 BMv2

下面用一个 v1model 程序验证连接、仲裁、查询版本、下发流水线、写表、读表和 PacketOut/PacketIn，最后删除表项并核对结果。实验不绑定网卡，表项读写与 CPU 回送分别验证，不涉及物理端口转发。

### 16.1 环境与依赖

本文实际验证的工具版本为：

```text
Python                    3.12.3
p4runtime                 1.5.0
protobuf                  3.20.3
googleapis-common-protos  1.56.4
grpcio                    1.59.3
p4c                       1.2.5.10
simple_switch_grpc        1.15.0
```

`p4runtime==1.5.0` 发布包中的 Python 文件由较旧版本 protoc 生成，配合不兼容的新版 protobuf 运行时会出现 `Descriptors cannot be created directly`。示例固定已验证的依赖组合。在仓库根目录执行以下命令，这组命令需要本机已有 P4 编译器和 BMv2：

```bash
python3 -m venv .venv
source .venv/bin/activate
python -m pip install \
  'p4runtime==1.5.0' \
  'protobuf==3.20.3' \
  'googleapis-common-protos==1.56.4' \
  'grpcio==1.59.3'
```

### 16.2 P4 程序

程序包含一张 IPv4 LPM 表和两个控制器头部。BMv2 要求普通头部总位宽是 8 的倍数，因此 9 位端口字段后补了 7 位 `reserved`。这个字段也会进入 P4Info，控制器必须将它作为元数据发送，生成 PacketIn 时程序将其置零。字段名由头部声明提供，不需要额外的 `@controller_metadata` 注解。

普通端口的 IPv4 路径只处理无选项的头部，检查解析错误、校验和、版本、IHL 和最小总长度。转发动作显式设置源、目的 MAC，递减 TTL。TTL 不足、非 IPv4 包和未命中表项的包会被丢弃。这里没有实现 ICMP、ARP 或完整的 IPv4 路由器功能。

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

这里没有绑定真实网卡。控制器把 PacketOut 的出口指定为 CPU 端口，P4 程序为它生成 PacketIn。该分支直接指定出口，绕过 IPv4 表，不会命中刚写入的路由。真实网络中需要绑定 BMv2 的端口，并由表动作或 punt/clone 逻辑将所需数据包送往 CPU 端口。CPU 端口 510 是本例约定，不是 P4Runtime 规定的通用端口号。

### 16.4 Python 控制器

控制器直接使用官方 protobuf/gRPC 绑定，没有借助更高层封装，所以每个协议步骤都能在源码中找到。

<<< @/sdn/codes/p4runtime/controller.py

在另一个终端，从仓库根目录运行：

```bash
source .venv/bin/activate
cd docs/sdn/codes/p4runtime
python controller.py
```

本机成功运行的输出如下。这里的 BMv2 1.15.0 构建报告 API 为 1.3.0，客户端库版本和服务器实现版本需要分别确认。

```text
arbitration: this client is primary
server P4Runtime API: 1.3.0
pipeline: VERIFY_AND_COMMIT succeeded
table: inserted one entry, read back 1 entry/entries
stream: PacketIn received, ingress_port=510, payload=34 bytes
table: entry deleted, read back 0 entries
```

多控制器测试中，本机这个 BMv2 构建在主节点断开后，会自动将较小 ID 的备节点提升为主节点，没有遵循第 5 节所述的历史最大选举 ID 规则。该构建还会接受缺少 `reserved` 元数据的 PacketOut 并回送，而规范要求丢弃这种包。本例的成功运行说明基本交互可用，验证 v1.5.0 协议一致性还需要核验服务器实现。

### 16.5 对照代码再走一遍

1. `P4InfoHelper` 按完整名称或别名找对象，拒绝零个或多个匹配。
2. `Client` 建立 gRPC 通道，在后台持续读取双向流，并记录流错误。
3. `become_primary()` 发送首条仲裁消息，核对设备、默认角色、状态和完整的 128 位选举 ID。
4. `set_pipeline()` 同时下发 P4Info、JSON 和 cookie，使用 `VERIFY_AND_COMMIT` 重建流水线状态。
5. `make_table_entry()` 从 P4Info 获取表、字段、动作和参数 ID。
6. `Write(INSERT)` 写入 `10.0.0.2/32`，`Read` 读回后按字段值核对键和动作，忽略重复字段的顺序差异。
7. PacketOut 携带 P4Info 声明的全部元数据。返回的 PacketIn 按元数据 ID 解码，并核对入口端口、保留字段和 payload。
8. `DELETE` 使用只含表项键的请求，再读表确认条目已不存在。
9. `finally` 无论中途是否报错都会关闭流和 channel。

默认选举 ID 为 1，适用于新启动的单控制器实验。连接已有服务器时，可用 `--election-id` 指定经过协调的值，不能重复占用其他活动客户端的 ID。生产控制器还需要管理身份变化、重连、期望状态和重试，本例不实现这些机制。代码使用明文的 loopback gRPC 通道，跨主机部署需要按服务器的认证和传输安全配置建立连接。

## 17 常见误区

### P4Info 是设备配置吗？

P4Info 描述控制面接口，设备配置是目标专用的流水线配置。本例由同一次编译生成二者，安装时需要匹配。固定功能设备的 P4Info 也可以由供应方提供，未必存在供控制器下发的编译产物。

### 有了 gRPC stub 就能随便写表吗？

stub 提供 RPC 调用和消息类型。对象 ID、匹配字段、动作参数和位宽需要来自当前 P4Info，写操作还要通过主备仲裁与角色权限检查。

### `Write` 成功就表示整批原子切换吗？

不一定。默认 `CONTINUE_ON_ERROR` 既不保证全成全败，也不保证数据包看不到中间状态。需要更强语义时必须请求并确认目标支持相应原子性。

### PacketIn 是抓包接口吗？

它是由 P4 程序和目标共同定义的控制面报文通道。哪些包上送、带哪些元数据、是否截断，都取决于流水线与设备配置，不等同于无条件镜像所有流量。

### API v1.5 客户端能控制 v1.3 服务器吗？

可能可以，但只能使用双方兼容的消息和能力。应先调用 `Capabilities`，再避免发送旧服务器不理解或不支持的可选语义。本文示例就是这种兼容交互，不代表所有 v1.5 功能都能在该 BMv2 版本上验证。

## 18 术语速查

| 术语                | 一句话理解                                      |
| ------------------- | ----------------------------------------------- |
| Target              | 被控制的 P4 设备或软件交换机                    |
| Client / Controller | 发起 P4Runtime 请求的一方                       |
| Server              | 实现 P4Runtime 服务并操作目标的一方             |
| P4Info              | P4 程序暴露给控制面的接口描述                   |
| Entity              | 表项、计数器等可读写对象的统一容器              |
| Role                | 控制器负责的管理域                              |
| Election ID         | 同一设备、同一角色下参与主控制器仲裁的 128 位值 |
| Preamble            | P4Info 对象共有的 ID、名称、别名和注解          |
| Device config       | 目标相关的流水线配置字节串                      |
| Bytestring          | 按 P4 位宽编码的大端 Protobuf `bytes`           |
| PRE                 | 负责组播和克隆的包复制引擎                      |
| StreamChannel       | 承载仲裁、Packet I/O 和异步通知的双向流         |
