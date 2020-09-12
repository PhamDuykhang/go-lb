package datastructure

import (
	"github.com/PhamDuyKhang/go-lb/internal/services"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	DefaultCapacity = 10
	MaximumElement  = 50
)

type (

	//BackendRingBuffer is a structure help hold backend list
	//inside the structure we have mutex to prevent race condition
	BackendRingBuffer struct {
		mx         sync.Mutex
		a          []services.Backend
		numElement int64
		size       int64
		tail       int64
		head       int64
	}
)

//NewRing create the ring buffer with Backend interface is type of ring buffer
//every struct that implement backend interface can used
func NewRing(maximumLen int64) BackendRingBuffer {
	return BackendRingBuffer{
		mx:         sync.Mutex{},
		size:       maximumLen,
		a:          make([]services.Backend, maximumLen),
		numElement: 0,
		head:       0,
		tail:       0,
	}
}

//EnQueue add a backend to buffer
func (r *BackendRingBuffer) EnQueues(b services.Backend) {
	if !r.IsFull() {
		r.mx.Lock()
		defer r.mx.Unlock()
		t := (r.tail + 1) % int64(len(r.a))
		r.a[t] = b
		if t == 0 {
			r.tail = t
		} else {
			r.tail++
		}
		r.numElement++
		return
	}
}

//DeQueue take the backend from buffer
func (r *BackendRingBuffer) DeQueue() services.Backend {
	if !r.IsEmpty() {
		r.mx.Lock()
		defer r.mx.Unlock()
		h := (r.head + 1) % int64(len(r.a))
		element := r.a[h]
		r.head = h
		r.numElement--
		return element
	}
	return nil
}

func (r *BackendRingBuffer) IsFull() bool {
	r.mx.Lock()
	defer r.mx.Unlock()
	if r.numElement == r.size {
		return true
	}
	return false
}

func (r *BackendRingBuffer) Len() int {
	r.mx.Lock()
	defer r.mx.Unlock()
	return int(r.numElement)
}

func (r *BackendRingBuffer) IsEmpty() bool {
	r.mx.Lock()
	defer r.mx.Unlock()
	if r.numElement <= 0 {
		return true
	}
	return false
}

func (r *BackendRingBuffer) helperGetRing() BackendRingBuffer {
	logrus.Infof("backend data %+v", r)
	return BackendRingBuffer{
		numElement: r.numElement,
		head:       r.head,
		tail:       r.tail,
		size:       r.size,
		a:          r.a,
	}
}
