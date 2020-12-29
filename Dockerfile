FROM alpine:latest

MAINTAINER liushaobo <liushaobo101@gmail.com>

RUN mkdir -p  /dev-work

COPY dist/scanport-linux-amd64/scanport /usr/local/bin
RUN chmod +x /usr/local/bin/scanport

CMD scanport -h