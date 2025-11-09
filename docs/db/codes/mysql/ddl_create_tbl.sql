-- CREATE TABLE [IF NOT EXISTS] 表名 (
--     字段名1 字段类型1 [COMMENT '字段注释'],
--     字段名2 字段类型2 [COMMENT '字段注释'],
--     字段名3 字段类型3 [COMMENT '字段注释'],
--     ...
--     字段名n 字段类型n [COMMENT '字段注释']
-- ) [COMMENT '表注释'];

CREATE TABLE IF NOT EXISTS tb_user (     
    `id`     INT              COMMENT '编号',
    `name`   VARCHAR(50)      COMMENT '姓名',
    `age`    TINYINT UNSIGNED COMMENT '年龄',
    `gender` CHAR(1)          COMMENT '性别'
) COMMENT '用户表';