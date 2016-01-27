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

func (s *relationServer) handlerAddBlackReq(con *connectServer, r *aproto.ServiceRequest, d []byte) {
	req, err := aproto.UnmarshalAddBlackReq(d)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerAddBlackReq", "from:", req.From, "black:", req.Black)

	hnl := storage.GetStorage().GetRelationHandler(req.From)

	err = hnl.DelUsrFriend(req.From, req.Black)
	if err != nil {
		golog.Error(err)
	}

	err = hnl.AddUsrBlack(req.From, req.Black)
	if err != nil {
		golog.Error(err)
	}

	res, err := aproto.MarshalAddBlackRes(req.From, req.Black, aproto.STATUS_OK)
	if err != nil {
		golog.Error(err)
		return
	}

	SendMsg(con, req.From, r, res, aproto.ADDBLACKRES_TYPE)
}

func (s *relationServer) handlerDelFriendReq(con *connectServer, r *aproto.ServiceRequest, d []byte) {
	req, err := aproto.UnmarshalDelFriendReq(d)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerDelFriendReq", "from:", req.From, "friend:", req.Friend)

	h1 := storage.GetStorage().GetRelationHandler(req.From)
	h2 := storage.GetStorage().GetRelationHandler(req.Friend)

	err = h1.DelUsrFriend(req.From, req.Friend)
	if err != nil {
		golog.Error(err)
	}

	err = h2.DelUsrFriend(req.Friend, req.From)
	if err != nil {
		golog.Error(err)
	}

	res, err := aproto.MarshalDelFriendRes(req.From, req.Friend, aproto.STATUS_OK)
	if err != nil {
		golog.Error(err)
		return
	}

	SendMsg(con, req.From, r, res, aproto.DELFRIENDRES_TYPE)
}

func (s *relationServer) handlerAddFriendRes(con *connectServer, r *aproto.ServiceRequest, d []byte) {
	res, err := aproto.UnmarshalAddFriendRes(d)
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

	SendMsg(nil, res.From, r, d, aproto.ADDFRIENDRES_TYPE)
}

func (s *relationServer) handlerAddFriendReq(con *connectServer, r *aproto.ServiceRequest, d []byte) {
	req, err := aproto.UnmarshalAddFriendReq(d)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerAddFriendReq", "from:", req.From, "friend:", req.Friend)

	hnl := storage.GetStorage().GetRelationHandler(req.From)

	exist, err := hnl.IsExistFriend(req.From, req.Friend)
	if err != nil {
		golog.Error(err)
		return
	}

	if exist {
		return
	}

	hn := storage.GetStorage().GetSessionHandler(req.Friend)

	//判断对方是否在线
	online, err := hn.IsExistSession(req.Friend)
	if err != nil {
		golog.Error(err)
		return
	}

	if !online {
		l, err := hnl.GetRelationMsgsSize(req.Friend)
		if err != nil {
			golog.Error(err)
			return
		}

		//对方的关系消息队列已满
		if l > 50 {
			res, err := aproto.MarshalAddFriendRes(0, req.From, req.Friend, aproto.MESSAGE_FULL)
			if err != nil {
				golog.Error(err)
				return
			}

			SendMsg(con, req.From, r, res, aproto.ADDFRIENDRES_TYPE)
			return
		}
	}

	//获取消息ID
	id, err := hn.IncreMsgId(req.Friend)
	if err != nil {
		golog.Error(err)
		return
	}

	req.Id = int64(id)

	msg, err := req.Marshal()
	if err != nil {
		golog.Error(err)
		return
	}

	//添加消息到消息队列
	hn.AddRelationMsg(req.Friend, base64.StdEncoding.EncodeToString(msg), id)

	if online {
		//直接发送交友申请
		SendMsg(nil, req.Friend, r, d, aproto.ADDFRIENDREQ_TYPE)
	}

	//给予等待回复
	res, err := aproto.MarshalAddFriendRes(req.Id, req.From, req.Friend, aproto.WAITTING_RESPONSE)
	if err != nil {
		golog.Error(err)
		return
	}

	SendMsg(con, req.From, r, res, aproto.ADDFRIENDRES_TYPE)
}
