//@auther foolbread
//@time 2015-12-22
//@file aqua/login_server/storage/storage_manager.go
package storage

import (
	astorage "aqua/common/storage"
	"aqua/login_server/config"
	"crypto/md5"
	"strings"

	"github.com/foolbread/fbcommon/golog"
)

const default_count = 5

func InitStorageManager() {
	golog.Info("initing login storage manager...")
	g_storage = newStorageManager()

	info := config.GetConfig().GetSessionDBInfos()
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
			g_storage.session_storages[k] = append(g_storage.session_storages[k], astorage.NewStorageHandler(handler))
		}
	}

}

func GetStorage() *storageManager {
	return g_storage
}

var g_storage *storageManager

type storageManager struct {
	session_storages [][]*astorage.StorageHandler
}

func newStorageManager() *storageManager {
	ret := new(storageManager)
	ret.session_storages = make([][]*astorage.StorageHandler, len(config.GetConfig().GetSessionDBInfos()))

	return ret
}

func (s *storageManager) GetSessionHandler(cid string) *astorage.StorageHandler {
	k := md5.Sum([]byte(cid))
	as := s.session_storages[int(k[0])%len(s.session_storages)]
	return as[int(k[0])%len(as)]
}
