//@auther foolbread
//@time 2015-12-24
//@file aqua/common/storage/storage_handler.go
package storage

import (
	"strconv"
	"strings"

	"github.com/foolbread/fbcommon/golog"
)

const login_session_format = "_LOGIN_SESSION"
const message_id_format = "_MSG_ID"
const message_format = "_PEER_MSG"

type StorageHandler struct {
	handler *RedisHandler
}

func NewStorageHandler(info string) *StorageHandler {
	r := new(StorageHandler)

	var err error
	idx := strings.LastIndex(info, ":")
	addr := info[:idx]
	pwd := info[idx+1:]
	r.handler, err = NewRedisHandler(addr, pwd)
	if err != nil {
		golog.Error(err)
		return nil
	}

	return r
}

/*func NewStorageHandler(h *RedisHandler) *StorageHandler {
	r := new(StorageHandler)
	r.handler = h

	return r
}*/

func (s *StorageHandler) SetUsrSession(cid string, session string) error {
	key := cid + login_session_format

	return s.handler.SetKey(key, session)
}

func (s *StorageHandler) GetUsrSession(cid string) (string, error) {
	key := cid + login_session_format

	return s.handler.GetKey(key)
}

func (s *StorageHandler) DelUsrSession(cid string) error {
	key := cid + login_session_format

	return s.handler.DelKey(key)
}

func (s *StorageHandler) IsExistSession(cid string) (bool, error) {
	key := cid + login_session_format

	return s.handler.ExistsKey(key)
}

func (s *StorageHandler) IncreMsgId(cid string) (int, error) {
	key := cid + message_id_format

	return s.handler.IncreKey(key)
}

func (s *StorageHandler) AddPeerMsg(cid string, msg string, id int) error {
	key := cid + message_format

	return s.handler.SetHash(key, strconv.Itoa(id), msg)
}

func (s *StorageHandler) DelPeerMsg(cid string, ids []int64) error {
	key := cid + message_format
	var fileds []string
	for k, _ := range ids {
		fileds = append(fileds, strconv.Itoa(int(ids[k])))
	}

	return s.handler.DelHash(key, fileds)
}

func (s *StorageHandler) GetPeerMsgs(cid string) (map[string]string, error) {
	key := cid + message_format

	return s.handler.GetAllHash(key)
}

func (s *StorageHandler) GetPeerMsgsSize(cid string) (int, error) {
	key := cid + message_format

	return s.handler.GetHashLen(key)
}
