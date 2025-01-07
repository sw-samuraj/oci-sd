FROM golang:1.18-alpine3.16 as builder
WORKDIR /app
COPY . .
RUN go mod download && \
    go build .
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/oci-sd /usr/bin/
CMD ["oci-sd"]
