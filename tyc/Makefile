export GOPROXY=https://goproxy.cn

build: cmd/main.go
	mkdir -p build/bin
	go build cmd/main.go && mv main build/bin

codegen: api/tyc.proto
	mkdir -p api/gen/go && mkdir -p api/gen/swagger
	protoc -I.. -I. --go_out=plugins=grpc,paths=source_relative:api/gen/go/ $<
	protoc -I.. -I. --grpc-gateway_out=logtostderr=true,paths=source_relative:api/gen/go $<
	protoc -I.. -I. --swagger_out=logtostderr=true:api/gen/swagger $<
