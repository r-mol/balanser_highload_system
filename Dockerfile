FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o balancer ./cmd/main.go

EXPOSE 8080

CMD ["./balancer", "start", "--config=/app/config/example-config.yaml", "--address=:8080"]