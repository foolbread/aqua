//@auther foolbread
//@time 2015-12-29
//@file aqua/login_server/server/connect_manager.go
package server

import (
	"container/list"
	"sync"

	"github.com/foolbread/fbcommon/golog"
)

type connectManager struct {
	lo          *sync.RWMutex
	connect_svr *list.List
}

func newConnectManager() *connectManager {
	r := new(connectManager)
	r.connect_svr = list.New()
	r.lo = new(sync.RWMutex)

	return r
}

func (s *connectManager) getConnectSvr(token []byte) *connectServer {
	var ret *connectServer
	s.lo.RLock()
	idx := int(token[0]) % s.connect_svr.Len()
	i := 0
	for e := s.connect_svr.Front(); e != nil; e = e.Next() {
		if i == idx {
			ret = e.Value.(*connectServer)
			break
		}
		i++
	}
	s.lo.RUnlock()

	return ret
}

func (s *connectManager) addConnectSvr(svr *connectServer) {
	s.lo.Lock()
	s.connect_svr.PushBack(svr)
	s.lo.Unlock()
}

func (s *connectManager) exitConnectSvr(Id uint32) {
	s.lo.Lock()
	for e := s.connect_svr.Front(); e != nil; e = e.Next() {
		v := e.Value.(*connectServer)
		if v.Id == Id {
			golog.Info("Id:", Id, "addr:", v.Addr, "connect server exit...")
			s.connect_svr.Remove(e)
			break
		}
	}
	s.lo.Unlock()

}
