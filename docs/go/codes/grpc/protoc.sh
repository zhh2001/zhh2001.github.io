protoc --proto_path=. \
  --go_out=. --go_opt=paths=source_relative \
  proto/msg.proto
