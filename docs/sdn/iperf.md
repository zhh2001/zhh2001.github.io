---
outline: [2, 3]
---

# iPerf

在 Mininet 中，经常会需要用到 iPerf 来模拟网络流量。iPerf 是一款网络性能测试工具，用于测量网络带宽、延迟、抖动和丢包率等指标。

## 安装

```shell
sudo apt install iperf
```

## 启动

在服务端运行下面的命令启动 iPerf 服务器：

```shell
iperf -s
```

在客户端运行以下命令：

```shell
iperf -c [IP]
```

后面的例子都在 Mininet 环境中，我开了两台主机 `h1` 和 `h2`，IP 分别是 `10.0.1.1` 和 `10.0.1.2`，并且在 `h2` 上启动了 iPerf 服务器：

```shell
mininet> h2 iperf -s &
```

## 常规选项

### `-f, --format [bkmaBKMA]`

指定打印出来的带宽的数字格式：

| 字母 |                格式 | \|  | 字母 |                 格式 |
| ---- | ------------------: | --- | ---- | -------------------: |
| `b`  |          `bits/sec` | \|  | `B`  |          `Bytes/sec` |
| `k`  |         `Kbits/sec` | \|  | `K`  |         `KBytes/sec` |
| `m`  |         `Mbits/sec` | \|  | `M`  |         `MBytes/sec` |
| `g`  |         `Gbits/sec` | \|  | `G`  |         `GBytes/sec` |
| `a`  | `adaptive bits/sec` | \|  | `A`  | `adaptive Bytes/sec` |

自适应格式根据需要在千（kilo-）和兆（mega-）之间进行选择。

默认值为 `a`。

::: tip 注意
`K = Kilo = 1024`  
`M = Mega = 1024 * 1024`  
`G = Giga = 1024 * 1024 * 1024`
:::

#### 示例：使用默认值的输出格式

执行命令：

```shell
mininet> h1 iperf -c h2
```

输出信息：

```text{7}
-----------------------------------------------------------------
Client connecting to 10.0.1.2, TCP port 5001
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
[  1] local 10.0.1.1 port 55414 connected with 10.0.1.2 port 5001
[ ID] Interval            Transfer     Bandwidth
[  1] 0.0000-10.3188 sec  46.5 MBytes  37.8 Mbits/sec
```

可以看到带宽的单位是 `Mbits/sec`。

#### 示例：指定参数为 `K` 的输出格式

执行命令：

```shell
mininet> h1 iperf -c h2 -f K
```

输出信息：

```text{7}
-----------------------------------------------------------------
Client connecting to 10.0.1.2, TCP port 5001
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
[  1] local 10.0.1.1 port 55414 connected with 10.0.1.2 port 5001
[ ID] Interval            Transfer      Bandwidth
[  1] 0.0000-10.2313 sec  43520 KBytes  4254 KBytes/sec
```

可以看到带宽的单位是 `KBytes/sec`。

### `-i, --interval [n]`

设置周期性带宽、抖动和丢失报告之间的间隔时间（秒）。如果非零，每隔 `n` 秒带宽就会生成一份报告。如果为 `0`，则不打印定期报告。

默认值为 `0`。

#### 示例：指定每两秒打印一次报告

执行命令：

```shell
mininet> h1 iperf -c h2 -i 2
```

输出信息：

```text{7-12}
-----------------------------------------------------------------
Client connecting to 10.0.1.2, TCP port 5001
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
[  1] local 10.0.1.1 port 55414 connected with 10.0.1.2 port 5001
[ ID] Interval            Transfer     Bandwidth
[  1] 0.0000- 2.0000 sec  10.9 MBytes  45.6 Mbits/sec
[  1] 2.0000- 4.0000 sec  10.0 MBytes  41.9 Mbits/sec
[  1] 4.0000- 6.0000 sec  9.00 MBytes  37.7 Mbits/sec
[  1] 6.0000- 8.0000 sec  9.00 MBytes  37.7 Mbits/sec
[  1] 8.0000-10.0000 sec  9.50 MBytes  39.8 Mbits/sec
[  1] 0.0000-10.4072 sec  48.5 MBytes  39.1 Mbits/sec
```

### `-l, --len [n][KM]`

要读写的缓冲区长度。TCP 默认为 `8KB`，UDP 为 `1470B`。

