-- 创建联合索引
CREATE INDEX `idx_user_pro_age_sta` ON `tb_user` (`profession`, `age`, `status`);

/* 验证范围查询的情况 */

-- 部分失效，只使用到了 profession 和 age 的索引
SELECT * FROM `tb_user`
WHERE `profession` = '电子信息' AND `age` > 22 AND `status` = '0';

-- 完全生效，三个字段的索引全生效
SELECT * FROM `tb_user`
WHERE `profession` = '电子信息' AND `age` >= 22 AND `status` = '0';