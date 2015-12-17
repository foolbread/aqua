//@auther foolbread
//@time 2015-12-01
//@file aqua/connect_server/server/outer_conn.go
package server

import (
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	"fbcommon/golog"
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

func (s *Client) run() {
	var buf [4096]byte
	var ser *logicServer
	for {
		req, err := s.readRequest(buf[:])
		if err != nil {
			golog.Error(err)
			return
		}

		ser = g_innerserver.getServer(int(req.ServiceType), s.Token)
		if ser == nil {
			return
		}
	}
}

func (s *Client) readRequest(buf []byte) (*aproto.ServiceRequest, error) {
	n, err := anet.RecvPacket(s.Con, buf)
	if err != nil {
		return nil, err
	}

	req, err := aproto.UnmarshalServiceReq(buf[:n])
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *Client) sendResponse(d []byte) error {

	return anet.SendPacket(s.Con, d)
}
