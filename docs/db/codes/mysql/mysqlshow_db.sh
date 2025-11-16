# 查询 test 库中每个表中的字段数，及行数
mysqlshow -uroot -p123456 test --count
# Database: test
# +-------------+----------+------------+
# |   Tables    | Columns  | Total Rows |
# +-------------+----------+------------+
# | account     |        3 |          2 |
# | my_memory   |        2 |          0 |
# | my_myisam   |        2 |          0 |
# | tb_employee |        7 |          7 |
# | tb_user     |        5 |          4 |
# | user_v_1    |        2 |          2 |
# +-------------+----------+------------+