package mocks

type Proxy interface {
	New() (*Proxy, error)
	GetHost() string
	GetLoad() int32
	IsAvailable() bool
	SetHealthCheck()
	Stop()
	run()
}
