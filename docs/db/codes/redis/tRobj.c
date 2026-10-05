// Redis 8.2.3 对象头中的淘汰字段，省略其他成员
struct redisObject {
    unsigned type : 4;
    unsigned encoding : 4;
    unsigned lru : 24;
    // LRU：秒级时钟，存在回绕
    // LFU：高 16 位为分钟级时间，低 8 位为概率计数
    // 其他成员见前面的 redisObject 定义
};
