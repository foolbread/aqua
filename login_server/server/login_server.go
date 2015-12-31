//@auther foolbread
//@time 2015-11-12
//@file aqua/login_server/server/login_server.go
package server

import (
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	"aqua/login_server/config"
	"aqua/login_server/storage"
	"container/list"
	"crypto/md5"
	"fbcommon/golog"
	"fmt"
	"net"
	"sync"
)

var session_format string = "%02X_%d"

type redirectInfo struct {
	status int
	addr   string
	token  [16]byte
}

type loginServer struct {
	lo          *sync.RWMutex
	connect_svr *list.List
}

func newLoginServer() *loginServer {
	r := new(loginServer)
	r.lo = new(sync.RWMutex)
	r.connect_svr = list.New()

	return r
}

func (s *loginServer) startListen() {
	go s.startInnerListen()

	go s.startOuterListen()
}

func (s *loginServer) startOuterListen() {
	addr := config.GetConfig().GetOuterAddr()

	li, err := net.Listen("tcp", addr)
	if err != nil {
		golog.Critical(err)
	}

	golog.Info("start listen:", addr)

	for {
		c, err := li.Accept()
		if err != nil {
			golog.Critical(err)
		}

		go s.handlerConnect(c)
	}
}

func (s *loginServer) startInnerListen() {
	addr := config.GetConfig().GetInnerAddr()

	li, err := net.Listen("tcp", addr)
	if err != nil {
		golog.Critical(err)
	}

	golog.Info("start listen:", addr)

	for {
		c, err := li.Accept()
		if err != nil {
			golog.Critical(err)
		}

		go s.handlerRegister(c)
	}
}

func (s *loginServer) handlerRegister(c net.Conn) {
	var buf [1024]byte
	var Id uint32
	defer func() {
		c.Close()
		g_conmanager.exitConnectSvr(Id)
	}()

	//request
	l, _, err := anet.RecvPacket(c, buf[:])
	if err != nil {
		golog.Error(err)
		return
	}

	req, err := aproto.UnmarshalConnectRegisterReq(buf[aproto.HEAD_LEN:l])
	if err != nil {
		golog.Error(err)
		return
	}

	Id = req.Id

	//response
	da, err := aproto.MarshalConnectRegisterRes(aproto.STATUS_OK)
	if err != nil {
		golog.Error(err)
		return
	}

	err = anet.SendPacket(c, da)
	if err != nil {
		golog.Error(err)
		return
	}

	//register info
	csvr := newConnectServer(req.Id, req.ListenAddr, c)
	g_conmanager.addConnectSvr(csvr)
	golog.Info("id:", Id, "addr:", c.RemoteAddr().String(), "connect server register success!")

	csvr.run()
}

func (s *loginServer) handlerConnect(c net.Conn) {
	//TO DO
	var buf [1024]byte
	defer c.Close()
	//recv client data
	l, _, err := anet.RecvPacket(c, buf[:])
	if err != nil {
		golog.Error(err)
		return
	}

	//parse client data
	h := aproto.UnmarshalHead(buf[:aproto.HEAD_LEN])
	switch h.Cmd {
	case aproto.LOGINREQ_CMD:
		s.handlerLogin(c, buf[:l])
	default:
		golog.Info("unknow cmd:", h.Cmd)
	}
}

func (s *loginServer) handlerLogin(c net.Conn, d []byte) {
	//parse login
	req, err := aproto.UnmarshalLoginReq(d[aproto.HEAD_LEN:])
	if err != nil {
		golog.Error(err)
		return
	}

	handler := storage.GetStorage().GetSessionHandler(req.Cid)
	/*//check token
	b, err := handler.IsExistSession(req.Cid)
	if err != nil {
		golog.Error(err)
		return
	}*/

	var info redirectInfo
	info.status = aproto.STATUS_OK
	/*if b {
		info.status = aproto.ALREADY_LOGIN
		s.handlerRedirect(c, &info)
		return
	}*/

	//token
	info.token = md5.Sum([]byte(req.Cid + c.RemoteAddr().String()))
	//addr
	csvr := g_conmanager.getConnectSvr(info.token[:])
	if csvr == nil {
		info.status = aproto.SERVICE_ERROR
		s.handlerRedirect(c, &info)
		return
	}
	info.addr = csvr.Addr
	//store session
	session := fmt.Sprintf(session_format, info.token, csvr.Id)
	handler.SetUsrSession(req.Cid, session)
	golog.Info("redirect new client:", c.RemoteAddr().String(), "to:", csvr.Addr, "cid:", req.Cid, "session:", session)

	s.handlerRedirect(c, &info)
}

func (s *loginServer) handlerRedirect(c net.Conn, info *redirectInfo) {
	da, err := aproto.MarshalRedirect(info.status, info.token[:], info.addr)
	if err != nil {
		golog.Error(err)
		return
	}

	err = anet.SendPacket(c, da)
	if err != nil {
		golog.Error(err)
	}
}
