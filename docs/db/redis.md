---
outline: [2, 4]
---

# Redis

Redis 诞生于 2009 年，全称是 <span style="color:red;">Re</span>mote <span style="color:red;">Di</span>ctionary <span style="color:red;">S</span>erver。

## 1 通用命令

### 1.1 `KEYS`

- 语法：`KEYS pattern`
- 功能：查看符合模版的所有 `key`，不建议在生产环境使用
- 时间复杂度：`O(N)`，其中 `N` 是数据库中的键数

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
- 功能：判断指定的 `key` 是否存在

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

- 语法：`EXPIRE key seconds`
- 功能：设置 `key` 的过期时长

### 1.5 `TTL`

- 语法：`TTL key`
- 功能：查看指定 `key` 的剩余有效时长（秒）

如果 `key` 存在但是没有设置过期时长，返回 `-1`。如果 `key` 不存在返回 `-2`。

## 2 String 类型

String 类型，也就是字符串类型，是 Redis 中最简单的存储类型。

其 value 是字符串，不过根据字符串的格式不同，又可以分为 3 类：

- string：普通字符串
- int：整数类型，可以进行自增自减
- float：浮点类型，可以进行自增自减

不管哪种格式，底层都是字节数组形式存储，只不过是编码方式不同。字符串类型的最大空间不能超过 512m

### 2.1 `SET`

- 语法：`SET key value [NX | XX] [EX seconds | KEEPTTL]`
- 功能：添加或修改一个 String 类型的键值对
- 可选项：
  - `NX`：只有 `key` 不存在时才设置
  - `XX`：只有 `key` 已存在时才设置
  - `EX`：设置过期时长
  - `KEEPTTL`：保留 `key` 原有的过期时长

### 2.2 `GET`

- 语法：`GET key`
- 功能：根据 `key` 获取 String 类型的 `value`

### 2.3 `MSET`

- 语法：`MSET key value [key value ...]`
- 功能：批量添加多个 String 类型的键值对

### 2.4 `MGET`

- 语法：`MGET key [key ...]`
- 功能：批量获取多个 `key` 的 `value`

### 2.5 `INCR`

- 语法：`INCR key`
- 功能：整型自增 `1`

### 2.6 `INCRBY`

- 语法：`INCRBY key increment`
- 功能：整型自增 `increment`

### 2.7 `INCRBYFLOAT`

- 语法：`INCRBYFLOAT key increment`
- 功能：浮点型自增 `increment`

### 2.8 `SETNX`

> 弃用，推荐采用 `SET key value NX`

- 语法：`SETNX key value`
- 功能：如果 `key` 不存在才新增。

`SETNX` 是 **SET** if **N**ot e**X**ists 的简写。

### 2.9 `SETEX`

> 弃用，推荐采用 `SET key value EX seconds`

- 语法：`SETEX key seconds value`
- 功能：新增 `key` 并设置有效时长

## 3 Key 的层级格式

Redis 的 `key` 允许有多个单词形成层级结构，多个单词间用 `:` 隔开。

## 4 Hash 类型

Hash 类型，也叫散列，其 `value` 是一个无需字典，类似于 Java 中的 HashMap 结构。

String 结构将对象的所有字段保存为一整个字符串，如果要修改其中某个字段很不方便。

Hash 结构可以将对象中的每次字段独立存储，可以针对单个字段做 CRUD。

### 4.1 `HSET`

- 语法：`HSET key field value [field value ...]`
- 功能：添加或修改 hash 类型 `key` 的 `filed` 的值

### 4.2 `HGET`

- 语法：`HGET key field`
- 功能：获取一个 hash 类型 `key` 的 `filed` 的值

### 4.3 `HMSET`

> 弃用，采用 `HSET` 效果一样

- 语法：`HMSET key field value [field value ...]`
- 功能：批量添加多个 hash 类型 `key` 的 `filed` 的值

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
- 功能：批量获取多个 hash 类型 `key` 的 `filed` 的值

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
- 功能：获取一个 hash 类型的 `key` 中的所有 `filed` 和对应 `value`

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
- 功能：获取一个 hash 类型的 `key` 中的所有的 `filed`

```bash
HSET myhash field1 "Hello" field2 "World"
# (integer) 2
HKEYS myhash
# 1) "field1"
# 2) "field2"
```

### 4.7 `HVALS`

- 语法：`HVALS key`
- 功能：获取一个 hash 类型的 `key` 中的所有的 `value`

```bash
HSET myhash field1 "Hello" field2 "World"
# (integer) 2
HVALS myhash
# 1) "Hello"
# 2) "World"
```

### 4.8 `HINCRBY`

- 语法：`HINCRBY key field increment`
- 功能：让一个 hash 类型的 `key` 的字段值增加指定步长

```bash
HSET myhash field 5
# (integer) 1
HINCRBY myhash field 1
# (integer) 6
HINCRBY myhash field -10
# (integer) -4
```

### 4.9 `HINCRBY`

- 语法：`HINCRBY key field increment`
- 功能：让一个 hash 类型的 `key` 的字段值增加指定步长

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

Redis 的 List 类型与 Java 的 LinkedList 类似，可以看作双向链表。既可以正向检索也可以反向检索。

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

Redis 的 Set 结构与 Java 中的 HashSet 类似，可以看作是一个 `value` 为 `null` 的 HashMap。因为也是一个 hash 表，因此具备与 HashSet 类似的特征：

- 无序
- 元素不重复
- 查找快
- 支持交集、并集、差集等功能

### 6.1 `SADD`

- 语法：`SADD key member [member ...]`
- 功能：往集合中添加元素

```bash
SADD myset "Hello" "World"
# (integer) 2
SADD myset "World"
# (integer) 0
```

### 6.2 `SREM`

- 语法：`SREM key member [member ...]`
- 功能：移除集合中的指定元素

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
- 功能：获取集合中的全部元素

```bash
SADD myset Hello World
# (integer) 2
SMEMBERS myset
# 1) "Hello"
# 2) "World"
```

### 6.6 `SINTER`

- 语法：`SINTER key [key ...]`
- 功能：求交集（intersection）

```bash
SADD s1 a b c d
# (integer) 4
SADD s2 c
# (integer) 1
SADD s2 a c e
# (integer) 3
SINTER s1 s2 s3
# 1) "c"
```

### 6.7 `SDIFF`

- 语法：`SDIFF key [key ...]`
- 功能：求差集（difference set）

