FROM centos:centos7

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >> /etc/timezone

COPY docker/ /
WORKDIR /rest_grpc
CMD [ "bin/tyc_server", "-c", "configs/tyc_server.json" ]