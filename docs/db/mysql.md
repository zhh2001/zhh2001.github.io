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

## 11 InnoDB 引擎

### 11.1 逻辑存储结构

- 表空间（ibd 文件），一个 MySQL 实例可以对应多个表空间，用于存储记录、索引等数据。
- 段，分为数据段（Leaf node segment）、索引段（Non-leaf node segment）、回滚段（Rollback segment），InnoDB 是索引组织表，数据段就是 B+ 树的叶子节点，索引段即为 B+ 树的非叶子节点。段用来管理多个 Extent（区）。
- 区，表空间的单元结构，每个区的大小为 1M。默认情况下，InnoDB 的存储引擎页大小为 16K，即一个区中一共有 64 个连续的页。
- 页，是 InnoDB 存储引擎磁盘管理的最小单元，每个页的大小默认为 16KB。为了保证页的连续性，InnoDB 存储引擎每次从磁盘申请 4-5 个区。
- 行，InnoDB 存储引擎数据是按行进行存放的。

::: tip 行的隐藏字段

- `Trx_id`：每次对某条记录进行改动时，都会把对应的事务 id 赋值给 trx_id 隐藏列。
- `Roll_pointer`：每次对某条引记录进行改动时，都会把旧的版本写入到 undo 日志中，然后这个隐藏列就相当于一个指针，可以通过它找到该纪录修改前的信息。
  :::

### 11.2 架构

MySQL 5.5 版本开始，默认使用 InnoDB 存储引擎，它擅长事务处理，具有崩溃恢复特性，在日常开发中使用非常广泛。下面是 InnoDB 架构图，左侧为内存结构，右侧为磁盘结构。

#### 11.2.1 内存架构

Buffer Pool：缓冲池是主内存中的一个区域，里面可以缓存磁盘上经常操作的真实数据，在执行增删改查操作时，先操作缓冲池中的数据（若缓冲池没有数据，则从磁盘加载并缓存），然后再以一定频率刷新到磁盘，从而减少磁盘 IO，加快处理速度。

缓冲池以 Page 页为单位，底层采用链表数据结构管理 Page。根据状态，将 Page 分为三种类型：

- free page：空闲 page，未被使用。
- clean page：被使用 page，数据没有被修改过。
- dirty page：脏页，被使用 page，数据被修改过，页中数据与磁盘数据不一致

Change Buffer：更改缓冲区（针对非唯一二级索引页），在执行 DML，如果这些数据 Page 没有在 Buffer Pool 中，不会直接操作磁盘，而会将数据变更存在更改缓冲区 Change Buffer 中，在未来数据被读取时，再将数据合并恢复到 Buffer Pool 中，再将合并后的数据刷新到磁盘中。

Adaptive Hash Index：自适应 Hash 索引，用于优化对 Buffer Pool 数据的查询。InnoDB 存储引擎会监控对表上各索引页的查询，如果观察到 hash 索引可以提升速度，则建立 hash 索引，称之为自适应 hash 索引。

Log Buffer：日志缓冲区，用来保存要写入到磁盘中的 log 日志数据（redo log、undo log），默认大小为 16 MB，日志缓冲区的日志会定期刷新到磁盘中。如果需要更新、插入或删除许多行的事务，增加日志缓冲区的大小可以节省磁盘 IO。

#### 11.2.2 磁盘架构

System Tablespace：系统表空间是更改缓冲区的存储区域。如果表是在系统表空间而不是每个表文件或通用表空间中创建的，它也可能包含表和索引数据。（在 MySQL 5.x 版本中还包含 InnoDB 数据字典、undolog 等）

File-Per-Table Tablespaces：每个表的文件表空间包含单个 InnoDB 表的数据和索引，并存储在文件系统上的单个数据文件中。

General Tablespaces：通用表空间，需要通过 `CREATE TABLESPACE` 语法创建通用表空间，在创建表时，可以指定该表空间。

<<< @/db/codes/mysql/create_tablespace.sql

Undo Tablespaces：撤销表空间，MySQL 实例在初始化时会自动创建两个默认 undo 表空间（初始大小 16M），用于存储 undo log 日志。

