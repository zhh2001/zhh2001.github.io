"""自定义拓扑示例

两个直接连接的交换机，每个交换机连接一个主机：

   主机 --- 交换机 --- 交换机 --- 主机

添加一个包含键值对的 'topos' 字典来生成新定义的拓扑，
使得可以从命令行传入 '--topo=mytopo' 来调用该拓扑。
"""

from mininet.topo import Topo

class MyTopo( Topo ):

    def build( self ):
        # 新增主机和交换机
        leftHost = self.addHost( 'h1' )
        rightHost = self.addHost( 'h2' )
        leftSwitch = self.addSwitch( 's3' )
        rightSwitch = self.addSwitch( 's4' )

        # 增加链路
        self.addLink( leftHost, leftSwitch )
        self.addLink( leftSwitch, rightSwitch )
        self.addLink( rightSwitch, rightHost )


topos = { 'mytopo': ( lambda: MyTopo() ) }