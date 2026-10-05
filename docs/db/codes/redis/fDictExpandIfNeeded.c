#define DICT_OK 0
#define DICT_ERR 1

typedef enum {
    DICT_RESIZE_ENABLE,
    DICT_RESIZE_AVOID,
    DICT_RESIZE_FORBID,
} dictResizeEnable;

int dictExpandIfNeeded(dict *d) {
    // 如果正在 rehash，则返回 OK
    if (dictIsRehashing(d)) return DICT_OK;

    // 如果哈希表为空，则初始化哈希表大小为默认大小 4
    if (DICTHT_SIZE(d->ht_size_exp[0]) == 0) {
        dictExpand(d, DICT_HT_INITIAL_SIZE);
        return DICT_OK;
    }

    // 允许调整时负载因子达到 1，或未禁止调整时达到强制阈值 4
    // dict_force_resize_ratio 在 Redis 8.2.3 中为 4
    if ((dict_can_resize == DICT_RESIZE_ENABLE &&
         d->ht_used[0] >= DICTHT_SIZE(d->ht_size_exp[0])) ||
        (dict_can_resize != DICT_RESIZE_FORBID &&
         d->ht_used[0] >= dict_force_resize_ratio * DICTHT_SIZE(d->ht_size_exp[0]))) {
        // 扩容大小为 used+1，底层会对扩容大小进行判断，实际上找的是首个 >= used+1 的 2^n
        if (dictTypeResizeAllowed(d, d->ht_used[0] + 1))
            dictExpand(d, d->ht_used[0] + 1);
        return DICT_OK;
    }
    return DICT_ERR;
}
