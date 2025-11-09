-- SELECT 字段列表 FROM 表名 WHERE 条件列表;
SELECT * FROM `tb_employee` WHERE `age` >= 20;
SELECT * FROM `tb_employee` WHERE `age` BETWEEN 40 AND 50;
SELECT * FROM `tb_employee` WHERE `idcard` IS NULL;
SELECT * FROM `tb_employee` WHERE `id` IN (1, 2, 3);
SELECT * FROM `tb_employee` WHERE `name` LIKE 'z_';