Temporary Tablespaces：InnoDB 使用会话临时表空间和全局临时表空间。存储用户创建的临时表等数据。

Doublewrite Buffer Files：双写缓冲区，InnoDB 引擎将数据页从 Buffer Pool 刷新到磁盘前，先将数据页写入双写缓冲区文件中，便于系统异常时恢复数据。

Redo Log：重做日志，是用来实现事务的持久性。该日志由两部分组成：重做日志缓冲区（redo log buffer）以及重做日志文件（redo log），前者位于内存中，后者存储在磁盘中。当事务提交后，会将所有修改信息记录到该日志里，用于在刷新脏页到磁盘时发生错误的场景下，进行数据恢复。

#### 11.2.3 后台线程

1. Master Thread  
   核心后台线程，负责调度其他线程，还负责将缓冲池中的数据异步刷新到磁盘中，保持数据的一致性，还包括脏页的刷新、合并插入缓存、undo 页的回收。

2. IO Thread  
   在 InnoDB 存储引擎中大量使用了 AIO 来处理 IO 请求，这样可以极大地提高数据库的性能，而 IO Thread 主要负责这些 IO 请求的回调。

| 线程类型             | 默认个数 | 职责                         |
| -------------------- | -------- | ---------------------------- |
| Read thread          | 4        | 负责读操作                   |
| Write thread         | 4        | 负责写操作                   |
| Log thread           | 1        | 负责将日志缓冲区刷新到磁盘   |
| Insert buffer thread | 1        | 负责将写缓冲区内容刷新到磁盘 |

<<< @/db/codes/mysql/show_engine.sql

3. Purge Thread  
   主要用于回收事务已经提交了的 undo log，在事务提交之后，undo log 可能不用了，就用它来回收。

4. Page Cleaner Thread  
   协助 Master Thread 刷新脏页到磁盘的线程，它可以减轻 Master Thread 的工作压力，减少阻塞。

### 11.3 事务原理

#### 11.3.1 redo log

重做日志，记录的是事务提交时数据页的物理修改，是用来实现事务的持久性。

该日志文件由两部分组成：重做日志缓冲（redo log buffer）以及重做日志文件（redo log file），前者是在内存中，后者在磁盘中。当事务提交之后会把所有修改信息都存到该日志文件中，用于刷新脏页到磁盘，发生错误时，进行数据恢复使用。

#### 11.3.2 undo log

回滚日志，用于记录数据被修改前的信息，作用包含两个：提供回滚和 MVCC（多版本并发控制）。

undo log 和 redo log 记录物理日志不一样，它是逻辑日志。可以认为当 `DELETE` 一条记录时，undo log 中会记录一条对应的 `INSERT` 记录，反之亦然，当 `UPDATE` 一条记录时，它记录一条对应相反的 `UPDATE` 记录。当执行 `ROLLBACK` 时，就可以从 undo log 中的逻辑记录读到相应的内容并进行回滚。

Undo log 销毁：undo log 在事务执行时产生，事务提交时，并不会立即删除 undo log，因为这些日志可能还用于 MVCC。

Undo log 存储：undo log 采用段的方式进行管理和记录，存放在 rollback segment 回滚段中，内部包含 1024 个 undo log segment。

### 11.4 MVCC

- 当前读

  读取的是记录的最新版本，读取时还要保证其他并发事务不能修改当前记录，会对读取的记录进行加锁。对于我们日常的操作，如：`SELECT ... LOCK IN SHARE MODE`（共享锁），`SELECT ... FOR UPDATE`、`UPDATE`、`INSERT`、`DELETE`（排他锁）都是一种当前读。

- 快照读

  简单的 `SELECT`（不加锁）就是快照读，读取的是记录数据的可见版本，有可能是历史数据，不加锁，是非阻塞读。

  - Read Committed：每次 SELECT，都生成一个快照读。
  - Repeatable Read：开启事务后第一个 SELECT 语句才是快照的地方。
  - Serializable：快照读退化为当前读。

- MVCC

  全称 Multi-Version Concurrency Control，多版本并发控制。指维护一个数据的多个版本，使得读写操作没有冲突，快照读为 MySQL 实现 MVCC 提供了一个非阻塞读功能。MVCC 的具体实现，还需要依赖于数据库记录的三个隐式字段、undo log、readView。

