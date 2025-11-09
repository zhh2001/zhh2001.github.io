CREATE PROCEDURE p4(IN p_n INT)
BEGIN
    DECLARE v_total INT DEFAULT 0;
    WHILE p_n > 0
        DO
            SET v_total := v_total + p_n;
            SET p_n := p_n - 1;
        END WHILE;
    SELECT v_total AS sum_result;
END;

CALL p4(100); -- 5050