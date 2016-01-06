//@auther foolbread
//@time 2016-01-04
//@file aqua/singlechat_server/server/singlechat_server.go
package server

import (
	aproto "aqua/common/proto"
	"fbcommon/golog"
)

type singlechatServer struct {
}

func newSinglechatServer() *singlechatServer {
	r := new(singlechatServer)

	return r
}

func (s *singlechatServer) handlerSendMsgReq(d []byte) {
	req, err := aproto.UnmarshalSendPMsgReq(d)
	if err != nil {
		golog.Error(err)
		return
	}

}
