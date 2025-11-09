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