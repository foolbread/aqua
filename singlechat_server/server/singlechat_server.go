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

func (s *singlechatServer) handlerRecvMsgRes(con *connectServer, r *aproto.ServiceRequest, d []byte) {
	req, err := aproto.UnmarshalRecvMsgRes(d)
	if err != nil {
		golog.Error(err)
		return
	}
	golog.Info("handlerRecvMsgRes", "cid:", req.Cid)

	hnl := storage.GetStorage().GetSingleHandler(req.Cid)

	err = hnl.DelPeerMsg(req.Cid, req.Id)
	if err != nil {
		golog.Error(err)
	}
}

func (s *singlechatServer) handlerGetMsgReq(con *connectServer, r *aproto.ServiceRequest, d []byte) {
	req, err := aproto.UnmarshalGetPMsgReq(d)
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
	var ms [MESSAGES_MAX]*aproto.PeerMessage
	for _, v := range msgs {
		data, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			golog.Error(err)
			continue
		}

		msg, err := aproto.UnmarshalPeerMessage(data)
		if err != nil {
			golog.Error(err)
			continue
		}

		ms[cnt] = msg
		cnt++

		if cnt >= MESSAGES_MAX {
			cnt = 0
			dd, err := aproto.MarshalGetPMsgRes(ms[:])
			if err != nil {
				golog.Error(err)
				continue
			}
			data, err := aproto.MarshalPeerPacket(aproto.GETMSGRES_TYPE, dd)
			if err != nil {
				golog.Error(err)
				continue
			}
			SendServiceMsg(con, r, data)
		}
	} //end for

	if cnt > 0 {
		dd, err := aproto.MarshalGetPMsgRes(ms[:cnt])
		if err != nil {
			golog.Error(err)
			return
		}
		data, err := aproto.MarshalPeerPacket(aproto.GETMSGRES_TYPE, dd)
		if err != nil {
			golog.Error(err)
			return
		}
		SendServiceMsg(con, r, data)
	}
}

func (s *singlechatServer) handlerSendMsgReq(con *connectServer, r *aproto.ServiceRequest, d []byte) {
	req, err := aproto.UnmarshalSendPMsgReq(d)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("handlerSendMsgReq", "from:", req.Msg.From, "to:", req.Msg.To)

	hnl := storage.GetStorage().GetSingleHandler(req.Msg.To)
	//判断接收方是否在线
	exist, err := hnl.IsExistSession(req.Msg.To)
	if err != nil {
		golog.Error(err)
		return
	}

	if !exist {
		l, err := hnl.GetPeerMsgsSize(req.Msg.To)
		if err != nil {
			golog.Error(err)
			return
		}

		if l > 10 {
			//用户消息队列已满，发送错误信息
			d, err = aproto.MarshalSendPMsgRes(aproto.MESSAGE_FULL, req.Msg.Sn)
			if err != nil {
				golog.Error(err)
				return
			}
			p, err := aproto.MarshalPeerPacket(aproto.SENDMSGRES_TYPE, d)
			if err != nil {
				golog.Error(err)
			}

			SendServiceMsg(con, r, p)
			return
		}
	}

	//获取消息ID
	id, err := hnl.IncreMsgId(req.Msg.To)
	if err != nil {
		golog.Error(err)
		return
	}

	req.Msg.Id = int64(id)

	msg, err := req.Msg.Marshal()
	if err != nil {
		golog.Error(err)
		return
	}

	//把消息添加到消息队列
	hnl.AddPeerMsg(req.Msg.To, base64.StdEncoding.EncodeToString(msg), id)

	if exist {
		SendPushMsgToUsr(r, req)
	}

	d, err = aproto.MarshalSendPMsgRes(aproto.STATUS_OK, req.Msg.Sn)
	if err != nil {
		golog.Error(err)
		return
	}
	p, err := aproto.MarshalPeerPacket(aproto.SENDMSGRES_TYPE, d)
	if err != nil {
		golog.Error(err)
		return
	}

	SendServiceMsg(con, r, p)
}
