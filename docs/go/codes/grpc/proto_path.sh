protoc \
    --proto_path=. \
    --proto_path=../third_party \
    --go_out=. \
    --go-grpc_out=. \
	hello.proto