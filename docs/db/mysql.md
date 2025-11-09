---
outline: [2, 3]
---

# MySQL

<<< @/db/codes/mysql/conn.sh

## 1 SQL

| 分类 | 全称                       | 说明                                                   |
| ---- | -------------------------- | ------------------------------------------------------ |
| DDL  | Data Definition Language   | 数据定义语言，用来定义数据库对象（数据库、表、字段）   |
| DML  | Data Manipulation Language | 数据操作语言，用来对数据库表中的数据进行增删改         |
| DQL  | Data Query Language        | 数据查询语言，用来查询数据库中表的记录                 |
| DCL  | Data Control Language      | 数据控制语言，用来创建数据库用户、控制数据库的访问权限 |

### 1.1 DDL

#### 1.1.1 数据库操作

- 查询

<<< @/db/codes/mysql/ddl_query.sql

- 创建

<<< @/db/codes/mysql/ddl_create.sql

- 删除

<<< @/db/codes/mysql/ddl_drop.sql

- 使用

<<< @/db/codes/mysql/ddl_use.sql

#### 1.1.2 表操作

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

- 删除表

<<< @/db/codes/mysql/ddl_drop_tbl.sql

- 删除指定表，并重新创建该表

<<< @/db/codes/mysql/ddl_truncate_tbl.sql

#### 1.1.3 数据类型

数值类型：

| 类型               | 大小    |
| ------------------ | ------- |
| `TINYINT`          | 1 byte  |
| `SMALLINT`         | 2 bytes |
| `MEDIUMINT`        | 3 bytes |
| `INT` 或 `INTEGER` | 4 byte  |
| `BIGINT`           | 8 byte  |
| `FLOAT`            | 4 byte  |
| `DOUBLE`           | 8 byte  |
| `DECIMAL`          |         |

字符串类型：`CHAR`、`VARCHAR`、`TINYBLOB`、`TINYTEXT`、`BLOB`、`TEXT`、`MEDIUMBLOB`、`MEDIUMTEXT`、`LONGBLOB`、`LONGTEXT`

日期类型：

| 类型        | 大小    | 格式                  |
| ----------- | ------- | --------------------- |
| `DATE`      | 3 bytes | `YYYY-MM-DD`          |
| `TIME`      | 3 bytes | `HH:MM:SS`            |
| `YEAR`      | 1 byte  | `YYYY`                |
| `DATETIME`  | 8 bytes | `YYYY-MM-DD HH:MM:SS` |
| `TIMESTAMP` | 4 bytes | `YYYY-MM-DD HH:MM:SS` |

### 1.2 DML

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

| 逻辑运算符     | 功能                     |
| -------------- | ------------------------ |
| `AND` 或 `&&`  | 并且（多个条件同时成立） |
| `OR` 或 `\|\|` | 或者（任一条件成立）     |
| `NOT` 或 `!`   | 非，不是                 |

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

`NULL` 值不参与聚合运算。

#### 1.3.4 分组查询

<<< @/db/codes/mysql/dql_group.sql

#### 1.3.5 排序查询

<<< @/db/codes/mysql/dql_order.sql

#### 1.3.6 分页查询

<<< @/db/codes/mysql/dql_limit.sql

### 1.4 DCL

#### 1.4.1 用户管理

- 查询用户

<<< @/db/codes/mysql/dcl_select_user.sql

- 创建用户

<<< @/db/codes/mysql/dcl_create_user.sql

- 修改用户密码

<<< @/db/codes/mysql/dcl_change_pwd.sql

- 删除用户

<<< @/db/codes/mysql/dcl_drop_user.sql

#### 1.4.2 权限控制

| 常用权限                | 说明               |
| ----------------------- | ------------------ |
| `ALL`，`ALL PRIVILEGES` | 所有权限           |
| `SELECT`                | 查询数据           |
| `INSERT`                | 插入数据           |
| `UPDATE`                | 修改数据           |
| `DELETE`                | 删除数据           |
| `ALTER`                 | 修改表             |
| `DROP`                  | 删除数据库/表/视图 |
| `CREATE`                | 创建数据库/表      |

- 查询权限

<<< @/db/codes/mysql/dcl_show_grants.sql

- 授予权限

<<< @/db/codes/mysql/dcl_grant.sql

- 撤销权限

<<< @/db/codes/mysql/dcl_revoke.sql

## 2 函数

### 2.1 字符串函数

| 常用函数                     | 说明                                                          |
| ---------------------------- | ------------------------------------------------------------- |
| `CONCAT(S1, S2, ..., Sn)`    | 字符串拼接                                                    |
| `LOWER(str)`                 | 将字符串全部转为小写                                          |
| `UPPER(str)`                 | 将字符串全部转为大写                                          |
| `LPAD(str, n, pad)`          | 左填充，用 `pad` 对 `str` 左边进行填充，达到 `n` 个字符串长度 |
| `RPAD(str, n, pad)`          | 右填充，用 `pad` 对 `str` 右边进行填充，达到 `n` 个字符串长度 |
| `TRIM(str)`                  | 去掉字符串头部和尾部的空格                                    |
| `SUBSTRING(str, start, len)` | 返回字符串 `str` 从 `start` 位置起的 `len` 个长度的子字符串   |

