# Load balancer for grpc requests

## Local building & running 💻

### Building ⚙️

```bash
go build -o balancer ./cmd/main.go           
```

### Running 🚀

```bash
./balancer start --config=./config/example-config.yaml --address=0.0.0.0:8080 --prometheus_address=0.0.0.0:1234
```


## Docker building & running 🐳

### Building ⚙️

```bash
docker build -t loadbalance-api .               
```

### Running 🚀

```bash
 docker run --mount type=bind,source="$(pwd)"/config/example-config.yaml,target=/app/config/example-config.yaml,readonly --name api --rm -p 8080:8080 loadbalance-api    
```