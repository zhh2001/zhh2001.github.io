-- 查询与 Howard 的薪资和直属领导相同的员工信息
SELECT * FROM `tb_employee`
WHERE (`salary`, `manager_id`) = (SELECT `salary`, `manager_id` FROM `tb_employee`
                                  WHERE `name` = 'Howard');