-- 计算从 1 到 n 的和
CREATE PROCEDURE p6(IN p_n INT)
BEGIN
    DECLARE v_total INT DEFAULT 0;
    sum:
    LOOP
        IF p_n <= 0 THEN
            LEAVE sum;
        END IF;
        SET v_total := v_total + p_n;
        SET p_n := p_n - 1;
    END LOOP sum;
    SELECT v_total AS sum_result;
END;

-- 计算从 1 到 n 中间偶数的和
CREATE PROCEDURE p7(IN p_n INT)
BEGIN
    DECLARE v_total INT DEFAULT 0;
    sum:
    LOOP
        IF p_n <= 0 THEN
            LEAVE sum;
        END IF;
        IF p_n % 2 <> 0 THEN
            SET p_n := p_n - 1;
            ITERATE sum;
        END IF;
        SET v_total := v_total + p_n;
    END LOOP sum;
    SELECT v_total AS sum_result;
END;

CALL p6(100); -- 5050
CALL p7(100); -- 2550
