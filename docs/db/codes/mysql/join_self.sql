-- SELECT 字段列表 FROM 表A 别名A JOIN 表A 别名B ON 条件 ...;
SELECT `a`.`name` AS `员工`, `b`.`name` AS `领导`
FROM `tb_employee` `a`
LEFT OUTER JOIN `tb_employee` `b` ON `a`.`manager_id` = `b`.`id`;