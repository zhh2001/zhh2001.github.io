typedef struct quicklistNode {
    struct quicklistNode *prev; // 前一个节点指针
    struct quicklistNode *next; // 下一个节点指针
    unsigned char *entry; // 当前节点的 ZipList 指针
    size_t sz; // 当前节点的 ZipList 的字节大小
    unsigned int count : 16; // 当前节点的 ZipList 的 entry 个数
    unsigned int encoding : 2; // 编码方式：1，ZipList；2，lzf压缩模式
    unsigned int container : 2; // 数据容器类型
    unsigned int recompress : 1; // 是否被解压缩
    unsigned int attempted_compress : 1; // 测试用
    unsigned int dont_compress : 1;
    unsigned int extra : 9; // 预留字段
} quicklistNode;