---
title: MySQL 查询、索引与事务机制
description: 整理 MySQL 的 SQL 语句、约束、多表查询和事务，介绍 InnoDB 索引、锁与日志机制，说明执行计划、查询优化、主从复制及数据库管理的使用条件。
outline: [2, 3]
---

# MySQL

本文以 **MySQL 8.4.x LTS** 为主，SQL 示例在 MySQL Community Server **8.4.6** 上验证。存储引擎相关内容主要讨论 InnoDB，与旧版行为不同的地方会单独注明。命令行示例使用 `mysql` 客户端，配置文件片段放在 `[mysqld]` 配置组下。

<<< @/db/codes/mysql/conn.sh

## 1. SQL

| 分类 | 全称                       | 说明                                                   |
| ---- | -------------------------- | ------------------------------------------------------ |
| DDL  | Data Definition Language   | 数据定义语言，用来定义数据库对象（数据库、表、字段）   |
| DML  | Data Manipulation Language | 数据操作语言，用来对数据库表中的数据进行增删改         |
| DQL  | Data Query Language        | 数据查询语言，用来查询数据库中表的记录                 |
| DCL  | Data Control Language      | 数据控制语言，用来创建数据库用户、控制数据库的访问权限 |

这是便于学习的分类，`SELECT` 在 MySQL 手册中归入数据操作语句。`START TRANSACTION`、`COMMIT`、`ROLLBACK` 属于事务控制语句。

### 1.1 DDL

#### 1.1.1 数据库操作

先创建并选择 `demo_sql`，删除数据库的语句在完成练习后执行。

- 查询

<<< @/db/codes/mysql/ddl_query.sql

- 创建

<<< @/db/codes/mysql/ddl_create.sql

- 删除

<<< @/db/codes/mysql/ddl_drop.sql

- 使用

<<< @/db/codes/mysql/ddl_use.sql

#### 1.1.2 表操作

先创建 `tb_user`，其余语句以相应对象存在为前提。改表时需要保留的默认值、注释和约束应在完整列定义中一并指定。

- 查询当前数据库所有表

<<< @/db/codes/mysql/ddl_show_tbl.sql

- 查询表结构

<<< @/db/codes/mysql/ddl_desc_tbl.sql

- 查询指定表的建表语句

<<< @/db/codes/mysql/ddl_show_create_tbl.sql

- 创建

<<< @/db/codes/mysql/ddl_create_tbl.sql

- 添加字段

<<< @/db/codes/mysql/ddl_add_field.sql

- 修改数据类型

<<< @/db/codes/mysql/ddl_modify_field.sql

- 修改字段名和字段类型

<<< @/db/codes/mysql/ddl_change_field.sql

- 删除字段

<<< @/db/codes/mysql/ddl_drop_field.sql

- 修改表名

<<< @/db/codes/mysql/ddl_rename_tbl.sql

- 清空表并重置自增计数

<<< @/db/codes/mysql/ddl_truncate_tbl.sql

- 删除表

<<< @/db/codes/mysql/ddl_drop_tbl.sql

