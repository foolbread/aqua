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

type relationServer struct {
}

func newRelationServer() *relationServer {
	r := new(relationServer)

	return r
}

func (s *relationServer) handlerDelBlackReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	req, err := aproto.UnmarshalDelBlackReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerDelBlackReq", "from:", req.From, "black:", req.Black)

	hnl := storage.GetStorage().GetRelationHandler(req.From)

	var status int32 = aproto.STATUS_OK
	err = hnl.DelUsrBlack(req.From, req.Black)
	if err != nil {
		golog.Error(err)
		status = aproto.SERVICE_ERROR
	}

	res, err := aproto.MarshalDelBlackRes(req.From, req.Black, status)
	if err != nil {
		golog.Error(err)
		return
	}

	rp := aproto.MarshalRelationPacketEx(aproto.DELBLACKRES_TYPE, 0, res)

	SendMsg(con, req.From, r, rp)
}

func (s *relationServer) handlerAddBlackReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	req, err := aproto.UnmarshalAddBlackReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerAddBlackReq", "from:", req.From, "black:", req.Black)

	hnl := storage.GetStorage().GetRelationHandler(req.From)

	var status int32 = aproto.STATUS_OK
	err = hnl.DelUsrFriend(req.From, req.Black)
	if err != nil {
		golog.Error(err)
		status = aproto.SERVICE_ERROR
	}

	err = hnl.AddUsrBlack(req.From, req.Black)
	if err != nil {
		golog.Error(err)
		status = aproto.SERVICE_ERROR
	}

	res, err := aproto.MarshalAddBlackRes(req.From, req.Black, status)
	if err != nil {
		golog.Error(err)
		return
	}

	rp := aproto.MarshalRelationPacketEx(aproto.ADDBLACKRES_TYPE, 0, res)

	SendMsg(con, req.From, r, rp)
}

func (s *relationServer) handlerDelFriendReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	req, err := aproto.UnmarshalDelFriendReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerDelFriendReq", "from:", req.From, "friend:", req.Friend)

	h1 := storage.GetStorage().GetRelationHandler(req.From)
	h2 := storage.GetStorage().GetRelationHandler(req.Friend)

	var status int32 = aproto.STATUS_OK
	err = h1.DelUsrFriend(req.From, req.Friend)
	if err != nil {
		golog.Error(err)
		status = aproto.SERVICE_ERROR
	}

	err = h2.DelUsrFriend(req.Friend, req.From)
	if err != nil {
		golog.Error(err)
		status = aproto.SERVICE_ERROR
	}

	res, err := aproto.MarshalDelFriendRes(req.From, req.Friend, status)
	if err != nil {
		golog.Error(err)
		return
	}

	rp := aproto.MarshalRelationPacketEx(aproto.DELFRIENDRES_TYPE, 0, res)

	SendMsg(con, req.From, r, rp)
}

func (s *relationServer) handlerAddFriendRes(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	res, err := aproto.UnmarshalAddFriendRes(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerAddFriendRes", "from:", res.From, "friend:", res.Friend, "status:", res.Status)

	if res.Status == aproto.AGREE_REQUEST {
		h1 := storage.GetStorage().GetRelationHandler(res.From)
		h2 := storage.GetStorage().GetRelationHandler(res.Friend)

		err = h1.AddUsrFriend(res.From, res.Friend)
		if err != nil {
			golog.Error(err)
		}

		err = h2.AddUsrFriend(res.Friend, res.From)
		if err != nil {
			golog.Error(err)
		}
	}

	hn_session := storage.GetStorage().GetSessionHandler(res.From)
	online, err := hn_session.IsExistSession(res.From)
	if err != nil {
		golog.Error(err)
		return
	}

	rp := aproto.MarshalRelationPacketEx(aproto.ADDFRIENDRES_TYPE, 0, pp.Data)
	if online {
		SendMsg(nil, res.From, r, rp)
	} else {
		hn_relation := storage.GetStorage().GetRelationHandler(res.From)
		id, err := hn_session.IncreMsgId(res.From)
		if err != nil {
			golog.Error(err)
			return
		}

		rp.Id = int64(id)
		msg, err := rp.Marshal()
		if err != nil {
			golog.Error(err)
		}

		hn_relation.AddRelationMsg(res.From, base64.StdEncoding.EncodeToString(msg), id)
	}

}

func (s *relationServer) handlerAddFriendReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.RelationPacket) {
	req, err := aproto.UnmarshalAddFriendReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerAddFriendReq", "from:", req.From, "friend:", req.Friend)

	hnl_relation := storage.GetStorage().GetRelationHandler(req.From)

	exist, err := hnl_relation.IsExistFriend(req.From, req.Friend)
	if err != nil {
		golog.Error(err)
		return
	}

	if exist {
		return
	}

	hn_session := storage.GetStorage().GetSessionHandler(req.Friend)

	//判断对方是否在线
	online, err := hn_session.IsExistSession(req.Friend)
	if err != nil {
		golog.Error(err)
		return
	}

	if !online {
		l, err := hnl_relation.GetRelationMsgsSize(req.Friend)
		if err != nil {
			golog.Error(err)
			return
		}

		//对方的关系消息队列已满
		if l > 50 {
			res, err := aproto.MarshalAddFriendRes(req.From, req.Friend, aproto.MESSAGE_FULL)
			if err != nil {
				golog.Error(err)
				return
			}

			rp := aproto.MarshalRelationPacketEx(aproto.ADDFRIENDRES_TYPE, 0, res)

			SendMsg(con, req.From, r, rp)
			return
		}
	}

	//获取消息ID
	id, err := hn_session.IncreMsgId(req.Friend)
	if err != nil {
		golog.Error(err)
		return
	}

	pp.Id = int64(id)

	msg, err := pp.Marshal()
	if err != nil {
		golog.Error(err)
		return
	}

	//添加消息到消息队列
	hn_session.AddRelationMsg(req.Friend, base64.StdEncoding.EncodeToString(msg), id)

	if online {
		//直接发送交友申请
		SendMsg(nil, req.Friend, r, pp)
	}

	//给予等待回复
	res, err := aproto.MarshalAddFriendRes(req.From, req.Friend, aproto.WAITTING_RESPONSE)
	if err != nil {
		golog.Error(err)
		return
	}

	rp := aproto.MarshalRelationPacketEx(aproto.ADDFRIENDRES_TYPE, pp.Id, res)

	SendMsg(con, req.From, r, rp)
}
