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
