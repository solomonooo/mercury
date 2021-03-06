//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

package client

import (
	"fmt"
	"github.com/solomonooo/mercury"
	"gopkg.in/fatih/pool.v2"
	"net"
	"time"
)

const (
	BUFF_SIZE = 32 * 1024
)

type TcpClient struct {
	conn    net.Conn
	buff    []byte
	recved  uint32
	timeout uint32

	Reqid uint64
	ip    string
	port  uint32
	cost  uint32
}

func (client *TcpClient) Init(Reqid uint64, ip string, port uint32, timeout uint32) error {
	var err error
	client.Reqid = Reqid
	client.timeout = timeout
	client.ip = ip
	client.port = port
	addr := fmt.Sprintf("%s:%d", client.ip, client.port)
	client.conn, err = net.Dial("tcp", addr)
	if err != nil {
		mercury.Warn("logid[%d] connect to server failed, addr[%s], err[%s]", Reqid, addr, err.Error())
		return mercury.ERR_CLIENT_CONN
	}
	client.buff = make([]byte, BUFF_SIZE)
	client.recved = 0
	return nil
}

func (client *TcpClient) InitWithPool(Reqid uint64, p pool.Pool, timeout uint32) error {
	var err error
	client.Reqid = Reqid
	client.timeout = timeout
	client.conn, err = p.Get()
	if err != nil {
		mercury.Warn("logid[%d] connect to server failed, err[%s]", Reqid, err.Error())
		return mercury.ERR_CLIENT_CONN
	}
	client.buff = make([]byte, BUFF_SIZE)
	client.recved = 0
	return nil
}

func (client *TcpClient) MarkUnusable() {
	if p, ok := client.conn.(*pool.PoolConn); ok {
		p.MarkUnusable()
	}
}

func (client *TcpClient) Close() {
	client.conn.Close()
}

func (client *TcpClient) Send(req []byte, check func([]byte) (int, error)) ([]byte, error) {
	var err error
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		client.cost = uint32((end - start) / 1000000)
	}()

	timeout := time.Duration(uint64(client.timeout) * 1000 * 1000)
	client.conn.SetWriteDeadline(time.Now().Add(timeout))
	client.conn.SetReadDeadline(time.Now().Add(timeout))

	mercury.Debug("logid[%d] timeout[%v] time[%s]", client.Reqid, timeout, time.Now().Add(timeout).Format("2006-01-02 15:04:05"))

	//send
	ret, err := client.conn.Write(req)
	if nil != err {
		mercury.Warn("logid[%d] write data to server failed, err[%s]", client.Reqid, err.Error())
		return nil, mercury.ERR_CLIENT_WRITE
	} else if ret != len(req) {
		mercury.Warn("logid[%d] write data error, real[%d], need[%d]", ret, len(req))
		return nil, mercury.ERR_CLIENT_WRITE
	}

	//recv
	needRecv := true
	var packSize int = 0
	for {
		if needRecv {
			ret, err = client.conn.Read(client.buff[client.recved:])
			if nil != err {
				mercury.Warn("logid[%d] read data from server failed, err[%s]", client.Reqid, err.Error())
				return nil, mercury.ERR_CLIENT_READ
			} else if ret == 0 {
				continue
			}
			client.recved += uint32(ret)
		}

		packSize, err = check(client.buff[0:client.recved])
		if err != nil {
			mercury.Warn("logid[%d] check rsp from server failed, err[%s]", client.Reqid, err.Error())
			return nil, mercury.ERR_CLIENT_READ
		} else if packSize < 0 {
			mercury.Warn("logid[%d] check rsp from server failed, size[%d]", client.Reqid, packSize)
			return nil, mercury.ERR_CLIENT_READ
		} else if packSize == 0 {
			mercury.Debug("logid[%d] rsp from server incomplete, recved[%d]", client.Reqid, client.recved)
			if client.recved == uint32(len(client.buff)) {
				newbuff := make([]byte, len(client.buff)*2)
				copy(newbuff, client.buff)
				client.buff = newbuff
			}
			needRecv = true
			continue
		}
		break
	}
	rsp := make([]byte, packSize)
	copy(rsp, client.buff[0:packSize])
	copy(client.buff, client.buff[packSize:client.recved])
	client.recved -= uint32(packSize)
	return rsp, nil
}

func (client *TcpClient) GetCost() uint32 {
	return client.cost
}