```bash
SADD s1 a b c d
# (integer) 4
SADD s2 c
# (integer) 1
SADD s2 a c e
# (integer) 3
SDIFF s1 s2 s3
# 1) "d"
# 2) "b"
```

### 6.8 `SUNION`

- 语法：`SUNION key [key ...]`
- 功能：求并集（union）

```bash
SADD s1 a b c d
# (integer) 4
SADD s2 c
# (integer) 1
SADD s2 a c e
# (integer) 3
SUNION s1 s2 s3
# 1) "c"
# 2) "e"
# 3) "b"
# 4) "d"
# 5) "a"
```

## 7 SortedSet 类型

Redis 的 SortedSet 是一个可排序的集合，与 Java 中的 TreeSet 有些类似，但底层数据结构差别很大。SortedSet 中每个元素都带有一个 score 属性，可以基于 score 属性对元素排序，底层实现是一个跳表（SkipList）加 hash 表。

SortedSet 具备下列特性：

- 可排序
- 元素不重复
- 查询速度快

因为 SortedSet 的可排序特性，经常被用来实现排行榜这样的功能。

### 7.1 `ZADD`

- 语法：`ZADD key [NX | XX] score member [score member ...]`
- 功能：添加元素到有序集合，如果已存在则更新 `score`
- 可选项：
  - `XX`：仅更新已存在的元素。不添加新元素。
  - `NX`：只添加新元素。不更新现有元素。

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

- 语法：`ZRANK key member`
- 功能：获取指定元素在有序集合中的排名

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

- 语法：`ZRANGE key start stop [BYSCORE | BYLEX] [REV] [WITHSCORES]`
- 功能：按照 `score` 升序排序后，获取指定排名范围内的元素，排名从 0 开始
- 参数：
  - `REV`：是否降序顺序
  - `WITHSCORES`：是否返回 `score`

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

> 弃用，推荐 `ZRANGE`

- 语法：`ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]`
- 功能：按照 `score` 排序后，获取指定 `score` 范围内的元素

### 7.10 `ZDIFF`

- 语法：`ZDIFF numkeys key [key ...] [WITHSCORES]`
- 功能：求差集（difference set）

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

### 7.11 `ZDIFF`

- 语法：`ZDIFF numkeys key [key ...] [WITHSCORES]`
- 功能：求差集（difference set）

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

### 7.12 `ZINTER`

- 语法：`ZINTER numkeys key [key ...] [WITHSCORES]`
- 功能：求交集（intersection）

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

### 7.13 `ZUNION`

- 语法：`ZUNION numkeys key [key ...] [WITHSCORES]`
- 功能：求并集（union）

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

缓存穿透是指客户端请求的数据在缓存中和数据库中都不存在，这样缓存永远不会生效，这些请求都会打到数据库。

常见的解决方案有两种：

- 缓存空对象
- 布隆过滤

## 9 缓存雪崩

缓存雪崩是指同一时段大量的缓存 key 同时失效或者 Redis 服务宕机，导致大量请求到达数据库，带来巨大压力。

解决方案：

- 给不同的 key 的 TTL 添加随机值
- 利用 Redis 集群提高服务的可用性
- 给缓存业务添加降级限流策略
- 给业务添加多级缓存

## 10 缓存击穿

缓存击穿问题也叫热点 Key 问题，就是一个被高并发访问并且缓存重建业务比较复杂的 key 突然失效了，无数的请求访问会在瞬间给数据库带来巨大的冲击。

解决方案：

- 互斥锁
- 逻辑过期

## 11 执行脚本

`EVAL` 命令：

- 语法：`EVAL script numkeys [key [key ...]] [arg [arg ...]]]`
- 功能：执行 Lua 脚本

<<< @/db/codes/redis/eval.sh

## 12 消息队列

### 12.1 基于 List

消息队列（Message Queue），字面意思就是存放消息的队列。而 Redis 的 List 数据结构是一个双向链表，很容易模拟出队列效果。

队列是入口和出口不在同一边，因此我们可以利用：`LPUSH` 结合 `RPOP`、或者 `RPUSH` 结合 `LPOP` 来实现。

不过要注意的是，当队列中没有消息时，`RPOP` 或 `LPOP` 操作会返回 `nil`，并不像 JVM 的阻塞队列那样会阻塞并等待消息。因此这里应该使用 `BRPOP` 或者 `BLPOP` 来实现阻塞效果。

优点：

- 利用 Redis 存储，不受限于 JVM 内存上限
- 基于 Redis 的持久化机制，数据安全性有保证
- 可以满足消息的有序性

缺点：

- 无法避免消息丢失
- 只支持单消费者

### 12.2 基于 PubSub

PubSub（发布订阅）是 Redis 2.0 版本引入的消息传递模型。顾名思义，消费者可以订阅一个或多个 channel，生产者向对应 channel 发送消息后，所有订阅者都能收到相关消息。

- `SUBSCRIBE channel [channel ...]`：订阅频道
- `PUBLISH channel message`：向一个频道发送消息
- `PSUBSCRIBE pattern [pattern ...]`：订阅与 `pattern` 格式匹配的所有频道

优点：

- 采用发布订阅模型，支持多生产、多消费

缺点：

- 不支持数据持久化
- 无法避免消息丢失
- 消息堆积有上限，超出时数据丢失

### 12.3 基于 Stream

Stream 是 Redis 5.0 引入的一种新数据类型，可以实现一个功能非常完善的消息队列。

特点：

- 消息可回溯
- 一个消息可以被多个消费者读取
- 可以阻塞读取

#### 12.3.1 `XADD`

- 语法：`XADD key [NOMKSTREAM] [KEEPREF | DELREF | ACKED] [<MAXLEN | MINID> [= | ~] threshold [LIMIT count]] <* | id> field value [field value ...]`
- 功能：发送消息
- 参数：
  - `NOMKSTREAM`：如果队列不存在是否创建队列，默认创建。使用该参数表示不创建
  - `<MAXLEN | MINID> [= | ~] threshold [LIMIT count]`：设置消息队列的最大消息数量
  - `<* | id>`：消息的 ID，\* 代表由 Redis 自动生成。格式是 `时间戳-递增数字`
  - `field value`：发送到队列中的消息，称为 Entry。格式就是多个 key-value 键值对

```bash
## 创建名为 users 的队列，并向其中发送一个消息，并使用 Redis 自动生成 ID
XADD users * name "zhh" age 18
# "1760774662027-0"
```

#### 12.3.2 `XLEN`

