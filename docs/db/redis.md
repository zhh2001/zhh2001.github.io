---
title: Redis 数据类型、缓存与底层实现
description: 整理 Redis 常用数据类型与命令，介绍缓存问题、持久化、主从复制、哨兵和分片，结合实现说明数据结构、网络模型、通信协议及内存回收。
outline: [2, 3]
---

# Redis

本文以 Redis 8.2.x 为主，命令、配置和源码片段使用 Redis 8.2.3 核对。涉及旧版实现时单独注明版本。底层 C 源码片段用于说明结构，部分省略了无关字段或分支。

数据操作示例在 `redis-cli` 交互提示符内输入，配置和建群章节标注的 shell 命令在终端执行。各组数据操作在独立的空数据库中分别执行，展示的是 `redis-cli` 的 RESP2 格式化输出。

命令参数可查阅[官方命令参考](https://redis.io/docs/latest/commands/)，版本范围以正文说明为准。

## 1 通用命令

### 1.1 `KEYS`

- 语法：`KEYS pattern`
- 功能：一次返回所有匹配模式的 `key`
- 时间复杂度：`O(N)`，其中 `N` 是数据库中的键数

`KEYS` 会遍历整个数据库，大键空间下可能阻塞其他请求。日常遍历优先使用 `SCAN cursor [MATCH pattern] [COUNT count]`，从游标 `0` 开始，直到返回游标再次为 `0`。`COUNT` 是工作量提示，单次返回数量不固定，也可能返回空集合但尚未结束。完整迭代可能重复返回键，遍历期间发生的增删不能按一致性快照理解。

```bash
MSET firstname Jack lastname Stuntman age 35
# "OK"
KEYS *name*
# 1) "firstname"
# 2) "lastname"
KEYS a??
# 1) "age"
KEYS *
# 1) "age"
# 2) "firstname"
# 3) "lastname"
```

### 1.2 `DEL`

- 语法：`DEL key [key ...]`
- 功能：删除指定的 `key`，如果 `key` 不存在则忽略

返回实际删除的键数。删除大型集合可能需要同步释放大量内存，Redis 4.0 起可用 `UNLINK` 将键移出键空间，再按对象情况异步释放内存。

```bash
SET key1 "Hello"
# "OK"
SET key2 "World"
# "OK"
DEL key1 key2 key3
# (integer) 2
```

### 1.3 `EXISTS`

- 语法：`EXISTS key [key ...]`
- 功能：返回参数中存在的键数，重复传入同一个存在的键会重复计数

```bash
SET key1 "Hello"
# "OK"
EXISTS key1
# (integer) 1
EXISTS nosuchkey
# (integer) 0
SET key2 "World"
# "OK"
EXISTS key1 key2 nosuchkey
# (integer) 2
```

### 1.4 `EXPIRE`

- 语法：`EXPIRE key seconds [NX | XX | GT | LT]`
- 功能：设置 `key` 的剩余存活时间，单位为秒，成功返回 `1`，键不存在或条件不满足返回 `0`

Redis 7.0 起支持条件选项：`NX` 仅给没有过期时间的键设置，`XX` 仅更新已有过期时间，`GT` 和 `LT` 分别要求新的过期时间更晚或更早。比较时，没有过期时间的键视为无限期。设置成功时，非正时长会立即删除键。`INCR`、`HSET` 等修改值的操作通常保留键的 TTL，`SET` 等覆盖操作则默认移除 TTL。

### 1.5 `TTL`

- 语法：`TTL key`
- 功能：查看指定 `key` 的剩余有效时长（秒）

如果 `key` 存在但是没有设置过期时长，返回 `-1`。如果 `key` 不存在返回 `-2`。

## 2 String 类型

String 类型，也就是字符串类型，是 Redis 中最简单的存储类型。

String 保存二进制安全的字节序列，也能保存文本形式的整数或浮点数。整数、浮点数不是独立的 Redis 数据类型，只是部分命令会按数值解释字符串内容。

Redis 8.2.x 默认将单个字符串及请求中的 bulk string 限制为 512 MiB，相关限制由 `proto-max-bulk-len` 控制。底层可以采用 `int`、`embstr` 或 `raw` 编码，整数编码直接保存数值。

### 2.1 `SET`

- 语法：`SET key value [NX | XX] [GET] [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | KEEPTTL]`
- 功能：添加或修改一个 String 类型的键值对
- 可选项：
  - `NX`：只有 `key` 不存在时才设置
  - `XX`：只有 `key` 已存在时才设置
  - `GET`：返回写入前的字符串值，不存在时返回空值。旧值不是 String 时返回错误
  - `EX`、`PX`：分别按秒、毫秒设置剩余存活时间
  - `EXAT`、`PXAT`：分别按秒、毫秒指定绝对到期时间
  - `KEEPTTL`：保留 `key` 原有的过期时长

普通 `SET` 可以覆盖其他类型的键，成功后默认清除原 TTL。条件不满足时不写入，不带 `GET` 时返回空值。`KEEPTTL` 在 Redis 6.0 加入，`GET`、`EXAT`、`PXAT` 在 6.2 加入，`NX` 与 `GET` 从 7.0 起可以一起使用。

### 2.2 `GET`

- 语法：`GET key`
- 功能：根据 `key` 获取 String 类型的 `value`

键不存在时返回空值，键存在但不是 String 时返回 `WRONGTYPE` 错误。

### 2.3 `MSET`

- 语法：`MSET key value [key value ...]`
- 功能：批量添加多个 String 类型的键值对

同一条 `MSET` 原子地写入所有键，覆盖旧值并移除被覆盖键的 TTL。Redis Cluster 要求这些键属于同一个槽。

### 2.4 `MGET`

- 语法：`MGET key [key ...]`
- 功能：批量获取多个 `key` 的 `value`

结果与参数顺序对应，不存在的键和非 String 类型的键都返回空值。

### 2.5 `INCR`

- 语法：`INCR key`
- 功能：整型自增 `1`

键不存在时按 `0` 初始化。`INCR` 和 `INCRBY` 使用有符号 64 位整数，内容不能解析为整数或计算溢出时返回错误。

### 2.6 `INCRBY`

- 语法：`INCRBY key increment`
- 功能：整型自增 `increment`

步长可以为负数，也可用 `DECR`、`DECRBY` 表达减法。

### 2.7 `INCRBYFLOAT`

- 语法：`INCRBYFLOAT key increment`
- 功能：浮点型自增 `increment`

步长可以为负数，键不存在时按 `0` 初始化。计算结果不能是 NaN 或无穷大。不应依赖浮点运算提供十进制金额的精确表示。

### 2.8 `SETNX`

> Redis 2.6.12 起不推荐用于新代码，改用 `SET key value NX`。命令仍可使用。

- 语法：`SETNX key value`
- 功能：如果 `key` 不存在才新增。

`SETNX` 是 **SET** if **N**ot e**X**ists 的简写。

需要同时设置过期时间时使用 `SET key value NX EX seconds`，避免将 `SETNX` 和 `EXPIRE` 分成两次操作。

### 2.9 `SETEX`

> Redis 2.6.12 起不推荐用于新代码，改用 `SET key value EX seconds`。命令仍可使用。

- 语法：`SETEX key seconds value`
- 功能：新增 `key` 并设置有效时长

## 3 Key 的层级格式

通常用 `业务:对象:ID` 命名键，例如 `login:user:101`。冒号只是便于阅读的命名约定，Redis 键空间本身是平面的，不会建立目录或父子关系。键也是二进制安全的字节序列。

## 4 Hash 类型

Hash 类型保存无序的字段与值映射。字段和值都是二进制安全的字符串，同一个 Hash 中字段名唯一，不保证遍历顺序。

对象可以序列化后整体保存为 String，也可以将字段分别放入 Hash。后者便于单独读取和更新字段。

### 4.1 `HSET`

- 语法：`HSET key field value [field value ...]`
- 功能：添加或修改 Hash 的字段值，返回新增字段数，更新已有字段不计入返回值

### 4.2 `HGET`

- 语法：`HGET key field`
- 功能：读取字段值，键或字段不存在时返回空值

### 4.3 `HMSET`

> Redis 4.0 起不推荐用于新代码，改用支持多个字段的 `HSET`。两者写入效果相同，但 `HMSET` 返回 `OK`，`HSET` 返回新增字段数。

- 语法：`HMSET key field value [field value ...]`
- 功能：批量添加多个 hash 类型 `key` 的 `field` 的值

```bash
HMSET myhash field1 "Hello" field2 "World"
# "OK"
HGET myhash field1
# "Hello"
HGET myhash field2
# "World"
```

### 4.4 `HMGET`

- 语法：`HMGET key field [field ...]`
- 功能：按参数顺序读取多个字段，不存在的键或字段对应空值

```bash
HSET myhash field1 "Hello" field2 "World"
# (integer) 2
HMGET myhash field1 field2 nofield
# 1) "Hello"
# 2) "World"
# 3) (nil)
```

### 4.5 `HGETALL`

- 语法：`HGETALL key`
- 功能：返回所有字段和值。RESP2 是字段和值交替排列的数组，RESP3 是映射，均不保证顺序

```bash
HSET myhash field1 "Hello" field2 "World"
# (integer) 2
HGETALL myhash
# 1) "field1"
# 2) "Hello"
# 3) "field2"
# 4) "World"
```

### 4.6 `HKEYS`

- 语法：`HKEYS key`
- 功能：返回所有字段名，不保证顺序

```bash
HSET myhash field1 "Hello" field2 "World"
# (integer) 2
HKEYS myhash
# 1) "field1"
# 2) "field2"
```

### 4.7 `HVALS`

- 语法：`HVALS key`
- 功能：返回所有字段值，不保证顺序

```bash
HSET myhash field1 "Hello" field2 "World"
# (integer) 2
HVALS myhash
# 1) "Hello"
# 2) "World"
```

### 4.8 `HINCRBY`

- 语法：`HINCRBY key field increment`
- 功能：将字段值按有符号 64 位整数增加指定步长，字段不存在时从 `0` 开始，不能解析为整数或溢出时返回错误

```bash
HSET myhash field 5
# (integer) 1
HINCRBY myhash field 1
# (integer) 6
HINCRBY myhash field -10
# (integer) -4
```

### 4.9 `HSETNX`

- 语法：`HSETNX key field value`
- 功能：只有这个 `key` 的字段不存在才能设置

```bash
HSETNX myhash field "Hello"
# (integer) 1
HSETNX myhash field "World"
# (integer) 0
HGET myhash field
# "Hello"
```

## 5 List 类型

List 是有序的字符串列表，允许重复元素，支持从两端插入和弹出，也支持按索引访问。底层编码见后面的数据结构章节。

### 5.1 `LPUSH`

- 语法：`LPUSH key element [element ...]`
- 功能：向列表左侧插入元素，`key` 不存在则会创建

```bash
LPUSH mylist "world"
# (integer) 1
LPUSH mylist "hello"
# (integer) 2
```

### 5.2 `RPUSH`

- 语法：`RPUSH key element [element ...]`
- 功能：向列表右侧插入元素，`key` 不存在则会创建

```bash
RPUSH mylist "one" "two" "three" "four" "five"
# (integer) 5
```

### 5.3 `LPOP`

- 语法：`LPOP key [count]`
- 功能：从列表左侧移除元素

`count` 在 Redis 6.2 加入。不带 `count` 时返回单个元素，带 `count` 时返回数组。键不存在时返回空值。

```bash
RPUSH mylist "one" "two" "three" "four" "five"
# (integer) 5
LPOP mylist
# "one"
LPOP mylist 2
# 1) "two"
# 2) "three"
```

### 5.4 `RPOP`

- 语法：`RPOP key [count]`
- 功能：从列表右侧移除元素

```bash
RPUSH mylist "one" "two" "three" "four" "five"
# (integer) 5
RPOP mylist
# "five"
RPOP mylist 2
# 1) "four"
# 2) "three"
```

### 5.5 `LRANGE`

- 语法：`LRANGE key start stop`
- 功能：返回索引在 `[start stop]` 内的所有元素

范围包含两端，索引从 `0` 开始，`-1` 表示最后一个元素。超出列表边界的范围会被截取，起始位置在列表末尾之外时返回空数组。

```bash
RPUSH mylist "one" "two" "three"
# (integer) 3
LRANGE mylist 0 0
# 1) "one"
LRANGE mylist -3 2
# 1) "one"
# 2) "two"
# 3) "three"
LRANGE mylist -100 100
# 1) "one"
# 2) "two"
# 3) "three"
LRANGE mylist 5 10
# (empty array)
```

## 6 Set 类型

Set 是无序的字符串集合，可以采用 intset、listpack 或哈希表编码。它在逻辑上与 Java 的 HashSet 类似，具有以下特征：

- 无序
- 元素不重复
- 查找快
- 支持交集、并集、差集等功能

### 6.1 `SADD`

- 语法：`SADD key member [member ...]`
- 功能：添加成员，返回新增成员数，重复成员不计入

```bash
SADD myset "Hello" "World"
# (integer) 2
SADD myset "World"
# (integer) 0
```

### 6.2 `SREM`

- 语法：`SREM key member [member ...]`
- 功能：移除指定成员，返回实际移除数

```bash
SADD myset "one" "two" "three"
# (integer) 3
SREM myset "one" "three"
# (integer) 2
SREM myset "four"
# (integer) 0
```

### 6.3 `SCARD`

- 语法：`SCARD key`
- 功能：返回集合中的元素数量

```bash
SADD myset "one" "two" "three"
# (integer) 3
SCARD myset
# (integer) 3
```

### 6.4 `SISMEMBER`

- 语法：`SISMEMBER key member`
- 功能：判断元素是否在集合中

```bash
SADD myset "one"
# (integer) 1
SISMEMBER myset "one"
# (integer) 1
SISMEMBER myset "two"
# (integer) 0
```

### 6.5 `SMEMBERS`

- 语法：`SMEMBERS key`
- 功能：返回所有成员，不保证顺序。大型集合可用 `SSCAN` 分批遍历，游标语义与 `SCAN` 相同

```bash
SADD myset Hello World
# (integer) 2
SMEMBERS myset
# 1) "Hello"
# 2) "World"
```

### 6.6 `SINTER`

- 语法：`SINTER key [key ...]`
- 功能：求所有集合的交集，任意输入键不存在时交集为空

```bash
SADD s1 a b c d
# (integer) 4
SADD s2 c
# (integer) 1
SADD s3 a c e
# (integer) 3
SINTER s1 s2 s3
# 1) "c"
```

### 6.7 `SDIFF`

- 语法：`SDIFF key [key ...]`
- 功能：返回第一个集合中不属于后续任何集合的成员

```bash
SADD s1 a b c d
# (integer) 4
SADD s2 c
# (integer) 1
SADD s3 a c e
# (integer) 3
SDIFF s1 s2 s3
# 1) "d"
# 2) "b"
```

### 6.8 `SUNION`

- 语法：`SUNION key [key ...]`
- 功能：求所有集合的并集，相同成员只返回一次

```bash
SADD s1 a b c d
# (integer) 4
SADD s2 c
# (integer) 1
SADD s3 a c e
# (integer) 3
SUNION s1 s2 s3
# 1) "c"
# 2) "e"
# 3) "b"
# 4) "d"
# 5) "a"
```

## 7 SortedSet 类型

SortedSet（ZSet）保存成员与分数的映射。成员唯一，分数可以重复，按分数升序排列，分数相同时按成员的二进制字典序排列。小型集合使用 listpack，较大的集合使用跳表和字典。

分数使用 IEEE 754 双精度浮点数，区间 `[-2^53, 2^53]` 内的整数可以精确表示，NaN 不能作为分数。

SortedSet 具备下列特性：

- 可排序
- 元素不重复
- 查询速度快

因为 SortedSet 的可排序特性，经常被用来实现排行榜这样的功能。

### 7.1 `ZADD`

- 语法：`ZADD key [NX | XX] [GT | LT] [CH] [INCR] score member [score member ...]`
- 功能：添加元素到有序集合，如果已存在则更新 `score`
- 可选项：
  - `XX`：仅更新已存在的元素。不添加新元素。
  - `NX`：只添加新元素。不更新现有元素。
  - `GT`、`LT`：仅在新分数更大或更小时更新，不限制新增成员。不能与 `NX` 一起使用
  - `CH`：返回新增和分数发生变化的成员总数，默认仅返回新增数
  - `INCR`：按 `ZINCRBY` 方式增加分数，此时只能指定一组分数与成员

```bash
ZADD myzset 1 one 1 uno 2 two 3 three
# (integer) 4
```

### 7.2 `ZREM`

- 语法：`ZREM key member [member ...]`
- 功能：删除有序集合中的指定元素

```bash
ZADD myzset 1 one 1 uno 2 two 3 three
# (integer) 4
ZREM myzset two
# (integer) 1
```

### 7.3 `ZSCORE`

- 语法：`ZSCORE key member`
- 功能：查询指定元素的 `score`

```bash
ZADD myzset 1 one
# (integer) 1
ZSCORE myzset one
# "1"
```

### 7.4 `ZRANK`

- 语法：`ZRANK key member [WITHSCORE]`
- 功能：获取指定元素按升序排列的排名，从 `0` 开始。`WITHSCORE` 从 Redis 7.2 起可用，同时返回分数

```bash
ZADD z 1 one 2 two 3 three
# (integer) 3
ZRANK z three
# (integer) 2
ZRANK z four
# (nil)
```

### 7.5 `ZCARD`

- 语法：`ZCARD key`
- 功能：获取有序集合中的元素数量

```bash
ZADD z 1 one 2 two 3 three
# (integer) 3
ZCARD z
# (integer) 3
```

### 7.6 `ZCOUNT`

- 语法：`ZCOUNT key min max`
- 功能：获取有序集合中的 `score` 在 `[min, max]` 内的元素数量

边界默认包含，可以用 `(` 排除边界，例如 `ZCOUNT z (1 3`。`-inf` 和 `+inf` 表示无穷边界。

```bash
ZADD z 1 one 2 two 3 three
# (integer) 3
ZCOUNT z -inf +inf
# (integer) 3
ZCOUNT z 2 3
# (integer) 2
```

### 7.7 `ZINCRBY`

- 语法：`ZINCRBY key increment member`
- 功能：让有序集合中指定元素的 `score` 自增，步长为 `increment`

```bash
ZADD myzset 1 "one" 2 "two"
# (integer) 2
ZINCRBY myzset 2 "one"
# "3"
```

### 7.8 `ZRANGE`

- 语法：`ZRANGE key start stop [BYSCORE | BYLEX] [REV] [LIMIT offset count] [WITHSCORES]`
- 功能：按照 `score` 升序排序后，获取指定排名范围内的元素，排名从 0 开始
- 参数：
  - 默认：`start`、`stop` 是包含两端的排名范围，支持负索引
  - `BYSCORE`：改按分数边界查询，支持 `(`、`-inf` 和 `+inf`
  - `BYLEX`：改按成员字典序边界查询，使用 `[`、`(`、`-`、`+` 表示边界。应保证所有成员分数相同，否则结果不具备该模式预期的语义
  - `REV`：反向查询。配合分数或字典序查询时，先写较大边界，再写较小边界
  - `LIMIT`：仅用于 `BYSCORE`、`BYLEX` 查询
  - `WITHSCORES`：同时返回分数

```bash
ZADD z 1 one 2 two 3 three
# (integer) 3
ZRANGE z 1 99
# 1) "two"
# 2) "three"
ZRANGE z 0 1 WITHSCORES
# 1) "one"
# 2) "1"
# 3) "two"
# 4) "2"
```

### 7.9 `ZRANGEBYSCORE`

> Redis 6.2 起不推荐用于新代码，改用 `ZRANGE key min max BYSCORE [WITHSCORES] [LIMIT offset count]`。

- 语法：`ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]`
- 功能：按照 `score` 排序后，获取指定 `score` 范围内的元素

### 7.10 `ZDIFF`

- 语法：`ZDIFF numkeys key [key ...] [WITHSCORES]`
- 功能：求第一个集合相对其余集合的差集，保留第一个集合的分数。Redis 6.2 加入

```bash
ZADD zset1 1 "one" 2 "two" 3 "three"
# (integer) 3
ZADD zset2 1 "one" 2 "two"
# (integer) 2
ZDIFF 2 zset1 zset2
# 1) "three"
ZDIFF 2 zset1 zset2 WITHSCORES
# 1) "three"
# 2) "3"
```

### 7.11 `ZINTER`

- 语法：`ZINTER numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE SUM | MIN | MAX] [WITHSCORES]`
- 功能：求交集。Redis 6.2 加入，分数默认相加，`WEIGHTS` 指定各集合的权重，`AGGREGATE` 指定聚合方式

```bash
ZADD zset1 1 "one" 2 "two"
# (integer) 2
ZADD zset2 1 "one" 2 "two" 3 "three"
# (integer) 3
ZINTER 2 zset1 zset2
# 1) "one"
# 2) "two"
ZINTER 2 zset1 zset2 WITHSCORES
# 1) "one"
# 2) "2"
# 3) "two"
# 4) "4"
```

### 7.12 `ZUNION`

- 语法：`ZUNION numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE SUM | MIN | MAX] [WITHSCORES]`
- 功能：求并集。Redis 6.2 加入，分数默认相加，支持权重及 `SUM`、`MIN`、`MAX` 聚合

```bash
ZADD zset1 1 "one" 2 "two"
# (integer) 2
ZADD zset2 1 "one" 2 "two" 3 "three"
# (integer) 3
ZUNION 2 zset1 zset2
# 1) "one"
# 2) "three"
# 3) "two"
ZUNION 2 zset1 zset2 WITHSCORES
# 1) "one"
# 2) "2"
# 3) "three"
# 4) "3"
# 5) "two"
# 6) "4"
```

## 8 缓存穿透

请求的数据在缓存和数据库中都不存在，如果每次缓存未命中都查询数据库，重复的无效请求就会持续占用数据库资源。

常用处理方式：

- 缓存空结果，设置较短的 TTL，并在数据创建或更新后失效相关缓存。应区分“确实不存在”和查询失败，避免把暂时故障缓存为空结果。
- 使用布隆过滤器预先判断。过滤器可能误判为存在，仍需访问缓存或数据库。没有误判为不存在的前提是相关数据已正确加入过滤器，数据更新和过滤器维护必须协调。
- 校验请求参数，对异常请求限流。

## 9 缓存雪崩

大量键在短时间内过期，或者缓存服务不可用，可能使请求集中到达数据库。

可按触发原因采取措施：

- 在合理范围内给 TTL 增加随机偏移，错开批量预热和刷新时间。
- 配置主从复制与自动故障转移，降低缓存中断的概率。故障切换期间仍可能出现请求失败或数据丢失。
- 给回源请求限流、合并并发请求，并准备降级响应。
- 必要时增加本地缓存，同时处理一致性与失效问题。

## 10 缓存击穿

高并发访问的热点键失效后，大量请求同时重建同一份缓存，给数据库带来瞬时压力。重建越慢，并发回源越容易累积。

常用处理方式：

- 合并同一进程内的并发回源，或使用带超时的互斥锁。获得锁后再次检查缓存，锁值使用唯一标识，释放时原子地校验持有者，避免删除其他请求的锁。
- 使用逻辑过期，由少量请求触发后台刷新，其余请求暂时读取旧值。该方式允许短期陈旧数据，需要安排刷新失败后的重试或降级。

## 11 执行脚本

- 语法：`EVAL script numkeys [key [key ...]] [arg [arg ...]]`
- 功能：执行 Lua 脚本，键参数通过 `KEYS` 读取，普通参数通过 `ARGV` 读取

<<< @/db/codes/redis/eval.sh

脚本在主线程执行，期间其他命令不能与其交错执行，因此应限制脚本耗时。原子执行不等于出错后回滚，运行时错误之前已经完成的写入会保留。脚本访问的键应全部显式传入，Redis Cluster 中这些键必须属于同一个槽，不应在脚本中拼接出未声明的键名。参见[官方脚本说明](https://redis.io/docs/latest/develop/programmability/eval-intro/)。

## 12 消息队列

### 12.1 基于 List

可以用 `LPUSH` 与 `RPOP`，或 `RPUSH` 与 `LPOP` 构造 FIFO 队列。队列为空时，普通弹出命令立即返回空值，`BRPOP`、`BLPOP` 则能阻塞等待。

多个消费者可以竞争弹出同一个 List 的消息，每条消息只会被其中一个弹出。弹出顺序不代表并发消费者的处理完成顺序。弹出后如果消费者崩溃，消息已经离开队列，不能自动重投。

需要处理中消息的恢复机制时，可以用 Redis 6.2 加入的 `LMOVE`、`BLMOVE` 将消息原子地移到处理中列表，处理成功后再移除，并自行实现超时恢复与去重。List 本身没有消费者组和确认协议。持久化与复制可以降低数据丢失风险，但不保证所有已返回成功的消息都能在故障后恢复。

### 12.2 基于 PubSub

Pub/Sub 是发布订阅模型。消费者订阅 channel，生产者向 channel 发布消息，在线的匹配订阅者各收到一份。

- `SUBSCRIBE channel [channel ...]`：订阅频道
- `PUBLISH channel message`：向频道发布消息
- `PSUBSCRIBE pattern [pattern ...]`：按模式订阅频道

Pub/Sub 提供至多一次投递，不保存可重放的消息，也没有确认和重试机制。断线期间的消息无法补读。慢消费者会积累输出缓冲，达到 `client-output-buffer-limit pubsub` 限制时连接可能被关闭，后续消息无法接收。需要历史回放或消费确认时应考虑 Stream。参见[官方 Pub/Sub 说明](https://redis.io/docs/latest/develop/pubsub/)。

### 12.3 基于 Stream

Stream 在 Redis 5.0 加入，以递增 ID 保存字段与值组成的消息条目，支持历史读取、阻塞等待和消费者组。

同一条消息可以被多个独立读取者或多个消费者组读取。可回溯的范围受裁剪、删除和持久化结果限制，故障恢复和重新投递也可能造成消息丢失或重复。

#### 12.3.1 `XADD`

- 语法：`XADD key [NOMKSTREAM] [KEEPREF | DELREF | ACKED] [<MAXLEN | MINID> [= | ~] threshold [LIMIT count]] <* | id> field value [field value ...]`
- 功能：向 Stream 追加条目，返回条目 ID
- 参数：
  - `NOMKSTREAM`：键不存在时不创建 Stream，返回空值
  - `MAXLEN`：按条目数量裁剪。`=` 为精确裁剪，`~` 为近似裁剪，近似结果可能超过阈值
  - `MINID`：按 ID 裁剪，删除小于阈值 ID 的条目
  - `LIMIT`：限制近似裁剪的工作量，仅与 `~` 配合使用
  - `KEEPREF`：默认行为，裁剪条目时保留消费者组 PEL 中的引用
  - `DELREF`：裁剪条目时同时移除所有消费者组中相应的 PEL 引用
  - `ACKED`：仅裁剪已被所有消费者组读取并确认的条目，可能无法达到数量或 ID 阈值
  - `*`：自动生成 ID，格式为 `毫秒时间戳-序号`。显式 ID 必须大于 Stream 已记录的最大 ID，且不能为 `0-0`

`NOMKSTREAM`、`MINID`、`LIMIT` 在 Redis 6.2 加入，`毫秒时间戳-*` 形式在 7.0 加入，`KEEPREF`、`DELREF`、`ACKED` 在 8.2 加入。后者描述的是裁剪时如何处理消费者组引用，与消息正文的字段无关。参见 [`XADD` 官方参考](https://redis.io/docs/latest/commands/xadd/)。

```bash
# 创建 users，并追加一条消息。下面的返回 ID 仅作示意
XADD users * name "zhh" age 18
# "1760774662027-0"
```

#### 12.3.2 `XLEN`

- 语法：`XLEN key`
- 功能：返回 Stream 当前保存的条目数，键不存在时返回 `0`

<<< @/db/codes/redis/xlen.sh

#### 12.3.3 `XREAD`

- 语法：`XREAD [COUNT count] [BLOCK milliseconds] STREAMS key [key ...] id [id ...]`
- 功能：从一个或多个 Stream 读取 ID 大于指定 ID 的条目
- 参数：
  - `COUNT`：每个 Stream 最多返回的条目数
  - `BLOCK`：没有可返回条目时等待，单位毫秒，`0` 表示无限等待，超时返回空值
  - `STREAMS`：先列出所有键，再列出与其一一对应的 ID
  - `0`：读取当前保存的历史条目
  - `$`：以执行命令时的最后一个 ID 为起点，只等待后续条目，不返回当时已有的最后一条

连续读取时，下一次请求应使用上次返回的最后一个 ID。`$` 仅适合首次开始监听新消息，反复使用会跳过两次请求之间产生的条目。Redis 7.4 起还支持 `+`，读取对应 Stream 的最后一条，此时该 Stream 忽略 `COUNT`。参见 [`XREAD` 官方参考](https://redis.io/docs/latest/commands/xread/)。

<<< @/db/codes/redis/xread.sh

### 12.4 Stream 消费者组

消费者组把新条目分配给组内消费者，并维护最后投递 ID（`last-delivered-id`）和待确认列表（PEL）。最后投递不等于处理完成。

默认情况下，投递后条目进入 PEL。业务处理成功后调用 `XACK key group id [id ...]`，从该组的 PEL 中移除引用，Stream 条目本身仍然保留。消费者故障后的 pending 消息需要显式恢复，可以用 `XPENDING` 检查，再用 `XCLAIM` 或 Redis 6.2 加入的 `XAUTOCLAIM` 接管。重新投递可能造成重复处理，业务应按消息 ID 或业务标识去重。参见[官方消费者组说明](https://redis.io/docs/latest/commands/xreadgroup/)。

#### 12.4.1 `XGROUP CREATE`

- 语法：`XGROUP CREATE key group <id | $> [MKSTREAM] [ENTRIESREAD entries-read]`
- 功能：创建消费者组，组已存在时返回错误
- 参数：
  - `id`：初始化最后投递 ID，`0` 表示从已保存的历史条目开始
  - `$`：跳过创建组时已有的条目，只消费后续条目
  - `MKSTREAM`：键不存在时创建空 Stream
  - `ENTRIESREAD`：Redis 7.0 加入，用于初始化消费进度统计，不替代起始 ID

#### 12.4.2 `XGROUP DESTROY`

- 语法：`XGROUP DESTROY key group`
- 功能：删除消费者组及其 PEL，保留 Stream。组存在时返回 `1`，否则返回 `0`

#### 12.4.3 `XREADGROUP`

- 语法：`XREADGROUP GROUP group consumer [COUNT count] [BLOCK milliseconds] [NOACK] STREAMS key [key ...] id [id ...]`
- 功能：以指定消费者身份从组内读取条目，消费者不存在时自动创建
- 参数：
  - `>`：获取尚未投递给该组任何消费者的新条目
  - 数值 ID：读取属于当前消费者、ID 大于该值的 pending 条目，此时 `BLOCK`、`NOACK` 不生效
  - `NOACK`：新条目投递后不进入 PEL，适用于允许丢失、不需要确认的场景

```bash
XADD jobs 1000-0 task "send-email"
# "1000-0"
XGROUP CREATE jobs workers 0
# OK
XREADGROUP GROUP workers worker-1 COUNT 1 STREAMS jobs >
# 返回 1000-0 及其字段，此时该消息在 PEL 中
XACK jobs workers 1000-0
# (integer) 1
```

## 13 GEO

GEO 是 Geolocation 的简写形式，代表地理坐标。Redis 在 3.2 版本中加入了对 GEO 的支持，允许存储地理坐标信息，帮助我们根据经纬度来检索数据。

GEO 底层使用 ZSet，分数编码地理位置。经度范围为 `[-180, 180]`，纬度范围约为 `[-85.05112878, 85.05112878]`，单位均为度。编码会产生量化误差，距离采用球面近似，不能作为精密测绘结果。下面的地点名称和坐标用于演示查询。

### 13.1 `GEOADD`

- 语法：`GEOADD key [NX | XX] [CH] longitude latitude member [longitude latitude member ...]`
- 功能：添加或更新位置，默认返回新增成员数，`CH` 返回新增或位置变化的成员数。`NX`、`XX`、`CH` 在 Redis 6.2 加入

<<< @/db/codes/redis/geoadd.sh

### 13.2 `GEODIST`

- 语法：`GEODIST key member1 member2 [M | KM | FT | MI]`
- 功能：返回两个点之间的距离
- 参数：
  - `M`：以米为单位，默认值
  - `KM`：以千米为单位
  - `FT`：以英尺为单位
  - `MI`：以英里为单位

<<< @/db/codes/redis/geodist.sh{10-15}

### 13.3 `GEOHASH`

- 语法：`GEOHASH key member [member ...]`
- 功能：返回 hash 字符串形式的 `member` 坐标

<<< @/db/codes/redis/geohash.sh{5,6}

### 13.4 `GEOPOS`

- 语法：`GEOPOS key member [member ...]`
- 功能：返回 `member` 的坐标

<<< @/db/codes/redis/geopos.sh{5-7}

### 13.5 `GEORADIUS`

> Redis 6.2 起不推荐用于新代码，改用 `GEOSEARCH` 和 `GEOSEARCHSTORE`。命令仍可使用。

- 常用形式：`GEORADIUS key longitude latitude radius <M | KM | FT | MI>`
- 功能：指定圆心、半径，找到该圆内包含的所有 `member`

### 13.6 `GEOSEARCH`

- 语法：`GEOSEARCH key <FROMMEMBER member | FROMLONLAT longitude latitude> <BYRADIUS radius <M | KM | FT | MI> | BYBOX width height <M | KM | FT | MI>> [ASC | DESC] [COUNT count [ANY]] [WITHCOORD] [WITHDIST] [WITHHASH]`
- 功能：在圆形或矩形范围内搜索成员。默认不保证结果顺序，`ASC`、`DESC` 才按距离排序。`COUNT count ANY` 可以提前结束搜索，不能保证返回最近的 `count` 个成员。Redis 6.2 加入

<<< @/db/codes/redis/geosearch.sh{10-19}

### 13.7 `GEOSEARCHSTORE`

- 语法：`GEOSEARCHSTORE destination source <FROMMEMBER member | FROMLONLAT longitude latitude> <BYRADIUS radius <M | KM | FT | MI> | BYBOX width height <M | KM | FT | MI>> [ASC | DESC] [COUNT count [ANY]] [STOREDIST]`
- 功能：将搜索结果写入目标 ZSet，覆盖目标键。默认保留地理编码分数，使用 `STOREDIST` 时改为保存指定单位下的距离分数。Redis 6.2 加入

## 14 Bitmap

Bitmap 使用 String 的各个二进制位，不是独立的数据类型。Redis 8.2.x 默认上限为 512 MiB，对应位偏移 `0` 到 `2^32 - 1`，限制与 `proto-max-bulk-len` 相关。对很大的偏移执行 `SETBIT` 会分配并填充中间空间，可能增加内存和延迟。

### 14.1 `SETBIT`

- 语法：`SETBIT key offset value`
- 功能：在 `offset` 处存入一个 `1` 或 `0`

<<< @/db/codes/redis/setbit.sh

### 14.2 `GETBIT`

- 语法：`GETBIT key offset`
- 功能：读取 `offset` 处的位。键不存在或偏移超出已存储范围时返回 `0`

<<< @/db/codes/redis/getbit.sh{3-6}

### 14.3 `BITFIELD`

- 常用读取形式：`BITFIELD key GET encoding offset`
- 功能：把连续的位解释为整数，也支持 `SET encoding offset value` 和 `INCRBY encoding offset increment`

`encoding` 使用 `i` 或 `u` 加位宽，例如 `i8`、`u4`。有符号位宽为 1 至 64，无符号位宽为 1 至 63。`OVERFLOW WRAP | SAT | FAIL` 控制后续数值写入的溢出行为，默认 `WRAP`。只读场景还可使用 Redis 6.0 加入的 `BITFIELD_RO`。

<<< @/db/codes/redis/bitfield.sh

### 14.4 `BITCOUNT`

- 语法：`BITCOUNT key [start end [BYTE | BIT]]`
- 功能：统计 Bitmap 中值为 `1` 的 bit 位的数量

范围默认按字节解释，Redis 7.0 起支持 `BIT`，范围包含两端，支持负索引。

<<< @/db/codes/redis/bitcount.sh

### 14.5 `BITPOS`

- 语法：`BITPOS key bit [start [end [BYTE | BIT]]]`
- 功能：查询 bit 数组中指定范围内第一个 0 或 1 出现的位置

返回相对整个字符串的位偏移，找不到时通常返回 `-1`。查找 `0` 且未指定结束位置时，可以把字符串末尾之外视为零填充，因此全为 `1` 的字符串也可能返回末尾之后的位置。`BIT` 单位选项在 Redis 7.0 加入。

<<< @/db/codes/redis/bitpos.sh

## 15 Redis 持久化

### 15.1 RDB

RDB 是 Redis 的二进制快照格式，记录生成快照时的数据集。恢复只能回到已保存的状态，最近一次快照之后的写入可能丢失。

<<< @/db/codes/redis/rdb.sh

正常关闭时是否保存 RDB 取决于持久化配置和 `SHUTDOWN SAVE | NOSAVE` 等参数。不能把正常关闭的行为推广到崩溃、断电或强制终止。

Redis 内部有触发 RDB 的机制，可以在 `redis.conf` 文件中找到，格式如下：

<<< @/db/codes/redis/rdb.conf

`BGSAVE` 使用 `fork` 创建子进程。子进程写快照，主进程继续服务，但创建子进程本身仍可能造成暂停。父子进程最初共享物理页，写时复制（copy-on-write）在内存页被修改时复制相关页，不是每次写入都复制整个数据集。高写入负载下应为这些额外页面预留内存。

### 15.2 AOF

AOF（Append Only File）记录恢复数据集所需的写操作。记录形式可能经过转换，不一定与客户端提交的命令完全相同。

AOF 默认关闭，下面给出开启 AOF 并采用每秒同步策略的配置：

<<< @/db/codes/redis/aof.conf

`appendfsync` 控制文件同步策略，不是控制所有写入何时才进入 AOF 缓冲。`always` 每批写入都同步，`everysec` 通常每秒同步，故障时通常可能丢失约一秒写入，`no` 由操作系统安排同步。实际损失窗口受系统负载、存储和故障类型影响。

Redis 7.0 起采用多文件 AOF，由基础文件、增量文件和 manifest 组织。默认 `aof-use-rdb-preamble yes`，基础文件可以采用 RDB 格式，因此不能直接把整个 AOF 当成一个纯文本文件读取。

`BGREWRITEAOF` 按当前数据集重建基础文件，去掉不再需要的历史操作，期间新增的写入继续记录在增量文件中。文件大小取决于数据、历史写入和基础文件格式。RDB 和 AOF 都启用时，重启优先使用 AOF 恢复。参见[官方持久化说明](https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/)。

Redis 也会在触发阈值时自动去重写 AOF 文件。阈值也可以在 `redis.conf` 中配置：

<<< @/db/codes/redis/auto-aof.conf

## 16 Redis 主从

复制将主节点的数据同步到副本。读写分离可以分担读取负载，但写入仍集中在主节点，副本读也可能读到尚未更新的数据。复制本身不提供自动故障转移。

### 16.1 搭建主从实例

以下是在同一台机器上的本地实验，使用 Redis 8.2.x 的 `redis-server`、`redis-cli`，端口为 `7001`、`7002`、`7003`。实例仅监听本机，数据、配置和日志目录分别独立。实验关闭自动快照和 AOF，以便观察复制，不作为生产持久化配置。

#### 16.1.1 配置与启动

<<< @/db/codes/redis/master-slaver_conf.sh

#### 16.1.2 配置复制关系

使用 `REPLICAOF host port` 把 `7002`、`7003` 配置为 `7001` 的副本。`SLAVEOF` 从 Redis 5.0 起不推荐用于新代码，但仍兼容。

<<< @/db/codes/redis/master-slaver_slaveof.sh

命令修改运行状态，若要在重启后保留，可在配置文件中写入 `replicaof 127.0.0.1 7001`，或在配置文件可写时执行 `CONFIG REWRITE`。`INFO replication` 中部分字段仍沿用 `slave` 命名。

实验结束后逐个关闭实例：

```bash
for port in 7001 7002 7003; do
  redis-cli -p "$port" SHUTDOWN NOSAVE
done
```

### 16.2 数据同步原理

#### 16.2.1 全量同步

副本通过 `PSYNC` 请求同步。没有可用的历史复制状态，或主节点无法提供所需增量时，会全量同步：主节点生成 RDB，副本加载快照，再接收快照生成期间积累的复制流。

复制状态主要由两个值描述：

- replication ID：标识一段复制历史。副本跟随主节点继承 ID，主节点角色变化时可以保留第二个 ID，用于一定范围内的部分同步。
- replication offset：复制流中的字节偏移，不是命令数量，也不是 backlog 当前占用大小。

相同 replication ID 下的相同 offset 对应相同的数据状态，仅比较 ID 或仅比较 offset 都不够。

#### 16.2.2 部分同步

断线重连或重启后能否部分同步，取决于副本是否保留有效复制状态、主节点是否接受该 replication ID，以及缺失的复制流是否仍在 backlog 中。重启后不保证一定走部分同步。

`repl-backlog-size` 控制历史复制流缓冲容量，旧内容会被覆盖。所需区间已经被覆盖时，需要全量同步。增大容量可以容忍更长的断连，但会消耗更多内存。

Redis 8.2.x 默认启用 `repl-diskless-sync yes`，全量同步可以由子进程直接把 RDB 写入网络，减少主节点的临时磁盘写入，仍有 fork、快照生成和网络传输开销。复制规模应结合数据集大小、带宽和恢复时间评估。级联复制可以分担主节点负载，也会增加同步路径和故障依赖。

复制通常是异步的，已返回成功的写入也可能在故障转移后丢失。`WAIT` 可以等待副本确认，Redis 7.2 加入的 `WAITAOF` 可以等待相应 AOF 同步，但都不能把 Redis 复制变成强一致系统。参见[官方复制说明](https://redis.io/docs/latest/operate/oss_and_stack/management/replication/)。

## 17 Redis 哨兵

Sentinel 负责监控主从实例、协调自动故障转移，并向客户端提供当前主节点地址。客户端必须支持 Sentinel 服务发现，并在切换后重新连接。Sentinel 不代理业务请求，也不消除异步复制的数据丢失窗口。

### 17.1 服务状态监控

Sentinel 定期检查实例。某实例超过 `down-after-milliseconds` 没有有效响应时，当前 Sentinel 判定其主观下线（SDOWN）。对于主节点，当至少 `quorum` 个 Sentinel 同意其下线，才判定为客观下线（ODOWN）。副本的主观下线不使用同样的 ODOWN 投票流程。

`quorum` 决定下线判断门槛，执行故障转移还需要 Sentinel 多数派授权。二者不能混为一谈。例如 5 个 Sentinel 配置 `quorum 2` 时，两个可以形成 ODOWN 判断，但仍至少需要三个授权故障转移。一般部署至少 3 个 Sentinel，并分布在独立故障域。

### 17.2 选择新的主节点

执行故障转移的 Sentinel 先获得选举授权，再从合格副本中选择要提升的节点。选择条件包括：

1. 排除不可达、状态不合格、`replica-priority 0` 的副本。
2. 排除与主节点断连过久的副本。Redis 8.2.x 的阈值考虑 `down-after-milliseconds × 10`，还加上主节点已处于 SDOWN 的时长。
3. 按 `replica-priority` 升序选择，数值越小越优先。
4. 优先选择 replication offset 更大的副本。
5. 仍相同时按 run ID 的字典序升序选择，不把它解释为普通数值。

### 17.3 故障转移流程

Sentinel 将选中的副本提升为主节点，再让其他副本跟随新的主节点。原主节点恢复后也被重新配置为副本。

Redis 8.2.x 的 Sentinel 内部仍发送兼容命令 `SLAVEOF NO ONE`、`SLAVEOF host port`，其效果对应 `REPLICAOF`。切换过程中客户端可能暂时失败或连接到旧主节点，需要设置合理的超时、重连和重试策略。参见[官方 Sentinel 文档](https://redis.io/docs/latest/operate/oss_and_stack/management/sentinel/)。

## 18 Redis 分片

### 18.1 Redis Cluster

Redis Cluster 将数据分配给多个主节点，每个主节点可配置副本，节点间使用 cluster bus 协调成员状态和故障转移。分片可以扩展总数据量和写入能力，但具体收益取决于键和请求是否均衡。

客户端请求了错误节点时，服务端通常返回 `MOVED` 或 `ASK` 重定向，由客户端访问目标节点，不会自动代理该业务命令。命令行使用 `redis-cli -c` 跟随重定向，应用使用支持 Redis Cluster 的客户端。Cluster 只支持数据库 `0`。

下面的实验使用 `7101` 至 `7106`，建立三个主节点及各一个副本，与前面的主从实验使用不同端口。所有节点在同一机器上仅用于演示，不能用于验证跨机器高可用。

<<< @/db/codes/redis/cluster-setup.sh

### 18.2 哈希槽

Redis Cluster 有 `16384` 个槽，编号 `0` 至 `16383`。槽分配给主节点，每个键映射到一个槽：

```text
slot = CRC16_XMODEM(key 的有效部分) mod 16384
```

hash tag 的规则是找到键中第一个 `{`，再找其后的第一个 `}`。两者之间非空时，只对其中内容计算，否则对整个键计算。不会在遇到空花括号后继续寻找下一组。

- `{user:101}:name` 与 `{user:101}:email` 使用相同有效部分 `user:101`。
- `foo{bar}{baz}` 使用 `bar`。
- `foo{}{bar}` 按整个键计算。

多键命令、事务和 Lua 脚本通常要求所涉及的键在同一槽。hash tag 可以实现这一点，过度集中也可能造成热点。参见[官方 Cluster 规范](https://redis.io/docs/latest/operate/oss_and_stack/reference/cluster-spec/)。

### 18.3 集群伸缩

新实例须启用 `cluster-enabled yes`，使用独立目录和 `cluster-config-file`，不能直接把前面普通主从实例当作 Cluster 节点加入。

下面在同一个实验 shell 中继续使用 `labdir`，启动 `7107` 并加入集群：

<<< @/db/codes/redis/add-node.sh

`redis-cli --cluster add-node new_host:new_port existing_host:existing_port` 默认添加主节点。添加副本可以使用 `--cluster-slave --cluster-master-id <id>`，这是 `redis-cli` 的参数名。

新主节点最初没有槽，需要用 `redis-cli --cluster reshard host:port` 迁移槽。交互中指定槽数、目标节点 ID 和源节点 ID，检查迁移计划后确认：

<<< @/db/codes/redis/reshard.sh

实验结束后关闭本次启动的节点：

```bash
for port in 7101 7102 7103 7104 7105 7106 7107; do
  redis-cli -p "$port" SHUTDOWN NOSAVE
done
```

### 18.4 故障转移

主节点被足够多的主节点判断为故障后，合格副本可发起选举，获得持有槽的主节点多数派授权后接管槽。是否能够恢复服务取决于多数派是否可达、有无合格副本以及槽覆盖配置。

`cluster-require-full-coverage yes` 是默认配置，槽覆盖不完整时集群停止提供普通数据服务。自动切换不保证零中断，也不保证异步复制的最后一批写入不会丢失。

#### 18.4.1 手动故障转移

- 语法：`CLUSTER FAILOVER [FORCE | TAKEOVER]`
- 执行位置：准备提升的副本节点

普通 `CLUSTER FAILOVER` 先与主节点协调暂停写入，等待副本追平复制偏移，再请求选举授权并接管槽。它不会让主节点宕机，也不是槽的 reshard 操作。

`FORCE` 跳过与主节点的协调，可在主节点不可达时使用，但仍需要多数派授权，可能丢失尚未复制的写入。`TAKEOVER` 进一步跳过正常的多数派授权，自行接管，可能造成配置冲突，应明确理解分区和数据风险后使用。参见 [`CLUSTER FAILOVER` 官方参考](https://redis.io/docs/latest/commands/cluster-failover/)。

### 18.5 Go 访问 Redis Cluster

本文的 go-redis 示例使用 v9.12.1 核对，两个文件各有自己的 `main`，应在独立 Go 模块中分别运行。先执行 `go mod init example/redis-demo` 和 `go get github.com/redis/go-redis/v9@v9.12.1`，再用 `go run 文件名.go` 运行所选示例。

go-redis v9 使用 `ClusterClient` 根据槽分配路由，处理 `MOVED`、`ASK` 等响应。示例连接上面的本地 Cluster，未设置认证，实际部署应按服务端配置传入凭据和 TLS 参数。

<<< @/db/codes/redis/go-redis_shard.go

`Ring` 是另一种客户端：它在多个独立 Redis 实例间按一致性哈希分配键，不使用 Redis Cluster 的槽协议。节点增减时映射会变化，Ring 不负责把已有数据迁移到新位置，不能替代 `ClusterClient`。参见 [go-redis 官方客户端说明](https://github.com/redis/go-redis)。

## 19 最佳实践

### 19.1 键值设计

通常采用 `[业务名]:[对象名]:[ID]`，例如 `login:user:101`，方便辨认业务范围并减少命名冲突。控制键名长度有助于降低内存和传输开销，但不应为缩短键名牺牲必要的含义。

键没有“必须小于 44 字节”或“不能含特殊字符”的通用要求。44 字节是普通字符串对象采用 `embstr` 的一个阈值，8.2.x 中键空间对象还要考虑键名和过期元数据，不能直接套到键名上。涉及 Cluster 多键操作时，单独设计 hash tag。

### 19.2 批处理

小命令频繁往返时，网络延迟可能占主要成本。批处理可以减少等待响应的次数，但复杂命令、脚本或大型集合操作也可能由服务端 CPU 主导，应通过测量判断瓶颈。

#### 19.2.1 多参数命令

可以优先使用 `MSET`、多字段 `HSET`、多成员 `SADD` 等命令。`HMSET` 已不推荐用于新代码。

单条命令的原子性不代表批量越大越好。过大的命令会长时间占用执行线程，并增加网络和内存压力，应根据值大小及延迟目标限制批量。

#### 19.2.2 Pipeline

Pipeline 将多个命令连续写入连接，再读取对应响应，减少逐条等待往返的开销。它不要求命令属于同一种数据类型。

下面的 Go 示例默认连接 `127.0.0.1:6379`，可用 `REDIS_ADDR` 调整地址，在实验实例执行，每次运行会增加 `pipeline-demo:counter` 并设置一小时 TTL。

<<< @/db/codes/redis/go-redis_pipe.go

Pipeline 不保证多个命令整体原子执行。应检查执行错误及各命令结果，连接中断时可能已有部分命令执行。对 `INCR` 等非幂等命令盲目重试可能重复修改数据。参见[官方流水线说明](https://redis.io/docs/latest/develop/using-commands/pipelining/)。

#### 19.2.3 Cluster 中的批处理

| 操作                                      | 槽要求                     | 行为                                     |
| ----------------------------------------- | -------------------------- | ---------------------------------------- |
| 单条 `MSET` 等多键命令                    | 所有键在同一槽             | 不同槽时返回 `CROSSSLOT`                 |
| `ClusterClient` Pipeline 中的多个独立命令 | 各命令分别满足自身的槽要求 | 客户端按目标节点组织请求，可以跨节点发送 |
| 使用相同 hash tag                         | 相关键位于同一槽           | 方便事务和多键操作，但可能造成热点       |

跨节点 Pipeline 不能保证全局顺序或原子性，也不会让一条跨槽 `MSET` 自动合法化。耗时由节点负载、网络延迟和批量共同决定，不能固定简化为一次网络往返。

### 19.3 服务端配置

#### 19.3.1 持久化与部署

持久化配置取决于允许丢失的数据量和恢复时间。可完全重建的缓存可以关闭持久化，但需要评估重启后的预热与回源压力。重要状态应根据需求选择 RDB、AOF 或同时启用，并保留独立备份、验证恢复流程。副本会同步误删和错误写入，不能替代备份。

重写阈值要结合磁盘空间和写入速率设置。`no-appendfsync-on-rewrite yes` 只是暂停后台保存或重写期间的 AOF fsync，写入仍会追加，代价是可能扩大故障丢失窗口，不能理解为“禁止做 AOF”。默认 `no`。

为数据集、复制缓冲、内存碎片、持久化期间的写时复制和操作系统预留容量。单实例大小应结合 fork 延迟、迁移时间和恢复目标确定，没有统一的 4 GiB 或 8 GiB 上限。避免与其他服务争抢关键 CPU、内存和磁盘资源。

#### 19.3.2 慢查询

SLOWLOG 记录命令执行耗时，不包含客户端网络传输等时间，不能用来代表请求的完整延迟。

- `slowlog-log-slower-than`：阈值，单位微秒，Redis 8.2.x 默认 `10000`。`0` 记录所有命令，负值关闭记录。
- `slowlog-max-len`：最多保留的记录数，默认 `128`。

根据观测需求调整阈值和容量，不必一律改为 `1000`。命令参数可能出现在日志中，应按实际数据内容管理访问权限。

- `SLOWLOG LEN`：查看记录数
- `SLOWLOG GET [count]`：读取最近记录，未指定数量时默认 `10` 条
- `SLOWLOG RESET`：清空记录

## 20 数据结构

以下实现以 Redis 8.2.3 为准。内部布局和阈值是实现细节，不能只根据逻辑数据类型推断实际编码，可用 `OBJECT ENCODING key` 检查。

### 20.1 动态字符串 SDS

C 的普通字符串约定以 `\0` 结尾，`strlen` 需要遍历。数组本身可以保存二进制数据，也可以修改，但基于终止符的字符串函数不能把嵌入的 `\0` 当成普通内容，字符串字面量也不能修改。

<<< @/db/codes/redis/string.c

SDS（Simple Dynamic String）通过头部记录长度和容量，再保存字节内容。结尾仍保留 `\0` 以兼容部分 C API，但判断内容长度不依赖它。

<<< @/db/codes/redis/sSDS.c

`sdshdr5` 用 flags 的高 5 位记录短字符串长度，不单独保存容量。Redis 8.2.x 仍使用这种布局，只是代码通常直接访问 flags，不把它解释为该结构体，不能据此认为该编码已弃用。

SDS 长度查询为 `O(1)`，内容二进制安全，扩容时可以预留空间。Redis 8.2.3 的贪心扩容按所需长度计算：小于 1 MiB 时通常扩大到两倍，达到或超过 1 MiB 时通常额外预留 1 MiB。头部、终止符及分配器取整另计，也有不预分配的扩容路径，不能把公式当成所有 SDS 分配的固定大小。参见 [Redis 8.2.3 的 SDS 实现](https://github.com/redis/redis/blob/8.2.3/src/sds.c)。

键名通常以 SDS 布局保存，字符串值也可能使用 SDS。`int` 编码值则直接保存整数，Redis 8.2.x 还使用带键布局的 `kvobj`，不能从一条 `SET` 推断固定创建两个独立 SDS 分配块。

### 20.2 IntSet

IntSet 是 Set 的一种紧凑编码，把整数按升序保存在连续数组中。查找采用二分搜索，时间复杂度为 `O(log N)`，插入和删除可能移动后续元素，最坏为 `O(N)`。

<<< @/db/codes/redis/sIntSet.c

整数宽度有三种，`encoding` 记录每个元素占用的字节数：

<<< @/db/codes/redis/dIntSet.c

加入超出当前宽度范围的整数时会升级。例如从 16 位升级到 32 位，先设置新编码并扩容，再从后向前按旧宽度读取、按新宽度写入，避免覆盖未搬移的数据。触发升级的负数放在最前面，正数放在最后面。删除元素不会自动降低整数宽度。

<<< @/db/codes/redis/fIntsetAdd.c

参见 [Redis 8.2.3 的 IntSet 实现](https://github.com/redis/redis/blob/8.2.3/src/intset.c)。

### 20.3 Dict

Dict 使用哈希桶与冲突链保存映射。Redis 8.2.x 的 `dict` 直接维护两张表的指针、已用元素数和容量指数，旧版独立 `dictht` 的字段布局不再适用。

<<< @/db/codes/redis/sDictht.c

下面是普通的带值 entry。某些只需要键的字典使用更紧凑的无值 entry，不能把它推广为所有节点的固定布局。

<<< @/db/codes/redis/sDictEntry.c

容量为 `2^ht_size_exp`，桶索引由哈希值与容量掩码计算。冲突较多时查找成本会上升，需要扩容或调整负载。

Redis 8.2.3 的一般扩容条件如下，还受字典类型回调和暂停自动调整等条件影响：

- 允许调整容量（`DICT_RESIZE_ENABLE`）时，负载因子达到 `1` 触发扩容。
- 避免调整（`DICT_RESIZE_AVOID`）时，负载因子达到强制阈值 `4` 仍可扩容。
- 禁止调整（`DICT_RESIZE_FORBID`）时，不按这些条件扩容。

后台持久化期间通常避免调整，以降低写时复制开销。阈值 `4` 属于本版实现细节。

<<< @/db/codes/redis/fDictExpandIfNeeded.c

扩容或收缩时建立第二张表，逐步把旧表中的桶迁移过去，再释放旧表。迁移期间读操作需要考虑两张表，新插入元素进入新表。迁移由字典操作和后台维护分批推进，不保证每次操作恰好迁移一个桶，也可能被安全迭代等条件暂停。

这种渐进式 rehash 分摊了搬移成本，但分配新表等操作仍可能造成延迟。参见 [Redis 8.2.3 的字典实现](https://github.com/redis/redis/blob/8.2.3/src/dict.c)。

### 20.4 ZipList 与 Listpack

ZipList 是 Redis 6.2.x 等旧版常见的紧凑容器。Redis 7.0 起，List、Hash、ZSet 的相关编码改用 listpack。8.2.x 中不应再用 ziplist 描述这些类型的现行存储方式。

ZipList 的条目可概括为：

| prevlen            | encoding     | content          |
| ------------------ | ------------ | ---------------- |
| 前一条目的字节长度 | 当前内容编码 | 字符串或整数内容 |

`prevlen` 在前一条目长度小于 254 字节时占 1 字节，否则占 5 字节。前一条目变长可能让后续 `prevlen` 也变长，引发连锁更新。容器使用连续内存，插入、删除可能发生搬移或重新分配，两端操作也可能达到 `O(N)`。

Listpack 同样紧凑存储，条目保存自身的编码和内容，并在末尾记录用于反向遍历的长度信息（backlen）。后续条目不依赖前一条目的长度，避免了 ziplist 的这类连锁更新。它仍可能在修改时移动连续内存，并非所有操作都为常数时间。参见 [Redis 8.2.3 的 Listpack 实现](https://github.com/redis/redis/blob/8.2.3/src/listpack.c)。

### 20.5 QuickList

QuickList 在 Redis 3.2 加入，是由紧凑容器节点组成的双向链表。早期节点使用 ziplist，Redis 7.0 起使用 listpack。Redis 8.2.x 的节点可以保存打包的 listpack，也可以用 plain 节点保存单个大元素。

`list-max-listpack-size` 控制打包节点的大小。正值按元素数量限制，负值对应以下字节大小目标：

| 配置值 | 大小目标      |
| ------ | ------------- |
| `-1`   | 4 KiB         |
| `-2`   | 8 KiB，默认值 |
| `-3`   | 16 KiB        |
| `-4`   | 32 KiB        |
| `-5`   | 64 KiB        |

这不是整个 List 的容量上限，也不表示单个大元素绝不能超过该大小。旧配置名 `list-max-ziplist-size` 在 8.2.x 仍作为兼容别名，正文使用现行名称。

<<< @/db/codes/redis/list-max-listpack-size.sh

`list-compress-depth` 控制首尾各保留多少个不压缩节点，`0` 表示关闭压缩，`1` 表示首尾各保留一个，依此类推。中间符合条件的打包节点使用 LZF，过小或压缩收益不足等节点不一定被压缩。

<<< @/db/codes/redis/list-compress-depth.sh

<<< @/db/codes/redis/sQuickList.c

<<< @/db/codes/redis/sQuickListNode.c

参见 [Redis 8.2.3 的 QuickList 定义](https://github.com/redis/redis/blob/8.2.3/src/quicklist.h)和[配置文件](https://github.com/redis/redis/blob/8.2.3/redis.conf)。

### 20.6 SkipList

跳表用随机层高建立多级前向指针。Redis 的 ZSet 跳表先按分数排序，同分数时按成员字典序排序，还用 backward 支持反向遍历，用 span 支持排名计算。

查找、插入和删除的期望时间复杂度为 `O(log N)`，不是严格的最坏情况保证。范围查询还与返回的元素数有关。

<<< @/db/codes/redis/sZSkipList.c

<<< @/db/codes/redis/sZSkipListNode.c

参见 [Redis 8.2.3 的跳表实现](https://github.com/redis/redis/blob/8.2.3/src/t_zset.c)。

### 20.7 RedisObject

`redisObject`（`robj`）记录值的逻辑类型、内部编码、引用计数和淘汰元信息。不能把每个键都描述为单独封装在普通 robj 中。Redis 8.2.x 使用 `kvobj` 布局组织键和值，并以 `iskvobj`、`expirable` 等标志表示相关属性。

<<< @/db/codes/redis/sRedisObject.c

编码常量保留了旧格式的编号。下面列出正文涉及的现行编码：

<<< @/db/codes/redis/sRedisObjEncoding.c

| 逻辑类型 | Redis 8.2.x 的编码                       |
| -------- | ---------------------------------------- |
| String   | `int`、`embstr`、`raw`                   |
| List     | `listpack`、`quicklist`                  |
| Set      | `intset`、`listpack`、`hashtable`        |
| ZSet     | `listpack`、`skiplist`，后者同时维护字典 |
| Hash     | `listpack`、`listpackex`、`hashtable`    |
| Stream   | `stream`，基于 radix tree 与 listpack    |

Bitmap 使用 String，GEO 使用 ZSet，没有独立对象编码。参见 [Redis 8.2.3 的对象定义](https://github.com/redis/redis/blob/8.2.3/src/server.h)和[对象操作](https://github.com/redis/redis/blob/8.2.3/src/object.c)。

### 20.8 常见类型的编码选择

#### 20.8.1 String

Redis 8.2.x 常见选择如下，实际编码还受写入命令和对象转换路径影响：

- 能解析为本机 `long` 范围内整数的值，可用 `int` 编码，把整数存入 `ptr` 字段，而不是指向一份整数 SDS。64 位环境下通常对应有符号 64 位范围。
- 普通独立字符串对象的 `embstr` 长度阈值为 44 字节，对象头和 SDS 一次分配。
- 写入键空间时，Redis 8.2.x 会创建 `kvobj`。是否继续嵌入字符串，还取决于键名和过期元数据能否满足 `kvobjSet` 的 64 字节大小判断。常见 64 位构建下，短键且没有附加过期时间时，判断式约为 `键长 + 值长 <= 41`，附加过期时间时还需减去 8 字节。因此不能仅按值长预测 `OBJECT ENCODING`。
- 较长字符串或需要转换的值使用 `raw`，对象头与 SDS 分开分配。

例如单字节键 `a` 写入 40 字节值时可保持 `embstr`，41 字节值则为 `raw`。这与普通对象的 44 字节阈值并不矛盾，前者还包含了键的布局。`APPEND` 等修改也可能将 `embstr` 转为 `raw`。长度均按字节计算。

<<< @/db/codes/redis/obj_encoding_string.sh

#### 20.8.2 List

Redis 3.2 至 6.2.x 的 List 使用带 ziplist 节点的 quicklist，Redis 7.0 改用带 listpack 节点的 quicklist。Redis 7.4 起，小型 List 可以直接使用单个 listpack，较大时转换为 quicklist。Redis 8.2.x 延续这两种编码，也会在满足条件时把缩小后的 quicklist 转回 listpack。

编码选择与 `list-max-listpack-size` 及元素大小有关。

#### 20.8.3 Set

全部成员是合适的整数且数量不超过 `set-max-intset-entries` 时，通常使用 intset，默认阈值 `512`。Redis 7.2 起还支持 listpack，8.2.x 默认限制为 `set-max-listpack-entries 128` 和 `set-max-listpack-value 64`，分别约束成员数和成员字节长度。

超过相应条件后可转换为 hashtable。字典保存成员，不需要一般映射中的值，但实现可以使用紧凑无值节点，不能按完整 `key + null value` 估算所有内存开销。编码转换路径也受已有编码影响，不是每次修改都重新选择最小容器。

#### 20.8.4 ZSet

小型 ZSet 使用 listpack，成员和分数相邻保存。Redis 8.2.x 默认限制为 `zset-max-listpack-entries 128` 和 `zset-max-listpack-value 64`，后者约束成员字节长度。

较大时转换为 skiplist 编码，其 `zset` 同时维护字典和跳表：字典按成员查询分数，跳表按分数或排名查询。不是 `hashtable` 与 `skiplist` 两种独立的 ZSet 对象编码。

<<< @/db/codes/redis/sZset.c

下面的函数创建 skiplist 编码对象，不代表所有 ZSet 初始都采用该编码：

<<< @/db/codes/redis/fCreateZsetObject.c

#### 20.8.5 Hash

普通小型 Hash 使用 listpack，字段与值相邻保存。默认 `hash-max-listpack-entries 512` 限制字段与值的对数，`hash-max-listpack-value 64` 限制单个字段或值的字节长度，超过条件时转换为 hashtable。

<<< @/db/codes/redis/hash-max-listpack.sh

Redis 7.4 加入字段级过期。Redis 8.2.x 中，为 listpack Hash 的字段设置过期时间时，会使用带过期元数据的 `listpackex`，hashtable Hash 也支持字段过期。字段 TTL 与整个键的 TTL 是不同层次，不能混用。参见 [`HEXPIRE` 官方参考](https://redis.io/docs/latest/commands/hexpire/)和 [Redis 8.2.3 的 Hash 实现](https://github.com/redis/redis/blob/8.2.3/src/t_hash.c)。

## 21 网络模型

### 21.1 用户态与内核态

用户态和内核态区分执行权限。应用通常通过系统调用请求内核管理的网络、文件和内存资源，不能把地址空间简单理解为固定切成两半。x86 常用 Ring 3、Ring 0 表示相应特权级，其他体系结构的机制不同。

普通 socket 读写通常涉及用户缓冲区与内核缓冲区之间的数据复制。内核报告“可读”或“可写”，表示此时相应操作可以推进，不表示一个完整应用消息已经到达，也不表示所有待发送数据都能一次写完。

下面区分阻塞、非阻塞、多路复用、信号驱动和异步 IO。Redis 的客户端连接主要使用非阻塞 socket 与事件循环，多路复用负责发现就绪连接。

### 21.2 阻塞 IO

阻塞 socket 上的 `recv` 可以等待数据到达后才返回。若事件循环直接在一个尚无数据的连接上阻塞，就不能及时处理其他连接。连接关闭和错误也会让调用返回，不是只在读到数据时返回。

### 21.3 非阻塞 IO

非阻塞 socket 上，暂时没有可读取数据时，`recv` 返回 `-1` 并设置 `EAGAIN` 或 `EWOULDBLOCK`。应用可以配合多路复用等待下一次就绪通知，不必持续轮询。非阻塞只改变等待行为，普通 `recv` 的数据复制仍在调用中完成。

### 21.4 IO 多路复用

文件描述符（FD）标识进程打开的资源，socket 也使用 FD。多路复用让一个线程等待多个 FD 的就绪事件，再对相应连接执行读写。

| 接口     | 描述符限制                                                    | 传入方式             | 返回后处理                 |
| -------- | ------------------------------------------------------------- | -------------------- | -------------------------- |
| `select` | Linux/glibc 常见 `FD_SETSIZE` 为 1024，描述符数值必须小于该值 | 每次传入 fd_set      | 检查返回集合中的置位描述符 |
| `poll`   | 没有 fd_set 的固定上限，仍受资源限制                          | 每次传入 pollfd 数组 | 遍历数组检查 revents       |
| `epoll`  | 没有 fd_set 的固定上限，仍受描述符、内存及监听限制            | 维护已注册的监听集合 | 遍历返回的就绪事件         |

#### select

`select` 会修改传入集合，下次调用通常需要重新准备。固定大小的 fd_set 和扫描成本限制了其扩展性。1024 是常见库实现的限制，不是所有平台的统一连接数上限。参见 [`select(2)`](https://man7.org/linux/man-pages/man2/select.2.html)。

#### poll

`poll` 使用 `pollfd` 数组，每个元素指定 FD 与关注事件，返回时在 `revents` 中给出状态。调用和返回处理都涉及整个数组，连接数量增加时扫描开销也随之增加。

<<< @/db/codes/redis/poll.c

返回值是有事件的数组元素数，`0` 表示超时，`-1` 表示错误。应用还要检查挂断和错误等事件，不应只处理 `POLLIN`。参见 [`poll(2)`](https://man7.org/linux/man-pages/man2/poll.2.html)。

#### epoll

`epoll_create1` 创建实例，`epoll_ctl` 增删改监听，`epoll_wait` 等待并返回就绪事件。监听集合保存在内核，不需要每次传入完整集合。事件交付仍有内核与用户态之间的复制和处理成本。

<<< @/db/codes/redis/epoll.c

#### LT 与 ET

- LT（Level Triggered）：默认模式，只要相应就绪条件持续满足，后续等待仍可以报告该事件。例如缓冲区有未读数据时可以继续报告可读。
- ET（Edge Triggered）：关注状态变化，不能依赖未处理完的就绪状态被反复报告。通常使用非阻塞 socket，循环读写直到 `EAGAIN`，并正确处理短读、短写及断连。

ET 不保证总是比 LT 快，也不等同于某连接只通知一次。Redis 8.2.x 的 epoll 后端未设置 `EPOLLET`，使用 LT。参见 [`epoll(7)`](https://man7.org/linux/man-pages/man7/epoll.7.html)和 [Redis 8.2.3 的 epoll 后端](https://github.com/redis/redis/blob/8.2.3/src/ae_epoll.c)。

### 21.5 信号驱动 IO

通过异步通知配置，可以让就绪事件产生 `SIGIO` 等信号，应用再进行读写。普通信号可能合并，信号处理函数也受异步信号安全约束，实现复杂度不能简单归因于“信号种类不够”。参见 [`signal(7)`](https://man7.org/linux/man-pages/man7/signal.7.html)。

### 21.6 异步 IO

异步 IO 先提交操作，再通过完成通知或结果队列获取结果。就绪通知表示“现在可以尝试读写”，完成通知表示“提交的操作已经完成”，这是它与多路复用的区别。

Linux 的 POSIX AIO、原生 AIO 和 io_uring 是不同接口，能力与实现也不同。glibc 的 POSIX AIO 通常使用用户态线程实现。Redis 8.2.x 的常规客户端事件循环在 Linux 上使用 epoll，并非改用 io_uring。参见 [`aio(7)`](https://man7.org/linux/man-pages/man7/aio.7.html)。

### 21.7 Redis 的线程与进程

Redis 8.2.x 的常规数据命令主要由主线程串行执行，但网络读写、协议解析和后台工作可以分配给其他线程。不同客户端之间的执行顺序不能简单等同于请求在网络中的到达顺序。

需要区分：

- IO 线程处理连接读写和部分协议解析，常规数据结构操作仍在主线程。
- 后台线程承担部分 AOF fsync、关闭文件和延迟释放对象等工作。
- `BGSAVE`、`BGREWRITEAOF` 使用子进程，不是后台线程。

Redis 6.0 引入可配置的多线程 IO，旧版可用 `io-threads-do-reads yes` 开启读取处理。Redis 8.0 起重构了 IO 线程机制，8.2.x 启用 `io-threads N` 后同时处理读、写及协议解析，`io-threads-do-reads` 已不再生效。默认 `io-threads 1`，相当于不启用额外 IO 线程。

是否增加线程应依据 CPU 与网络测量。大型集合操作、长 Lua 脚本等可能占用主线程较长时间，此时 CPU 也会成为瓶颈。参见 [Redis 8.2.3 的 IO 线程实现](https://github.com/redis/redis/blob/8.2.3/src/iothread.c)和[配置说明](https://github.com/redis/redis/blob/8.2.3/redis.conf)。

## 22 Redis 通信协议

### 22.1 RESP

RESP（REdis Serialization Protocol）是 Redis 客户端与服务端之间的序列化协议，常见传输是 TCP，也支持 Unix socket。协议头使用 ASCII 控制字符，bulk string 按字节长度传输，可以包含任意二进制数据，不能简单称为纯文本协议。

Redis 2.0 起使用 RESP2，Redis 6.0 加入 RESP3 和 `HELLO` 协商。Redis 8.2.x 的新连接初始使用 RESP2，客户端可以发送 `HELLO 3` 切换，部分驱动会主动协商 RESP3，因此应用实际使用的版本取决于驱动配置。

### 22.2 RESP2 数据类型

RESP2 用首字节区分五类数据，头部以 CRLF（`\r\n`）终止：

| 首字节 | 类型             | 示例                           |
| ------ | ---------------- | ------------------------------ |
| `+`    | 简单字符串       | `+OK\r\n`                      |
| `-`    | 错误             | `-ERR message\r\n`             |
| `:`    | 有符号 64 位整数 | `:1024\r\n`                    |
| `$`    | bulk string      | `$5\r\nhello\r\n`              |
| `*`    | 数组             | `*2\r\n$1\r\na\r\n$1\r\nb\r\n` |

`$0\r\n\r\n` 是空字符串，`$-1\r\n` 是空 bulk string，`*-1\r\n` 是空数组，与 `*0\r\n` 的零元素数组不同。空值不只表示键不存在，也可能表示阻塞读取超时等结果。数组元素可以嵌套，也可以是错误。

Redis 请求通常是 bulk string 数组。`SET name zhh` 可以表示为下面的字节序列，展示中的转义符表示实际 CR、LF 字节，长度按字节计算：

```text
*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$3\r\nzhh\r\n
```

`redis-cli --no-raw` 输出的是供人阅读的格式，不是线上的 RESP 原文。需要观察原始响应时可用 TCP 工具或抓包，例如本地未启用 TLS 的实例：

```bash
printf '*1\r\n$4\r\nPING\r\n' | nc -N 127.0.0.1 6379
# 响应字节为 +PONG\r\n。-N 是 OpenBSD netcat 的选项
```

RESP3 增加映射、集合、布尔值、浮点数和 push 等类型。解析器必须按已协商的版本处理响应。参见[官方 RESP 规范](https://redis.io/docs/latest/develop/reference/protocol-spec/)。

### 22.3 用 Go 模拟客户端

下面的教学示例只实现 RESP2，处理简单字符串、错误、整数、bulk string 和递归数组。它区分空值与空字符串，按长度读取二进制内容，并在数组中保留错误元素，避免未消费完响应而破坏后续读取。

运行 `go run docs/db/codes/redis/resp.go`，默认连接本地 `127.0.0.1:6379`，也可通过 `REDIS_ADDR` 指定地址。示例使用 `resp-demo:*` 键，需在实验实例执行，未实现认证、TLS、连接池或重试。

<<< @/db/codes/redis/resp.go

## 23 Redis 内存回收

64 位 Redis 默认 `maxmemory 0`，表示不设置数据内存上限。可以配置：

```conf
maxmemory 4gb
```

这个值只是示例，应按可用内存与额外开销选择。`maxmemory` 不是进程 RSS 的硬上限，复制和 AOF 缓冲、内存碎片、写时复制等仍会占用内存，不能直接设置成机器全部内存。

过期删除与内存淘汰是不同机制。过期删除根据到期时间运行，并不要求先达到 `maxmemory`。内存淘汰则在超出相应内存限额时按 `maxmemory-policy` 处理，可能删除尚未过期甚至没有 TTL 的键。

### 23.1 过期删除

#### DB 与到期时间

Redis 8.2.x 的 `redisDb` 通过 `kvstore *keys` 管理键空间，通过 `kvstore *expires` 索引带键级过期时间的对象。到期时间以绝对毫秒时间戳附在 `kvobj` 的相应布局中，不是单独保存一个剩余 TTL。`hexpires` 用于组织 Hash 字段过期的维护。

<<< @/db/codes/redis/sRedisDb.c

这些索引、对象及元数据都有内存开销，不能断言 TTL 只额外占用一个 8 字节整数。普通实例默认有 16 个数据库，可通过 `databases` 配置，Cluster 只使用数据库 `0`。参见 [Redis 8.2.3 的键空间实现](https://github.com/redis/redis/blob/8.2.3/src/db.c)。

#### 惰性删除

访问键时检查到期时间，已经过期的键按不存在处理，再按实例角色和删除策略完成相应清理。它避免给每个键单独设置定时器，但不访问的过期键仍需要主动清理，否则可能继续占用内存。

#### 主动删除

Redis 主动扫描过期索引，用时间预算限制单次工作量。Redis 8.2.3 默认 `active-expire-effort 1` 时：

- SLOW 周期由 `serverCron` 触发，基础时间预算是调度周期的 25%。当实际 `hz` 为 10 时，约为 25 ms。动态 hz 及负载会影响调度，不能认为每轮固定占用这些时间。
- 每轮数据库扫描以约 20 个带 TTL 的键为目标，也限制扫描的桶数，不是每个桶抽 20 个键。
- 估计过期占比超过 10% 时可继续扫描，直到扫描完成或时间预算耗尽。
- FAST 周期在事件循环的 `beforeSleep` 中按条件触发，基础预算 1 ms，间隔至少为该预算的两倍，并非每轮都执行。

`active-expire-effort` 可调整扫描力度和预算。Hash 字段还有独立的过期维护，不能把所有过期行为都套入键级扫描步骤。参见 [Redis 8.2.3 的过期实现](https://github.com/redis/redis/blob/8.2.3/src/expire.c)。

### 23.2 内存淘汰

判断淘汰时，Redis 会从分配器统计的内存中扣除部分不计入淘汰的复制、AOF 缓冲等开销，`INFO memory` 的 `mem_not_counted_for_evict` 可用于观察，因此不是直接按 RSS 或完整 `used_memory` 比较上限。

下面是 `processCommand` 中的相关片段，省略其他命令检查：

<<< @/db/codes/redis/fProcessCommand.c

Redis 8.2.x 有以下八种策略。候选不足、无法释放足够内存时，仍可能拒绝受内存限制的命令。

| 策略              | 候选范围    | 选择依据                                                    |
| ----------------- | ----------- | ----------------------------------------------------------- |
| `noeviction`      | 不主动淘汰  | 默认策略，超过限额后拒绝可能增加内存且受 OOM 检查约束的命令 |
| `volatile-lru`    | 有 TTL 的键 | 近似 LRU                                                    |
| `allkeys-lru`     | 所有键      | 近似 LRU                                                    |
| `volatile-lfu`    | 有 TTL 的键 | 近似 LFU                                                    |
| `allkeys-lfu`     | 所有键      | 近似 LFU                                                    |
| `volatile-random` | 有 TTL 的键 | 随机                                                        |
| `allkeys-random`  | 所有键      | 随机                                                        |
| `volatile-ttl`    | 有 TTL 的键 | 优先选择采样中更早到期的键                                  |

没有带 TTL 的键时，`volatile-*` 没有可淘汰候选，行为接近 `noeviction`。读取及部分删除等减少内存的操作通常仍能执行。参见[官方淘汰说明](https://redis.io/docs/latest/develop/reference/eviction/)。

#### LRU 与 LFU

LRU（Least Recently Used）考虑最近访问时间，LFU（Least Frequently Used）考虑近期访问频率。Redis 使用采样和候选池近似选择，并不维护完整的全局 LRU 链表或精确频次排序。`maxmemory-samples` 默认 `5`，影响采样精度与开销。

对象头中的 `lru` 字段占 24 bit：LRU 模式保存秒级时钟信息，LFU 模式把高 16 bit 用于分钟级时间，低 8 bit 用于概率计数。

<<< @/db/codes/redis/tRobj.c

Redis 8.2.3 的 LFU 计数初值为 `5`，增加概率为：

```text
base = max(counter - 5, 0)
P = 1 / (base × lfu-log-factor + 1)
```

`lfu-log-factor` 默认 `10`。访问时以概率 `P` 增加计数，计数最大为 `255`，不是每访问一次就加一。

`lfu-decay-time` 默认 `1` 分钟。访问或评估淘汰候选时，根据经过的衰减周期减少计数，最低为 `0`，不是后台定时器逐分钟遍历所有对象减一。配置为 `0` 时不衰减。LFU 只是对近期频率的近似估计，不是累计访问次数，也不保证选出严格的全局最低频键。参见 [Redis 8.2.3 的淘汰与 LFU 实现](https://github.com/redis/redis/blob/8.2.3/src/evict.c)。
