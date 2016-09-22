//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

package mercury

import (
	"fmt"
	"runtime"
)

const (
	DEFAULT_STATUS_CYCLE = 30 * 1000
)

type Timer interface {
	GetName() string
	GetCycle() uint32
	Process() error
}

func RegisterTimer(name string, timer Timer) error {
	if nil == timer {
		Error("register invalid timer, name:%s", name)
		err := fmt.Errorf("register invalid timer, name:%s", name)
		panic(err)
	}
	mercury.timers[name] = timer
	return nil
}

//内置的几个timer
//1. 状态timer
type StatusTimer struct {
}

func (timer StatusTimer) GetName() string {
	return "status timer"
}

func (timer StatusTimer) GetCycle() uint32 {
	return mercury.config.StatusCycle
}

func (timer StatusTimer) Process() error {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	self := memStat.Alloc
	Info("go routine : %d, mem use : %d", runtime.NumGoroutine(), self)
	return nil
}

//3. fake
type FakeTimer struct {
}

func (timer FakeTimer) GetName() string {
	return "fake"
}

func (timer FakeTimer) GetCycle() uint32 {
	return 1000
}

func (timer FakeTimer) Process() error {
	fmt.Println("this is fake")
	return nil
}
