//@auther foolbread
//@time 2015-12-22
//@file aqua/login_server/storage/storage_manager.go
package storage

import (
	astorage "aqua/common/storage"
	"aqua/login_server/config"
	"crypto/md5"
	"fbcommon/golog"
	"strings"
)

const default_count = 5

func InitStorageManager() {
	golog.Info("initing login storage manager...")
	g_storage = newStorageManager()

	info := config.GetConfig().GetDBInfos()
	for k, v := range info {
		idx := strings.LastIndex(v, ":")
		addr := v[:idx]
		pwd := v[idx+1:]
		golog.Info("addr:", addr, "pwd:", pwd)
		for i := 0; i < default_count; i++ {
			handler, err := astorage.NewRedisHandler(addr, pwd)
			if err != nil {
				golog.Critical(err)
			}
			g_storage.storage_handlers[k] = append(g_storage.storage_handlers[k], astorage.NewStorageHandler(handler))
		}
	}

}

func GetStorage() *storageManager {
	return g_storage
}

var g_storage *storageManager

type storageManager struct {
	storage_handlers [][]*astorage.StorageHandler
}

func newStorageManager() *storageManager {
	ret := new(storageManager)
	ret.storage_handlers = make([][]*astorage.StorageHandler, len(config.GetConfig().GetDBInfos()))

	return ret
}

func (s *storageManager) GetStorageHandler(cid string) *astorage.StorageHandler {
	k := md5.Sum([]byte(cid))
	as := s.storage_handlers[int(k[0])%len(s.storage_handlers)]
	return as[int(k[0])%len(as)]
}
