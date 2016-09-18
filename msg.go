//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

package mercury

type Msg interface {
	GetReqId() uint64
}

type DefaultMsg struct {
	reqid uint64
}

func (msg DefaultMsg) GetReqId() uint64 {
	return msg.reqid
}