<<< @/db/codes/mysql/func_string.sql

### 2.2 数值函数

| 常用函数      | 说明                                   |
| ------------- | -------------------------------------- |
| `CEIL(x)`     | 向上取整                               |
| `FLOOR(x)`    | 向下取整                               |
| `MOD(x, y)`   | 返回 `x % y`                           |
| `RAND()`      | 返回 `0~1` 的随机数                    |
| `ROUND(x, y)` | 返回 `x` 的四舍五入值，保留 `y` 位小数 |

<<< @/db/codes/mysql/func_math.sql

### 2.3 日期函数

| 常用函数                             | 说明 |
| ------------------------------------ | ---- |
| `CURDATE()`                          |      |
| `CURTIME()`                          |      |
| `NOW()`                              |      |
| `YEAR(date)`                         |      |
| `MONTH(date)`                        |      |
| `DAY(date)`                          |      |
| `DATE_ADD(date, INTERVAL expr type)` |      |
| `DATEDIFF(date1, date2)`             |      |

<<< @/db/codes/mysql/func_date.sql

### 2.4 流程函数

| 常用函数                                                     | 说明                                                                   |
| ------------------------------------------------------------ | ---------------------------------------------------------------------- |
| `IF(value, t, f)`                                            | 如果 `value` 为 `true`，则返回 `t`，否则返回 `f`                       |
| `IFNULL(value1, value2)`                                     | 如果 `value1` 不为 `NULL`，返回 `value1`，否则返回 `value2`            |
| `CASE WHEN [val1] THEN [res1] ... ELSE [default] END`        | 如果 `val1` 为 `true`，返回 `res1`，...，否则返回 `default` 默认值     |
| `CASE [expr] WHEN [val1] THEN [res1] ... ELSE [default] END` | 如果 `expr` 的值为 `val1`，返回 `res1`，...，否则返回 `default` 默认值 |

<<< @/db/codes/mysql/func_flow.sql

## 3 约束

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

外键用来让两张表的数据之间建立连接，从而保证数据的的一致性和完整性。

- 添加外键

方式一：

<<< @/db/codes/mysql/foreign_key_1.sql

方式二：

<<< @/db/codes/mysql/foreign_key_2.sql

外键的删除/更新行为：

| 行为          | 说明                                                                                                                       |
| ------------- | -------------------------------------------------------------------------------------------------------------------------- |
| `NO ACTION`   | 当在父表中删除/更新对应记录时，首先检查该记录是否有对应外键，如果有则不允许删除/更新。（与 `RESTRICT` 一致）               |
| `RESTRICT`    | 当在父表中删除/更新对应记录时，首先检查该记录是否有对应外键，如果有则不允许删除/更新。（与 `NO ACTION` 一致）              |
| `CASCADE`     | 当在父表中删除/更新对应记录时，首先检查该记录是否有对应外键，如果有则也删除/更新外键在子表中的记录。                       |
| `SET NULL`    | 当在父表中删除对应记录时，首先检查该记录是否有对应外键，如果有则设置子表中该外键值为 `NULL`。（这个字段需要允许为 `NULL`） |
| `SET DEFAULT` | 父表有变更时，子表将外键列设置成一个默认的值。（Innodb 不支持）                                                            |

<<< @/db/codes/mysql/foreign_key_action.sql

## 4 多表查询

### 4.1 内连接

- 隐式内连接

<<< @/db/codes/mysql/join_inner_where.sql

- 显示内连接

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

### 4.5 子查询

SQL 语句中嵌套 `SELECT` 语句，成为嵌套查询，又称子查询。

#### 4.5.1 标量子查询

子查询返回的结果是单个值（数字、字符串、日期等），这种子查询成为标量子查询。

<<< @/db/codes/mysql/subquery_scalar.sql

#### 4.5.2 列子查询

子查询返回的结果是一列，这种子查询成为列子查询。

<<< @/db/codes/mysql/subquery_column.sql

#### 4.5.3 行子查询

子查询返回的结果是一行，这种子查询成为行子查询。

<<< @/db/codes/mysql/subquery_row.sql

#### 4.5.4 表子查询

子查询返回的结果是多行多列，这种子查询成为表子查询。

<<< @/db/codes/mysql/subquery_table.sql

## 5 事务

事务是一组操作的集合，它是一个不可分割的工作单位，事务会把所有的操作作为一个整体一起向操作系统提交或撤销操作请求，即这些操作要么同时成功，要么同时失败。

默认 MySQL 的事务是自动提交的，当执行一条 DML 语句，MySQL 会立即隐式地提交事务。

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

