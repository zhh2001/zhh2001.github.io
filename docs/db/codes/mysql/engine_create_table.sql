-- CREATE TABLE 表名 (
--     字段1, 字段1类型 [COMMENT 字段1注释],
--     ......
--     字段n, 字段n类型 [COMMENT 字段n注释]
-- ) ENGINE = 存储引擎 [COMMENT 表注释];

-- 创建表 my_myisam，并指定 MyISAM 存储引擎
CREATE TABLE `my_myisam` (
    `id`   int         NOT NULL,
    `name` varchar(10) NOT NULL
) ENGINE = MyISAM;

-- 创建表 my_memory，并指定 MEMORY 存储引擎
CREATE TABLE `my_memory` (
    `id`   int         NOT NULL,
    `name` varchar(10) NOT NULL
) ENGINE = MEMORY;