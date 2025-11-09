-- 查询销售部的所有员工信息
SELECT * FROM `tb_employee`
WHERE `dept_id` = (SELECT `id` FROM `tb_dept` WHERE `name` = '销售部');