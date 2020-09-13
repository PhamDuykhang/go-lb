package pool

import (
	"errors"
	"github.com/PhamDuyKhang/go-lb/internal/datastructure"
	"github.com/PhamDuyKhang/go-lb/internal/services"
	"github.com/PhamDuyKhang/go-lb/internal/util"
	"net/http"
)

type (
	//Poller is a interface of a pool's func need to have,
	Poller interface {
		Next() services.Backend
		AddNewNodeToPool(bk services.Backend)
	}
	//RoundRobinStrategies is an one of simpler load balancing strategies
	RoundRobinStrategies struct {
		backendList datastructure.BackendRingBuffer
	}
)

//NewRoundRobinStrategies create the round robin
func NewRoundRobinStrategies() RoundRobinStrategies {
	return RoundRobinStrategies{
		backendList: datastructure.NewRing(10),
	}
}

func (rs *RoundRobinStrategies) InitBackend(bks []services.Backend) error {

	if len(bks) == 0 {
		return errors.New("backend list is empty")
	}
	logger.Debug("adding backend to pool")

	for _, bk := range bks {
		rs.AddNewNodeToPool(bk)
	}
	return nil
}

//AddNewNodeToPool add new backend service
func (rs *RoundRobinStrategies) AddNewNodeToPool(bk services.Backend) {
	if bk.IsAlive() {
		logger.Debugf("adding container %s to pool", bk.GetID())
		beforeLen := rs.backendList.Len()
		bk.ErrorHandle(func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Errorf("[%s] Error:%s", bk.GetID(), err.Error())
			// handle retry and mark the backend is down
			util.JSONWrite(w, http.StatusServiceUnavailable, nil)
		})
		rs.backendList.EnQueues(bk)
		logger.Debugf("adding container successfully len: %d grow to: %d", beforeLen, rs.backendList.Len())
	}
	logger.Debugf("adding new backend to round robin pool is success now we have %d service inside", rs.backendList.Len())
	return
}

func (rs *RoundRobinStrategies) LoadBalancing(w http.ResponseWriter, r *http.Request) {
	logger.Debug("starting lb for request")
	for {
		b := rs.backendList.DeQueue()
		logger.Debugf("the backend %+v", b)
		if b != nil && b.IsAlive() {
			logger.Infof("forward request %s --------------> %s stating", r.URL.String(), b.Stat().URL)
			b.Serve(w, r)
			logger.Infof("forward request %s --------------> %s completed", r.URL.String(), b.Stat().URL)
			logger.Debugf("return backend to pool")
		}
		rs.backendList.EnQueues(b)
		return
	}
}
