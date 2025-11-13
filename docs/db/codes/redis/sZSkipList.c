typedef struct zskiplist {
    struct zskiplistNode *header, *tail; // 头尾节点指针
    unsigned long length; // 节点数量
    int level; // 最大的索引层级，默认是 1
} zskiplist;