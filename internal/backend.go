package golb

import (
	"net/http"
	"net/http/httputil"
	"sync"
)

type (
	//Backend is a service instand when it is scaled
	Backend struct {
		TagertURL string
		mu        sync.RWMutex
		rpx       *httputil.ReverseProxy
		isAlive   bool
	}
)

//SetDead mark a backend service status to down
func (b *Backend) SetDead() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.isAlive = false
}

//IsAlive to check the backend is a alive or not
func (b *Backend) IsAlive() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.isAlive

}

//HeathCheck call a special endpoint for checking the backend donw or on air
func (b *Backend) HeathCheck(hcPatterm string) {
	req, err := http.NewRequest(http.MethodGet, b.TagertURL+"/"+hcPatterm, nil)
	if err != nil {
		b.SetDead()
		return
	}
	cl := http.DefaultClient
	rs, err := cl.Do(req)
	if err != nil {
		b.SetDead()
		return
	}
	if rs.StatusCode != http.StatusOK {
		b.SetDead()
	}
}
