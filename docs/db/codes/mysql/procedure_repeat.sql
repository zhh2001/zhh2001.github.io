DELIMITER $$
CREATE PROCEDURE p5(IN p_n INT)
BEGIN
    DECLARE v_total BIGINT DEFAULT 0;
    IF p_n IS NULL OR p_n < 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'n must be nonnegative';
    END IF;
    IF p_n > 0 THEN
        REPEAT
            SET v_total := v_total + p_n;
            SET p_n := p_n - 1;
        UNTIL p_n = 0 END REPEAT;
    END IF;
    SELECT v_total AS sum_result;
END$$
DELIMITER ;

CALL p5(100); -- 5050
