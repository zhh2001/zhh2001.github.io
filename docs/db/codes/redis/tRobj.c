struct redisObject
{
    unsigned type : 4;       // 对象类型
    unsigned encoding : 4;   // 编码方式
    unsigned lru : LRU_BITS; // LRU：以秒为单位记录最近一次访问时间，长度 24 bit
                             // LFU：高 16 位以分钟为单位记录最近一次访问时间，低 8 位记录逻辑访问次数

    unsigned iskvobj : 1;   // 1 if this struct serves as a kvobj base
    unsigned expirable : 1; // 1 if this key has expiration time attached.
                            // If set, then this object is of type kvobj

    unsigned refcount : OBJ_REFCOUNT_BITS; // 引用计数，计数为 0 则可以回收
    void *ptr;                             // 数据指针，指向真实数据
};