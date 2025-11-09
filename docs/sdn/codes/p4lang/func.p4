extern void f(in bit<32> x, out bit<16> y);
bit<32> xa = 0;
bit<16> ya;
f(xa, ya);          // 按照位置匹配参数
f(x = xa, y = ya);  // 按照名称匹配参数
f(y = ya, x = xa);  // 按照名称匹配参数，顺序任意
f(x = xa);          // 错误：参数不够
f(x = xa, x = ya);  // 错误：参数重复指定
f(x = xa, ya);      // 错误：只有部分参数指定了名称
f(z = xa, w = yz);  // 错误：没有名为 z 或 w 的参数
f(x = xa, y = 0);   // 错误：y 必须是一个左值