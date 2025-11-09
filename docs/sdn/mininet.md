---
outline: deep
---

# Mininet

Mininet 是一种网络仿真工具，可以创建包含虚拟主机、交换机、控制器和链路的网络。Mininet 主机运行标准的 Linux 网络软件，其交换机支持 OpenFlow，实现高度灵活的定制路由和软件定义网络（SDN）。

## 安装

准备好 Ubuntu 或者 Debian 操作系统，然后获取 Mininet 的源码：

```shell
git clone https://github.com/mininet/mininet
```

上述 `git` 命令将获取最新版本的 Mininet。如果想切换其他特定版本，可以明确检出该版本：

```shell
cd mininet
git tag  # 列出可用版本
git checkout -b mininet-2.3.0 2.3.0  # 希望安装的版本
cd ..
```

安装命令为：

```shell
mininet/util/install.sh [选项]
```

`install.sh` 脚本的常用选项包括：

- `-a`：安装所有内容，包括 Open vSwitch 等依赖项，以及 OpenFlow Wireshark 解码器和 POX 等附加工具。默认情况下，这些工具将安装在主目录中创建的文件夹内。
- `-nfv`：安装 Mininet、OpenFlow 参考交换机和 Open vSwitch。
- `-s mydir`：在其他选项之前使用此选项，将源码放置在指定目录 mydir 中，而不是放在主目录中。

可以通过以下命令了解其他有用的选项（例如安装 OpenFlow Wireshark 解码器，如果安装的 Wireshark 版本中尚未包含该解码器）：

```shell
install.sh -h
```

安装完成后，测试 Mininet 的基本功能：

```shell
sudo mn --switch ovsbr --test pingall
```

## 入门

### 1 Mininet 日常使用

#### 1.1 查看帮助信息

```shell
sudo mn -h
```

#### 1.2 启动 Wireshark

```shell
sudo wireshark &
```

##### 1.2.1 如果未安装 Wireshark（显示“找不到命令”错误）

如果没有安装 Wireshark 和 OpenFlow 插件，可以尝试使用 Mininet 的 `install.sh` 脚本来安装它们，如下所示：

```shell
cd ~
git clone https://github.com/mininet/mininet  # 如果之前没有克隆
mininet/util/install.sh -w
```

#### 1.3 与主机和交换机交互

启动一个最小拓扑并进入 Mininet 客户端：

```shell
sudo mn
```

默认拓扑是 `minimal` 拓扑，包括一个连接两个主机的 OpenFlow 内核交换机和一个 OpenFlow 参考控制器。此拓扑也可以在命令行中使用 `--topo=minimal` 来指定。Mininet 还提供其他现成的拓扑，可以查看 `mn -h` 输出中的 `--topo` 部分。

此时，虚拟机中运行着四个实体：2 个主机进程、1 个交换机进程、1 个基本控制器。

同时，会进入 Mininet 命令行。

在 Wireshark 窗口中，可以看到内核交换机连接到参考控制器。

显示 Mininet CLI 命令：

```shell
mininet> help
```

显示节点：

```shell
mininet> nodes
```

显示链接：

```shell
mininet> net
```

转储所有节点的信息：

```shell
mininet> dump
```

此时应该会看到列出的交换机和两个主机。

如果在 Mininet CLI 中输入的第一个字符串是主机、交换机或控制器的名称，则该命令将在该节点上执行。在主机进程上运行命令：

```shell
mininet> h1 ifconfig -a
```

此时应该会看到主机的 `h1-eth0` 和回环（`lo`）接口。请注意，当在主系统中运行 `ifconfig` 时，无法看到此接口（`h1-eth0`），因为它是特定于主机进程的网络命名空间。

相反，交换机默认运行在根网络命名空间中，因此在“交换机”上运行命令与在常规终端中运行命令相同：

```shell
mininet> s1 ifconfig -a
```

这将显示交换机接口以及虚拟机的外部连接（`eth0`）。

#### 1.4 测试主机之间的连通性

验证是否可以从主机 `h1` ping 到主机 `h2`：

```shell
mininet> h1 ping -c 1 h2
```

另一种更简便的测试方法是使用 Mininet CLI 内置的 `pingall` 命令，它会对所有节点成对 `ping` 测试：

```shell
mininet> pingall
```

#### 1.5 运行一个简单的 Web 服务器和客户端

`ping` 并不是可以在主机上运行的唯一命令。Mininet 主机可以运行任何在底层 Linux 系统（或虚拟机）及其文件系统中可用的命令或应用程序，还可以输入 `bash` 命令，包括作业控制（`&`、`jobs`、`kill` 等）。

