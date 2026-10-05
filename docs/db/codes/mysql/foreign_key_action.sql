-- ALTER TABLE 表名
--     ADD CONSTRAINT 外键名 FOREIGN KEY (外键字段名) REFERENCES 主表名(主表字段名)
--         ON UPDATE CASCADE ON DELETE CASCADE;
-- 先移除前一例中创建的约束，再演示级联动作。
ALTER TABLE `tb_user` DROP FOREIGN KEY fk_user_dept_id;
ALTER TABLE `tb_user`
    ADD CONSTRAINT fk_user_dept_id FOREIGN KEY (dept_id) REFERENCES tb_dept (id)
        ON UPDATE CASCADE ON DELETE CASCADE;
