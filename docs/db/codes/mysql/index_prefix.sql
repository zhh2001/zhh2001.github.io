-- CREATE INDEX 索引名 ON 表名 (字段名(前缀长度));
CREATE INDEX `idx_user_email_5` ON `tb_user` (`email`(5));