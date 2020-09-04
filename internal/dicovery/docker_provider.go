package dicovery

import (
	"context"
	"fmt"
	"github.com/PhamDuyKhang/go-lb/internal/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func GetListBackend(networkName string) []string {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	containerList, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		logrus.Error("can't get list container")
		return nil
	}
	var backendList []string
	for _, container := range containerList {
		isEnable := container.Labels[config.DiscoveryLabel]
		logrus.Debug(isEnable)
		if isEnable != "enable" {
			continue
		}
		if container.State != "running" {
			continue
		}
		logrus.Debug("the network")
		ipConfig := container.NetworkSettings.Networks["docker_kapp"]
		for k, v := range container.NetworkSettings.Networks {
			logrus.Debugf("key:%vx value:%vx", k, v)
		}
		address := ipConfig.IPAddress
		port := container.Ports[0].PrivatePort
		if address != "" && port != 0 {
			url := fmt.Sprintf("http://%s:%d", address, port)
			backendList = append(backendList, url)
		}
	}

	return backendList
}
