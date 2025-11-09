-- CREATE PROCEDURE 存储过程([参数列表])
-- BEGIN
--     SQL 语句
-- END;
CREATE PROCEDURE p1()
BEGIN
    SELECT COUNT(*) FROM `tb_user`;
END;