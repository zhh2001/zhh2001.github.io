-- 记录上一页最后一条的 `id`，直接从该 `id` 往后查
EXPLAIN SELECT * FROM `tb_sku`
WHERE `id` > 2000000 ORDER BY `id` LIMIT 10;