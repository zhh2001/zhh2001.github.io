control noargs();
package top(noargs c1, noargs c2);

control c() {
    @name(".foo.bar") table t { /* 主体省略 */ }
    apply { /* 主体省略 */ }
}
top(c(), c()) main;