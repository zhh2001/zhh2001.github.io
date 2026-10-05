// Redis 8.2.3，省略暂停标志和 metadata 等字段
struct dict {
    dictType *type;
    dictEntry **ht_table[2];      // 两张桶数组，第二张在 rehash 时使用
    unsigned long ht_used[2];   // 两张表的元素数
    long rehashidx;              // -1 表示不在 rehash
    // 其他字段省略
    signed char ht_size_exp[2];  // 容量为 2^exp，-1 表示未分配
};
