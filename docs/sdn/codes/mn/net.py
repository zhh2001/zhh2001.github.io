"""运行独立的 Mininet UDP 多流实验，无需修改 Mininet 源码。"""

import argparse
import random

from mininet.log import setLogLevel
from mininet.net import Mininet
from mininet.node import OVSBridge
from mininet.topo import SingleSwitchTopo

from cli import UdpCLI
from iperf import run_udp_pairs


class UdpMininet(Mininet):
    udp_log_dir = '/tmp/mn-udp'

    def iperf_multi(self, udp_bw='10M', seconds=10, seed=0):
        if len(self.hosts) < 2:
            raise ValueError('UDP 多流测试至少需要两台主机')
        rng = random.Random(seed)
        pairs = [(client, rng.choice([host for host in self.hosts
                                      if host is not client]))
                 for client in self.hosts]
        return run_udp_pairs(pairs, udp_bw, seconds, self.udp_log_dir)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--hosts', type=int, default=3)
    parser.add_argument('--bandwidth', default='10M')
    parser.add_argument('--seconds', type=float, default=10)
    parser.add_argument('--seed', type=int, default=0)
    parser.add_argument('--log-dir', default='/tmp/mn-udp')
    parser.add_argument('--test', action='store_true', help='执行测试后退出')
    args = parser.parse_args()
    if args.hosts < 2:
        parser.error('--hosts 至少为 2')

    setLogLevel('info')
    net = UdpMininet(topo=SingleSwitchTopo(k=args.hosts),
                    switch=OVSBridge, controller=None,
                    autoSetMacs=True, autoStaticArp=True)
    net.udp_log_dir = args.log_dir
    try:
        net.start()
        if net.pingAll(timeout='1'):
            raise RuntimeError('主机未全部连通')
        if args.test:
            net.iperf_multi(args.bandwidth, args.seconds, args.seed)
        else:
            UdpCLI(net)
    finally:
        net.stop()


if __name__ == '__main__':
    main()
