package dicovery

import (
	"context"
	"errors"
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

func GetDockerContainerIP(dockerID string) ([]string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	ctn, err := cli.ContainerInspect(context.TODO(), dockerID)
	if err != nil {
		logrus.Error("can't get container", err)
		return nil, err
	}
	//js ,_ := json.Marshal(ctn)
	//logrus.Info(string(js))
	var urls []string
	port := ctn.NetworkSettings.Ports
	if ctn.NetworkSettings.Networks["docker_kapp"].IPAddress == "" {
		logrus.Error("can't get IP address")
		return nil, errors.New("ip nil")
	}

	var ports []string
	for k, _ := range port {
		ports = append(ports, k.Port())
	}
	for _, p := range ports {
		url := fmt.Sprintf("http://%s:%s", ctn.NetworkSettings.Networks["docker_kapp"].IPAddress, p)
		urls = append(urls, url)
	}

	return urls, nil
}
