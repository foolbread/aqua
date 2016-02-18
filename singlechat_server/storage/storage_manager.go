//@auther foolbread
//@time 2015-12-22
//@file aqua/singlechat_server/storage/storage_manager.go
package storage

import (
	astorage "aqua/common/storage"
	"aqua/singlechat_server/config"

	"github.com/foolbread/fbcommon/golog"
)

const default_count = 5

func InitStorageManager() {
	golog.Info("initing singlechat storage manager...")
	g_storage = newStorageManager()

	infos := config.GetConfig().GetSessionDBInfos()
	for k, v := range infos {
		for i := 0; i < default_count; i++ {
			hnl := astorage.NewSessionHandler(v)
			if hnl != nil {
				g_storage.session_storages[k] = append(g_storage.session_storages[k], hnl)
			}
		}
	}

	infos = config.GetConfig().GetSinglechatDBInfos()
	for k, v := range infos {
		for i := 0; i < default_count; i++ {
			hnl := astorage.NewSingleChatHandler(v)
			if hnl != nil {
				g_storage.singlechat_storages[k] = append(g_storage.singlechat_storages[k], hnl)
			}
		}
	}
}

func GetStorage() *storageManager {
	return g_storage
}

var g_storage *storageManager

type storageManager struct {
	session_storages    [][]*astorage.SessionHandler
	singlechat_storages [][]*astorage.SingleChatHandler
}

func newStorageManager() *storageManager {
	ret := new(storageManager)
	ret.session_storages = make([][]*astorage.SessionHandler, len(config.GetConfig().GetSessionDBInfos()))
	ret.singlechat_storages = make([][]*astorage.SingleChatHandler, len(config.GetConfig().GetSinglechatDBInfos()))

	return ret
}

func (s *storageManager) GetSessionHandler(cid string) *astorage.SessionHandler {
	by := astorage.Md5ToByte(cid)
	as := s.session_storages[int(by)%len(s.session_storages)]
	return as[int(by)%len(as)]
}

func (s *storageManager) GetSingleHandler(cid string) *astorage.SingleChatHandler {
	by := astorage.Md5ToByte(cid)
	as := s.singlechat_storages[int(by)%len(s.singlechat_storages)]
	return as[int(by)%len(as)]
}
