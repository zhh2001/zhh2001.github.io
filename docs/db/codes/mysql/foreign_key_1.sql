-- CREATE TABLE 表名 (
--     字段名 数据类型,
--     ...,
--     [CONSTRAINT] [外键名称] FOREIGN KEY(外键字段名) REFERENCES 主表(主表列名)
-- );

CREATE TABLE tb_member (
    id INT PRIMARY KEY,
    dept_id INT,
    CONSTRAINT fk_member_dept_id FOREIGN KEY (dept_id) REFERENCES tb_dept(id)
) ENGINE=InnoDB;
