DELIMITER $$
CREATE PROCEDURE p4(IN p_n INT)
BEGIN
    DECLARE v_total BIGINT DEFAULT 0;
    IF p_n IS NULL OR p_n < 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'n must be nonnegative';
    END IF;
    WHILE p_n > 0
        DO
            SET v_total := v_total + p_n;
            SET p_n := p_n - 1;
        END WHILE;
    SELECT v_total AS sum_result;
END$$
DELIMITER ;

CALL p4(100); -- 5050
