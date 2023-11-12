FROM golang:latest

WORKDIR /app

COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go build -o wolgate

CMD ["/app/wolgate"]