接下来，尝试在 `h1` 上启动一个简单的 HTTP 服务器，从 `h2` 发起请求，然后关闭该 Web 服务器：

```shell
mininet> h1 python -m http.server 80 &
mininet> h2 wget -O - h1
...
mininet> h1 kill %python
```

::: warning 注意
对于 Python 3，HTTP 服务器命令是 `http.server`；对于 Python 2，它是 `SimpleHTTPServer`。请确保使用与运行的 Mininet 版本对应的正确命令。要确定 Mininet 使用的 Python 版本，可以输入以下命令：
```shell
mininet> py sys.version
'3.12.3 (main, Sep 11 2024, 14:17:37) [GCC 13.2.0]'
```
:::

退出 CLI：

```shell
mininet> exit
```

#### 1.6 清理环境

如果 Mininet 因某种原因崩溃，执行清理操作可能会解决问题：

```shell
sudo mn -c
```

### 2 高级启动选项

#### 2.1 回归测试

例如：

```shell
sudo mn --test pingpair
```

此命令创建了一个最小拓扑，启动了 OpenFlow 参考控制器，运行了全对 ping 测试，然后销毁了拓扑和控制器。

另一个有用的测试是 `iperf`（大约需要 10 秒完成）：

```shell
sudo mn --test iperf
```

此命令创建了相同的 Mininet，分别在一台主机上运行 iperf 服务器，在另一台主机上运行 iperf 客户端。

#### 2.2 拓扑大小和类型

默认拓扑是一个连接两个主机的单交换机。可以使用 `--topo` 更改为不同的拓扑，并为该拓扑的创建传递参数。例如，要验证一个交换机和三个主机之间的全对 ping 连通性：

运行回归测试：

```shell
sudo mn --test pingall --topo single,3
```

另一个示例是线性拓扑（每个交换机有一个主机，所有交换机连接成一条直线）：

```shell
sudo mn --test pingall --topo linear,4
```

#### 2.3 链路参数

Mininet 2.0 允许设置链路参数，这些参数甚至可以通过命令行自动设置：

```shell
sudo mn --link tc,bw=10,delay=10ms
mininet> iperf
...
mininet> h1 ping -c10 h2
```

如果每条链路的延迟为 10 毫秒，那么往返时间 (RTT) 应大约为 40 毫秒，因为 ICMP 请求经过两条链路（一个到交换机，一个到目标），ICMP 回复则经过返回的两条链路。

可以使用 Mininet 的 Python API 自定义每条链路。

#### 2.4 详细信息

默认的详细程度是 `info`，它会在启动和清理过程中打印 Mininet 的操作。使用 `-v` 参数可以查看完整的 `debug` 输出：

```shell
sudo mn -v debug
...
mininet> exit
```

这会打印出大量的详细信息。尝试使用 `output` 级别，它仅打印 CLI 输出，几乎没有其他信息：

```shell
sudo mn -v output
mininet> exit
```

在 CLI 之外，还可以使用其他详细程度级别，例如 warning，该级别用于回归测试中以隐藏不必要的功能输出。

#### 2.5 自定义拓扑

使用简单的 Python API 也可以轻松定义自定义拓扑。在如下示例中，两个交换机直接连接，每个交换机连接一个主机：

<<< @/sdn/codes/mn/topo.py

提供自定义 Mininet 文件时，可以向命令行添加新的拓扑、交换机类型和测试。例如：

```shell
sudo mn --custom ./custom.py --topo mytopo --test pingall
```

#### 2.6 ID = MAC

默认情况下，主机启动时会被随机分配 MAC 地址。这会增加调试难度，因为每次创建 Mininet 时，MAC 地址都会更改。

`--mac` 选项可以将主机的 MAC 和 IP 地址设置为简短、唯一且易于读取的 ID。

不使用 `--mac`：

```shell
sudo mn
...
mininet> h1 ifconfig
"h1-eth0  Link encap:Ethernet  HWaddr f6:9d:5a:7f:41:42
          inet addr:10.0.0.1  Bcast:10.255.255.255  Mask:255.0.0.0
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:6 errors:0 dropped:0 overruns:0 frame:0
          TX packets:6 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:392 (392.0 B)  TX bytes:392 (392.0 B)"
mininet> exit
```

使用 `--mac`：

