//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

package mercury

import (
	"io"
	"net"
	"os"
	"time"
)

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func RecvReq(conn net.Conn, data []byte) (int, error) {
	timeout := time.Duration(mercury.config.RecvTimeout * 1000 * 1000)
	conn.SetReadDeadline(time.Now().Add(timeout))
	ret, err := conn.Read(data)
	if err == io.EOF {
		Debug("client close, remote[%s]", conn.RemoteAddr().String())
	} else if err != nil {
		Warn("recv req error, remote:%s, err:%s", conn.RemoteAddr().String(), err.Error())
		return ret, err
	}
	if ret != 0 {
		Debug("recv req success, remote:%s size:%d", conn.RemoteAddr().String(), ret)
	}
	return ret, err
}

func SendRsp(conn net.Conn, reqid uint64, data []byte) (int, error) {
	timeout := time.Duration(mercury.config.SendTimeout * 1000 * 1000)
	conn.SetWriteDeadline(time.Now().Add(timeout))
	ret, err := conn.Write(data)
	if err != nil {
		Warn("logid[%d] send rsp failed, remote[%s], err[%s]", reqid, conn.RemoteAddr().String(), err.Error())
		return ret, err
	}
	Debug("logid[%d] send rsp success, remote[%s], size[%d]", reqid, conn.RemoteAddr().String(), ret)
	return ret, err
}
