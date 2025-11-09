-- 创建联合索引
CREATE INDEX `idx_user_pro_age_sta` ON `tb_user` (`profession`, `age`, `status`);

/* 验证最左前缀原则 */

-- 生效，三个索引字段全生效
SELECT * FROM `tb_user`
WHERE `profession` = '电子信息' AND `age` = 22 AND `status` = '0';

-- 生效，前两个索引字段生效
SELECT * FROM `tb_user`
WHERE `profession` = '电子信息' AND `age` = 22;

-- 生效，第一个字段全生效
SELECT * FROM `tb_user`
WHERE `profession` = '电子信息';

SELECT * FROM `tb_user` WHERE `age` = 22 AND `status` = '0'; -- 未生效
SELECT * FROM `tb_user` WHERE `status` = '0'; -- 未生效

-- 部分失效，profession 字段索引生效，跳过了 age 字段导致 status 字段未生效
SELECT * FROM `tb_user` WHERE `profession` = '电子信息' AND `status` = '0';