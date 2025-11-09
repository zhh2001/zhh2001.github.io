header H { bit<32> x; bit<32> y; }
H h;
h = { 10, 12 };
h = { y = 12, x = 10 };