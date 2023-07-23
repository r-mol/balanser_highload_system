package balancer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/r-mol/balanser_highload_system/internal/proxy"
	data_transfer_api "github.com/r-mol/balanser_highload_system/protos"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type LoadBalancer struct {
	data_transfer_api.UnimplementedKeyValueServiceServer
	proxies weightedProxiesBunch

	mu       sync.Mutex
	current  int32
	reqCount int32
	metrics  *Metrics
	Logger   *log.Logger
}

func (lb *LoadBalancer) GetValue(ctx context.Context, request *data_transfer_api.GetValueRequest) (*data_transfer_api.GetValueResponse, error) {
	p, err := lb.Next()
	if err != nil {
		return nil, fmt.Errorf("the server didn't respond: %s", err)
	}

	defer func() {
		lb.metrics.Stat("GetValue", p.GetHost(), time.Now(), err)
	}()

	conn, err := grpc.Dial(p.GetHost(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial grpc connection with \"%s\": %w", p.GetHost(), err)
	}
	defer conn.Close()

	client := data_transfer_api.NewKeyValueServiceClient(conn)

	response, err := client.GetValue(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get value from host \"%s\": %w", p.GetHost(), err)
	}

	lb.Logger.Infof("get value from host \"%s\"\n", p.GetHost())

	return response, nil
}

func (lb *LoadBalancer) StoreValue(ctx context.Context, request *data_transfer_api.StoreValueRequest) (*data_transfer_api.StoreValueResponse, error) {
	p, err := lb.Next()
	if err != nil {
		return nil, fmt.Errorf("the server didn't respond: %s", err)
	}
	defer func() {
		lb.metrics.Stat("StoreValue", p.GetHost(), time.Now(), err)
	}()

	conn, err := grpc.Dial(p.GetHost(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial grpc connection with \"%s\": %w", p.GetHost(), err)
	}
	defer conn.Close()

	client := data_transfer_api.NewKeyValueServiceClient(conn)

	response, err := client.StoreValue(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to store value from host \"%s\": %w", p.GetHost(), err)
	}

	lb.Logger.Infof("store value to host \"%s\"\n", p.GetHost())

	return response, nil
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
		Logger:  opts.logger,
		metrics: opts.metrics,
	}, nil
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
