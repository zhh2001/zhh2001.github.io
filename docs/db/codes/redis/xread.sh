XADD users * name zhh age 18
# "1760775173913-0"
XADD users * name howard age 17
# "1760775244371-0"
XADD users * name John age 19
# "1760775294634-0"

XREAD COUNT 2 STREAMS users 0
# 1) 1) "users"
#    2) 1) 1) "1760775173913-0"
#          2) 1) "name"
#             2) "zhh"
#             3) "age"
#             4) "18"
#       2) 1) "1760775244371-0"
#          2) 1) "name"
#             2) "howard"
#             3) "age"
#             4) "17"

## $ 是指最新发送来的消息，下面这样是读不到的
XREAD COUNT 1 STREAMS users $
# (nil)

## 需要加上参数阻塞等待有人发送消息
XREAD BLOCK 5000 STREAMS users $
# 1) 1) "users"
#    2) 1) 1) "1760776311849-0"
#          2) 1) "name"
#             2) "howard"
#             3) "age"
#             4) "17"
# (2.66s)