- 语法：`XLEN key`
- 功能：返回 Stream 中消息条目的数量

<<< @/db/codes/redis/xlen.sh

#### 12.3.3 `XREAD`

- 语法：`XREAD [COUNT count] [BLOCK milliseconds] STREAMS key [key ...] id [id ...]`
- 功能：读取消息
- 参数：
  - `COUNT count`：每次读取消息的最大数量
  - `BLOCK milliseconds`：当没有消息时是否阻塞
  - `STREAMS key [key ...]`：要从哪个队列读消息，`key` 就是队列名
  - `id [id ...]`：起始 ID，只返回大于该 ID 的消息
    - `0`：代表从第一个消息开始
    - `$`：代表从最新消息开始

<<< @/db/codes/redis/xread.sh

### 12.4 Stream-消费者组

将多个消费者划分到一个组中，监听同一个队列。具备下列特点：

1. 消息分流：队列中的消息会分流给组内的不同消费者，而不是重复消费，从而加快消息处理的速度
2. 消息标识：消费者组会维护一个表示，记录最后一个被处理的消息，哪怕消费者宕机重启，还会从标识之后读取消息。确保每个消息都被消费
3. 消息确认：消费者获取消息后，消息处于 pending 状态，并存入一个 pending-list。当处理完成后需要通过 XACK 来确认消息，标记消息为已处理，才会从 pending-list 移除

#### 12.4.1 `XGROUP CREATE`

- 语法：`XGROUP CREATE key group <id | $> [MKSTREAM] [ENTRIESREAD entries-read]`
- 功能：创建消费者组
- 参数：
  - `key`：队列名称
  - `group`：消费者组名称
  - `<id | $>`：起始 ID 标识
  - `MKSTREAM`：不存在时自动创建

#### 12.4.2 `XGROUP DESTROY`

- 语法：`XGROUP DESTROY key group`
- 功能：删除消费者组

#### 12.4.3 `XREADGROUP GROUP`

- 语法：`XREADGROUP GROUP group consumer [COUNT count] [BLOCK milliseconds] [NOACK] STREAMS key [key ...] id [id ...]`
- 功能：从消费者组读取消息
- 参数：
  - `group`：消费者组名称
  - `consumer`：消费者名称，如果不存在则自动创建

## 13 GEO

GEO 是 Geolocation 的简写形式，代表地理坐标。Redis 在 3.2 版本中加入了对 GEO 的支持，允许存储地理坐标信息，帮助我们根据经纬度来检索数据。

### 13.1 `GEOADD`

- 语法：`GEOADD key [NX | XX] longitude latitude member [longitude latitude member ...]`
- 功能：添加一个地理空间信息

<<< @/db/codes/redis/geoadd.sh

### 13.2 `GEODIST`

- 语法：`GEODIST key member1 member2 [M | KM]`
- 功能：返回两个点之间的距离
- 参数：
  - `M`：以米为单位
  - `KM`：以千米为单位

<<< @/db/codes/redis/geodist.sh{10-15}

### 13.3 `GEOHASH`

- 语法：`GEOHASH key [member [member ...]]`
- 功能：返回 hash 字符串形式的 `member` 坐标

<<< @/db/codes/redis/geohash.sh{5,6}

### 13.4 `GEOPOS`

- 语法：`GEOPOS key [member [member ...]]`
- 功能：返回 `member` 的坐标

<<< @/db/codes/redis/geopos.sh{5-7}

### 13.5 `GEORADIUS`

> 废弃。推荐采用 `GEOSEARCH` 和 `GEOSEARCHSTORE`

- 语法：`GEORADIUS key longitude latitude radius <M | KM | FT | MI>`
- 功能：指定圆心、半径，找到该圆内包含的所有 `member`

### 13.6 `GEOSEARCH`

- 语法：`GEOSEARCH key <FROMMEMBER member | FROMLONLAT longitude latitude> <BYRADIUS radius <M | KM | FT | MI> | BYBOX width height <M | KM | FT | MI>> [ASC | DESC] [COUNT count [ANY]] [WITHCOORD] [WITHDIST] [WITHHASH]`
- 功能：在指定范围内搜索 `member`，并按照与指定点的距离排序后返回。范围可以是圆形或矩形

<<< @/db/codes/redis/geosearch.sh{10-19}

### 13.6 `GEOSEARCHSTORE`

- 语法：`GEOSEARCHSTORE destination source <FROMMEMBER member | FROMLONLAT longitude latitude> <BYRADIUS radius <M | KM | FT | MI> | BYBOX width height <M | KM | FT | MI>> [ASC | DESC] [COUNT count [ANY]] [STOREDIST]`
- 功能：与 `GEOSEARCH` 功能一致，不过可以把结果存到一个指定的 `key`

## 14 Bitmap

Redis 中是利用 string 类型数据结构实现 Bitmap，因此最大上限是 512 M。

### 14.1 `SETBIT`

- 语法：`SETBIT key offset value`
- 功能：在 `offset` 处存入一个 `1` 或 `0`

<<< @/db/codes/redis/setbit.sh

### 14.2 `GETBIT`

- 语法：`GETBIT key offset`
- 功能：在 `offset` 处存入一个 `1` 或 `0`

<<< @/db/codes/redis/getbit.sh{3-6}

### 14.3 `BITFIELD`

- 语法：`BITFIELD key GET encoding offset`
- 功能：查询 Bitmap 中 bit 数组指定位置的值

<<< @/db/codes/redis/bitfield.sh

### 14.4 `BITCOUNT`

- 语法：`BITCOUNT key`
- 功能：统计 Bitmap 中值为 `1` 的 bit 位的数量

<<< @/db/codes/redis/bitcount.sh

### 14.5 `BITPOS`

- 语法：`BITPOS key bit [start [end [BYTE | BIT]]]`
- 功能：查询 bit 数组中指定范围内第一个 0 或 1 出现的位置

<<< @/db/codes/redis/bitpos.sh

## 15 Redis 持久化

### 15.1 RDB

全称 Redis Database Backup file（Redis 数据备份文件），也被叫做 Redis 数据快照。简单来说就是把内存中的所有数据都记录到磁盘中。当 Redis 实例故障重启后，从磁盘读取快照文件，恢复数据。

快照文件称为 RDB 文件。

<<< @/db/codes/redis/rdb.sh

Redis 停机时会执行一次 RDB。

Redis 内部有触发 RDB 的机制，可以在 `redis.conf` 文件中找到，格式如下：