#### 11.4.1 隐藏字段

| 隐藏字段    | 含义                                                                  |
| ----------- | --------------------------------------------------------------------- |
| DB_TRX_ID   | 最近修改事务 ID，记录插入这条记录或最后一次修改该记录的事务 ID        |
| DB_ROLL_PTR | 回滚指针，指向这条记录的上一个版本，用于配合 undo log，指向上一个版本 |
| DB_ROW_ID   | 隐藏主键，如果表结构没有指定主键，将会生成该隐藏字段                  |

#### 11.4.2 undo log

回滚日志，在 `INSERT`、`UPDATE`、`DELETE` 的时候产生的便于数据回滚的日志。

当 `INSERT` 的时候，产生的 undo log 日志只在回滚时需要，在事务提交后，可被立即删除。

而 `UPDATE`、`DELETE` 的时候，产生的 undo log 日志不仅在回滚时需要，在快照读也需要，不会立即删除。

#### 11.4.3 undo log 版本链

不同事务或相同事务对同一记录进行修改，会导致该记录的 undolog 生成一条记录版本链表，链表头部是最新的旧纪录，链表尾部是最早的旧纪录。

#### 11.4.4 readview

ReadView（读视图）是快照读 SQL 执行时 MVCC 提取数据的依据，记录并维护系统当前活跃的事务（未提交的）id。

ReadView 中包含了四个核心字段:

| 字段             | 含义                                                     |
| ---------------- | -------------------------------------------------------- |
| `m_ids`          | 当前活跃的事务 ID 集合                                   |
| `min_trx_id`     | 最小活跃事务 ID                                          |
| `max_trx_id`     | 预分配事务 ID，当前最大事务 ID+1（因为事务 ID 是自增的） |
| `creator_trx_id` | ReadView 创建者的事务 ID                                 |

1. `trx_id == creator_trx_id`：<span style="color:green;">可以访问该版本</span>
2. `trx_id < min_trx_id`：<span style="color:green;">可以访问该版本</span>
3. `trx_id > max_trx_id`：<span style="color:red;">不可以访问该版本</span>
4. `min_trx_id <= trx_id <= max_trx_id`：如果 `trx_id` 不在 `m_ids` 中是<span style="color:green;">可以</span>访问该版本的

不同的隔离级别，生成 ReadView 的时机不同：

- READ COMMITTED：在事务中每一次执行快照读时生成 ReadView。
- REPEATABLE READ：仅在事务中第一次执行快照读时生成 ReadView，后续复用该 ReadView。

## 12 MySQL 管理

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
  - `-p, --password [=name]`：指定密码
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
  - `-d, --database=name`：指定数据库名称，只列出指定的数据库相关操作。
  - `-o, --offset=#`：忽略掉日志中的前 `n` 行命令。
  - `-r, --result-file=name`：将输出的文本格式日志输出到指定文件。
  - `-s, --short-form`：显示简单格式，省略掉一些信息。
  - `--start-datetime=date1 --stop-datetime=date2`：指定日期间隔内的所有日志。
  - `--start-position=pos1 --stop-position=pos2`：指定位置间隔内的所有日志。

<<< @/db/codes/mysql/mysqlbinlog.sh

#### 12.2.4 `mysqlshow`

`mysqlshow` 客户端对象查找工具，用来快速查找存在哪些数据库、数据库中的表、表中的列或者索引。

- 语法：`mysqlshow [options] [db_name [table_name [col_name]]]`
- 选项：
  - `--count`：显示数据库及表的统计信息（数据库，表均可以不指定）
  - `-i`：显示指定数据库或者指定表的状态信息

示例：

<<< @/db/codes/mysql/mysqlshow.sh

<<< @/db/codes/mysql/mysqlshow_db.sh

<<< @/db/codes/mysql/mysqlshow_db_tb.sh

#### 12.2.5 `mysqldump`

`mysqldump` 是 MySQL 的客户端备份工具，主要用于：

- **数据库备份**：生成包含建表语句、数据插入语句的 SQL 文件；
- **数据迁移**：在不同 MySQL 实例（或数据库）间转移数据。

