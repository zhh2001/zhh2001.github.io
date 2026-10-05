-- 使用 7.3 已创建的 (profession, age, status) 联合索引。

/* 验证范围查询的情况 */

-- 观察范围扫描以及 status 条件的过滤位置
EXPLAIN FORMAT=JSON SELECT * FROM `tb_user`
WHERE `profession` = '电子信息' AND `age` > 22 AND `status` = '0';

-- >= 不保证后续列都能收窄整个范围，比较实际计划
EXPLAIN FORMAT=JSON SELECT * FROM `tb_user`
WHERE `profession` = '电子信息' AND `age` >= 22 AND `status` = '0';
