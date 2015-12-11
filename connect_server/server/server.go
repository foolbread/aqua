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
	mcon map[string]*outerConn
}

func newOuterConnSet() *outerConnSet {
	r := new(outerConnSet)
	r.lo = new(sync.RWMutex)
	r.mcon = make(map[string]*outerConn)

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

		go s.handlerOuterCon(c)
	}
}

func (s *outerConnectServer) handlerOuterCon(c net.Conn) {
	defer c.Close()
	//handler login
	con, err := s.handlerLogin(c)
	if err != nil {
		golog.Error(err)
		return
	}

	con.handlerService()
}

func (s *outerConnectServer) handlerLogin(c net.Conn) (*outerConn, error) {
	var buf [1024]byte
	var ret *outerConn = nil
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
		ret = &outerConn{req.Cid, key, req.Token, c}
		set.mcon[key] = ret
	}
	set.lo.Unlock()

	//construct login response
	d, err := aproto.MarshalLoginRes(req.Token, LOGIN_SUCCESS, req.Cid)
	if err != nil {
		s.exitConn(req.Token, key)
		return nil, err
	}

	err = anet.SendPacket(c, d)
	if err != nil {
		s.exitConn(req.Token, key)
		return nil, err
	}

	return ret, nil
}

func (s *outerConnectServer) exitConn(key []byte, keystr string) {
	set := s.cons[key[0]%ARRARY_LEN]

	set.lo.Lock()
	delete(set.mcon, keystr)
	set.lo.Unlock()

	golog.Info("token:", keystr, "connect exit!")
}

////////////////////////////////////////////////////////////////////////////////////

type innerLogicServer struct {
	con net.Conn
}

type innerMsg struct {
	t int
	d []byte
}

type innerConnectServer struct {
	server map[int][]*innerLogicServer
	lo     *sync.RWMutex
	ch     chan *innerMsg
}

func newInnerConnectServer() *innerConnectServer {
	r := new(innerConnectServer)
	r.ch = make(chan *innerMsg, 1024)
	r.server = make(map[int][]*innerLogicServer)
	r.lo = new(sync.RWMutex)

	return r
}

func (s *innerConnectServer) startListen() {
	iaddr := config.GetConfig().GetInnerAddr()

	go s.innerServer(iaddr)
}

func (s *innerConnectServer) pushMsg(m *innerMsg) {
	s.lo.RLock()
	if _, ok := s.server[m.t]; ok {
	}
	s.lo.RUnlock()
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

		go s.handlerInnerCon(c)
	}
}

func (s *innerConnectServer) handlerInnerCon(c net.Conn) {

}
