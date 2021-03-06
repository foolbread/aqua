//@auther foolbread
//@time 2016-01-21
//@file aqua/relation_server/storage/storage_manager.go
package storage

import (
	astorage "aqua/common/storage"
	"aqua/relation_server/config"

	"github.com/foolbread/fbcommon/golog"
)

const default_count = 5

func InitStorageManager() {
	golog.Info("initing relation storage manager...")
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

	infos = config.GetConfig().GetRelationDBInfos()
	for k, v := range infos {
		for i := 0; i < default_count; i++ {
			hnl := astorage.NewRelationHandler(v)
			if hnl != nil {
				g_storage.relation_storages[k] = append(g_storage.relation_storages[k], hnl)
			}
		}
	}
}

func GetStorage() *storageManager {
	return g_storage
}

var g_storage *storageManager

type storageManager struct {
	session_storages  [][]*astorage.SessionHandler
	relation_storages [][]*astorage.RelationHandler
}

func newStorageManager() *storageManager {
	ret := new(storageManager)
	ret.session_storages = make([][]*astorage.SessionHandler, len(config.GetConfig().GetSessionDBInfos()))
	ret.relation_storages = make([][]*astorage.RelationHandler, len(config.GetConfig().GetRelationDBInfos()))

	return ret
}

func (s *storageManager) GetSessionHandler(cid string) *astorage.SessionHandler {
	by := astorage.Md5ToByte(cid)
	as := s.session_storages[int(by)%len(s.session_storages)]
	return as[int(by)%len(as)]
}

func (s *storageManager) GetRelationHandler(cid string) *astorage.RelationHandler {
	by := astorage.Md5ToByte(cid)
	as := s.relation_storages[int(by)%len(s.relation_storages)]
	return as[int(by)%len(as)]
}
