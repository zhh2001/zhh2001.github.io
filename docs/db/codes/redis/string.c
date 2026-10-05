int main(void) {
    const char *s = "hello"; // 字符串字面量不能修改
    char editable[] = "hello";
    editable[0] = 'H';       // 字符数组可以修改
    return 0;
}