- 原子性（<span style="color:red;">A</span>tomicity）：事务是不可分割的最小操作单元，要么全部成功，要么全部失败。
- 一致性（<span style="color:red;">C</span>onsistency）：事务完成时，必须使所有的数据都保持一致状态。
- 隔离性（<span style="color:red;">I</span>solation）：数据库系统提供的隔离机制，保证事务在不受外部并发操作影响的独立环境下进行。
- 持久性（<span style="color:red;">D</span>urability）：事务一旦提交或回滚，它对数据库中的数据的改变就是永久的。

### 5.3 并发事务问题

| 问题       | 描述                                                         |
| ---------- | ------------------------------------------------------------ |
| 脏读       | 一个事务读取了另一个事务尚未提交的修改数据                   |
| 不可重复读 | 同一事务内，多次读取同一数据，但结果不一致                   |
| 幻读       | 同一事务内，多次执行相同的查询条件，但返回的结果集行数不一致 |

### 5.4 事务隔离级别

| 隔离级别                | 脏读 | 不可重复读 | 幻读 |
| ----------------------- | :--: | :--------: | :--: |
| Read uncommitted        |  √   |     √      |  √   |
| Read committed          |  ×   |     √      |  √   |
| Repeatable Read（默认） |  ×   |     ×      |  √   |
| Serializable            |  ×   |     ×      |  ×   |

- 查看事务隔离级别

<<< @/db/codes/mysql/transaction_isolation_select.sql

- 设置事务隔离级别

<<< @/db/codes/mysql/transaction_isolation_set.sql

`SESSION` 表示当前会话，`GLOBAL` 表示全局。

## 6 存储引擎

- 查询当前数据库支持的存储引擎

<<< @/db/codes/mysql/engine_show.sql

|       Engine       | Support | Comment                                                        | Transactions |  XA  | Savepoints |
| :----------------: | :-----: | -------------------------------------------------------------- | :----------: | :--: | :--------: |
|      ARCHIVE       |   YES   | Archive storage engine                                         |      NO      |  NO  |     NO     |
|     BLACKHOLE      |   YES   | /dev/null storage engine (anything you write to it disappears) |      NO      |  NO  |     NO     |
|     MRG_MYISAM     |   YES   | Collection of identical MyISAM tables                          |      NO      |  NO  |     NO     |
|     FEDERATED      |   NO    | Federated MySQL storage engine                                 |     NULL     | NULL |    NULL    |
|       MyISAM       |   YES   | MyISAM storage engine                                          |      NO      |  NO  |     NO     |
| PERFORMANCE_SCHEMA |   YES   | Performance Schema                                             |      NO      |  NO  |     NO     |
|       InnoDB       | DEFAULT | Supports transactions, row-level locking, and foreign keys     |     YES      | YES  |    YES     |
|       MEMORY       |   YES   | Hash based, stored in memory, useful for temporary tables      |      NO      |  NO  |     NO     |
|        CSV         |   YES   | CSV storage engine                                             |      NO      |  NO  |     NO     |

- 在创建表时指定存储引擎

<<< @/db/codes/mysql/engine_create_table.sql

### 6.1 InnoDB

InnoDB 是一种兼顾高可靠性和高性能的通用存储引擎，在 MySQL 5.5 之后，InnoDB 是 MySQL 默认的存储引擎。

特点：

- DML 操作遵循 ACID 模型，支持事务
- 行级锁，提高并发访问性能
- 支持外键 `FOREIGN KEY` 约束，保证数据的完整性和一致性

文件：每个 InnoDB 引擎的表都会对应一个表空间文件 `tb_name.ibd`，存储该表的表结构（frm、sdi）、数据和索引。

`.ibd` 文件存在的前提是启用了独立表空间（默认启用），可通过以下命令确认：

<<< @/db/codes/mysql/engine_var_idb.sql

Ubuntu 的 `.ibd` 文件默认放在 `/var/lib/mysql/数据库名` 路径下。

例如，我的 `test` 数据库下有两张使用 InnoDB 引擎的表 `tb_user`、`tb_account`。

那么在 `/var/lib/mysql/test` 路径下，就会存在 `tb_user.ibd`、`tb_account.ibd` 两个文件。

`.ibd` 文件是二进制文件，无法直接打开，需要使用命令 `ibd2sdi 表名.ibd` 才能查看。

逻辑存储结构：

1. Tablespace：表空间，一个表空间包含多个段
2. Segment：段，一个段包含多个区
3. Extend：区，一个区大小为 1 M，包含多个页（64 个页）
4. Page：页，一个页大小为 16 K，包含多个行
5. Row：行

### 6.2 MyISAM

MyISAM 是 MySQL 早期的默认存储引擎。

特点：

- 不支持事务，不支持外键
- 支持表锁，不支持行锁
- 访问速度快

文件：

