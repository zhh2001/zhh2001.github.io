SELECT @@autocommit;
SET SESSION autocommit = 0;
-- 本组练习结束后提交或回滚，再恢复默认值：
ROLLBACK;
SET SESSION autocommit = 1;
