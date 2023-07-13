package balancer

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/r-mol/balanser_highload_system/internal/proxy"
)

type LoadBalancer struct {
	proxies weightedProxiesBunch

	mu       sync.Mutex
	current  int32
	reqCount int32
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
