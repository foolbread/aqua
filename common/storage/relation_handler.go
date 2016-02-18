//@auther foolbread
//@time 2016-02-18
//@file aqua/common/storage/relation_handler.go
package storage

import (
	"strconv"
	"strings"

	"github.com/foolbread/fbcommon/golog"
)

const relation_message_format = "_RELATION_MSG"
const friend_list_format = "_FRIEND_LIST"
const black_list_format = "_BLACK_LIST"
const relation_msg_id_format = "_RELATION_MSG_ID"

type RelationHandler struct {
	handler *RedisHandler
}

func NewRelationHandler(info string) *RelationHandler {
	r := new(RelationHandler)

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

func (s *RelationHandler) AddRelationMsg(cid string, msg string, id int) error {
	key := cid + relation_message_format

	return s.handler.SetHash(key, strconv.Itoa(id), msg)
}

func (s *RelationHandler) DelRelationMsg(cid string, id int64) error {
	key := cid + relation_message_format

	return s.handler.DelHash(key, strconv.Itoa(int(id)))
}

func (s *RelationHandler) DelRelationMsgs(cid string, ids []int64) error {
	key := cid + relation_message_format
	var fileds []string
	for _, v := range ids {
		fileds = append(fileds, strconv.Itoa(int(v)))
	}

	return s.handler.DelHashs(key, fileds)
}

func (s *RelationHandler) GetRelationMsgs(cid string) (map[string]string, error) {
	key := cid + relation_message_format

	return s.handler.GetAllHash(key)
}

func (s *RelationHandler) GetRelationMsgsSize(cid string) (int, error) {
	key := cid + relation_message_format

	return s.handler.GetHashLen(key)
}

func (s *RelationHandler) AddUsrBlack(cid string, black string) error {
	set := cid + black_list_format

	return s.handler.AddSet(set, black)
}

func (s *RelationHandler) DelUsrBlack(cid string, black string) error {
	set := cid + black_list_format

	return s.handler.DelSet(set, black)
}

func (s *RelationHandler) IsExistBlack(cid string, black string) (bool, error) {
	set := cid + black_list_format

	return s.handler.SismemberSet(set, black)
}

///////////////////////////////////////////////////////////////////////
func (s *RelationHandler) AddUsrFriend(cid string, friend string) error {
	set := cid + friend_list_format

	return s.handler.AddSet(set, friend)
}

func (s *RelationHandler) DelUsrFriend(cid string, friend string) error {
	set := cid + friend_list_format

	return s.handler.DelSet(set, friend)
}

func (s *RelationHandler) IsExistFriend(cid string, friend string) (bool, error) {
	set := cid + friend_list_format

	return s.handler.SismemberSet(set, friend)
}

///////////////////////////////////////////////////////////////////////
func (s *RelationHandler) IncreRelationMsgId(cid string) (int, error) {
	key := cid + relation_msg_id_format

	return s.handler.IncreKey(key)
}
