//@auther foolbread
//@time 2015-11-12
//@file aqua/connect_server/server/connect_server.go
package server

import (
	aeer "aqua/common/error"
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	"aqua/connect_server/config"
	"fbcommon/golog"
	"fmt"
	"net"
)

type connectServer struct {
	cons [ARRARY_LEN]*clientManager
}

func (s *connectServer) startListen() {
	oaddr := config.GetConfig().GetOuterAddr()

	go s.server(oaddr)
}

func (s *connectServer) server(a string) {
	li, err := net.Listen("tcp", a)
	if err != nil {
		golog.Critical(err)
	}

	golog.Info("start outer listen:", a)

	for {
		c, err := li.Accept()
		if err != nil {
			golog.Critical(err)
		}

		go s.handlerClientCon(c)
	}
}

func (s *connectServer) handlerClientCon(c net.Conn) {
	defer c.Close()
	//handler login
	cli, err := s.handlerClientLogin(c)
	if err != nil {
		golog.Error(err)
		return
	}

	cli.run()

	//exit client
	s.exitClient(cli.Token, cli.TokenStr)
}

func (s *connectServer) handlerClientLogin(c net.Conn) (*Client, error) {
	var buf [1024]byte
	var ret *Client = nil
	l, err := anet.RecvPacketEx(c, buf[:], default_timeout)
	if err != nil {
		return nil, err
	}

	req, err := aproto.UnmarshalLoginReq(buf[aproto.HEAD_LEN:l])
	if err != nil {
		return nil, err
	}

	//check token
	set := s.cons[req.Token[0]%ARRARY_LEN]
	key := fmt.Sprintf("%02X", req.Token)
	set.lo.RLock()
	_, ok := set.mclient[key]
	if ok {
		set.lo.RUnlock()
		return nil, aeer.ErrConnExsit
	}
	set.lo.RUnlock()

	set.lo.Lock()
	//double check
	_, ok = set.mclient[key]
	if !ok {
		ret = newClient(req.Cid, req.Token, c)
		set.mclient[key] = ret
	} else {
		set.lo.Unlock()
		return nil, aeer.ErrConnExsit
	}
	set.lo.Unlock()

	//construct login response
	d, err := aproto.MarshalLoginRes(req.Token, aproto.STATUS_OK, req.Cid)
	if err != nil {
		s.exitClient(req.Token, key)
		return nil, err
	}

	err = anet.SendPacket(c, d)
	if err != nil {
		s.exitClient(req.Token, key)
		return nil, err
	}

	return ret, nil
}

func (s *connectServer) exitClient(key []byte, keystr string) {
	set := s.cons[key[0]%ARRARY_LEN]

	set.lo.Lock()
	delete(set.mclient, keystr)
	set.lo.Unlock()

	golog.Info("token:", keystr, "connect exit!")
}
