//@auther foolbread
//@time 2016-01-22
//@file aqua/relation_server/server/relation_server.go
package server

import (
	aproto "aqua/common/proto"
	"aqua/relation_server/storage"
	"encoding/base64"

	"github.com/foolbread/fbcommon/golog"
)

const MESSAGES_MAX = 10

type relationServer struct {
}

func newRelationServer() *relationServer {
	r := new(relationServer)

	return r
}

func (s *relationServer) handlerRecvPMsg(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	req, err := aproto.UnmarshalRecvPMsg(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerRecvPMsg", "cid:", req.Cid, "msg_id:", req.Id)

	if len(req.Id) == 0 {
		return
	}

	hnl := storage.GetStorage().GetRelationHandler(req.Cid)

	if len(req.Id) == 1 {
		err = hnl.DelRelationMsg(req.Cid, req.Id[0])
		if err != nil {
			golog.Error(err)
			return
		}
	} else {
		err = hnl.DelRelationMsgs(req.Cid, req.Id)
		if err != nil {
			golog.Error(err)
			return
		}
	}

}

func (s *relationServer) handlerGetRMsgReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	req, err := aproto.UnmarshalGetRMsgReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerGetRMsgReq", "cid:", req.Cid)

	hnl := storage.GetStorage().GetRelationHandler(req.Cid)

	msgs, err := hnl.GetRelationMsgs(req.Cid)
	if err != nil {
		golog.Error(err)
		return
	}

	//每10条消息一个包
	cnt := 0
	var ms [MESSAGES_MAX]*aproto.RelationPacket
	for _, v := range msgs {
		data, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			golog.Error(err)
			continue
		}

		rp, err := aproto.UnmarshalRelationPacket(data)
		if err != nil {
			golog.Error(err)
			continue
		}

		ms[cnt] = rp
		cnt++

		if cnt >= MESSAGES_MAX {
			cnt = 0
			res, err := aproto.MarshalGetRMsgRes(ms[:])
			if err != nil {
				golog.Error(err)
				continue
			}

			rp := aproto.MarshalRelationPacketEx(aproto.GETRMSGRES_TYPE, 0, res)

			SendMsg(con, req.Cid, r, rp, false)
		}
	} //end for

	if cnt > 0 {
		res, err := aproto.MarshalGetRMsgRes(ms[:cnt])
		if err != nil {
			golog.Error(err)
			return
		}

		rp := aproto.MarshalRelationPacketEx(aproto.GETRMSGRES_TYPE, 0, res)

		SendMsg(con, req.Cid, r, rp, false)
	}
}

func (s *relationServer) handlerDelBlackReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	req, err := aproto.UnmarshalDelBlackReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerDelBlackReq", "from:", req.From, "black:", req.Black)

	from_relation := storage.GetStorage().GetRelationHandler(req.From)

	var status int32 = aproto.STATUS_OK
	err = from_relation.DelUsrBlack(req.From, req.Black)
	if err != nil {
		golog.Error(err)
		status = aproto.SERVICE_ERROR
	}

	res, err := aproto.MarshalDelBlackRes(req.From, req.Black, status)
	if err != nil {
		golog.Error(err)
		return
	}

	from_session := storage.GetStorage().GetSessionHandler(req.From)
	id, err := from_session.IncreMsgId(req.From)
	if err != nil {
		golog.Error(err)
		return
	}

	rp := aproto.MarshalRelationPacketEx(aproto.DELBLACKRES_TYPE, int64(id), res)

	SendMsg(con, req.From, r, rp, true)
}

func (s *relationServer) handlerAddBlackReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	req, err := aproto.UnmarshalAddBlackReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerAddBlackReq", "from:", req.From, "black:", req.Black)

	from_relation := storage.GetStorage().GetRelationHandler(req.From)

	var status int32 = aproto.STATUS_OK
	err = from_relation.DelUsrFriend(req.From, req.Black)
	if err != nil {
		golog.Error(err)
		status = aproto.SERVICE_ERROR
	}

	err = from_relation.AddUsrBlack(req.From, req.Black)
	if err != nil {
		golog.Error(err)
		status = aproto.SERVICE_ERROR
	}

	res, err := aproto.MarshalAddBlackRes(req.From, req.Black, status)
	if err != nil {
		golog.Error(err)
		return
	}

	from_session := storage.GetStorage().GetSessionHandler(req.From)
	id, err := from_session.IncreMsgId(req.From)
	if err != nil {
		golog.Error(err)
		return
	}

	rp := aproto.MarshalRelationPacketEx(aproto.ADDBLACKRES_TYPE, int64(id), res)

	SendMsg(con, req.From, r, rp, true)
}

