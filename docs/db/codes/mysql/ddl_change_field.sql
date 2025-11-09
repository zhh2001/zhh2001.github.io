-- ALTER TABLE 表名 CHANGE 旧字段名 新字段名 数据类型 [COMMENT '注释'] [约束];
ALTER TABLE `tb_user` CHANGE `nickname` `nick` VARCHAR(30) COMMENT '用户昵称';