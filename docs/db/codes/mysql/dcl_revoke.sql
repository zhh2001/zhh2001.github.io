-- REVOKE 权限列表 ON 数据库名.表名 FROM '用户名'@'主机名';
REVOKE SELECT, INSERT, UPDATE, DELETE ON demo_sql.* FROM 'demo_user'@'localhost';
