USE demo_constraint;
CREATE TABLE tb_dept (id INT PRIMARY KEY, name VARCHAR(50)) ENGINE=InnoDB;
INSERT INTO tb_dept VALUES (1, '研发部');
ALTER TABLE tb_user ADD COLUMN dept_id INT;

-- CREATE TABLE 的示例使用 tb_member，ALTER TABLE 的示例使用 tb_user。
