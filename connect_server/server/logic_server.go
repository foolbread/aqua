//@auther foolbread
//@time 2015-12-01
//@file aqua/connect_server/server/inner_conn.go
package server

import (
	aerr "aqua/common/error"
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	fbatomic "fbcommon/atomic"
	"fbcommon/golog"
	fbtime "fbcommon/time"
	"net"
	"time"
)

var logic_timer fbtime.Timer
var checkalive_time time.Duration = 15 * time.Second

type logicPacket struct {
	token []byte
	cmd   uint32
	sn    string
	data  []byte
}

var keepalive logicPacket

type logicServer struct {
	Service_type int
	Con          net.Conn
	Addr         string
	alive        fbatomic.AtomicBool
}

func newLogicServer(t int, c net.Conn) *logicServer {
	r := new(logicServer)
	r.Service_type = t
	r.Con = c
	r.Addr = c.RemoteAddr().String()
	r.alive.Set(true)

	return r
}

func (s *logicServer) run() {
	var buf [4096]byte

	logic_timer.NewTimer(checkalive_time, s.checkKeepalive)

	for {
		pa, err := s.read(buf[:])
		if err != nil {
			golog.Error(err)
			g_logicmanager.exitServer(s.Service_type, s.Addr)
			return
		}

		switch pa.cmd {
		case aproto.SERVICERES_CMD:
			g_logicmanager.handlerLogicRes(pa)
		case aproto.KEEPALIVE_CMD:
			s.alive.Set(true)
		}
	}
}

func (s *logicServer) read(buf []byte) (*logicPacket, error) {
	n, cmd, err := anet.RecvPacket(s.Con, buf)
	if err != nil {
		return nil, err
	}

	switch cmd {
	case aproto.SERVICERES_CMD:
		req, err := aproto.UnmarshalServiceReq(buf[aproto.HEAD_LEN:n])
		if err != nil {
			return nil, err
		}

		return &logicPacket{req.Token, cmd, req.Sn, buf[:n]}, nil
	case aproto.KEEPALIVE_CMD:
		return &keepalive, nil
	}

	return nil, aerr.ErrUnknowCmd
}

func (s *logicServer) send(d []byte) error {
	return anet.SendPacket(s.Con, d)
}

func (s *logicServer) checkKeepalive() {
	if !s.alive.Get() {
		s.Con.Close()
		return
	}

	s.alive.Set(false)
	logic_timer.NewTimer(checkalive_time, s.checkKeepalive)
}
