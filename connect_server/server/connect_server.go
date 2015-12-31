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
	"time"
)

const keyformat string = "%02X"

var keepalive_time time.Duration = 10 * time.Second

type connectServer struct {
	id      uint32
	con     net.Conn
	addr    string
	clients [ARRARY_LEN]*clientManager
}

func (s *connectServer) startListen() {
	oaddr := config.GetConfig().GetOuterAddr()

	go s.server(oaddr)
}

func (s *connectServer) startRegister() {
	addr := config.GetConfig().GetLoginAddr()

	go s.register(addr)
}

func (s *connectServer) register(a string) {
	var err error
	for {
		s.con, err = net.Dial("tcp", a)
		if err != nil {
			golog.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}

	data, err := aproto.MarshalConnectRegisterReq(s.id, s.addr)
	if err != nil {
		golog.Critical(err)
	}

	err = anet.SendPacket(s.con, data)
	if err != nil {
		golog.Error(err)
		s.con.Close()
		s.startRegister()
		return
	}

	var buf [1024]byte
	l, _, err := anet.RecvPacket(s.con, buf[:])
	if err != nil {
		golog.Error(err)
		s.con.Close()
		s.startRegister()
		return
	}

	res, err := aproto.UnmarshalConnectRegisterRes(buf[aproto.HEAD_LEN:l])
	if err != nil {
		golog.Critical(err)
	}

	if res.Status == aproto.STATUS_OK {
		golog.Info("connect to login server:", a, "success!")
		go s.keepalive()
	}
}

func (s *connectServer) keepalive() {
	err := anet.SendPacket(s.con, aproto.KeepAlive[:])
	if err != nil {
		golog.Error(err)
		s.con.Close()
		s.startRegister()
		return
	}

	connect_timer.NewTimer(keepalive_time, s.keepalive)
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
	s.exitClient(cli.Token)
}

func (s *connectServer) handlerClientLogin(c net.Conn) (*Client, error) {
	var buf [1024]byte
	var ret *Client = nil
	l, _, err := anet.RecvPacketEx(c, buf[:], default_timeout)
	if err != nil {
		return nil, err
	}

	req, err := aproto.UnmarshalLoginReq(buf[aproto.HEAD_LEN:l])
	if err != nil {
		return nil, err
	}

	//check token
	set := s.clients[req.Token[0]%ARRARY_LEN]
	key := fmt.Sprintf(keyformat, req.Token)
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
		s.exitClient(req.Token)
		return nil, err
	}

	err = anet.SendPacket(c, d)
	if err != nil {
		s.exitClient(req.Token)
		return nil, err
	}

	return ret, nil
}

func (s *connectServer) handlerClientRes(pa *logicPacket) {
	clis := s.clients[pa.token[0]%ARRARY_LEN]

	clis.pushResponse(pa)
}

func (s *connectServer) exitClient(key []byte) {
	clis := s.clients[key[0]%ARRARY_LEN]
	keystr := fmt.Sprintf(keyformat, key)
	clis.exitClient(keystr)
}
