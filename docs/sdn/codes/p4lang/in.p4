extern void f(in bit a, in bit<3> b = 2, in bit<5> c);

void g() {
    f(a = 1, b = 2, c = 3);  // 合法
    f(a = 1, c = 3);         // 合法，等价于上一个调用，b 使用默认值
    f(1, 2, 3);              // 合法，等价于上一个调用
    f(1, 3);                 // 非法
}