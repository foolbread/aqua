//@auther foolbread
//@time 2015-11-12
//@file aqua/login_server/server/login_server.go
package server

import (
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	"aqua/login_server/config"
	"fbcommon/golog"
	"net"
)

type loginServer struct {
}

func (s *loginServer) startListen() {
	addr := config.GetConfig().GetListenAddr()

	li, err := net.Listen("tcp", addr)
	if err != nil {
		golog.Critical(err)
	}

	golog.Info("start listen:", addr)

	for {
		c, err := li.Accept()
		if err != nil {
			golog.Critical(err)
		}

		go s.handleConnect(c)
	}
}

func (s *loginServer) handleConnect(c net.Conn) {
	//TO DO
	var buf [1024]byte
	defer c.Close()
	//recv client data
	l, _, err := anet.RecvPacket(c, buf[:])
	if err != nil {
		golog.Error(err)
		return
	}
	//parse client data
	h := aproto.UnmarshalHead(buf[:l])
	switch h.Cmd {
	case aproto.LOGINREQ_CMD:
		req, err := aproto.UnmarshalLoginReq(buf[:l])
		if err != nil {
			golog.Error(err)
			return
		}
		s.handleLogin(c, req)
	}
}

func (s *loginServer) handleLogin(c net.Conn, req *aproto.LoginRequest) {
	//check token

	//get addr
	var addr string
	//get new token
	var token []byte
	//construct redirect
	d, err := aproto.MarshalRedirect(token, addr)
	if err != nil {
		golog.Error(err)
		return
	}

	err = anet.SendPacket(c, d)
	if err != nil {
		golog.Error(err)
	}

}
