SETBIT bits 1 1
# (integer) 0
SETBIT bits 3 1
# (integer) 0
SETBIT bits 4 1
# (integer) 0
SETBIT bits 7 1
# (integer) 0

## 此时 bits = 01011001
BITPOS bits 1
# (integer) 1
BITPOS bits 0
# (integer) 0
BITPOS bits 1 2 8 BIT
# (integer) 3