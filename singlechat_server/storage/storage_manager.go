//@auther foolbread
//@time 2015-12-22
//@file aqua/singlechat_server/storage/storage_manager.go
package storage

import (
	astorage "aqua/common/storage"
	"aqua/singlechat_server/config"
	"crypto/md5"
	"strings"

	"github.com/foolbread/fbcommon/golog"
)

const default_count = 5

func InitStorageManager() {
	golog.Info("initing login storage manager...")
	g_storage = newStorageManager()

	infos := config.GetConfig().GetSessionDBInfos()
	g_storage.initSessionStorages(infos)

	infos = config.GetConfig().GetSinglechatDBInfos()
	g_storage.initSinglechatStorages(infos)
}

func GetStorage() *storageManager {
	return g_storage
}

var g_storage *storageManager

type storageManager struct {
	session_storages    [][]*astorage.StorageHandler
	singlechat_storages [][]*astorage.StorageHandler
}

func newStorageManager() *storageManager {
	ret := new(storageManager)
	ret.session_storages = make([][]*astorage.StorageHandler, len(config.GetConfig().GetSessionDBInfos()))

	return ret
}

func (s *storageManager) initSessionStorages(infos []string) {
	for k, v := range infos {
		idx := strings.LastIndex(v, ":")
		addr := v[:idx]
		pwd := v[idx+1:]
		golog.Info("session_storage", "addr:", addr, "pwd:", pwd)
		for i := 0; i < default_count; i++ {
			handler, err := astorage.NewRedisHandler(addr, pwd)
			if err != nil {
				golog.Critical(err)
			}
			g_storage.session_storages[k] = append(g_storage.session_storages[k], astorage.NewStorageHandler(handler))
		}
	}
}

func (s *storageManager) initSinglechatStorages(infos []string) {
	for k, v := range infos {
		idx := strings.LastIndex(v, ":")
		addr := v[:idx]
		pwd := v[idx+1:]
		golog.Info("singlechat_storage", "addr:", addr, "pwd:", pwd)
		for i := 0; i < default_count; i++ {
			handler, err := astorage.NewRedisHandler(addr, pwd)
			if err != nil {
				golog.Critical(err)
			}
			g_storage.singlechat_storages[k] = append(g_storage.singlechat_storages[k], astorage.NewStorageHandler(handler))
		}
	}
}

func (s *storageManager) GetSessionHandler(cid string) *astorage.StorageHandler {
	by := s.md5Byte(cid)
	as := s.session_storages[int(by)%len(s.session_storages)]
	return as[int(by)%len(as)]
}

func (s *storageManager) GetSingleHandler(cid string) *astorage.StorageHandler {
	by := s.md5Byte(cid)
	as := s.singlechat_storages[int(by)%len(s.singlechat_storages)]
	return as[int(by)%len(as)]
}

func (s *storageManager) md5Byte(str string) byte {
	k := md5.Sum([]byte(str))
	return k[0]
}
