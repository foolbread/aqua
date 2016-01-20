//@auther foolbread
//@time 2015-12-01
//@file aqua/connect_server/server/client.go
package server

import (
	aerr "aqua/common/error"
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	"fmt"
	"net"

	"github.com/foolbread/fbcommon/golog"
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
		req, data, err := s.readData(buf[:])
		if err != nil {
			golog.Error(err)
			return
		}

		golog.Info("recv services from client:", s.Con.RemoteAddr().String(), "data_len:", len(data))

		ser = g_logicmanager.getServer(int(req.ServiceType), s.Token)
		if ser == nil {
			golog.Error(aerr.ErrNoLogicSvr)
			return
		}

		err = ser.send(data)
		if err != nil {
			golog.Error(err)
		}
	}
}

func (s *Client) readData(buf []byte) (*aproto.ServiceRequest, []byte, error) {
	n, _, err := anet.RecvPacket(s.Con, buf)
	if err != nil {
		return nil, nil, err
	}

	req, err := aproto.UnmarshalServiceReq(buf[aproto.HEAD_LEN:n])
	if err != nil {
		return nil, nil, err
	}

	return req, buf[:n], nil
}

func (s *Client) sendData(d []byte) error {
	golog.Info("send services  to [client]:", s.Con.RemoteAddr().String(), "[token]:", s.TokenStr, "[data_len]:", len(d))
	return anet.SendPacket(s.Con, d)
}
