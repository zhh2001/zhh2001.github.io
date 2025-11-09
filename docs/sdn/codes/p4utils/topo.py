from p4utils.utils.helper import load_topo
from p4utils.utils.sswitch_p4runtime_API import SimpleSwitchP4RuntimeAPI

# 加载拓扑文件
topo = load_topo('topology.json')

# 初始化控制器
controller = SimpleSwitchP4RuntimeAPI(
    topo['s1']['device_id'],
    topo['s1']['grpc_port'],
    p4rt_path=topo['s1']['p4rt_path'],
    json_path=topo['s1']['json_path']
)

# 遍历 s1 的所有邻居节点
for neigh in topo.get_neighbors('s1'):
    if topo.isHost(neigh):
        # 添加转发表项
        controller.table_add(
            'dmac',
            'forward',
            [topo.get_host_mac(neigh)],
            [str(topo.node_to_node_port_num('s1', neigh))]
        )