package balancer

import (
	"fmt"

	"github.com/r-mol/balanser_highload_system/internal/proxy"
)

type proxiesBunch interface {
	Len() int
	Get(idx int) *proxy.Proxy
}

type proxyWithWeight struct {
	*proxy.Proxy
	weight int32
}

type weightedProxiesBunch []*proxyWithWeight

func (b weightedProxiesBunch) Len() int                 { return len(b) }
func (b weightedProxiesBunch) Get(idx int) *proxy.Proxy { return b[idx].Proxy }

func getAvailableProxy(proxies proxiesBunch, marker int) (*proxy.Proxy, error) {
	for i := 0; i < proxies.Len(); i++ {
		tryProxy := (marker + i) % proxies.Len()
		p := proxies.Get(tryProxy)
		if p.IsAvailable() {
			return p, nil
		}
	}
	return nil, fmt.Errorf("all proxies are unavailable")
}
