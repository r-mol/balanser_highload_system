package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/r-mol/balanser_highload_system/internal/proxy/health"
)

type Proxy struct {
	health *health.ProxyHealth
	proxy  *httputil.ReverseProxy
	load   int32
}

func New(o ...Option) (*Proxy, error) {
	opts := &Options{}

	for _, option := range o {
		option(opts)
	}

	switch {
	case opts.proxy == nil:
		return nil, fmt.Errorf("\"proxy\" is not provided")
	case opts.health == nil:
		return nil, fmt.Errorf("\"health\" is not provided")
	}

	return &Proxy{
		proxy:  opts.proxy,
		health: opts.health,
		load:   opts.load,
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&p.load, 1)
	defer atomic.AddInt32(&p.load, -1)
	p.proxy.ServeHTTP(w, r)
}

func (p *Proxy) IsAvailable() bool {
	return p.health.IsAvailable()
}

func (p *Proxy) GetHealth() *health.ProxyHealth {
	return p.health
}

func (p *Proxy) SetHealthCheck(check func(addr *url.URL) bool, period time.Duration) {
	p.health.SetHealthCheck(check, period)
}

func (p *Proxy) GetLoad() int32 {
	return atomic.LoadInt32(&p.load)
}
