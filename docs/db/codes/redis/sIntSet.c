typedef struct intset {
    uint32_t encoding; // 编码方式，支持存放 16、32、64 位整数
    uint32_t length; // 元素个数
    int8_t contents[]; // 整数数组，保存集合元素
} intset;

intset *intsetNew(void);

intset *intsetAdd(intset *is, int64_t value, uint8_t *success);

intset *intsetRemove(intset *is, int64_t value, int *success);

uint8_t intsetFind(intset *is, int64_t value);

int64_t intsetRandom(intset *is);

int64_t intsetMax(intset *is);

int64_t intsetMin(intset *is);

uint8_t intsetGet(intset *is, uint32_t pos, int64_t *value);

uint32_t intsetLen(const intset *is);

size_t intsetBlobLen(intset *is);

int intsetValidateIntegrity(const unsigned char *is, size_t size, int deep);