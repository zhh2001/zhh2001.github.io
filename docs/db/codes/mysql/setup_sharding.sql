CREATE DATABASE demo_sharding CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
CREATE USER 'demo_sharding'@'localhost' IDENTIFIED BY 'DemoPass_846!';
GRANT CREATE, ALTER, INDEX, SELECT, INSERT, UPDATE, DELETE
ON demo_sharding.* TO 'demo_sharding'@'localhost';
