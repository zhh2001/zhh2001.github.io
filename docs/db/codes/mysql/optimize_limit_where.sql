-- 记录上一页最后一条的 `id`，直接从该 `id` 往后查
SET @last_id = 2000000; -- 上一页实际返回的最后一个 id，示例值
EXPLAIN SELECT * FROM `tb_sku`
WHERE `id` > @last_id ORDER BY `id` LIMIT 10;
