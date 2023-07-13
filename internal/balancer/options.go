package balancer

import (
	"sync"

	"github.com/r-mol/balanser_highload_system/internal/proxy"
)

type Options struct {
	proxies weightedProxiesBunch
	mu      sync.Mutex
}

type Option func(*Options)

func WithProxies(proxies map[*proxy.Proxy]int32) Option {
	return func(opts *Options) {
		opts.proxies = make(weightedProxiesBunch, 0, len(proxies))
		for p, w := range proxies {
			opts.proxies = append(opts.proxies, &proxyWithWeight{
				Proxy:  p,
				weight: w,
			})
		}
	}
}

func WithMutex(mu sync.Mutex) Option {
	return func(opts *Options) {
		opts.mu = mu
	}
}
