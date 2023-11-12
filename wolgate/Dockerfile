# WE BUILD THE APP HERE
FROM golang:1.21-alpine as builder

WORKDIR /app

COPY application/. .

RUN go mod tidy
RUN go build -o wolgate

# WE USE ALPINE TO RUN THE APP
FROM alpine:latest

RUN apk update && apk upgrade && apk add bash jq jsonlint

COPY --from=builder /app/wolgate /wolgate
COPY run.sh /run.sh

RUN chmod a+x /run.sh

CMD ["/bin/sh", "-c", "/run.sh"]