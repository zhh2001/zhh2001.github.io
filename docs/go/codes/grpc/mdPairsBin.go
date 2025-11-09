md := metadata.Pairs(
    "key", "string value",
    "key-bin", string([]byte{96, 102}),  // 二进制数据在发送前将被编码（base64），并在传输后被解码。
)