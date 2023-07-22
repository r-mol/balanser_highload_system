FROM golang:alpine

RUN apk update && \
    apk add protoc protobuf protobuf-dev && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN cd protos && \
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative data_transfer_api.proto

RUN go build -o balancer ./cmd/main.go

EXPOSE 8080 1234

ENTRYPOINT ["./balancer", "start", "--config=/app/config/example-config.yaml", "--address=0.0.0.0:8080", "--prometheus_address=0.0.0.0:1234"]