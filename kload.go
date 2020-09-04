package main

import (
	"flag"
	"fmt"
	golb "github.com/PhamDuyKhang/go-lb/internal"
	"github.com/PhamDuyKhang/go-lb/internal/dicovery"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("starting load balancing")

	np := flag.String("np", "", "the network name of docker provide")

	flag.Parse()

	if np == nil {
		logrus.Panic("can't get important parameter")
	}

	backendList := dicovery.GetListBackend(*np)

	lbPool, err := golb.NewLoadBalancingPool(backendList)
	if err != nil {
		logrus.Panic(err)
	}
	KLB := golb.NewLoadBalancer(lbPool)

	mainServer := http.Server{
		Addr:    fmt.Sprintf(":8080"),
		Handler: http.HandlerFunc(KLB.LoadBalance),
	}

	err = mainServer.ListenAndServe()
	if err != nil {
		logrus.Panic(err)
	}
}