<<< @/db/codes/redis/rdb.conf

`BGSAVE` 开始时会 fork 主进程得到子进程，子进程<span style="color:red;">共享</span>主进程的内存数据。完成 fork 后读取内存数据并写入 RDB 文件。

fork 采用的是 copy-on-write 技术：

- 当主进程执行读操作时，访问共享内存；
- 当主进程执行写操作时，则会拷贝一份数据，执行写操作。

### 15.2 AOF

全称 Append Only File（追加文件）。Redis 处理的每一个写命令都会记录在 AOF 文件，可以看作是命令日志文件。

AOF 默认时关闭的，需要修改 `redis.conf` 配置文件来开启 AOF：

<<< @/db/codes/redis/aof.conf

因为是记录命令，AOF 文件会比 RDB 文件大很多。而且 AOF 会记录对同一个 key 的多次写操作，但只有最后一次写操作才有意义。通过执行 `BGREWRITEAOF` 命令，可以让 AOF 文件执行重写功能，用最少的命令达到相同的效果。

Redis 也会在触发阈值时自动去重写 AOF 文件。阈值也可以在 `redis.conf` 中配置：

<<< @/db/codes/redis/auto-aof.conf

## 16 Redis 主从

单节点 Redis 的并发能力是有上限的，要进一步提高 Redis 的并发能力，就需要搭建主从集群，实现读写分离。

### 16.1 搭建主从集群

在一台机器上开三个 Redis 搭建主从集群。

#### 16.1.1 配置

开启 RDB，关闭 AOF。

创建 3 个文件夹，把配置文件放进去，并且修改 `port` 为三个不同的端口，然后启动三个 Redis 服务：

<<< @/db/codes/redis/master-slaver_conf.sh

#### 16.1.2 开启

现在三个实例还没有任何关系，要配置主从可以使用 `REPLICAOF` 或者 `SLAVEOF` 命令。

有临时和永久两种模式：

- 修改配置文件（永久生效）
  - 在 `redis.conf` 添加一行配置：`replicaof <masterip> <masterport>`
- 使用 redis-cli 客户端连接到 redis 服务，执行 `SLAVEOF` 命令（重启后失效）

  <<< @/db/codes/redis/master-slaver_slaveof.sh

### 16.2 数据同步原理

#### 16.2.1 全量同步

主从第一次同步是全量同步。

master 判断 slave 是否首次同步数据，需依赖两个核心概念：

- Replication id（简称 replid）：是数据集的唯一标记，id 一致代表属于同一数据集。每个 master 有专属 replid，slave 会继承 master 的 replid。
- offset（偏移量）：随 repl_backlog 中记录的数据增加而递增。slave 完成同步时会记录当前 offset，若 slave 的 offset 小于 master 的 offset，说明 slave 数据落后，需要更新。

因此 slave 进行数据同步时，必须向 master 声明自身的 replication id 和 offset，master 才能据此判断需要同步的内容。

#### 16.2.2 增量同步

主从第一次同步是全量同步，但如果 slave 重启后同步，则执行增量同步。

::: warning
repl_backlog 的大小存在上限，当写满后会覆盖最早的数据。

若 slave 断开连接的时间过久，导致其未备份的数据被 repl_backlog 覆盖，此时无法基于日志进行增量同步，只能再次执行全量同步。
:::

可以从以下几个方面来优化 Redis 主从集群：

- 在 master 节点配置 `repl-diskless-sync yes`，全量同步时无需将 RDB 文件写入磁盘，直接通过网络传输给 slave，避免磁盘 IO 开销。
- 减少 Redis 单节点的内存使用量，降低生成 RDB 文件时的磁盘 IO 压力。
- 适当增大 repl_backlog 的大小，同时及时恢复故障 slave 节点，尽可能避免全量同步。
- 限制单个 master 的 slave 数量；若 slave 过多，采用 “主-从-从” 的链式结构，让部分 slave 从其他 slave 同步数据，减轻 master 的负载。

## 17 Redis 哨兵

Redis 通过哨兵机制保障主从集群的自动故障恢复。哨兵的结构和作用如下：

- **监控**：持续检查 master 和 slave 节点是否正常运行。
- **自动故障恢复**：当 master 节点故障时，自动将一个 slave 节点提升为新的 master；待原故障 master 恢复后，会以新 master 的 slave 身份重新加入集群。
- **通知**：作为客户端的服务发现组件，在集群发生故障转移（主从切换）时，将最新的集群节点信息推送至客户端。

### 17.1 服务状态监控

Redis Sentinel 通过心跳机制监测集群实例状态：每隔 1 秒向集群中每个实例发送 `ping` 命令：

- 主观下线：单个 Sentinel 节点发现某实例未在规定时间内响应 `ping` 命令时，则认为该实例主观下线（仅当前 Sentinel 的独立判断）。
- 客观下线：若超过指定数量（quorum）的 Sentinel 节点都判定该实例为主观下线，则该实例客观下线。quorum 值最好超过 Sentinel 实例数量的一半。

### 17.2 选举新的 master

当 Sentinel 检测到 master 故障后，会从 slave 节点中选举新的 master，其选择依据如下：

- 若 slave 与原 master 的断开时间超过指定阈值（`down-after-milliseconds \* 10`），则直接排除该 slave。
- 比较 slave 的 `replica-priority` 配置值，数值越小优先级越高；若该值为 0，该 slave 不参与选举。
- 若 `replica-priority` 相同，比较 slave 的 offset，数值越大说明数据越新，优先级越高。
- 最后比较 slave 的运行 ID，数值越小优先级越高。

### 17.3 故障转移的实现流程

在 Sentinel 选定新的 master（以 slave1 为例）后，通过以下步骤完成故障转移：

1. Sentinel 向 slave1 发送 `slaveof no one` 命令，使其脱离从节点身份，成为新的 master。
2. Sentinel 向集群中剩余的所有 slave 发送 `slaveof [新master的IP] [新master的端口]` 命令，让这些 slave 将新 master 作为主节点，开始从新 master 同步数据。
3. Sentinel 将原故障的 master 标记为 slave，待其恢复后，会自动以 slave 身份连接新 master。

## 18 Redis 分片

### 18.1 分片集群结构

主从与哨兵机制可实现集群的高可用及高并发读，但仍存在两个核心局限：

- 无法承载海量数据存储
- 无法支撑高并发写请求

分片集群是针对上述问题的解决方案，其核心特征如下：

