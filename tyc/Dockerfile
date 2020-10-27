FROM centos:centos7 as builder

RUN yum install -y kde-l10n-Chinese && yum reinstall -y glibc-common && localedef -c -f UTF-8 -i zh_CN zh_CN.utf8
RUN yum install -y epel-release
RUN yum install -y docker-io
RUN yum install -y make

RUN yum install -y make \
    && yum install -y git gcc wget glibc-headers gcc-c++ \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" >> /etc/timezone

# golang
ENV GOPROXY=https://goproxy.cn
RUN curl -OL https://dl.google.com/go/go1.15.3.linux-amd64.tar.gz
RUN tar -xzvf go1.15.3.linux-amd64.tar.gz && mv go /usr/local
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:$GOROOT/bin
ENV GOPATH=/gopath
ENV GOBIN=$GOPATH/bin
ENV PATH=$PATH:$GOBIN
ENV GO111MODULE=on

# grpc support
RUN yum install -y unzip
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.10.0/protoc-3.10.0-linux-x86_64.zip
RUN unzip protoc-3.10.0-linux-x86_64.zip
RUN go version

RUN mkdir -p /src/google.golang.org/grpc && \
    cd /src/google.golang.org/grpc && \
    export GOPATH=/ && \
    export GO111MODULE=on && \
    go get -u google.golang.org/grpc && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
    go get github.com/gogo/protobuf/proto && \
    go get github.com/gogo/protobuf/jsonpb && \
    go get github.com/gogo/protobuf/protoc-gen-gogo && \
    go get github.com/gogo/protobuf/gogoproto && \
    go get github.com/gogo/protobuf/protoc-gen-gofast && \

# 添加代码和编译
RUN cd / && git clone https://github.com/realwrtoff/rest_grpc.git
RUN cd rest_grpc/tyc && make all

FROM centos:centos7
COPY --from=builder /rest_grpc/tyc/output/ /
EXPOSE 7060

WORKDIR /rest_grpc
CMD [ "bin/tyc_server", "-c", "configs/tyc_server.json" ]