- `tb_name.sdi`：存储表结构信息
- `tb_name.MYD`：存储数据
- `tb_name.MYI`：存储索引

### 6.3 MEMORY

Memory 引擎的表数据存储在内存中，由于收到硬件问题或断电问题的影响，只能将这些表作为临时表或缓存使用。

特点：

- 内存存放
- Hash 索引（默认）

文件：`tb_name.sdi` 存储表结构信息。

### 6.4 对比

|   特点    | InnoDB | MyISAM | MEMORY |
| :-------: | :----: | :----: | :----: |
|   事务    |  支持  |   -    |   -    |
|   外键    |  支持  |   -    |   -    |
|  锁机制   |  行锁  |  表锁  |  表锁  |
| B+ 树索引 |  支持  |  支持  |  支持  |
| Hash 索引 |   -    |   -    |  支持  |

## 7 索引

### 7.1 结构

MySQL 采用 B+ 树作为索引。

| 特性             | B 树                                | B+ 树                                         |
| ---------------- | ----------------------------------- | --------------------------------------------- |
| **存储内容**     | 非叶子节点和叶子节点均存储键+值     | 仅叶子节点存储键+值，非叶子节点仅存键（索引） |
| **叶子节点**     | 孤立存在，无链接                    | 所有叶子节点通过双向链表连接                  |
| **查询效率**     | 随机查询可能更快（找到键即可返回）  | 随机查询略慢（必须遍历到叶子节点）            |
| **范围查询**     | 需回溯父节点，效率低                | 利用叶子节点链表，一次遍历完成，效率高        |
| **节点存储密度** | 低（键+值占用空间大，单节点键数少） | 高（非叶子节点仅存键，单节点键数更多）        |
| **IO 次数**      | 较多（层级可能更深）                | 较少（层级更浅，因单节点键数多）              |
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
| **二级索引（Secondary Index）** | 将数据与所有分开存储，索引结构的叶子节点关联的是对应的主键 | 可以存在多个         |

聚簇索引选取规则：

- 如果存在主键，主键索引就是聚簇索引
- 如果不存在主键，将使用第一个唯一索引作为聚簇索引
- 如果表没有主键和唯一索引，则 InnoDB 会自动生成一个 `rowid` 作为隐藏的聚簇索引

### 7.3 语法

- 创建索引

<<< @/db/codes/mysql/index_create.sql

- 查看索引

<<< @/db/codes/mysql/index_show.sql

- 删除索引

<<< @/db/codes/mysql/index_drop.sql

### 7.4 SQL 性能分析

#### 7.4.1 SQL 执行频率

MySQL 客户端连接成功后，通过 `SHOW [SESSION | GLOBAL] STATUS` 命令可以提供服务器状态信息。通过如下指令，可以查看当前数据库的 `INSERT`、`UPDATE`、`DELETE`、`SELECT` 的访问频次：

<<< @/db/codes/mysql/status_show.sql

#### 7.4.2 慢查询日志

慢查询日志记录了所有执行时间超过指定参数（`long_query_time`，单位：秒，默认为 `10`）的所有 SQL 语句的日志。

MySQL 的慢查询日志默认没有开启，需要在 MySQL 的配置文件（Ubuntu 默认在 `/etc/mysql/mysql.conf.d/mysqld.cnf`）配置如下信息：

<<< @/db/codes/mysql/slow_query_log.cnf

#### 7.4.3 PROFILE 详情

`SHOW PROFILES` 能够在做 SQL 优化时帮助我们了解时间都耗到哪里去了。通过 `have_profiling` 参数，能够看到当前 MySQL 是否支持 profile 操作：

<<< @/db/codes/mysql/have_profiling.sql

默认 `profiling` 是关闭的，可以通过 `SET` 语句在 `SESSION/GLOBAL` 级别开启 `profiling`：

<<< @/db/codes/mysql/set_profiling.sql

执行一系列的业务 SQL 操作，然后通过如下指令查看执行耗时：

<<< @/db/codes/mysql/profile_show.sql

#### 7.4.4 `EXPLAIN` 执行计划

`EXPLAIN` 或者 `DESC` 命令获取 MySQL 如何执行 `SELECT` 语句的信息，包括在 `SELECT` 语句执行过程中表如何连接和连接的顺序。

语法：直接在 `SELECT` 语句前加关键字 `EXPLAIN` / `DESC`

<<< @/db/codes/mysql/explain.sql

`EXPLAIN` 执行计划各字段含义：

- **`id`**：`SELECT` 查询的序列号，表示查询中执行 `SELECT` 子句或者是操作表的顺序（大 `id` 先执行，同 `id` 按顺序执行）
- **`select_type`**：表示 `SELECT` 的类型，常见取值：
  - `SIMPLE`：简单查询，无表连接或子查询
  - `PRIMARY`：主查询，即外层的查询
  - `UNION`：`UNION` 中的第二个或者后面的查询语句
  - `SUBQUERY`：`SELECT` / `WHERE` 之后包含了子查询
