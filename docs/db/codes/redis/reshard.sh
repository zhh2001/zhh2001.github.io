redis-cli -p 7107 CLUSTER MYID
# 记录 7107 的节点 ID，供 reshard 的目标参数使用
redis-cli --cluster reshard 127.0.0.1:7101
# 依提示填写迁移槽数、目标 ID、源 ID，核对计划后确认
