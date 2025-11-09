-- 2. 开启从本地加载文件导入数据的开关
SELECT @@local_infile; -- 0
SET GLOBAL local_infile = 1;
SELECT @@local_infile; -- 1

-- 3. 执行 LOAD 指令将准备好的数据，加载到表结构中
LOAD DATA LOCAL INFILE '/root/sql.log' INTO TABLE `tb_user`
FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n';