- 集群包含多个 master 节点，每个 master 负责存储不同的数据分片
- 每个 master 可关联多个 slave 节点
- master 节点之间通过 ping 机制互检状态
- 客户端可访问集群任意节点，请求会被自动转发至对应数据所在的 master 节点

### 18.2 散列插槽

Redis 分片集群中，所有 master 节点对应 0-16383 共 16384 个插槽，数据 Key 不与节点直接绑定，而是与插槽绑定，绑定逻辑由 Key 的有效部分计算得出。

根据 Key 中是否包含 `{}`，分为两种情况：

1. **包含 `{}` 且内部有字符**：以 `{}` 中的内容作为有效部分（例如 Key 为 `{seee}num`，有效部分是 `seee`）；
2. **不包含 `{}`**：以整个 Key 作为有效部分。

插槽值的计算方式：

1. 基于 Key 的有效部分，通过 CRC16 算法生成哈希值；
2. 将哈希值对 16384 取余，结果即为该 Key 对应的插槽（Slot）值。

### 18.3 集群伸缩

添加节点的命令：`add-node new_host:new_port existing_host:existing_port`

默认为 master 节点，通过参数 `--cluster-slave --cluster-master-id <arg>` 可指定为 slave 节点。

<<< @/db/codes/redis/add-node.sh

此时 `7004` 节点上没有插槽。需要通过 `reshard <host:port> or <host> <port>` 命令分配。

<<< @/db/codes/redis/reshard.sh

### 18.4 故障转移

当集群中有一个 master 宕机会发生什么：

1. 首先是该实例与其它实例失去连接
2. 然后是疑似宕机
3. 最后是确认下线，自动提升一个 slave 为新的 master

#### 18.4.1 数据迁移

利用 `CLUSTER FAILOVER` 命令可以手动让集群中的某个 master 宕机，切换到执行 `CLUSTER FAILOVER` 命令的这个 slave 节点，实现无感知的数据迁移。

- 语法：`CLUSTER FAILOVER [FORCE | TAKEOVER]`
- 时间复杂度：`O(1)`
- 可选项：
  - `FORCE`：省略对 offset 的一致性校验
  - `TAKEOVER`：忽略数据一致性、忽略 master 状态和其他 master 的意见

### 18.5 Go Redis 访问分片集群

Ring 分片客户端，是采用了一致性 HASH 算法在多个 redis 服务器之间分发 key，每个节点承担一部分 key 的存储。

Ring 客户端会监控每个节点的健康状况，并从 Ring 中移除掉宕机的节点，当节点恢复时，会再加入到 Ring 中。

<<< @/db/codes/redis/go-redis_shard.go

## 19 最佳实践

### 19.1 键值设计

Redis 的 Key 虽然可以自定义，但遵循以下最佳实践约定，可提升可读性、维护性与性能：

- 遵循基本格式：`[业务名]:[数据名]:[ID]`
- 长度不超过 44 字节
- 不包括特殊字符

例如，我们的登录业务，保存用户信息，其 Key 是这样的：`login:user:101`

优点：

1. 可读性强
2. 避免 Key 冲突
3. 方便管理
4. 更节省内存：Key 是 string 类型，底层编码包含 int、embstr 和 raw 三种。embstr 在小于 44 字节使用，采用连续内存空间，内存占用更小

### 19.2 批处理优化

网络传输耗时远大于 Redis 执行命令的耗时。

多个命令可以在一次网络请求中发送，而不是每次请求发送一个命令。

#### 19.2.1 MSET

Redis 提供了很多以 `Mxxx` 命名的命令，可实现数据的批量插入操作，常见的包括：

- `MSET`
- `HMSET`

::: warning
也不要在一次批处理中传输太多命令，否则单次命令占用带宽过多，会导致网络阻塞。
:::

#### 19.2.2 Pipeline

`MSET` 虽然可以批处理，但只能操作部分数据类型，因此如果有对复杂数据类型的批处理需要，建议使用 Pipeline 功能：

<<< @/db/codes/redis/go-redis_pipe.go

::: warning
Pipeline 的多个命令之间不具备原子性。
:::

#### 19.2.3 集群下的批处理

在 Redis 集群环境中，MSET、Pipeline 这类批处理命令需要在一次请求中携带多条命令，但存在一个核心约束：批处理涉及的多个 key 必须被映射到同一个插槽。

若多个 key 分散在不同插槽，批处理命令会执行失败。

|          | 串行命令                       | 串行 slot                                                                                                | 并行 slot                                                                                                | hash_tag                                                    |
| -------- | ------------------------------ | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------- |
| 实现思路 | for 循环遍历，依次执行每个命令 | 在客户端计算每个 key 的 slot，将 slot 一致分为一组，每组都利用 Pipeline 批处理。<br>**串行执行各组命令** | 在客户端计算每个 key 的 slot，将 slot 一致分为一组，每组都利用 Pipeline 批处理。<br>**并行执行各组命令** | 将所有 key 设置相同的 hash_tag，则所有 key 的 slot 一定相同 |
| 耗时     | N 次网络耗时 + N 次命令耗时    | m 次网络耗时 + N 次命令耗时<br>m = key 的 slot 个数                                                      | 1 次网络耗时 + N 次命令耗时                                                                              | 1 次网络耗时 + N 次命令耗时                                 |
| 优点     | 实现简单                       | 耗时较短                                                                                                 | 耗时非常短                                                                                               | 耗时非常短、实现简单                                        |
| 缺点     | 耗时非常久                     | 实现稍复杂<br>slot 越多，耗时越久                                                                        | 实现复杂                                                                                                 | 容易出现数据倾斜                                            |

### 19.3 服务端优化

#### 19.3.1 持久化配置

Redis 持久化可保障数据安全，但会增加额外开销，建议遵循以下配置规则：

1. 用来做缓存的 Redis 实例尽量不要开启持久化功能
2. 建议关闭 RDB 持久化功能，使用 AOF 持久化
3. 利用脚本定期在 slave 节点做 RDB，实现数据备份
4. 设置合理的 rewrite 阈值，避免频繁的 bgrewrite
5. 配置 `no-appendfsync-on-rewrite = yes`，禁止在 rewrite 期间做 aof，避免因 AOF 引起的阻塞

部署有关建议：

1. Redis 实例的物理机要预留足够内存，应对 fork 和 rewrite
2. 单个 Redis 实例内存上限不要太大，例如 4G 或 8G。可以加快 fork 的速度、减少主从同步、数据迁移压力
3. 不要与 CPU 密集型应用部署在一起
4. 不要与高硬盘负载应用一起部署。例如：数据库、消息队列

