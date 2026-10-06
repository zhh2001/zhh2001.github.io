---
title: iPerf 网络吞吐测试与结果解读
description: 在 Mininet 等 Linux 实验环境中使用 iperf 2 和 iperf 3 测试 TCP、UDP 流量，比较版本与参数差异，解读吞吐量、UDP 丢包和抖动，并说明统计口径。
outline: [2, 3]
---

# iPerf 网络吞吐测试 {#iperf}

在 Mininet 中验证转发规则或调整链路参数后，可以用 iperf 生成 TCP、UDP 流量，观察端到端吞吐量、UDP 丢包和抖动。报告统计的是 socket 层的数据传输，不包含 IP、TCP、UDP 和链路层头部开销，不能直接当作物理链路利用率。

本文以 Linux 上的 iperf 2.1.9 为主，命令名为 `iperf`。版本对照和 JSON 示例采用 iperf 3.16，命令名为 `iperf3`。参数依据 [iperf 2.1.9 官方发行源码](https://downloads.sourceforge.net/project/iperf2/iperf-2.1.9.tar.gz)中的 `man/iperf.1`、参数解析与报告实现，以及 [iperf 3.16 手册](https://github.com/esnet/iperf/blob/f9481e1cd35159929458513692e4a8f9fdd1bd6f/src/iperf3.1)。这两个版本用于说明本文的实验环境，安装其他版本时应重新核对两端的 `--help` 和版本号。

## iperf 2 与 iperf 3

两者独立开发，测试协议不兼容。客户端和服务端需要使用同一系列工具，不能用 `iperf` 连接 `iperf3` 服务端。

| 项目         | iperf 2.1.9（Linux）                                     | iperf 3.16                        |
| ------------ | -------------------------------------------------------- | --------------------------------- |
| 默认端口     | `5001`                                                   | `5201`                            |
| 反向测试     | `-R` 或 `--reverse`                                      | `-R`                              |
| 同时双向     | `--full-duplex`，同一 socket 收发。`-d` 使用额外反向连接 | `--bidir`，两个方向分别使用数据流 |
| 顺序双向     | `-r`                                                     | 分别运行普通测试和 `-R` 测试      |
| `-P n`       | 客户端并行流数。服务端含义见后文                         | 客户端并行流数                    |
| 结构化输出   | `-y C` 输出 CSV                                          | `-J` 输出 JSON                    |
| TCP 拥塞算法 | `-Z algo`                                                | `-C algo`                         |
| IPv6         | `-V`                                                     | `-6`，`-V` 表示详细输出           |
| `-Z`         | 选择 TCP 拥塞算法                                        | 使用零拷贝发送                    |
| 线程模型     | 本文使用的构建支持 POSIX 线程                            | 从 3.16 起每条测试流使用独立线程  |

iperf 3.16 之前采用单线程处理测试流。这个变化记录在 [3.16 发布说明](https://github.com/esnet/iperf/blob/f9481e1cd35159929458513692e4a8f9fdd1bd6f/RELNOTES.md#iperf-316-2023-11-30)中，不能再把“iperf 3 是单线程”作为各版本通用的结论。iperf 3 的 `-d` 是调试输出，也不是 iperf 2 的双向测试选项。

## 安装与实验环境

Debian、Ubuntu 可以通过包管理器安装：

```shell
sudo apt update
sudo apt install iperf iperf3
iperf --version
iperf3 --version
```

本文测试环境中的 iperf 2 输出如下：

```text
iperf version 2.1.9 (14 March 2023) pthreads
```

`pthreads` 表示这个构建启用了 POSIX 线程。包管理器提供的版本随发行版而变化，记录实验结果时要保存客户端和服务端的实际版本。

下面的 iperf 2 示例使用 Mininet 2.3.0 的两主机拓扑。显式指定 IP 地址段，便于复现：

```shell
sudo mn --topo minimal --switch ovsbr --controller none --ipbase 10.0.1.0/24
```

进入 CLI 后检查地址和连通性：

```shell
mininet> h1 ip -brief address
mininet> h2 ip -brief address
mininet> pingall
```

h1、h2 的 IPv4 地址分别为 `10.0.1.1`、`10.0.1.2`。这里采用 OVS 独立桥，不依赖控制器。默认链路没有配置限速，因此测试值主要反映本机内核和 CPU 的处理能力。需要模拟链路带宽、时延或丢包时，可参考 [Mininet](./mininet.md) 中的 TCLink 配置。

`mininet>` 是提示符，不属于命令。Mininet CLI 会把节点命令中的 h2 替换为对应 IP，但普通 shell 不会自动识别这个名字。下文使用完整 IP 地址。

## 最小 TCP 测试

在 h2 启动服务端，并保留进程对象和日志文件：

```shell
mininet> px import subprocess
mininet> px tcp_log = open('/tmp/h2-iperf-tcp.log', 'w')
mininet> px tcp_server = h2.popen(['iperf', '-s', '-p', '5001', '-i', '1'], stdout=tcp_log, stderr=subprocess.STDOUT)
mininet> h2 ss -ltn 'sport = :5001'
```

确认输出中出现监听记录后，在 h1 发起测试：

```shell
mininet> h1 iperf -c 10.0.1.2 -p 5001 -t 5 -i 1
```

默认方向为客户端 h1 向服务端 h2 发送数据。省略 `-t` 时，客户端默认发送 10 秒。后续 TCP 示例继续使用这个服务端，全部测试结束后再清理。

普通 TCP 报告包含三个主要字段：

| 字段        | 含义                               |
| ----------- | ---------------------------------- |
| `Interval`  | 本行统计的时间区间                 |
| `Transfer`  | 本端在该区间内写入或读出的数据量   |
| `Bandwidth` | 该数据量除以统计时长得到的平均速率 |

TCP 客户端通常报告写入 socket 的数据，服务端报告从 socket 读出的数据。两端的统计区间和结束时刻可能不同，短测试中数值也可能有差异。分析交付到接收端的吞吐量时，应同时保存服务端日志。

## 通用选项

### `-f` 输出单位

`-f` 改变报告的显示单位，不改变发送速率：

| 取值               | 速率单位与换算                                 |
| ------------------ | ---------------------------------------------- |
| `b`、`k`、`m`、`g` | bit/s、Kbit/s、Mbit/s、Gbit/s，倍率为 1000     |
| `B`、`K`、`M`、`G` | Byte/s、KByte/s、MByte/s、GByte/s，倍率为 1024 |
| `a`、`A`           | 自动选择合适量级，分别显示 bit/s、Byte/s       |

默认使用 `a`。报告中的 `Transfer` 仍以字节计量，例如 `MBytes` 按 2²⁰ 字节换算，不能与十进制 `Mbits/sec` 混为一谈。

```shell
mininet> h1 iperf -c 10.0.1.2 -t 5 -f m
mininet> h1 iperf -c 10.0.1.2 -t 5 -f K
```

数值参数的后缀是另一套规则。iperf 2.1.9 的 `-b`、`-n`、`-l`、`-w` 使用小写 `k/m/g` 表示 10³、10⁶、10⁹，大写 `K/M/G` 表示 2¹⁰、2²⁰、2³⁰。参数本身决定单位：`-b 10M` 是 10,485,760 bit/s，`-n 100M` 是 104,857,600 字节。

iperf 3.16 的速率参数 `-b` 使用十进制倍率，`-b 10M` 是 10,000,000 bit/s。跨工具比较时，可以直接写 `-b 10000000`，避免后缀差异。

### `-i` 报告间隔

`-i n` 每隔 n 秒输出一次区间统计，结束时再输出汇总。iperf 2 默认不打印周期报告：

```shell
mininet> h1 iperf -c 10.0.1.2 -t 10 -i 2
```

区间平均值会掩盖区间内的短暂突发。调整 `-i` 也不会把吞吐测试变成逐包时延测量。

### `-l` 应用缓冲区和 UDP 数据报长度

普通 TCP 测试默认使用 `128K`，即 131,072 字节的应用读写缓冲区。它不是 socket 缓冲区，也不表示每个 TCP 分段的长度。内核根据 MSS、拥塞控制和卸载设置处理 TCP 分段。

UDP 模式下，`-l` 指定 UDP 载荷长度，其中包含 iperf 自己的测量头。iperf 2.1.9 的 IPv4 默认值为 1470 字节，IPv6 客户端默认值为 1450 字节。

是否超出 MTU，要按路径上的 IP MTU 和实际头部计算。路径 MTU 为 1500，且没有 IP 选项、IPv6 扩展头或额外封装时，IPv4 UDP 载荷上限为 1472 字节，IPv6 为 1452 字节。默认值不是所有路径都能安全使用的保证，隧道等场景需要重新计算。

过大的数据报可能在发送时被拒绝。如果系统允许分片，也可能被拆成多个 IP 分片。[IPv6 路由器不执行分片](https://www.rfc-editor.org/rfc/rfc8200.html#section-4.5)。若实验希望避免分片，应设置适合路径的长度，并通过抓包确认。

### `-m` 查看 MSS

`-m` 打印 TCP MSS，`-M n` 则请求设置 `TCP_MAXSEG`，两者作用不同：

```shell
mininet> h1 iperf -c 10.0.1.2 -t 5 -m -e
```

计算 TCP MSS 选项值时，只扣除固定 IP 和 TCP 头部。IP MTU 为 1500 时，IPv4 为 1460 字节，IPv6 为 1440 字节。发送数据时还要为实际使用的 IP、TCP 选项留出空间，具体规则见 [RFC 6691 §2](https://www.rfc-editor.org/rfc/rfc6691.html#section-2)。

iperf 2.1.9 的普通 TCP 客户端在 `connect()` 前生成[设置报告](https://sourceforge.net/p/iperf2/code/ci/7ab75afd882a6595b1187d2e7fa6b4019aa45c17/tree/src/Client.cpp?format=raw)，`-m` 此时读取的 MSS 可能尚未反映连接状态。本环境中，`-m` 打印 `MSS size 536 bytes`，连接行中的 `icwnd/mss/irtt` 却显示 MSS 为 1448。前者不能直接代表连接建立后的分段上限，排查时应结合连接报告、`ss -ti` 和抓包。

### `-p` 端口

`-p n` 指定服务端监听端口或客户端目标端口，两端应一致。iperf 2 默认为 5001。相同地址、端口上的 TCP 和 UDP 是不同 socket，启动 TCP 服务端不会同时提供 UDP 接收服务。

### `-n` 与 `-t` 结束条件

`-n` 指定传输字节数，`-t` 指定发送时长，通常选择其中一种：

```shell
mininet> h1 iperf -c 10.0.1.2 -n 100M
mininet> h1 iperf -c 10.0.1.2 -t 60 -i 1
```

iperf 2.1.9 不会因同时出现 `-n`、`-t` 而一律报错，这两个选项按解析顺序切换结束模式。为避免歧义，不建议同时指定。连接建立、数据排空和结果交换还需要时间，进程总运行时间不一定等于 `-t`。

### `-w` socket 缓冲区

`-w n` 请求设置内核 socket 缓冲区。iperf 2.1.9 的普通正向测试中，客户端设置发送缓冲区，服务端设置接收缓冲区。UDP 也使用相应的发送或接收缓冲区，并非只影响接收端。

```shell
mininet> h1 iperf -c 10.0.1.2 -t 5 -w 64K
```

TCP 接收窗口、拥塞窗口、RTT 和主机处理能力都会影响吞吐量。`-w` 的请求值不能直接当作 TCP 实际通告窗口。Linux 对这些缓冲区值的记账还包含额外空间，读回值可能是请求值的两倍，也可能受内核上限约束。本环境中请求 `64K` 后报告显示 `128 KByte`。需要检查两端实际设置和内核限制，而不能仅凭这行警告判断实验失败。

### `-N` 关闭 Nagle

`-N` 设置 `TCP_NODELAY`，禁用 Nagle 算法。它主要影响小块 TCP 写入的合并行为，不能保证吞吐量或时延一定改善。对比实验应保持该选项一致。普通吞吐报告也不提供小报文 RTT，若要测往返时延，需要另选合适的测试方法。

### `-B` 绑定本地地址或接口

服务端使用 `-B` 可限定监听地址，客户端使用它可指定源地址。例如：

```shell
mininet> h1 iperf -c 10.0.1.2 -B 10.0.1.1 -t 5
```

指定源地址不会替代路由选择。多网卡或多链路环境中，应先检查路由，必要时配置策略路由。Linux 上可使用 `%dev` 请求绑定设备，通常需要 root 或相应能力：

```shell
mininet> h1 ip route get 10.0.1.2 from 10.0.1.1
mininet> h1 iperf -c 10.0.1.2 -B 10.0.1.1%h1-eth0 -t 5
```

UDP 组播服务端还可以通过 `-B` 指定加入的组地址。组播转发是否可用，需要另行检查接口、组成员关系和网络配置。

### `-V` 使用 IPv6

iperf 2 使用 `-V`。两端需要已有可达的 IPv6 地址，单独增加这个选项不会自动配置网络。例如，在本拓扑中显式添加地址：

```shell
mininet> h1 ip -6 addr add fd00:1::1/64 dev h1-eth0 nodad
mininet> h2 ip -6 addr add fd00:1::2/64 dev h2-eth0 nodad
mininet> px ipv6_server = h2.popen(['iperf', '-s', '-u', '-V', '-p', '5004'], stdout=subprocess.DEVNULL, stderr=subprocess.STDOUT)
mininet> h2 ss -lun 'sport = :5004'
mininet> h1 iperf -c fd00:1::2 -V -u -p 5004 -b 1m -l 1400 -t 5
mininet> px ipv6_server.terminate()
mininet> py ipv6_server.wait(timeout=3)
```

`nodad` 用于这个地址已明确分配的隔离实验，省去等待重复地址检测。实际网络应正常执行地址检测。iperf 3 的对应选项是 `-6`。

## UDP 测试与报告

在 h2 启动 UDP 服务端。这里使用 5002，与前面的 TCP 服务区分：

```shell
mininet> px udp_log = open('/tmp/h2-iperf-udp.log', 'w')
mininet> px udp_server = h2.popen(['iperf', '-s', '-u', '-p', '5002', '-i', '1'], stdout=udp_log, stderr=subprocess.STDOUT)
mininet> h2 ss -lun 'sport = :5002'
mininet> h1 iperf -c 10.0.1.2 -u -p 5002 -b 10m -l 1400 -t 2 -i 1
```

确认监听就绪后再运行客户端。一次实测的汇总如下，实际结果取决于运行环境：

```text
[  1] 0.0000-2.0016 sec  2.39 MBytes  10.0 Mbits/sec
[  1] Sent 1790 datagrams
[  1] Server Report:
[ ID] Interval       Transfer     Bandwidth        Jitter   Lost/Total Datagrams
[  1] 0.0000-2.0015 sec  2.39 MBytes  10.0 Mbits/sec   0.006 ms 0/1789 (0%)
```

发送端的速率说明它写出了多少数据，接收端报告则用于判断有多少数据到达应用。`Server Report` 依靠 UDP 结束交互返回，可能因丢包而缺失，禁用结束交互时也不会返回。收集结果时应保存服务端日志。

### 丢包计数

`Lost/Total Datagrams` 是接收端按数据报序号推算的丢失数和预期数据报数。括号中为丢失比例。它不等同于对链路逐包抓取后得到的精确发送计数，也不应直接把客户端 `Sent` 与服务端分母之差当作额外丢包。

上例中的 1790 和 1789 就存在结束报文计数口径差异。乱序、测试末尾数据报和结束交互也会影响统计。UDP 丢失既可能发生在网络中，也可能发生在接收主机的队列或 socket 缓冲区，需要结合接口计数、抓包和主机负载定位。

### Jitter

iperf 的常规 UDP jitter 采用 RTP 到达间隔抖动的平滑估计，算法可对照 [RFC 3550 §6.4.1](https://www.rfc-editor.org/rfc/rfc3550.html#section-6.4.1)。它比较相邻数据报的发送间隔与接收间隔，再平滑两者差值的绝对值，报告单位为 ms。iperf 2.1.9 回传的汇总 jitter 是测试期间这些平滑估计值的平均值，具体可查看[报告实现](https://sourceforge.net/p/iperf2/code/ci/7ab75afd882a6595b1187d2e7fa6b4019aa45c17/tree/src/Reports.c?format=raw)。

这个值不是 RTT、平均单向时延或时延标准差。稳定的时钟偏移可在间隔差中抵消，但时钟速率差、主机调度和排队都会影响结果。不能凭一个固定的 jitter 阈值就断定链路拥塞。

### `-b` 目标速率

普通正向客户端的 `-b n` 设置应用层目标发送速率，单位为 bit/s，可用于 TCP 或 UDP。目标值不保证接收端达到同样的吞吐量，也不包含协议头部开销。

iperf 2.1.9 的 UDP 默认目标为 1,048,576 bit/s，约 1.049 Mbit/s。TCP 默认不设置这个应用限速。为了便于比较，UDP 实验应显式指定速率：

```shell
mininet> h1 iperf -c 10.0.1.2 -u -p 5002 -b 10000000 -l 1400 -t 10
```

iperf 2.1.9 的 TCP 在低目标速率下还会检查 `-l` 是否过大，例如 `-b 1m` 与默认 `-l 128K` 会产生不兼容提示。可以减小应用写入长度：

```shell
mininet> h1 iperf -c 10.0.1.2 -b 1m -l 8K -t 5
```

对于 `-R` 的 TCP 测试，这个版本的 `-b` 限制的是客户端的读取速率，通过 TCP 流控间接影响反向发送。它与普通正向客户端限制写入速率的行为不同。

## 测试方向与并行流

### `-R` 反向单向

```shell
mininet> h1 iperf -c 10.0.1.2 -R -t 5
```

连接仍由 h1 发起，数据改由 h2 发向 h1。iperf 2.1.9 使用已建立的 socket 完成角色反转，通常比需要额外回连的测试更容易用于 NAT 或防火墙路径。本文讨论 Linux，旧 Windows 构建中 `-R` 移除服务的含义不能套用到这里。

### `--full-duplex` 与 `-d` 同时双向

```shell
mininet> h1 iperf -c 10.0.1.2 --full-duplex -t 5
mininet> h1 iperf -c 10.0.1.2 -d -L 5011 -t 5
```

`--full-duplex` 在同一 socket 上同时收发，`-d` 另外建立反向连接。`-L` 指定客户端接收回连的监听端口，需要保证 h2 能连到 h1 的这个端口。

两个方向的负载会同时竞争主机或路径资源。应分别记录各方向结果，双向速率之和不能直接代表某个方向的可用带宽，也不能据此判断底层设备是否支持物理全双工。

### `-r` 顺序双向

```shell
mininet> h1 iperf -c 10.0.1.2 -r -L 5012 -t 5
```

先测 h1 到 h2，再通过额外连接测 h2 到 h1。每个方向各发送约 5 秒，总运行时间还包括连接和收尾。两个方向不同时打流，但测试时刻也不同，结果仍可能受到背景负载变化影响。

### `-P` 并行流

```shell
mininet> h1 iperf -c 10.0.1.2 -P 4 -t 10 -i 1
mininet> h1 iperf -c 10.0.1.2 -u -p 5002 -P 2 -b 1m -t 5
```

`-P n` 在客户端创建 n 条并行流，`[SUM]` 汇总其速率。iperf 2.1.9 和 iperf 3.16 的 `-b` 都按每条流设置，上面的两条 UDP 流总目标约为 2 Mbit/s。

并行流可能提高总吞吐量，也会改变 CPU 开销、拥塞控制竞争和多路径散列结果。做单流与多流比较时，需要记录流数，不能把多流汇总值当作单连接能力。

### `-Z` 选择 TCP 拥塞算法

Linux 上可以先查看内核可用算法，再指定其中一种：

```shell
mininet> h1 cat /proc/sys/net/ipv4/tcp_available_congestion_control
mininet> h1 iperf -c 10.0.1.2 -Z reno -t 5
```

`-Z` 设置本端 socket 的拥塞算法。测试反向或同时双向发送时，也要检查实际发送端的配置，不能假定客户端选项会自动修改另一端。iperf 3 的对应选项为 `-C`。

### `-T` 与 `-F`

`-T n` 用于设置组播 TTL，默认组播 TTL 为 1。iperf 2.1.9 的 Linux 实现也可设置 IPv4 单播 TTL，但这不会让 iperf 具备 traceroute 的逐跳探测功能。

`-F file` 从文件读取发送数据。接收端只统计测试数据，不会按原文件名保存文件，也不提供文件校验和或完整性验证。报告可能包含 iperf 测试头部，不能用 `Transfer` 核验文件大小。需要复制文件时，应使用文件传输工具。

## 服务端与日志

`-s` 启动服务端。Linux 上的 `-D` 让服务端后台化，客户端不能使用这个选项。Mininet 实验中使用前面的 `Popen` 对象更便于保存进程身份和清理。

`-o file` 可以把 iperf 2 的报告写入文件，Linux 也支持，并非 Windows 专用：

```shell
mininet> h1 iperf -c 10.0.1.2 -t 5 -o /tmp/h1-iperf.log
```

iperf 2.1.9 的服务端实现仍处理 `-P n`，但它限制的是接受连接的累计次数，达到次数后停止继续监听，已经启动的测试线程可以完成。它不是最大同时连接数。若希望一次只处理一个客户端，可查看 `-1` 的串行处理选项。

不要把 `iperf -s -c host` 当作来源白名单。在本文版本的实测中，这种写法没有拒绝其他源地址的客户端。需要限制访问范围时，应配置监听地址或防火墙规则，监听地址本身也不是客户端身份认证。

## iperf 3 的 JSON 输出

需要自动处理结果时，可以使用 iperf 3 的 `-J`，本文版本在测试结束后输出完整 JSON。下面仍在 h2 启动服务端，在 h1 执行普通正向 TCP 测试：

```shell
mininet> px server3 = h2.popen(['iperf3', '-s', '-1', '-p', '5201'], stdout=subprocess.DEVNULL, stderr=subprocess.STDOUT)
mininet> h2 ss -ltn 'sport = :5201'
mininet> h1 iperf3 -c 10.0.1.2 -p 5201 -t 5 -J --get-server-output > /tmp/iperf3-tcp.json
mininet> py server3.wait(timeout=3)
```

确认监听后再执行客户端。iperf 3 的 `-1` 表示完成一次测试后退出，与 iperf 2 的 `-1` 含义不同。要运行下一次测试，需要重新启动这个服务端。

普通 shell 中可这样读取接收端汇总速率：

```python
import json
from pathlib import Path

report = json.loads(Path('/tmp/iperf3-tcp.json').read_text())
if 'error' in report:
    raise RuntimeError(report['error'])

rate = report['end']['sum_received']['bits_per_second']
print(f'{rate / 1_000_000:.3f} Mbit/s')
```

`sum_sent` 和 `sum_received` 分别汇总发送端与接收端统计。iperf 3.16 的 UDP 接收汇总中还有 `jitter_ms`、`lost_packets`、`packets`、`lost_percent` 等字段，同时双向测试则增加反向汇总字段。解析脚本应记录版本、协议和方向，不能把某一种模式的字段结构直接套用到所有测试。

iperf 3 即使测试 UDP，也需要 TCP 控制连接来协商参数并交换结果。防火墙规则应允许控制连接和相应的数据流。它与 iperf 2 的 UDP 结束报告交互不同。

## 速查与收尾

下表中的 iperf 2 客户端命令以对应服务端仍在运行为前提：

| 场景                  | 客户端命令                                          |
| --------------------- | --------------------------------------------------- |
| 普通 TCP              | `iperf -c 10.0.1.2 -t 10 -i 1`                      |
| UDP，十进制 10 Mbit/s | `iperf -c 10.0.1.2 -u -p 5002 -b 10m -l 1400 -t 10` |
| 反向单向              | `iperf -c 10.0.1.2 -R -t 10`                        |
| 同一 socket 同时双向  | `iperf -c 10.0.1.2 --full-duplex -t 10`             |
| 顺序双向              | `iperf -c 10.0.1.2 -r -L 5012 -t 10`                |
| 四条 TCP 流           | `iperf -c 10.0.1.2 -P 4 -t 10`                      |
| 固定传输量            | `iperf -c 10.0.1.2 -n 100M`                         |

在 Mininet 中执行时，客户端命令前加 `h1`。全部测试完成后，关闭前面创建的两个服务端并回收进程：

```shell
mininet> px tcp_server.terminate()
mininet> py tcp_server.wait(timeout=3)
mininet> px tcp_log.close()
mininet> px udp_server.terminate()
mininet> py udp_server.wait(timeout=3)
mininet> px udp_log.close()
mininet> exit
```

默认 Mininet Host 共享 PID 空间，按名称执行 `killall iperf` 或 `pkill iperf` 可能影响其他节点和宿主机进程。实验日志位于共享文件系统中，多个服务端也应使用不同文件名。
