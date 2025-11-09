table t {
    actions = {
        a,               // 可以出现在任何地方
        @tableonly b,    // 只能出现在表中
        @defaultonly c,  // 只能出现在默认动作中
    }
    /* 省略主体 */
}