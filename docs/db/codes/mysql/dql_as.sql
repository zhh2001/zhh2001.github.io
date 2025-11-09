-- SELECT 字段1 [AS 别名1], 字段2 [AS 别名2], ... FROM 表名;
SELECT `id` AS `uid`, `name` AS `uname`, `age` FROM `tb_employee`;

-- AS 可省略
SELECT `id` `uid`, `name` `uname`, `age` FROM `tb_employee`;