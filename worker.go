//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

package mercury

import (
	"fmt"
	"net"
)

type Worker interface {
	//
	Complete(buf []byte) (int, error)
	Receive(buf []byte) (Msg, error)
	Process(conn net.Conn, msg Msg) error
}

func RegisterWorker(name string, worker Worker) error {
	if nil == worker {
		Error("register invalid worker, name:%s", name)
		err := fmt.Errorf("register invalid worker, name:%s", name)
		panic(err)
	}
	mercury.workers[name] = worker
	return nil
}

type DefaultWorker struct {
}

func (worker *DefaultWorker) Complete(buf []byte) (int, error) {
	return len(buf), nil
}

func (worker *DefaultWorker) Receive(buf []byte) (Msg, error) {
	return DefaultMsg{}, nil
}

func (worker *DefaultWorker) Process(conn net.Conn, msg Msg) error {
	Info("one msg process, hello world")
	return nil
}
