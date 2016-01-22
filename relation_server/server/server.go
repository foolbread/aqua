//@auther foolbread
//@time 2016-01-21
//@file aqua/relation_server/server/server.go
package server

import (
	"time"

	"github.com/foolbread/fbcommon/golog"

	fbtime "github.com/foolbread/fbcommon/time"
)

func InitServer() {
	golog.Info("initing relation server......")
	g_conmanager = newConnectManager()
	g_relation = newRelationServer()

	go logic_timer.Start()
}

var g_conmanager *connectManager
var g_relation *relationServer
var logic_timer *fbtime.Timer = fbtime.New(1 * time.Second)
