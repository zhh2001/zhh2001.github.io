type userIDKey struct{}

ctx := context.WithValue(context.Background(), userIDKey{}, 12345)
userID, ok := ctx.Value(userIDKey{}).(int)
fmt.Println(userID, ok) // 12345 true
