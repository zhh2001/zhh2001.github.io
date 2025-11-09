-- 没创建索引时，根据 age、phone 进行排序
EXPLAIN SELECT `id`, `age`, `phone` FROM `tb_user`
ORDER BY `age` ASC, `phone` ASC; -- Using filesort

-- 创建索引
CREATE INDEX `idx_user_age_phone_aa` ON `tb_user` (`age` ASC, `phone` ASC);

-- 创建索引后，根据 age、phone 进行升序排序
EXPLAIN SELECT `id`, `age`, `phone` FROM `tb_user`
ORDER BY `age` ASC, `phone` ASC; -- Using index

-- 创建索引后，根据 age、phone 进行降序排序
EXPLAIN SELECT `id`, `age`, `phone` FROM `tb_user`
ORDER BY `age` DESC, `phone` DESC; -- Backward index scan; Using index

-- 创建索引后，根据 phone、age 进行升序排序
EXPLAIN SELECT `id`, `age`, `phone` FROM `tb_user`
ORDER BY `phone` ASC, `age` ASC; -- Using index; Using filesort

-- 创建索引后，根据 age 进行升序排序，根据 phone 进行降序排序
EXPLAIN SELECT `id`, `age`, `phone` FROM `tb_user`
ORDER BY `age` ASC, `phone` DESC; -- Using index; Using filesort

-- 创建索引，age 从小到大排，phone 从大到小排
CREATE INDEX `idx_user_age_phone_ad` ON `tb_user` (`age` ASC, `phone` DESC);

-- 创建索引 idx_user_age_phone_ad 后，根据 age 进行升序排，根据 phone 进行降序排
EXPLAIN SELECT `id`, `age`, `phone` FROM `tb_user`
ORDER BY `age` ASC, `phone` DESC; -- Using index