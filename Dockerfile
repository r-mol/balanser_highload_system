FROM golang:alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o balancer cmd/main.go

EXPOSE 5000

CMD ["./balancer --config=config/example-config.yaml --address=127.0.0.1:5000"]










