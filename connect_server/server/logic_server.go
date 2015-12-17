//@auther foolbread
//@time 2015-12-01
//@file aqua/connect_server/server/inner_conn.go
package server

import (
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	"net"
)

type logicServer struct {
	Service_type int
	Con          net.Conn
	Addr         string
}

func newLogicServer(t int, c net.Conn) *logicServer {
	r := new(logicServer)
	r.Service_type = t
	r.Con = c
	r.Addr = c.RemoteAddr().String()

	return r
}

func (s *logicServer) run() {
	var buf [4096]byte
	for {
		s.readResponse(buf[:])
	}
}

func (s *logicServer) readResponse(buf []byte) (*aproto.ServiceResponse, uint32, error) {
	n, err := anet.RecvPacket(s.Con, buf)
	if err != nil {
		return nil, 0, err
	}

	res, err := aproto.UnmarshalServiceRes(buf[aproto.HEAD_LEN:n])
	if err != nil {
		return nil, 0, err
	}

	return res, n, nil
}

func (s *logicServer) sendRequest(d []byte) error {
	return anet.SendPacket(s.Con, d)
}
