binary=sycm_server
dockeruser=realwrtoff
gituser=realwrtoff
repository=rest_grpc
sub_path=sycm
version=$(shell git describe --tags)

export PATH:=${PATH}:${GOPATH}/bin:
export GOPOXY=https://goproxy.cn

.PHONY: all
all: vendor build

build: cmd/*.go internal/*/*.go api/gen/go/api/*.go scripts/version.sh Makefile vendor
	@echo "compile"
	@go build -ldflags "-X 'main.AppVersion=`sh scripts/version.sh`'" cmd/main.go && \
	mkdir -p build/${repository}/bin && mv main build/${repository}/bin/${binary} && \
	mkdir -p build/${repository}/configs && cp configs/* build/${repository}/configs && \
	mkdir -p build/${repository}/log

vendor: go.mod
	@echo "install golang dependency"
	go mod tidy
	go mod vendor

codegen: api/sycm.proto
	mkdir -p api/gen/go && mkdir -p api/gen/swagger
	protoc -I.. -I. --gofast_out=plugins=grpc,paths=source_relative:api/gen/go/ $<
	protoc -I.. -I. --grpc-gateway_out=logtostderr=true,paths=source_relative:api/gen/go $<
	protoc -I.. -I. --swagger_out=logtostderr=true:api/gen/swagger $<

.PHONY: dockerenv
dockerenv:
	if [ -z "$(shell docker network ls --filter name=testnet -q)" ]; then \
		docker network create -d bridge testnet; \
	fi
	if [ -z "$(shell docker ps -a --filter name=go-build-env -q)" ]; then \
		docker run --name go-build-env --network testnet -d daocloud.io/realwrtoff/go-build-env:latest tail -f /dev/null; \
	fi

.PHONY: cleandockerenv
cleandockerenv:
	if [ ! -z "$(shell docker ps -a --filter name=go-build-env -q)" ]; then \
		docker stop go-build-env  && docker rm go-build-env; \
	fi
	if [ ! -z "$(shell docker network ls --filter name=testnet -q)" ]; then \
		docker network rm testnet; \
	fi

.PHONY: image
image: dockerenv
	rm -rf docker
	docker exec go-build-env rm -rf /data/src/${gituser}/${repository}
	docker exec go-build-env mkdir -p /data/src/${gituser}/
	docker cp ../ go-build-env:/data/src/${gituser}/
	docker exec go-build-env bash -c "cd /data/src/${gituser}/${repository}/${sub_path} && make codegen && make build"
	docker cp go-build-env:/data/src/${gituser}/${repository}/${sub_path}/build/ docker/
	docker build --tag=realwrtoff/${repository}/${version} .

.PHONY: clean
clean:
	rm -rf build docker

.PHONY: deep_clean
deep_clean:
	rm -rf build docker log vendor scripts/tmp go.sum