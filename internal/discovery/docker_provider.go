package discovery

import (
	"context"
	"errors"
	"fmt"

	"github.com/PhamDuyKhang/go-lb/internal/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("package", "discovery")

type DockerContainer struct {
	ContainerName    string
	ContainerID      string
	ContainerAddress string
	ContainerPort    string
}

func GetListBackend(networkName string) []DockerContainer {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	containerList, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		logger.Error("can't get list container")
		return nil
	}
	var backendList []DockerContainer
	for _, container := range containerList {
		isEnable := container.Labels[config.DiscoveryLabel]
		if isEnable != "enable" {
			continue
		}
		if container.State != "running" {
			continue
		}
		ipConfig := container.NetworkSettings.Networks["docker_kapp"]
		for k, v := range container.NetworkSettings.Networks {
			logger.Debugf("key:%vx value:%vx", k, v)
		}

		address := ipConfig.IPAddress
		port := container.Ports[0].PublicPort
		if address != "" && port != 0 {
			url := fmt.Sprintf("http://%s:%d", "localhost", port)
			d := DockerContainer{
				ContainerName:    container.Names[0],
				ContainerAddress: url,
				ContainerID:      container.ID,
				ContainerPort:    fmt.Sprintf("%d", port),
			}
			logger.Debugf("append backend %+v", d)
			backendList = append(backendList, d)
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
		logger.Error("can't get container", err)
		return nil, err
	}

	var urls []string
	port := ctn.NetworkSettings.Ports
	if ctn.NetworkSettings.Networks["docker_kapp"].IPAddress == "" {
		logger.Error("can't get IP address")
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
