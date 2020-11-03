FROM centos:centos7

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >> /etc/timezone

COPY docker/ /
WORKDIR /rest_grpc
CMD [ "bin/sycm_server", "-c", "configs/sycm_server.json" ]