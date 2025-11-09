-- SELECT 字段列表 FROM 表1 [INNER] JOIN 表2 ON 连接条件 ...;
SELECT `tb_employee`.`name`, `tb_dept`.`name`
FROM `tb_employee`
         INNER JOIN `tb_dept` ON `tb_employee`.`dept_id` = `tb_dept`.`id`;