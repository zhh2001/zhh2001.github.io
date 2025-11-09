const bit<4> x = 1;
control p() {
    const bit<8> x = 8;   // 局部变量 x 的声明覆盖了全局 x
    const bit<4> y = .x;  // 当前使用的是顶层 x
    const bit<8> z = x;   // 当前为 p 的局部变量 x
    apply {}
}