FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY app/ .

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o http-server-projeto-korp .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /build/http-server-projeto-korp .

EXPOSE 8080

CMD ["./http-server-projeto-korp"]
