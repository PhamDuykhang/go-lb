package golb

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type (

	//LoadBalancePool is a struct hold all information of a backend.
	//And distribute the request
	LoadBalancePool struct {
		bk      []*Backend
		current uint64
	}
	//LoadBalancer hold a pool of backend server
	LoadBalancer struct {
		pool *LoadBalancePool
	}
)

//NewLoadBalaningPool creat a load balacing pool form array of inside backend we have
func NewLoadBalaningPool(backendList []string) (*LoadBalancePool, error) {
	lbP := LoadBalancePool{}
	for i := range backendList {
		url, err := url.Parse(backendList[i])
		if err != nil {
			return nil, err
		}
		rveProxy := httputil.NewSingleHostReverseProxy(url)
		bk := Backend{
			TagertURL: url.String(),
			rpx:       rveProxy,
		}
		if bk.IsAlive() {
			lbP.bk = append(lbP.bk, &bk)
		}
	}
	return &lbP, nil
}

//NextIdx get next backend to forward request
func (lp *LoadBalancePool) NextIdx() int {
	return int(atomic.AddUint64(&lp.current, uint64(1)) % uint64(len(lp.bk)))
}

//Next get a avaliable backed for net request
//The function get a backend form a start of nextIdx to end nextIdx-1
func (lp *LoadBalancePool) Next() *Backend {
	n := lp.NextIdx()
	l := len(lp.bk) + n
	for i := n; i < l; i++ {
		idx := i % len(lp.bk)
		if lp.bk[idx].IsAlive() {
			if i != n {
				atomic.StoreUint64(&lp.current, uint64(idx))
			}
			return lp.bk[idx]
		}
	}
	return nil
}

//NewLoadBalancer get a backend pool and start the forward handler
func NewLoadBalancer(sevPool *LoadBalancePool) LoadBalancer {
	return LoadBalancer{
		pool: sevPool,
	}
}

//LoadBalance get a request and forward to one of our backends
func (lb LoadBalancer) LoadBalance(w http.ResponseWriter, r *http.Request) {
	p := lb.pool.Next()
	if p != nil {
		p.rpx.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
	return
}
