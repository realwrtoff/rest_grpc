FROM centos:centos7 as builder

RUN yum install -y make \
    && yum install -y git gcc \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" >> /etc/timezone

# golang
ENV GOPROXY=https://goproxy.cn
RUN curl -OL https://dl.google.com/go/go1.15.3.linux-amd64.tar.gz
RUN tar -xzvf go1.15.3.linux-amd64.tar.gz && mv go /usr/local
ENV PATH=$PATH:/usr/local/go/bin
ENV GOROOT=/usr/local/go
ENV GO111MODULE=on
RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get github.com/gogo/protobuf/proto
RUN go get github.com/gogo/protobuf/jsonpb
RUN go get github.com/gogo/protobuf/protoc-gen-gogo
RUN go get github.com/gogo/protobuf/gogoproto
RUN go get github.com/gogo/protobuf/protoc-gen-gofast

# 添加代码和编译
RUN git clone https://github.com/realwrtoff/rest_grpc.git
RUN cd rest_grpc/tyc && make codegen
RUN make output

FROM centos:centos7
COPY --from=builder /rest_grpc/tyc/output/ /
EXPOSE 7060

WORKDIR /rest_grpc
CMD [ "bin/tyc_server", "-c", "configs/tyc_server.json" ]