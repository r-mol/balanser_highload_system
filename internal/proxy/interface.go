package proxy

import (
	"net/url"
	"time"
)

type Client interface {
	GetLoad() int32
	GetHost() string
	IsAvailable() bool
	SetHealthCheck(check func(addr *url.URL) bool, period time.Duration)
	Stop()
}
