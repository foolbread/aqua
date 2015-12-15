//@auther foolbread
//@time 2015-11-12
//@file aqua/connect_server/server/server.go
package server

import (
	aeer "aqua/common/error"
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	"aqua/connect_server/config"
	"fbcommon/golog"
	"fmt"
	"net"
	"sync"
)

const ARRARY_LEN = 255
const (
	LOGIN_SUCCESS = 0
)

func InitServer() {
	golog.Info("initing connect server......")
	g_outserver = new(outerConnectServer)
	for i := 0; i < ARRARY_LEN; i++ {
		g_outserver.cons[i] = newOuterConnSet()
	}
	g_outserver.startListen()

	g_innerserver = new(innerConnectServer)
	g_innerserver.startListen()
}

var g_outserver *outerConnectServer
var g_innerserver *innerConnectServer

type outerConnSet struct {
	lo   *sync.RWMutex
	mcon map[string]*Client
}

func newOuterConnSet() *outerConnSet {
	r := new(outerConnSet)
	r.lo = new(sync.RWMutex)
	r.mcon = make(map[string]*Client)

	return r
}

type outerConnectServer struct {
	cons [ARRARY_LEN]*outerConnSet
}

func (s *outerConnectServer) startListen() {
	oaddr := config.GetConfig().GetOuterAddr()

	go s.outerServer(oaddr)
}

func (s *outerConnectServer) outerServer(a string) {
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

func (s *outerConnectServer) handlerClientCon(c net.Conn) {
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

func (s *outerConnectServer) handlerClientLogin(c net.Conn) (*Client, error) {
	var buf [1024]byte
	var ret *Client = nil
	l, err := anet.RecvPacket(c, buf[:])
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
	_, ok := set.mcon[key]
	if ok {
		set.lo.RUnlock()
		return nil, aeer.ErrConnExsit
	}
	set.lo.RUnlock()

	set.lo.Lock()
	_, ok = set.mcon[key]
	if !ok {
		ret = newClient(req.Cid, req.Token, c)
		set.mcon[key] = ret
	}
	set.lo.Unlock()

	//construct login response
	d, err := aproto.MarshalLoginRes(req.Token, LOGIN_SUCCESS, req.Cid)
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

func (s *outerConnectServer) exitClient(key []byte, keystr string) {
	set := s.cons[key[0]%ARRARY_LEN]

	set.lo.Lock()
	delete(set.mcon, keystr)
	set.lo.Unlock()

	golog.Info("token:", keystr, "connect exit!")
}

////////////////////////////////////////////////////////////////////////////////////
type innerConnectServer struct {
	server map[int][]*logicServer
	lo     *sync.RWMutex
}

func newInnerConnectServer() *innerConnectServer {
	r := new(innerConnectServer)
	r.server = make(map[int][]*logicServer)
	r.lo = new(sync.RWMutex)

	return r
}

func (s *innerConnectServer) startListen() {
	iaddr := config.GetConfig().GetInnerAddr()

	go s.innerServer(iaddr)
}

func (s *innerConnectServer) getServer(t int, key []byte) *logicServer {
	var ret *logicServer
	s.lo.RLock()
	v, ok := s.server[t]
	if ok {
		ret = v[int(key[0])%len(v)]
	}
	s.lo.RUnlock()

	return ret
}

func (s *innerConnectServer) innerServer(a string) {
	li, err := net.Listen("tcp", a)
	if err != nil {
		golog.Critical(err)
	}

	golog.Info("start inner listen:", a)

	for {
		c, err := li.Accept()
		if err != nil {
			golog.Critical(err)
		}

		go s.handlerLogicCon(c)
	}
}

func (s *innerConnectServer) handlerLogicCon(c net.Conn) {
	ser, err := s.handlerLogicLogin(c)
	if err != nil {
		golog.Error(err)
		return
	}

	ser.run()

	s.exitServer(ser.Service_type, ser.Addr)
}

func (s *innerConnectServer) handlerLogicLogin(c net.Conn) (*logicServer, error) {
	var buf [1024]byte
	n, err := anet.RecvPacket(c, buf[:])
	if err != nil {
		return nil, err
	}

	req, err := aproto.UnmarshalSvrRegisterReq(buf[:n])
	if err != nil {
		return nil, err
	}

	l := newLogicServer(int(req.ServiceType), c)

	s.lo.Lock()
	val := s.server[int(req.ServiceType)]
	val = append(val, l)
	s.lo.Unlock()

	return l, nil
}

func (s *innerConnectServer) exitServer(t int, a string) {
	s.lo.Lock()
	ser, ok := s.server[t]
	if ok {
		for k, v := range ser {
			if a == v.Addr {
				if k == 0 {
					ser = ser[1:]
				} else if k == len(ser)-1 {
					ser = ser[:len(ser)-1]
				} else {
					t := ser[k+1:]
					ser = ser[:k]
					ser = append(ser, t...)
				}

				s.server[t] = ser
				golog.Info("service type:", t, "service addr:", a, "exit...")
				break
			}
		}
	}
	s.lo.Unlock()
}
