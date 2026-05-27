---
outline: deep
---

# Mininet

Mininet 把若干个 Linux network namespace 拼起来当作虚拟主机，再用 veth pair 接到 OVS 之类的软件交换机上，于是一台普通笔记本就能跑出几十节点的网络拓扑。所有主机跑的都是真实 Linux 协议栈，交换机说的是 OpenFlow，所以它常被拿来做 SDN 实验、复现论文拓扑或者教学。

下面笔记基于 Mininet 2.3.x（Python 3，Ubuntu 22.04/24.04 验证）。

## 安装

官方仓库源码安装最稳：

```shell
git clone https://github.com/mininet/mininet
cd mininet
git tag                                  # 查看可用版本
git checkout -b mininet-2.3.0 2.3.0      # 可选：切到稳定版
cd ..
mininet/util/install.sh [选项]
```

`install.sh` 常用开关：

| 选项 | 含义 |
| --- | --- |
| `-a` | 全家桶：Mininet + OVS + OpenFlow 参考交换机 + POX + Wireshark 解码器 |
| `-nfv` | 精简版：仅 Mininet、参考交换机、OVS |
| `-s DIR` | 把源码放到 `DIR`，默认装在 `$HOME` 下 |
| `-n` | 只重装 Mininet 核心 Python 包（修改源码后用得最多） |
| `-w` | 单独装 Wireshark 与 OpenFlow 解码器 |
| `-h` | 列出全部子项 |

装完冒烟测试一下：

```shell
sudo mn --switch ovsbr --test pingall
```

`ovsbr` 是 OVS 的独立桥模式，不挂控制器也能转发，适合验证安装本身有没有问题。

## 一点底层背景：network namespace

理解 Mininet 的关键就一句话：**每台 host 各自一个 netns，交换机默认留在 root netns**。

- `h1 ifconfig` 看到的 `h1-eth0` 在主机 root shell 里 `ifconfig` 找不到，因为它属于 `h1` 的命名空间；
- `s1 ifconfig` 实际等价于在 root netns 跑，所以它能看到外网卡 `eth0`；
- veth pair 一端在 host 的 netns、另一端落在 OVS 上，OVS 按 OpenFlow 流表转发；
- `sudo mn -c` 之所以能修复一切，原因是它把残留的 netns、veth、OVS 桥统统清掉。

这条认知顺手解决了 90% 的“为什么我看不到 / ping 不通”问题。

## 一次最小启动

```shell
sudo mn
```

不带任何参数时，默认是 `minimal` 拓扑：两个主机挂在同一个 OpenFlow 交换机上，外加一个参考控制器。开起来后会落到 `mininet>` 提示符。

四个常用观察命令依次跑一遍：

```shell
mininet> nodes    # 列出所有节点
mininet> net      # 列出链路
mininet> dump     # 列出节点详细信息（IP、PID、所属 netns）
mininet> help     # 全部 CLI 命令
```

在节点上执行命令，把节点名当前缀写就行：

```shell
mininet> h1 ifconfig -a
mininet> h1 ping -c 1 h2
mininet> pingall              # 任意两主机间的全对测试
```

退出：

```shell
mininet> exit
sudo mn -c                    # 退出后习惯性清一次
```

## 内建拓扑模板

`--topo` 选哪种拓扑、跟几个参数。最常用的六种：

| 模板 | 形状 |
| --- | --- |
| `minimal` | 默认。1 个 switch + 2 个 host |
| `single,N` | 1 个 switch 挂 N 个 host |
| `linear,N` | N 个 switch 串成一条线，每个挂 1 个 host |
| `reversed,N` | 同 `single,N`，但主机编号与端口顺序相反 |
| `tree,depth=D,fanout=F` | 深度 D、扇出 F 的多叉树（叶子节点是 host） |
| `torus,M,N` | M×N 环面网格 |

举例：

```shell
sudo mn --topo single,3 --test pingall
sudo mn --topo linear,4 --test pingall
sudo mn --topo tree,depth=2,fanout=3       # 1 → 3 → 9 主机
```

