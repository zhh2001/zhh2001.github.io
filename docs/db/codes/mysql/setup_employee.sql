USE demo_sql;
CREATE TABLE tb_employee (
    id INT PRIMARY KEY,
    workno VARCHAR(20),
    name VARCHAR(50),
    gender CHAR(1),
    age TINYINT UNSIGNED,
    idcard CHAR(18),
    entrydate DATE,
    city VARCHAR(50) DEFAULT '广州'
) ENGINE=InnoDB;