#### 12.2.6 `mysqlimport`/`source`

`mysqlimport` 是客户端数据导入工具，用于导入 `mysqldump` 加 `-T` 参数后导出的文本文件（`-T` 会生成表结构 SQL 和对应的数据文本文件）。

`source` 用于导入 SQL 文件（如 `mysqldump` 导出的完整备份 SQL）。

## 13 日志

### 13.1 错误日志

错误日志是 MySQL 中最重要的日志之一，它记录了 mysqld 启动和停止时，以及服务器在运行过程中发生任何严重错误时的相关信息。当数据库出现任何故障无法正常使用时，建议首先查看此日志。

该日志默认开启，默认存放在 `/var/log/mysql/`，默认日志文件名为 `mysqld.log`。查看日志位置：

<<< @/db/codes/mysql/show_log_err.sh

### 13.2 二进制日志

二进制日志（BINLOG）记录了所有的 DDL（数据定义语言）语句和 DML（数据操纵语言）语句，但不包含查询（`SELECT`、`SHOW`）语句。

作用：

1. 灾难时的数据恢复；
2. MySQL 的主从复制。在 MySQL 8 版本中，默认二进制日志是开启的，涉及到的参数如下：

<<< @/db/codes/mysql/show_log_bin.sh

MySQL 服务器中提供了多种格式来记录二进制日志，具体格式及特点如下：

| 日志格式  | 含义                                                                                              |
| :-------: | ------------------------------------------------------------------------------------------------- |
| STATEMENT | 基于 SQL 语句的日志记录，记录的是 SQL 语句，对数据进行修改的 SQL 都会记录在日志文件中。           |
|    ROW    | 基于行的日志记录，记录的是每一行的数据变更。（默认）                                              |
|   MIXED   | 混合了 STATEMENT 和 ROW 两种格式，优先采用 STATEMENT，在某些特殊情况下会自动切换为 ROW 进行记录。 |

查看当前二进制日志格式的语句：

<<< @/db/codes/mysql/show_log_bin_fmt.sh

#### 13.2.1 日志查看

由于日志是以二进制方式存储的，不能直接读取，需要通过二进制日志查询工具 `mysqlbinlog` 来查看。

- 语法：`mysqlbinlog [参数选项] logfilename`
- 参数选项：
  - `-d`：指定数据库名称，只列出指定的数据库相关操作。
  - `-o`：忽略掉日志中的前 `n` 行命令。
  - `-v`：将行事件（数据变更）重构为 SQL 语句。
  - `-vv`：将行事件（数据变更）重构为 SQL 语句，并输出注释信息。

#### 13.2.2 日志删除

对于比较繁忙的业务系统，每天生成的 binlog 数据巨大，如果长时间不清除，将会占用大量磁盘空间。可以通过以下几种方式清理日志：

| 指令                                               | 含义                                                                |
| -------------------------------------------------- | ------------------------------------------------------------------- |
| `reset master`                                     | 删除全部 binlog 日志，删除之后，日志编号将从 binlog.000001 重新开始 |
| `purge master logs to 'binlog.******'`             | 删除 `******` 编号之前的所有日志                                    |
| `purge master logs before 'yyyy-mm-dd hh24:mi:ss'` | 删除指定日期时间之前产生的所有日志                                  |

也可以在 MySQL 的配置文件中配置二进制日志的过期时间，设置了以后，二进制日志过期会自动删除。

<<< @/db/codes/mysql/show_log_bin_expire.sh

### 13.3 查询日志

查询日志记录了客户端的所有操作语句，而二进制日志不包含查询数据的 SQL 语句。默认情况下，查询日志是未开启的。如果需要开启查询日志，可以设置以下配置：

<<< @/db/codes/mysql/show_log_general.sh

修改 MySQL 的配置文件，添加以下内容：

<<< @/db/codes/mysql/log_general.cnf

### 13.4 慢查询日志

