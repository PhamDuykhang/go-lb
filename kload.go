package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PhamDuyKhang/go-lb/internal/discovery"
	"github.com/PhamDuyKhang/go-lb/internal/pool"
	"github.com/PhamDuyKhang/go-lb/internal/services"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})

}

func main() {
	logrus.Infof("starting load balancing at %s", time.Now().Format(time.ANSIC))

	listBackend := discovery.GetListBackend("docker_kapp")
	logrus.Infof("%d available backend", len(listBackend))

	loadBalancer := pool.NewRoundRobinStrategies()

	var servicesList []services.Backend

	for _, container := range listBackend {
		bk := services.NewDockerEnvContainer(container.ContainerAddress, container.ContainerID, container.ContainerName)
		if err := bk.Create(); err != nil {
			continue
		}
		servicesList = append(servicesList, bk)
	}
	logrus.Infof("staring init backend list")
	err := loadBalancer.InitBackend(servicesList)
	if err != nil {
		logrus.Panic("can't init backend", err)
	}

	logrus.Infof("staring http server")

	mainServer := http.Server{
		Addr:    fmt.Sprintf(":8080"),
		Handler: http.HandlerFunc(loadBalancer.LoadBalancing),
	}

	logrus.Info("starting load balancing service")
	go func() {
		err = mainServer.ListenAndServe()
		if err != nil {
			logrus.Panic(err)
		}
	}()
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	logrus.Info("service has been started")
	<-killSignal
	logrus.Info("service has been stopped")

}