::: tip 注意
对于 UDP 来说，这是数据报大小，当使用 IPv6 寻址时，需要将其降低到 `1450B` 或更低，避免碎片化。
:::

### `-m, --print_mss`

打印 TCP MSS（Maximum Segment Size，最大段大小）的大小，以及观察到的读取大小。

MSS 通常等于 MTU（Maximum Transmission Unit，最大传输单元）减去 40 字节（用于 TCP/IP 头部）。由于 IP 选项的额外头部空间，一般会报告一个略小的 MSS。

#### 示例

执行命令：

```shell
mininet> h1 iperf -c h2 -m
```

输出信息：

```text{3}
-----------------------------------------------------------------
Client connecting to 10.0.1.2, TCP port 5001
MSS size 536 bytes
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
[  1] local 10.0.1.1 port 53052 connected with 10.0.1.2 port 5001
[ ID] Interval            Transfer     Bandwidth
[  1] 0.0000-10.3415 sec  52.9 MBytes  42.9 Mbits/sec
```

### `-p, --port [n]`

服务器监听的端口，以及客户端连接的服务器端口。

默认值为 `5001`。

### `-u, --udp`

使用 UDP 协议进行传输。

#### 示例

执行命令：

```shell
mininet> h1 iperf -c h2 -u
```

输出信息：

```text
--------------------------------------------------------------------
Client connecting to 10.0.1.2, UDP port 5001
Sending 1470 byte datagrams, IPG target: 11215.21 us (kalman adjust)
UDP buffer size:  208 KByte (default)
--------------------------------------------------------------------
[  1] local 10.0.1.1 port 60142 connected with 10.0.1.2 port 5001
[ ID] Interval            Transfer     Bandwidth
[  1] 0.0000-10.0157 sec  1.25 MBytes  1.05 Mbits/sec
```

### `-w, --window [n][KM]`

将套接字缓冲区大小设置为 `n`。对 TCP 来说，这会设置 TCP 窗口大小。但对于 UDP，它只是接收数据报的缓冲区，限制最大的可接收数据报大小。

### `-V`

iPerf 1.6 版本开始支持。

绑定到一个 IPv6 地址。

服务端命令：

```shell
iperf -s -V
```

客户端命令：

```shell
iperf -c [服务端的IPv6地址] -V
```

### `-v, --version`

打印版本信息。

#### 示例

执行命令：

```shell
iperf -v
```

输出结果：

```text
iperf version 2.1.9 (14 March 2023) pthreads
```

末尾的 `pthreads` 意思是 `POSIX threads`。此外还可能是别的值，如 `win32 threads` 意思是 `Microsoft Win32 threads`，又或者 `single threaded`。

## 服务端选项

### `-s, --server`

启动 iPerf 服务器。

### `-D`

iPerf v1.2 版本开始支持。

在 Unix 上，作为守护进程（daemon）运行。在支持服务的 Win32 上，作为服务运行。

### `-R`

iPerf v1.2 版本开始支持。仅在 Windows 上有效。

移除 iPerf 服务。

### `-o`

iPerf v1.2 版本开始支持。仅在 Windows 上有效。

将输出的信息重定向到指定文件。

### `-c, --client [host]`

iPerf 服务器指定 `-c [host]` 将会限制服务器只就受指定主机的连接。

### `-P, --parallel [n]`

用于指定服务器的最大并行连接数。

默认值为 `0`。表示不限制并行连接数。服务器会接受来自客户端的任意数量的连接（仅受系统资源限制）。

## 客户端选项

### `-b, --bandwidth [n][KM]`

指定带宽，单位为 `bit/sec`。

#### 示例

执行命令：

```shell
mininet> h1 iperf -c h2 -b 2M
```

输出结果：

```text
-----------------------------------------------------------------
Client connecting to 10.0.1.2, TCP port 5001
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
[  1] local 10.0.1.1 port 57368 connected with 10.0.1.2 port 5001
[ ID] Interval            Transfer     Bandwidth
[  1] 0.0000-10.5852 sec  2.75 MBytes  2.18 Mbits/sec
```

### `-c, --client [host]`

在客户端模式下运行 iPerf，连接到指定地址的 iPerf 服务器。

```shell
mininet> h1 iperf -c h2
```

### `-d, --dualtest`

