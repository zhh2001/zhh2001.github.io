-- 以下 Extra 为可能的表现，实际结果取决于数据与成本选择。
-- 未建索引的表可能需要临时表，比较时先移除相应索引
EXPLAIN SELECT `profession`, COUNT(*) FROM `tb_user`
GROUP BY `profession`; -- 可能使用临时表

-- 使用 7.3 已创建的联合索引，需要比较前后计划时单独移除、重建。

-- 创建索引后，根据 profession 字段分组
EXPLAIN SELECT `profession`, COUNT(*) FROM `tb_user`
GROUP BY `profession`;  -- Using index

-- 创建索引后，根据 age 字段分组。不满足最左前缀原则
EXPLAIN SELECT `age`, COUNT(*) FROM `tb_user`
GROUP BY `age`; -- Using index; Using temporary

-- 创建索引后，根据 profession、age 两个字段分组
EXPLAIN SELECT `profession`, `age`, COUNT(*) FROM `tb_user`
GROUP BY `profession`, `age`; -- Using index

-- 创建索引后，先查询 profession 再根据 age 字段分组
EXPLAIN SELECT `age`, COUNT(*) FROM `tb_user`
WHERE `profession` = '网络安全' GROUP BY `age`; -- Using index
