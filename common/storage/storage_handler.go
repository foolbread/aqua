//@auther foolbread
//@time 2015-12-24
//@file aqua/common/storage/storage_handler.go
package storage

import ()

const login_token_format = "_login_token"

type StorageHandler struct {
	handler *RedisHandler
}

func NewStorageHandler(h *RedisHandler) *StorageHandler {
	r := new(StorageHandler)
	r.handler = h

	return r
}

func (s *StorageHandler) SetUsrToken(cid string, token string) error {
	key := cid + login_token_format

	return s.handler.SetKey(key, token)
}

func (s *StorageHandler) GetUsrToken(cid string) (string, error) {
	key := cid + login_token_format

	return s.handler.GetKey(key)
}

func (s *StorageHandler) IsExistToken(cid string) (bool, error) {
	key := cid + login_token_format

	return s.handler.ExistsKey(key)
}