用于启用双向测试模式 ​（dual test mode），即同时测试上行和下行带宽。这意味着客户端和服务器会同时发送和接收数据，从而模拟双向网络流量。

#### 示例

执行命令：

```shell
mininet> h1 iperf -c h2 -d
```

输出结果：

```text
-----------------------------------------------------------------
Client connecting to 10.0.1.2, TCP port 5001
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
-----------------------------------------------------------------
Server listening on TCP port 5001
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
[  1] local 10.0.1.1 port 43076 connected with 10.0.1.2 port 5001
[  2] local 10.0.1.1 port 5001 connected with 10.0.1.2 port 43448
[ ID] Interval            Transfer     Bandwidth
[  2] 0.0000-10.6654 sec  24.0 MBytes  18.9 Mbits/sec
[  1] 0.0000-10.8313 sec  31.0 MBytes  24.0 Mbits/sec
```

### `-r, --tradeoff`

用于启用双向交替测试模式 ​（bidirectional alternating test mode）。与 `-d` 参数（同时双向测试）不同，`-r` 参数会先测试一个方向，然后再测试另一个方向，而不是同时进行。

#### 示例

执行命令：

```shell
mininet> h1 iperf -c h2 -r
```

输出结果：

```text
-----------------------------------------------------------------
Server listening on TCP port 5001
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
-----------------------------------------------------------------
Client connecting to 10.0.1.2, TCP port 5001
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
[  1] local 10.0.1.1 port 47896 connected with 10.0.1.2 port 5001
[ ID] Interval            Transfer     Bandwidth
[  1] 0.0000-10.5505 sec  47.5 MBytes  37.8 Mbits/sec
[  2] local 10.0.1.1 port 5001 connected with 10.0.1.2 port 47030
[ ID] Interval            Transfer     Bandwidth
[  2] 0.0000-10.4139 sec  46.9 MBytes  37.8 Mbits/sec
```

### `-t, --time [n]`

用于指定测试的持续时间（以秒为单位）。

默认值为 `10`。

#### 示例

在前面其他的示例中可以发现，每次测试时间都是 10 秒出头。接下来测试一个 60 秒的持续时间。

执行命令：

```shell
mininet> h1 iperf -c h2 -t 60
```

输出结果：

```text
-----------------------------------------------------------------
Client connecting to 10.0.1.2, TCP port 5001
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
[  1] local 10.0.1.1 port 55366 connected with 10.0.1.2 port 5001
[ ID] Interval             Transfer     Bandwidth
[  1] 0.0000-60.6649 sec   320 MBytes  44.2 Mbits/sec
```

### `-P, --parallel [n]`

用于指定**并行连接数**​（parallel connections），即同时发起多个连接进行测试。通过增加并行连接数，可以更充分地利用网络带宽，尤其是在高带宽或高延迟的网络环境中。

#### 示例：建立 4 个连接开始并行测试

执行命令：

```shell
mininet> h1 iperf -c h2 -P 4
```

输出结果：

```text
-----------------------------------------------------------------
Client connecting to 10.0.1.2, TCP port 5001
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
[  4] local 10.0.1.1 port 35008 connected with 10.0.1.2 port 5001
[  3] local 10.0.1.1 port 34990 connected with 10.0.1.2 port 5001
[  1] local 10.0.1.1 port 34996 connected with 10.0.1.2 port 5001
[  2] local 10.0.1.1 port 34974 connected with 10.0.1.2 port 5001
[ ID] Interval            Transfer     Bandwidth
[  1] 0.0000-10.8666 sec  11.5 MBytes  8.88 Mbits/sec
[  3] 0.0000-10.8663 sec  14.6 MBytes  11.3 Mbits/sec
[  2] 0.0000-10.9576 sec  10.1 MBytes  7.75 Mbits/sec
[  4] 0.0000-10.9563 sec  13.3 MBytes  10.1 Mbits/sec
[SUM] 0.0000-10.3611 sec  49.5 MBytes  40.1 Mbits/sec
```

### `-T, --ttl [n]`

用于为测试流量设置 ​TTL 值。该参数通常用于测试网络路径或排查路由问题。

### `-F`

iPerf v1.2 版本开始支持。

用于指定一个文件作为测试数据源。iPerf 会将指定文件的内容作为测试数据发送到服务器端，而不是使用随机生成的数据。这种方式可以模拟真实文件传输的场景，测试网络的文件传输性能。
