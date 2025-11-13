/* Insert an integer in the intset */
intset *intsetAdd(intset *is, int64_t value, uint8_t *success) {
    uint8_t valenc = _intsetValueEncoding(value); // 获取当前值编码
    uint32_t pos; // 要插入的位置
    if (success) *success = 1;

    // 判断编码是不是超过了当前 intset 的编码
    if (valenc > intrev32ifbe(is->encoding)) {
        /* This always succeeds, so we don't need to curry *success. */
        return intsetUpgradeAndAdd(is, value);
    } else {
        // 在当前 intset 中查找值与 value 一样的元素的角标 pos
        if (intsetSearch(is, value, &pos)) {
            if (success) *success = 0;
            return is;
        }

        // 数组扩容
        is = intsetResize(is, intrev32ifbe(is->length) + 1);
        if (pos < intrev32ifbe(is->length)) intsetMoveTail(is, pos, pos + 1);
    }

    _intsetSet(is, pos, value); // 插入新元素
    is->length = intrev32ifbe(intrev32ifbe(is->length) + 1);
    return is;
}

/* Upgrades the intset to a larger encoding and inserts the given integer. */
static intset *intsetUpgradeAndAdd(intset *is, int64_t value) {
    uint8_t curenc = intrev32ifbe(is->encoding); // 获取当前 intset 编码
    uint8_t newenc = _intsetValueEncoding(value); // 获取新编码
    int length = intrev32ifbe(is->length); // 获取元素个数
    int prepend = value < 0 ? 1 : 0;

    // 重置编码与数组大小
    is->encoding = intrev32ifbe(newenc);
    is = intsetResize(is, intrev32ifbe(is->length) + 1);

    // 倒序遍历，逐个搬运元素到新位置，_intsetGetEncoded 按照旧编码方式查找旧元素
    while (length--)
        _intsetSet(is, length + prepend, _intsetGetEncoded(is, length, curenc));

    // 插入新元素，prepend 决定是队首还是队尾
    if (prepend)
        _intsetSet(is, 0, value);
    else
        _intsetSet(is, intrev32ifbe(is->length), value);
    is->length = intrev32ifbe(intrev32ifbe(is->length) + 1);
    return is;
}