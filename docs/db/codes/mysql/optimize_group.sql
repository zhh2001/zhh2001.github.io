-- 创建索引前，根据 profession 字段分组
EXPLAIN SELECT `profession`, COUNT(*) FROM `tb_user`
GROUP BY `profession`; -- Using temporary（临时表，效率较低）

-- 创建索引
CREATE INDEX `idx_user_pro_age_sta` ON `tb_user` (`profession`, `age`, `status`);

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