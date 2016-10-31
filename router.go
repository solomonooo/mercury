//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

package mercury

import (
	"fmt"
)

type Router interface {
	//check if can route
	Ready(buf []byte) (bool, error)
	Route(buf []byte) (string, error)
}

func SetRouter(router Router) error {
	if nil == router {
		Error("register invalid router")
		err := fmt.Errorf("register invalid router")
		panic(err)
	}
	mercury.router = router
	return nil
}

//默认router, 从当前worker中遍历选择一个worker执行
type DefaultRouter struct {
}

func (router *DefaultRouter) Ready(buf []byte) (bool, error) {
	return true, nil
}

func (router *DefaultRouter) Route(buf []byte) (string, error) {
	for name, _ := range mercury.workers {
		return name, nil
	}
	return "", ERR_ROUTER_FAILLD
}
