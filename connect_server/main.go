//@auther foolbread
//@time 2015-11-12
//@file aqua/connect_server/main.go
package main

import (
	"aqua/connect_server/config"
	"aqua/connect_server/server"
)

func init() {
	config.InitConfig()
	server.InitServer()
}

func main() {
	select {}
}