- **`type`**：连接类型，性能从好到差依次为：
  1. `NULL`：不涉及表（如 `SELECT 1+1`）
  2. `system`：表中只有一行数据（`const` 的特例，如系统表）
  3. `const`：通过主键/唯一索引匹配到单行
  4. `eq_ref`：多表连接中，被连接表通过唯一索引匹配
  5. `ref`：通过普通索引匹配多行
  6. `range`：索引范围扫描（如 `BETWEEN`、`IN`、`>` 等）
  7. `index`：扫描整个索引树（不扫描数据行，比 `ALL` 快）
  8. `ALL`：全表扫描（性能最差）
- **`possible_keys`**：显示可能应用在这张表上的索引，一个或多个
- **`key`**：实际使用的索引，如果为 `NULL`，则没有使用索引
- **`key_len`**：表示索引中使用的字节数，该值为索引字段的最大可能长度，并非实际使用长度，在不损失精确性的前提下，长度越短越好
- **`rows`**：MySQL 认为必须要执行查询的行数，在 InnoDB 引擎的表中，是一个估计值，不一定准确
- **`filtered`**：表示返回结果的行数占需读取行数的百分比，该值越大越好

### 7.5 使用规则

#### 7.5.1 最左前缀原则

如果索引了多列（联合索引），查询从索引的最左列开始，并且不跳过索引中的列。如果跳过某一列，后面的字段索引将失效

<<< @/db/codes/mysql/combined_index.sql

#### 7.5.2 范围查询

联合索引中，出现范围查询（`>`、`<`），范围查询右侧的列索引失效

<<< @/db/codes/mysql/combined_index_range.sql

#### 7.5.3 索引失效

以下情况会导致索引失效：

- **索引字段参与运算**

  如 `WHERE id + 1 = 10`（对 `id` 运算）、`WHERE SUBSTR(name, 1, 3) = "abc"`（函数处理）。

- **字符串不加引号**

  如 `WHERE name = 123`，会触发隐式类型转换，导致索引失效，应改为 `"123"`。

- **`LIKE` 以通配符开头**

  如 `WHERE name LIKE "%abc"`（前导模糊匹配），索引失效；而 `WHERE name LIKE "abc%"`（后缀模糊）可命中索引。

- **使用 `OR` 连接非索引字段**

  如 `WHERE idx_col = 1 OR no_idx_col = 2`（`no_idx_col` 无索引），会导致整个查询无法使用索引。

- **数据分布影响**

  如果 MySQL 评估使用索引比全表更慢，则不使用索引

#### 7.5.4 SQL 提示

SQL 提示，是优化数据库的一个重要手段，简单来说，就是在 SQL 语句中加入一些人为的提示来达到优化操作的目的。

**`USE INDEX`**：

<<< @/db/codes/mysql/index_use.sql

**`IGNORE INDEX`**：

<<< @/db/codes/mysql/index_ignore.sql

**`FORCE INDEX`**：

<<< @/db/codes/mysql/index_force.sql

#### 7.5.5 覆盖索引

尽量使用覆盖索引（查询使用了索引，并且需要返回的列在该索引中已经全部能够找到），减少 `SELECT *`。

#### 7.5.6 前缀索引

当字段类型为字符串（`VARCHAR`、`TEXT` 等）时，有时候需要索引很长的字符串，这会让索引变得更大，查询时，浪费大量磁盘 IO，影响查询效率。此时可以只将字符串的一小部分前缀，建立索引，这样大大节约索引空间，从而提高索引效率。

<<< @/db/codes/mysql/index_prefix.sql

可以根据索引的选择性来决定**前缀长度**，而选择性是指不重复的索引数（基数）和数据表的记录总数的比值，索引选择性越高则查询效率越高，唯一索引的选择性是 `1`，这是最好的选择性，性能也是最好的。

<<< @/db/codes/mysql/index_prefix_len.sql

## 8 SQL 优化

### 8.1 `INSERT` 优化

#### 8.1.1 小批量插入数据

<<< @/db/codes/mysql/optimize_insert_old.sql

优化 1：改为批量插入

<<< @/db/codes/mysql/optimize_insert_multi.sql

优化 2：手动提交事务

默认情况下，每句 DML 都是自动提交事务的。大量的插入语句会频繁的开启事务与提交事务。

<<< @/db/codes/mysql/optimize_insert_transaction.sql

优化 3：主键顺序插入

```txt
主键乱序插入：8 1 9 21 88 2 4 15 89 5 7 3
主键顺序插入：1 2 3 4 5 7 8 9 15 21 88 89
```

#### 8.1.2 大批量插入数据

如果一次性需要插入大批量数据，使用 `INSERT` 语句插入性能极低，此时可以使用 MySQL 数据库提供的 `LOAD` 指令进行插入。操作如下：

<<< @/db/codes/mysql/load_infile.sh
<<< @/db/codes/mysql/load_infile.sql