#### 19.3.2 慢查询

在 Redis 执行时耗时超过某个阈值的命令，称为慢查询。

慢查询的阈值可以通过配置指定：

- `slowlog-log-slower-than`：慢查询阈值，单位是微秒。默认是 10000，建议 1000。

慢查询会被放入慢查询日志中，日志的长度有上限，可以通过配置指定：

- `slowlog-max-len`：慢查询日志（本质是一个队列）长度。默认是 128，建议 1000。

查看慢查询日志列表：

- `SLOWLOG LEN`：查询慢查询日志长度
- `SLOWLOG GET [count]`：读取 `count` 条慢查询日志
- `SLOWLOG RESET`：清空慢查询列表

## 20 数据结构

### 20.1 动态字符串 SDS

Redis 中保存的 key 是字符串，value 往往是字符串或者字符串的集合。可见字符串是 Redis 中最常用的一种数据结构。

不过 Redis 没有使用 C 语言中的字符串，因为 C 语言字符串存在很多问题：

- 获取字符串的长度需要通过运算
- 非二进制安全
- 不可修改

<<< @/db/codes/redis/string.c

Redis 构建了一种新的字符串结构，称为**简单动态字符串**（Simple Dynamic String），简称 SDS。

例如，我们执行命令：`SET name zhh`

那么 Redis 底层就创建了两个 SDS，其中一个是包含 `name` 的 SDS，另一个是包含 `zhh` 的 SDS。

Redis 是 C 语言实现的，其中 SDS 是一个结构体，源码如下：

<<< @/db/codes/redis/sSDS.c

SDS 被称为动态字符串，核心特性是具备动态扩容能力。

内存预分配规则：

- 若新字符串长度 < 1M：`新空间 = 扩展后长度 × 2 + 1`
- 若新字符串长度 > 1M：`新空间 = 扩展后长度 + 1M + 1`

优点：

1. 获取字符串长度的时间复杂度为 `O(1)`
2. 支持动态扩容
3. 减少内存分配次数
4. 二进制安全（可存储任意二进制数据，不依赖 `\0` 标识结尾）

### 20.2 IntSet

IntSet 是 Redis 中 set 集合的一种实现方式，基于整数数组来实现，并且具备长度可变、有序等特征。

结构如下：

<<< @/db/codes/redis/sIntSet.c

其中的 `encoding` 包含三种模式，表示存储的整数大小不同：

<<< @/db/codes/redis/dIntSet.c

为了方便查找， Redis 会将 intset 中所有整数按照升序依次保存在 `contents` 数组中。

IntSet 升级流程（例如从 `INTSET_ENC_INT16` 升到 `INTSET_ENC_INT32`）：

1. 升级编码为 `INTSET_ENC_INT32`，每个整数占 4 字节，并按照新的编码方式及元素个数扩容数组
2. 倒序依次将数组中的元素拷贝到扩容后的正确位置
3. 将待添加的元素放入数组末尾
4. 最后，将 intset 的 `encoding` 改为 `INTSET_ENC_INT32`。

添加元素的函数源码：

<<< @/db/codes/redis/fIntsetAdd.c

### 20.3 Dict

Redis 是一个键值型（Key-Value Pair）的数据库，可以根据键实现快速的增删改查。而键与值的映射关系正是通过 Dict 来实现的。

Dict 由三个部分组成，分别是：哈希表（DictHashTable）、哈希节点（DictEntry）、字典（Dict）

<<< @/db/codes/redis/sDictht.c

<<< @/db/codes/redis/sDictEntry.c

当我们向 Dict 添加键值对时，Redis 首先根据 key 计算出 hash 值 `h`，然后利用 `h & sizemask` 来计算元素应该存储到数组中的哪个索引位置。

Dict 中的 HashTable 就是数组结合单向链表的实现，当集合中元素较多时，必然导致哈希冲突增多，链表过长，则查询效率会大大降低。

Dict 在每次新增键值对时都会检查**负载因子**（`LoadFactor = used / size`），满足以下两种情况时就会触发**哈希表扩容**：

- 哈希表的 `LoadFactor >= 1`，并且服务器没有执行 `BGSAVE` 或者 `BGREWRITEAOF` 等后台进程；
- 哈希表的 `LoadFactor > 5`

<<< @/db/codes/redis/fDictExpandIfNeeded.c

不管是扩容还是收缩，必定会创建新的哈希表，导致哈希表的 `size` 和 `sizemask` 变化，而 `key` 的查询与 `sizemask` 有关。因此必须对哈希表中的每一个 `key` 重新计算索引，插入新的哈希表，这个过程称为 **rehash**。过程是这样的：

1. 计算新 hash 表的 `realSize`，值取决于当前要做的是扩容还是收缩：
   - 如果是扩容，则新 `size` 为第一个大于等于 `used + 1` 的 $2^n$
   - 如果是收缩，则新 `size` 为第一个大于等于 `used` 的 $2^n$（不得小于 4）
2. 按照新的 `realSize` 申请空间，创建 `dictht`，并复制给 `dict.ht_table[1]`
3. 设置 `dict.rehashidx = 0`，表示开始 rehash
4. ~~将 `dict.ht_table[0]` 中的每一个 `dictEntry` 都 rehash 到 `dict.ht_table[1]`~~
5. 将 `dict.ht_table[1]` 赋值给 `dict.ht_table[0]`，给 `dict.ht_table[1]`初始化为空哈希表，释放原来 `dict.ht_table[0]` 的内存

Dict 的 rehash 并不是一次性完成的。试想一下，如果 Dict 中包含数百万的 entry，要在一次 rehash 完成，极有可能导致主线程阻塞。因此 Dict 的 rehash 是分多次、渐进式的完成，因此称为**渐进式 rehash**。上面的流程第 4 点应为：

4. 每次执行增删改查操作时，都检查一下 `dict.rehashidx` 是否大于 `-1`，如果是则将 `dict.ht_table[0].table[rehashidx]` 的 entry 链表 rehash 到 `dict.ht_table[1]`，并且将 `rehashidx++`。直至 `dict.ht_table[0]` 的所有数据都 rehash 到 `dict.ht_table[1]`

### 20.4 ZipList

ZipList 是一种特殊的“双端列表”，由一系列特殊编码的连续内存块组成。可以在任意一端进行压入/弹出操作，并且该操作的时间复杂度为 `O(1)`。

