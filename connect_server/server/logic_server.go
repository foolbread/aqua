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
	"net"
)

type logicPacket struct {
	token []byte
	cmd   uint32
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
	for {
		pa, err := s.read(buf[:])
		if err != nil {
			golog.Error(err)
			continue
		}

		switch pa.cmd {
		case aproto.SERVICERES_CMD:
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

		return &logicPacket{req.Token, cmd, buf[:n]}, nil
	case aproto.KEEPALIVE_CMD:
		return &keepalive, nil
	}

	return nil, aerr.ErrUnknowCmd
}

func (s *logicServer) send(d []byte) error {
	return anet.SendPacket(s.Con, d)
}

func (s *logicServer) checkKeepalive() {

}
