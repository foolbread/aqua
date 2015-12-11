//@auther foolbread
//@time 2015-12-01
//@file aqua/connect_server/server/outer_conn.go
package server

import (
	"fmt"
	"net"
)

type Client struct {
	Cid      string
	TokenStr string
	Token    []byte
	Con      net.Conn
}

func newClient(cid string, token []byte, con net.Conn) *Client {
	r := new(Client)
	r.Cid = cid
	r.Token = token
	r.TokenStr = fmt.Sprintf("%02X", token)
	r.Con = con

	return r
}

func (s *Client) handlerService() {

}
