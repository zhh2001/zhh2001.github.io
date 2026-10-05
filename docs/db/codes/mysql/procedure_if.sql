DELIMITER $$
CREATE PROCEDURE p3()
BEGIN
    DECLARE score INT DEFAULT 58;
    DECLARE result VARCHAR(10);

    IF score >= 85 THEN
        SET result := '优秀';
    ELSEIF score >= 60 THEN
        SET result := '及格';
    ELSE
        SET result := '不及格';
    END IF;

    SELECT result;
END$$
DELIMITER ;

CALL p3(); -- 不及格