自定义拓扑见后文 [Python API 骨架](#python-api-骨架)。

## 常用命令行选项

按使用频率分组，**带 `★` 的几乎每次都得手敲**。

### 拓扑与节点

| 选项 | 作用 |
| --- | --- |
| ★ `--topo` | 选拓扑模板，见上一节 |
| `--custom file.py` | 加载自定义拓扑/交换机/测试 |
| `--mac` | 把 MAC、IP 改成 `00:00:00:00:00:0X` / `10.0.0.X`，调试时强烈建议加 |
| `--arp` | 预填全网 ARP 表，省掉首包广播解析丢失 |
| `--ipbase 10.0.0.0/8` | 改 IP 段 |

### 交换机与控制器

| 选项 | 作用 |
| --- | --- |
| ★ `--switch ovs / ovsk / user / ovsbr / lxbr` | 选交换机后端，详见下方对比表 |
| ★ `--controller default / ref / ovsc / none` | 选内建控制器 |
| ★ `--controller remote,ip=X,port=6633` | 接外部 SDN 控制器（ONOS / Ryu / Floodlight / 自己写的 P4Runtime 控制器） |
| `--listenport / --nolistenport` | OVS 是否监听被动连接 |

### 链路与流量

| 选项 | 作用 |
| --- | --- |
| ★ `--link tc,bw=10,delay=10ms,loss=1,jitter=2ms` | 通过 `tc` 给链路加带宽/延迟/丢包/抖动 |
| `--link default` | 不限速，默认 |

### 运行控制

| 选项 | 作用 |
| --- | --- |
| `-x` | 启动后自动给每个节点开一个 xterm |
| `--test pingall / pingpair / iperf / iperfudp` | 跑完测试直接退出，做回归很顺手 |
| `--pre / --post script` | 启动前/退出后跑指定脚本 |
| `-v debug / info / output / warning` | 日志级别，info 是默认 |
| `-c` | 清理上次残留（不启动新拓扑） |
| `--nat` | 给 Mininet 网络挂一个 NAT 网关，让 host 能上外网 |
| `-h` | 全部选项 |

## 链路参数怎么算

`--link tc,bw=10,delay=10ms` 给**每一条**链路设定 10 Mbps 带宽、10 ms 单程传播延迟。

以默认 `minimal` 为例，h1 → h2 的路径是 `h1 — s1 — h2`，两段链路。一次 ping 要走过去再走回来，共 4 段单程，所以 RTT 大约 40 ms：

```shell
sudo mn --link tc,bw=10,delay=10ms
mininet> h1 ping -c 5 h2
# rtt min/avg/max/mdev = 40.x/40.x/40.x ms
```

`loss=1` 表示 1% 丢包，`jitter=2ms` 在延迟上叠加高斯抖动。

## 交换机后端对比

| 后端 | 内核/用户态 | 控制器 | 性能 | 典型用途 |
| --- | --- | --- | --- | --- |
| `ovs` | OVS 内核模块 | 需要 | 高 | 默认值，绝大多数实验用它 |
| `ovsk` | 同上 | 需要 | 高 | `ovs` 的旧别名 |
| `ovsbr` | OVS（独立桥模式） | 不需要 | 高 | 拿来跑连通性 smoke test |
| `user` | 用户态参考实现 | 需要 | 低 | 调试 OpenFlow 协议本身、看清协议交互 |
| `lxbr` | Linux bridge | 不需要 | 中 | 不需要 OpenFlow、只要二层桥接的场景 |

```shell
sudo mn --switch user --test iperf   # 带宽明显低、抖动大，正常现象
sudo mn --switch ovs  --test iperf   # 接近 Gbps，参考基线
```

## 外接 SDN 控制器

把内建控制器关掉、指向远端，是做 SDN 实验最常用的姿势：

```shell
# 一台 switch + 3 个 host，控制器在本机 6633 端口（ONOS / Ryu / 自写）
sudo mn --topo single,3 \
        --switch ovs \
        --controller remote,ip=127.0.0.1,port=6633 \
        --mac --arp
```

如果控制器先没起来，OVS 默认会装一条 fail-mode=secure 的全丢策略，`pingall` 全 X 是正常的；先 `controller` 后 `pingall` 即可。

## Mininet CLI 速查

启动后停在 `mininet>` 时，下面这些命令最常用：

| 命令 | 作用 |
| --- | --- |
| `help` | 列全部命令 |
| `nodes` / `net` / `dump` | 节点、链路、详细信息 |
| `<节点> <shell 命令>` | 在该节点的 netns 里执行 |
| `pingall` / `pingpair` | 全对 ping / 一对 ping |
| `iperf [h1 h2]` | TCP 带宽测试 |
| `iperfudp BW [h1 h2]` | UDP 带宽测试，指定速率 |
| `link s1 h1 down / up` | 临时拔/插链路，会触发 OF 端口状态变更 |
| `xterm h1 h2` | 给指定节点弹 xterm |
| `dpctl dump-flows` | 在所有 OF 交换机上看流表 |
| `dpctl add-flow s1 ...` | 手工下流表 |
| `py <Python 表达式>` | 直接在 mininet 进程里跑 Python |
| `sh <shell 命令>` | 在 root netns 跑 shell |
| `exit` / Ctrl-D | 退出并销毁拓扑 |

`py` 这个口子很强大，因为 `net` 就是当前的 Mininet 实例：

```shell
mininet> py net.hosts
mininet> py h1.IP(), h1.MAC()
mininet> py h1.cmd('uname -a')
mininet> py s1.dpctl('dump-flows')
mininet> py dir(s1)               # 查 Node 对象有哪些方法
```

## --mac 前后对比

不加 `--mac`，每次启动 host 都拿到一个随机 MAC，调流表时认不出来谁是谁：

```text
h1-eth0  HWaddr f6:9d:5a:7f:41:42   inet 10.0.0.1
```

加上 `--mac` 后规整成 `00:00:00:00:00:0X`：

```text
h1-eth0  HWaddr 00:00:00:00:00:01   inet 10.0.0.1
h2-eth0  HWaddr 00:00:00:00:00:02   inet 10.0.0.2
```

写 P4 / OpenFlow 流表里要按 MAC 匹配的话，几乎都得开它。

## Wireshark 抓 OpenFlow

```shell
sudo wireshark &
```

监听 `lo` 接口，过滤器填 `openflow_v4`（OF 1.3）或 `of`，就能看到控制器与 OVS 之间的 Hello / FeaturesRequest / PacketIn / FlowMod 流程。

如果没装解码器：

```shell
mininet/util/install.sh -w
```

## 简单 Web 服务测试

```shell
mininet> h1 python3 -m http.server 80 &
mininet> h2 wget -qO- 10.0.0.1
mininet> h1 kill %python3
```

Mininet host 是一个常驻 bash，所以 `kill %python3` 这种 job-control 写法能命中之前 `&` 放后台的进程。如果不放心，用 `pkill -f http.server` 也行。

::: tip
要确认 Mininet 用的 Python 版本：`mininet> py sys.version`。
:::

## Python API 骨架

CLI 灵活有限，复杂拓扑（fat-tree、leaf-spine、自定义带宽矩阵）一律走 Python API。最小可运行例子已经在仓库里：

<<< @/sdn/codes/mn/topo.py

配合 `--custom` 暴露给 `mn`：

```shell
sudo mn --custom ./topo.py --topo mytopo --test pingall
```

或者直接写 `net = Mininet(topo=MyTopo(), ...)` 自己跑：

```python
from mininet.net import Mininet
from mininet.node import RemoteController, OVSSwitch
from mininet.cli import CLI
from mininet.log import setLogLevel

if __name__ == '__main__':
    setLogLevel('info')
    net = Mininet(topo=MyTopo(),
                  switch=OVSSwitch,
                  controller=lambda name: RemoteController(name, ip='127.0.0.1', port=6633),
                  autoSetMacs=True)
    net.start()
    CLI(net)
    net.stop()
```

## 给 Mininet 加自定义 CLI 命令

源码改动 4 个地方，最后重装：

1. `mininet/net.py`：实现新功能；
2. `mininet/cli.py`：把它注册成 CLI 命令；
3. `bin/mn`：登记到 `TESTS` / `ALTSPELLING` 表里，命令行 `--test` 才能调用；
4. `mininet/util/install.sh -n`：只重装核心 Python 包，不再重编 OVS、控制器。

下面以批量随机 UDP 打流为例走一遍。

### `iperfSingle`：两节点点对点 UDP 测试

写在 `Mininet` 类里：

<<< @/sdn/codes/mn/iperf.py

### `iperfMulti`：每台 host 随机挑一台 server

复用 `iperfSingle`，让全网同时打流，看交换机/控制器在压力下的行为：

<<< @/sdn/codes/mn/net.py

### 注册到 CLI

<<< @/sdn/codes/mn/cli.py

### 暴露给 `--test`

```python{8,16}
TESTS = {
    name: True
    for name in (
        'pingall',
        'pingpair',
        'iperf',
        'iperfudp',
        'iperfmulti'   # new
    )
}

ALTSPELLING = {
    'pingall': 'pingAll',
    'pingpair': 'pingPair',
    'iperfudp': 'iperfUdp',
    'iperfmulti': 'iperfMulti'   # new
}
```

### 重装并验证

```shell
mininet/util/install.sh -n
sudo mn
mininet> iperf-m<Tab>      # 应当补全成 iperfmulti
mininet> iperfmulti
```

::: tip
若 `install.sh` 默认走的是 Python 2，可手动改 `util/install.sh` 里调用 `python` 的地方为 `python3`。
:::

## 常见报错速查

| 现象 | 大概原因 / 处理 |
| --- | --- |
| `Cannot find required executable controller` | 没装参考控制器；要么 `install.sh -a`，要么改用 `--controller=ovsc` 或 `remote` |
| `RTNETLINK answers: File exists` | 上次没 clean，先 `sudo mn -c` |
| `pingall` 全部 X | 控制器没起 / 没接上；OVS 进了 fail-secure 默认丢包 |
| `iperf` 一直停在 `Server listening` | 服务端起得比客户端慢，或者两端不在一条路径上 |
| `h1 ifconfig` 看不到 `h1-eth0` | 多半是 `--innamespace` 配置异常，先 `sudo mn -c` 重来 |
| OVS 报 `database connection failed` | `ovsdb-server` / `ovs-vswitchd` 没启动：`sudo service openvswitch-switch start` |

## 收尾

养成两个习惯能省下不少时间：

```shell
sudo mn -c    # 每次实验结束都跑一次
ps aux | grep -E 'mn|ovs|controller'   # 排查残留进程
```
