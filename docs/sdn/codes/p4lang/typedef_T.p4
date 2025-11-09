struct S<T> {
    T field;
}

typedef S X;           // -- 非法：S没有类型参数
typedef S<bit<32>> X;  // -- 合法