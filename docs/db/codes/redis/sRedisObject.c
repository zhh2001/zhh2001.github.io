struct redisObject {
    // 分别是 string、hash、list、set 和 zset，占 4 位
    unsigned type:4;

    // 底层编码方式，共有 11 种
    unsigned encoding:4;

    // lru 表示该对象最后一次被访问的时间，其占用 24 位。便于判断空闲太久的 key
    unsigned lru:LRU_BITS;

    unsigned iskvobj : 1;   /* 1 if this struct serves as a kvobj base */
    unsigned expirable : 1; /* 1 if this key has expiration time attached.
                             * If set, then this object is of type kvobj */

    // 对象引用计数器，计数器为 0 则说明对象无人引用，可以被回收
    unsigned refcount : OBJ_REFCOUNT_BITS;

    // 指针，指向存放实际数据的空间
    void *ptr;
};

/* The actual Redis Object */
#define OBJ_STRING 0    /* String object. */
#define OBJ_LIST 1      /* List object. */
#define OBJ_SET 2       /* Set object. */
#define OBJ_ZSET 3      /* Sorted set object. */
#define OBJ_HASH 4      /* Hash object. */