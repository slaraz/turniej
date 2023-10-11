protoc --go_out=../gra_go/proto --go_opt=paths=source_relative \
    --go-grpc_out=../gra_go/proto --go-grpc_opt=paths=source_relative \
    gra.proto