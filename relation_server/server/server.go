//@auther foolbread
//@time 2016-01-21
//@file aqua/relation_server/server/server.go
package server

import (
	"time"

	"github.com/foolbread/fbcommon/golog"

	"aqua/relation_server/config"

	fbtime "github.com/foolbread/fbcommon/time"
)

func InitServer() {
	golog.Info("initing relation server......")
	g_conmanager = newConnectManager()
	g_relation = newRelationServer()

	csvr := config.GetConfig().GetConnetServer()
	for _, v := range csvr {
		newConnectServer(v)
	}

	go logic_timer.Start()
}

var logic_timer *fbtime.Timer = fbtime.New(1 * time.Second)

var g_conmanager *connectManager
var g_relation *relationServer
