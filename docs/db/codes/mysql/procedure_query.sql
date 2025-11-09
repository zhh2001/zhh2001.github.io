-- 查询指定数据库的存储过程及状态信息：
--   SELECT * FROM information_schema.ROUTINES
--   WHERE ROUTINE_SCHEMA = '数据库名';
SELECT * FROM information_schema.ROUTINES
WHERE ROUTINE_SCHEMA = 'test';

-- 查询某个存储过程的定义：
--   SHOW CREATE PROCEDURE 存储过程名称;
SHOW CREATE PROCEDURE p1;