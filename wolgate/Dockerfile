FROM golang:1.21-alpine as builder

WORKDIR /app

COPY application/. .

RUN go mod tidy
RUN go build -o wolgate

FROM alpine:latest
COPY --from=builder /app/wolgate /wolgate
CMD ["/wolgate"]