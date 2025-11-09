-- SET 变量名 = 值;
-- SET 变量名 := 值;
-- SELECT 字段名 INTO 变量名 FROM 表名 ...;
CREATE PROCEDURE p2()
BEGIN
    DECLARE stu_count INT DEFAULT 0;
    SELECT COUNT(*) INTO stu_count FROM `tb_student`;
END;