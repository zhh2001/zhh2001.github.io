# 查询 test 库中 tb_user 表的详细情况
mysqlshow -uroot -p123456 test tb_user --count
# Database: test  Wildcard: tb_user
# +---------+----------+------------+
# | Tables  | Columns  | Total Rows |
# +---------+----------+------------+
# | tb_user |        5 |          4 |
# +---------+----------+------------+
# 1 row in set.