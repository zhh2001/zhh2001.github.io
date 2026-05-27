---
outline: [2, 3]
---

# iPerf

跑 Mininet 拓扑时，经常需要一个轻量的打流工具来看看链路能跑多少带宽、丢多少包。iPerf 就是干这件事的，两端一开，几秒钟就能拿到吞吐、抖动、丢包数据。

下面的笔记基于 iPerf2（`iperf` 命令，版本 2.1.9）。所有示例都跑在同一套环境里：Mininet 起两台主机 `h1` 和 `h2`，IP 分别是 `10.0.1.1`、`10.0.1.2`。

::: warning iPerf2 ≠ iPerf3
两者是两套互不兼容的二进制（`iperf` 与 `iperf3`），同一个选项字母经常含义完全不同，混用会踩坑。要点对照：

| 项目         | iPerf2 (`iperf`)        | iPerf3 (`iperf3`)                       |
| ------------ | ----------------------- | --------------------------------------- |
| 默认端口     | `5001`                  | `5201`                                  |
| `-R`         | 移除 Windows 服务       | **反向测试**（server → client 发数据）  |
| 双向         | `-d` 同时 / `-r` 顺序   | 都不支持，改用 `--bidir`                |
| `-P` 并行流  | 客户端、服务端语义略不同| 仅客户端有效，复用单条控制连接          |
| JSON 输出    | 无                      | `-J`（脚本化首选）                      |
| 拥塞算法切换 | 无                      | `-C cubic / bbr / reno…`                |
| 多线程模型   | pthreads                | 单线程 + 单 TCP 控制流                  |

经验法则：自动化采集脚本走 iPerf3 + `-J`；手工跑 Mininet 拓扑、要看 dualtest 或 UDP 抖动报告走 iPerf2。
:::

## 安装

```shell
sudo apt install iperf
```

## 起一个最小测试

服务端守着 5001 端口：

```shell
iperf -s
```

客户端打过去：

```shell
iperf -c [IP]
```

放到 Mininet 里就是：

```shell
mininet> h2 iperf -s &        # h2 当服务端，丢后台
mininet> h1 iperf -c 10.0.1.2 # h1 打流过去
```

默认会跑 10 秒 TCP，结束时给一行总带宽。

## 通用选项

### `-f` 切换单位

`-f [bkmgaBKMGA]`：小写是 bit/sec，大写是 Byte/sec；`k/m/g` 千兆万兆，`a/A` 让 iPerf 自己挑合适的量级。默认 `a`。

约定：`K = 1024`，`M = 1024²`，`G = 1024³`。

```shell
mininet> h1 iperf -c h2        # 默认自适应 → 37.8 Mbits/sec
mininet> h1 iperf -c h2 -f K   # 强制 KByte/sec → 4254 KBytes/sec
```

### `-i` 周期性汇报

`-i n` 让 iPerf 每隔 `n` 秒打印一次中间结果。`0` 表示只在结束时打一次总报告（默认）。

```shell
mininet> h1 iperf -c h2 -i 2
```

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

### `-l` 读写缓冲区

`-l [n][KM]`：TCP 默认 `8K`，UDP 默认 `1470B`。UDP 模式下这个值就是数据报大小；走 IPv6 的话要降到 `1450B` 以下，否则会被分片。

### `-m` 打印 MSS

加上 `-m` 会多输出一行 `MSS size`。粗算：MSS ≈ MTU − 40B（TCP + IP 头），IP 选项还会再吃几个字节，因此实测往往比理论值小一点。

```text{3}
Client connecting to 10.0.1.2, TCP port 5001
...
MSS size 536 bytes
TCP window size: 85.3 KByte (default)
-----------------------------------------------------------------
[  1] local 10.0.1.1 port 53052 connected with 10.0.1.2 port 5001
[ ID] Interval            Transfer     Bandwidth
[  1] 0.0000-10.3415 sec  52.9 MBytes  42.9 Mbits/sec
```

### `-p` 端口

`-p n`：服务端监听端口/客户端连接端口，默认 `5001`。两端要一致。

### `-u` 改走 UDP

```shell
mininet> h1 iperf -c h2 -u -b 10M
```

客户端发完会回收一段 **Server Report**，比 TCP 多了抖动和丢包：

```text
[ ID] Interval        Transfer     Bandwidth
[  1] 0.0-10.0 sec    11.9 MBytes  10.0 Mbits/sec
[  1] Sent 8505 datagrams
[  1] Server Report:
[ ID] Interval        Transfer     Bandwidth        Jitter   Lost/Total Datagrams
[  1] 0.0-10.0 sec    11.9 MBytes  10.0 Mbits/sec   0.045 ms 2/8505 (0.024%)
```

四个 UDP 特有字段的读法：

- **Jitter**：RFC 1889 定义的到达间隔抖动，单位 ms。值越小越稳；几十 ms 以上通常意味着排队或链路拥塞。
- **Lost / Total Datagrams**：丢失数 / 总发送数，括号里是百分比。
- 注意只有 **服务端** 才知道实际收到多少，所以这行总会以 `Server Report` 形式回传到客户端打印。
- UDP 模式不给 `-b` 就只跑约 1 Mbps，不是带宽打不上去，而是 iPerf 默认就限到这个速率。

