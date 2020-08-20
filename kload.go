package main

import (
	"fmt"
	"net/http"

	golb "github.com/PhamDuyKhang/go-lb/internal"
)

func main() {
	listBackend := []string{
		"http://localhost:8082",
		"http://localhost:8083",
		"http://localhost:8084",
		"http://localhost:8085",
		"http://localhost:8086",
	}

	lbPool, err := golb.NewLoadBalaningPool(listBackend)
	if err != nil {
		panic(err)
	}
	KLB := golb.NewLoadBalancer(lbPool)

	mainServer := http.Server{
		Addr:    fmt.Sprintf(":8080"),
		Handler: http.HandlerFunc(KLB.LoadBalance),
	}

	err = mainServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
