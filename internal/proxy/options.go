package proxy

import (
	"net/http/httputil"

	"github.com/r-mol/balanser_highload_system/internal/proxy/health"
)

type Options struct {
	health *health.ProxyHealth
	proxy  *httputil.ReverseProxy
	load   int32
}

type Option func(*Options)

func WithHealth(health *health.ProxyHealth) Option {
	return func(opts *Options) {
		opts.health = health
	}
}

func WithProxy(proxy *httputil.ReverseProxy) Option {
	return func(opts *Options) {
		opts.proxy = proxy
	}
}

func WithLoad(load int32) Option {
	return func(opts *Options) {
		opts.load = load
	}
}