### 8.2 主键优化

在 InnoDB 存储引擎中，表数据都是根据主键顺序组织存放的，这种存储方式的表称为<span style="color:red;">索引组织表</span>（Index Organized Table，<span style="color:red;">IOT</span>）。

#### 8.2.1 页分裂

页可以为空，也可以填充一半，也可以填充 100%。每个页包含了 2~N 行数据（如果一行数据过大，会行溢出），根据主键排列。

#### 8.2.2 页合并

当删除一行记录时，实际上并没有被物理删除，只是记录被标记（flaged）为删除并且它的空间允许被其他记录声明使用。

当页中删除的记录达到 `MERGE_THRESHOLD`（默认为页的 50%），InnoDB 会开始寻找最靠近的页（前或后）看看是否可以将两个页合并以优化空间使用。

::: tip `MERGE_THRESHOLD`
合并页的阈值，可以自己设置，在创建表或者创建索引时指定。
:::

#### 8.2.3 主键设计原则

- 满足业务需求的情况下，尽量降低主键长度
- 插入数据时，尽量顺序插入，选择 `AUTO_INCREMENT` 自增主键
- 尽量不要使用 UUID 做主键或者是其他自然主键，如身份证号
- 业务操作时，避免对主键的修改

### 8.3 `ORDER BY` 优化

1. `Using filesort`：通过表的索引或全表扫描，读取满足条件的数据行，然后在排序缓冲区 sort buffer 中完成排序操作，所有不是通过索引直接返回排序结果的排序都叫 FileSort 排序
2. `Using index`：通过有序索引顺序扫描直接返回有序数据，不需要额外排序，操作效率高

<<< @/db/codes/mysql/optimize_order.sql

- 根据排序字段建立合适的索引，多字段排序时，也遵守最左前缀法则
- 尽量使用覆盖索引
- 多字段排序，一个升序一个降序，此时需要注意联合索引在创建时的规则（`ASC` / `DESC`）
- 如果不可避免的出现 filesort，大数据量排序时，可以适当增大排序缓冲区大小 `sort_buffer_size`（默认 256K）

  <<< @/db/codes/mysql/sort_buffer_size.sql

### 8.4 `GROUP BY` 优化

<<< @/db/codes/mysql/optimize_group.sql

- 在分组操作时，可以通过索引来提升效率
- 分组操作时，索引的使用也是满足最左前缀原则的

### 8.4 `LIMIT` 优化

`LIMIT 2000000, 10` 会查询前 2000010 数据，仅返回 2000001~2000010 的记录，其他记录丢弃，代价很大。

优化思路 1：采用覆盖索引，避免回表查询

<<< @/db/codes/mysql/optimize_limit_cover.sql

优化思路 2：基于索引有序性避免大偏移量（最有效）

<<< @/db/codes/mysql/optimize_limit_where.sql

### 8.5 `COUNT` 优化

- MyISAM 引擎把一个表的总行数存在了磁盘上，因此执行 `COUNT(*)` 的时候会直接返回这个数，效率很高
- InnoDB 引擎在执行 `COUNT(*)` 的时候，需要把数据一行一行读出来，然后累计个数

优化思路：自己计数

#### 8.5.1 `COUNT` 的几种用法

`COUNT()` 是一个聚合函数，对于返回的结果集，一行行地判断，如果参数不为 `NULL`，累计值加一，最后返回累计值。

用法：

- `COUNT(主键)`：InnoDB 引擎会遍历整张表，取出每一行的主键值，返回给服务层。服务层拿到主键后，直接按行进行累加（主键不可能为 `NULL`）
- `COUNT(字段)`：
  - 没有 `NOT NULL` 约束：InnoDB 引擎会遍历整张表，取出每一行的字段值，返回给服务层。服务层判断是否为 `NULL`，不是则累加
  - 有 `NOT NULL` 约束：InnoDB 引擎会遍历整张表，取出每一行的字段值，返回给服务层，直接按行进行累加
- `COUNT(1)`：InnoDB 引擎会遍历整张表，但不取值。服务层对于返回的每一行，放个数字 `1` 进去，直接按行进行累加
- `COUNT(*)`：InnoDB 引擎会遍历整张表，InnoDB 专门做了优化，不会取值，直接按行进行累加

效率排序：`COUNT(*)` ≈ `COUNT(1)` > `COUNT(主键)` > `COUNT(字段)`

### 8.6 `UPDATE` 优化

InnoDB 的行锁是针对索引加的锁，而是针对记录加的锁，并且该索引不能失效，否则会从行锁升级为表锁。

## 9 视图/存储过程/触发器

### 9.1 视图

视图（View）是一种虚拟存在的表。视图中的数据并不在数据库中实际存在，行和列数据来自定义视图的查询中使用的表，并且是在使用视图时动态生成的。

