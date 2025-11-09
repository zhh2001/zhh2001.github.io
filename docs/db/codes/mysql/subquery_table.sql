-- 查询与张三、李四的职位和薪资相同的员工信息
SELECT * FROM tb_employee
WHERE (`job`, `salary`) IN (SELECT `job`, `salary` FROM `tb_employee`
                            WHERE `name` = '张三' OR `name` = '李四');

-- 查询入职日期是 2010-01-01 之后的员工信息，及其部门信息
SELECT `e`.*, `d`.*
FROM (SELECT * FROM tb_employee WHERE `entrydate` > '2010-01-01') `e`
         LEFT JOIN `tb_dept` `d` ON `e`.`dept_id` = `d`.`id`;