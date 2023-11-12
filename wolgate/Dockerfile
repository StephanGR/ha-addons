FROM golang:1.21-alpine as builder

WORKDIR /app

COPY application/main.go .
COPY application/go.mod .
COPY application/go.sum .

RUN go build -o wolgate

FROM alpine:latest

COPY --from=builder /app/wolgate /wolgate
CMD ["/wolgate"]