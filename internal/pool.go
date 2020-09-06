package golb

import (
	"context"
	"fmt"
	"github.com/PhamDuyKhang/go-lb/internal/dicovery"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

var mu = sync.Mutex{}

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

//NewLoadBalancingPool creat a load balancing pool form array of inside backend we have
func NewLoadBalancingPool(backendList []string) (*LoadBalancePool, error) {
	lbP := LoadBalancePool{}
	for i := range backendList {
		urlFull, err := url.Parse(backendList[i])
		if err != nil {
			return nil, err
		}
		rveProxy := httputil.NewSingleHostReverseProxy(urlFull)
		rveProxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
			logrus.Error("error ", err)
			_, _ = fmt.Fprint(writer, "error")
		}
		bk := Backend{
			TargetURL: urlFull.String(),
			rpx:       rveProxy,
		}
		bk.HeathCheck("check")
		if bk.IsAlive() {
			logrus.Infof("service %s ------------- OK", backendList[i])
			lbP.bk = append(lbP.bk, &bk)
		}
	}
	return &lbP, nil
}

func (lbP *LoadBalancePool) AddBackend(backendList []string) error {
	for i := range backendList {
		urlFull, err := url.Parse(backendList[i])
		if err != nil {
			return err
		}
		rveProxy := httputil.NewSingleHostReverseProxy(urlFull)
		rveProxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
			logrus.Error("error ", err)
			_, _ = fmt.Fprint(writer, "error")
		}
		bk := Backend{
			TargetURL: urlFull.String(),
			rpx:       rveProxy,
		}
		bk.HeathCheck("check")
		if bk.IsAlive() {
			logrus.Infof("service %s ------------- OK", backendList[i])
			mu.Lock()
			defer mu.Unlock()
			lbP.bk = append(lbP.bk, &bk)
		}
	}
	return nil
}

//NextIdx get next backend to forward request
func (lbP *LoadBalancePool) NextIdx() int {
	return int(atomic.AddUint64(&lbP.current, uint64(1)) % uint64(len(lbP.bk)))
}

//Next get a available backed for net request
//The function get a backend form a start of nextIdx to end nextIdx-1
func (lbP *LoadBalancePool) Next() *Backend {
	n := lbP.NextIdx()
	l := len(lbP.bk) + n
	for i := n; i < l; i++ {
		idx := i % len(lbP.bk)
		if lbP.bk[idx].IsAlive() {
			if i != n {
				atomic.StoreUint64(&lbP.current, uint64(idx))
			}
			return lbP.bk[idx]
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
		logrus.Infof("get request %v -----> %s", r.URL, p.TargetURL)
		p.rpx.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
	return
}

func (lbP LoadBalancePool) WatchChange() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	filter := filters.NewArgs()
	filter.Add("type", "container")
	filter.Add("event", "start")
	filter.Add("event", "die")
	eventsOpp := types.EventsOptions{
		Filters: filter,
	}
	cC, Ec := cli.Events(context.TODO(), eventsOpp)
	for {
		select {
		case evenData := <-cC:
			if evenData.Action == "start" {
				logrus.Info("a new node is started")
				listURl, err := dicovery.GetDockerContainerIP(evenData.ID)
				if err != nil {
					logrus.Error("can't get url")
				} else {
					err := lbP.AddBackend(listURl)
					if err != nil {
						logrus.Error("can't add new ")
					}
				}

			}
		case err := <-Ec:
			logrus.Errorf("fail to listen even %v", err)
		}
	}
}
