FROM alpine:edge
MAINTAINER ruicai.li@changhong.com
RUN echo "Asia/Shanghai" > /etc/timezones
WORKDIR /workspace
COPY ./docker/log .
RUN chmod +x log
ENTRYPOINT ["./log"]

