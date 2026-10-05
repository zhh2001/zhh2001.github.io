XADD users 1000-0 name zhh age 18
# "1000-0"
XLEN users
# (integer) 1
XADD users 1001-0 name howard age 17
# "1001-0"
XLEN users
# (integer) 2
