//@auther foolbread
//@time 2015-11-12
//@file aqua/connect_server/server/logic_manager.go
package server

import (
	anet "aqua/common/net"
	aproto "aqua/common/proto"
	"aqua/connect_server/config"
	"fbcommon/golog"
	"net"
	"sync"
)

type logicManager struct {
	mserver map[int][]*logicServer
	lo      *sync.RWMutex
}

func newLogicManager() *logicManager {
	r := new(logicManager)
	r.mserver = make(map[int][]*logicServer)
	r.lo = new(sync.RWMutex)

	return r
}

func (s *logicManager) startListen() {
	iaddr := config.GetConfig().GetInnerAddr()

	go s.server(iaddr)
}

func (s *logicManager) getServer(t int, key []byte) *logicServer {
	var ret *logicServer
	s.lo.RLock()
	v, ok := s.mserver[t]
	if ok {
		ret = v[int(key[0])%len(v)]
	}
	s.lo.RUnlock()

	return ret
}

func (s *logicManager) server(a string) {
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

func (s *logicManager) handlerLogicCon(c net.Conn) {
	ser, err := s.handlerLogicLogin(c)
	if err != nil {
		golog.Error(err)
		return
	}

	ser.run()

	s.exitServer(ser.Service_type, ser.Addr)
}

func (s *logicManager) handlerLogicLogin(c net.Conn) (*logicServer, error) {
	var buf [1024]byte
	//request
	n, _, err := anet.RecvPacketEx(c, buf[:], default_timeout)
	if err != nil {
		return nil, err
	}

	req, err := aproto.UnmarshalSvrRegisterReq(buf[:n])
	if err != nil {
		return nil, err
	}

	//response
	da, err := aproto.MarshalSvrRegisterRes(aproto.STATUS_OK)
	if err != nil {
		return nil, err
	}

	err = anet.SendPacket(c, da)
	if err != nil {
		return nil, err
	}

	//register
	l := newLogicServer(int(req.ServiceType), c)

	s.lo.Lock()
	val := s.mserver[int(req.ServiceType)]
	val = append(val, l)
	s.lo.Unlock()

	return l, nil
}

func (s *logicManager) handlerLogicRes(pa *logicPacket) {
	g_conserver.handlerClientRes(pa)
}

func (s *logicManager) exitServer(t int, a string) {
	s.lo.Lock()
	ser, ok := s.mserver[t]
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

				s.mserver[t] = ser
				golog.Info("service type:", t, "service addr:", a, "exit...")
				break
			}
		}
	}
	s.lo.Unlock()
}
