EXPLAIN SELECT *
FROM `tb_sku` `t`,
     (SELECT `id` FROM `tb_sku` ORDER BY `id` LIMIT 2000000, 10) `a`
WHERE `t`.`id` = `a`.`id`;