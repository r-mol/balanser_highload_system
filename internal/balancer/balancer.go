package balancer

import (
	"fmt"
	"github.com/r-mol/balanser_highload_system/internal/proxy"
	"net/http"
	"sync"
)

type Balancer struct {
	proxies weightedProxiesBunch

	mu       sync.Mutex
	current  int32
	reqCount int32
}

func New(o ...Option) (*Balancer, error) {
	opts := &Options{}

	for _, option := range o {
		option(opts)
	}

	switch {
	case opts.proxies == nil:
		return nil, fmt.Errorf("\"proxies\" is not provided")
	}

	return &Balancer{
		proxies: opts.proxies,
		mu:      opts.mu,
	}, nil
}

func (l *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, err := l.Next()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("The server didn't respond: %s", err)))
		return
	}
	p.ServeHTTP(w, r)
}

func (l *Balancer) Next() (*proxy.Proxy, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	currentProxy := l.proxies[l.current]
	if l.reqCount < currentProxy.weight {
		l.reqCount++
	} else {
		l.current = (l.current + 1) % int32(len(l.proxies))
		l.reqCount = 1
	}
	return getAvailableProxy(l.proxies, int(l.current))
}
