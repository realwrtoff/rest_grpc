binary=tyc_server
dockeruser=realwrtoff
gituser=realwrtoff
repository=rest_grpc
version=$(shell git describe --tags)

export PATH:=${PATH}:${GOPATH}/bin:$(shell pwd)/third/protobuf/bin
export GOPOXY=https://goproxy.cn

.PHONY: all
all: vendor output

output: cmd/*.go internal/*/*.go api/gen/go/api/*.go scripts/version.sh Makefile vendor
	@echo "compile"
	@go build -ldflags "-X 'main.AppVersion=`sh scripts/version.sh`'" cmd/main.go && \
	mkdir -p output/${repository}/bin && mv main output/${repository}/bin/${binary} && \
	mkdir -p output/${repository}/configs && cp configs/* output/${repository}/configs && \
	mkdir -p output/${repository}/log

vendor: go.mod
	@echo "install golang dependency"
	go mod tidy
	go mod vendor

codegen: api/tyc.proto
	mkdir -p api/gen/go && mkdir -p api/gen/swagger
	protoc -I.. -I. --gofast_out=plugins=grpc,paths=source_relative:api/gen/go/ $<
	protoc -I.. -I. --grpc-gateway_out=logtostderr=true,paths=source_relative:api/gen/go $<
	protoc -I.. -I. --swagger_out=logtostderr=true:api/gen/swagger $<

.PHONY: protoc
protoc:
	@hash protoc 2>/dev/null || { \
		echo "install protobuf codegen tool protoc" && \
		mkdir -p third && cd third && \
		wget https://github.com/google/protobuf/releases/download/v3.2.0/protobuf-cpp-3.2.0.tar.gz && \
		tar -xzvf protobuf-cpp-3.2.0.tar.gz && \
		cd protobuf-3.2.0 && \
		./configure --prefix=`pwd`/../protobuf && \
		make -j8 && \
		make install && \
		cd ../.. && \
		protoc --version; \
	}

.PHONY: clean
clean:
	rm -rf output

.PHONY: deep_clean
deep_clean:
	rm -rf output api/gen