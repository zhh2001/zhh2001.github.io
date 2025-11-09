struct S {
    bit<32> a;
    bit<32> b;
}

const S x = { 10, 20 };
const S x = { a = 10, b = 20 };
const S x = (S) { a = 10, b = 20 };