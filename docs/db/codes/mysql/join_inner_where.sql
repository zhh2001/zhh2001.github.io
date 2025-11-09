-- SELECT 字段列表 FROM 表1, 表2 WHERE 条件 ...;
SELECT `tb_employee`.`name`, `tb_dept`.`name`
FROM `tb_employee`, `tb_dept`
WHERE `tb_employee`.`dept_id` = `tb_dept`.`id`;