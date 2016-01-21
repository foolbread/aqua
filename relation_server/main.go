//@auther foolbread
//@time 2016-01-21
//@file aqua/relation_server/main.go
package main

import (
	"aqua/relation_server/config"
	"aqua/relation_server/storage"
)

func init() {
	config.InitConfig()
	storage.InitStorageManager()
}

func main() {

}
