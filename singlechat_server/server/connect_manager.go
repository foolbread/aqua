//@auther foolbread
//@time 2016-01-04
//@file aqua/singlechat_server/server/connect_manager.go
package server

import (
	"fbcommon/golog"
	"sync"
)

type connectManager struct {
	con_map map[int]*connectServer
	lo      *sync.RWMutex
}

func newConnectManager() *connectManager {
	r := new(connectManager)
	r.con_map = make(map[int]*connectServer)
	r.lo = new(sync.RWMutex)

	return r
}

func (s *connectManager) getConnectSvr(id int) *connectServer {
	s.lo.RLock()
	r := s.con_map[id]
	s.lo.RUnlock()

	return r
}

func (s *connectManager) addConnectSvr(id int, c *connectServer) {
	s.lo.Lock()
	s.con_map[id] = c
	s.lo.Unlock()

	golog.Info("add connect server:", c.addr, "id:", id)
}

func (s *connectManager) exitConnectSvr(id int) {
	s.lo.Lock()
	delete(s.con_map, id)
	s.lo.Unlock()

	golog.Info("id:", id, "connect server exit!")
}
