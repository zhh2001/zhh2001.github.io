SELECT * FROM tb_user;
SELECT * FROM tb_user WHERE id = 1;
SELECT * FROM tb_user WHERE name = 'zhang';

SHOW PROFILES;
-- 根据上一条命令实际输出的 Query_ID 选择查询，不固定写为 6。
-- SHOW PROFILE FOR QUERY query_id;
-- SHOW PROFILE CPU FOR QUERY query_id;
