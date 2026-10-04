struct S {
    bit<32> a;
    bit<32> b;
}

const S x1 = { 10, 20 };
const S x2 = { a = 10, b = 20 };
const S x3 = (S) { a = 10, b = 20 };
