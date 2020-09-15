package pool

import (
	"errors"
	"net/http"
	"time"

	"github.com/PhamDuyKhang/go-lb/internal/config"
	"github.com/PhamDuyKhang/go-lb/internal/datastructure"
	"github.com/PhamDuyKhang/go-lb/internal/services"
	"github.com/PhamDuyKhang/go-lb/internal/util"
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
			retries := util.GetRetry(r)
			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					logger.Debugf("retry %d times at node %s ", retries, bk.GetID())
					bk.Serve(w, util.SetRetry(r, retries+1))
				}
				return
			}
			bk.SetHealth(true)
			logger.Debugf("node %s has been marked down", bk.GetID())
			errRes := config.ForwardError{
				SourceURL: r.URL.String(),
				Status:    "Fail",
				Message:   "Service you chose has been die",
			}
			util.JSONWrite(w, http.StatusServiceUnavailable, errRes)
			return
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
		if b != nil && b.IsAlive() {
			rs.backendList.EnQueues(b)
			return
		}
		logger.Debugf("the node %s is down remove that node forever", b.GetID())
		return
	}
}

func (rs *RoundRobinStrategies) AddListenerDiscovery() {

}
