CREATE TABLE IF NOT EXISTS tb_user (
    `id`     INT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `name`   VARCHAR(10) NOT NULL UNIQUE COMMENT '姓名',
    `age`    INT CHECK ( `age` BETWEEN 0 AND 150) COMMENT '年龄',
    `status` CHAR(1) DEFAULT '1' COMMENT '状态',
    `gender` CHAR(1) COMMENT '性别'
) COMMENT '用户表';

INSERT INTO `tb_user` (`name`, `age`, `status`, `gender`)
VALUES ('zhang1', 18, '1', '男'), -- ID: 1
       ('zhang2', 22, '0', '男'); -- ID: 2

INSERT INTO `tb_user` (`name`, `age`, `status`, `gender`)
VALUES ('zhang1', 20, '1', '男');
-- Error: Duplicate entry 'zhang1' for key 'tb_user.name'

INSERT INTO `tb_user` (`name`, `age`, `status`, `gender`)
VALUES (NULL, 18, '1', '男');
-- Error: Column 'name' cannot be null

INSERT INTO `tb_user` (`name`, `age`, `status`, `gender`)
VALUES ('zhang3', 160, '1', '男');
-- Error: Check constraint 'tb_user_chk_1' is violated.

INSERT INTO `tb_user` (`name`, `age`, `gender`)
VALUES ('zhang3', 24, '女'); -- status: '1'