// Redis 8.2.3
// QL_FILL_BITS、QL_COMP_BITS、QL_BM_BITS 与指针宽度有关
// 64 位构建中分别为 16、16、4
typedef struct quicklist {
    quicklistNode *head;
    quicklistNode *tail;
    unsigned long count;                    // 总元素数
    unsigned long len;                      // 节点数
    signed int fill : QL_FILL_BITS;
    unsigned int compress : QL_COMP_BITS;
    unsigned int bookmark_count : QL_BM_BITS;
    quicklistBookmark bookmarks[];
} quicklist;
