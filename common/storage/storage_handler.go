//@auther foolbread
//@time 2015-12-24
//@file aqua/common/storage/storage_handler.go
package storage

const login_session_format = "_LOGIN_SESSION"

type StorageHandler struct {
	handler *RedisHandler
}

func NewStorageHandler(h *RedisHandler) *StorageHandler {
	r := new(StorageHandler)
	r.handler = h

	return r
}

func (s *StorageHandler) SetUsrSession(cid string, session string) error {
	key := cid + login_session_format

	return s.handler.SetKey(key, session)
}

func (s *StorageHandler) GetUsrSession(cid string) (string, error) {
	key := cid + login_session_format

	return s.handler.GetKey(key)
}

func (s *StorageHandler) IsExistSession(cid string) (bool, error) {
	key := cid + login_session_format

	return s.handler.ExistsKey(key)
}
