package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"time"
)

var cli *client.Client

var Running []types.Container

func init() {
	var err error

	// Init docker cli
	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// Init running containers
	err = reload()
	if err != nil {
		panic(err)
	}

	// Listen for docker events
	go listen()
}

func main() {
	fmt.Print("listening docker engine")
	time.Sleep(10 * time.Minute)
}

func Ps() {
	fmt.Println("Running containers")
	for _, container := range Running {
		fmt.Printf("%s\n", container.ID)
	}
}

// Reload the list of running containers
func reload() error {
	var err error
	Running, err = cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: false,
	})
	if err != nil {
		return err
	}
	return nil
}

// Listen for docker events
func listen() {
	filter := filters.NewArgs()
	filter.Add("type", "container")
	filter.Add("event", "start")
	filter.Add("event", "die")

	msg, errChan := cli.Events(context.Background(), types.EventsOptions{
		Filters: filter,
	})

	for {
		select {
		case err := <-errChan:
			panic(err)
		case d := <-msg:
			jsond, _ := json.Marshal(d)
			fmt.Println(string(jsond))
			reload()
		}
	}
}
