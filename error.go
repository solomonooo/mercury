//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

package mercury

import (
	"errors"
)

const (
	CODE_SUCCESS = iota
	CODE_INVALID_PARAM
	CODE_ROUTER_FAILED
	CODE_INVALID_WORKER
	CODE_CLIENT_CONN  = 100
	CODE_CLIENT_READ  = 101
	CODE_CLIENT_WRITE = 102
	CODE_UNKNOWN      = 9999
)

var (
	ERR_INVALID_PARAM  = errors.New("invalid param")
	ERR_ROUTER_FAILLD  = errors.New("request route failed")
	ERR_INVALID_WORKER = errors.New("invalid worker")
	ERR_CLIENT_CONN    = errors.New("client connect failed")
	ERR_CLIENT_READ    = errors.New("client read failed")
	ERR_CLIENT_WRITE   = errors.New("client write failed")
	ERR_UNKNOWN        = errors.New("unknown error")
)

var (
	codeList map[error]int32
	errList  map[int32]error
)

func init() {
	errList = make(map[int32]error)
	errList[CODE_INVALID_PARAM] = ERR_INVALID_PARAM
	errList[CODE_ROUTER_FAILED] = ERR_ROUTER_FAILLD
	errList[CODE_INVALID_WORKER] = ERR_INVALID_WORKER
	errList[CODE_CLIENT_CONN] = ERR_CLIENT_CONN
	errList[CODE_CLIENT_READ] = ERR_CLIENT_READ
	errList[CODE_CLIENT_WRITE] = ERR_CLIENT_WRITE
	errList[CODE_UNKNOWN] = ERR_UNKNOWN

	codeList = make(map[error]int32)
	for k, v := range errList {
		codeList[v] = k
	}
}

func Code2Error(code int32) error {
	if code == 0 {
		return nil
	}
	if v, ok := errList[code]; ok {
		return v
	}
	return ERR_UNKNOWN
}

func Error2Code(err error) int32 {
	if err == nil {
		return CODE_SUCCESS
	}
	if v, ok := codeList[err]; ok {
		return v
	}
	return CODE_UNKNOWN
}
