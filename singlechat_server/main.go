//@auther foolbread
//@time 2015-11-12
//@file aqua/singlechat_server/main.go
package main

import (
	"aqua/singlechat_server/config"
	"aqua/singlechat_server/server"
	"aqua/singlechat_server/storage"
	"runtime"
)

func init() {
	config.InitConfig()
	server.InitServer()
	storage.InitStorageManager()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	select {}
}
