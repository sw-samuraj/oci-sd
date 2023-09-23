FROM golang:alpine AS builder

WORKDIR /oci-sd

COPY ./oci-sd .

RUN apk update && \
    apk add --no-cache curl gettext make cmake jq && \
    make build && \
    rm -rf Makefile adapter config.go && \
    rm -rf config_test.go go.mod main.go oci vendor && \
    apk del curl make cmake jq && \
    rm -rf /var/cache/apk/*

FROM alpine

RUN apk update && \
    apk add --no-cache gettext

WORKDIR /oci-sd

COPY --from=builder /oci-sd .

COPY start.sh oci-sd.tmpl ./

CMD ["./start.sh"]