通俗的讲，视图只保存了查询的 SQL 逻辑，不保存查询结果。所以我们在创建视图的时候，主要的工作就落在创建这条 SQL 查询语句上。

- 创建

<<< @/db/codes/mysql/view_create.sql

- 查询

<<< @/db/codes/mysql/view_query.sql

- 修改

<<< @/db/codes/mysql/view_modify.sql

- 删除

<<< @/db/codes/mysql/view_drop.sql

### 9.2 存储过程

存储过程是事先经过编译并存储在数据库中的一段 SQL 语句的集合，调用存储过程可以简化应用开发人员的很多工作，减少数据在数据库和应用服务器之间的传输，对于提高数据处理的效率是有好处的。

存储过程思想上很简单，就是数据库 SQL 语言层面的代码封装与重用。

特点：

- 封装、复用
- 可以接受参数，也可以返回数据
- 减少网络交互，效率提升

#### 9.2.1 基本语法

- 创建

<<< @/db/codes/mysql/procedure_create.sql

- 调用

<<< @/db/codes/mysql/procedure_call.sql

- 查看

<<< @/db/codes/mysql/procedure_query.sql

- 删除

<<< @/db/codes/mysql/procedure_drop.sql

#### 9.2.2 系统变量

系统变量是 MySQL 服务器提供，不是用户定义的，属于服务器层面。分为全局变量（`GLOBAL`）、会话变量（`SESSION`，默认）。

- 查看系统变量

<<< @/db/codes/mysql/variables_query.sql

- 设置系统变量

<<< @/db/codes/mysql/variables_set.sql

#### 9.2.3 用户定义变量

用户定义变量是用户根据需要自己定义的变量，用户变量不用提前声明，在用的时候直接用 `@变量名` 使用就可以。其作用域为当前连接。

#### 9.2.4 局部变量

局部变量是根据需要定义的在局部生效的变量，访问之前，需要 `DECLARE` 声明。可用作存储过程内的局部变量和输入参数，局部变量的范围是在其声明的 `BEGIN ... END` 块。

- 声明

<<< @/db/codes/mysql/variables_declare.sql{4}

- 赋值

<<< @/db/codes/mysql/variables_local_set.sql{7}

#### 9.2.5 `IF`

<<< @/db/codes/mysql/procedure_if.sql{6-12}

#### 9.2.6 `WHILE`

<<< @/db/codes/mysql/procedure_while.sql{4-8}

#### 9.2.7 `REPEAT`

`REPEAT` 是有条件的循环控制语句，当满足条件的时候退出循环。

<<< @/db/codes/mysql/procedure_repeat.sql{4-7}

#### 9.2.8 `LOOP`

`LOOP` 实现简单的循环，如果不在 SQL 逻辑中增加退出循环的条件，可以用其来实现简单的死循环。`LOOP` 可以配合以下两个语句使用：

- `LEAVE`：配合循环使用，退出循环
- `ITERATE`：必须用在循环中，作用是跳过当前循环剩下的语句，直接进入下一次循环

<<< @/db/codes/mysql/procedure_loop.sql{5-12,20-30}

### 9.3 触发器

触发器是与数据库表有关的数据对象，指在 `INSERT` / `UPDATE` / `DELETE` 之前或之后，触发并执行触发器中定义的 SQL 语句集合。触发器的这种特性可以协助应用在数据库端确保数据的完整性、日志记录、数据校验等操作。

使用别名 `OLD` 和 `NEW` 来引用触发器中发生变化的记录内容，这与其他的数据库是相似的。现在触发器还只支持行级触发，不支持语句级触发。

| 触发器类型        | `OLD` 和 `NEW`                                           |
| ----------------- | -------------------------------------------------------- |
| `INSERT` 型触发器 | `NEW` 表示将要或已经新增的数据                           |
| `UPDATE` 型触发器 | `OLD` 表示修改之前的数据，`NEW` 表示将要或已经修改的数据 |
| `DELETE` 型触发器 | `OLD` 表示将要或已经删除的数据                           |

## 10 锁

按锁的粒度可分为全局锁、表级锁、行级锁。

### 10.1 全局锁

全局锁就是对整个数据库实例加锁，加锁后整个实例就处于只读状态，后续的 DML、DDL 语句，已经更新操作的事务提交语句都将被阻塞。

典型使用场景是做全库的逻辑备份，对所有表进行锁定，从而获取一致性视图，保证数据完整性。

1. 加锁

<<< @/db/codes/mysql/lock_global_lock.sql

2. 备份

<<< @/db/codes/mysql/mysqldump.sh

3. 解锁

<<< @/db/codes/mysql/lock_global_unlock.sql

数据库中加全局锁，是一个比较重的操作，存在以下问题：

1. 如果在主库上备份，那么在备份期间都不能执行更新，业务基本上就得停摆。
2. 如果在从库上备份，那么在备份期间从库不能执行主库同步过来的二进制日志（binlog），会导致主从延迟。

