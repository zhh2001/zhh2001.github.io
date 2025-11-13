// zset 结构
typedef struct zset {
    dict *dict; // Dict 指针
    zskiplist *zsl; // SkipList 指针
} zset;