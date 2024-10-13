# MiniNet

## MiniNet自定义命令拓展实现

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

```python
class Mininet(object):
    ...  # 原始代码

    def iperf_single(self, hosts=None, udp_bw='10M', period=60):
        if not hosts:
            return
        assert len(hosts) == 2
        client, server = hosts
        filename = client.name[1:] + '.out'
        output(f'*** iperf: testing bandwidth between {client.name} and {server.name}\n')
        iperf_args = 'iperf -u'
        bw_args = f'-b {udp_bw}'

        print('*** Start server')
        server.cmd(f'{iperf_args} -s -i 1 > ~/mn_log/{filename}&')

        print('*** Start client')
        client.cmd(f'{iperf_args} -t {period} -c {server.IP()} {bw_args} > ~/mn_log/client{filename}&')

```

### 2 `mininet/net.py` 中定义 `iperfMulti` 函数

依次为每台主机随机选择另一台主机作为 `iperf` 的服务端，通过调用 `iperfSingle`，自身以客户端按照指定参数发送 UDP 流。

代码编写在 `iperfSingle` 之后：

```python
import random


class Mininet(object):
    ...  # 原始代码

    def iperfSingle(self, hosts=None, udp_bw='10M', period=60):
        ...  # 前面实现的方法

    def iperfMulti(self, bw, period=60):
        for client in self.hosts:
            server = client
            while server == client:
                server = random.choice(self.hosts)
            self.iperfSingle(hosts=[client, server], udp_bw=bw, period=period)
            sleep(0.05)

        sleep(period)
        print('test has done')
```

### 3 `mininet/cli.py` 中注册 `iperfmulti` 命令

解析用户输入命令，`net.py` 定义的 `iperf_multi` 函数需要在 `CLI` 类中注册成命令：

```python
class CLI(Cmd):
    ...  # 原始代码

    def do_iperfmulti(self, line):
        """Multi iperf UDP test between nodes"""
        args = line.split()
        if len(args) == 1:
            self.mn.iperfMulti(args[0])
        elif len(args) == 2:
            self.mn.iperfMulti(args[0], float(args[1]))
        else:
            error('invalid number of args: iperfMulti udpBw period\n' +
                  'udpBw examples:1M 120\n')
```

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
