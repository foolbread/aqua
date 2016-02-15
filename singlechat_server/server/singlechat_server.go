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

	hnl := storage.GetStorage().GetSingleHandler(req.Cid)

	err = hnl.DelPeerMsgs(req.Cid, req.Id)
	if err != nil {
		golog.Error(err)
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

			SendMsg(con, req.Cid, r, pr)
		}
	} //end for

	if cnt > 0 {
		res, err := aproto.MarshalGetPMsgRes(ms[:cnt])
		if err != nil {
			golog.Error(err)
			return
		}

		pr := aproto.MarshalPeerPacketEx(aproto.GETPMSGRES_TYPE, 0, res)

		SendMsg(con, req.Cid, r, pr)
	}
}

func (s *singlechatServer) handlerSendPMsgReq(con *connectServer, r *aproto.ServiceRequest, pp *aproto.PeerPacket) {
	req, err := aproto.UnmarshalSendPMsgReq(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerSendMsgReq", "from:", req.Msg.From, "to:", req.Msg.To)

	to_single := storage.GetStorage().GetSingleHandler(req.Msg.To)
	to_session := storage.GetStorage().GetSessionHandler(req.Msg.To)
	//判断接收方是否在线
	online, err := to_session.IsExistSession(req.Msg.To)
	if err != nil {
		golog.Error(err)
		return
	}

	if !online {
		l, err := to_single.GetPeerMsgsSize(req.Msg.To)
		if err != nil {
			golog.Error(err)
			return
		}

		if l > 10 {
			//用户消息队列已满，发送错误信息
			res, err := aproto.MarshalSendPMsgRes(req.Msg.From, aproto.MESSAGE_FULL, req.Msg.Sn)
			if err != nil {
				golog.Error(err)
				return
			}

			pr := aproto.MarshalPeerPacketEx(aproto.SENDPMSGRES_TYPE, 0, res)

			SendMsg(con, req.Msg.From, r, pr)
			return
		}
	}

	//获取消息ID
	id, err := to_session.IncreMsgId(req.Msg.To)
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

	//把消息添加到消息队列
	to_single.AddPeerMsg(req.Msg.To, base64.StdEncoding.EncodeToString(msg), id)

	if online {
		SendMsg(nil, req.Msg.To, r, pp)
	}

	res, err := aproto.MarshalSendPMsgRes(req.Msg.From, aproto.STATUS_OK, req.Msg.Sn)
	if err != nil {
		golog.Error(err)
		return
	}

	pr := aproto.MarshalPeerPacketEx(aproto.SENDPMSGRES_TYPE, int64(id), res)

	SendMsg(con, req.Msg.From, r, pr)
}

func (s *singlechatServer) handlerSendPMsgRes(con *connectServer, r *aproto.ServiceRequest, pp *aproto.PeerPacket) {
	res, err := aproto.UnmarshalSendPMsgRes(pp.Data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerSendPMsgRes", "cid:", res.Cid, "msg_id:", pp.Id)

	hnl := storage.GetStorage().GetSingleHandler(res.Cid)

	err = hnl.DelPeerMsg(res.Cid, pp.Id)
	if err != nil {
		golog.Error(err)
	}
}