```shell
sudo mn --mac
...
mininet> h1 ifconfig
"h1-eth0  Link encap:Ethernet  HWaddr 00:00:00:00:00:01
          inet addr:10.0.0.1  Bcast:10.255.255.255  Mask:255.0.0.0
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)"
mininet> exit
```

#### 2.7 其他交换机类型

可以使用其他类型的交换机。例如，要运行用户空间交换机：

```shell
sudo mn --switch user --test iperf
```

::: warning 注意
与之前使用内核交换机相比，此时 iperf 报告的 TCP 带宽要低得多。
:::

如果进行 ping 测试，延迟明显增大，因为现在数据包必须经历内核到用户空间的转换。ping 时间也会更加不稳定，因为主机的用户空间进程可能会被操作系统调度进出。

另一种交换机类型是 Open vSwitch（OVS），它在 Mininet 虚拟机中已预装。iperf 报告的 TCP 带宽与 OpenFlow 内核模块相似，甚至可能更快：

```shell
sudo mn --switch ovsk --test iperf
```

### 3 Mininet CLI 命令

#### 3.1 显示命令列表

要查看命令行界面（CLI）选项列表，先启动一个最小拓扑并保持其运行。构建 Mininet：

```shell
sudo mn
```

显示命令列表：

```shell
mininet> help
```

#### 3.2 Python 解释器

如果在 Mininet 命令行中第一个词是 `py`，则该命令将使用 Python 执行。每个主机、交换机和控制器都有一个关联的 Node 对象。

在 Mininet CLI 中运行：

```shell
mininet> py 'hello ' + 'world'
```

打印可访问的局部变量：

```shell
mininet> py locals()
```

接下来，使用 `dir()` 函数查看节点的可用方法和属性：

```shell
mininet> py dir(s1)
```

可以使用 `help()` 函数查看节点上可用方法的在线文档：

```shell
mininet> py help(h1) (按 "q" 退出阅读文档)
```

还可以执行变量的方法：

```shell
mininet> py h1.IP()
```

#### 3.3 链路启用/禁用

禁用虚拟以太网对的两个部分：

```shell
mininet> link s1 h1 down
```

此时，会看到一个 OpenFlow 端口状态更改通知被生成。要重新启用链路：

```shell
mininet> link s1 h1 up
```

#### 3.4 启动 XTerm

为 h1 和 h2 启动一个 xterm 终端窗口：

```shell
mininet> xterm h1 h2
```

## Mininet自定义命令拓展实现

1. 修改 `mininet/net.py`，编写自定义的功能代码；
2. 修改 `mininet/cli.py`，注册命令；
3. 修改 `bin/mn`，加入到可执行文件中；
4. 重新安装 MiniNet 核心文件：

```shell
mininet/util/install.sh -n
```

## 增加网络性能测试命令

### 1 `mininet/net.py` 中定义 `iperfSingle` 函数

在两个主机间进行 `iperf udp` 测试，并在 server 端记录，实现该功能的 `iperfSingle` 函数写在 `Mininet` 类下，代码如下：

<<< @/sdn/codes/mn/iperf.py

### 2 `mininet/net.py` 中定义 `iperfMulti` 函数

依次为每台主机随机选择另一台主机作为 `iperf` 的服务端，通过调用 `iperfSingle`，自身以客户端按照指定参数发送 UDP 流。

代码编写在 `iperfSingle` 之后：

<<< @/sdn/codes/mn/net.py

### 3 `mininet/cli.py` 中注册 `iperfmulti` 命令

解析用户输入命令，`net.py` 定义的 `iperf_multi` 函数需要在 `CLI` 类中注册成命令：

<<< @/sdn/codes/mn/cli.py

### 4 `bin/mn` 中加入 `iperf-multi` 命令

修改后的 `bin/mn` 部分代码如下：

```python{8,16}
TESTS = {
    name: True
    for name in (
        'pingall',
        'pingpair',
        'iperf',
        'iperfudp',
        'iperfmulti'  # new
    )
}

ALTSPELLING = {
    'pingall': 'pingAll',
    'pingpair': 'pingPair',
    'iperfudp': 'iperfUdp',
    'iperfmulti': 'iperfMulti'  # new
}
```

### 5 重新编译 `mininet`

重新编译安装 `mininet`：

```shell
mininet/util/install.sh -n
```

:::tip
`mininet` 的安装默认使用 Python2，可以修改 `util/install.sh` 文件，指定 Python 版本。
:::

重新创建网络，输入 `iperf-m`，可用 `Tab` 键补全 `iperfmulti`，从而可使用 `iperfmulti` 进行随机流量测试。