### `-n` 按字节数测

`-n [n][KM]`：固定传输总量后停止，和 `-t` 互斥。想跑"传 100 MB 要多久"这种场景就用它。

```shell
mininet> h1 iperf -c h2 -n 100M
```

### `-w` 窗口/缓冲区

`-w [n][KM]`：TCP 下是 socket 窗口大小，直接决定带宽时延积上限；UDP 下只是接收缓冲区。

### `-N` 关闭 Nagle

仅 TCP 生效。Nagle 算法会把小包攒起来一起发，做延迟敏感测试（比如交互式协议、小报文 RTT）时务必加 `-N` 关掉，否则吞吐曲线会被攒包人为压平。

### `-B` 绑定本地地址/接口

`-B host[%dev]`：

- **服务端**：只在指定 IP 上 listen，多网卡机器避免误监听到外网。
- **客户端**：强制从指定源地址发包，等价于挑一条出接口；Mininet 多链路拓扑里几乎必加，否则路由表会替你选，结果不一定走你想测的那条链路。
- **UDP 组播**：服务端 `-B 239.x.x.x` 直接加入组播组监听。

```shell
mininet> h1 iperf -c 10.0.1.2 -B 10.0.1.1   # 强制走 h1-eth0
```

### `-V` IPv6（v1.6+）

两端都加 `-V`，客户端 `-c` 后面跟 IPv6 地址即可。

### `-v` 版本

```shell
$ iperf -v
iperf version 2.1.9 (14 March 2023) pthreads
```

末尾的 `pthreads` 表示线程模型用的是 POSIX 线程；其他可能值有 `win32 threads`、`single threaded`。

## 服务端选项

### `-s` 起服务端

老朋友了，不多说。

### `-D`（v1.2+）

后台化运行。Unix 上变 daemon，Windows 上注册成服务。

### `-R` / `-o`（Windows-only，v1.2+）

`-R` 移除已注册的服务；`-o` 把输出重定向到指定文件。

### `-c host`（限制来源）

服务端加 `-c [host]` 后只接受该主机的连接，相当于一道简易白名单。

### `-P n`

服务端的最大并发连接数，默认 `0`（不限，吃满系统资源为止）。

## 客户端选项

### `-b` 限速（UDP 必备）

`-b [n][KM]`：单位 bit/sec。TCP 模式下控制发送速率，UDP 模式下基本是必填项，否则只跑 1 Mbps。

```shell
mininet> h1 iperf -c h2 -b 2M
# → 2.18 Mbits/sec
```

### `-c host`（连接目标）

```shell
mininet> h1 iperf -c h2
```

### `-d` 同时双向（dualtest）

上下行同时打。结果里会看到两个流 ID，分别对应两个方向的吞吐：

```text
[  1] 0.0000-10.83 sec  31.0 MBytes  24.0 Mbits/sec   # h1 → h2
[  2] 0.0000-10.67 sec  24.0 MBytes  18.9 Mbits/sec   # h2 → h1
```

### `-r` 先后双向（tradeoff）

和 `-d` 的差别在于：`-r` 先跑完一个方向，再跑反向，互不抢带宽，结果更"干净"。带宽紧的环境用这个更能反映单向极限。

### `-t` 测多久

`-t n`：测试时长，默认 `10` 秒。想看长时段的稳定性就拉长，比如 `-t 60`。

### `-P n` 并行流

客户端侧的并行连接数。链路 BDP 大、单流吃不满带宽时，把它开到 4 或 8 通常能看见明显的总吞吐提升。报告末尾会多一行 `[SUM]` 汇总：

```text
[  1] 0.0000-10.87 sec  11.5 MBytes  8.88 Mbits/sec
[  2] 0.0000-10.96 sec  10.1 MBytes  7.75 Mbits/sec
[  3] 0.0000-10.87 sec  14.6 MBytes  11.3 Mbits/sec
[  4] 0.0000-10.96 sec  13.3 MBytes  10.1 Mbits/sec
[SUM] 0.0000-10.36 sec  49.5 MBytes  40.1 Mbits/sec
```

### `-T n` TTL

设置发送报文的 TTL，常用来配合 traceroute 排查多跳路径。

### `-F file`（v1.2+）

用一份真实文件当数据源，而不是随机字节。想模拟"传一个 1G 镜像"的场景就用它。

## 速查

| 场景            | 命令                                |
| --------------- | ----------------------------------- |
| 看单向 TCP 上限 | `iperf -c h2`                       |
| 打 UDP 并指定速 | `iperf -c h2 -u -b 100M`            |
| 中途看进度      | `iperf -c h2 -i 1 -t 30`            |
| 多流压满带宽    | `iperf -c h2 -P 8`                  |
| 同时双向        | `iperf -c h2 -d`                    |
| 顺序双向        | `iperf -c h2 -r`                    |
| 改端口避开占用  | `iperf -s -p 5555` / `-c h2 -p 5555`|
