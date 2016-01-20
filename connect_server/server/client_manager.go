//@auther foolbread
//@time 2015-11-12
//@file aqua/connect_server/server/client_manager.go
package server

import (
	"fmt"
	"sync"

	"github.com/foolbread/fbcommon/golog"
)

type clientManager struct {
	lo      *sync.RWMutex
	mclient map[string]*Client
	ch      chan *logicPacket
}

func newClientManager() *clientManager {
	r := new(clientManager)
	r.lo = new(sync.RWMutex)
	r.mclient = make(map[string]*Client) //key->token
	r.ch = make(chan *logicPacket, 1024)
	go r.run()

	return r
}

func (s *clientManager) run() {
	var pa *logicPacket
	var key string
	for {
		pa = <-s.ch
		if pa == nil {
			continue
		}

		key = fmt.Sprintf(keyformat, pa.token)
		s.lo.RLock()
		cli, ok := s.mclient[key]
		s.lo.RUnlock()

		if !ok {
			golog.Error("sn:", pa.sn, "token:", key, "client already exit...")
			continue
		}

		cli.sendData(pa.data)
	}
}

func (s *clientManager) pushResponse(pa *logicPacket) {
	s.ch <- pa
}

func (s *clientManager) exitClient(token string) {
	s.lo.Lock()
	delete(s.mclient, token)
	s.lo.Unlock()

	golog.Info("token:", token, "client exit...")
}