见 [7.4.2 慢查询日志](#_7-4-2-慢查询日志)

## 14 主从复制

主从复制是指将主数据库的 DDL 和 DML 操作通过二进制日志传到从库服务器中，然后在从库上对这些日志重新执行（也叫重做），从而使得从库和主库的数据保持同步。

MySQL 支持一台主库同时向多台从库进行复制，从库同时也可以作为其他从服务器的主库，实现链状复制。

MySQL 复制的优点包含：

1. 主库出现问题，可以快速切到从库提供服务
2. 实现读写分离，降低主库的访问压力
3. 可以在从库中执行备份，避免备份期间影响主库服务

原理：

1. Master 主库在提交事务时，会把数据变更记录在二进制日志文件 Binlog 中。
2. 从库读取主库的二进制文件 Binlog，写入到从库的中继日志 Relay Log。
3. Slave 重做中继日志中的事件，将改变反映到它自己的数据。

## 15 分库分表

随着互联网的发展，应用系统的数据量也呈指数增长，若采用单数据库进行数据存储，存在以下性能瓶颈：

1. IO 瓶颈：热点数据太多，数据库缓存不足，产生大量磁盘 IO，效率较低。请求数据太多，带宽不够，网络 IO 瓶颈。
2. CPU 瓶颈：排序、分组、连接查询、聚合统计等 SQL 会耗费大量的 CPU 资源，请求数太多了，CPU 出现瓶颈。

分库分表的中心思想都是将数据分散存储，使得单一数据库/表的数据量变小来缓解单一数据库的性能问题，从而达到提升数据库性能的目的。

### 15.1 拆分策略

#### 15.1.1 垂直分库

以表为依据，根据业务将不同的表拆分到不同库中。特点如下：

1. 每个库的表结构都不一样。
2. 每个库的数据也不一样。
3. 所有库的并集是全量数据。

#### 15.1.2 垂直分表

以字段为依据，根据字段属性将不同字段拆分到不同表中。特点如下：

1. 每个表的结构都不一样。
2. 每个表的数据也不一样，一般通过一列（主键/外键）关联。
3. 所有表的并集是全量数据。

#### 15.1.3 水平分库

特点：

1. 每个库的表结构都一样。
2. 每个库的数据都不一样。
3. 所有库的并集是全量数据。

#### 15.1.4 水平分表

特点：

1. 每个表的表结构都一样。
2. 每个表的数据都不一样。
3. 所有表的并集是全量数据。

### 15.2 实现技术

- **shardingJDBC**：基于 AOP 原理，在应用程序中对本地执行的 SQL 进行拦截、解析、改写、路由处理。需要自行编码配置实现，仅支持 Java 语言，性能较高。
- **MyCat**：数据库分库分表中间件，无需修改业务代码（通过中间件配置）即可实现分库分表，支持多种语言，性能不及 shardingJDBC。
- **Gorm Sharding**：高性能的数据库分表中间件。

### 15.3 Gorm Sharding

它基于 Conn 层做 SQL 拦截、AST 解析、分表路由、自增主键填充，带来的额外开销极小。对开发者友好、透明，使用上与普通 SQL、Gorm 查询无差别，只需要额外注意一下分表键条件。

特性：

1. 非侵入式设计。加载插件，指定配置，即可完成。
2. 速度极快。无需基于网络的中间件。
3. 支持多种数据库（PostgreSQL、MySQL）。
4. 集成主键生成器（Snowflake、PostgreSQL 序列、自定义等）。

#### 15.3.1 安装

<<< @/db/codes/mysql/sharding_install.sh

#### 15.3.2 用法

配置分片中间件，注册需要分片的表：

<<< @/db/codes/mysql/sharding_config.go

::: warning
在多个节点上使用默认的 Snowflake 生成器可能会导致主键冲突，请使用自定义主键生成器，或在发生冲突时重新生成主键。
:::

### 15.4 分片方式

1. 范围分片：根据指定的字段所在的范围与节点的对应关系，来决定该数据属于哪一个分片。
2. 取模分片：根据指定的字段值与节点数量进行求模运算，根据运算结果，来决定该数据属于哪一个分片。
3. 一致性 hash：所谓一致性哈希，相同的哈希因子计算值总是被划分到相同的分区表中，不会因为分区节点的增加而改变原来数据的分区位置。
4. 枚举分片：通过设置可能得到枚举值，指定数据分布到不同的数据节点上，本规则适用于按照省份、性别、状态拆分数据等业务。
5. 按日期分片
