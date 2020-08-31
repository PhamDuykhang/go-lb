package dicovery

import (
	"errors"

	"github.com/docker/docker/client"
)

type (
	Discovery struct {
		client  *client.Client
		changeC chan ServiceMetadata
		errC    chan error
	}

	BackendManager interface {
		Connect() error
		ListenChange()
		HandlerChange()
	}
)

func NewDockerDiscovery() Discovery {
	cOut := make(chan ServiceMetadata)
	eOut := make(chan error)
	return Discovery{
		changeC: cOut,
		errC:    eOut,
	}
}

func (dd Discovery) GetOutChan() (chan ServiceMetadata, error) {
	if dd.changeC == nil {
		return nil, errors.New("out chanel is nil")
	}
	return dd.changeC, nil
}

func (dd Discovery) Connect() error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	dd.client = cli
	return nil
}

func (dd Discovery) ListenChange() {
	//filter := filters.NewArgs()
	//filter.Add("type", "container")
	//filter.Add("event", "start")
	//filter.Add("event", "die")
	//
	//o,errChan := dd.client.Events(context.TODO(), types.EventsOptions{
	//	Filters: filter,
	//})
	//for {
	//	select {
	//	case err := <-errChan:
	//		//logger here
	//		dd.errC <-err
	//	case d := <-o:
	//		switch d.Action {
	//		case string(Die):
	//			container,err := dd.client.ContainerInspect(context.Background(),d.Actor.ID)
	//			if err!= nil{
	//				dd.errC <-fmt.Errorf("can't get target container %s ",short(d.Actor.ID))
	//				break
	//			}
	//			container.HostConfig
	//		}
	//		case string(Start):
	//
	//		}
	//	}
}

func short(hexID string) string {
	return hexID[0:7]
}
