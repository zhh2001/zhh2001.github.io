-- 使用 7.3 已创建的 (profession, age, status) 联合索引。

/* 验证最左前缀原则 */

-- 可用三列等值条件构造查找范围
EXPLAIN FORMAT=JSON SELECT * FROM `tb_user`
WHERE `profession` = '电子信息' AND `age` = 22 AND `status` = '0';

-- 可用前两列构造查找范围
EXPLAIN FORMAT=JSON SELECT * FROM `tb_user`
WHERE `profession` = '电子信息' AND `age` = 22;

-- 可用首列构造查找范围
EXPLAIN FORMAT=JSON SELECT * FROM `tb_user`
WHERE `profession` = '电子信息';

EXPLAIN FORMAT=JSON SELECT * FROM `tb_user` WHERE `age` = 22 AND `status` = '0'; -- 缺少首列，观察优化器选择
EXPLAIN FORMAT=JSON SELECT * FROM `tb_user` WHERE `status` = '0'; -- 缺少首列，观察优化器选择

-- 跳过 age，profession 仍可用于定位，status 可能参与过滤
EXPLAIN FORMAT=JSON SELECT * FROM `tb_user` WHERE `profession` = '电子信息' AND `status` = '0';
