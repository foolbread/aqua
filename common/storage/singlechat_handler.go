//@auther foolbread
//@time 2016-02-18
//@file aqua/common/storage/singlechat_handler.go
package storage

import (
	"strconv"
	"strings"

	"github.com/foolbread/fbcommon/golog"
)

const peer_message_format = "_PEER_MSG"
const peer_msg_id_format = "_PEER_MSG_ID"

type SingleChatHandler struct {
	handler *RedisHandler
}

func NewSingleChatHandler(info string) *SingleChatHandler {
	r := new(SingleChatHandler)

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

func (s *SingleChatHandler) AddPeerMsg(cid string, msg string, id int) error {
	key := cid + peer_message_format

	return s.handler.SetHash(key, strconv.Itoa(id), msg)
}

func (s *SingleChatHandler) DelPeerMsg(cid string, id int64) error {
	key := cid + peer_message_format

	return s.handler.DelHash(key, strconv.Itoa(int(id)))
}

func (s *SingleChatHandler) DelPeerMsgs(cid string, ids []int64) error {
	key := cid + peer_message_format
	var fileds []string
	for k, _ := range ids {
		fileds = append(fileds, strconv.Itoa(int(ids[k])))
	}

	return s.handler.DelHashs(key, fileds)
}

func (s *SingleChatHandler) GetPeerMsgs(cid string) (map[string]string, error) {
	key := cid + peer_message_format

	return s.handler.GetAllHash(key)
}

func (s *SingleChatHandler) GetPeerMsgsSize(cid string) (int, error) {
	key := cid + peer_message_format

	return s.handler.GetHashLen(key)
}

///////////////////////////////////////////////////////////////////////
func (s *SingleChatHandler) IncrePeerMsgId(cid string) (int, error) {
	key := cid + peer_msg_id_format

	return s.handler.IncreKey(key)
}
