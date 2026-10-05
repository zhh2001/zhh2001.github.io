DROP DATABASE IF EXISTS demo_query;
CREATE DATABASE demo_query CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE demo_query;

CREATE TABLE tb_dept (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
) ENGINE=InnoDB;
INSERT INTO tb_dept VALUES (1, '研发部'), (2, '销售部'), (3, '市场部'), (4, '财务部');

CREATE TABLE tb_employee (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    gender CHAR(1),
    age TINYINT UNSIGNED,
    salary DECIMAL(10, 2),
    job VARCHAR(50),
    manager_id INT,
    dept_id INT,
    entrydate DATE,
    FOREIGN KEY (dept_id) REFERENCES tb_dept(id)
) ENGINE=InnoDB;
INSERT INTO tb_employee VALUES
    (1, 'Howard', '男', 28, 8000, '工程师', 5, 1, '2020-09-01'),
    (2, '张三', '男', 32, 6000, '销售', 5, 2, '2018-01-01'),
    (3, '李四', '女', 26, 5500, '运营', 5, 3, '2021-03-01'),
    (4, '王五', '男', 55, 5000, '销售', 5, 2, '2009-01-01'),
    (5, '赵六', '女', 40, 12000, '经理', NULL, 1, '2010-01-01'),
    (6, '孙七', '女', 29, 8000, '工程师', 5, 1, '2022-06-01'),
    (7, '周八', '男', 30, 6500, '顾问', NULL, NULL, '2023-01-01');

CREATE TABLE tb_user (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    age TINYINT UNSIGNED,
    phone VARCHAR(20),
    profession VARCHAR(50),
    status CHAR(1),
    email VARCHAR(100)
) ENGINE=InnoDB;
INSERT INTO tb_user VALUES
    (1, 'zhang', 22, '13800000001', '电子信息', '0', 'zhang@example.com'),
    (2, 'li', 23, '13800000002', '电子信息', '1', 'li@example.com'),
    (3, 'wang', 22, '13800000003', '网络安全', '0', 'wang@example.com'),
    (4, 'zhao', 30, '13800000004', '网络安全', '1', 'zhao@example.com'),
    (5, 'sun', 25, '13800000005', '电子信息', '0', NULL);

CREATE TABLE tb_sku (id BIGINT PRIMARY KEY, name VARCHAR(100), price DECIMAL(10, 2));
INSERT INTO tb_sku VALUES (1, '键盘', 199), (2, '鼠标', 99), (3, '显示器', 1299);
