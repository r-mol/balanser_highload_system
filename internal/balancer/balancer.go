package balancer

import (
	"context"
	"fmt"
	data_transfer_api "github.com/r-mol/balanser_highload_system/protos"
	"google.golang.org/grpc"
	"net/http"
	"sync"

	"github.com/r-mol/balanser_highload_system/internal/proxy"
)

type LoadBalancer struct {
	data_transfer_api.UnimplementedKeyValueServiceServer
	proxies weightedProxiesBunch

	mu       sync.Mutex
	current  int32
	reqCount int32
}

func (lb *LoadBalancer) GetValue(ctx context.Context, request *data_transfer_api.GetValueRequest) (*data_transfer_api.GetValueResponse, error) {
	p, err := lb.Next()
	if err != nil {
		return nil, fmt.Errorf("the server didn't respond: %s", err)
	}

	conn, err := grpc.Dial(p.GetHealth().Origin.Hostname())
	if err != nil {
		panic(err)
	}

	client := data_transfer_api.NewKeyValueServiceClient(conn)

	return client.GetValue(ctx, request)
}

func (lb *LoadBalancer) StoreValue(ctx context.Context, request *data_transfer_api.StoreValueRequest) (*data_transfer_api.StoreValueResponse, error) {
	p, err := lb.Next()
	if err != nil {
		return nil, fmt.Errorf("the server didn't respond: %s", err)
	}

	conn, err := grpc.Dial(p.GetHealth().Origin.Host, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := data_transfer_api.NewKeyValueServiceClient(conn)

	return client.StoreValue(ctx, request)
}

func New(o ...Option) (*LoadBalancer, error) {
	opts := &Options{}

	for _, option := range o {
		option(opts)
	}

	switch {
	case opts.proxies == nil:
		return nil, fmt.Errorf("\"proxies\" is not provided")
	}

	return &LoadBalancer{
		proxies: opts.proxies,
		mu:      opts.mu,
	}, nil
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, err := lb.Next()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("The server didn't respond: %s", err)))
		return
	}
	p.ServeHTTP(w, r)
}

func (lb *LoadBalancer) Next() (*proxy.Proxy, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	currentProxy := lb.proxies[lb.current]
	if lb.reqCount < currentProxy.weight {
		lb.reqCount++
	} else {
		lb.current = (lb.current + 1) % int32(len(lb.proxies))
		lb.reqCount = 1
	}
	return getAvailableProxy(lb.proxies, int(lb.current))
}
