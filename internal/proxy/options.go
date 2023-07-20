package proxy

import (
	"net/url"
	"sync"
	"time"
)

type Options struct {
	origin *url.URL

	mu          sync.Mutex
	check       func(addr *url.URL) bool
	period      time.Duration
	cancel      chan struct{}
	isAvailable bool
	load        int32
}

type Option func(*Options)

func WithOrigin(origin *url.URL) Option {
	return func(opts *Options) {
		opts.origin = origin
	}
}

func WithLoad(load int32) Option {
	return func(opts *Options) {
		opts.load = load
	}
}

func WithMutex(mu sync.Mutex) Option {
	return func(opts *Options) {
		opts.mu = mu
	}
}

func WithCheck(check func(addr *url.URL) bool) Option {
	return func(opts *Options) {
		opts.check = check
	}
}

func WithPeriod(period time.Duration) Option {
	return func(opts *Options) {
		opts.period = period
	}
}

func WithCancel(cancel chan struct{}) Option {
	return func(opts *Options) {
		opts.cancel = cancel
	}
}

func WithIsAvailable(isAvailable bool) Option {
	return func(opts *Options) {
		opts.isAvailable = isAvailable
	}
}
