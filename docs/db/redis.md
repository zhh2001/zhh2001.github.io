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
