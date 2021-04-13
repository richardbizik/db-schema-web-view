FROM alpine:latest

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /db-shema-web-view
RUN mkdir -p ./_webapp/dist

ADD db-schema-web-view ./db-shema-web-view
ADD connections.yaml ./
ADD _webapp/dist ./_webapp/dist

RUN chmod +x ./db-shema-web-view

EXPOSE 9090

CMD ["./db-shema-web-view"]