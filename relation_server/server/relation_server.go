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

func (s *relationServer) handlerAddFriendRes(con *connectServer, r *aproto.ServiceRequest, d []byte) {
	res, err := aproto.UnmarshalAddFriendRes(d)
	if err != nil {
		golog.Error(err)
		return
	}

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

	SendFriendRes(nil, res.From, r, d)
}

func (s *relationServer) handlerAddFriendReq(con *connectServer, r *aproto.ServiceRequest, d []byte) {
	req, err := aproto.UnmarshalAddFriendReq(d)
	if err != nil {
		golog.Error(err)
		return
	}

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

			SendFriendRes(con, req.From, r, res)
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
		SendFriendReq(nil, req.Friend, r, d)
	}

	//给予等待回复
	res, err := aproto.MarshalAddFriendRes(req.Id, req.From, req.Friend, aproto.WAITTING_RESPONSE)
	if err != nil {
		golog.Error(err)
		return
	}
	SendFriendRes(con, req.From, r, res)
}
