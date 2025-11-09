header H { bit<32> x; bit<32> y; }
H h;
h = {#};  // 这会使报头 h 变为无效
if (h == {#}) {  // 这相当于条件 !h.isValid()
    // ...
}