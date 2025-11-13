typedef struct quicklist {
    quicklistNode *head; // 头节点指针
    quicklistNode *tail; // 尾节点指针
    unsigned long count; // 所有 ZipList 的 entry 的数量
    unsigned long len; // ZipLists 的总数量
    signed int fill: QL_FILL_BITS; // ZipLists 的 entry 上限，默认 -2
    unsigned int compress: QL_COMP_BITS; // 首尾不压缩的节点数量
    unsigned int bookmark_count: QL_BM_BITS; // 内存重分配时的书签数量，一般用不到
    quicklistBookmark bookmarks[]; // 内存重分配时的书签数组，一般用不到
} quicklist;