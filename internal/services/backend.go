package services

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	HealthCheck = "/check"
	StatusDown  = "down"
	StatusUp    = "up"
)

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)

type (
	//BackendCommonInformation hold a basic information
	BackendCommonInformation struct {
		URL  string
		Port string
		ID   string
		Name string
	}
	//Backend is a set of function we need to get net instance
	Backend interface {
		IsAlive() bool
		HealthCheck() string
		SetHealth(isDeed bool)
		Create() error
		GetID() string
		Stat() BackendCommonInformation
		ErrorHandle(halerFunc ErrorHandlerFunc)
		Serve(w http.ResponseWriter, r *http.Request)
	}
	//DockerEnvContainer for docker env
	DockerEnvContainer struct {
		ContainerID   string
		ContainerName string
		url           string
		Port          string
		prx           *httputil.ReverseProxy
		State         string
		isAlive       bool
		weigh         int
	}
)

//NewDockerEnvContainer if docker env
func NewDockerEnvContainer(url string, id string, name string) *DockerEnvContainer {
	return &DockerEnvContainer{
		url:           url,
		ContainerID:   id,
		ContainerName: name,
	}
}

//IsAlive check the status of container
func (dc *DockerEnvContainer) Create() error {
	address, err := url.Parse(dc.url)
	if err != nil {
		return err
	}
	if dc.HealthCheck() == StatusUp {
		prx := httputil.NewSingleHostReverseProxy(address)
		dc.prx = prx
		dc.SetHealth(false)
	}
	return nil

}

//IsAlive check the status of container
func (dc *DockerEnvContainer) IsAlive() bool {
	return dc.isAlive
}

//SetHealth is used for health checking
func (dc *DockerEnvContainer) SetHealth(isDeed bool) {
	if isDeed {
		dc.isAlive = false
		return
	}
	dc.isAlive = true
}

//GetID get a unique identify of service instance
func (dc *DockerEnvContainer) GetID() string {
	return dc.ContainerID
}

//Stat get the information of container
func (dc *DockerEnvContainer) Stat() BackendCommonInformation {
	return BackendCommonInformation{
		URL:  dc.url,
		Port: dc.Port,
		Name: dc.ContainerName,
		ID:   dc.ContainerID,
	}
}

//ErrorHandle set a callback func if the reverse proxy got error
func (dc *DockerEnvContainer) ErrorHandle(errorHandlerFunc ErrorHandlerFunc) {
	dc.prx.ErrorHandler = errorHandlerFunc
}

//Serve forward the request in to container
func (dc *DockerEnvContainer) Serve(w http.ResponseWriter, r *http.Request) {
	dc.prx.ServeHTTP(w, r)
}

//HealthCheck check the service was up or not
func (dc *DockerEnvContainer) HealthCheck() string {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", dc.url, HealthCheck), nil)

	if err != nil {
		return ""
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return StatusDown
	}
	if res.StatusCode != http.StatusOK {
		return StatusDown
	}
	return StatusUp

}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareTwo")
		if r.URL.Path == "/foo" {
			return
		}

		next.ServeHTTP(w, r)
		log.Println("Executing middlewareTwo again")
	})
}
