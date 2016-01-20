//@auther foolbread
//@time 2016-01-04
//@file aqua/singlechat_server/server/connect_server.go
package server

import (
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	"aqua/singlechat_server/config"
	"net"
	"time"

	"github.com/foolbread/fbcommon/golog"
)

var keepalive_time time.Duration = 10 * time.Second

type connectServer struct {
	addr string
	con  net.Conn
	id   uint32
}

func newConnectServer(a string) *connectServer {
	r := new(connectServer)
	r.addr = a

	go r.register()

	return r
}

func (s *connectServer) register() {
	var err error
	for {
		s.con, err = net.Dial("tcp", s.addr)
		if err != nil {
			golog.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}

	t := config.GetConfig().GetServiceType()
	d, err := aproto.MarshalLogicRegisterReq(t)
	if err != nil {
		golog.Critical(err)
	}

	err = anet.SendPacket(s.con, d)
	if err != nil {
		golog.Error(err)
		s.con.Close()
		go s.register()
		return
	}

	var buf [1024]byte
	l, _, err := anet.RecvPacket(s.con, buf[:])
	if err != nil {
		golog.Error(err)
		s.con.Close()
		go s.register()
		return
	}

	res, err := aproto.UnmarshalLogicRegisterRes(buf[aproto.HEAD_LEN:l])
	if err != nil {
		golog.Critical(err)
	}

	s.id = res.Id
	g_conmanager.addConnectSvr(int(s.id), s)
	golog.Info("register connect server:", s.addr, "scuccess!")

	go s.keepalive()

	s.ReadFromCon()

	g_conmanager.exitConnectSvr(int(s.id))
}

func (s *connectServer) keepalive() {
	golog.Info(s.addr, "connect server keepalive")
	err := anet.SendPacket(s.con, aproto.KeepAlive[:])
	if err != nil {
		golog.Error(err)
		s.con.Close()
		go s.register()
		return
	}

	logic_timer.NewTimer(keepalive_time, s.keepalive)
}

func (s *connectServer) ReadFromCon() {
	var buf [4096]byte
	for {
		l, _, err := anet.RecvPacket(s.con, buf[:])
		if err != nil {
			golog.Error(err)
			return
		}

		req, err := aproto.UnmarshalServiceReq(buf[aproto.HEAD_LEN:l])
		if err != nil {
			golog.Error(err)
			continue
		}

		go s.handlerPeerCmd(req)
	}
}

func (s *connectServer) SendToCon(d []byte) error {
	return anet.SendPacket(s.con, d)
}

func (s *connectServer) handlerPeerCmd(req *aproto.ServiceRequest) {
	pg, err := aproto.UnmarshalPeerPacket(req.Payload)
	if err != nil {
		golog.Error(err)
		return
	}

	switch pg.PacketType {
	case aproto.SENDMSGREQ_TYPE:
		g_singlechat.handlerSendMsgReq(s, req, pg.Data)
	case aproto.GETMSGREQ_TYPE:
		g_singlechat.handlerGetMsgReq(s, req, pg.Data)
	case aproto.RECVMSGRES_TYPE:
		g_singlechat.handlerRecvMsgRes(s, req, pg.Data)
	}
}
