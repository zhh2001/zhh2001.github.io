SETBIT bits 1 1
# (integer) 0
SETBIT bits 3 1
# (integer) 0
SETBIT bits 4 1
# (integer) 0
SETBIT bits 7 1
# (integer) 0

## 此时 bits = 01011001
## 从 bits[1] 开始读，读 4 位（1011），作为一个无符号数
BITFIELD bits GET u4 1
# 1) (integer) 11