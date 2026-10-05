-- SET 变量名 = 值;
-- SET 变量名 := 值;
-- SELECT 字段名 INTO 变量名 FROM 表名 ...;
DELIMITER $$
CREATE PROCEDURE p2_count()
BEGIN
    DECLARE stu_count INT DEFAULT 0;
    SELECT COUNT(*) INTO stu_count FROM `tb_user`;
    SELECT stu_count;
END$$
DELIMITER ;
