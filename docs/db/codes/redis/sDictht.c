typedef struct dict {
    // entry 数组
    dictEntry **table; // 数组中保存的是指向 entry 的指针
    dictType *type;
    unsigned long size; // 哈希表大小
    unsigned long sizemask; // 哈希表大小的掩码，总等于 size-1
    unsigned long used; // entry 个数
    void *privdata;
} dict;