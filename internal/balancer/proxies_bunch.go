package balancer

import (
	"fmt"

	"github.com/r-mol/balanser_highload_system/internal/proxy"
)

type proxiesBunch interface {
	Len() int
	Get(idx int) proxy.Client
}

type proxyWithWeight struct {
	proxy.Client
	weight int32
}

type weightedProxiesBunch []*proxyWithWeight

func (b weightedProxiesBunch) Len() int                 { return len(b) }
func (b weightedProxiesBunch) Get(idx int) proxy.Client { return b[idx].Client }

func getAvailableProxy(proxies proxiesBunch, marker int) (proxy.Client, error) {
	for i := 0; i < proxies.Len(); i++ {
		tryProxy := (marker + i) % proxies.Len()
		p := proxies.Get(tryProxy)
		if p.IsAvailable() {
			return p, nil
		}
	}
	return nil, fmt.Errorf("all proxies are unavailable")
}
