---
outline: deep
---

# P4 INT

带内网络遥测（In-band Network Telemetry，简称 INT）

## 1 引言

INT 是一种无需控制平面干预即可由数据平面直接采集和上报网络状态的框架。在 INT 架构中，数据包中可以携带由网络设备解释为“遥测指令”的特殊报文头字段。INT 流量的生成端（包括应用程序、终端主机网络协议栈、虚拟机管理器、网卡、发送端的顶层汇聚交换机等）可以将这些指令嵌入普通数据包、克隆数据包，或专门构造的探测包中。另一种方式是直接在数据平面中预设规则，对匹配的流量自动执行指定的遥测指令。

这些指令的作用是明确告知支持 INT 功能的网络设备需要收集哪些状态信息。所采集的网络状态可由数据平面直接导出至遥测监控系统，或在数据包转发过程中动态写入包头。当状态信息以嵌入形式随数据包在网络中传播时，流量接收端便可提取这些信息，并按需上报，从而实现对数据包实际经过路径中数据平面状态的精确监控。

流量接收端可根据具体需求执行以下几类典型操作：

- **运维监控（OAM）**：[接收端](#sink)可提取并上报嵌入的数据平面状态，上传至外部控制器。这些信息可以是未经处理的原始数据，也可进行初步加工，例如压缩、去重或截断。
- **实时控制与反馈机制**：接收端可基于采集的遥测信息，向流量源反馈控制信号，驱动其进行流量工程调整或转发路径优化。例如，显式拥塞通知（ECN）机制就是一种典型的反馈控制手段。
- **网络事件检测**：若路径状态表明存在严重拥塞或违反某些数据平面不变量的情况，[接收端](#sink)可立即触发响应动作，以应对当前的网络事件。这类响应可以是集中式的，也可以是类似 TCP 的全分布式反馈控制。

::: tip <span id="sink">接收端 sink</span>
虽然这些操作通常由接收端节点（Sink nodes）完成，但中继节点（Transit nodes）也可能会发起运维管理（OAM）操作，或执行网络事件检测。
:::

INT 架构模型具有高度的通用性，能够支持多种高层次的网络应用场景，包括但不限于以下几类：

- **网络故障排查与性能监测**：例如路径追踪（Traceroute）、微突发检测（micro-burst detection）以及数据包历史记录（又称为“[明信片](#postcards)”机制）。
    ::: info <span id="postcards">明信片 postcards</span>
    Handigol, Nikhil, Brandon Heller, Vimalkumar Jeyakumar, David Mazières, and Nick McKeown. "I know what your packet did last hop: Using packet histories to troubleshoot networks." In 11th USENIX Symposium on Networked Systems Design and Implementation (NSDI 14), pp. 71-85. 2014.
    :::
- **高级拥塞控制策略**：利用实时遥测数据进行更灵敏和精准的拥塞响应与调节。
- **智能路由机制**：如基于链路利用率感知的路由方法，代表性方案包括 [HULA](#HULA) 和 [CLOVE](#CLOVE) 等。
    ::: info <span id="HULA">HULA</span>
    Naga Katta, Mukesh Hira, Changhoon Kim, Anirudh Sivaraman, and Jennifer Rexford. 2016. HULA: Scalable Load Balancing Using Programmable Data Planes. In Proceedings of the Symposium on SDN Research (SOSR '16). Association for Computing Machinery, New York, NY, USA, Article 10, 1–12. https://doi.org/10.1145/2890955.2890968
    :::
    ::: info <span id="CLOVE">CLOVE</span>
    Naga Katta, Aditi Ghag, Mukesh Hira, Isaac Keslassy, Aran Bergman, Changhoon Kim, and Jennifer Rexford. 2017. Clove: Congestion-Aware Load Balancing at the Virtual Edge. In Proceedings of the 13th International Conference on emerging Networking EXperiments and Technologies (CoNEXT '17). Association for Computing Machinery, New York, NY, USA, 323–335. https://doi.org/10.1145/3143361.3143401
    :::
- **数据平面验证**：通过遥测信息对网络数据平面的行为与一致性进行验证。

上述应用的具体用例及评估可参考《[Millions of Little Minions](#Millions_of_Little_Minions)》一文中的详细论述。

::: info <span id="Millions_of_Little_Minions">Millions of Little Minions</span>
Vimalkumar Jeyakumar, Mohammad Alizadeh, Yilong Geng, Changhoon Kim, and David Mazières. 2014. Millions of little minions: using packets for low latency network programming and visibility. In Proceedings of the 2014 ACM conference on SIGCOMM (SIGCOMM '14). Association for Computing Machinery, New York, NY, USA, 3–14. https://doi.org/10.1145/2619239.2626292
:::

## 2 术语说明

**监控系统（Monitoring System）**：负责收集来自不同网络设备的遥测数据的系统。该系统的组件在物理上可能是分布式的，但逻辑上应视为集中统一的实体。

**INT 报文头（INT Header）**：用于携带 INT 信息的数据包头部，主要分为三类：嵌入数据类型（MD-type）、嵌入指令类型（MX-type）和目标类型（Destination-type），具体定义参见[第 5.1 节](#sec51)。

**INT 报文（INT Packet）**：携带 INT 报文头的数据包。

**INT 节点（INT Node）**：具备 INT 功能的网络设备，参与 INT 数据平面操作，并执行如下至少一项任务：插入、附加、移除或处理 INT 报文头中的指令。根据部署情境，典型的 INT 节点包括路由器、交换机、网卡等。

**INT 指令（INT Instruction）**：用于指定在各 INT 节点处需采集哪些 INT 元数据（定义见下文）的遥测指令。该指令可通过配置 INT 节点中的流监控表（Flow Watchlist）进行设置，或直接写入数据包的 INT 报文头中。

**流监控表（Flow Watchlist）**：位于数据平面中的一张匹配表，基于数据包头字段进行匹配，对命中的流执行相应的 INT 指令插入或应用操作。一个流通常指具有相同关键字段值的一组数据包。

**INT 源（INT Source）**：受信任的实体，负责生成并插入 INT 报文头至其发送的数据包中。可通过配置流监控表来选择哪些流应嵌入 INT 报文头。

**INT 接收端（INT Sink）**：受信任的实体，负责从数据包中提取 INT 报文头，收集其中包含的路径状态信息，并在必要时将该报文头移除，以确保 INT 对上层协议透明。值得注意的是，这并不排除在多级或嵌套 INT 域中部署接收端。INT Sink 可将提取的数据上报至遥测监控系统。

**INT 中继跳点（INT Transit Hop）**：受信任的实体，根据 INT 指令从数据平面中采集元数据。数据既可直接导出至监控系统，也可写入数据包的 INT 报文头中。

需要指出的是，一个物理设备在同一或不同流上可能同时承担多个角色，如 INT 源、中继、接收端。例如，某 INT 源节点可将自身元数据嵌入数据包，因而在逻辑上也充当了中继跳点的角色。

**INT 元数据（INT Metadata）**：由 INT 源节点或中继节点插入至 INT 报文头或遥测报告中的信息。典型元数据示例详见第 4 节。

**INT 域（INT Domain）**：由一组处于同一管理体系下、互联的 INT 节点构成的域。当前规范定义了不同厂商设备在同一 INT 域内交互时的行为规范与报文头格式，以保障互操作性。域内所有 INT 节点需进行一致性配置，建议在域边界部署 INT 接收端能力，以防止遥测信息泄露至域外。

## 3 INT 的运行模式

自 INT 于 2015 年在 P4.org 社区首次提出以来，其架构已衍生出多种变体，并在 IETF 及业界社区中持续被探讨与演进。同时，“INT”这一术语也逐渐被广义使用，用以泛指所有基于数据平面的遥测方法，而不再仅限于最初提出的“经典 INT”模式——即在数据包中同时嵌入遥测指令与元数据的方式。

为此，依据数据包所经历的修改程度，特别是嵌入内容的不同，对 INT 的运行模式进行分类与定义。下文将详细介绍各类运行模式，并在[图 1](#fig1) 中进行汇总展示。

<p id="fig1" align="center">
    <img width="95%" src="/int/fig1.png" alt="Various modes of INT operation." />
    <span>图 1. INT 操作的各种模式</span>
</p>

### 3.1 INT 应用模式

在 INT 应用场景中，原始数据包将被监测，并可能被修改以携带遥测指令和元数据。根据对数据包修改的程度不同，INT 模式可分为以下三类：

- **INT-XD（eXport Data，导出数据）模式**：INT 节点根据其流监控表（Flow Watchlist）中配置的 INT 指令，从数据平面中直接导出元数据至遥测监控系统，无需对数据包本身进行任何修改。  
该模式在早期遥测报告规范中也被称为“明信片（Postcard）”模式，灵感来源于 [Handigol 等人](#handigol2014know)提出的机制。
    ::: info <span id="handigol2014know">明信片 postcards</span>
    Handigol, Nikhil, Brandon Heller, Vimalkumar Jeyakumar, David Mazières, and Nick McKeown. "I know what your packet did last hop: Using packet histories to troubleshoot networks." In 11th USENIX Symposium on Networked Systems Design and Implementation (NSDI 14), pp. 71-85. 2014.
    :::

- **INT-MX（eMbed instructions，嵌入指令）模式**：在此模式下，INT 源节点将遥测指令嵌入数据包头中。随后，INT 源节点、中继节点（Transit）以及接收节点（Sink）依据这些嵌入的指令，分别将采集到的元数据直接发送至遥测监控系统。接收端在转发前会移除指令报文头，以确保对后续处理的透明性。  
此模式下，对数据包的修改仅限于插入遥测指令头部，因此即便数据包经过多个中继节点，其长度也不会不断增长。  
INT-MX 亦支持“源端插入”元数据作为域特定指令（Domain-Specific Instructions）的一部分，使得源节点可嵌入额外元数据，供其他节点或监控系统使用。  
此模式的设计灵感来自 [IOAM 提出的“直接导出（Direct Export）”机制](#IOAM)。

    ::: info <span id="IOAM">IOAM</span>
    F. Brockners, S. Bhandari, and T. Mizrahi. 2022. RFC 9197: Data Fields for In Situ Operations, Administration, and Maintenance (IOAM). RFC Editor, USA.

    H. Song, B. Gafni, F. Brockners, S. Bhandari, and T. Mizrahi. 2022. RFC 9326: In Situ Operations, Administration, and Maintenance (IOAM) Direct Exporting. RFC Editor, USA.
    :::

- **INT-MD（eMbed Data，嵌入数据）模式**：该模式是经典的逐跳遥测方式，数据包中同时嵌入遥测指令与采集的元数据。其具体流程为：
  1. INT 源节点嵌入遥测指令；
  2. 源节点及各中继节点依照指令将本地元数据写入数据包中；
  3. INT 接收端移除数据包中的指令及聚合后的元数据，并视需要将结果发送至监控系统。

    在三种模式中，INT-MD 对数据包的修改程度最大，但由于所有遥测信息集中封装在报文中传递，因此监控系统无需从多个节点分别聚合遥测报告，从而大幅降低系统端的处理负载。  
    自 INT 规范 2.0 版本起，INT-MD 模式也引入了“仅源端插入”元数据的能力，作为域特定指令的一部分，使源节点可嵌入附加遥测数据，供接收端或监控系统解析使用。

::: warning 注意
除非特别说明，后续内容均默认采用 INT-MD 模式作为基础运行模式。
:::

### 3.2 INT 在合成流量中的应用

INT 源节点可以通过克隆原始数据包或生成特定的探测包（probe packet）来创建携带 INT 标记的合成流量（Synthetic Traffic）。对于这类合成流量，所有中继节点均应以与处理常规流量完全相同的方式应用 INT 操作。

合成流量与真实业务流量之间的主要区别在于：INT 接收节点在提取完路径中收集到的遥测数据后，通常需要丢弃合成流量，而非继续将其转发。为实现这一机制，INT 报文头中的 D 位（Discard Bit） 可被用作标识，表明对应数据包为克隆包或探测包，并应在 INT Sink 节点被丢弃。

所有 INT 模式均可应用于此类合成或探测数据包，具体模式由 INT 源节点决定。其中，若将 INT-MD（嵌入数据）模式 应用于探测包，则可实现与 [IFA](#IFA) 类似的遥测功能。

::: info <span id="IFA">IFA</span>
```bibtex
@techreport{kumar-ippm-ifa-08,
    number =    {draft-kumar-ippm-ifa-08},
    type =      {Internet-Draft},
    institution =   {Internet Engineering Task Force},
    publisher = {Internet Engineering Task Force},
    note =      {Work in Progress},
    url =       {https://datatracker.ietf.org/doc/draft-kumar-ippm-ifa/08/},
    title =     {{Inband Flow Analyzer}},
    pagetotal = 33,
    year =      2024,
    month =     apr,
    day =       26,
}
```
:::

通常而言，由克隆操作生成的合成流量在到达 INT Sink 后将被丢弃，而探测包则可依据具体用途被设置为转发或丢弃。正确标记这些数据包以明确接收端行为（转发或丢弃）是 INT 源节点的职责，确保 INT Sink 节点在完成遥测数据提取后作出正确处理。

## 4 监测内容的定义 {#sec4}

从理论上讲，INT 方法可用于定义并采集任意设备内部状态信息。然而在实际部署中，更具可行性和通用性的是确定一组精简的基础元数据集合，该集合可在多种网络设备上实现并广泛使用。本节所列出的元数据即构成了这样的基础集合。随着 INT 规范的不断演进，预计该集合将持续扩展，新增更多遥测字段以满足日益增长的需求。

需要注意的是，以下各类元数据的精确定义（例如时间戳的单位、跳延迟的计算方式、队列占用量或缓冲区使用率的量化方式等）可能因设备架构差异、功能集限制或资源约束等多种因素而有所不同。

相反，默认在具体部署过程中，各类 INT 设备的元数据语义应通过带外机制（Out-of-Band）与数据解释或分析组件进行协商与共享，从而确保遥测数据能被准确解析与应用。

### 4.1 设备级信息

- **节点 ID（Node ID）**
 
    表示 INT 节点的唯一标识符，通常由网络管理员进行配置。该标识在 INT 域内必须具备唯一性，用于在遥测报告中标明数据来源的设备节点。

### 4.2 入接口信息

- **入接口标识符（Ingress Interface Identifier）**

    表示 INT 数据包被接收时所在的接口。一个数据包可能通过一系列接口结构堆栈进入设备，该堆栈通常起始于物理端口。例如，数据包可能首先被接收于某个物理端口，该端口隶属于链路聚合组，该聚合组又作为三层交换虚拟接口的一部分，同时该数据包在三层还可能处于隧道封装状态。

    尽管理论上可以监测整个接口堆栈，但建议最多监测两个层级的入接口标识符：
    - 第一层级用于识别接收数据包的物理端口，占用 16 位（即一个 4 字节元数据字段的一半）；
    - 第二层级用于识别逻辑接口，占用完整的 4 字节字段，32 位空间足以支持每个网络设备上的大量逻辑接口。

    需要注意的是，不同设备的接口标识符语义可能存在差异，每个 INT 节点可自行决定在两个层级中选择何种类型的接口进行报告。

- **入接口时间戳（Ingress Timestamp）**  

    记录数据包在设备入方向（物理或逻辑接口）被接收时的本地时间。

### 4.3 出接口信息

- **出接口标识符（Egress Interface Identifier）**

    表示 INT 数据包发送时所使用的接口。数据包可能从一组接口结构堆栈中被转发出去，该堆栈最终连接至物理端口。例如，数据包可能通过隧道传出，经过三层交换虚拟接口，再经过链路聚合组中的某个物理端口。

    类似地，建议最多记录两个层级的出接口标识符：
    - 第一层级表示具体的物理端口，占用 16 位；
    - 第二层级用于逻辑接口，占用 32 位空间，以支持丰富的接口层次。

    每个 INT 节点可根据设备实际结构与配置，自主选择在每一层级中报告的接口类型。

- **出接口时间戳（Egress Timestamp）**

    表示数据包在设备出方向（物理或逻辑端口）被处理或发送的本地时间。

- **跳延迟（Hop Latency）**

    表示数据包在本设备内从接收至转发所经历的处理时延。

- **出接口链路利用率（Egress Interface TX Link Utilization）**

    表示数据包被发送时，出接口的实时带宽利用情况。设备可使用不同方法来估算此利用率，例如时间桶计数法或滑动平均算法等。尽管滑动平均效果优于时间桶计数法，INT 框架并不限定具体实现方式，由设备厂商自行决定。

- **队列占用量（Queue Occupancy）**

    表示数据包在转发过程中观察到的出队列内积压的流量量，可采用字节数、单元数或包数等不同计量单位。该元数据字段为 4 字节，其具体格式与单位依赖于设备实现。YANG 模型中的元数据语义定义应描述该字段的编码方式与度量单位。

- **缓冲区占用量（Buffer Occupancy）**

    表示数据包在转发过程中观察到的缓冲区中的积压流量情况，适用于多个队列共享缓冲区的场景。该字段也是 4 字节长度，格式与队列占用量类似，由实现决定具体含义，需由元数据语义模型进行说明。

目前，IETF 正在制定一套[基于 YANG 的元数据语义模型](#yang)，用于支持网络设备报告其所采集元数据的格式、单位及语义定义，从而提升 INT 报告的可解释性与互操作性。

::: info <span id="yang">YANG model</span>
p4-dtel-metadata-semantics, https://github.com/p4lang/p4-applications/blob/master/telemetry/code/models/p4-dtel-metadata-semantics.yang
:::

## 5 INT 头部格式

本节主要讲 INT 头部（INT Headers）的格式与位置，适用于 INT-MX 与 INT-MD 模式。这两种模式下，INT 指令（在 MD 模式下还包括元数据栈）被嵌入至数据包中，因此需要在数据平面对 INT 头部进行处理。

### 5.1 INT 头部类型 {#sec51}

INT 头部分为三种类型：**MD 类型**、**MX 类型**与 **Destination 类型**。每个 INT 数据包可能携带 MD 或 MX 类型中的一种头部，以及/或者 Destination 类型的头部。当同时存在 Destination 与 MD/MX 类型头部时，MD 或 MX 类型头部必须位于 Destination 类型头部之前。

- **MD 类型头部（INT Header Type 1）**

    此类型的头部需要由中间节点（即 INT Transit Hop）进行处理。其详细格式将在[第 5.8 节](#sec58)中讲解。

- **Destination 类型头部（INT Header Type 2）**

    此类型头部仅由 INT Sink 节点处理，中间节点必须忽略该头部内容。该类型头部用于支持 INT 源节点与终节点之间的边缘到边缘（Edge-to-Edge）通信。例如：
    - INT 源节点可添加序列号，以检测 INT 数据包是否丢失；
    - INT 源节点可写入原始 IP TTL 值与 INT 剩余跳数（Remaining Hop Count），INT Sink 可通过对比 TTL 与剩余跳数的递减值，检测路径中是否存在不支持 INT 的网络设备（假设每个设备都是 L3 跳点）。

    Destination 类型头部的具体格式将在未来版本的规范中定义。值得注意的是，某些边缘到边缘的 INT 用例也可通过 MD 或 MX 类型头部中的“仅源节点插入（source-only）”或“源节点嵌入（source-inserted）”的领域专用指令（Domain Specific Instructions）加以实现。

- **MX 类型头部（INT Header Type 3）**

    此类型头部也需由中间节点进行处理，节点应根据其中嵌入的指令将指定元数据报告至监控系统。MX 类型头部的详细格式将在[第 5.9 节](#59)中说明。

### 5.2 每跳节点的头部操作

#### 5.2.1 INT 源节点（INT Source Node）

在 INT-MD 和 INT-MX 模式下，位于数据包转发路径起点的 INT 源节点负责创建相应的 INT-MD 或 INT-MX 头部。

在 INT-MD 模式中，源节点在生成头部后，需附加自身的 INT 元数据。为防止因转发环路或其他异常导致的头部空间耗尽，强烈建议通过适当设置 INT 头部中的 **剩余跳数（Remaining Hop Count）** 字段，限制后续 Transit Hop 节点添加的总元数据字段数量。

INT-MD 与 INT-MX 头部的具体结构将在后续章节中详细描述。

#### 5.2.2 INT 中继节点（INT Transit Hop Node）

在 INT-MD 模式中，路径上的每个转发节点会**按需动态扩展** INT-MD 头部，以添加本节点的 INT 元数据。同样地，为避免由于环路或其他异常情况导致头部空间耗尽，每个 Transit Hop 节点必须**适当递减“剩余跳数”字段的值**。

在 INT-MX 模式中，每个 Transit Hop 节点依据 INT-MX 头部中嵌入的指令执行相应操作：收集本节点的特定元数据，并通过遥测报告（Telemetry Report）方式将其导出。

Transit Hop 节点可以修改 INT-MD 或 INT-MX 头部中的 DS Flags 字段，**但不得修改**以下字段：
- Hop ML（Hop Metadata Length）
- Instruction Bitmap
- Domain Specific ID
- DS Instruction

#### 5.2.3 INT 汇聚节点（INT Sink Node）

在 INT-MD 模式下，INT Sink 节点负责**剥离**数据包中的 INT 头部与元数据栈，并决定是否将收集到的信息上报。

在 INT-MX 模式下，INT Sink 节点**移除** INT-MX 头部，收集本地元数据，并决定是否对这些数据进行报告。

### 5.3 MTU 设置

在 INT-MX 和 INT-MD 模式下，INT 源节点插入 INT 头部时，可能会超出出口链路的最大传输单元（Maximum Transmission Unit，MTU）。

在 INT-MD 模式下，由于每个跳节点都需要为 INT 头部创建额外的空间以添加其元数据，数据包的大小会增加。这可能会导致在 INT 节点处超出出口链路的 MTU。

为解决这一问题，可以采取以下几种方式：

- **建议**：INT 源节点与 INT 汇聚节点之间的链路 MTU 应设置为高于前置链路（如服务器/虚拟机网卡的 MTU）一定值。可以通过配置如下的 MTU 差异值来防止因 INT 元数据插入而导致的 MTU 超出问题：
  - `[每跳元数据长度 * 4 * INT 跳数 + 固定 INT 头部长度]` 字节（对于 INT-MX 模式，只需考虑 [固定 INT 头部长度] 字节）。基于保守估计的 INT 跳数和每跳元数据长度配置该差值，能够有效防止因 INT 跳节点插入元数据而导致的出口 MTU 超出。固定 INT 头部长度为 INT 元数据头长度（12 字节）与封装特定的 shim/选项头部大小（4 字节）的和，具体讲解参见[第 5.7 节](#sec57)。
- **动态路径 MTU 发现**：INT 源节点或中继节点可以选择参与监控流的路径 MTU 动态发现过程。该过程通过向流量源发送 ICMP 消息，按照相应 L3 协议（IPv4 参照 RFC 1191，IPv6 参照 RFC 1981）的路径 MTU 发现机制进行。
  - 在这种情况下，INT 源节点或中继节点可能会报告一个保守的 MTU 值，假设数据包将通过最大允许的 INT 跳数（即剩余跳数递减至零），并考虑到所有 INT 跳节点上元数据的累积插入，假设所有下游 INT 跳节点的出口 MTU 与本节点的出口 MTU 相同。此方式有助于路径 MTU 发现源更快地收敛到一个路径 MTU 估算值，尽管该估算值较为保守。
  - 另外，每个 INT 跳节点也可以只报告其插入的元数据所需的 MTU，从而使路径 MTU 发现源精确计算路径 MTU，但此方法的代价是需要接收来自每个 INT 跳节点的更多 ICMP 消息。

无论是否参与路径 MTU 发现，如果 INT 中继节点由于插入所有请求的元数据会导致数据包长度超出出口链路 MTU，它必须：

- **不插入任何元数据**，并在 INT 头部设置 **M 位**，表示在 INT 跳节点处已超出出口 MTU，或
- **报告从前几个跳节点收集到的元数据堆栈**（如果生成的是 Telemetry Report 2.0 数据包，则设置 Intermediate Report 位），并从数据包中删除已报告的元数据堆栈，且将当前跳节点的元数据（包括该跳节点的元数据）包含在报告中或嵌入 INT-MD 元数据头部中。

INT 源节点插入 12 字节的固定 INT 头部，并且还可以插入每跳元数据长度 * 4 字节的其自有元数据。如果插入固定头部导致超出链路 MTU，则该数据包无法启动 INT。如果 INT 源节点被编程插入自有 INT 元数据，并且在数据包中有足够空间插入固定 INT 头部，但没有足够的空间插入其自有元数据，则源节点必须启动 INT 并设置 INT 头部的 M 位。

理论上，INT 中继节点可以通过 IPv4 分片来克服因插入元数据而导致的出口 MTU 限制。然而，IPv4 分片可能对应用程序产生不利影响。此外，IPv6 数据包在中继节点不能进行分片。另外，在 INT 中继节点处进行数据包分片（无论是否复制先前的 INT 元数据到分片中）会增加 INT 监控引擎对分片的关联复杂度。考虑到这些因素，建议 INT 节点在附加 INT 信息时**不对数据包进行分片**。

### 5.4 拥塞考量

INT 封装的使用不应加剧网络中的拥塞影响。虽然许多传输层协议（如 TCP、SCTP、DCCP、QUIC）本身具备拥塞控制机制，但另一些协议（如 UDP）不具备。在使用 UDP 的场景下，应用程序应自行实现拥塞控制机制或限制流量规模。

**建议不要**对已知**不具备拥塞控制机制的应用流量**启用 INT（详见 RFC 8085 第 3.1.11 节）。为此，应提供基于 IP 协议号与 L4 端口的**数据包过滤机制**（如访问控制列表 ACL），以便对特定协议流量执行匹配与筛选。

由于 INT 封装的源节点与汇聚节点通常处于同一管理域内，运营者可以选择允许对未具备拥塞控制机制的流量进行 INT 封装。在此情形下，运营者应仔细评估拥塞可能带来的影响，并采取适当的拥塞控制与缓解机制，包括容量规划、流量工程、速率限制等手段。

### 5.5 INT 可封装于任意协议 {#sec55}

INT 的规范文件**并未对 INT 头部的具体位置做出限制**：INT 头部可以作为任何封装类型的选项或有效载荷进行插入。唯一的要求是封装协议需提供足够的空间以携带 INT 信息，并且所有 INT 节点（源节点、中继节点和汇聚节点）需就 INT 头部的位置达成一致。

以下为若干在常见协议栈中可行的封装方式，部署方也可以根据其具体需求与环境选择更适合的封装格式：

- **INT over VXLAN**：作为 VXLAN 有效载荷插入，依 GPE 扩展规范；
- **INT over Geneve**：作为 Geneve 协议选项字段插入；
- **INT over NSH**：作为 NSH 协议的有效载荷插入；
- **INT over TCP**：作为 TCP 有效载荷插入；
- **INT over UDP**：作为 UDP 有效载荷插入；
- **INT over GRE**：作为 shim 层插入于 GRE 头部与其封装的有效载荷之间。

### 5.6 更新校验和

如[第 5.5 节](#sec55)所述，INT 头部与元数据可能承载于诸如 TCP 或 UDP 等 L4 协议，或包含 L4 头部的封装协议（如 VXLAN）中。当 INT 节点插入或移除 INT 头部与元数据，从而修改 L4 有效载荷时，**必须更新 TCP 或 UDP 的校验和字段**。但也存在一些例外情况：

- 对于运行在 IPv4 之上的 UDP，根据 [RFC 768]，可将校验和设置为 0，使接收端忽略该字段；
- 对于 IPv6 中的 UDP，某些特定应用场景亦允许设置为零校验和，详见 [RFC 6936]。

INT 的源节点、中继节点与接收节点在处理 L4 有效载荷时**必须遵循 IETF 的相关协议标准**。例如，若某节点收到带有零校验和的 UDP 流量，则不得修改该字段，以符合 RFC 768 与 RFC 6936 所定义的行为。

在需要进行校验和更新的场景中，INT 节点可采用以下两种方式之一：

- **直接更新**：修改 L4 校验和字段，使其匹配插入/删除 INT 元数据之后的新数据包；
- **使用校验和补偿（Checksum Complement）**：若 INT 源节点在指令位中设置了允许校验和中性更新（Checksum-neutral Update），则源节点或中继节点可计算并插入一个校验和补偿字段，使修改后的 L4 有效载荷在不变更原始 L4 校验和的前提下保持一致性。

使用校验和补偿的动因在于：某些硬件实现以流水线方式串行处理数据包，当 INT 字段位于 L4 校验和字段之后时，直接修改原始校验和将增加处理复杂度。为简化此过程，Checksum Complement 元数据通常为堆栈中的**最后一项元数据**。

若存在 Checksum Complement 元数据，INT 节点仍可以选择直接修改 L4 校验和字段。此时，Checksum Complement 字段应被赋予保留值 `0xFFFFFFFF`。对于执行 L4 校验的接收主机而言，无论是否使用了校验和补偿机制，最终的 L4 校验和均应保持正确性，不影响正常校验逻辑。

需要注意的是，**INT 汇聚节点（Sink）不可使用校验和补偿机制**，因为其会移除所有 INT 头部，必须通过直接更新 L4 校验和字段以完成最终修改。

无论是通过修改 L4 校验和字段，还是使用 Checksum Complement 元数据进行中性更新，建议采用**增量校验和算法**以检测潜在数据损坏。若节点需执行完整校验和重计算，则应先进行已有校验和的验证，以避免掩盖前一跳可能发生的损坏。

### 5.7 INT 头部位置 {#sec57}

INT 规范定义了四种典型的 INT 封装格式，以适配虚拟化与非虚拟化网络环境下的多种部署场景：

1. **INT over IPv4/GRE**

    INT 头部插入于 GRE 头部与其封装负载之间；

2. **INT over TCP/UDP**

    在 TCP 或 UDP 头部之后插入 shim 层头部，INT 头部位于 shim 层与原始 TCP/UDP 有效载荷之间。自 INT v2.0 起，亦支持另一种方式：在原有 L4 头部之前插入一个新的 UDP 头部，随后插入 INT 头部。该方式无需依赖隧道或虚拟化机制，适用于原生与虚拟化流量。

3. **INT over VXLAN**

    利用 VXLAN 的通用协议扩展（GPE），在 VXLAN 头部与其封装负载之间插入 INT 头部；

4. **INT over Geneve**

    Geneve 作为一种可扩展隧道协议，允许以 Geneve 选项字段的形式插入 INT 头部。

### 5.8 INT-MD 元数据头格式 {#sec58}

下面是 INT-MD（In-band Network Telemetry Metadata）元数据头的格式及其元数据堆栈结构。

```text
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|Ver = 2|D|E|M|       Reserved        | Hop ML  |RemainingHopCnt|
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|        Instruction Bitmap     |        Domain Specific ID     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           DS Instruction      |            DS Flags           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          INT 元数据堆栈（每跳插入 Hop ML * 4 字节）             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         ......                                |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                     最后一跳 INT 元数据                        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

- **INT-MD 元数据头** 总长为 12 字节，后跟可变长的元数据堆栈。每条元数据长度为 4 字节或 8 字节，每个 INT 跳点（除源节点包含的“源节点专用”元数据外）插入相同长度的元数据。由于每个数据包可能经历不同数量的跳点，元数据堆栈长度具有可变性。
- **Ver（4b）**：INT 元数据头的版本号，本版本固定为 2。
- **D（1b）**：丢弃位（Discard）。若设置，INT Sink 节点应在提取元数据后丢弃该数据包。
- **E（1b）**：最大跳数超限（Exceeded）。若 RemainingHopCnt 已减至 0，且该节点无法继续添加元数据，应设置此标志位。该位由 INT 源节点初始化为 0。
- **M（1b）**：MTU 超限（MTU Exceeded）。若某节点无法插入元数据，因插入后数据包长度将超过其出端口的 MTU，应设置该位，并不插入任何元数据。该标志用于诊断是否需要调整路径上的 MTU 设置。
- **Reserved（12b）**：保留位，INT 源节点应设置为 0，其他节点忽略。
- **Hop ML（5b）**：每跳元数据长度，以 4 字节为单位。该值由 INT 源节点设置，定义每个转发节点应插入的元数据长度。不包含源节点专用元数据部分。
- **RemainingHopCnt（8b）**：剩余可插入元数据的跳数。INT 源节点初始化该值，沿路径逐跳递减。若该值为 0，当前节点不得插入元数据，并应设置 E 位。
- **Instruction Bitmap（32b）**：指令位图，用于指示需要插入的标准元数据字段。每个置位的 bit 表示一个特定字段（如 Node ID、接口 ID、延迟等），并对应插入 4 字节或 8 字节的元数据。详见[第 4 节](#sec4)。
- **Domain Specific ID（16b）**：域标识符，用于区分不同的 INT 域。0x0000 表示默认域，DS Instruction 字段在该值下无效。其余值由运营者在 0x0001 到 0xFFFF 范围内分配。
- **DS Instruction（16b）**：域专用指令位图，用于指示插入的域专用元数据。某些位可被定义为“源节点专用”，仅由源节点插入。INT Transit 节点应根据 Hop ML 和已置位的指令比特插入等量元数据。
- **DS Flags（16b）**：域专用控制位，配合 DS Instruction 使用。
- **INT 元数据堆栈**：由各跳节点按照指令位图（Instruction Bitmap 和 DS Instruction）中的置位次序插入元数据。各跳必须将自己的元数据插入堆栈顶端，采用“前插”方式（即类似堆栈的 push 操作）。Checksum Complement 必须位于堆栈最后。
- 若某节点无法提供某项元数据，应插入全 1 保留值（0xFFFFFFFF 或 8 字节 0xFF），表示“无效”。若遇到保留位被置位（如与更新版本兼容），节点可插入保留值或不插入元数据。
- **若 DSID 不匹配本地已知域**，节点可选择：
  1. 插入保留值填充域专用元数据长度；
  2. 跳过 INT 处理，完全不插入元数据。
- **资源受限的 INT 节点** 可在无法处理全部指令时插入固定长度的元数据（Hop ML × 4 字节），但用保留值填充；若无法插入完整元数据段，应跳过插入。
- **未插入元数据的节点** 不得递减 RemainingHopCnt。

INT Sink 节点可选以下两种方式之一插入其本地遥测元数据：

1. 按照 Transit 节点方式插入堆栈中，元数据将出现在遥测报告（通常为数据包截断片段）中。
2. 插入至遥测报告的“可选基础元数据”或“可选域专用元数据”字段中，类似 INT-MX 节点（详见第 5.9 节）：
   - 若使用 Report Header，不适用 Checksum Complement，其位应清零；
   - 源节点专用元数据仍保留在截断堆栈中，而不出现在可选域专用字段。

::: warning 注意
INT 节点实现应仅支持上述两种方式之一，而遥测收集系统应支持解析二者。
:::

**字段使用总结：**

   - **INT 源节点** 需设置以下字段：
     - Ver、D、M、Hop ML、RemainingHopCnt、Instruction Bitmap；
     - 所有保留位应设置为 0；
     - 可选设置 Domain Specific 字段。
   - **中间转发节点（Transit）** 可设置：
     - E、M、RemainingHopCnt、DS Flags；
     - 不得修改 Hop ML、Instruction Bitmap、Domain Specific ID、DS Instruction。

**元数据堆栈长度要求：**

INT 元数据堆栈长度必须为 Hop ML × 4 字节 的整数倍，加上源节点专用元数据长度（若存在）。总长度可通过以下公式计算：

```text
INT 元数据总长度 = (INT Shim Header 长度 × 4) - 12 字节（INT-MD 固定头长度）
```

### 5.9 INT-MX 报文头格式 {#sec59}

INT-MX 报文头的格式如下所示，该报文头总长度为 12 字节，用于指示需要采集的标准和特定领域的元数据项：

```text
  0                   1                   2                   3
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |Ver = 2|D|                  Reserved                           |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |      Instruction Bitmap       |    Domain Specific ID (16b)   |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |      DS Instruction (16b)     |        DS Flags               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |         可选的特定领域 Source-Inserted Metadata（变长）         |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

**字段说明：**

- **Ver（4 位）**：INT-MX 报文头的版本号，当前应设为 2。
- **D（1 位）**：丢弃标志位。若设置为 1，表示 INT Sink 节点在发送遥测报告后应丢弃该数据包。
- **Reserved（27 位）**：保留位，应由 INT Source 设置为 0，其他节点忽略。
- **Instruction Bitmap（32 位）**：指示标准 INT 元数据项的位图，每一位对应一种标准元数据（详见[第 4 节](#sec4)），如：
  - bit0：节点 ID
  - bit1：一级入口/出口端口 ID
  - bit2：跳数延迟
  - bit3：队列 ID + 队列占用率
  - bit4/5：入/出时间戳（各占 8 字节）
  - bit6：二级端口 ID（入口 + 出口）
  - bit7：出口端口利用率
  - bit8：Buffer ID + Buffer 占用率
  - 其余位保留。位 0~14 为基础 INT 指令，每设置一位对应采集 4 或 8 字节的元数据。
- **Domain Specific ID（16 位）**：表示 INT 域的唯一标识。若与当前节点所识别的域 ID 匹配，则继续解析下方的 DS Instruction 和 DS Flags 字段。默认值为 0x0000，表示无特定域信息。
- **DS Instruction（16 位）**：为特定 INT 域定制的元数据指令位图，支持以下两种模式：
  - **Export 模式**：节点根据设置的位采集并上报相应的特定域元数据；
  - **Source-Inserted 模式**：源节点在数据包中插入对应的元数据，其使用要求由下列策略控制：
    - **All Nodes**：所有节点均需上报；
    - **Sink Node**：仅由终节点上报；
    - **None**：仅供网络中转节点消费，无需上报。
  - **Source-Inserted Metadata Mutability（可变性）**：
    - **Source-Only**：仅源节点可设置，中间和终端节点不得更改；
    - **Cumulative**：中转和终端节点可更新或覆盖其值。
- **Optional Domain Specific Source-Inserted Metadata**：即“源插入元数据”，位于 DS Instruction 与 DS Flags 之后，长度必须为 4 字节的整数倍。

**报文处理逻辑：**

对于支持 INT-MX 的路径节点（源节点、转发节点和接收节点），每个节点会根据 INT-MX 头部中的 **Instruction Bitmap** 和 **DS Instruction** 字段，生成并上报自身的元数据。处理流程如下：

1. **复制并修改指令位图**：

    每个节点将 INT-MX 头部中的 `Instruction Bitmap` 和 `DS Instruction` 字段复制到遥测报告中的 `RepMdBits` 和 `DSMdBits`。之后根据以下规则对它们进行必要的修改。

2. **指令位图的特殊处理规则**：

   - 如果 `Instruction Bitmap` 中设置了 bit 0（表示 Node ID），则需要清除 `RepMdBits` 中的该位，因为 Node ID 已包含在遥测报告通用头中。
   - 如果节点无法提供某些指定的元数据（例如元数据不可用或对应位为保留位），必须清除遥测报告中对应的 `RepMdBits` 或 `DSMdBits`。

3. **域匹配与处理逻辑**：

   - 若 `Domain Specific ID` 不属于该节点已知的域，节点需执行以下两种策略之一：
     - 仅发送与 `Instruction Bitmap` 对应的标准元数据，并确保遥测报告中的 `DSMdBits` 全部清零；
     - 完全不向监控系统发送任何元数据。

4. **处理 Source-Inserted Metadata（源插入元数据）**：

   - 如果该元数据要求 "All Nodes" 或 "Sink Node" 上报，INT 节点应尽力将其包含进遥测报告，方法有两种：
     1. 将原始 INT-MX 头部（含源插入元数据）嵌入遥测报告的 Individual Report 内部内容中；
     2. 将该元数据解析后填入遥测报告的 Variable Optional Domain Specific Metadata 字段中（需了解对应 `Domain Specific ID` 和 `DSMdBits` 的定义，以确保元数据顺序正确）。

    ::: warning 注意
    若未采用方法 2，必须清除 `DSMdBits` 中对应的位。两个方法可并存，但每个节点通常只使用其中之一。

    出于兼容性考虑，域管理员在未来定义新的 "source-inserted" 元数据时，应谨慎处理保留位的使用。
    :::

5. **资源限制处理**：
   - 若节点受限于指令处理能力或单报文可收集的元数据长度，应尽量发送可处理的部分，并更新 `RepMdBits` 和 `DSMdBits` 表示实际包含的元数据内容。

6. **元数据长度字段计算**：
   - 报告中 `MD Length` 字段的值应基于实际设置的 `RepMdBits` 和 `DSMdBits` 计算。
   - `Variable Optional Baseline Metadata` 与 `Variable Optional Domain Specific Metadata` 字段内容同样由这两个位图决定。

**字段使用小结：**

- **INT 源节点必须设置**：
  - `Ver`, `D`, `Instruction Bitmap` 字段；
  - 所有保留位应置零；
  - 可选择性设置域特定字段（如 `Domain Specific ID`、`DS Instruction`、`Source-Inserted Metadata`）。
- **中间转发节点可设置**：
  - 仅可设置 `DS Flags` 字段；
  - **不得修改** `Instruction Bitmap`、`Domain Specific ID` 和 `DS Instruction` 字段。

INT 元数据总长度必须是 **4 字节的倍数**，否则不符合格式规范。
