-- CREATE [UNIQUE | FULLTEXT] INDEX 索引名 ON 表名 (字段名, ...);
CREATE INDEX `idx_user_name` ON `tb_user` (`name`);
CREATE INDEX `idx_user_pro_age_sta` ON `tb_user` (`profession`, `age`, `status`);
CREATE UNIQUE INDEX `idx_user_phone` ON `tb_user` (`phone`);