package health

import (
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"
)

// TODO
var (
	defaultHealthCheckTimeout = 10 * time.Second
	defaultHealthCheckPeriod  = 10 * time.Second
)

var defaultHealthCheck = func(addr *url.URL) bool {
	// Dial a connection to the target server
	conn, err := net.DialTimeout("tcp", addr.Host, defaultHealthCheckTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()

	return true
}

type ProxyHealth struct {
	Origin *url.URL

	mu          sync.Mutex
	check       func(addr *url.URL) bool
	period      time.Duration
	cancel      chan struct{}
	isAvailable bool
}

func New(o ...Option) (*ProxyHealth, error) {
	opts := &Options{
		check:  defaultHealthCheck,
		period: defaultHealthCheckPeriod,
	}

	for _, option := range o {
		option(opts)
	}

	switch {
	case opts.origin == nil:
		return nil, fmt.Errorf("\"Origin\" is not provided")
	case !opts.isAvailable:
		opts.isAvailable = opts.check(opts.origin)
	}

	h := &ProxyHealth{
		Origin:      opts.origin,
		mu:          opts.mu,
		check:       opts.check,
		period:      opts.period,
		cancel:      opts.cancel,
		isAvailable: opts.isAvailable,
	}
	h.run()

	return h, nil
}

func (h *ProxyHealth) IsAvailable() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.isAvailable
}

func (h *ProxyHealth) SetHealthCheck(check func(addr *url.URL) bool, period time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.stop()
	h.check = check
	h.period = period
	h.cancel = make(chan struct{})
	h.isAvailable = h.check(h.Origin)
	h.run()
}

func (h *ProxyHealth) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.stop()
}

func (h *ProxyHealth) run() {
	checkHealth := func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		isAvailable := h.check(h.Origin)
		h.isAvailable = isAvailable
	}

	go func() {
		t := time.NewTicker(h.period)
		for {
			select {
			case <-t.C:
				checkHealth()
			case <-h.cancel:
				t.Stop()
				return
			}
		}
	}()
}

func (h *ProxyHealth) stop() {
	if h.cancel != nil {
		h.cancel <- struct{}{}
		close(h.cancel)
		h.cancel = nil
	}
}
