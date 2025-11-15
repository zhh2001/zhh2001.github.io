CREATE TABLESPACE my_tablespace 
ADD DATAFILE 'my_tablespace.ibd' ENGINE = InnoDB;

CREATE TABLE a (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(10)
) ENGINE=InnoDB TABLESPACE my_tablespace;