//@auther foolbread
//@time 2016-01-04
//@file aqua/singlechat_server/server/server.go
package server

import (
	"aqua/singlechat_server/config"
	"fbcommon/golog"
	fbtime "fbcommon/time"
	"time"
)

func InitServer() {
	golog.Info("initing singlechat server ......")
	g_singlechat = newSinglechatServer()
	g_conmanager = newConnectManager()

	csvrs := config.GetConfig().GetConnetServer()
	for _, v := range csvrs {
		newConnectServer(v)
	}

	go logic_timer.Start()
}

var logic_timer *fbtime.Timer = fbtime.New(1 * time.Second)
var default_timeout time.Duration = 2 * time.Minute

var g_singlechat *singlechatServer
var g_conmanager *connectManager
