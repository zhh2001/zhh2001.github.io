// Objects encoding
#define OBJ_ENCODING_RAW        0 // raw 编码的动态字符串
#define OBJ_ENCODING_INT        1 // long 类型的整数的字符串
#define OBJ_ENCODING_HT         2 // hash table
#define OBJ_ENCODING_ZIPMAP     3 // 弃用
#define OBJ_ENCODING_LINKEDLIST 4 // 双端链表
#define OBJ_ENCODING_ZIPLIST    5 // 压缩列表
#define OBJ_ENCODING_INTSET     6 // 整数集合
#define OBJ_ENCODING_SKIPLIST   7 // 跳表
#define OBJ_ENCODING_EMBSTR     8 // embstr 的动态字符串
#define OBJ_ENCODING_QUICKLIST  9 // 快速列表
#define OBJ_ENCODING_STREAM    10 // Stream