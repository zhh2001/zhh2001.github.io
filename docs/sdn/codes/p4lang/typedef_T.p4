struct S<T> {
    T field;
}

// typedef S X;       // 非法：S 缺少类型实参
typedef S<bit<32>> X;  // 合法
