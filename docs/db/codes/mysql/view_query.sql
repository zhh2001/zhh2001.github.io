-- 查看创建视图语句：SHOW CREATE VIEW 视图名称;
SHOW CREATE VIEW `user_v_1`;

-- 查看视图数据：SELECT * FROM 视图名称 ...;
SELECT * FROM `user_v_1`;
SELECT * FROM `user_v_1` WHERE `id` < 3;