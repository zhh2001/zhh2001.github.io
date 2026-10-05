// Redis 8.2.3，省略压缩状态标志
typedef struct quicklistNode {
    struct quicklistNode *prev;
    struct quicklistNode *next;
    unsigned char *entry;       // listpack、单个大元素或压缩数据
    size_t sz;                  // 未压缩数据的字节数
    unsigned int count : 16;    // 元素数
    unsigned int encoding : 2;  // RAW=1，LZF=2
    unsigned int container : 2; // PLAIN=1，PACKED=2
    // 其他标志省略
} quicklistNode;
