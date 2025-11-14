typedef struct redisDb
{
    kvstore *keys;       // The keyspace for this DB. As metadata, holds keysizes histogram
    kvstore *expires;    // 存放每一个 key 及其对应的 TTL 存活时间，只包含设置了 TTL 的 key
    ebuckets hexpires;   // Hash expiration DS. Single TTL per hash (of next min field to expire)
    dict *blocking_keys; // Keys with clients waiting for data (BLPOP)
    dict *blocking_keys_unblock_on_nokey;
    dict *ready_keys;   // Blocked keys that received a PUSH
    dict *watched_keys; // WATCHED keys for MULTI/EXEC CAS
    int id;             // Database ID，0-15
    long long avg_ttl;  // 记录平均 TTL 时长
    unsigned long expires_cursor;
} redisDb;