ZipList 中的 Entry 并不像普通链表那样记录前后节点的指针，因为记录两个指针要占用 16 个字节，浪费内存。而是采用了下面的结构：

| prevrawlen | encoding | content |
| ---------- | -------- | ------- |

- `prevrawlen`：前一节点的长度，占 1 或 5 个字节
  - 如果前一字节的长度小于 254 字节，则采用 1 个字节保存这个长度值
  - 如果前一字节的长度大于等于 254 字节，则采用 5 个字节保存这个长度值，第一个字节为 `0xfe`，后四个字节才是真实长度数据
- `encoding`：编码属性，记录 `content` 的数据类型（字符串还是整数）以及长度，占用 1、2、5 个字节
- `content`：负责保存节点的数据，可以是字符串或整数

### 20.5 QuickList

问题 1：ZipList 虽然节省内存，但申请内存必须是连续空间，如果内存占用过多，申请内存效率很低。怎么办？

- 为了缓解这个问题，我们必须限制 ZipList 的长度和 entry 大小。

问题 2：但是我们要存储大量数据，超过了 ZipList 最佳的上限该怎么办？

- 我们可以创建多个 ZipList 来分片存储数据。

问题 3：数据拆分后比较分散，不方便管理和查找，这多个 ZipList 怎么建立联系？

- Redis 在 3.2 版本引入了新的数据结构 QuickList，它是一个双端链表，只不过链表中的每个节点都是一个 ZipList。

为了避免 QuickList 中的每个 ZipList 中 entry 过多，Redis 提供了一个配置项：`list-max-ziplist-size` 来限制。

- 如果值为正，则代表 ZipList 的允许的 entry 个数的最大值
- 如果值为负，则代表 ZipList 的最大内存大小，分 5 种情况：
  1. `-1`：每个 ZipList 的内存占用不能超过 4kb
  2. `-2`：每个 ZipList 的内存占用不能超过 8kb
  3. `-3`：每个 ZipList 的内存占用不能超过 16kb
  4. `-4`：每个 ZipList 的内存占用不能超过 32kb
  5. `-5`：每个 ZipList 的内存占用不能超过 64kb

其默认值为 `-2`：

<<< @/db/codes/redis/list-max-ziplist-size.sh

除了控制 ZipList 的大小，QuickList 还可以对节点的 ZipList 做压缩。通过配置项 `list-compress-depth` 来控制。因为链表一般都是从首尾访问较多，所以首尾是不压缩的。这个参数是控制首尾不压缩的节点个数：

- `0`：特殊值，代表不压缩
- `1`：标示 QuickList 的首尾各有 1 个节点不压缩，中间节点压缩
- `2`：标示 QuickList 的首尾各有 2 个节点不压缩，中间节点压缩
- 以此内推

默认值：

<<< @/db/codes/redis/list-compress-depth.sh

以下是 QuickList 和 QuickListNode 的结构源码：

<<< @/db/codes/redis/sQuickList.c

<<< @/db/codes/redis/sQuickListNode.c

### 20.6 SkipList

SkipList（跳表）首先是链表，但是与传统链表相比有几点差异：

- 元素按照升序排序存储
- 节点可能包含多个指针，指针跨度不同

<<< @/db/codes/redis/sZSkipList.c

<<< @/db/codes/redis/sZSkipListNode.c

### 20.7 RedisObject

Redis 中的任意数据类型的键和值都会被封装为一个 RedisObject，也叫做 Redis 对象，源码如下：

<<< @/db/codes/redis/sRedisObject.c

Redis 种会根据存储的数据类型不同，选择不同的编码方式，共包含 11 种不同类型：

<<< @/db/codes/redis/sRedisObjEncoding.c

每种数据类型的使用的编码方式如下：

| 数据类型     | 编码方式                                                                    |
| ------------ | --------------------------------------------------------------------------- |
| `OBJ_STRING` | `OBJ_ENCODING_INT`、`OBJ_ENCODING_EMBSTR`、`OBJ_ENCODING_RAW`               |
| `OBJ_LIST`   | `OBJ_ENCODING_LINKEDLIST`、`OBJ_ENCODING_ZIPLIST`、`OBJ_ENCODING_QUICKLIST` |
| `OBJ_SET`    | `OBJ_ENCODING_INTSET`、`OBJ_ENCODING_HT`                                    |
| `OBJ_ZSET`   | `OBJ_ENCODING_ZIPLIST`、`OBJ_ENCODING_HT`、`OBJ_ENCODING_SKIPLIST`          |
| `OBJ_HASH`   | `OBJ_ENCODING_ZIPLIST`、`OBJ_ENCODING_HT`                                   |

### 20.8 五种数据类型

#### 20.8.1 String

String 是 Redis 中最常见的数据存储类型：

- 其基本编码方式是 **RAW**，基于简单动态字符串（SDS）实现，存储上限为 512mb。
- 如果存储的 SDS 长度小于 44 字节，则会采用 **EMBSTR** 编码，此时 object head 与 SDS 是一段连续空间。申请内存时只需要调用一次内存分配函数，效率更高。
- 如果存储的字符串是整数值，并且大小在 `LONG_MAX` 范围内，则会采用 **INT** 编码：直接将数据保存在 RedisObject 的 `ptr` 指针位置（刚好 8 字节），不再需要 SDS 了。

<<< @/db/codes/redis/obj_encoding_string.sh

#### 20.8.2 List

Redis 的 List 类型可以从首、尾操作列表中的元素。

哪一个数据结构能满足上述特征？

- LinkedList：普通链表，可以从双端访问，内存占用较高，内存碎片较多
- ZipList：压缩列表，可以从双端访问，内存占用低，存储上限低
- QuickList：LinkedList + ZipList，可以从双端访问，内存占用较低，包含多个 ZipList，存储上限高

在 3.2 版本之后，Redis 统一采用 QuickList 来实现 List。

#### 20.8.3 Set

Set 是 Redis 中的单列集合，满足以下特点：

- 不保证有序性
- 保证元素唯一（可以判断元素是否存在）
- 求交集、差集、并集

为了查询效率和唯一性，set 采用了 HT 编码（Dict）。Dict 中的 key 用来存储元素，value 统一为 null。

当存储的所有数据都是整数，并且元素数量不超过 `set-max-intset-entries` 时，Set 会采用 IntSet 编码，以节省内存。

#### 20.8.4 ZSet

ZSet 也就是 SortedSet，其中每一个元素都需要指定一个 score 值和 member 值：

