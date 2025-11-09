CREATE PROCEDURE p5(IN p_n INT)
BEGIN
    DECLARE v_total INT DEFAULT 0;
    REPEAT
        SET v_total := v_total + p_n;
        SET p_n := p_n - 1;
    UNTIL p_n <= 0 END REPEAT;
    SELECT v_total AS sum_result;
END;

CALL p5(100); -- 5050