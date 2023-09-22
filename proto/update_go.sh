protoc --go_out=../go/gra_proto --go_opt=paths=source_relative \
    --go-grpc_out=../go/gra_proto --go-grpc_opt=paths=source_relative \
    gra.proto