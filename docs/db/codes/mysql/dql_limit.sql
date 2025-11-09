-- SELECT 字段列表 FROM 表名 LIMIT 偏移量, 查询记录数;

-- 查询第 1 页员工数据，每页展示 20 条数据
SELECT * FROM `tb_employee` LIMIT 0, 20;
-- 如果偏移量为 0 则可以省略，下面语句效果和上面一样
SELECT * FROM `tb_employee` LIMIT 20;

-- 查询第 2 页员工数据，每页展示 20 条数据
SELECT * FROM `tb_employee` LIMIT 20, 20;