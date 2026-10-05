-- 服务器默认关闭 LOCAL，需要管理权限启用。
SET @old_local_infile = @@GLOBAL.local_infile;
SET GLOBAL local_infile = ON;

LOAD DATA LOCAL INFILE '/tmp/tb_insert.csv' INTO TABLE tb_insert
CHARACTER SET utf8mb4
FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n'
(id, name);

-- 在同一会话中恢复练习前的设置。
SET GLOBAL local_infile = @old_local_infile;
