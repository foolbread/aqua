//@auther foolbread
//@time 2015-11-12
//@file aqua/login_server/server/server.go
package server

import (
	"fbcommon/golog"
	fbtime "fbcommon/time"
	"time"
)

func InitServer() {
	golog.Info("initing login server......")
	g_server = newLoginServer()
	g_conmanager = newConnectManager()

	go login_timer.Start()

	go g_server.startListen()
}

var g_server *loginServer
var g_conmanager *connectManager

var login_timer *fbtime.Timer = fbtime.New(1 * time.Second)
