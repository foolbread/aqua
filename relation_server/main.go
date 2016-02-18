//@auther foolbread
//@time 2016-01-21
//@file aqua/relation_server/main.go
package main

import (
	"aqua/relation_server/config"
	"aqua/relation_server/server"
	"aqua/relation_server/storage"
	"runtime"
)

func init() {
	config.InitConfig()
	storage.InitStorageManager()
	server.InitServer()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	select {}
}
