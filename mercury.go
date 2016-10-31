//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

//mecury is a server framwwork for go.
package mercury

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type Mercury struct {
	config  *MercuryConfig
	router  Router
	workers map[string]Worker
	timers  map[string]Timer
}

var (
	mercury Mercury
)

func init() {
	//set gp
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	//init config
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	mercury.config = NewConfig()
	//find conf in curr dir
	configFilePath := filepath.Join(workPath, CONFIG_FILE)
	if !FileExists(configFilePath) {
		configFilePath = filepath.Join(workPath, "conf", CONFIG_FILE)
		if !FileExists(configFilePath) {
			panic("can't find mercury.conf")
		}
	}
	err = mercury.config.Init(configFilePath)
	if err != nil {
		panic(err)
	}

	//init log
	setLogDir(mercury.config.LogDir)
	setLogLevel(mercury.config.LogLevel)

	//init stderr output
	if mercury.config.StdErr2File {
		stdErr2File()
	}

	//init workers
	mercury.workers = make(map[string]Worker)
	//init timers
	mercury.timers = make(map[string]Timer)
}

func Init() error {
	SetRouter(&DefaultRouter{})
	RegisterTimer("status", StatusTimer{})
	return nil
}

func Run() error {
	//start timer
	for name, timer := range mercury.timers {
		Info("start timer : %s", name)
		go func(t Timer) {
			cycle := time.Duration(uint64(t.GetCycle()) * 1000 * 1000)
			ticker := time.NewTicker(cycle)
			for _ = range ticker.C {
				err := t.Process()
				if err != nil {
					Warn("one cycle %s failed, err:%s", t.GetName(), err.Error())
				}
			}
		}(timer)
	}

	tcpAddr := fmt.Sprintf("%s:%d", mercury.config.Ip, mercury.config.Port)
	listener, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		Error("mercury listen failed, addr:%s, err:%s", tcpAddr, err.Error())
		return err
	}
	defer listener.Close()

	Info("mercury run...")
	//
	for {
		conn, err := listener.Accept()
		if err != nil {
			Warn("mercury accept failed, err:%s", err.Error())
			continue
		}
		Debug("new conn[%s]", conn.RemoteAddr().String())
		go mercury.handleConn(conn)
	}

	return nil
}

func (mercury *Mercury) handleConn(conn net.Conn) error {
	bufSize := mercury.config.RecvBuffSize
	buf := make([]byte, bufSize)
	var recved uint32 = 0
	defer func() {
		conn.Close()
		Debug("conn[%s] close", conn.RemoteAddr().String())
	}()

	needRecv := true
	var ret int
	var err error
	var ready bool
	var workerName string
	for {
		if needRecv {
			ret, err = RecvReq(conn, buf[recved:])
			if err != nil {
				return err
			} else if ret == 0 {
				continue
			}
			recved += uint32(ret)
		}

		//route
		ready, err = mercury.router.Ready(buf[0:recved])
		if err != nil {
			Warn("request ready route failed")
			return err
		} else if false == ready {
			needRecv = true
			continue
		}

		workerName, err = mercury.router.Route(buf[0:recved])
		if err != nil {
			Warn("request route failed")
			return err
		} else if workerName == "" {
			needRecv = true
			continue
		} else if _, ok := mercury.workers[workerName]; false == ok {
			Warn("invalid worker name[%s]", workerName)
			return ERR_INVALID_WORKER
		}
		Debug("router success, worker_name[%s]", workerName)

		//
		worker := mercury.workers[workerName]
		packSize, err := worker.Complete(buf[0:recved])
		if err != nil {
			Warn("check pack complete error, remote:%s, err:%s", conn.RemoteAddr().String(), err.Error())
			return err
		} else if packSize < 0 {
			Warn("pack complete size invalid, remote:%s, size:%d", conn.RemoteAddr().String(), packSize)
			return errors.New("pack complete size invalid")
		} else if packSize == 0 {
			//not complete
			if recved == bufSize {
				newBuf := make([]byte, bufSize*2)
				copy(newBuf, buf)
				buf = newBuf
				bufSize = bufSize * 2
			}
			needRecv = true
			continue
		}

		//process
		newMsg, err := worker.Receive(buf[0:packSize])
		if err != nil {
			Warn("recv msg failed, remote:%s, err:%s", conn.RemoteAddr().String(), err.Error())
		} else {
			go func() {
				err := worker.Process(conn, newMsg)
				if err != nil {
					Warn("logid[%d] worker process failed, err:%s", newMsg.GetReqId(), err.Error())
				} else {
					Debug("logid[%d] worker process success", newMsg.GetReqId())
				}
			}()
		}
		copy(buf, buf[packSize:recved])
		recved -= uint32(packSize)
		if recved > 0 {
			needRecv = false
		} else {
			needRecv = true
		}
	}
}

func stdErr2File() {
	procName := os.Args[0]
	procName = procName[strings.LastIndex(procName, "/")+1:]
	filePath := fmt.Sprintf("%s/%s_%d_%d.err", mercury.config.LogDir, procName, os.Getpid(), time.Now().Unix())
	f, err := os.Create(filePath)
	if err != nil {
		Error("create error file failed, file[%s]", filePath)
		return
	}
	err = syscall.Dup2(int(f.Fd()), 2)
	if err != nil {
		Error("dup2 stderr to file failed, file[%s], err[%s]", filePath, err.Error())
	}
}