- 可以根据 score 值排序
- member 必须唯一
- 可以根据 member 查询分数

因此，zset 底层数据结构必须满足**键值存储、键必须唯一、可排序**这几个需求。哪种编码结构可以满足？

- SkipList：可以排序，并且可以同时存储 score 和 member
- HT（Dict）：可以键值存储，并且可以根据 key 找 value

<<< @/db/codes/redis/sZset.c

<<< @/db/codes/redis/fCreateZsetObject.c

#### 20.8.5 Hash

Hash 结构与 Redis 中的 Zset 非常类似：

- 都是键值存储
- 都需根据键获得值
- 键必须唯一

区别如下：

- zset 的键是 member，值是 score；hash 的键和值都是任意值
- zset 要根据 score 排序；hash 则无需排序

因此，hash 底层采用的编码与 Zset 也基本一致，只需要把排序 SkipList 去掉即可：

- Hash 结构默认采用 ZipList 编码，用以节省内存。ZipList 中相邻的两个 entry 分别保存 field 和 value
- 当数据量较大时，Hash 结构会转为 HT 编码，也就是 Dict，触发条件有两个：
  1. ZipList 中的元素数量超过了 `hash-max-ziplist-entries`（默认 512）
  2. ZipList 中的任意 entry 大小超过了 `hash-max-ziplist-value`（默认 64 字节）

<<< @/db/codes/redis/hash-max-ziplist.sh

## 21 Redis 网络模型

### 21.1 用户空间和内核空间

任何 Linux 发行版，其系统内核都是 Linux。我们的应用都需要通过 Linux 内核与硬件交互。

为了避免用户应用导致冲突甚至内核崩溃，用户应用与内核是分离的：

- 进程的寻址空间会划分为两部分：**内核空间、用户空间**
- **用户空间**只能执行受限的命令（Ring3），而且不能直接调用系统资源，必须通过内核提供的接口来访问
- **内核空间**可以执行特权命令（Ring0），调用一切资源

Linux 系统为了提高 IO 效率，会在用户空间和内核空间都加入缓冲区：

- 写数据时，要把用户缓冲数据拷贝到内核缓冲区，然后写入设备
- 读数据时，要从设备读取数据到内核缓冲区，然后拷贝到用户缓冲区

在《UNIX 网络编程》一书中，总结归纳了 5 种 IO 模型：

- 阻塞 IO（Blocking IO）
- 非阻塞 IO（Nonblocking IO）
- IO 多路复用（IO Multiplexing）
- 信号驱动 IO（Signal Driven IO）
- 异步 IO（Asynchronous IO）

### 21.2 阻塞 IO

顾名思义，阻塞 IO 就是两个阶段都必须阻塞等待。

### 21.3 非阻塞 IO

顾名思义，非阻塞 IO 的 `recvfrom` 操作会立即返回结果而不是阻塞用户进程。

虽然是非阻塞，但是性能并没有得到提升。而且忙等机制会导致 CPU 空转，CPU 使用率暴增。

### 21.4 IO 多路复用

无论是阻塞 IO 还是非阻塞 IO，用户都需要调用 `recvfrom` 来获取数据，差别在于无数据时的处理方案：

- 如果调用 `recvfrom` 时，恰好**没有**数据，阻塞 IO 会使进程阻塞，非阻塞 IO 使 CPU 空转，都不能充分发挥 CPU 的作用。
- 如果调用 `recvfrom` 时，恰好**有**数据，则用户进程可以直接进入第二阶段，读取并处理数据。

比如服务端处理客户端 Socket 请求时，在单线程情况下，只能依次处理每一个 socket，如果正在处理的 socket 恰好未就绪（数据不可读或不可写），线程就会被阻塞，所有其他客户端 socket 都必须等待，性能自然会很差。

**文件描述符**（File Descriptor）：简称 FD，是一个从 0 开始递增的无符号整数，用来关联 Linux 中的一个文件。在 Linux 中，一切皆文件，例如常规文件、视频、硬件设备等，当然也包括网络套接字（Socket）。

IO 多路复用：利用单个线程来同时监听多个 FD，并在某个 FD 可读、可写时得到通知，从而避免无效的等待，充分利用 CPU 资源。

不过监听 FD 的方式、通知的方式又有多种实现，常见的有：

- select
- poll
- epoll

差异：

- select 和 poll 只会通知用户进程有 FD 就绪，但不确定具体是哪个 FD，需要用户进程逐个遍历 FD 来确认
- epoll 则会在通知用户进程 FD 就绪的同时，把已就绪的 FD 写入用户空间

#### 21.4.1 select

select 是 Linux 中最早的 I/O 多路复用实现方案。存在的问题：

- 需要将整个 `fd_set` 从用户空间拷贝到内核空间，select 结束还要再次拷贝回用户空间
- select 无法得知具体是哪个 fd 就绪，需要遍历整个 `fd_set`
- `fd_set` 监听的 fd 数量不能超过 1024

#### 21.4.2 poll

poll 模式对 select 模式做了简单改进，但性能提升不明显，部分关键代码如下：

<<< @/db/codes/redis/poll.c

IO 流程：

1. 创建 `pollfd` 数组，向其中添加关注的 `fd` 信息，数组大小自定义
2. 调用 `poll` 函数，将 `pollfd` 数组拷贝到内核空间，转链表存储，无上限
3. 内核遍历 `fd`，判断是否就绪
4. 数据就绪或超时后，拷贝 `pollfd` 数组到用户空间，返回就绪 `fd` 数量 n
5. 用户进程判断 n 是否大于 0
6. 大于 0 则遍历 `pollfd` 数组，找到就绪的 `fd`

与 select 相比：

- select 模式中的 `fd_set` 大小固定为 1024，而 `pollfd` 在内核中采用链表，理论上无上限
- 监听 FD 越多，每次遍历消耗时间也越久，性能反而会下降

#### 21.4.3 epoll

epoll 模式是对 select 和 poll 的改进，它提供了三个函数：

<<< @/db/codes/redis/epoll.c

##### 23.4.3.1 事件通知机制

当 FD 有数据可读时，我们调用 `epoll_wait` 就可以得到通知。但是事件通知的模式有两种：

- LevelTriggered：简称 LT。当 FD 有数据可读时，会重复通知多次，直至数据处理完成。是 Epoll 的默认通知。
- EdgeTriggered：简称 ET。当 FD 有数据可读时，只会被通知一次，不管数据是否处理完成。
