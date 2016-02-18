//@auther foolbread
//@time 2016-02-18
//@file aqua/common/storage/session_handler.go
package storage

import (
	"strings"

	"github.com/foolbread/fbcommon/golog"
)

const login_session_format = "_LOGIN_SESSION"
const peer_msg_id_format = "_PEER_MSG_ID"
const relation_msg_id_format = "_RELATION_MSG_ID"

type SessionHandler struct {
	handler *RedisHandler
}

func NewSessionHandler(info string) *SessionHandler {
	r := new(SessionHandler)

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

///////////////////////////////////////////////////////////////////////
func (s *SessionHandler) SetUsrSession(cid string, session string) error {
	key := cid + login_session_format

	return s.handler.SetKey(key, session)
}

func (s *SessionHandler) GetUsrSession(cid string) (string, error) {
	key := cid + login_session_format

	return s.handler.GetKey(key)
}

func (s *SessionHandler) DelUsrSession(cid string) error {
	key := cid + login_session_format

	return s.handler.DelKey(key)
}

func (s *SessionHandler) IsExistSession(cid string) (bool, error) {
	key := cid + login_session_format

	return s.handler.ExistsKey(key)
}

///////////////////////////////////////////////////////////////////////
func (s *SessionHandler) IncrePeerMsgId(cid string) (int, error) {
	key := cid + peer_msg_id_format

	return s.handler.IncreKey(key)
}

func (s *SessionHandler) IncreRelationMsgId(cid string) (int, error) {
	key := cid + relation_msg_id_format

	return s.handler.IncreKey(key)
}
