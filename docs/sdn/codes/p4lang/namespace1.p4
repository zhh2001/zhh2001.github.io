const bit<32> x = 2;
control c() {
    int<32> x = 0;
    apply {
        x = x + (int<32>).x;  // x 是 int<32> 的局部变量，
                              // .x 是顶层的 bit<32> 变量
    }
}