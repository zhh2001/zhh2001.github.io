---
outline: deep
---

# Mininet

Mininet 是基于 Linux 的网络仿真（network emulation）工具。它通过进程、网络命名空间和虚拟链路构建主机与交换机，在同一台机器上运行真实的协议栈和应用，常用于 SDN 教学、原型验证和实验复现。

本文按 [Mininet 2.3.0 源码](https://github.com/mininet/mininet/tree/d7f399d7a1200b602bd6a95ae88bae1009664d4a) 说明命令和 API。验证环境为 Ubuntu 24.04（WSL2）、Python 3.12、Mininet 2.3.0、Open vSwitch 3.3.9 和 iperf 2.1.9。其他发行版的软件包和默认配置可能不同，实验前应记录实际版本。

所有节点共享宿主机内核及计算资源。节点数量和吞吐量受 CPU、内存、调度、交换机实现和链路配置影响，不能将 Mininet 的测量结果直接当作独立物理设备的性能。

## 安装

Ubuntu/Debian 可以先使用发行版软件包：

```shell
sudo apt update
sudo apt install mininet openvswitch-switch iperf iproute2 net-tools
sudo systemctl start openvswitch-switch
mn --version
ovs-vsctl --version
iperf --version
```

不使用 systemd 的环境可以通过相应服务管理方式启动 OVS，例如 `sudo service openvswitch-switch start`。Mininet 的内置带宽测试调用 `iperf`（iperf 2），安装 `iperf3` 不能替代这个程序。

先使用不依赖外部控制器的独立桥模式检查安装：

```shell
sudo mn --switch ovsbr --controller none --test pingall
```

若需要阅读或修改源码，可取固定版本：

```shell
git clone https://github.com/mininet/mininet.git
cd mininet
git checkout d7f399d7a1200b602bd6a95ae88bae1009664d4a
PYTHON=python3 util/install.sh -h
```

[该版本的安装脚本](https://github.com/mininet/mininet/blob/d7f399d7a1200b602bd6a95ae88bae1009664d4a/util/install.sh) 提供以下选项：

| 选项     | 含义                                                                                     |
| -------- | ---------------------------------------------------------------------------------------- |
| `-a`     | 安装 all 函数列出的组件，包括 Mininet、OVS、OpenFlow 参考实现、Wireshark、POX 和测试工具 |
| `-nfv`   | 安装 Mininet 及其依赖、OpenFlow 参考实现和 OVS                                           |
| `-s DIR` | 指定依赖源码和构建目录，需放在执行安装的选项之前                                         |
| `-n`     | 安装 Mininet 的依赖和核心文件，不是只重装 Python 包                                      |
| `-w`     | 安装 Wireshark，并按其版本处理旧版 OpenFlow 解码器                                       |
| `-h`     | 显示选项说明                                                                             |

例如 `PYTHON=python3 util/install.sh -n` 会选择 Python 3。脚本包含发布时的依赖包名，较新发行版需要检查包名和 Python 包安装方式。本文实测采用发行版软件包，未将这份历史安装脚本作为 Ubuntu 24.04 的已验证安装流程。

## network namespace 与资源共享

默认情况下，每台 Host 的 shell 运行在独立网络命名空间中，交换机和控制器位于启动 Mininet 的网络命名空间。可以分别查看接口、路由和命名空间标识：

```shell
mininet> h1 ip -brief address
mininet> h1 ip route
mininet> h1 readlink /proc/self/ns/net
mininet> s1 readlink /proc/self/ns/net
mininet> sh readlink /proc/self/ns/net
```

通常 h1 的命名空间标识与 s1、`sh` 不同。Mininet 使用的命名空间一般没有注册为 `ip netns` 的命名对象，因此 `ip netns list` 不一定能列出这些 Host。

链路通常由 veth pair 构成。连接 Host 和 OVS 时，一端位于 Host 的网络命名空间，另一端作为 OVS 端口接入交换机。软件交换机可以采用 OpenFlow，也可以使用独立桥模式，Mininet 本身没有要求所有交换机都使用 OpenFlow。

默认 Host 只隔离网络命名空间，文件系统和 PID 空间仍共享。某个 Host 中执行 `pkill` 或 `killall`，可能影响其他 Host 或宿主机的进程。启动后台服务时应保存本次创建的 PID 或 `Popen` 对象，按对象清理。

`--innamespace` 会采用额外的控制网络配置，不适用于所有交换机后端。尤其是这里使用的 OVS 内核交换机不支持简单地通过这个开关放入独立命名空间。

## 最小启动与观察

```shell
sudo mn --switch ovsbr --controller none
```

默认 `minimal` 拓扑包含两个主机和一个交换机。明确使用 `ovsbr`，可以先验证二层连通性，避免引入控制器依赖。

直接运行 `sudo mn` 时，Mininet 2.3.0 的 [启动逻辑](https://github.com/mininet/mininet/blob/d7f399d7a1200b602bd6a95ae88bae1009664d4a/bin/mn#L324-L357) 会寻找本地参考控制器或 OVS 控制器。若找不到，而且交换机选项仍是 `default`，则回退为 OVSBridge。显式选择 `--switch ovs` 时不会采用这个回退，通常需要同时给出控制器配置。

进入 `mininet>` 后，可以观察拓扑并执行节点命令：

```shell
mininet> nodes
mininet> net
mininet> dump
mininet> help
mininet> h1 ip -brief address
mininet> h1 ping -c 1 h2
mininet> pingall
mininet> exit
```

`dump` 输出节点类型、接口/IP 和 shell PID 等信息，不直接列出网络命名空间标识。节点命令中的 h2 会被 CLI 替换为相应 IP，`pingall` 则测试主机间两个方向的可达性。

正常退出会调用网络的停止流程。异常终止后的清理见文末，不必在每次正常退出后都执行全局清理。

## 内建拓扑

| 模板                    | 含义                                                     |
| ----------------------- | -------------------------------------------------------- |
| `minimal`               | 1 个交换机、2 个主机                                     |
| `single,N`              | 1 个交换机、N 个主机                                     |
| `linear,N`              | N 个交换机串联，每个交换机连接 1 个主机                  |
| `reversed,N`            | 与 single 类似，但主机编号与交换机端口编号的对应顺序反转 |
| `tree,depth=D,fanout=F` | 深度为 D、扇出为 F 的树，叶子为主机                      |
| `torus,M,N`             | M×N 个交换机构成二维环面，默认每个交换机连接 1 个主机    |

```shell
sudo mn --switch ovsbr --controller none --topo single,3 --test pingall
sudo mn --switch ovsbr --controller none --topo linear,4 --test pingall
sudo mn --switch ovsbr --controller none --topo tree,depth=2,fanout=3 --test pingall
```

最后一个例子包含 4 个交换机和 9 个主机。`torus` 要求两个维度均至少为 3，而且存在二层环路，需要启用并等待 STP 收敛，或使用能够处理环路的控制器策略，不能直接套用树形网络的普通学习转发逻辑。

自定义拓扑见 [Python API 骨架](#python-api-骨架)。

## 常用命令行选项

### 拓扑与主机

| 选项                    | 作用                                                 |
| ----------------------- | ---------------------------------------------------- |
| `--topo NAME,...`       | 选择拓扑及其参数                                     |
| `--custom file.py`      | 执行自定义 Python 文件，加载拓扑及相关类或配置       |
| `--mac`                 | 为默认 Host 分配可预测的 MAC，不改变默认 IP 分配策略 |
| `--arp`                 | 配置主机间静态 ARP 项，适用于 IPv4 地址固定的实验    |
| `--ipbase 10.20.0.0/16` | 设置默认主机 IP 地址段                               |

### 交换机与控制器

| 选项                                         | 作用                                                   |
| -------------------------------------------- | ------------------------------------------------------ |
| `--switch ovs`                               | 使用 OVSSwitch，默认 failMode 为 secure                |
| `--switch ovsbr`                             | 使用 OVS 独立桥模式                                    |
| `--controller default`                       | 寻找可用的本地控制器                                   |
| `--controller ref` / `ovsc`                  | 分别使用参考控制器或 OVSController，需要相应可执行文件 |
| `--controller none`                          | 不启动或连接控制器，不自动保证交换机能够转发           |
| `--controller remote,ip=127.0.0.1,port=6653` | 指向已在该地址监听的外部 OpenFlow 控制器               |
| `--listenport 6654` / `--nolistenport`       | 设置交换机的被动 OpenFlow 监听基准端口，或关闭这种监听 |

`RemoteController` 只描述连接目标，不负责启动外部控制器。标准 OVS 的这个连接使用 OpenFlow，P4Runtime 控制器不能直接接到此处。使用 BMv2 等 P4 目标时，需要相应的 Mininet 交换机类及独立的 P4Runtime 连接配置。

### 链路与运行

| 选项                                                 | 作用                                                                  |
| ---------------------------------------------------- | --------------------------------------------------------------------- |
| `--link tc,bw=10,delay=10ms,loss=1,jitter=2ms`       | 使用 TCLink 配置带宽、延迟、丢包和延迟抖动                            |
| `--link default`                                     | 使用默认虚拟链路，不设置这些流量控制参数                              |
| `--test pingall` / `pingpair` / `iperf` / `iperfudp` | 运行指定测试后停止网络                                                |
| `--pre script`                                       | 在拓扑已构建、`net.start()` 之前执行 Mininet CLI 脚本                 |
| `--post script`                                      | 在测试或交互 CLI 结束后、`net.stop()` 之前执行 Mininet CLI 脚本       |
| `-x`                                                 | 为节点启动 xterm，需要可用的图形显示环境                              |
| `-v debug` / `info` / `output` / `warning`           | 设置日志级别，默认 info                                               |
| `--nat`                                              | 添加 NAT 节点和外部连通配置，会修改宿主机的路由、防火墙及转发相关状态 |
| `-c`                                                 | 清理残留并退出                                                        |
| `-h`                                                 | 查看完整选项                                                          |

`--pre` 和 `--post` 文件中的内容是 CLI 命令。需要执行宿主机 shell 命令时，应使用 CLI 的 `sh`，不能直接把任意 shell 脚本作为这两个选项的输入。NAT 还需要可用的外部路由与 DNS，开启该选项并不保证所有网络环境都能访问互联网。

## 链路参数与测量

```shell
sudo mn --switch ovsbr --controller none --link tc,bw=10,delay=10ms
mininet> h1 ping -c 5 h2
```

TCLink 默认会在链路两端的接口出方向配置流量控制。`bw=10` 表示每个方向的带宽限制为 10 Mbit/s，`delay=10ms` 表示对应出方向队列施加 10 ms 延迟。

在 `h1 -> s1 -> h2` 路径上，报文经过两段链路，往返共经过四个受延迟控制的出接口，因此低负载下 RTT 约为 40 ms。本次环境实测平均值约为 40.5 ms，额外排队、调度和协议处理还会影响结果。

`loss=1` 表示 netem 配置的随机丢包概率为 1%，不是整条路径固定丢失 1% 的报文。`jitter=2ms` 配置延迟变化幅度，Mininet 不会仅凭这个参数自动选择高斯分布。[TCIntf 的实现](https://github.com/mininet/mininet/blob/d7f399d7a1200b602bd6a95ae88bae1009664d4a/mininet/link.py#L292-L310) 将这些参数交给 Linux `tc netem`。

`bw` 是流量控制配置值，iperf 的应用有效吞吐量还会受到协议开销和宿主机资源影响。UDP 的 `-b` 是目标发送速率，也不等于接收端的成功接收速率，速率后缀的倍率应以所用 iperf 版本为准。

## 交换机后端

| 后端           | 实现                                       | 转发方式及依赖                                 |
| -------------- | ------------------------------------------ | ---------------------------------------------- |
| `ovs` / `ovsk` | 默认使用 OVS 内核数据路径，ovsk 为兼容别名 | 通过 OpenFlow 流表控制，也可手工配置规则       |
| `ovsbr`        | OVSBridge，failMode 为 standalone          | 使用普通 MAC 学习转发，不依赖控制器            |
| `user`         | 用户态 OpenFlow 参考实现                   | 需要 ofdatapath、ofprotocol 等旧版参考实现组件 |
| `lxbr`         | LinuxBridge                                | 使用 Linux 二层桥，需要 bridge-utils 等依赖    |

后端性能需要在相同拓扑、链路和负载条件下测量。没有统一的“OVS 接近 Gbps”或固定高、中、低等级，也不能将 Mininet 的 `user` 参考交换机与 OVS 的用户态数据路径视为同一种实现。

## 外接 OpenFlow 控制器

先在目标地址启动控制器，并启用适用于该拓扑的转发应用。下面的例子要求控制器监听 TCP 6653 并支持 OpenFlow 1.3：

```shell
sudo mn --topo single,3 \
        --switch ovs,protocols=OpenFlow13 \
        --controller remote,ip=127.0.0.1,port=6653 \
        --mac --arp
```

连接后可以检查 OVS 的控制器配置和连接状态：

```shell
mininet> sh ovs-vsctl get-controller s1
mininet> sh ovs-vsctl list Controller
mininet> sh ovs-ofctl -O OpenFlow13 dump-flows s1
mininet> pingall
```

OVSSwitch 默认采用 `fail_mode=secure`。按 [OVS 3.3.9 的定义](https://github.com/openvswitch/ovs/blob/1a3fefbcd0b0f70ced621b92e286fef3dffaac09/vswitchd/vswitch.xml#L1353-L1377)，控制器不可用时，它不会自行建立学习转发规则。已有规则仍按各自的超时和匹配条件处理报文，没有有效规则时无法连通。连接控制器和安装转发规则是两件事，CLI 中也没有用于启动远端控制器的 `controller` 命令。

## Mininet CLI 速查

| 命令                              | 作用                                                    |
| --------------------------------- | ------------------------------------------------------- |
| `help`                            | 列出命令，`help CMD` 查看具体命令用法                   |
| `nodes` / `net` / `dump`          | 查看节点、连接和节点详细信息                            |
| `<节点> <shell 命令>`             | 在该节点的网络环境中执行命令                            |
| `pingall` / `pingpair`            | 测试主机间可达性，或第一台与最后一台主机的可达性        |
| `iperf h1 h2`                     | 内置 TCP 带宽测试                                       |
| `iperfudp` / `iperfudp 10M h1 h2` | 默认 UDP 测试，或指定目标发送速率和两台主机             |
| `link s1 h1 down` / `up`          | 设置链路两端接口的状态                                  |
| `xterm h1 h2`                     | 为指定节点启动终端                                      |
| `dpctl dump-flows`                | 在所有交换机上调用对应的 dpctl 工具，OVS 使用 ovs-ofctl |
| `py <表达式>`                     | 求值 Python 表达式                                      |
| `px <语句>`                       | 执行 Python 语句                                        |
| `sh <shell 命令>`                 | 在 Mininet 进程的宿主网络环境中执行                     |
| `exit` / Ctrl-D                   | 结束 CLI，随后由启动脚本停止网络                        |

`iperfudp 10M` 不是有效的内置 CLI 写法，要同时给出两个节点。当前验证环境中，内置 UDP 测试正常，但内置 TCP 测试出现等待服务端统计输出的情况，不能将它一概归因于服务端未启动。需要检查 Mininet 与 iperf 版本的输出兼容性，也可以通过 `Node.popen()` 显式启动客户端和服务端。

`dpctl` 会将相同参数传给所有交换机，没有 `dpctl add-flow s1 ...` 这种选择单台交换机的语法。针对指定 OVS 交换机，可以使用：

```shell
mininet> py s1.dpctl('dump-flows')
mininet> sh ovs-ofctl -O OpenFlow13 dump-flows s1
```

第二条适用于已配置 OpenFlow 1.3 的交换机。手工加入规则时同样需要选择正确协议版本，例如在前一节的 s1 上执行：

```shell
mininet> sh ovs-ofctl -O OpenFlow13 add-flow s1 'priority=100,in_port=1,actions=output:2'
```

这只增加一条从端口 1 到端口 2 的单向规则，不能据此认为其他方向或其他端口已经连通。控制器也可能随后修改这些规则。

Python 命令可以直接引用当前网络和节点：

```shell
mininet> py net.hosts
mininet> py h1.IP(), h1.MAC()
mininet> py h1.cmd('uname -r')
mininet> py sys.version
mininet> py dir(s1)
```

## --mac 的作用

默认 IP 分配按网络创建主机的顺序进行，通常从 `10.0.0.1` 开始，与 `--mac` 无关。启用该选项后，默认主机 MAC 也按序分配：

```text
h1: IP 10.0.0.1，MAC 00:00:00:00:00:01
h2: IP 10.0.0.2，MAC 00:00:00:00:00:02
```

这有利于编写和检查匹配规则，但不是 P4 或 OpenFlow 的必需条件。通过 Python API 显式设置的 IP、MAC 可以覆盖默认值。

## Wireshark 抓取 OpenFlow

Wireshark 已内置 OpenFlow 解码器，新版本通常无需安装旧的外部插件。可以按发行版的 dumpcap 权限配置，以普通用户运行：

```shell
wireshark
```

控制器与 OVS 都在本机时通常选择 `lo`，远端控制器则应选择连接实际经过的接口。显示过滤器可以使用 `tcp.port == 6653`、`openflow_v1`（OpenFlow 1.0）或 `openflow_v4`（OpenFlow 1.3）。OpenFlow 1.3 的过滤器名称可对照 [Wireshark v4.2.2 的解码器定义](https://github.com/wireshark/wireshark/blob/40459284278611128aac5cef35a563218933f8da/epan/dissectors/packet-openflow_v4.c)。`of` 不是这里使用的通用 OpenFlow 过滤器名称。非标准端口可能需要手动 Decode As，TLS 连接则不能直接按明文 OpenFlow 消息解析。

## 简单 Web 服务

在前面的两主机独立桥拓扑中，可以用 Python 的 `Popen` 对象管理服务进程：

```shell
mininet> px import subprocess
mininet> px web = h1.popen(['python3', '-m', 'http.server', '8080'], stdout=subprocess.DEVNULL, stderr=subprocess.STDOUT)
mininet> h2 wget --timeout=3 --tries=3 --retry-connrefused --waitretry=1 -qO- http://10.0.0.1:8080/
mininet> px web.terminate()
mininet> py web.wait(timeout=3)
```

HTTP 服务使用共享文件系统中的当前工作目录。停止时只操作 `web` 对应的进程，比按 `http.server` 名称全局查杀更明确。

## Python API 骨架

自定义拓扑位于 `docs/sdn/codes/mn/topo.py`：

<<< @/sdn/codes/mn/topo.py

在仓库根目录，可以通过 `topos` 字典加载该拓扑：

```shell
sudo mn --custom docs/sdn/codes/mn/topo.py --topo mytopo \
        --switch ovsbr --controller none --test pingall
```

也可以在与 `topo.py` 相同的目录保存下面的 `run_topo.py`，再执行 `sudo python3 run_topo.py`：

```python
from mininet.cli import CLI
from mininet.log import setLogLevel
from mininet.net import Mininet
from mininet.node import OVSBridge
from topo import MyTopo

if __name__ == '__main__':
    setLogLevel('info')
    net = Mininet(topo=MyTopo(), switch=OVSBridge,
                  controller=None, autoSetMacs=True)
    try:
        net.start()
        CLI(net)
    finally:
        net.stop()
```

这里明确导入了 `MyTopo`，并在退出或运行异常时调用 `net.stop()`。如果切换为 OVSSwitch 和 RemoteController，还需要配置控制器地址、协议版本及实际转发逻辑。

## UDP 多流实验与自定义 CLI

下面的三个文件采用 Mininet 和 CLI 子类，不需要修改上游源码或重装 Mininet。每台主机向另一台主机发送一条 UDP 流，随机种子控制流向选择。运行前检查主机数量，每条流分配独立端口，避免多个客户端选中同一个服务端时发生监听冲突。

### 流的启动与清理

`iperf.py` 使用 iperf 2。先确认所有服务端 UDP 端口已监听，再启动客户端。客户端结束或出现异常后，清理本次创建的进程，日志保存到独立目录。

<<< @/sdn/codes/mn/iperf.py

### 自定义 CLI

`cli.py` 添加 `iperfmulti BW [SECONDS [SEED]]` 命令。Python 的 `cmd` 机制会识别 `do_iperfmulti`，交互补全前缀是 `iperfm`，不是带连字符的 `iperf-m`。

<<< @/sdn/codes/mn/cli.py

### 网络与启动入口

`net.py` 负责构建网络、选择流向和停止拓扑：

<<< @/sdn/codes/mn/net.py

在仓库根目录执行批量测试：

```shell
sudo python3 docs/sdn/codes/mn/net.py --test \
    --hosts 3 --bandwidth 1M --seconds 5 --seed 0
```

交互运行则使用：

```shell
sudo python3 docs/sdn/codes/mn/net.py
mininet> iperfmulti 1M 5 0
mininet> exit
```

这项扩展属于示例脚本，系统自带 `mn --test iperfmulti` 不会因此获得同名测试。若要集成到 `mn`，需按其测试注册机制另行配置。

每次运行会打印 `/tmp/mn-udp/run-*` 日志目录。`flows.json` 记录流向、端口、目标速率和时长，`flow-N-client.log`、`flow-N-server.log` 分别保存两端输出。脚本以 root 运行时，读取该目录通常也需要 sudo。固定种子可以复现流向，不能保证调度顺序和性能数值完全相同。

## 常见问题与清理

| 现象                             | 检查方向                                                                                                 |
| -------------------------------- | -------------------------------------------------------------------------------------------------------- |
| 找不到默认控制器                 | 确认是否需要 OpenFlow 控制器，或显式使用 ovsbr 与 controller none。ovsc 也需要对应程序，不能保证自动可用 |
| `RTNETLINK answers: File exists` | 检查重名接口和未结束的实验，确认残留后再清理                                                             |
| `pingall` 出现 X                 | 检查地址、ARP、接口状态、控制器连接、流表和链路丢包配置                                                  |
| 内置 iperf 等待不返回            | 检查连通性、进程、端口和版本输出兼容性，必要时改用显式进程管理                                           |
| Host 中看不到预期接口            | 用 nodes、net、dump 和 ip link 检查命名及创建结果，不能仅凭这一现象判断 --innamespace 异常               |
| OVS 数据库连接失败               | 检查 ovsdb-server、ovs-vswitchd 和 openvswitch-switch 服务                                               |
| xterm/Wireshark 无法显示         | 检查图形环境及 DISPLAY，远程 SSH 还需相应显示转发配置                                                    |

异常退出留下资源时，可以使用：

```shell
sudo mn -c
```

[Mininet 2.3.0 的 cleanup 实现](https://github.com/mininet/mininet/blob/d7f399d7a1200b602bd6a95ae88bae1009664d4a/mininet/clean.py) 会查杀多类进程、移除 OVS 桥、删除匹配命名规则的接口及部分临时文件，范围不限于最近一次实验。共享环境中要先检查其他实验和 OVS 桥，不能把这个命令当作每次退出后的固定步骤。

正常实验应先停止自己创建的应用，再调用 `net.stop()`。配置错误、控制器策略、内核能力和链路丢包也可能导致不连通，清理残留只处理其中一类原因。
