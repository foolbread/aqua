//@auther foolbread
//@time 2015-11-12
//@file aqua/connect_server/server/server.go
package server

import (
	"fbcommon/golog"
)

const ARRARY_LEN = 255

func InitServer() {
	golog.Info("initing connect server......")
	g_conserver = new(connectServer)
	for i := 0; i < ARRARY_LEN; i++ {
		g_conserver.cons[i] = newClientManager()
	}
	g_conserver.startListen()

	g_logicmanager = new(logicManager)
	g_logicmanager.startListen()
}

var g_conserver *connectServer
var g_logicmanager *logicManager
