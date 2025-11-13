struct dictEntry {
    void *key; // 键

    union {
        void *val;
        uint64_t u64;
        int64_t s64;
        double d;
    } v; // 值

    struct dictEntry *next; // 下一个 Entry 的指针
};