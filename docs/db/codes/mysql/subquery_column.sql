-- 查询销售部和市场部的所有员工信息
SELECT * FROM `tb_employee`
WHERE `dept_id` IN (SELECT `id` FROM `tb_dept`
                    WHERE `name` = '销售部' OR `name` = '市场部');