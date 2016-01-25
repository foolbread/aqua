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
const peer_message_format = "_PEER_MSG"
const relation_message_format = "_SERVER_MSG"
const friend_list_format = "_FRIEND_LIST"
const black_list_format = "_BLACK_LIST"

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

func (s *StorageHandler) AddRelationMsg(cid string, msg string, id int) error {
	key := cid + relation_message_format

	return s.handler.SetHash(key, strconv.Itoa(id), msg)
}

func (s *StorageHandler) DelRelationMsg(cid string, ids []int64) error {
	key := cid + relation_message_format
	var fileds []string
	for _, v := range ids {
		fileds = append(fileds, strconv.Itoa(int(v)))
	}

	return s.handler.DelHash(key, fileds)
}

func (s *StorageHandler) GetRelationMsgs(cid string) (map[string]string, error) {
	key := cid + relation_message_format

	return s.handler.GetAllHash(key)
}

func (s *StorageHandler) GetRelationMsgsSize(cid string) (int, error) {
	key := cid + relation_message_format

	return s.handler.GetHashLen(key)
}

///////////////////////////////////////////////////////////////////////
func (s *StorageHandler) AddUsrBlack(cid string, black string) error {
	set := cid + black_list_format

	return s.handler.AddSet(set, black)
}

func (s *StorageHandler) DelUsrBlack(cid string, black string) error {
	set := cid + black_list_format

	return s.handler.DelSet(set, black)
}

func (s *StorageHandler) IsExistBlack(cid string, black string) (bool, error) {
	set := cid + black_list_format

	return s.handler.SismemberSet(set, black)
}

///////////////////////////////////////////////////////////////////////
func (s *StorageHandler) AddUsrFriend(cid string, friend string) error {
	set := cid + friend_list_format

	return s.handler.AddSet(set, friend)
}

func (s *StorageHandler) DelUsrFriend(cid string, friend string) error {
	set := cid + friend_list_format

	return s.handler.DelSet(set, friend)
}

func (s *StorageHandler) IsExistFriend(cid string, friend string) (bool, error) {
	set := cid + friend_list_format

	return s.handler.SismemberSet(set, friend)
}

///////////////////////////////////////////////////////////////////////
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

///////////////////////////////////////////////////////////////////////
func (s *StorageHandler) IncreMsgId(cid string) (int, error) {
	key := cid + message_id_format

	return s.handler.IncreKey(key)
}

///////////////////////////////////////////////////////////////////////
func (s *StorageHandler) AddPeerMsg(cid string, msg string, id int) error {
	key := cid + peer_message_format

	return s.handler.SetHash(key, strconv.Itoa(id), msg)
}

func (s *StorageHandler) DelPeerMsg(cid string, ids []int64) error {
	key := cid + peer_message_format
	var fileds []string
	for k, _ := range ids {
		fileds = append(fileds, strconv.Itoa(int(ids[k])))
	}

	return s.handler.DelHash(key, fileds)
}

func (s *StorageHandler) GetPeerMsgs(cid string) (map[string]string, error) {
	key := cid + peer_message_format

	return s.handler.GetAllHash(key)
}

func (s *StorageHandler) GetPeerMsgsSize(cid string) (int, error) {
	key := cid + peer_message_format

	return s.handler.GetHashLen(key)
}