`TRUNCATE TABLE` 保留表定义，不逐行执行 `DELETE`，也不触发 `DELETE` 触发器。它与多数建表、改表语句一样会隐式提交，不能用普通事务的 `ROLLBACK` 撤销。InnoDB 的原子 DDL 保证 DDL 在崩溃恢复后的完整性，不表示 DDL 可以随业务事务回滚。参见 [隐式提交](https://dev.mysql.com/doc/refman/8.4/en/implicit-commit.html)。

#### 1.1.3 数据类型

数值类型：

| 类型               | 大小               |
| ------------------ | ------------------ |
| `TINYINT`          | 1 byte             |
| `SMALLINT`         | 2 bytes            |
| `MEDIUMINT`        | 3 bytes            |
| `INT` 或 `INTEGER` | 4 bytes            |
| `BIGINT`           | 8 bytes            |
| `FLOAT`            | 4 bytes            |
| `DOUBLE`           | 8 bytes            |
| `DECIMAL`          | 变长（按精度存储） |

字符串类型：`CHAR`、`VARCHAR`、`TINYBLOB`、`TINYTEXT`、`BLOB`、`TEXT`、`MEDIUMBLOB`、`MEDIUMTEXT`、`LONGBLOB`、`LONGTEXT`

`CHAR(M)`、`VARCHAR(M)` 的 `M` 表示字符数，实际字节数取决于字符集。`VARCHAR` 还受 65,535 字节的行大小限制，不能把这个上限直接当作可存储的字符数。`BLOB` 存储二进制数据，`TEXT` 使用字符集和排序规则。金额通常用定点类型 `DECIMAL`，避免二进制浮点数带来的表示误差。

日期类型：

| 类型        | 大小    | 格式                  |
| ----------- | ------- | --------------------- |
| `DATE`      | 3 bytes | `YYYY-MM-DD`          |
| `TIME`      | 3 bytes | `HH:MM:SS`            |
| `YEAR`      | 1 byte  | `YYYY`                |
| `DATETIME`  | 5 bytes | `YYYY-MM-DD HH:MM:SS` |
| `TIMESTAMP` | 4 bytes | `YYYY-MM-DD HH:MM:SS` |

表中未计小数秒。`TIME`、`DATETIME`、`TIMESTAMP` 支持 0 至 6 位小数秒，分别额外占用 0 至 3 字节。`DATETIME` 的非小数部分从 MySQL 5.6.4 起使用 5 字节。`TIMESTAMP` 在存取时进行会话时区与 UTC 的转换，8.4 的范围上限仍在 2038 年，`DATETIME` 不做这种时区转换。`TIME` 也可表示时长，其范围超过一天。参见[存储空间](https://dev.mysql.com/doc/refman/8.4/en/storage-requirements.html)和[日期时间类型](https://dev.mysql.com/doc/refman/8.4/en/datetime.html)。

### 1.2 DML

先建立练习表：

<<< @/db/codes/mysql/setup_employee.sql

- 给指定字段添加数据

<<< @/db/codes/mysql/dml_insert.sql

- 给全部字段添加数据

<<< @/db/codes/mysql/dml_insert_all.sql

- 批量添加数据

<<< @/db/codes/mysql/dml_insert_batch.sql

- 修改数据

<<< @/db/codes/mysql/dml_update.sql

- 删除数据

<<< @/db/codes/mysql/dml_delete.sql

### 1.3 DQL

#### 1.3.1 基础查询

- 查询多个字段

<<< @/db/codes/mysql/dql_base.sql

- 设置别名

<<< @/db/codes/mysql/dql_as.sql

- 去重

<<< @/db/codes/mysql/dql_distinct.sql

#### 1.3.2 条件查询

| 比较运算符            | 功能                                             |
| --------------------- | ------------------------------------------------ |
| `>`                   | 大于                                             |
| `>=`                  | 大于等于                                         |
| `<`                   | 小于                                             |
| `<=`                  | 小于等于                                         |
| `=`                   | 等于                                             |
| `<>` 或 `!=`          | 不等于                                           |
| `BETWEEN ... AND ...` | 在某个范围内（含最大、最小值）                   |
| `IN (...)`            | 在 `IN` 之后的列表中的值                         |
| `LIKE 占位符`         | 模糊匹配（`_` 匹配单个字符，`%` 匹配任意个字符） |
| `IS NULL`             | 是 `NULL`                                        |

| 逻辑运算符 | 功能                     |
| ---------- | ------------------------ |
| `AND`      | 并且（多个条件同时成立） |
| `OR`       | 或者（任一条件成立）     |
| `NOT`      | 非，不是                 |

使用标准关键字可避免 SQL 模式差异。例如启用 `PIPES_AS_CONCAT` 后，`||` 表示字符串拼接。字符串字面量使用单引号，反引号用于标识符。

与 `NULL` 的普通比较结果是未知值，判断空值应使用 `IS NULL` 或 `IS NOT NULL`。`WHERE` 只保留条件为真的行。尤其注意 `NOT IN` 列表或子查询含 `NULL` 时的结果。

<<< @/db/codes/mysql/dql_where.sql

#### 1.3.3 聚合函数

将一列数据作为一个整体，进行纵向计算。

| 常见函数 | 功能     |
| -------- | -------- |
| `COUNT`  | 统计数量 |
| `MAX`    | 最大值   |
| `MIN`    | 最小值   |
| `AVG`    | 平均值   |
| `SUM`    | 求和     |

<<< @/db/codes/mysql/dql_agg.sql

`COUNT(*)` 统计结果集的行数，`COUNT(列)` 只统计该列非 `NULL` 的行。这里的 `SUM`、`AVG`、`MIN`、`MAX` 忽略 `NULL`，没有非空输入时返回 `NULL`，而 `COUNT` 返回 `0`。

#### 1.3.4 分组查询

<<< @/db/codes/mysql/dql_group.sql

8.4 默认启用 `ONLY_FULL_GROUP_BY`。选择列表中的非聚合列必须出现在 `GROUP BY` 中，或满足 MySQL 可识别的函数依赖等条件。`WHERE` 在分组前筛选行，`HAVING` 筛选分组后的结果。参见[分组规则](https://dev.mysql.com/doc/refman/8.4/en/group-by-handling.html)。

#### 1.3.5 排序查询

<<< @/db/codes/mysql/dql_order.sql

#### 1.3.6 分页查询

<<< @/db/codes/mysql/dql_limit.sql

分页需要明确的 `ORDER BY`，排序列有重复值时再用唯一列打破平局，否则不同执行计划可能返回不同的分页结果。

### 1.4 DCL

#### 1.4.1 用户管理

MySQL 账号由 `'用户名'@'主机'` 共同确定，同名但主机不同的是不同账号。下面的本地练习账号需要由具备相应管理权限的账号创建，示例密码仅供练习。

- 查询用户

<<< @/db/codes/mysql/dcl_select_user.sql

- 创建用户

<<< @/db/codes/mysql/dcl_create_user.sql

- 修改用户密码

<<< @/db/codes/mysql/dcl_change_pwd.sql

8.4 默认使用 `caching_sha2_password`。`mysql_native_password` 从 8.0.34 起弃用，8.4 默认禁用，9.0 移除，不能直接沿用旧版修改密码示例。参见[认证插件](https://dev.mysql.com/doc/refman/8.4/en/native-pluggable-authentication.html)。

- 删除用户，在权限练习结束后执行

<<< @/db/codes/mysql/dcl_drop_user.sql

#### 1.4.2 权限控制

| 常用权限                | 说明                                                 |
| ----------------------- | ---------------------------------------------------- |
| `ALL`，`ALL PRIVILEGES` | 指定范围内的权限集合，不含 `GRANT OPTION` 和 `PROXY` |
| `SELECT`                | 查询数据                                             |
| `INSERT`                | 插入数据                                             |
| `UPDATE`                | 修改数据                                             |
| `DELETE`                | 删除数据                                             |
| `ALTER`                 | 修改表                                               |
| `DROP`                  | 删除数据库/表/视图                                   |
| `CREATE`                | 创建数据库/表                                        |

- 查询权限

<<< @/db/codes/mysql/dcl_show_grants.sql

- 授予权限

<<< @/db/codes/mysql/dcl_grant.sql

- 撤销权限

<<< @/db/codes/mysql/dcl_revoke.sql

`GRANT` 的目标账号必须先存在。`库名.*` 表示该库的对象权限，`*.*` 表示全局范围。通过 `CREATE USER`、`ALTER USER`、`GRANT`、`REVOKE` 修改账号或权限后，无需额外执行 `FLUSH PRIVILEGES`。参见[账号和授权](https://dev.mysql.com/doc/refman/8.4/en/creating-accounts.html)。

## 2. 函数

### 2.1 字符串函数

| 常用函数                     | 说明                                                        |
| ---------------------------- | ----------------------------------------------------------- |
| `CONCAT(S1, S2, ..., Sn)`    | 字符串拼接                                                  |
| `LOWER(str)`                 | 将字符串全部转为小写                                        |
| `UPPER(str)`                 | 将字符串全部转为大写                                        |
| `LPAD(str, n, pad)`          | 在左侧填充至 `n` 个字符，原字符串过长时截断                 |
| `RPAD(str, n, pad)`          | 在右侧填充至 `n` 个字符，原字符串过长时截断                 |
| `TRIM(str)`                  | 去掉字符串头部和尾部的空格                                  |
| `SUBSTRING(str, start, len)` | 返回字符串 `str` 从 `start` 位置起的 `len` 个长度的子字符串 |

<<< @/db/codes/mysql/func_string.sql

`SUBSTRING` 的正起点从 `1` 开始，负起点从末尾计算，起点为 `0` 时返回空串。`CONCAT` 任一参数为 `NULL` 时返回 `NULL`。参见[字符串函数](https://dev.mysql.com/doc/refman/8.4/en/string-functions.html)。

### 2.2 数值函数

| 常用函数      | 说明                                   |
| ------------- | -------------------------------------- |
| `CEIL(x)`     | 向上取整                               |
| `FLOOR(x)`    | 向下取整                               |
| `MOD(x, y)`   | 返回 `x % y`                           |
| `RAND()`      | 返回满足 `0 <= x < 1` 的随机数         |
| `ROUND(x, y)` | 返回 `x` 的四舍五入值，保留 `y` 位小数 |

<<< @/db/codes/mysql/func_math.sql

`ROUND` 对精确值采用中间值远离零的舍入规则，对近似浮点值则依赖底层库，不能统一理解为十进制四舍五入。

### 2.3 日期函数

| 常用函数                             | 说明                                       |
| ------------------------------------ | ------------------------------------------ |
| `CURDATE()`                          | 当前会话时区下的日期                       |
| `CURTIME()`                          | 当前会话时区下的时间                       |
| `NOW()`                              | 当前语句开始时的日期和时间                 |
| `YEAR(date)`                         | 提取年份                                   |
| `MONTH(date)`                        | 提取月份                                   |
| `DAY(date)`                          | 提取日号，等价于 `DAYOFMONTH`              |
| `DATE_ADD(date, INTERVAL expr type)` | 加上指定时间间隔                           |
| `DATEDIFF(date1, date2)`             | 前一个日期减后一个日期的天数，忽略时间部分 |

<<< @/db/codes/mysql/func_date.sql

### 2.4 流程函数

| 常用函数                                                     | 说明                                                                   |
| ------------------------------------------------------------ | ---------------------------------------------------------------------- |
| `IF(value, t, f)`                                            | 如果 `value` 为 `true`，则返回 `t`，否则返回 `f`                       |
| `IFNULL(value1, value2)`                                     | 如果 `value1` 不为 `NULL`，返回 `value1`，否则返回 `value2`            |
| `CASE WHEN [val1] THEN [res1] ... ELSE [default] END`        | 如果 `val1` 为 `true`，返回 `res1`，...，否则返回 `default` 默认值     |
| `CASE [expr] WHEN [val1] THEN [res1] ... ELSE [default] END` | 如果 `expr` 的值为 `val1`，返回 `res1`，...，否则返回 `default` 默认值 |

<<< @/db/codes/mysql/func_flow.sql

## 3. 约束

约束是作用于表中字段上的规则，用于限制存储在表中的数据。

| 约束     | 关键字        | 描述                                                     |
| -------- | ------------- | -------------------------------------------------------- |
| 非空约束 | `NOT NULL`    | 限制该字段数据不能为 `NULL`                              |
| 唯一约束 | `UNIQUE`      | 保证该字段数据唯一、不重复                               |
| 主键约束 | `PRIMARY KEY` | 主键是一行数据的唯一标识，要求非空且唯一                 |
| 默认约束 | `DEFAULT`     | 保存数据时，如果未指定该字段的值，则采用默认值           |
| 检查约束 | `CHECK`       | 保证字段值满足一个条件                                   |
| 外键约束 | `FOREIGN KEY` | 用来让两张表的数据之间建立连接，保证数据的一致性和完整性 |

<<< @/db/codes/mysql/constraint.sql

`UNIQUE` 允许多个 `NULL`，要禁止空值还需 `NOT NULL`。`CHECK` 从 8.0.16 起真正执行检查，结果为假时拒绝写入，为未知值时通过，因此 `CHECK (age >= 0)` 不能替代 `NOT NULL`。示例中故意违反约束的语句已标明，逐条执行可观察相应错误。参见[CHECK 约束](https://dev.mysql.com/doc/refman/8.4/en/create-table-check-constraints.html)。

外键约束要求子表的非空外键值在父表中存在。InnoDB 的关联列需要兼容的数据类型，整数的有无符号属性也要一致。8.4 默认限制引用非唯一键等非标准外键，示例统一引用父表主键。

先建立父表并为练习表添加外键列：

<<< @/db/codes/mysql/setup_foreign_key.sql

- 添加外键

方式一：

<<< @/db/codes/mysql/foreign_key_1.sql

方式二：

<<< @/db/codes/mysql/foreign_key_2.sql

外键的删除/更新行为：

| 行为          | 说明                                                                                                          |
| ------------- | ------------------------------------------------------------------------------------------------------------- |
| `NO ACTION`   | 当在父表中删除/更新对应记录时，首先检查该记录是否有对应外键，如果有则不允许删除/更新。（与 `RESTRICT` 一致）  |
| `RESTRICT`    | 当在父表中删除/更新对应记录时，首先检查该记录是否有对应外键，如果有则不允许删除/更新。（与 `NO ACTION` 一致） |
| `CASCADE`     | 当在父表中删除/更新对应记录时，首先检查该记录是否有对应外键，如果有则也删除/更新外键在子表中的记录。          |
| `SET NULL`    | 父表删除记录或修改被引用键时，将相关子表外键设为 `NULL`，外键列必须允许空值                                   |
| `SET DEFAULT` | InnoDB 不支持，指定该动作的外键定义会被拒绝                                                                   |

上述 `NO ACTION` 与 `RESTRICT` 等价指的是 InnoDB，它不支持将外键检查延迟到事务提交。

<<< @/db/codes/mysql/foreign_key_action.sql

## 4. 多表查询

本节及后续查询示例使用以下练习表。初始化脚本清空并重建 `demo_query`，需要重做练习时再执行。

<<< @/db/codes/mysql/setup_query.sql

### 4.1 内连接

- 隐式内连接

<<< @/db/codes/mysql/join_inner_where.sql

- 显式内连接

<<< @/db/codes/mysql/join_inner_on.sql

### 4.2 外连接

- 左外连接

<<< @/db/codes/mysql/join_left.sql

- 右外连接

<<< @/db/codes/mysql/join_right.sql

### 4.3 自连接

自连接查询可以是内连接，也可以是外连接。

<<< @/db/codes/mysql/join_self.sql

### 4.4 联合查询

关键字：

- `UNION ALL`：直接合并
- `UNION`：去重合并

<<< @/db/codes/mysql/join_union.sql

各查询的列数必须相同，对应列的数据类型要可兼容。`UNION ALL` 和 `UNION` 都不保证输出顺序，需要时在整个联合查询后添加 `ORDER BY`。

### 4.5 子查询

SQL 语句中嵌套 `SELECT` 语句，成为嵌套查询，又称子查询。

#### 4.5.1 标量子查询

子查询返回的结果是单个值（数字、字符串、日期等），这种子查询称为标量子查询。

标量子查询没有结果行时取值为 `NULL`，返回超过一行时报错。

<<< @/db/codes/mysql/subquery_scalar.sql

#### 4.5.2 列子查询

子查询返回的结果是一列，这种子查询称为列子查询。

常配合 `IN`、`ANY`、`ALL` 使用，结果可以有多行。

<<< @/db/codes/mysql/subquery_column.sql

#### 4.5.3 行子查询

行子查询用于将一行中的多个列值与行构造器比较，下面的等值比较要求子查询最多返回一行。

<<< @/db/codes/mysql/subquery_row.sql

#### 4.5.4 表子查询

子查询也可返回多行、多列，用于多列 `IN` 比较，或放在 `FROM` 中形成派生表。派生表需要指定别名。

<<< @/db/codes/mysql/subquery_table.sql

## 5. 事务

事务将相关操作组织成一个提交或回滚的单元。例如转账时，扣款和入账应在同一事务中完成。本节讨论 InnoDB，非事务存储引擎上的写入不能靠 `ROLLBACK` 撤销。

默认 `autocommit=1` 且未显式开启事务时，每条语句独立提交。显式开启事务后，直到 `COMMIT` 或 `ROLLBACK` 才结束这组操作：

```sql
CREATE TABLE account (id INT PRIMARY KEY, balance DECIMAL(12, 2) NOT NULL);
INSERT INTO account VALUES (1, 1000), (2, 1000);

START TRANSACTION;
UPDATE account SET balance = balance - 100 WHERE id = 1;
UPDATE account SET balance = balance + 100 WHERE id = 2;
COMMIT;
```

应用还需检查执行结果并决定提交或回滚。语句报错不一定自动回滚整个事务，例如默认的锁等待超时只回滚当前语句，死锁则回滚被选中的整个事务。参见 [InnoDB 错误处理](https://dev.mysql.com/doc/refman/8.4/en/innodb-error-handling.html)。

- 查看/设置事务的提交方式

<<< @/db/codes/mysql/transaction_auto_commit.sql

- 提交事务

<<< @/db/codes/mysql/transaction_commit.sql

- 回滚事务

<<< @/db/codes/mysql/transaction_rollback.sql

### 5.1 事务操作

- 开启事务

<<< @/db/codes/mysql/transaction_start.sql

或：

<<< @/db/codes/mysql/transaction_begin.sql

- 提交事务

<<< @/db/codes/mysql/transaction_commit.sql

- 回滚事务

<<< @/db/codes/mysql/transaction_rollback.sql

### 5.2 四大特性 ACID

- 原子性（Atomicity）：事务中的修改整体提交或回滚，InnoDB 用 undo log 支持回滚。
- 一致性（Consistency）：事务前后应满足定义的数据约束和业务不变量。数据库约束不能替代应用中的全部业务校验。
- 隔离性（Isolation）：隔离级别规定并发事务之间哪些变化可见，InnoDB 通过 MVCC 和锁实现相应行为。
- 持久性（Durability）：已提交的修改可在崩溃后恢复。保证程度依赖 redo、binlog 的刷盘配置及存储设备是否正确完成持久化。Doublewrite 用于防止数据页部分写入，不替代 redo log。

### 5.3 并发事务问题

| 问题       | 描述                                                           |
| ---------- | -------------------------------------------------------------- |
| 脏读       | 一个事务读取了另一个事务尚未提交的修改数据                     |
| 不可重复读 | 同一事务重复读取同一行时，因其他事务提交修改而读到不同值       |
| 幻读       | 重复执行同一条件查询，因其他事务的变更而出现或消失满足条件的行 |

### 5.4 事务隔离级别

| 隔离级别         |  脏读  | 不可重复读 |   幻读   |
| ---------------- | :----: | :--------: | :------: |
| Read uncommitted |  可能  |    可能    |   可能   |
| Read committed   | 不发生 |    可能    |   可能   |
| Repeatable read  | 不发生 |   不发生   | 标准允许 |
| Serializable     | 不发生 |   不发生   |  不发生  |

上表用于说明隔离级别允许的典型现象。InnoDB 默认是 `REPEATABLE READ`，一致性读复用快照，锁定读和写操作通常用 next-key lock 阻止范围内插入。快照读与锁定读看到的数据可能不同，同一事务也能看到自己的写入，不能笼统地说所有 RR 查询都读取同一份数据。参见 [InnoDB 隔离级别](https://dev.mysql.com/doc/refman/8.4/en/innodb-transaction-isolation-levels.html)。

- 查看事务隔离级别

<<< @/db/codes/mysql/transaction_isolation_select.sql

- 设置事务隔离级别

<<< @/db/codes/mysql/transaction_isolation_set.sql

`SESSION` 设置当前会话后续事务的默认隔离级别，不改变正在执行的事务。`GLOBAL` 设置后续新连接的默认值，不改变已有连接。省略作用域的 `SET TRANSACTION` 只针对下一次事务。

## 6. 存储引擎

- 查询当前数据库支持的存储引擎

<<< @/db/codes/mysql/engine_show.sql

以下是常见输出，支持状态以当前实例的编译和安装配置为准。

|       Engine       | Support | Comment                                                        | Transactions |  XA   | Savepoints |
| :----------------: | :-----: | -------------------------------------------------------------- | :----------: | :---: | :--------: |
|      ARCHIVE       |   YES   | Archive storage engine                                         |      NO      |  NO   |     NO     |
|     BLACKHOLE      |   YES   | /dev/null storage engine (anything you write to it disappears) |      NO      |  NO   |     NO     |
|     MRG_MYISAM     |   YES   | Collection of identical MyISAM tables                          |      NO      |  NO   |     NO     |
|     FEDERATED      |   NO    | Federated MySQL storage engine                                 |     NULL     | NULL  |    NULL    |
|       MyISAM       |   YES   | MyISAM storage engine                                          |      NO      |  NO   |     NO     |
| PERFORMANCE_SCHEMA |   YES   | Performance Schema                                             |      NO      |  NO   |     NO     |
|       InnoDB       | DEFAULT | Supports transactions, row-level locking, and foreign keys     |     YES      |  YES  |    YES     |
|       MEMORY       |   YES   | Hash based, stored in memory, useful for temporary tables      |      NO      |  NO   |     NO     |
|        CSV         |   YES   | CSV storage engine                                             |      NO      |  NO   |     NO     |

- 在创建表时指定存储引擎

<<< @/db/codes/mysql/engine_create_table.sql

### 6.1 InnoDB

InnoDB 是支持事务、行级锁和崩溃恢复的通用存储引擎，从 MySQL 5.5 起成为默认存储引擎。

特点：

- DML 操作遵循 ACID 模型，支持事务
- 行级锁，提高并发访问性能
- 支持外键 `FOREIGN KEY` 约束，保证数据的完整性和一致性

开启 `innodb_file_per_table` 后，新建的普通非分区表默认使用独立的 `tb_name.ibd`，存储数据、索引和序列化字典信息（SDI）。表也可以放在通用表空间等共享表空间中，分区表则有多个分区文件，不能把“一张表一个 `.ibd`”作为所有表的规则。

创建新表时是否默认使用独立表空间，可通过以下命令确认。该设置的变化不会自动搬迁已有表：

<<< @/db/codes/mysql/engine_var_idb.sql

数据目录由 `datadir` 决定，可以执行 `SELECT @@datadir` 查询。部分 Linux 发行版的软件包使用 `/var/lib/mysql`，这不是所有部署方式的固定路径。

MySQL 8.0 起不再使用 `.frm` 文件保存表定义，事务数据字典存储于 `mysql.ibd`。`ibd2sdi 表名.ibd` 提取的是 SDI 元数据，不是把表内的业务行转换成可读文本。

逻辑存储结构：

1. Tablespace：表空间，一个表空间包含多个段
2. Segment：段，管理分配给该段的零散页和区
3. Extent：区，默认 1 MiB，包含 64 个 16 KiB 页
4. Page：页，默认 16 KiB，是缓存和磁盘读写的基本单位
5. Row：行，较长的变长字段可能使用页外存储

页大小在实例初始化时确定。4、8、16 KiB 页对应 1 MiB 的区，32、64 KiB 页对应 2、4 MiB 的区。

### 6.2 MyISAM

MyISAM 是 MySQL 早期的默认存储引擎。

特点：

- 不支持事务，不支持外键
- 支持表锁，不支持行锁
- 适用于不需要事务和外键的特定场景，性能需根据工作负载测量

文件：

- `tb_name_<id>.sdi`：存储表结构的 SDI，文件名含内部表 ID
- `tb_name.MYD`：存储数据
- `tb_name.MYI`：存储索引

### 6.3 MEMORY

MEMORY 表将数据放在内存中，服务器重启后表定义仍在，数据会丢失。它适合可重新生成的数据，但不能与 MySQL 执行查询时使用的内部临时表混为一谈，8.4 的内部临时表通常使用 TempTable 或 InnoDB。

特点：

- 内存存放
- Hash 索引（默认）
- 可以显式建立 BTREE 索引，支持范围查找

文件：`tb_name_<id>.sdi` 存储表结构的 SDI，不保存内存中的业务数据。

### 6.4 对比

|    特点    |       InnoDB       | MyISAM | MEMORY |
| :--------: | :----------------: | :----: | :----: |
|    事务    |        支持        |   -    |   -    |
|    外键    |        支持        |   -    |   -    |
| 主要数据锁 | 行级锁及表级意向锁 |  表锁  |  表锁  |
| BTREE 索引 | 支持，内部为 B+ 树 |  支持  |  支持  |
| Hash 索引  |         -          |   -    |  支持  |

这里的 Hash 索引指用户显式创建的索引，InnoDB 的自适应哈希索引是另一种内部机制。参见 [MEMORY 引擎](https://dev.mysql.com/doc/refman/8.4/en/memory-storage-engine.html)。

## 7. 索引

### 7.1 结构

InnoDB 的普通索引采用 B+ 树。MySQL 还支持 MEMORY 的 Hash 索引、空间索引的 R-tree、全文索引等，索引结构取决于引擎与索引类型。

| 特性             | B 树                                | B+ 树                                         |
| ---------------- | ----------------------------------- | --------------------------------------------- |
| **存储内容**     | 非叶子节点和叶子节点均存储键+值     | 仅叶子节点存储键+值，非叶子节点仅存键（索引） |
| **叶子节点**     | 一般不依靠叶子链表组织数据          | 叶子节点按键顺序相连，适合顺序访问            |
| **查找路径**     | 可能在非叶子节点找到数据            | 查找数据需要到达叶子节点                      |
| **范围查询**     | 需沿树结构做有序遍历                | 定位首个叶子记录后，可沿叶子链继续遍历        |
| **节点存储密度** | 低（键+值占用空间大，单节点键数少） | 高（非叶子节点仅存键，单节点键数更多）        |
| **树高**         | 取决于节点大小、数据大小与填充率    | 非叶子节点不存行数据，通常能容纳更多分隔键    |
| **数据冗余**     | 无冗余（键仅出现一次）              | 有冗余（非叶子节点的键是叶子节点的副本）      |

### 7.2 分类

| 分类         | 含义                                                 | 特点                     | 关键字     |
| ------------ | ---------------------------------------------------- | ------------------------ | ---------- |
| **主键索引** | 针对于表中主键创建的索引                             | 默认自动创建, 只能有一个 | `PRIMARY`  |
| **唯一索引** | 避免同一个表中某数据列中的值重复                     | 可以有多个               | `UNIQUE`   |
| **常规索引** | 快速定位特定数据                                     | 可以有多个               |            |
| **全文索引** | 全文索引查找的是文本中的关键词，而不是比较索引中的值 | 可以有多个               | `FULLTEXT` |

在 InnoDB 存储引擎中，根据索引的存储形式，又可以分为以下两种：

| 分类                            | 含义                                                       | 特点                 |
| ------------------------------- | ---------------------------------------------------------- | -------------------- |
| **聚簇索引（Clustered Index）** | 将数据存储和索引放到了一块，索引结构的叶子节点保存了行数据 | 必须要，而且只有一个 |
| **二级索引（Secondary Index）** | 叶子记录保存索引列及聚簇键，需要其他列时再查聚簇索引       | 可以存在多个         |

聚簇索引选取规则：

- 如果存在主键，主键索引就是聚簇索引
- 如果不存在主键，选择第一个所有列均为 `NOT NULL` 的唯一索引
- 如果没有符合条件的索引，生成 6 字节的 `DB_ROW_ID`，建立隐藏聚簇索引 `GEN_CLUST_INDEX`

### 7.3 语法

本节及后续索引示例使用 `demo_query.tb_user`，同名索引只需创建一次。小型练习表可能更适合全表扫描，实际计划未使用索引并不等于索引定义无效。

- 创建索引

<<< @/db/codes/mysql/index_create.sql

- 查看索引

<<< @/db/codes/mysql/index_show.sql

- 删除索引

<<< @/db/codes/mysql/index_drop.sql

### 7.4 SQL 性能分析

#### 7.4.1 SQL 执行频率

`SHOW SESSION STATUS` 查看当前会话状态，`SHOW GLOBAL STATUS` 查看实例累计状态。`Com_*` 是命令执行次数，不是当前数据库的访问次数，也不是影响的行数。比较一段时间内计数的差值，可了解相应命令的执行频率。

<<< @/db/codes/mysql/status_show.sql

#### 7.4.2 慢查询日志

慢查询日志通常记录执行时间超过 `long_query_time` 且检查行数至少为 `min_examined_row_limit` 的语句。`long_query_time` 默认 10 秒，管理语句等是否记录还受其他选项控制，不能只看耗时条件。

MySQL 的慢查询日志默认没有开启，需要在 MySQL 的配置文件（Ubuntu 默认在 `/etc/mysql/mysql.conf.d/mysqld.cnf`）配置如下信息：

<<< @/db/codes/mysql/slow_query_log.cnf

文件路径需适合当前安装方式，并允许服务进程写入。修改全局 `long_query_time` 后，已有连接的会话值不会自动变化。参见[慢查询日志](https://dev.mysql.com/doc/refman/8.4/en/slow-query-log.html)。

#### 7.4.3 语句耗时

8.4 优先通过 Performance Schema 和 `sys` 视图分析语句耗时。下面按总耗时查看规范化语句的统计，统计范围是实例而不是当前连接：

<<< @/db/codes/mysql/statement_analysis.sql

旧的 `SHOW PROFILE`、`SHOW PROFILES` 在 8.4 仍可用，但已弃用，不适合作为新的诊断工具。需要阅读旧代码时，可先检查支持情况：

<<< @/db/codes/mysql/have_profiling.sql

`profiling` 只支持会话作用域，默认关闭：

<<< @/db/codes/mysql/set_profiling.sql

执行一系列的业务 SQL 操作，然后通过如下指令查看执行耗时：

<<< @/db/codes/mysql/profile_show.sql

#### 7.4.4 `EXPLAIN` 执行计划

`EXPLAIN` 或者 `DESC` 命令获取 MySQL 如何执行 `SELECT` 语句的信息，包括在 `SELECT` 语句执行过程中表如何连接和连接的顺序。

`EXPLAIN FORMAT=TRADITIONAL` 返回下述表格格式，`FORMAT=TREE` 或 `FORMAT=JSON` 提供不同的计划细节。`EXPLAIN ANALYZE` 从 8.0.18 起支持，它会实际执行查询并报告测量结果，不能把它当作只查看静态计划的命令。

<<< @/db/codes/mysql/explain.sql

`EXPLAIN` 执行计划各字段含义：

- **`id`**：查询块的标识，不能单凭数值大小推断完整的执行顺序
- **`select_type`**：表示 `SELECT` 的类型，常见取值：
  - `SIMPLE`：不含 `UNION` 或子查询的简单查询，可以包含表连接
  - `PRIMARY`：主查询，即外层的查询
  - `UNION`：`UNION` 中的第二个或者后面的查询语句
  - `SUBQUERY`：`SELECT` / `WHERE` 之后包含了子查询
- **`type`**：访问或连接类型，常见值如下，不能脱离数据量和总成本给它们做固定性能排名：
  1. `NULL`：不涉及表（如 `SELECT 1+1`）
  2. `system`：表中只有一行数据（`const` 的特例，如系统表）
  3. `const`：通过主键/唯一索引匹配到单行
  4. `eq_ref`：多表连接中，被连接表通过唯一索引匹配
  5. `ref`：通过普通索引匹配多行
  6. `range`：索引范围扫描（如 `BETWEEN`、`IN`、`>` 等）
  7. `index`：完整索引扫描，是否还需要读取完整行取决于索引是否覆盖
  8. `ALL`：全表扫描，小表或返回大部分行时可能是合理选择
- **`possible_keys`**：显示可能应用在这张表上的索引，一个或多个
- **`key`**：实际使用的索引，如果为 `NULL`，则没有使用索引
- **`key_len`**：计划使用的索引键长度，受类型、字符集、可空属性等影响，可辅助判断使用了哪些键部分
- **`rows`**：MySQL 认为必须要执行查询的行数，在 InnoDB 引擎的表中，是一个估计值，不一定准确
- **`filtered`**：估计的条件过滤后保留比例，`rows × filtered / 100` 可估算传给后续连接的行数，数值大不代表计划一定更好
- **`Extra`**：补充信息，例如覆盖索引、索引条件下推、临时表或额外排序

参见 [EXPLAIN 输出](https://dev.mysql.com/doc/refman/8.4/en/explain-output.html)。

### 7.5 使用规则

#### 7.5.1 最左前缀原则

联合索引 `(profession, age, status)` 可用连续的最左列构造常规查找范围。跳过 `age` 后，仍可按 `profession` 定位，再过滤 `status`，并不是整个索引失去作用。后续列还可能用于索引条件下推或覆盖查询。8.0.13 起提供的 Skip Scan 在满足条件时也可能使用缺少首列条件的索引。

<<< @/db/codes/mysql/combined_index.sql

#### 7.5.2 范围查询

构造多列 BTREE 范围时，遇到非等值范围条件后，通常不能继续用后续列完整收窄所有扫描区间。后续列仍可能参与边界处理、索引条件下推或覆盖查询。将 `>` 改成 `>=` 不保证三个字段都用于定位，需检查实际范围与过滤步骤。

<<< @/db/codes/mysql/combined_index_range.sql

参见[范围优化](https://dev.mysql.com/doc/refman/8.4/en/range-optimization.html)和[索引条件下推](https://dev.mysql.com/doc/refman/8.4/en/index-condition-pushdown-optimization.html)。

#### 7.5.3 难以使用索引定位的情况

以下情况可能妨碍普通索引构造高效的查找范围，是否改用全表扫描还由优化器成本判断决定：

- **索引字段参与运算**

  如 `WHERE id + 1 = 10` 或 `WHERE SUBSTR(name, 1, 3) = 'abc'`。8.0.13 起支持函数索引，表达式与相应索引匹配时可以使用，不能把所有函数条件都归为不可索引。

- **字符串不加引号**

  如用数值与字符串列比较的 `WHERE name = 123`，转换语义不能直接对应字符串索引的排序，应按字段类型写为 `WHERE name = '123'`。

- **`LIKE` 以通配符开头**

  如 `WHERE name LIKE '%abc'` 通常不能利用普通 BTREE 构造前缀范围，但仍可能扫描覆盖索引。`LIKE 'abc%'` 则可形成前缀范围。

- **使用 `OR` 连接非索引字段**

  如 `WHERE idx_col = 1 OR no_idx_col = 2`，未索引的分支通常妨碍 Index Merge 的联合访问。两个分支都有合适索引时可能使用 Index Merge，是否采用仍需看成本。

- **数据分布影响**

  如果 MySQL 评估使用索引比全表更慢，则不使用索引

#### 7.5.4 SQL 提示

索引提示可以限制优化器的候选索引。先确认统计信息和计划，再比较实际耗时，不能认为加入提示一定更快。

**`USE INDEX`**：

<<< @/db/codes/mysql/index_use.sql

**`IGNORE INDEX`**：

<<< @/db/codes/mysql/index_ignore.sql

**`FORCE INDEX`**：

<<< @/db/codes/mysql/index_force.sql

这些传统语法在 8.4 仍支持。8.4 还可使用 `INDEX`、`JOIN_INDEX`、`GROUP_INDEX`、`ORDER_INDEX` 等优化器提示，具体适用范围见 [索引提示](https://dev.mysql.com/doc/refman/8.4/en/index-hints.html)。

#### 7.5.5 覆盖索引

尽量使用覆盖索引（查询使用了索引，并且需要返回的列在该索引中已经全部能够找到），减少 `SELECT *`。

InnoDB 二级索引隐含聚簇键，因此二级索引覆盖的列不仅是定义时写出的列。`Extra` 中的 `Using index` 表示覆盖访问，不表示排序一定无需 filesort。

#### 7.5.6 前缀索引

当字段类型为字符串（`VARCHAR`、`TEXT` 等）时，有时候需要索引很长的字符串，这会让索引变得更大，查询时，浪费大量磁盘 IO，影响查询效率。此时可以只将字符串的一小部分前缀，建立索引，这样大大节约索引空间，从而提高索引效率。

<<< @/db/codes/mysql/index_prefix.sql

可比较不同前缀长度的去重比例，再结合查询模式、索引大小和写入成本选择长度。比例高意味着重复值较少，不足以单独推出查询性能最好。字符列的前缀长度按字符计算，二进制列按字节计算。

<<< @/db/codes/mysql/index_prefix_len.sql

前缀索引不能覆盖查询所需的完整字段值。若设为唯一前缀索引，约束的也是前缀，两个完整字符串不同但前缀相同的值仍会冲突。

## 8. SQL 优化

### 8.1 `INSERT` 优化

#### 8.1.1 小批量插入数据

下面比较同一张空表上的三种插入方式，每种方式分别执行。先建立练习表，重做时用 `TRUNCATE TABLE tb_insert` 清空：

```sql
USE demo_query;
CREATE TABLE tb_insert (id INT PRIMARY KEY, name VARCHAR(50)) ENGINE=InnoDB;
```

<<< @/db/codes/mysql/optimize_insert_old.sql

优化 1：改为批量插入

<<< @/db/codes/mysql/optimize_insert_multi.sql

优化 2：手动提交事务

自动提交且未显式开启事务时，多条插入语句分别提交。适度批量插入或将一批写入放入同一事务，可以减少通信和提交开销，但过大的批次也会增加锁持有时间、日志量与失败重试成本。

<<< @/db/codes/mysql/optimize_insert_transaction.sql

优化 3：主键顺序插入

```txt
主键乱序插入：8 1 9 21 88 2 4 15 89 5 7 3
主键顺序插入：1 2 3 4 5 7 8 9 15 21 88 89
```

#### 8.1.2 大批量插入数据

导入大量文本数据时，可比较多行 `INSERT` 与 `LOAD DATA` 的实际表现。以下示例将客户端文件中的两列导入 `tb_insert`，不与前面的插入示例同时执行。先准备无表头的 `/tmp/tb_insert.csv`，每行格式为 `id,name`。

<<< @/db/codes/mysql/load_infile.sh
<<< @/db/codes/mysql/load_infile.sql

`LOCAL` 从客户端读取文件，客户端和服务器都需允许本地导入。没有 `LOCAL` 时由服务器读取文件，需要 `FILE` 权限，并受 `secure_file_priv` 限制。参见 [LOAD DATA](https://dev.mysql.com/doc/refman/8.4/en/load-data.html)。

### 8.2 主键优化

在 InnoDB 存储引擎中，表数据都是根据主键顺序组织存放的，这种存储方式的表称为<span style="color:red;">索引组织表</span>（Index Organized Table，<span style="color:red;">IOT</span>）。

#### 8.2.1 页分裂

插入目标页空间不足时，可能分配新页并重新分配记录。聚簇索引按聚簇键排序，随机键插入通常更容易分散到不同页。每页能容纳多少记录取决于记录大小、页头开销和填充情况，没有通用的“至少两行”规则。

#### 8.2.2 页合并

InnoDB 通常先将记录标记为删除，在不再需要它支持回滚或一致性读后，由 purge 清理。不能把标记删除理解为该空间立即可复用。

删除或缩短记录使页的填充比例低于 `MERGE_THRESHOLD` 时，InnoDB 尝试与相邻页合并。默认阈值为 50%，指剩余填充比例，不是被删除记录数量的比例。参见[页合并阈值](https://dev.mysql.com/doc/refman/8.4/en/index-page-merge-threshold.html)。

::: tip `MERGE_THRESHOLD`
合并页的阈值，可以自己设置，在创建表或者创建索引时指定。
:::

#### 8.2.3 主键设计原则

- 满足业务需求的情况下，尽量降低主键长度
- 插入数据时，尽量顺序插入，选择 `AUTO_INCREMENT` 自增主键
- 比较候选键的宽度、顺序性和业务稳定性，随机 UUID 与较长自然键可能增加索引空间和随机写入
- 业务操作时，避免对主键的修改

### 8.3 `ORDER BY` 优化

1. `Using filesort`：通过表的索引或全表扫描，读取满足条件的数据行，然后在排序缓冲区 sort buffer 中完成排序操作，所有不是通过索引直接返回排序结果的排序都叫 FileSort 排序
2. 索引顺序满足 `ORDER BY` 时，可以省去额外排序。`Using index` 本身只说明覆盖访问，不能据此判断排序方式

filesort 可以全部在内存中完成，名称中的 file 不表示一定写磁盘。以下注释说明可能的计划，需要结合当前数据和统计信息确认。

<<< @/db/codes/mysql/optimize_order.sql

- 根据排序字段建立合适的索引，多字段排序时，也遵守最左前缀法则
- 尽量使用覆盖索引
- 多字段排序，一个升序一个降序，此时需要注意联合索引在创建时的规则（`ASC` / `DESC`）
- 多字段混合升降序可使用方向匹配的索引，InnoDB 从 8.0 起支持实际的降序索引
- `sort_buffer_size` 默认 256 KiB，按会话分配。应根据排序规模和并发量测量，不能一律调大

  <<< @/db/codes/mysql/sort_buffer_size.sql

参见 [ORDER BY 优化](https://dev.mysql.com/doc/refman/8.4/en/order-by-optimization.html)。

### 8.4 `GROUP BY` 优化

<<< @/db/codes/mysql/optimize_group.sql

- 在分组操作时，可以通过索引来提升效率
- 分组操作时，索引的使用也是满足最左前缀原则的

合适的索引可能减少临时表或启用松散索引扫描，但要满足具体查询条件。8.0 起 `GROUP BY` 不再隐式排序，需要有序输出时显式加 `ORDER BY`。

### 8.5 `LIMIT` 优化

对于按索引扫描的 `ORDER BY id LIMIT 2000000, 10`，仍需跳过前 2,000,000 个匹配项，偏移量越大，扫描成本通常越高。下列大偏移量仅用于说明问题，练习表数据不足时会返回空结果。

优化思路 1：先从覆盖索引取得这一页的主键，再读取完整行，减少回表次数，但仍需扫描并跳过偏移部分

<<< @/db/codes/mysql/optimize_limit_cover.sql

优化思路 2：记录上一页的最后一个排序键，用范围条件继续读取，避免深偏移。它适合顺序翻页，不能直接替代任意页码跳转

<<< @/db/codes/mysql/optimize_limit_where.sql

### 8.6 `COUNT` 优化

- MyISAM 保存总行数，单表、无 `WHERE` 等符合条件的 `COUNT(*)` 可以直接取得该值
- InnoDB 需要统计当前事务可见的记录，通常扫描最小的二级索引，没有二级索引时扫描聚簇索引，不一定读取每行的全部数据

频繁获取大表总数时，可考虑统计缓存或维护计数表。精确计数需要和业务修改保持事务一致性，缓存则要明确允许的延迟。

#### 8.6.1 `COUNT` 的几种用法

`COUNT()` 是一个聚合函数，对于返回的结果集，一行行地判断，如果参数不为 `NULL`，累计值加一，最后返回累计值。

用法：

- `COUNT(主键)`：主键不为空，语义上等于统计结果行数
- `COUNT(字段)`：统计该字段不为 `NULL` 的行数
- `COUNT(1)`：常量 `1` 不为空，统计结果行数
- `COUNT(*)`：直接统计结果行数，通常最清楚地表达意图

InnoDB 对 `COUNT(*)` 与 `COUNT(1)` 的处理相同，其他形式的性能受可用索引、可空性及查询条件影响，不适合写成固定排名。

### 8.7 `UPDATE` 优化

InnoDB 的记录锁加在索引记录上。`UPDATE` 不能用合适的索引缩小扫描范围时，可能锁住大量记录或区间，增加并发等待，但不会自动将行锁升级成表锁。实际锁范围还取决于隔离级别，RC 会释放不匹配行上的记录锁。

## 9. 视图 / 存储过程 / 触发器

### 9.1 视图

视图（View）是由查询定义的虚拟表，结果来自查询引用的基表或其他视图。

视图保存查询定义，不持久保存查询结果。执行时可合并到外层查询，也可能使用临时表，具体取决于视图定义和执行计划。并非所有视图都可更新，含聚合、分组等操作的视图通常不可更新。

- 创建

<<< @/db/codes/mysql/view_create.sql

- 查询

<<< @/db/codes/mysql/view_query.sql

- 修改

<<< @/db/codes/mysql/view_modify.sql

- 删除

<<< @/db/codes/mysql/view_drop.sql

### 9.2 存储过程

存储过程是保存在服务器端的程序，通过 `CALL` 调用。它可以集中处理相关 SQL，减少应用与数据库之间的通信次数。

特点：

- 支持 `IN`、`OUT`、`INOUT` 参数，可以返回结果集
- 服务器按会话缓存解析后的程序，相关元数据变化时可能重新解析，不能理解为永久缓存一个固定执行计划
- 是否放入数据库应根据事务边界、部署与调试方式、迁移需求决定

#### 9.2.1 基本语法

- 创建

<<< @/db/codes/mysql/procedure_create.sql

`DELIMITER` 是 `mysql` 客户端命令，用于避免客户端在过程体内的分号处提前结束输入。它不是服务器 SQL，通过驱动执行创建语句时不应把 `DELIMITER` 一起发送。下面的完整过程示例均在定义结束后恢复 `;`。

- 调用

<<< @/db/codes/mysql/procedure_call.sql

- 查看

<<< @/db/codes/mysql/procedure_query.sql

- 删除

<<< @/db/codes/mysql/procedure_drop.sql

#### 9.2.2 系统变量

系统变量是 MySQL 服务器提供，不是用户定义的，属于服务器层面。分为全局变量（`GLOBAL`）、会话变量（`SESSION`，默认）。

不同变量支持的作用域、是否可动态修改并不相同。多数会话变量在建立连接时继承全局值，修改全局值不等于修改全部已有会话。

- 查看系统变量

<<< @/db/codes/mysql/variables_query.sql

- 设置系统变量

<<< @/db/codes/mysql/variables_set.sql

#### 9.2.3 用户定义变量

用户定义变量是用户根据需要自己定义的变量，用户变量不用提前声明，在用的时候直接用 `@变量名` 使用就可以。其作用域为当前连接。

未初始化时值为 `NULL`。建议通过 `SET @变量名 = 值` 赋值，避免依赖同一表达式中读写变量的未定义求值顺序。

#### 9.2.4 局部变量

局部变量是根据需要定义的在局部生效的变量，访问之前，需要 `DECLARE` 声明。可用作存储过程内的局部变量和输入参数，局部变量的范围是在其声明的 `BEGIN ... END` 块。

- 声明

<<< @/db/codes/mysql/variables_declare.sql

- 赋值

<<< @/db/codes/mysql/variables_local_set.sql

#### 9.2.5 `IF`

<<< @/db/codes/mysql/procedure_if.sql

#### 9.2.6 `WHILE`

<<< @/db/codes/mysql/procedure_while.sql

#### 9.2.7 `REPEAT`

`REPEAT` 是有条件的循环控制语句，当满足条件的时候退出循环。

它先执行循环体，再检查退出条件，至少执行一次。下列求和过程将 `NULL`、负数作为非法输入，`0` 返回 `0`。

<<< @/db/codes/mysql/procedure_repeat.sql

#### 9.2.8 `LOOP`

`LOOP` 实现简单的循环，如果不在 SQL 逻辑中增加退出循环的条件，可以用其来实现简单的死循环。`LOOP` 可以配合以下两个语句使用：

- `LEAVE`：配合循环使用，退出循环
- `ITERATE`：必须用在循环中，作用是跳过当前循环剩下的语句，直接进入下一次循环

<<< @/db/codes/mysql/procedure_loop.sql

### 9.3 触发器

触发器是与数据库表有关的数据对象，指在 `INSERT` / `UPDATE` / `DELETE` 之前或之后，触发并执行触发器中定义的 SQL 语句集合。触发器的这种特性可以协助应用在数据库端确保数据的完整性、日志记录、数据校验等操作。

使用别名 `OLD` 和 `NEW` 来引用触发器中发生变化的记录内容，这与其他的数据库是相似的。现在触发器还只支持行级触发，不支持语句级触发。

| 触发器类型        | `OLD` 和 `NEW`                                           |
| ----------------- | -------------------------------------------------------- |
| `INSERT` 型触发器 | `NEW` 表示将要或已经新增的数据                           |
| `UPDATE` 型触发器 | `OLD` 表示修改之前的数据，`NEW` 表示将要或已经修改的数据 |
| `DELETE` 型触发器 | `OLD` 表示将要或已经删除的数据                           |

只有 `BEFORE` 触发器可通过 `SET NEW.列 = 值` 改写即将写入的值，`OLD` 始终只读。

## 10. 锁

按锁的粒度可分为全局锁、表级锁、行级锁。

### 10.1 全局锁

`FLUSH TABLES WITH READ LOCK` 获取全局读锁，在持锁期间阻止表的数据更新及相关提交、DDL 等操作，普通查询仍可执行。获取锁本身也可能等待正在执行的语句。

典型使用场景是做全库的逻辑备份，对所有表进行锁定，从而获取一致性视图，保证数据完整性。

1. 在交互式会话 A 中加锁，并保持连接

<<< @/db/codes/mysql/lock_global_lock.sql

2. 在另一个终端 B 中备份

<<< @/db/codes/mysql/mysqldump.sh

3. 回到会话 A 解锁

<<< @/db/codes/mysql/lock_global_unlock.sql

数据库中加全局锁，是一个比较重的操作，存在以下问题：

1. 在主库持有全局读锁期间，业务写入会等待锁释放。
2. 在副本持有全局读锁期间，复制事件的应用可能等待，增加复制延迟。

锁随连接结束而释放，不能用 `mysql -e 'FLUSH TABLES WITH READ LOCK'` 执行后，再在另一个连接中假定锁仍然存在。

对 InnoDB 数据，可以使用 `--single-transaction` 取得一致性快照，避免在整个备份期间持有全局读锁。它不保证 MyISAM、MEMORY 表的一致性，备份期间也需避免对被备份表执行可能破坏快照读取的 DDL。某些 GTID 或坐标相关选项仍可能需要短暂锁定，不能概括成所有配置下都完全不加锁。参见 [mysqldump](https://dev.mysql.com/doc/refman/8.4/en/mysqldump.html)。

<<< @/db/codes/mysql/mysqldump_single_transaction.sh

### 10.2 表级锁

表级锁覆盖整张表，具体用途与兼容关系取决于锁类型。MyISAM 和 MEMORY 主要通过表锁控制数据访问，InnoDB 也存在显式表锁、MDL 和意向锁。

对于表级锁，主要分为以下三类：

1. 表锁
2. 元数据锁（Meta Data Lock，MDL）
3. 意向锁

#### 10.2.1 表锁

对于表锁，分为两类：

1. 表共享读锁（read lock）。我加了读锁以后，我能读不能写，别人也是能读不能写。
2. 表独占写锁（write lock）。我加了写锁以后，我能读能写，别人既不能读也不能写。

语法：

1. 加锁：`LOCK TABLES 表名 ... READ/WRITE`。
2. 释放锁：`UNLOCK TABLES` / 客户端断开连接。

`LOCK TABLES` 会隐式提交当前事务，持锁连接访问的表也必须在锁定列表中列出，不宜直接夹在普通 `START TRANSACTION` 流程中使用。参见 [LOCK TABLES](https://dev.mysql.com/doc/refman/8.4/en/lock-tables.html)。

#### 10.2.2 元数据锁

访问表时系统自动获取相应的 MDL，用于协调表结构变更与数据访问，避免执行期间依赖的表定义被不兼容地修改。

MDL 从 MySQL 5.5 起提供。读写通常获得兼容的元数据锁，DDL 在某些阶段需要更强的锁。即使使用在线 DDL，也不意味着整个过程不需要 MDL。事务访问过的表通常要到事务结束后才释放相应元数据锁，长事务因此可能阻塞 DDL。

查看元数据锁：

<<< @/db/codes/mysql/mdl_select.sql

#### 10.2.3 意向锁

为了避免 DML 在执行时加的行锁与表锁冲突，在 InnoDB 中引入了意向锁，使得表锁不用检查每行数据是否加锁，使用意向锁来减少表锁的检查。

1. 意向共享锁（IS）：表明事务准备在表内获得共享记录锁，例如 `SELECT ... FOR SHARE`
2. 意向排他锁（IX）：表明事务准备在表内获得排他记录锁，例如 `SELECT ... FOR UPDATE`

IS 与 IX 彼此兼容，两个事务持有 IX 不代表它们的记录锁一定冲突。下面是 InnoDB 表级锁模式的兼容关系，S/X 指整表共享锁/排他锁：

| 已有锁 / 请求锁 | IS   | IX   | S    | X    |
| --------------- | ---- | ---- | ---- | ---- |
| IS              | 兼容 | 兼容 | 兼容 | 冲突 |
| IX              | 兼容 | 兼容 | 冲突 | 冲突 |
| S               | 兼容 | 冲突 | 兼容 | 冲突 |
| X               | 冲突 | 冲突 | 冲突 | 冲突 |

查看意向锁：

<<< @/db/codes/mysql/lock_raw_select.sql

### 10.3 行级锁

行级锁在较小范围内控制并发访问，可减少无关记录之间的冲突。热点记录上的竞争仍可能限制并发，不能仅凭锁粒度判断整体性能。

InnoDB 的记录锁实际作用在索引记录上，没有显式索引时也会使用隐藏的聚簇索引。常见的数据锁包括以下三类：

1. 记录锁（Record Lock）：锁定索引记录，分为共享和排他模式，RC、RR 等隔离级别均使用。
2. 间隙锁（Gap Lock）：阻止其他事务向指定间隙插入，不锁定间隙边界上的已有记录，RR 下常用于保护查询范围。
3. 临键锁（Next-Key Lock）：记录锁与该记录前方间隙锁的组合，RR 的锁定扫描通常使用它。

#### 10.3.1 行锁

InnoDB 实现了以下两种类型的行锁：

1. 共享锁（S）：允许一个事务去读一行，阻止其他事务获得相同数据集的排他锁。
2. 排他锁（X）：允许获取排他锁的事务更新数据，阻止其他事务获得相同数据集的共享锁和排他锁。

| SQL                           | 行锁类型   | 说明                                    |
| ----------------------------- | ---------- | --------------------------------------- |
| `INSERT ...`                  | 排他锁     | 自动加锁                                |
| `UPDATE ...`                  | 排他锁     | 自动加锁                                |
| `DELETE ...`                  | 排他锁     | 自动加锁                                |
| 普通 `SELECT`，RC/RR 一致性读 | 不加记录锁 | 仍获取 MDL，并受视图可见性规则约束      |
| `SELECT ... FOR SHARE`        | 共享锁     | 锁定读，旧语法为 `LOCK IN SHARE MODE`   |
| `SELECT ... FOR UPDATE`       | 排他锁     | 需要手动在 `SELECT` 之后加 `FOR UPDATE` |

RR 下的锁定读和写操作通常通过记录锁及临键锁保护访问范围。

1. 针对唯一索引进行检索时，对已存在的记录进行等值匹配时，将会自动优化为行锁。
2. 上述唯一匹配需要覆盖唯一索引的全部列。没有合适索引时，锁定扫描可能涉及大量记录和间隙，但 InnoDB 不会自动做行锁到表锁的升级。

需要在多条语句间保持锁时，应在显式事务中执行锁定读，并在完成后提交或回滚。自动提交下语句结束就释放事务锁。

查看行锁：

<<< @/db/codes/mysql/lock_raw_select.sql

#### 10.3.2 间隙锁/临键锁

锁范围需要结合隔离级别与访问路径判断，常见情况如下：

1. RR 下完整唯一键的等值锁定读命中已有记录时，通常只锁记录。若记录不存在，则可锁住其应插入的间隙。
2. 非唯一键等值查询或范围查询，通常需要 next-key lock 覆盖搜索范围。具体边界取决于谓词、索引和执行计划，不应把某次观察到的锁范围当作所有版本的规则。
3. RC 下通常不为普通搜索保留间隙锁，但外键检查和重复键检查等仍可能使用间隙锁。

参见 [InnoDB 锁类型](https://dev.mysql.com/doc/refman/8.4/en/innodb-locking.html)。

间隙锁用于阻止插入，多个事务的间隙锁可以共存，获得同一间隙的锁本身不会互相阻塞。

## 11. InnoDB 引擎

本节进一步讨论 InnoDB 的存储结构、内存组件与事务实现。

### 11.1 逻辑存储结构

InnoDB 把数据按下面 5 个层级组织，自顶向下：

| 层级                 | 大小        | 关系                                                 |
| -------------------- | ----------- | ---------------------------------------------------- |
| Tablespace（表空间） | 可增长      | 可独立于单表，也可供多表共享，部分表空间包含多个文件 |
| Segment（段）        | 可增长      | 每个索引分别管理叶子与非叶子节点的段                 |
| Extent（区）         | 默认 1 MiB  | 默认页大小下包含 64 个连续页                         |
| Page（页）           | 默认 16 KiB | 页大小在初始化时配置，区大小随页大小变化             |
| Row（行）            | 变长        | 索引页保存记录，较长字段可有页外部分                 |

每个索引的叶子与非叶子部分使用各自的段。小段可以先分配零散页，增长后再分配完整区，并非每次固定申请四五个区。回滚段管理 undo log，是与索引段不同的结构。

::: tip 行的隐藏列

- `DB_TRX_ID`：聚簇记录上最近一次插入或修改的事务 ID，占 6 字节
- `DB_ROLL_PTR`：回滚指针，用于通过 undo 重建历史版本，占 7 字节
- `DB_ROW_ID`：表没主键也没非空唯一索引时由 InnoDB 自动生成的 6 字节隐藏主键，作为聚簇索引使用。

:::

### 11.2 架构总览

InnoDB 在内存中缓存数据页和日志，在磁盘中管理表空间与日志文件，通过后台线程完成刷页、日志写入和历史版本清理。

#### 11.2.1 内存

**Buffer Pool**：缓存数据页和索引页，页大小与实例配置一致。使用带有老区和新区的 LRU 管理方法减少顺序扫描对热点缓存的冲击。可以按状态区分页：

- **free page**：可供分配的缓冲页
- **clean page**：与磁盘内容一致的缓存页
- **dirty page**：内容已修改且尚未写回磁盘的缓存页

这些状态不等于三个独立缓存池，空闲页、LRU 和脏页刷盘分别有对应的管理链表。

**Change Buffer**：启用后，可缓冲不在 Buffer Pool 中的部分非唯一二级索引页变更，在读取页面或后台合并时应用，以减少随机读取。不是所有索引都支持，例如包含降序列的二级索引不支持。**8.4 的 `innodb_change_buffering` 默认是 `none`，8.0 默认是 `all`**，不能按旧默认值假定它正在工作。

**Adaptive Hash Index**：启用后，InnoDB 根据访问模式为部分索引页建立内存哈希结构，尝试加速适用的查找。它也可能带来锁竞争，需要测量是否受益。**8.4 默认关闭，8.0 默认开启**，开关是 `innodb_adaptive_hash_index`。参见 [自适应哈希索引](https://dev.mysql.com/doc/refman/8.4/en/innodb-adaptive-hash.html)。

**Log Buffer**：redo log 的内存缓冲，由 `innodb_log_buffer_size` 控制，**8.4 默认 64 MiB，8.0 默认 16 MiB**。日志写入与事务提交不必一一对应，长事务也可能在提交前写出部分日志。undo log 不使用这个缓冲区。

#### 11.2.2 磁盘

| 文件 / 区域                               | 用途                                                                                                                   |
| ----------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| System Tablespace（通常为 `ibdata1`）     | 包含 Change Buffer 等内部结构，也可配置存放表数据                                                                      |
| Data Dictionary Tablespace（`mysql.ibd`） | 8.0 起的事务数据字典，以及 `mysql` 系统库中的相关表                                                                    |
| File-Per-Table（`tb_name.ibd`）           | 每张表一个数据文件，包含本表的聚簇索引、二级索引、数据                                                                 |
| General Tablespaces                       | 用户用 `CREATE TABLESPACE` 显式建出来共享给多张表                                                                      |
| Undo Tablespaces                          | 存放持久表事务的 undo log，初始化时创建两个默认表空间                                                                  |
| Temporary Tablespaces                     | 会话临时表、全局临时表的数据                                                                                           |
| Doublewrite Buffer                        | 默认启用，先保存待刷数据页的完整副本，用于恢复部分写入的数据页。8.0.20 起使用独立 doublewrite 文件，旧版存于系统表空间 |
| Redo Log（`#innodb_redo/` 目录）          | 8.0.30 起使用该目录，记录用于崩溃恢复的页修改，容量由 `innodb_redo_log_capacity` 控制                                  |

参见 [undo 表空间](https://dev.mysql.com/doc/refman/8.4/en/innodb-undo-tablespaces.html)、[Doublewrite](https://dev.mysql.com/doc/refman/8.4/en/innodb-doublewrite-buffer.html)和[redo 文件](https://dev.mysql.com/doc/refman/8.4/en/innodb-redo-log.html)。

通用表空间的语法可看下面这条：

<<< @/db/codes/mysql/create_tablespace.sql

#### 11.2.3 后台线程

InnoDB 是多线程引擎，关键的几个角色：

- **Master Thread**：协调后台活动，包括适用的 Change Buffer 合并等任务
- **IO Thread**：处理后台读写请求及完成通知，不能理解为所有磁盘 IO 都变成异步。8.4 的读 IO 线程默认按可用逻辑处理器数计算，受变量上限约束，不能再固定写为 4 个

| 配置项                    | 8.4 默认值                              | 用途           |
| ------------------------- | --------------------------------------- | -------------- |
| `innodb_read_io_threads`  | 可用逻辑处理器数的一半，至少 4，最多 64 | 后台读 IO 线程 |
| `innodb_write_io_threads` | 4                                       | 后台写 IO 线程 |
| `innodb_page_cleaners`    | 随 Buffer Pool 实例数确定               | 脏页刷盘线程   |

日志写入和持久化由专门的日志线程处理。实际配置可查 `SHOW VARIABLES`，运行状态可通过以下命令观察，输出随版本与负载变化：

<<< @/db/codes/mysql/show_engine.sql

- **Purge Thread**：清理不再需要的历史版本与删除标记，线程数由 `innodb_purge_threads` 控制
- **Page Cleaner Thread**：负责 Buffer Pool 的脏页刷盘

参见 [8.4 默认值变化](https://dev.mysql.com/doc/refman/8.4/en/mysql-nutshell.html)和[InnoDB 参数](https://dev.mysql.com/doc/refman/8.4/en/innodb-parameters.html)。

### 11.3 事务原理

InnoDB 的事务靠两套日志撑起来：redo log 负责崩溃恢复（D），undo log 负责回滚和 MVCC（A 与 I）。

#### 11.3.1 redo log

redo 描述页及记录的修改，属于物理层面的恢复日志，不是原始 SQL。WAL 要求相关 redo 先持久化，修改后的数据页才能持久化。简化流程如下：

1. 事务里的修改先改 Buffer Pool 中的页（这页变脏）。
2. 修改的物理动作写进 Log Buffer。
3. 事务 commit 时按 `innodb_flush_log_at_trx_commit` 的策略把 Log Buffer 刷到磁盘 redo log。
4. 脏页可以晚一点由 Page Cleaner 刷回 `.ibd`。

崩溃恢复会根据 redo 恢复尚未落入数据文件的修改，其中可能包括未提交事务的修改，再用 undo 回滚未完成事务。不能说 redo 只重放已提交事务。

`innodb_flush_log_at_trx_commit` 三档：

- `1`（默认）：每次提交要求 redo 写入并完成持久化，可通过组提交合并实际刷盘
- `2`：每次提交写入文件系统缓存，由后台周期性持久化，操作系统崩溃或断电可丢失近期提交
- `0`：提交时不要求写出，由后台周期性写入并持久化，服务进程崩溃也可能丢失近期提交

后台间隔受 `innodb_flush_log_at_timeout` 控制，调度延迟等使“一秒”不是严格的最大丢失窗口。需要 redo 与 binlog 都持久化时，通常同时使用 `innodb_flush_log_at_trx_commit=1` 和 `sync_binlog=1`，前提是操作系统与设备正确兑现刷盘要求。

#### 11.3.2 undo log

undo 保存撤销聚簇记录修改和重建旧版本所需的信息，通常归为逻辑日志。可以用反向操作理解回滚效果，但内部不是逐条保存对应的 `DELETE`、`INSERT` 或 `UPDATE` SQL。

它有两个用途：

- 事务回滚时按 undo log 逆推回去。
- MVCC 读历史版本时沿着回滚指针往老版本走。

undo log 由回滚段中的 undo slot 管理，槽数与页大小有关，不能简单认为一个槽永远等于一个事务。INSERT undo 在提交后不再需要，UPDATE undo 还用于一致性读，必须等历史版本不再被需要时才能 purge。参见 [undo log](https://dev.mysql.com/doc/refman/8.4/en/innodb-undo-logs.html) 和[多版本实现](https://dev.mysql.com/doc/refman/8.4/en/innodb-multi-versioning.html)。

### 11.4 MVCC

MVCC 让一致性读通过历史版本获得可见数据，通常不必与写操作争用记录锁。锁定读、写操作和元数据访问仍受相应锁约束，不能把它理解为所有读写都不阻塞。

#### 11.4.1 当前读 vs 快照读

| 类型                        | 看的是                                              | 加锁                 | 触发场景                                   |
| --------------------------- | --------------------------------------------------- | -------------------- | ------------------------------------------ |
| 锁定读 / 写操作，常称当前读 | 访问当前记录并遵守冲突锁，必要时等待其他事务        | 加数据锁             | DML、`SELECT ... FOR UPDATE` / `FOR SHARE` |
| 一致性读，常称快照读        | ReadView 可见的数据及本事务的修改，可能需要历史版本 | 不加记录锁，仍有 MDL | RC/RR 下的普通 `SELECT`                    |

不同隔离级别下，快照读的"快照"是什么时候定的：

- **Read committed**：每次一致性读建立新的快照，同一事务可能看到其他事务新提交的修改
- **Repeatable read**：第一次一致性读建立快照，后续一致性读复用。`START TRANSACTION WITH CONSISTENT SNAPSHOT` 可在事务开始时建立快照，它在 RR 下生效
- **Serializable**：关闭自动提交时，普通 `SELECT` 隐式转为 `FOR SHARE`。自动提交下独立的只读查询可以使用非锁定一致性读，并非所有 `SELECT` 都加记录锁

锁定读不会按上述 RR 规则建立用于后续一致性读的快照。参见[一致性读](https://dev.mysql.com/doc/refman/8.4/en/innodb-consistent-read.html)。

#### 11.4.2 隐藏字段（再回顾一次）

| 字段          | 作用                                            |
| ------------- | ----------------------------------------------- |
| `DB_TRX_ID`   | 写入这条记录的事务 ID                           |
| `DB_ROLL_PTR` | 指向 undo log 中的上一个版本，串起版本链        |
| `DB_ROW_ID`   | 没主键且没 NOT NULL 唯一索引时，InnoDB 兜底生成 |

#### 11.4.3 undo log 与版本链

聚簇索引保存当前记录，`DB_ROLL_PTR` 指向可用于重建更早状态的 undo 记录。读取时先判断当前版本，不可见才通过 undo 逐步重建历史版本，并非总从“最近修改前的版本”开始读取。

#### 11.4.4 ReadView 与可见性判断

ReadView 保存快照建立时的活跃读写事务信息和可见性边界。以下使用 MySQL 8.4 源码中的字段名，纯只读事务不一定获得普通读写事务 ID：

| 字段               | 含义                                            |
| ------------------ | ----------------------------------------------- |
| `m_ids`            | 快照建立时活跃的读写事务 ID 集合，不含创建者    |
| `m_up_limit_id`    | 最小活跃事务 ID，没有其他活跃事务时等于下方边界 |
| `m_low_limit_id`   | 建立快照时下一个将分配的读写事务 ID             |
| `m_creator_trx_id` | 创建该 ReadView 的事务 ID                       |

判断一个数据版本上的 `DB_TRX_ID` 是否对当前事务可见，按下面顺序逐条判：

| 情况                               | 结论                               |
| ---------------------------------- | ---------------------------------- |
| `trx_id == m_creator_trx_id`       | 自己的修改，可见                   |
| `trx_id < m_up_limit_id`           | 在快照建立前已经提交，可见         |
| `trx_id >= m_low_limit_id`         | 在快照建立时尚未分配的事务，不可见 |
| 位于两个边界之间，且属于 `m_ids`   | 建立快照时仍活跃，不可见           |
| 位于两个边界之间，且不属于 `m_ids` | 建立快照前已提交，可见             |

不可见时继续重建前一版本。如果不存在可见的记录版本，这一行不进入结果集，不是返回一行 SQL `NULL`。删除标记也按相应版本的可见性处理。源码可参考 [MySQL 8.4.6 的 ReadView](https://github.com/mysql/mysql-server/blob/mysql-8.4.6/storage/innobase/include/read0types.h)。

## 12. MySQL 管理

### 12.1 系统数据库

MySQL 数据库安装完成后，自带了以下四个数据库，具体作用如下：

| 数据库               | 含义                                                                                        |
| -------------------- | ------------------------------------------------------------------------------------------- |
| `mysql`              | 存储 MySQL 服务器正常运行所需要的各种信息（时区、主从、用户、权限等）                       |
| `information_schema` | 提供了访问数据库元数据的各种表和视图，包含数据库、表、字段类型及访问权限等                  |
| `performance_schema` | 为 MySQL 服务器运行时状态提供了一个底层监控功能，主要用于收集数据库服务器性能参数           |
| `sys`                | 包含了一系列方便 DBA 和开发人员利用 `performance_schema` 性能数据库进行性能调优和诊断的视图 |

### 12.2 常用工具

#### 12.2.1 `mysql`

该 `mysql` 不是指 MySQL 服务，而是指 MySQL 的客户端工具。

- 语法：`mysql [options] [database]`
- 选项：
  - `-u, --user=name`：指定用户名
  - `-p, --password[=password]`：省略值时交互式输入密码，短选项携带值时不能在 `-p` 后加空格
  - `-h, --host=name`：指定服务器 IP 或域名
  - `-P, --port=port`：指定连接端口
  - `-e, --execute=name`：执行 SQL 语句并退出

`-e` 选项可以在 MySQL 客户端连接数据库后执行 SQL 语句，执行完成后自动退出，对于一些批处理脚本，这种方式尤其方便。

示例：

<<< @/db/codes/mysql/mysqle.sh

#### 12.2.2 `mysqladmin`

`mysqladmin` 是一个执行管理操作的客户端程序。可以用它来检查服务器的配置和当前状态、创建并删除数据库等。

通过帮助文档查看选项：

<<< @/db/codes/mysql/mysqladmin_help.sh

示例：

<<< @/db/codes/mysql/mysqladmin_drop.sh

<<< @/db/codes/mysql/mysqladmin_v.sh

#### 12.2.3 `mysqlbinlog`

由于服务器生成的二进制日志文件以二进制格式保存，所以如果想要检查这些日志的文本格式，就会使用到 `mysqlbinlog` 日志管理工具。

- 语法：`mysqlbinlog [options] log-files1 log-files2 ...`
- 选项：
  - `-d, --database=name`：按数据库筛选，语句日志与行日志的筛选语义不同。
  - `-o, --offset=#`：跳过前 `n` 个日志事件，不是跳过文本行。
  - `-r, --result-file=name`：将输出的文本格式日志输出到指定文件。
  - `-s, --short-form`：显示简单格式，省略掉一些信息。
  - `--start-datetime=date1 --stop-datetime=date2`：按运行该工具的本地时区解释日期时间边界。
  - `--start-position=pos1 --stop-position=pos2`：指定位置间隔内的所有日志。

<<< @/db/codes/mysql/mysqlbinlog.sh

检查 ROW 事件时使用 `--base64-output=DECODE-ROWS -vv` 可显示带注释的伪 SQL。它用于阅读，不是可直接执行的恢复脚本。实际回放需要工具默认生成的可执行输出，并保留完整事务及所需日志文件。参见 [mysqlbinlog](https://dev.mysql.com/doc/refman/8.4/en/mysqlbinlog.html) 和[行事件显示](https://dev.mysql.com/doc/refman/8.4/en/mysqlbinlog-row-events.html)。

#### 12.2.4 `mysqlshow`

`mysqlshow` 客户端对象查找工具，用来快速查找存在哪些数据库、数据库中的表、表中的列或者索引。

- 语法：`mysqlshow [options] [db_name [table_name [col_name]]]`
- 选项：
  - `--count`：显示数据库及表的统计信息（数据库，表均可以不指定）
  - `-i`：显示指定数据库或者指定表的状态信息
  - `-k, --keys`：显示表索引

示例：

<<< @/db/codes/mysql/mysqlshow.sh

<<< @/db/codes/mysql/mysqlshow_db.sh

<<< @/db/codes/mysql/mysqlshow_db_tb.sh

`--count` 会执行计数，较大的表可能耗时，输出随当前数据变化。末尾参数中的 `_` 等字符会被当作通配符，按表名查看时应像示例一样转义，或添加列模式参数。参见 [mysqlshow](https://dev.mysql.com/doc/refman/8.4/en/mysqlshow.html)。

#### 12.2.5 `mysqldump`

`mysqldump` 是 MySQL 的客户端备份工具，主要用于：

- **数据库备份**：生成包含建表语句、数据插入语句的 SQL 文件
- **数据迁移**：在不同 MySQL 实例（或数据库）间转移数据。

#### 12.2.6 `mysqlimport`/`source`

`mysqlimport` 用 `LOAD DATA` 导入分隔文本文件，并根据文件名确定目标表名，不限于 `mysqldump -T` 生成的文件。例如 `tb_insert.csv` 对应已存在的 `tb_insert` 表，列分隔符需与文件一致。

`source /path/to/backup.sql` 是 `mysql` 客户端执行 SQL 文件的命令，也可在终端使用 `mysql -u root -p demo_query < backup.sql`。它与导入 CSV 的 `mysqlimport` 用途不同。参见 [mysqlimport](https://dev.mysql.com/doc/refman/8.4/en/mysqlimport.html)。

## 13. 日志

### 13.1 错误日志

错误日志是 MySQL 中最重要的日志之一，它记录了 mysqld 启动和停止时，以及服务器在运行过程中发生任何严重错误时的相关信息。当数据库出现任何故障无法正常使用时，建议首先查看此日志。

错误输出的目标取决于启动方式和 `log_error` 配置，可以是文件或标准错误，不存在跨平台统一的默认路径。查看配置：

<<< @/db/codes/mysql/show_log_err.sql

### 13.2 二进制日志

二进制日志（binlog）记录需要用于复制和恢复的数据变更及相关事件，通常不记录不改变数据的普通 `SELECT`、`SHOW`。不能将其理解为所有 SQL 的完整审计记录，是否记录还受 `sql_log_bin`、过滤配置及语句类型等影响。

作用：

1. 配合完整备份进行时间点恢复
2. MySQL 的主从复制。在 MySQL 8 版本中，默认二进制日志是开启的，涉及到的参数如下：

<<< @/db/codes/mysql/show_log_bin.sql

MySQL 服务器中提供了多种格式来记录二进制日志，具体格式及特点如下：

| 日志格式  | 含义                                                                                              |
| :-------: | ------------------------------------------------------------------------------------------------- |
| STATEMENT | 基于 SQL 语句的日志记录，记录的是 SQL 语句，对数据进行修改的 SQL 都会记录在日志文件中。           |
|    ROW    | 基于行的日志记录，记录的是每一行的数据变更。（默认）                                              |
|   MIXED   | 混合了 STATEMENT 和 ROW 两种格式，优先采用 STATEMENT，在某些特殊情况下会自动切换为 ROW 进行记录。 |

查看当前二进制日志格式的语句：

<<< @/db/codes/mysql/show_log_bin_fmt.sql

8.4 默认开启 binlog，并使用 `ROW`。`binlog_format` 从 8.0.34 起已弃用，虽然 8.4 仍支持上述三种格式，新的配置应以行日志为主。DDL 即使在 ROW 模式下仍主要记录为语句事件。参见[二进制日志](https://dev.mysql.com/doc/refman/8.4/en/binary-log.html)。

#### 13.2.1 日志查看

由于日志是以二进制方式存储的，不能直接读取，需要通过二进制日志查询工具 `mysqlbinlog` 来查看。

- 语法：`mysqlbinlog [参数选项] logfilename`
- 参数选项：
  - `-d`：按数据库过滤，具体语义见 12.2.3。
  - `-o`：跳过指定数量的事件。
  - `-v`：把行事件显示为注释形式的伪 SQL。
  - `-vv`：在上述显示中增加类型等元信息。

#### 13.2.2 日志删除

对于比较繁忙的业务系统，每天生成的 binlog 数据巨大，如果长时间不清除，将会占用大量磁盘空间。可以通过以下几种方式清理日志：

| 指令                                             | 含义                               |
| ------------------------------------------------ | ---------------------------------- |
| `PURGE BINARY LOGS TO 'binlog.000123'`           | 删除指定文件之前的日志，保留该文件 |
| `PURGE BINARY LOGS BEFORE '2026-01-01 00:00:00'` | 删除指定时间之前的日志文件         |

清理前应确认备份恢复和所有复制节点均不再需要这些日志。8.4 已移除 `RESET MASTER` 和 `PURGE MASTER LOGS`。`RESET BINARY LOGS AND GTIDS` 会删除全部 binlog 并清空 GTID 执行历史，不应作为日常日志清理命令。

也可以在 MySQL 的配置文件中配置二进制日志的过期时间，设置了以后，二进制日志过期会自动删除。

<<< @/db/codes/mysql/show_log_bin_expire.sql

8.4 默认 `binlog_expire_logs_seconds=2592000`，即 30 天，`binlog_expire_logs_auto_purge=ON`。自动清理通常在启动、日志轮转等时机检查，不是文件达到期限的瞬间删除。

### 13.3 查询日志

通用查询日志记录客户端连接、断开及服务器接收到的语句，不代表语句已经执行成功。8.4 默认关闭，可用于临时观察请求，与记录数据变更的 binlog 用途不同。

<<< @/db/codes/mysql/show_log_general.sql

修改 MySQL 的配置文件，添加以下内容：

<<< @/db/codes/mysql/log_general.cnf

日志目标还由 `log_output` 决定，可能为文件或表。持续记录会增加开销，应在诊断完成后恢复所需配置。参见[通用查询日志](https://dev.mysql.com/doc/refman/8.4/en/query-log.html)。

### 13.4 慢查询日志

见 [7.4.2 慢查询日志](#_7-4-2-慢查询日志)

## 14. 主从复制

传统复制由源服务器（source，主库）提供 binlog，副本（replica，从库）接收并应用事件。普通复制默认异步，读副本可能读到旧数据，主库提交成功不表示副本已经应用完成。一主多副本及链式复制均可配置。

常见用途：

1. **故障切换**：选择合适的副本提升为主库，普通复制本身不提供完整的自动选主和客户端切换
2. **读写分离**：将适合容忍复制延迟的读取分配给副本
3. **备份**：在副本执行备份可以减少主库负担，但仍会消耗副本资源，加锁时可能阻塞复制应用

简化工作流程：

1. 源服务器记录 binlog，通过连接对应的 binlog dump 线程向副本发送事件
2. 副本接收线程（receiver，传统称 IO Thread）取得事件并写入 relay log
3. 副本应用线程（applier，传统称 SQL Thread）读取 relay log，或由协调线程分派给多个 worker 并行应用

8.4 的 binlog 默认是 `ROW`，但 **`gtid_mode` 默认是 `OFF`**。GTID 复制需要另行配置 `gtid_mode=ON`、`enforce_gtid_consistency=ON` 等前提，启用后可通过 `SOURCE_AUTO_POSITION=1` 使用事务标识定位，不能把它当作默认开启的功能。

8.4 使用 `CHANGE REPLICATION SOURCE TO`、`START REPLICA`、`STOP REPLICA`、`SHOW REPLICA STATUS`。旧的 `CHANGE MASTER TO`、`START SLAVE`、`SHOW SLAVE STATUS` 等语法已经移除。

在副本查看 `SHOW REPLICA STATUS\G`，重点检查 `Replica_IO_Running`、`Replica_SQL_Running`、`Last_IO_Error`、`Last_SQL_Error`。`Seconds_Behind_Source` 只反映特定条件下的延迟估计，网络延迟、接收线程落后或断开等可能使它低估延迟或为 `NULL`，不能只凭该值为零判断主副本完全一致。

## 15. 分库分表

垂直扩容是提升单节点的 CPU、内存、存储等资源，水平扩容是增加节点分担数据或负载。分库分表是拆分数据的手段，可用于缓解单实例或单表的容量及访问压力，需先确认实际瓶颈：

- **IO 瓶颈**：工作集不能有效缓存，随机读写或网络传输成为限制
- **CPU 瓶颈**：查询处理、排序、分组等持续消耗处理能力

拆分不保证负载均匀，路由键不合适时热点仍可能集中。还需处理跨片查询、事务、全局约束、ID 生成、扩容迁移等问题，应先评估索引、查询及现有部署能否满足需求。

### 15.1 拆分策略

#### 15.1.1 垂直分库

按业务域将不同的表放在不同库中，便于独立部署。跨库关联与事务需要另外设计。

#### 15.1.2 垂直分表

按列拆出多个表，通常保留共同的主键。例如将低频的大字段从主表拆出，需要完整记录时通过键进行 `JOIN`。

#### 15.1.3 水平分库

按路由键把同类记录分配到不同库，相关分片表通常保持相同结构。跨库查询可能需要应用或中间件合并结果。

#### 15.1.4 水平分表

按路由键把同一逻辑表的记录分配到多个同结构的物理表。各片行集合共同组成逻辑表，唯一约束和外键只在物理表或实例范围内生效，不能自动保证跨片约束。

### 15.2 实现技术

- **Apache ShardingSphere-JDBC**：在 Java 应用的 JDBC 接口层处理 SQL 解析、路由、改写和结果归并。参见 [5.5.2 文档](https://shardingsphere.apache.org/document/5.5.2/en/overview/)。
- **MyCat**：独立数据库代理，客户端通过代理访问后端分片。是否需要调整业务取决于路由和 SQL 支持范围，性能需结合实际部署测量。这里仅比较架构，不涉及与 MySQL 8.4 的接入兼容性。参见 [MyCat 1.6.7.5](https://github.com/MyCATApache/Mycat-Server)。
- **GORM Sharding**：运行在 Go 应用内的 GORM 分表插件，通过修改访问目标表实现路由，本身不提供跨多个数据库节点的完整分库方案。

### 15.3 GORM Sharding

以下以 **`gorm.io/sharding v0.6.2`** 为例。插件在连接访问层解析 SQL 并改写表名，支持 PostgreSQL、MySQL 和主键生成器。它仍有 SQL 和配置限制，不能把普通 GORM 的所有操作都视为自动支持。参见 [v0.6.2 文档](https://github.com/go-gorm/sharding/blob/v0.6.2/README.md)。

特性：

1. 需要先准备物理分片表，再注册逻辑表和路由规则
2. 常规路由需要分片键的等值条件，仅包含范围条件或遗漏分片键时不能假定自动扫描所有分片
3. 默认 Snowflake 可从生成的主键解析分片，其他主键方案如需按主键查询要配置相应路由算法
4. v0.6.2 不支持 GORM 的 `PrepareStmt: true`

#### 15.3.1 安装

<<< @/db/codes/mysql/sharding_install.sh

#### 15.3.2 用法

先用管理账号创建空练习库和专用账号：

<<< @/db/codes/mysql/setup_sharding.sql

在独立的 Go module 中安装上述依赖，将代码保存为 `main.go`。设置 `MYSQL_DSN` 后执行 `go run .`，例如：

```sh
MYSQL_DSN='demo_sharding:DemoPass_846!@tcp(127.0.0.1:3306)/demo_sharding?charset=utf8mb4&parseTime=true' go run .
```

这是单进程练习，物理表通过迁移 API 创建，示例账号具有该库建表和读写权限。

<<< @/db/codes/mysql/sharding_config.go

::: warning
内置 Snowflake 的节点号由分片号决定，多个应用进程使用相同规则时可能发生 ID 冲突。生产环境应采用保证跨进程唯一性的主键方案，更换生成规则时，按主键路由的算法也需与之匹配。分片数量变化需要数据迁移，修改配置不会自动移动已有记录。
:::

### 15.4 分片方式

1. 范围分片：根据指定的字段所在的范围与节点的对应关系，来决定该数据属于哪一个分片。
2. 取模分片：根据指定的字段值与节点数量进行求模运算，根据运算结果，来决定该数据属于哪一个分片。
3. 一致性哈希：节点集合不变时，同一键映射保持稳定。增加或移除节点时仍有部分键重新映射，优势是通常减少迁移范围，已有数据也需要迁移。参见 [Karger et al., 1997](https://doi.org/10.1145/258533.258660)。
4. 枚举分片：将有限取值映射到指定节点，需评估各取值的数据量和访问量，避免热点集中。
5. 按日期分片：按时间范围拆分，便于归档与过期数据管理，但最新分片可能承受主要写入压力。

这些是可选的路由思路，不代表前述插件都自带相应实现。GORM Sharding 的自定义算法需同时维护路由规则、表后缀和迁移方案。
