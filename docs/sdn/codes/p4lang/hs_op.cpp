// 类型声明
struct hs_t {
    bit<32> nextIndex;
    bit<32> size;
    h[n] data; // 普通数组
}

// 实例声明和初始化
hs_t hs;
hs.nextIndex = 0;
hs.size = n;