func (s *relationServer) handlerDelFriendReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	req, err := aproto.UnmarshalDelFriendReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerDelFriendReq", "from:", req.From, "friend:", req.Friend)

	from_relation := storage.GetStorage().GetRelationHandler(req.From)
	friend_relation := storage.GetStorage().GetRelationHandler(req.Friend)

	var status int32 = aproto.STATUS_OK
	err = from_relation.DelUsrFriend(req.From, req.Friend)
	if err != nil {
		golog.Error(err)
		status = aproto.SERVICE_ERROR
	}

	err = friend_relation.DelUsrFriend(req.Friend, req.From)
	if err != nil {
		golog.Error(err)
		status = aproto.SERVICE_ERROR
	}

	res, err := aproto.MarshalDelFriendRes(req.From, req.Friend, status)
	if err != nil {
		golog.Error(err)
		return
	}

	from_session := storage.GetStorage().GetSessionHandler(req.From)
	id, err := from_session.IncreMsgId(req.From)
	if err != nil {
		golog.Error(err)
		return
	}

	rp := aproto.MarshalRelationPacketEx(aproto.DELFRIENDRES_TYPE, int64(id), res)

	SendMsg(con, req.From, r, rp, true)
}

func (s *relationServer) handlerAddFriendRes(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	res, err := aproto.UnmarshalAddFriendRes(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerAddFriendRes", "from:", res.From, "friend:", res.Friend, "status:", res.Status)

	from_relation := storage.GetStorage().GetRelationHandler(res.From)
	friend_relation := storage.GetStorage().GetRelationHandler(res.Friend)
	if res.Status == aproto.AGREE_REQUEST {
		err = from_relation.AddUsrFriend(res.From, res.Friend)
		if err != nil {
			golog.Error(err)
		}

		err = friend_relation.AddUsrFriend(res.Friend, res.From)
		if err != nil {
			golog.Error(err)
		}
	}

	from_session := storage.GetStorage().GetSessionHandler(res.From)
	id, err := from_session.IncreMsgId(res.From)
	if err != nil {
		golog.Error(err)
		return
	}

	rp := aproto.MarshalRelationPacketEx(aproto.ADDFRIENDRES_TYPE, int64(id), pp.Data)

	SendMsg(nil, res.From, r, rp, true)
}

func (s *relationServer) handlerAddFriendReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	req, err := aproto.UnmarshalAddFriendReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerAddFriendReq", "from:", req.From, "friend:", req.Friend)

	from_relation := storage.GetStorage().GetRelationHandler(req.From)
	exist, err := from_relation.IsExistFriend(req.From, req.Friend)
	if err != nil {
		golog.Error(err)
		return
	}

	if exist {
		return
	}

	friend_session := storage.GetStorage().GetSessionHandler(req.Friend)
	id, err := friend_session.IncreMsgId(req.Friend)
	if err != nil {
		golog.Error(err)
		return
	}
	pp.Id = int64(id)

	SendMsg(nil, req.Friend, r, pp, true)

	//给予等待回复
	res, err := aproto.MarshalAddFriendRes(req.From, req.Friend, aproto.WAITTING_RESPONSE)
	if err != nil {
		golog.Error(err)
		return
	}

	from_session := storage.GetStorage().GetSessionHandler(req.From)
	id, err = from_session.IncreMsgId(req.From)
	if err != nil {
		golog.Error(err)
		return
	}

	rp := aproto.MarshalRelationPacketEx(aproto.ADDFRIENDRES_TYPE, int64(id), res)

	SendMsg(con, req.From, r, rp, true)
}
