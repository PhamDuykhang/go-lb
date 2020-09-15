package datastructure

//
//import (
//	"github.com/sirupsen/logrus"
//	"testing"
//)
//
//var ring BackendRingBuffer
//
//func Test_RunRing(t *testing.T){
//	t.Parallel()
//
//	r:= NewRing(5)
//	ring = r
//
//	st1:= "Hello 1"
//	st2:= "Hello 2"
//	st3:= "Hello 3"
//	st4:= "Hello 4"
//	st5:= "Hello 5"
//
//	//st6:= "hello 17"
//
//	ring.EnQueues(st1)
//	ring.EnQueues(st2)
//	ring.EnQueues(st3)
//	ring.EnQueues(st4)
//	ring.EnQueues(st5)
//
//	logrus.Info(ring.DeQueue())
//	logrus.Info(ring.DeQueue())
//	logrus.Info(ring.DeQueue())
//	logrus.Info(ring.DeQueue())
//	logrus.Info(ring.DeQueue())
//	logrus.Info(ring.DeQueue())
//
//	t.Log()
//}
