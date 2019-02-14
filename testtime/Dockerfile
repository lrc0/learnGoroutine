FROM alpine:edge
MAINTAINER ruicai.li@changhong.com
RUN echo "Asia/Shanghai" > /etc/timezones
WORKDIR /workspace
COPY ./huanhang .
RUN chmod +x huanhang
ENTRYPOINT ["./huanhang"]