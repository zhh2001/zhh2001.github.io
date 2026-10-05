SET age 18
# OK
OBJECT ENCODING age
# "int"
SET name zhh
# OK
OBJECT ENCODING name
# "embstr"
# Redis 8.2.x 的键空间对象还要计入键名，以下使用单字节键
# 值长为 40 字节，键长与值长之和为 41
SET a aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
# OK
OBJECT ENCODING a
# "embstr"
# 值长增加到 41 字节，创建 kvobj 时改用 raw
SET a aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
# OK
OBJECT ENCODING a
# "raw"
