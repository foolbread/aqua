//@auther foolbread
//@time 2015-12-01
//@file aqua/connect_server/server/inner_conn.go
package server

import (
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

}
