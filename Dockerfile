FROM alpine:3.8
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 \
    && apk add --no-cache openssh-client

ENV PORT 9999
EXPOSE $PORT

COPY druid /
CMD ["/druid"]

VOLUME [ "/opt/druid" ]
