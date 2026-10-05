"""在 Mininet 主机间启动并清理多条 iperf 2 UDP 流。"""

import json
import math
import re
import subprocess
import tempfile
import time
from pathlib import Path

from mininet.log import output


def run_udp_pairs(pairs, udp_bw='10M', seconds=10, log_dir='/tmp/mn-udp'):
    """pairs 中每一项为 (client, server)，返回本次日志目录。"""
    pairs = list(pairs)
    seconds = float(seconds)
    if not pairs or any(client is server for client, server in pairs):
        raise ValueError('每条流需要两个不同的主机')
    if not math.isfinite(seconds) or seconds <= 0:
        raise ValueError('测试时长必须是有限的正数')
    if not re.fullmatch(r'([0-9]+(?:\.[0-9]+)?)([kKmMgG]?)', udp_bw):
        raise ValueError('发送速率格式示例：10M、500k')
    rate = float(re.match(r'[0-9.]+', udp_bw).group())
    if not math.isfinite(rate) or rate <= 0:
        raise ValueError('发送速率必须是有限的正数')
    if len(pairs) > 65535 - 5001 + 1:
        raise ValueError('流数量超过可分配端口范围')

    base = Path(log_dir).resolve()
    base.mkdir(parents=True, exist_ok=True)
    run_dir = Path(tempfile.mkdtemp(prefix='run-', dir=base))
    processes, files, clients = [], [], []
    plan = [dict(client=client.name, server=server.name, port=5001 + i)
            for i, (client, server) in enumerate(pairs)]
    (run_dir / 'flows.json').write_text(json.dumps(
        dict(udp_bw=udp_bw, seconds=seconds, flows=plan), indent=2) + '\n')

    try:
        # 先启动所有服务端，每条流使用独立端口。
        for i, (_, server) in enumerate(pairs):
            port = plan[i]['port']
            log = (run_dir / f'flow-{i + 1}-server.log').open('wb')
            files.append(log)
            process = server.popen(
                ['iperf', '-s', '-u', '-p', str(port), '-i', '1'],
                stdout=log, stderr=subprocess.STDOUT)
            processes.append(process)
            deadline = time.monotonic() + 3
            while True:
                if process.poll() is not None:
                    raise RuntimeError(f'服务端退出，请查看 {log.name}')
                sockets = server.cmd(f'ss -H -lun sport = :{port}')
                if any(len(fields) >= 4 and fields[3].endswith(f':{port}')
                       for fields in map(str.split, sockets.splitlines())):
                    break
                if time.monotonic() >= deadline:
                    raise RuntimeError(f'UDP 端口 {port} 未就绪，请查看 {log.name}')
                time.sleep(0.05)

        for i, (client, server) in enumerate(pairs):
            port = plan[i]['port']
            output(f'*** {client.name} -> {server.name}, UDP port {port}\n')
            log = (run_dir / f'flow-{i + 1}-client.log').open('wb')
            files.append(log)
            process = client.popen(
                ['iperf', '-c', server.IP(), '-u', '-p', str(port),
                 '-b', udp_bw, '-t', str(seconds), '-i', '1'],
                stdout=log, stderr=subprocess.STDOUT)
            processes.append(process)
            clients.append(process)

        deadline = time.monotonic() + seconds + 15
        for process in clients:
            remaining = max(0.1, deadline - time.monotonic())
            if process.wait(timeout=remaining) != 0:
                raise RuntimeError(f'客户端退出异常，请查看 {run_dir}')
    finally:
        # 只清理本次创建的进程，主机之间没有 PID 命名空间隔离。
        for process in processes:
            if process.poll() is None:
                process.terminate()
        for process in processes:
            try:
                process.wait(timeout=3)
            except subprocess.TimeoutExpired:
                process.kill()
                process.wait()
        for log in files:
            log.close()

    output(f'*** 日志目录：{run_dir}\n')
    return run_dir
