-- ALTER TABLE 表名
--     ADD CONSTRAINT 外键名 FOREIGN KEY (外键字段名) REFERENCES 主表名(主表字段名);
ALTER TABLE `tb_user`
    ADD CONSTRAINT fk_user_dept_id FOREIGN KEY (dept_id) REFERENCES tb_dept (id);