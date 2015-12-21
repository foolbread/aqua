//@auther foolbread
//@time 2015-11-12
//@file aqua/login_server/server/server.go
package server

import (
	"fbcommon/golog"
)

func InitServer() {
	golog.Info("initing login server......")
	g_server = new(loginServer)
	g_server.startListen()
}

var g_server *loginServer
