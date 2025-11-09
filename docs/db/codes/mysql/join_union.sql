-- SELECT 字段列表 FROM 表A ...
-- UNION [ALL]
-- SELECT 字段列表 FROM 表B ...;
SELECT * FROM `tb_employee` WHERE `salary` < 5500
UNION
SELECT * FROM `tb_employee` WHERE `age` > 50;