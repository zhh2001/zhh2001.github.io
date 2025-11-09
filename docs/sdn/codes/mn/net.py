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