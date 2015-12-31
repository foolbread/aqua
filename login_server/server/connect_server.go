//@auther foolbread
//@time 2015-12-29
//@file aqua/login_server/server/connect_server.go
package server

import (
	aerr "aqua/common/error"
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	fbatomic "fbcommon/atomic"
	"fbcommon/golog"
	"net"
	"time"
)

var check_time time.Duration = 15 * time.Second
var recv_time time.Duration = 30 * time.Second

type connectServer struct {
	Id    uint32
	Addr  string
	Alive fbatomic.AtomicBool
	Con   net.Conn
}

func newConnectServer(id uint32, addr string, con net.Conn) *connectServer {
	r := new(connectServer)
	r.Id = id
	r.Addr = addr
	r.Con = con
	r.Alive.Set(true)

	return r
}

func (s *connectServer) run() {
	var buf [1024]byte
	var cmd uint32
	var err error
	login_timer.NewTimer(check_time, s.checkAlive)

	for {
		_, cmd, err = anet.RecvPacket(s.Con, buf[:])
		if err != nil {
			golog.Error(err)
			break
		}

		if cmd != aproto.KEEPALIVE_CMD {
			golog.Error(aerr.ErrKeepalive)
			break
		}

		s.Alive.Set(true)
	}
}

func (s *connectServer) checkAlive() {
	if !s.Alive.Get() {
		s.Con.Close()
		return
	}

	s.Alive.Set(false)
	login_timer.NewTimer(check_time, s.checkAlive)
}
