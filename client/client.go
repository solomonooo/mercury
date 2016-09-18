//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

package client

type Client interface {
	GetCost() uint32
	GetLastInfo() (string, uint32)
}
