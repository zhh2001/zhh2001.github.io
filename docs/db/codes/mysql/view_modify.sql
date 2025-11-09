-- 方式一：通过 OR REPLACE 替换已有视图
--   CREATE [OR REPLACE] VIEW 视图名称[(列名列表)] AS
--   SELECT语句 [WITH [CASCADED | LOCAL] CHECK OPTION];
CREATE OR REPLACE VIEW `user_v_1` AS
SELECT `id`, `name` FROM `tb_user` WHERE `id` <= 3;

-- 方式二：
--   ALTER VIEW 视图名称[(列名列表)] AS
--   SELECT语句 [WITH [CASCADED | LOCAL] CHECK OPTION];
ALTER VIEW `user_v_1` AS
SELECT `id`, `name` FROM `tb_user` WHERE `id` <= 5;