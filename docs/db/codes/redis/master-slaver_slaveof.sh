# 在 7002 的 redis-cli 执行：
SLAVEOF 10.120.82.95 7001
# 在 7003 的 redis-cli 执行：
SLAVEOF 10.120.82.95 7001

# 在主节点 7001 上查看信息：
INFO REPLICATION
# # Replication
# role:master
# connected_slaves:2
# slave0:ip=10.120.82.95,port=7002,state=online,offset=280,lag=0
# slave1:ip=10.120.82.95,port=7003,state=online,offset=280,lag=0
# master_failover_state:no-failover
# master_replid:a25030c9b6542dfc6dd30f217121bc2b5e3b6a52
# master_replid2:0000000000000000000000000000000000000000
# master_repl_offset:280
# second_repl_offset:-1
# repl_backlog_active:1
# repl_backlog_size:1048576
# repl_backlog_first_byte_offset:1
# repl_backlog_histlen:280