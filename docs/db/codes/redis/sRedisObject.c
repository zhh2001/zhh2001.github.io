// Redis 8.2.3，LRU_BITS=24，OBJ_REFCOUNT_BITS=30
struct redisObject {
    unsigned type : 4;
    unsigned encoding : 4;
    unsigned lru : LRU_BITS;                 // LRU 时钟或 LFU 元信息
    unsigned iskvobj : 1;                    // 是否采用 kvobj 布局
    unsigned expirable : 1;                  // 是否附有键级到期时间
    unsigned refcount : OBJ_REFCOUNT_BITS;
    void *ptr;                              // 数据指针或 int 编码整数
};
