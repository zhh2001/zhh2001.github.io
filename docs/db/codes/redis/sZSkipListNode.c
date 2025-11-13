/* ZSETs use a specialized version of Skiplists */
typedef struct zskiplistNode {
    sds ele; // 节点存储的值
    double score; // 节点的分数，排序、查找用
    struct zskiplistNode *backward; // 前一个节点指针
    struct zskiplistLevel {
        struct zskiplistNode *forward; // 下一个节点指针
        unsigned long span; // 索引跨度
    } level[]; // 多级索引数组
} zskiplistNode;