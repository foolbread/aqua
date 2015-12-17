//@auther foolbread
//@time 2015-11-12
//@file aqua/connect_server/server/client_manager.go
package server

import (
	"fbcommon/golog"
	"sync"
)

type clientManager struct {
	lo      *sync.RWMutex
	mclient map[string]*Client
}

func newClientManager() *clientManager {
	r := new(clientManager)
	r.lo = new(sync.RWMutex)
	r.mclient = make(map[string]*Client) //key->token

	return r
}

func (s *clientManager) pushResponse(token string, d []byte) {
	s.lo.RLock()
	cli, ok := s.mclient[token]
	s.lo.RUnlock()
	if !ok {
		golog.Error("token:", token, "client already exit...")
		return
	}

	cli.sendResponse(d)
}

func (s *clientManager) exitClient(token string) {

}
