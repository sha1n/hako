FROM alpine:latest

EXPOSE 8080

ADD bin/hako-linux-amd64 /opt/hako

CMD ["/opt/hako", "start", "-p", "8080"]
