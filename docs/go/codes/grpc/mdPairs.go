md := metadata.Pairs(
    "key1", "val1",
    "key1", "val1-2", // "key1" 的值将变成 []string{"val1", "val1-2"}
    "key2", "val2",
)