"""带有 UDP 多流命令的 Mininet CLI。"""

from subprocess import TimeoutExpired

from mininet.cli import CLI
from mininet.log import error


class UdpCLI(CLI):
    def do_iperfmulti(self, line):
        """iperfmulti BW [SECONDS [SEED]]：每台主机发起一条 UDP 流。"""
        args = line.split()
        if not 1 <= len(args) <= 3:
            error('用法：iperfmulti BW [SECONDS [SEED]]\n')
            return
        try:
            seconds = float(args[1]) if len(args) >= 2 else 10
            seed = int(args[2]) if len(args) == 3 else 0
            self.mn.iperf_multi(args[0], seconds=seconds, seed=seed)
        except (ValueError, RuntimeError, OSError, TimeoutExpired) as exc:
            error(f'{exc}\n')
