//@auther foolbread
//@time 2015-11-12
//@file aqua/connect_server/server/server.go
package server

import (
	aproto "aqua/common/proto"
	"aqua/connect_server/config"
	"fbcommon/golog"
	"time"
)

const ARRARY_LEN = 255

var default_timeout time.Duration = 2 * time.Minute

func InitServer() {
	golog.Info("initing connect server......")
	g_conserver = new(connectServer)
	g_conserver.id = uint32(config.GetConfig().GetConnectId())
	for i := 0; i < ARRARY_LEN; i++ {
		g_conserver.clients[i] = newClientManager()
	}
	g_conserver.startListen()
	g_conserver.startRegister()

	g_logicmanager = new(logicManager)
	g_logicmanager.startListen()

	keepalive.cmd = aproto.KEEPALIVE_CMD
	go connect_timer.Start()
}

var g_conserver *connectServer
var g_logicmanager *logicManager
