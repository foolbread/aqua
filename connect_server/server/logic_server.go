//@auther foolbread
//@time 2015-12-01
//@file aqua/connect_server/server/logic_server.go
package server

import (
	aerr "aqua/common/error"
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	"net"
	"time"

	"github.com/foolbread/fbcommon/golog"

	fbatomic "github.com/foolbread/fbcommon/atomic"
	fbtime "github.com/foolbread/fbcommon/time"
)

var connect_timer *fbtime.Timer = fbtime.New(1 * time.Second)
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
	connect_timer.NewTimer(checkalive_time, s.checkKeepalive)

	for {
		var buf [4096]byte
		pa, err := s.read(buf[:])
		if err != nil {
			golog.Error(err)
			g_logicmanager.exitServer(s.Service_type, s.Addr)
			return
		}

		switch pa.cmd {
		case aproto.LOGICSERVICERES_CMD:
			g_logicmanager.handlerLogicToClient(pa)
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
	case aproto.LOGICSERVICERES_CMD:
		golog.Info("recv services from [logicserver]:", s.Addr, "[service_type]:", s.Service_type, "[data_len]:", n)
		res, err := aproto.UnmarshalServiceRes(buf[aproto.HEAD_LEN:n])
		if err != nil {
			return nil, err
		}

		return &logicPacket{res.Token, cmd, res.Sn, buf[:n]}, nil
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
	connect_timer.NewTimer(checkalive_time, s.checkKeepalive)
}
