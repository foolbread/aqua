//@auther foolbread
//@time 2016-01-04
//@file aqua/singlechat_server/server/singlechat_server.go
package server

import (
	aproto "aqua/common/proto"
	"aqua/singlechat_server/storage"
	"encoding/base64"

	"github.com/foolbread/fbcommon/golog"
)

const MESSAGES_MAX = 10

type singlechatServer struct {
}

func newSinglechatServer() *singlechatServer {
	r := new(singlechatServer)

	return r
}

func (s *singlechatServer) handlerRecvPMsgRes(con *connectServer, r *aproto.ServiceRequest, pp *aproto.PeerPacket) {
	req, err := aproto.UnmarshalRecvPMsg(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}
	golog.Info("handlerRecvMsgRes", "cid:", req.Cid, "msg_id:", req.Id)

	if len(req.Id) == 0 {
		return
	}

	hnl := storage.GetStorage().GetSingleHandler(req.Cid)

	if len(req.Id) > 1 {
		err = hnl.DelPeerMsgs(req.Cid, req.Id)
		if err != nil {
			golog.Error(err)
		}
	} else {
		err = hnl.DelPeerMsg(req.Cid, req.Id[0])
		if err != nil {
			golog.Error(err)
		}
	}

}

func (s *singlechatServer) handlerGetPMsgReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.PeerPacket) {
	req, err := aproto.UnmarshalGetPMsgReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}
	golog.Info("handlerGetMsgReq", "cid:", req.Cid)

	hnl := storage.GetStorage().GetSingleHandler(req.Cid)

	//获取用户离线消息
	msgs, err := hnl.GetPeerMsgs(req.Cid)
	if err != nil {
		golog.Error(err)
		return
	}

	//每10条消息一个包
	cnt := 0
	var ms [MESSAGES_MAX]*aproto.PeerPacket
	for _, v := range msgs {
		data, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			golog.Error(err)
			continue
		}

		rq, err := aproto.UnmarshalPeerPacket(data)
		if err != nil {
			golog.Error(err)
			continue
		}

		ms[cnt] = rq
		cnt++

		if cnt >= MESSAGES_MAX {
			cnt = 0
			res, err := aproto.MarshalGetPMsgRes(ms[:])
			if err != nil {
				golog.Error(err)
				continue
			}

			pr := aproto.MarshalPeerPacketEx(aproto.GETPMSGRES_TYPE, 0, res)

			SendMsg(con, req.Cid, r, pr, false)
		}
	} //end for

	if cnt > 0 {
		res, err := aproto.MarshalGetPMsgRes(ms[:cnt])
		if err != nil {
			golog.Error(err)
			return
		}

		pr := aproto.MarshalPeerPacketEx(aproto.GETPMSGRES_TYPE, 0, res)

		SendMsg(con, req.Cid, r, pr, false)
	}
}

func (s *singlechatServer) handlerSendPMsgReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.PeerPacket) {
	req, err := aproto.UnmarshalSendPMsgReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerSendMsgReq", "from:", req.Msg.From, "to:", req.Msg.To)

	//send to
	to_session := storage.GetStorage().GetSessionHandler(req.Msg.To)

	id, err := to_session.IncreMsgId(req.Msg.To)
	if err != nil {
		golog.Error(err)
		return
	}
	pp.Id = int64(id)

	SendMsg(nil, req.Msg.To, r, pp, true)

	//send from
	from_session := storage.GetStorage().GetSessionHandler(req.Msg.From)

	id, err = from_session.IncreMsgId(req.Msg.From)
	if err != nil {
		golog.Error(err)
		return
	}

	res, err := aproto.MarshalSendPMsgRes(req.Msg.From, aproto.STATUS_OK, req.Msg.Sn)
	if err != nil {
		golog.Error(err)
		return
	}

	pr := aproto.MarshalPeerPacketEx(aproto.SENDPMSGRES_TYPE, int64(id), res)

	SendMsg(con, req.Msg.From, r, pr, true)
}

func (s *singlechatServer) handlerSendPMsgRes(con *connectServer, r *aproto.ServiceRequest, pp *aproto.PeerPacket) {
	res, err := aproto.UnmarshalSendPMsgRes(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerSendPMsgRes", "cid:", res.Cid, "msg_id:", pp.Id)
}
