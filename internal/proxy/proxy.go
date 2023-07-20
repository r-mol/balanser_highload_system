package proxy

import (
	"fmt"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

var (
	defaultHealthCheckTimeout = 10 * time.Second
	defaultHealthCheckPeriod  = 10 * time.Second
)

var defaultHealthCheck = func(addr *url.URL) bool {
	conn, err := net.DialTimeout("tcp", addr.Host, defaultHealthCheckTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()

	return true
}

type Proxy struct {
	origin *url.URL

	mu          sync.Mutex
	check       func(addr *url.URL) bool
	period      time.Duration
	cancel      chan struct{}
	isAvailable bool
	load        int32
}

func New(o ...Option) (*Proxy, error) {
	opts := &Options{
		check:  defaultHealthCheck,
		period: defaultHealthCheckPeriod,
	}

	for _, option := range o {
		option(opts)
	}

	switch {
	case opts.origin == nil:
		return nil, fmt.Errorf("\"origin\" is not provided")
	case !opts.isAvailable:
		opts.isAvailable = opts.check(opts.origin)
	}

	h := &Proxy{
		origin:      opts.origin,
		mu:          opts.mu,
		check:       opts.check,
		period:      opts.period,
		cancel:      opts.cancel,
		isAvailable: opts.isAvailable,
	}
	h.run()

	return h, nil
}

func (p *Proxy) GetLoad() int32 {
	return atomic.LoadInt32(&p.load)
}

func (p *Proxy) GetHost() string {
	return p.origin.Host
}

func (p *Proxy) IsAvailable() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.isAvailable
}

func (p *Proxy) SetHealthCheck(check func(addr *url.URL) bool, period time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.stop()
	p.check = check
	p.period = period
	p.cancel = make(chan struct{})
	p.isAvailable = p.check(p.origin)
	p.run()
}

func (p *Proxy) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stop()
}

func (p *Proxy) run() {
	checkHealth := func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		isAvailable := p.check(p.origin)
		p.isAvailable = isAvailable
	}

	go func() {
		t := time.NewTicker(p.period)
		for {
			select {
			case <-t.C:
				checkHealth()
			case <-p.cancel:
				t.Stop()
				return
			}
		}
	}()
}

func (p *Proxy) stop() {
	if p.cancel != nil {
		p.cancel <- struct{}{}
		close(p.cancel)
		p.cancel = nil
	}
}
