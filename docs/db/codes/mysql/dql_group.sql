-- SELECT 字段列表 FROM 表名 [WHERE 条件] GROUP BY 分组字段名 [HAVING 分组过滤条件];
SELECT `gender`, COUNT(*) FROM `tb_employee` GROUP BY `gender`;
SELECT `gender`, AVG(`age`) FROM `tb_employee` GROUP BY `gender`;

-- 查询年龄小于 45 的员工，并且根据工作城市分组，获取员工数量大于 3 的城市
SELECT `city`, COUNT(*) AS `city_count`
FROM `tb_employee`
WHERE `age` < 45
GROUP BY `city`
HAVING `city_count` >= 3;