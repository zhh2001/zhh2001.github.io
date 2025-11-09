---
outline: deep
---

# P4-Utils

P4-Utils 是一个可以创建和测试包含 P4 交换机的虚拟网络的 Python 包。

## 安装

如果已经安装了所有必要的[依赖项](#依赖项)，可以使用以下命令安装 P4-Utils：

<<< @/sdn/codes/p4utils/install.sh

### 依赖项

1. [PI 库](https://github.com/p4lang/PI)：提供一个 P4Runtime 服务器的实现框架，仅在拓扑中使用 P4Runtime 交换机时需要。
2. [行为模型（bmv2）](https://github.com/p4lang/behavioral-model)：包含行为模型的多种软件实现（例如，`simple_switch` 和 `simple_switch_grpc`）。
3. [p4c](https://github.com/p4lang/p4c)：P4 语言的编译器，支持 **P4_14** 和 **P4_16**。
4. [Mininet](https://github.com/mininet/mininet)：允许在单台机器上创建一个真实的虚拟网络，运行实际的内核、交换机和应用代码。
5. [FRRouting](https://github.com/FRRouting/FRR)：一个用于 Linux 和 Unix 平台的免费开源 Internet 路由协议套件，实现了 BGP、OSPF、RIP、IS-IS、PIM、LDP、BFD、Babel、PBR、OpenFabric 和 VRRP，并对 EIGRP 和 NHRP 提供 alpha 支持。P4-Utils 中的路由器节点基于 FRRouting，仅在拓扑中包含路由器时需要。

由于 P4-Utils 所需的依赖项较多，手动安装过程相当漫长且繁琐。因此，P4-Utils 官方提供了一个 [Bash 脚本](#自动安装)，可自动执行每一步。

::: warning 注意
该脚本已在 **Ubuntu 20.04 和 Ubuntu 22.04** 以及编译器 **GCC 9.4** 上进行了测试。
:::

::: tip 重要提示
通过以下安装方法，可以下载并安装 _Mininet_ 和 P4-Tools 套件（包括 P4-Utils、P4-Learning 及其依赖项）到用户主目录中。
:::

### 自动安装

<<< @/sdn/codes/p4utils/install_auto.sh

## 用法

这个章节将简要概述 P4-Utils 的基本功能。

为了使说明更简单和具体，后面将使用以下网络示例。将逐步讲解配置文件，这些文件允许用户定义和使用这样的拓扑。

![l2 topology](/p4/l2_topology.png)

如上图所示，我们有四个主机连接到一个交换机。所有交换机都位于相同的子网 `10.0.0.0/16` 中。我们需要创建网络，并通过 L2 转发在主机之间建立连接。

### 网络创建

让我们创建一个名为 `network.py` 的文件。为了定义网络，首先需要导入所需的模块并创建一个 `NetworkAPI` 对象：

<<< @/sdn/codes/p4utils/network.py

我们还可以设置脚本执行期间显示日志的详细级别：

<<< @/sdn/codes/p4utils/log_level.py

另一个重要的选项是涉及主机 ARP 表的设置。可以通过使用 `disableArpTables()` 和 `disableGwArp()` 方法禁用相同子网内主机及其网关的静态 ARP 条目。这些选项不适用于我们的简单示例。

::: tip 重要提示
默认情况下，网络中主机的 ARP 表会在网络启动时以静态方式填充。这样，在操作网络时无需考虑 ARP 请求。
:::

可能的 `logLevel` 值如下（按详细程度递减顺序）：

- `debug`
- `info`
- `output`
- `warning`
- `error`
- `critical`

现在我们可以定义拓扑了。首先添加节点：

<<< @/sdn/codes/p4utils/node.py

如上所示，我们添加了一个名为 `s1` 的 P4 交换机和四个主机 `h1`、`h2`、`h3`、`h4`。

::: warning 注意
在添加节点时，请确保每个节点名称唯一。
:::

对于 P4 交换机，需要配置一个 P4 程序。假设我们有一个名为 `l2_forwarding.p4` 的 P4 程序文件，与 Python 脚本位于同一文件夹中。通过以下命令将其添加到 `s1`：

<<< @/sdn/codes/p4utils/p4_source.py

此文件将被编译并传递给交换机。

现在设置链路：

<<< @/sdn/codes/p4utils/link.py

::: warning 警告
链路必须在节点添加之后设置，因为 `addLink()` 方法会检查连接的节点是否实际存在于网络中。
:::

可以通过指定节点的端口号简化交换机配置：

<<< @/sdn/codes/p4utils/intf_port.py

::: tip 重要提示
如果未指定端口号，将进行自动分配。自动分配在不同网络脚本执行中是一致的。
:::

如果要将 `s1` 和 `h1` 之间的链路带宽限制为 5 Mbps，可以使用以下方法：

<<< @/sdn/codes/p4utils/bw.py

要一次性为所有链路设置带宽：

<<< @/sdn/codes/p4utils/bw_all.py

定义拓扑后，需要为节点分配 IP 和 MAC 地址。提供了以下三种方法：

- 默认设置：如果未指定，所有节点都位于 `10.0.0.0/8` 网络中，MAC 地址随机分配。
- 手动分配：
    - 使用 `setIntfIp()` 设置接口的 IP 地址：

        <<< @/sdn/codes/p4utils/intf_ip.py

    - 使用 `setIntfMac()` 设置接口的 MAC 地址：

        <<< @/sdn/codes/p4utils/intf_mac.py

- 预定义分配策略：
    - 使用 L2 策略：

        <<< @/sdn/codes/p4utils/l2.py

    - 使用混合策略：

        <<< @/sdn/codes/p4utils/mixed.py

    - 使用 L3 策略：

        <<< @/sdn/codes/p4utils/l3.py

在本例中，主机位于同一网络，因此可以使用 L2 策略。

可以启用 `.pcap` 文件记录和日志记录：

<<< @/sdn/codes/p4utils/pcap.py

::: warning 注意
也可以使用 `enablePcapDump()` 和 `enableLog()` 为特定交换机启用这些功能。
:::

最后，我们可以启用网络客户端并启动网络：

<<< @/sdn/codes/p4utils/start.py

要执行网络，只需以超级用户权限运行我们的 Python 脚本：

<<< @/sdn/codes/p4utils/start.sh

### 自动分配策略

为网络中每个接口指定地址可能非常繁琐。为此，可以使用自动分配策略，这些策略遵循简单规则，在大多数情况下非常实用。

::: warning 注意
以下所有策略都假设：
- 每个主机只连接到一个交换机。
- 仅允许交换机和主机。
- 不允许并行链路。
:::

#### L2 策略

**L2** 策略将所有设备放置在同一个 IPv4 网络（`10.0.0.0/16`）中。可以通过调用 `l2()` 方法实现。

如果主机和交换机的命名分别为 `h<ID>` 和 `s<ID>`（例如 `h1`、`h2`、`s1`、`s2`、...），IP 地址的分配规则如下：

- **主机 IP**: `10.0.x.y/16`，其中 `x` 和 `y` 分别是主机 ID 的高字节和低字节（即其二进制表示形式的上半部分和下半部分）。

#### Mixed 策略

**Mixed** 策略将连接到同一交换机的主机放置在同一子网中，而不同的交换机（即使它们之间有链路连接）则放置在不同的子网中。可以通过调用 `mixed()` 方法实现。

如果主机和交换机的命名分别为 `h<ID>` 和 `s<ID>`（例如 `h1`、`h2`、`s1`、`s2`、...），IP 地址的分配规则如下：

- **主机 IP**: `10.x.y.z/24`，其中 `x` 和 `y` 分别是网关交换机 ID 的高字节和低字节（即其二进制表示形式的上半部分和下半部分），`z` 是主机 ID。
- **直接连接到主机的交换机端口 IP**: `10.x.y.254/24`，其中 `x` 和 `y` 分别是网关交换机 ID 的高字节和低字节（即其二进制表示形式）。
- **交换机之间的接口 IP**: `20.sw1.sw2.<1,2>/24`，其中 `sw1` 是第一个交换机的 ID（按照链路定义中的顺序），`sw2` 是第二个交换机的 ID。最后一个字节为 `1` 表示 `sw1` 的接口，`2` 表示 `sw2` 的接口。

#### L3 策略

**L3** 策略将每个主机放置在一个独立的子网中，该子网与连接到主机的交换机端口共享 IP 地址。可以通过调用 `l3()` 方法实现。

如果主机和交换机的命名分别为 `h<ID>` 和 `s<ID>`（例如 `h1`、`h2`、`s1`、`s2`、...），IP 地址的分配规则如下：

- **主机 IP**: `10.x.y.2/24`，其中 `x` 是网关交换机的 ID，`y` 是主机的 ID。
- **直接连接到主机的交换机端口 IP**: `10.x.y.1/24`，其中 `x` 是网关交换机的 ID，`y` 是主机的 ID。
- **交换机之间的接口 IP**: `20.sw1.sw2.<1,2>/24`，其中 `sw1` 是第一个交换机的 ID（按照链路定义中的顺序），`sw2` 是第二个交换机的 ID。最后一个字节为 `1` 表示 `sw1` 的接口，`2` 表示 `sw2` 的接口。

### 网络客户端

网络客户端由 `p4utils.mininetlib.cli.P4CLI` 实现。如果启用了客户端，它会在网络启动后立即运行。

例如，如果想快速检查所有主机之间的连通性，可以使用以下命令：

<<< @/sdn/codes/p4utils/pingall.sh

还可以通过以下命令获取命令的摘要：

<<< @/sdn/codes/p4utils/help.sh

### 控制平面配置

在拥有已配置的 P4 交换机的工作拓扑后，需要通过转发信息填充数据平面以建立连接。这可以通过 Python 脚本以编程方式完成，也可以通过 Thrift 客户端以静态方式完成。第一种方法将在后面介绍，以下是第二种方法的说明。

#### Thrift 客户端

要启动 _Thrift_ 客户端并连接到 `s1`，需要知道其 _Thrift_ 服务器的 IP 和端口。在本例中，IP 是 `127.0.0.1`，端口是 `9090`。

::: warning 注意
_Thrift_ 端口号可以在网络配置中显式分配。如果未为网络中的任何 P4 交换机指定端口，则按顺序自动分配，第一台 P4 交换机的端口为 `9090`。
:::

启动网络后，可以执行以下命令连接到客户端：

<<< @/sdn/codes/p4utils/thrift.sh

以下选项可以传递给 `simple_switch_CLI`：

<<< @/sdn/codes/p4utils/thrift_opt.sh

::: tip 重要提示
如果未指定，`simple_switch_CLI` 命令会假设 IP 为 `127.0.0.1`，端口为 `9090`。P4-Utils 始终将 IP `127.0.0.1` 分配给所有 _Thrift_ 服务器，因此唯一变化的是每个交换机监听的端口。
:::

可以通过以下命令获取所有可用命令列表：

<<< @/sdn/codes/p4utils/thrift_question.sh

要检查特定命令的语法，可以使用以下方式：

<<< @/sdn/codes/p4utils/thrift_help.sh

假设我们要填充示例中交换机 `s1` 的转发表，可以运行以下命令：

<<< @/sdn/codes/p4utils/cmd_table_add.sh

每条命令会向表中添加一个匹配项，例如：

```text
Adding entry to exact match table dmac
match key:           EXACT-00:00:0a:00:00:01
action:              forward
runtime data:        00:01
Entry has been added with handle 0
```

#### 命令文件

为每个交换机逐条输入命令可能较为繁琐，因此可以将所有命令写入一个 `.txt` 文件，每行一条命令，然后将该文件传递给 P4-Utils，后者会使用 `simple_switch_CLI` 自动执行。

例如，将 `s1` 客户端中输入的命令写入一个名为 `s1-commands.txt` 的文件：

```txt
table_add dmac forward 00:00:0a:00:00:01 => 1
table_add dmac forward 00:00:0a:00:00:02 => 2
table_add dmac forward 00:00:0a:00:00:03 => 3
table_add dmac forward 00:00:0a:00:00:04 => 4
```

在 Python 网络配置脚本中添加以下行，将该命令文件传递给 P4-Utils（假设脚本和命令文件在同一文件夹中）：

<<< @/sdn/codes/p4utils/cli_input.py

如果使用 JSON 网络配置文件，可以通过修改 `topology` 字段中的交换机 `s1` 来指定 _Thrift_ 命令文件：

```json
"s1": {"cli_input": "s1-commands.txt"}
```

网络启动后，交换机会自动配置。

## 高级用法

在本章节中，我们将探讨一些以编程方式控制交换机和在网络节点上调度任务的方法。

让我们再次回顾[用法章节](#用法)中已经介绍的示例。假设我们已经按照说明创建了网络，并且可以启动它。接下来，我们需要配置交换机 `s1`，使其能够执行 L2 转发。

![l2 topology](/p4/l2_topology.png)

### 控制平面配置

在以下部分中，我们探讨两种替代 [_Thrift_ 命令行客户端的方法](#thrift-客户端)来控制交换机。

#### Thrift API

_Thrift API_ 可以用于所有 P4 交换机，基于 _Thrift_ 命令行客户端的代码实现。实际上，它提供了相同的方法来控制交换机。_Thrift API_ 的实现依赖于 `SimpleSwitchThriftAPI`。

开始使用，我们创建一个新的 Python 脚本，命名为 `controller.py`，并导入用于配置 P4 交换机的模块：

<<< @/sdn/codes/p4utils/thrift.py

我们还需要与运行在交换机上的服务器建立连接。我们知道 `s1` 的 _Thrift_ 服务器监听地址为 `127.0.0.1:9090`（见[相关说明](#thrift-客户端)），因此可以通过以下方式连接：

<<< @/sdn/codes/p4utils/thrift_conn.py

::: warning 注意
`SimpleSwitchThriftAPI` 默认假设 _Thrift_ 服务器的 IP 地址为 `127.0.0.1`。
:::

现在，我们可以使用控制器设置转发规则。调用 `table_add()` 方法：

<<< @/sdn/codes/p4utils/table_add.py

在网络启动后，可以运行控制器脚本来填充 `s1` 的转发表：

<<< @/sdn/codes/p4utils/controller.sh

#### P4Runtime API

_P4Runtime API_ 由 `SimpleSwitchP4RuntimeAPI` 实现，仅适用于支持 P4Runtime 的交换机。并非行为模型（behavioral-model）提供的所有二进制文件都支持此功能。因此，只有 `P4RuntimeSwitch` 支持 P4Runtime，而 `P4Switch` 不支持。

::: warning 警告
行为模型必须与 P4Runtime 依赖项一起构建才能使其正常工作。否则，仅支持非 P4Runtime 的目标设备。
:::

假设我们已正确启用了 P4Runtime，可以为交换机 `s1` 编写 Python P4Runtime 控制器脚本（`controller.py`）：

<<< @/sdn/codes/p4utils/p4runtime.py

::: tip 重要提示
P4Info 文件 `l2_forwarding_p4rt.txt` 和 P4 编译后的 JSON 文件 `l2_forwarding.json` 都是由 P4 编译器生成的。
:::

::: warning 注意
- `device_id`：
    - 可以在网络配置中显式分配给交换机。
    - 如果未显式分配，则按顺序自动分配，第一个交换机的 `device_id` 为 `1`。
- `grpc_port`：
    - 可以在网络配置中显式分配给 P4Runtime 交换机。
    - 如果未显式分配，则按顺序自动分配，第一个 P4Runtime 交换机的 `grpc_port` 为 `9559`。
:::

现在， 我们可以使用控制器设置转发规则，调用 `table_add()` 方法：

<<< @/sdn/codes/p4utils/table_add.py

在网络启动后，可以运行控制器脚本来填充 `s1` 的转发表：

<<< @/sdn/codes/p4utils/controller.sh

##### 使用 Python 启用 P4Runtime

针对我们的简单示例，P4 编译器需要知道我们正在使用 P4Runtime 交换机，以生成连接到交换机的 P4Runtime 服务器所需的 P4Info 文件。

<<< @/sdn/codes/p4utils/p4rt.py

接下来，需要指定正在使用 P4Runtime 交换机。在 Python 网络配置脚本中，可以通过以下方式实现：

<<< @/sdn/codes/p4utils/p4sw.py

##### 使用 JSON 启用 P4Runtime

对于 JSON 配置文件，可以通过在 `compiler_module` 选项中指定生成 P4Info 文件来告知 P4 编译器使用 P4Runtime。此外，可以使用 `switch_node` 选项设置 P4 交换机的默认类型为 `P4RuntimeSwitch`，从而使每个交换机都支持 P4Runtime。

将上述设置应用于简单网络示例的 JSON 配置文件后，结果如下：

```json
{
  "p4_src": "l2_forwarding.p4",
  "cli": true,
  "pcap_dump": true,
  "enable_log": true,
  "compiler_module": {
    "options": {
      "p4rt": true
    }
  },
  "switch_node": {
    "module_name": "p4utils.mininetlib.node",
    "object_name": "P4RuntimeSwitch"
  },
  "topology": {
    "assignment_strategy": "l2",
    "default": {
      "bw": 10
    },
    "links": [["h1", "s1"], ["h2", "s1"], ["h3", "s1"], ["h4", "s1"]],
    "hosts": {
      "h1": {},
      "h2": {},
      "h3": {},
      "h4": {}
    },
    "switches": {
      "s1": {}
    }
  }
}
```

### 拓扑数据库

到目前为止，我们已经看到了一些控制交换机的方法，但这些方法都依赖于用户提供的信息：用户需要知道所有的网络地址、端口号等。这在小型拓扑中是可行的，但在需要处理几十甚至上百个地址和端口号的大型拓扑中会变得非常困难。

为了解决这个问题，P4-Utils 提供了一个内置的拓扑数据库。在网络启动后，该数据库会自动生成并保存为一个 JSON 文件（通常命名为 `topology.json`），位于执行目录下。用户可以查询此文件以检索拓扑信息。此功能由 `NetworkGraph` 实现。

以我们简单的拓扑示例为例，可以在无需了解交换机任何详细信息的情况下，自动配置转发表：

<<< @/sdn/codes/p4utils/topo.py

### 任务调度器

任务调度器允许用户轻松地在不同节点上调度不同的任务（例如生成流量）。可以通过以下两种方式访问任务调度器：

- 在[网络客户端](#网络客户端)中添加任务。
- 将任务写入一个 `.txt` 文件（每行一个任务），由 P4-Utils 解析。

以下是一些简单示例。

#### 使用网络客户端调度任务

在网络启动后，可以在网络客户端中使用以下命令调度任务：

<<< @/sdn/codes/p4utils/task.sh

参数说明：

- `node`：节点名称。
- `exe`：要运行的可执行程序（可以是一个 shell 命令字符串，也可以是 Python 函数的名称）。
- `argX`：传递给函数的定位参数（可选）。
- `start`：任务延迟时间（以秒为单位，相对于当前时间）。
- `duration`：任务持续时间（以秒为单位，如果小于或等于 `0`，则任务没有时间限制）。
- `keyX` 和 `kwargX`：传递给函数的关键字参数（可选）。

::: tip 重要提示
默认查找函数的模块是 `p4utils.utils.traffic_utils`。可以使用 `--mod <module>` 在命令中指定其他模块。
:::

假设我们使用简单的网络拓扑，让 `h1` 对 `h2` 进行 10 秒的 ping，可以在客户端中输入以下命令：

<<< @/sdn/codes/p4utils/task_ping.sh

#### 使用文件调度任务

当需要处理多个任务时，使用一个文件集中管理所有任务比直接使用客户端更加方便。此文件会在网络启动时被解析，并在网络启动后立即分发任务。

::: warning 注意
对于任务文件，`start` 值表示任务相对于网络启动时间的延迟。
:::

任务文件的语法与网络客户端中的命令基本相同，只有一个区别：**不需要在任务行的开头添加 `task` 命令**。每行定义一个任务。

假设我们有一个 L2 转发示例，需要以下任务：

- `h1` 在网络启动 30 秒后对 `h2` 进行 10 秒的 ping。
- `h3` 在网络启动 10 秒后对 `h4` 进行 30 秒的 ping。

我们可以创建一个名为 `tasks.txt` 的文件，其内容如下：

```text
h1 30 10 "ping 10.0.0.2"
h3 10 30 "ping 10.0.0.4"
```

现在，我们需要将任务文件传递给 P4-Utils 框架。

如果使用的是 Python 网络配置脚本，可以通过添加以下代码来实现：

<<< @/sdn/codes/p4utils/task.py

如果使用的是 JSON 配置文件，可以在主字典中添加以下键值对（例如，可以将其放在 `p4_src` 选项之后）：

```json
"tasks_file": "tasks.txt"
```

现在一切准备就绪。网络启动后，每个任务都会根据任务文件中的信息自动调度并执行。
