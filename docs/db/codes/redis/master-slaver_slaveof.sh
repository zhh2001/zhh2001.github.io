# 在 shell 中配置两个副本
redis-cli -p 7002 REPLICAOF 127.0.0.1 7001
# OK
redis-cli -p 7003 REPLICAOF 127.0.0.1 7001
# OK

# 查看复制角色及连接状态，同步需要一定时间
redis-cli -p 7001 INFO replication
# role:master
# connected_slaves:2
# 其他字段省略
redis-cli -p 7002 INFO replication
# role:slave
# master_host:127.0.0.1
# master_port:7001
# master_link_status:up
# 其他字段省略
