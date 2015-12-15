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
}

func newLogicServer(t int, c net.Conn) *logicServer {
	r := new(logicServer)
	r.Service_type = t
	r.Con = c

	return r
}

func (s *logicServer) run() {

}
