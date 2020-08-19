package main

import (
	"fmt"
	"net/http"

	golb "github.com/PhamDuyKhang/go-lb/internal"
)

func main() {
	listBackend := []string{
		"localhost:8081",
		"localhost:8082",
		"localhost:8083",
		"localhost:8084",
		"localhost:8085",
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