在 InnoDB 引擎中，我们可以在备份时加上参数 `--single-transaction` 参数来完成不加锁的一致性数据备份。

<<< @/db/codes/mysql/mysqldump_single_transaction.sh

### 10.2 表级锁

每次操作锁住整张表。锁定粒度大，发生锁冲突的概率高，并发度低。应用在 MyISAM、InnoDB、BDB 等存储引擎中。

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

#### 10.2.2 元数据锁

MDL 加锁过程中是系统自动控制的，无需显式使用，在访问一张表的时候会自动加上。MDL 锁主要作用是维护表元数据的数据一致性，在表上有活动事务的时候，不可以对元数据进行写入操作。为了避免 DML 与 DDL 的冲突，保证读写的正确性。

在 MySQL 5.5 中引入了 MDL，当对一张表进行增删改查的时候，加 MDL 读锁（共享）；当对表结构进行变更操作的时候，加 MDL 写锁（排他）。

查看元数据锁：

<<< @/db/codes/mysql/mdl_select.sql

#### 10.2.3 意向锁

为了避免 DML 在执行时加的行锁与表锁冲突，在 InnoDB 中引入了意向锁，使得表锁不用检查每行数据是否加锁，使用意向锁来减少表锁的检查。

1. 意向共享锁（IS）：与表锁共享锁（read）兼容，与表排它锁（write）互斥。`SELECT ... LOCK IN SHARE MODE`
2. 意向排他锁（IX）：与表锁共享锁（read）及排它锁（write）都互斥。`... FOR UPDATE`

查看意向锁：

<<< @/db/codes/mysql/il_select.sql

### 10.3 行级锁

每次操作锁住对应的行数据。锁定粒度最小，发生锁冲突的概率最低，并发度最高。应用在 InnoDB 存储引擎中。

InnoDB 的数据是基于索引组织的，行锁是通过对索引上的索引项加锁来实现的，而不是对记录加锁。对于行级锁，主要分为以下三类：

1. 行锁（Record Lock）：锁定单个行记录的锁，防止其他事务对此行进行 `UPDATE` 和 `DELETE`。在 `RC`、`RR` 隔离级别下都支持。
2. 间隙锁（Gap Lock）：锁定索引记录间隙（不含该记录），确保索引记录间隙不变，防止其他事务在这个间隙 `INSERT`，产生幻读。在 `RR` 隔离级别下都支持。
3. 临键锁（Next-Key Lock）：行锁和间隙锁的组合，同时锁住数据，并且锁住数据前面的间隙 Gap。在 `RR` 隔离级别下支持。

#### 10.3.1 行锁

InnoDB 实现了以下两种类型的行锁：

1. 共享锁（S）：允许一个事务去读一行，阻止其他事务获得相同数据集的排他锁。
2. 排他锁（X）：允许获取排他锁的事务更新数据，阻止其他事务获得相同数据集的共享锁和排他锁。

| SQL                             | 行锁类型   | 说明                                            |
| ------------------------------- | ---------- | ----------------------------------------------- |
| `INSERT ...`                    | 排他锁     | 自动加锁                                        |
| `UPDATE ...`                    | 排他锁     | 自动加锁                                        |
| `DELETE ...`                    | 排他锁     | 自动加锁                                        |
| `SELECT`（正常）                | 不加任何锁 |                                                 |
| `SELECT ... LOCK IN SHARE MODE` | 共享锁     | 需要手动在 `SELECT` 之后加 `LOCK IN SHARE MODE` |
| `SELECT ... FOR UPDATE`         | 排他锁     | 需要手动在 `SELECT` 之后加 `FOR UPDATE`         |

默认情况下，InnoDB 在 `REPEATABLE READ` 事务隔离级别运行，InnoDB 使用临键锁进行搜索和扫描，以防止幻读。

1. 针对唯一索引进行检索时，对已存在的记录进行等值匹配时，将会自动优化为行锁。
2. InnoDB 的行锁是针对索引加的锁，不通过索引条件检索数据，那么 InnoDB 将对表中的所有记录加锁，此时就会升级为表锁。

查看行锁：

<<< @/db/codes/mysql/lock_raw_select.sql

#### 10.3.2 间隙锁/临键锁

默认情况下，InnoDB 在 `REPEATABLE READ` 事务隔离级别运行，InnoDB 使用临键锁进行搜索和扫描，以防止幻读。

1. 索引上的等值查询（唯一索引），给不存在的记录加锁时，将会优化为间隙锁。
2. 索引上的等值查询（普通索引），向右遍历时最后一个值不满足查询条件时，临键锁退化为间隙锁。
3. 索引上的范围查询（唯一索引），会访问到不满足条件的第一个值为止。

::: warning 注意
间隙锁唯一目的是防止其他事务插入间隙。间隙锁可以共存，一个事务采用的间隙锁不会阻止另一个事务在同一个间隙上采用间隙锁。
:::
