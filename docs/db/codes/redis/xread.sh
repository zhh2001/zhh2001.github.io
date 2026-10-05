XADD users 1000-0 name zhh age 18
# "1000-0"
XADD users 1001-0 name howard age 17
# "1001-0"
XADD users 1002-0 name John age 19
# "1002-0"

XREAD COUNT 2 STREAMS users 0
# 1) 1) "users"
#    2) 1) 1) "1000-0"
#          2) 1) "name"
#             2) "zhh"
#             3) "age"
#             4) "18"
#       2) 1) "1001-0"
#          2) 1) "name"
#             2) "howard"
#             3) "age"
#             4) "17"

# $ 以当前最后一条为起点，没有更新条目时立即返回空值
XREAD COUNT 1 STREAMS users $
# (nil)

# 连续读取时用上次收到的最后一个 ID，返回已有的 1002-0
XREAD COUNT 1 STREAMS users 1001-0
# 返回 1002-0 及其字段

# 阻塞等待后续消息。在第二个 redis-cli 终端于 5 秒内执行：
# XADD users 1003-0 name Alice age 20
XREAD BLOCK 5000 STREAMS users $
# 若第二个终端及时追加，返回 1003-0 及其字段，否则超时返回 (nil)
