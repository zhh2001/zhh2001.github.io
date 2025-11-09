// 服务端
func (s *Server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
    return &pb.UserResponse{...}, nil
}

// 客户端
res, err := client.GetUser(ctx, &pb.UserRequest{...})