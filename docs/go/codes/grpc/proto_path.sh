protoc \
  --proto_path=. \
  --proto_path=../third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/hello.proto
