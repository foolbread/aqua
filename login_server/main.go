//@auther foolbread
//@time 2015-11-12
//@file aqua/login_server/main.go
package main

import (
	"aqua/login_server/config"
	"aqua/login_server/server"
)

func init() {
	config.InitConfig()
	server.InitServer()
}

func main() {
	select {}
}
