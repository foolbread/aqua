//@auther foolbread
//@time 2015-12-01
//@file aqua/connect_server/server/inner_conn.go
package server

import (
	"net"
)

type innerConn struct {
	Service_type int
	Con          net.Conn
}
