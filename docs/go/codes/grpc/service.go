type UserServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{Name: "user-" + req.GetId()}, nil
}
