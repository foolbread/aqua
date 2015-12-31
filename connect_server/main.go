//@auther foolbread
//@time 2015-11-12
//@file aqua/connect_server/main.go
package main

import (
	"aqua/connect_server/config"
	"aqua/connect_server/server"
	"aqua/connect_server/storage"
)

func init() {
	config.InitConfig()
	server.InitServer()
	storage.InitStorageManager()
}

func main() {
	select {}
}
