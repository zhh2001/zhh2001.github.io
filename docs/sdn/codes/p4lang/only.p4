table t {
    actions = {
        a;               // 可用于普通表项或默认动作
        @tableonly b;    // 只能用于普通表项
        @defaultonly c;  // 只能用于默认动作
    }
    /* 省略主体 */
}
