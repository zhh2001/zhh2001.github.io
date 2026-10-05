// Redis 8.2.3，省略阻塞、WATCH、统计等字段
typedef struct redisDb {
    kvstore *keys;                   // 键空间
    kvstore *expires;                // 带键级到期时间的对象索引
    ebuckets hexpires;               // Hash 字段过期维护
    // 其他字段省略
    int id;
    long long avg_ttl;
    unsigned long expires_cursor;   // 主动过期扫描游标
} redisDb;
