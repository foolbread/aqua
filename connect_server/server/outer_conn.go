//@auther foolbread
//@time 2015-12-01
//@file aqua/connect_server/server/outer_conn.go
package server

import (
	"net"
)

type outerConn struct {
	Cid      string
	TokenStr string
	Token    []byte
	Con      net.Conn
}

func (s *outerConn) handlerService() {

}
